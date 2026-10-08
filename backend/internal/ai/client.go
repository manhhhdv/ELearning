// Package ai là lớp Service AI của hệ thống: dựng prompt, gọi API của nhà cung
// cấp, kiểm tra và chuẩn hoá kết quả trước khi trả về cho tầng HTTP.
//
// Client giữ một danh sách "kênh" (nhà cung cấp + khoá + model). Mỗi lần gọi,
// Client xoay vòng qua các kênh theo thứ tự ưu tiên: kênh nào vừa bị rate limit
// hoặc lỗi tạm thời sẽ bị treo một lúc và lượt gọi chuyển sang kênh kế tiếp.
// Nhờ vậy hệ thống vẫn chạy khi một nhà cung cấp hết hạn mức.
//
// Khoá API chỉ nằm ở máy chủ; frontend không bao giờ gọi thẳng nhà cung cấp.
package ai

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ErrDisabled được trả về khi chưa cấu hình khoá API nào.
var ErrDisabled = errors.New("chưa cấu hình khoá API cho nhà cung cấp AI")

// cooldownAfterFailure là thời gian treo một kênh sau khi kênh đó lỗi tạm thời.
// Đủ dài để không đập liên tục vào nhà cung cấp đang quá tải, nhưng vẫn ngắn
// hơn một tiết học nên admin không phải can thiệp tay.
const cooldownAfterFailure = 60 * time.Second

// cooldownRateLimited dài hơn vì hạn mức thường tính theo phút.
const cooldownRateLimited = 5 * time.Minute

// Credential là một kênh gọi AI: nhà cung cấp, khoá API và model cụ thể.
type Credential struct {
	// Label do admin đặt để phân biệt nhiều khoá của cùng một nhà cung cấp.
	Label    string `json:"label"`
	Provider string `json:"provider"`
	APIKey   string `json:"apiKey"`
	Model    string `json:"model"`
	// BaseURL ghi đè endpoint mặc định (dùng cho proxy hoặc nhà cung cấp "khác").
	BaseURL string `json:"baseUrl"`
	// Disabled=true thì kênh bị bỏ qua hoàn toàn, không cần xoá khỏi danh sách.
	Disabled bool `json:"disabled"`
}

// Normalize điền giá trị mặc định và cắt khoảng trắng.
func (c Credential) Normalize() Credential {
	c.Label = strings.TrimSpace(c.Label)
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.BaseURL = strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")
	if c.Provider = strings.ToLower(strings.TrimSpace(c.Provider)); c.Provider == "" {
		c.Provider = DefaultProvider
	}
	if c.Model = strings.TrimSpace(c.Model); c.Model == "" {
		c.Model = DefaultModelFor(c.Provider)
	}
	return c
}

// Name là tên hiển thị trong log và thông báo lỗi (không bao giờ chứa khoá).
func (c Credential) Name() string {
	if c.Label != "" {
		return c.Label
	}
	return c.Provider + "/" + c.Model
}

// channel là một kênh đang chạy, kèm trạng thái treo sau lỗi.
type channel struct {
	cred Credential
	prov Provider
	// blockedUntil: trước mốc này kênh bị bỏ qua khi xoay vòng.
	blockedUntil time.Time
}

// Client gọi API của các nhà cung cấp qua REST.
// An toàn khi dùng đồng thời từ nhiều request.
type Client struct {
	mu       sync.Mutex
	channels []*channel
	// next là con trỏ xoay vòng, để tải rải đều thay vì luôn đập vào kênh đầu.
	next int
	http *http.Client
}

// NewClient dựng client cho một kênh duy nhất. Giữ lại chữ ký cũ cho các chỗ
// chỉ cần một nhà cung cấp (ví dụ nút "Kiểm tra kết nối").
func NewClient(provider, apiKey, model string) *Client {
	return NewRotatingClient([]Credential{{Provider: provider, APIKey: apiKey, Model: model}})
}

// NewRotatingClient dựng client xoay vòng qua nhiều kênh theo đúng thứ tự
// truyền vào. Kênh không có khoá API hoặc bị tắt sẽ bị loại ngay.
func NewRotatingClient(creds []Credential) *Client {
	c := &Client{http: &http.Client{Timeout: 120 * time.Second}}
	for _, raw := range creds {
		cred := raw.Normalize()
		if cred.APIKey == "" || cred.Disabled {
			continue
		}
		prov, ok := FindProvider(cred.Provider)
		if !ok {
			slog.Warn("bỏ qua kênh AI với nhà cung cấp không hỗ trợ", "provider", cred.Provider)
			continue
		}
		if cred.BaseURL == "" {
			cred.BaseURL = prov.Endpoint
		}
		if cred.BaseURL == "" {
			slog.Warn("bỏ qua kênh AI thiếu địa chỉ endpoint", "kênh", cred.Name())
			continue
		}
		c.channels = append(c.channels, &channel{cred: cred, prov: prov})
	}
	return c
}

// Enabled cho biết hệ thống đã sẵn sàng gọi AI chưa.
func (c *Client) Enabled() bool { return c != nil && len(c.channels) > 0 }

// Provider trả về mã nhà cung cấp của kênh đầu tiên, để hiển thị cho admin.
func (c *Client) Provider() string {
	if !c.Enabled() {
		return DefaultProvider
	}
	return c.channels[0].cred.Provider
}

// Model trả về tên model của kênh đầu tiên, để hiển thị ở màn hình cấu hình.
func (c *Client) Model() string {
	if !c.Enabled() {
		return ""
	}
	return c.channels[0].cred.Model
}

// Channels trả về số kênh đang cấu hình.
func (c *Client) Channels() int {
	if c == nil {
		return 0
	}
	return len(c.channels)
}

// ---------------------------------------------------------------------------
// Kiểu dữ liệu chung, độc lập với nhà cung cấp
// ---------------------------------------------------------------------------

// Part là một mẩu nội dung: hoặc văn bản, hoặc tệp nhúng (PDF) dạng base64.
type Part struct {
	Text       string    `json:"text,omitempty"`
	InlineData *InlineDa `json:"inline_data,omitempty"`
}

type InlineDa struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

// Content là một lượt trong hội thoại. Role chỉ nhận "user" hoặc "model".
type Content struct {
	Role  string `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}

// Request là tham số một lần sinh nội dung.
type Request struct {
	// Hướng dẫn hệ thống: vai trò và ràng buộc chung cho cả cuộc gọi.
	System string
	// Toàn bộ lượt hội thoại, lượt cuối là câu hỏi hiện tại.
	Contents []Content
	// Bật khi cần model trả về đúng một khối JSON (dùng cho sinh câu hỏi, giáo án).
	JSON bool
	// 0 = dùng mặc định của model. Sinh câu hỏi dùng giá trị thấp cho ổn định.
	Temperature float64
}

// UserText tạo nhanh một lượt hội thoại chỉ có văn bản.
func UserText(text string) Content {
	return Content{Role: "user", Parts: []Part{{Text: text}}}
}

// UserWithPDF tạo một lượt hội thoại gồm văn bản kèm tài liệu PDF đã mã hoá base64.
func UserWithPDF(text, base64PDF string) Content {
	return Content{Role: "user", Parts: []Part{
		{InlineData: &InlineDa{MimeType: "application/pdf", Data: base64PDF}},
		{Text: text},
	}}
}

// ---------------------------------------------------------------------------
// Xoay vòng
// ---------------------------------------------------------------------------

// retryableError đánh dấu lỗi nên chuyển sang kênh khác thay vì trả về ngay:
// hết hạn mức, quá tải, lỗi máy chủ hoặc lỗi mạng.
type retryableError struct {
	err error
	// rateLimited=true thì treo kênh lâu hơn.
	rateLimited bool
}

func (e *retryableError) Error() string { return e.err.Error() }
func (e *retryableError) Unwrap() error { return e.err }

// retryable bọc một lỗi thành lỗi nên thử kênh khác.
func retryable(rateLimited bool, format string, args ...any) error {
	return &retryableError{err: fmt.Errorf(format, args...), rateLimited: rateLimited}
}

// statusRetryable cho biết một mã HTTP có đáng thử kênh khác không.
func statusRetryable(code int) (retry bool, rateLimited bool) {
	switch {
	case code == http.StatusTooManyRequests:
		return true, true
	case code == http.StatusRequestTimeout, code == http.StatusConflict:
		return true, false
	case code >= 500:
		return true, false
	}
	return false, false
}

// pickChannels trả về thứ tự kênh sẽ thử cho lần gọi này: bắt đầu từ con trỏ
// xoay vòng, các kênh đang bị treo xếp xuống cuối (vẫn thử nếu không còn kênh
// nào khác — thà gọi một kênh đang treo còn hơn trả lỗi cho giáo viên).
func (c *Client) pickChannels() []*channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	n := len(c.channels)
	if n == 0 {
		return nil
	}
	start := c.next % n
	c.next = (c.next + 1) % n

	now := time.Now()
	ready := make([]*channel, 0, n)
	blocked := make([]*channel, 0, n)
	for i := 0; i < n; i++ {
		ch := c.channels[(start+i)%n]
		if ch.blockedUntil.After(now) {
			blocked = append(blocked, ch)
		} else {
			ready = append(ready, ch)
		}
	}
	return append(ready, blocked...)
}

// block treo một kênh sau khi nó lỗi tạm thời.
func (c *Client) block(ch *channel, rateLimited bool) {
	d := cooldownAfterFailure
	if rateLimited {
		d = cooldownRateLimited
	}
	c.mu.Lock()
	ch.blockedUntil = time.Now().Add(d)
	c.mu.Unlock()
}

// Generate gửi yêu cầu tới kênh đầu tiên dùng được và trả về phần văn bản.
// Gặp rate limit hay lỗi tạm thời thì tự chuyển sang kênh kế tiếp; chỉ khi mọi
// kênh đều hỏng mới trả lỗi về cho người dùng.
func (c *Client) Generate(ctx context.Context, req Request) (string, error) {
	if !c.Enabled() {
		return "", ErrDisabled
	}

	var lastErr error
	for _, ch := range c.pickChannels() {
		text, err := c.generateOn(ctx, ch, req)
		if err == nil {
			return text, nil
		}
		// Người dùng huỷ hoặc hết hạn context thì đổi kênh cũng vô ích.
		if ctx.Err() != nil {
			return "", err
		}
		var re *retryableError
		if !errors.As(err, &re) {
			// Lỗi do chính nội dung hoặc cấu hình sai: kênh khác cũng hỏng như vậy.
			return "", err
		}
		c.block(ch, re.rateLimited)
		lastErr = err
		slog.Warn("kênh AI lỗi tạm thời, chuyển sang kênh kế tiếp",
			"kênh", ch.cred.Name(), "lỗi", err)
	}
	return "", fmt.Errorf("mọi kênh AI đều không dùng được (lỗi cuối: %w)", lastErr)
}

// generateOn gọi đúng một kênh, chọn bộ chuyển đổi theo giao thức.
func (c *Client) generateOn(ctx context.Context, ch *channel, req Request) (string, error) {
	switch ch.prov.Kind {
	case KindOpenAI:
		return c.generateOpenAI(ctx, ch, req)
	case KindAnthropic:
		return c.generateAnthropic(ctx, ch, req)
	default:
		return c.generateGemini(ctx, ch, req)
	}
}

// doJSON gửi một yêu cầu POST JSON và trả về thân phản hồi cùng mã HTTP.
// Lỗi mạng luôn là lỗi nên thử kênh khác.
func (c *Client) doJSON(ctx context.Context, url string, headers map[string]string, payload []byte) ([]byte, int, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytesReader(payload))
	if err != nil {
		return nil, 0, fmt.Errorf("dựng yêu cầu AI: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, 0, retryable(false, "gọi API AI thất bại: %w", err)
	}
	defer resp.Body.Close()

	raw, err := readLimited(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, retryable(false, "đọc phản hồi AI: %w", err)
	}
	return raw, resp.StatusCode, nil
}
