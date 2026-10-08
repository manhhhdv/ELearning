package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

func bytesReader(b []byte) *bytes.Reader { return bytes.NewReader(b) }

// readLimited đọc tối đa 16MB để một phản hồi hỏng không thổi bay bộ nhớ.
func readLimited(r io.Reader) ([]byte, error) { return io.ReadAll(io.LimitReader(r, 16<<20)) }

const maxOutputTokens = 65536

// quotaWords là các từ khoá cho biết lỗi thuộc về hạn mức, kể cả khi nhà cung
// cấp trả mã HTTP không phải 429.
var quotaWords = []string{"rate limit", "rate_limit", "quota", "overloaded", "too many requests", "resource_exhausted", "capacity"}

// looksRateLimited đoán lỗi hạn mức từ nội dung thông báo.
func looksRateLimited(msg string) bool {
	msg = strings.ToLower(msg)
	for _, w := range quotaWords {
		if strings.Contains(msg, w) {
			return true
		}
	}
	return false
}

// apiError dựng lỗi từ thông báo của nhà cung cấp: lỗi hạn mức/quá tải được
// đánh dấu để Generate chuyển sang kênh khác.
func apiError(name string, status int, message string) error {
	retry, rateLimited := statusRetryable(status)
	if looksRateLimited(message) {
		retry, rateLimited = true, true
	}
	if retry {
		return retryable(rateLimited, "%s báo lỗi: %s", name, message)
	}
	return fmt.Errorf("%s báo lỗi: %s", name, message)
}

// ---------------------------------------------------------------------------
// Gemini
// ---------------------------------------------------------------------------

type geminiRequest struct {
	SystemInstruction *Content          `json:"systemInstruction,omitempty"`
	Contents          []Content         `json:"contents"`
	GenerationConfig  *generationConfig `json:"generationConfig,omitempty"`
}

type generationConfig struct {
	ResponseMimeType string   `json:"responseMimeType,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	MaxOutputTokens  int      `json:"maxOutputTokens,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content      Content `json:"content"`
		FinishReason string  `json:"finishReason"`
	} `json:"candidates"`
	PromptFeedback struct {
		BlockReason string `json:"blockReason"`
	} `json:"promptFeedback"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (c *Client) generateGemini(ctx context.Context, ch *channel, req Request) (string, error) {
	body := geminiRequest{Contents: req.Contents}
	if req.System != "" {
		body.SystemInstruction = &Content{Parts: []Part{{Text: req.System}}}
	}
	cfg := generationConfig{MaxOutputTokens: maxOutputTokens}
	if req.JSON {
		cfg.ResponseMimeType = "application/json"
	}
	if req.Temperature > 0 {
		t := req.Temperature
		cfg.Temperature = &t
	}
	body.GenerationConfig = &cfg

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("dựng yêu cầu Gemini: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent", ch.cred.BaseURL, ch.cred.Model)
	// Khoá đi ở header thay vì query string để không lọt vào log truy cập.
	raw, status, err := c.doJSON(ctx, url, map[string]string{"x-goog-api-key": ch.cred.APIKey}, payload)
	if err != nil {
		return "", err
	}

	var parsed geminiResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", retryable(false, "phản hồi Gemini không phải JSON hợp lệ (HTTP %d)", status)
	}
	if parsed.Error != nil {
		slog.Error("Gemini API trả lỗi", "status", parsed.Error.Status, "message", parsed.Error.Message)
		return "", apiError("Gemini API", status, parsed.Error.Message)
	}
	if status != 200 {
		return "", apiError("Gemini API", status, fmt.Sprintf("HTTP %d", status))
	}
	if parsed.PromptFeedback.BlockReason != "" {
		return "", fmt.Errorf("nội dung bị Gemini từ chối xử lý (%s)", parsed.PromptFeedback.BlockReason)
	}
	if len(parsed.Candidates) == 0 {
		return "", retryable(false, "Gemini không trả về nội dung nào")
	}

	var sb strings.Builder
	for _, p := range parsed.Candidates[0].Content.Parts {
		sb.WriteString(p.Text)
	}
	return nonEmpty(sb.String(), "Gemini")
}

// ---------------------------------------------------------------------------
// OpenAI và các dịch vụ tương thích (Groq, OpenRouter, proxy nội bộ…)
// ---------------------------------------------------------------------------

type openAIMessage struct {
	Role string `json:"role"`
	// Content là chuỗi khi chỉ có văn bản, là mảng phần tử khi kèm tệp.
	Content any `json:"content"`
}

type openAIRequest struct {
	Model          string          `json:"model"`
	Messages       []openAIMessage `json:"messages"`
	Temperature    *float64        `json:"temperature,omitempty"`
	MaxTokens      int             `json:"max_completion_tokens,omitempty"`
	ResponseFormat map[string]any  `json:"response_format,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error"`
}

func (c *Client) generateOpenAI(ctx context.Context, ch *channel, req Request) (string, error) {
	msgs := make([]openAIMessage, 0, len(req.Contents)+1)
	if req.System != "" {
		msgs = append(msgs, openAIMessage{Role: "system", Content: req.System})
	}
	for _, content := range req.Contents {
		msgs = append(msgs, openAIMessage{Role: openAIRole(content.Role), Content: openAIContent(content.Parts)})
	}

	body := openAIRequest{Model: ch.cred.Model, Messages: msgs, MaxTokens: maxOutputTokens}
	if req.JSON {
		body.ResponseFormat = map[string]any{"type": "json_object"}
	}
	if req.Temperature > 0 {
		t := req.Temperature
		body.Temperature = &t
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("dựng yêu cầu %s: %w", ch.prov.Label, err)
	}

	raw, status, err := c.doJSON(ctx, ch.cred.BaseURL+"/chat/completions",
		map[string]string{"Authorization": "Bearer " + ch.cred.APIKey}, payload)
	if err != nil {
		return "", err
	}

	var parsed openAIResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", retryable(false, "phản hồi %s không phải JSON hợp lệ (HTTP %d)", ch.prov.Label, status)
	}
	if parsed.Error != nil {
		slog.Error("nhà cung cấp AI trả lỗi", "kênh", ch.cred.Name(), "message", parsed.Error.Message)
		return "", apiError(ch.prov.Label, status, parsed.Error.Message)
	}
	if status != 200 {
		return "", apiError(ch.prov.Label, status, fmt.Sprintf("HTTP %d", status))
	}
	if len(parsed.Choices) == 0 {
		return "", retryable(false, "%s không trả về nội dung nào", ch.prov.Label)
	}
	return nonEmpty(parsed.Choices[0].Message.Content, ch.prov.Label)
}

// openAIRole đổi role của Gemini ("model") sang role của OpenAI ("assistant").
func openAIRole(role string) string {
	if role == "model" || role == "assistant" {
		return "assistant"
	}
	return "user"
}

// openAIContent đổi các Part sang định dạng nội dung của OpenAI. Chỉ có văn bản
// thì trả về chuỗi cho gọn — nhiều dịch vụ tương thích không nhận dạng mảng.
func openAIContent(parts []Part) any {
	hasFile := false
	for _, p := range parts {
		if p.InlineData != nil {
			hasFile = true
			break
		}
	}
	if !hasFile {
		var sb strings.Builder
		for _, p := range parts {
			sb.WriteString(p.Text)
		}
		return sb.String()
	}

	out := make([]map[string]any, 0, len(parts))
	for _, p := range parts {
		if p.InlineData != nil {
			out = append(out, map[string]any{
				"type": "file",
				"file": map[string]any{
					"filename":  "tai-lieu.pdf",
					"file_data": "data:" + p.InlineData.MimeType + ";base64," + p.InlineData.Data,
				},
			})
			continue
		}
		if p.Text != "" {
			out = append(out, map[string]any{"type": "text", "text": p.Text})
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Anthropic
// ---------------------------------------------------------------------------

type anthropicRequest struct {
	Model       string             `json:"model"`
	System      string             `json:"system,omitempty"`
	Messages    []anthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature *float64           `json:"temperature,omitempty"`
}

type anthropicMessage struct {
	Role    string           `json:"role"`
	Content []map[string]any `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	StopReason string `json:"stop_reason"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) generateAnthropic(ctx context.Context, ch *channel, req Request) (string, error) {
	system := req.System
	if req.JSON {
		// Anthropic không có chế độ JSON riêng, phải nói rõ trong hướng dẫn.
		system = strings.TrimSpace(system + "\nChỉ trả về đúng một khối JSON hợp lệ, không kèm giải thích hay dấu ```.")
	}

	msgs := make([]anthropicMessage, 0, len(req.Contents))
	for _, content := range req.Contents {
		blocks := make([]map[string]any, 0, len(content.Parts))
		for _, p := range content.Parts {
			if p.InlineData != nil {
				blocks = append(blocks, map[string]any{
					"type": "document",
					"source": map[string]any{
						"type":       "base64",
						"media_type": p.InlineData.MimeType,
						"data":       p.InlineData.Data,
					},
				})
				continue
			}
			if p.Text != "" {
				blocks = append(blocks, map[string]any{"type": "text", "text": p.Text})
			}
		}
		if len(blocks) == 0 {
			continue
		}
		role := "user"
		if content.Role == "model" || content.Role == "assistant" {
			role = "assistant"
		}
		msgs = append(msgs, anthropicMessage{Role: role, Content: blocks})
	}

	// max_tokens của Anthropic là bắt buộc và giới hạn thấp hơn Gemini.
	body := anthropicRequest{Model: ch.cred.Model, System: system, Messages: msgs, MaxTokens: 16384}
	if req.Temperature > 0 {
		t := req.Temperature
		body.Temperature = &t
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("dựng yêu cầu Anthropic: %w", err)
	}

	raw, status, err := c.doJSON(ctx, ch.cred.BaseURL+"/messages", map[string]string{
		"x-api-key":         ch.cred.APIKey,
		"anthropic-version": "2023-06-01",
	}, payload)
	if err != nil {
		return "", err
	}

	var parsed anthropicResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", retryable(false, "phản hồi Anthropic không phải JSON hợp lệ (HTTP %d)", status)
	}
	if parsed.Error != nil {
		slog.Error("Anthropic API trả lỗi", "type", parsed.Error.Type, "message", parsed.Error.Message)
		return "", apiError("Anthropic API", status, parsed.Error.Message)
	}
	if status != 200 {
		return "", apiError("Anthropic API", status, fmt.Sprintf("HTTP %d", status))
	}

	var sb strings.Builder
	for _, b := range parsed.Content {
		if b.Type == "text" {
			sb.WriteString(b.Text)
		}
	}
	return nonEmpty(sb.String(), "Anthropic")
}

// nonEmpty từ chối phản hồi rỗng — thường là dấu hiệu model bị cắt giữa chừng,
// và lần gọi sang kênh khác có thể ra kết quả.
func nonEmpty(text, name string) (string, error) {
	if text = strings.TrimSpace(text); text == "" {
		return "", retryable(false, "%s trả về nội dung rỗng", name)
	}
	return text, nil
}
