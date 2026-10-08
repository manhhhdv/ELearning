package ai

import (
	"context"
	"fmt"
	"strings"
)

// Nhà cung cấp AI được hỗ trợ.
const (
	ProviderGemini     = "gemini"
	ProviderOpenAI     = "openai"
	ProviderAnthropic  = "anthropic"
	ProviderGroq       = "groq"
	ProviderOpenRouter = "openrouter"
	ProviderCustom     = "custom"
)

// Kiểu giao thức REST của nhà cung cấp. Nhiều dịch vụ dùng chung định dạng
// /chat/completions của OpenAI nên chỉ cần ba bộ chuyển đổi.
const (
	KindGemini    = "gemini"
	KindOpenAI    = "openai"
	KindAnthropic = "anthropic"
)

// DefaultProvider là nhà cung cấp dùng khi cấu hình chưa chỉ định.
const DefaultProvider = ProviderGemini

// Provider mô tả một nhà cung cấp AI cho màn hình cấu hình của admin.
type Provider struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	// Kind quyết định định dạng yêu cầu/phản hồi khi gọi API.
	Kind string `json:"kind"`
	// Endpoint mặc định; admin có thể ghi đè bằng baseUrl riêng.
	Endpoint string `json:"endpoint"`
	// Các model gợi ý trong ô chọn; admin vẫn gõ được tên model khác.
	Models []string `json:"models"`
	// Model dùng khi admin không chọn gì.
	DefaultModel string `json:"defaultModel"`
	// Nơi lấy khoá API, hiện ở phần hướng dẫn trên giao diện.
	APIKeyURL string `json:"apiKeyUrl"`
	// RequiresBaseURL: nhà cung cấp không có endpoint cố định, admin phải nhập.
	RequiresBaseURL bool `json:"requiresBaseUrl"`
}

// Providers là danh sách nhà cung cấp hiện được hỗ trợ.
var Providers = []Provider{{
	ID:       ProviderGemini,
	Label:    "Google Gemini",
	Kind:     KindGemini,
	Endpoint: "https://generativelanguage.googleapis.com/v1beta",
	Models: []string{
		"gemini-2.5-flash",
		"gemini-2.5-pro",
		"gemini-2.5-flash-lite",
		"gemini-2.0-flash",
	},
	DefaultModel: "gemini-2.5-flash",
	APIKeyURL:    "https://aistudio.google.com/apikey",
}, {
	ID:       ProviderOpenAI,
	Label:    "OpenAI",
	Kind:     KindOpenAI,
	Endpoint: "https://api.openai.com/v1",
	Models: []string{
		"gpt-4.1-mini",
		"gpt-4.1",
		"gpt-4o-mini",
		"gpt-4o",
	},
	DefaultModel: "gpt-4.1-mini",
	APIKeyURL:    "https://platform.openai.com/api-keys",
}, {
	ID:       ProviderAnthropic,
	Label:    "Anthropic Claude",
	Kind:     KindAnthropic,
	Endpoint: "https://api.anthropic.com/v1",
	Models: []string{
		"claude-sonnet-4-5",
		"claude-haiku-4-5",
		"claude-opus-4-1",
	},
	DefaultModel: "claude-sonnet-4-5",
	APIKeyURL:    "https://console.anthropic.com/settings/keys",
}, {
	ID:       ProviderGroq,
	Label:    "Groq",
	Kind:     KindOpenAI,
	Endpoint: "https://api.groq.com/openai/v1",
	Models: []string{
		"llama-3.3-70b-versatile",
		"llama-3.1-8b-instant",
	},
	DefaultModel: "llama-3.3-70b-versatile",
	APIKeyURL:    "https://console.groq.com/keys",
}, {
	ID:       ProviderOpenRouter,
	Label:    "OpenRouter",
	Kind:     KindOpenAI,
	Endpoint: "https://openrouter.ai/api/v1",
	Models: []string{
		"google/gemini-2.5-flash",
		"openai/gpt-4.1-mini",
		"anthropic/claude-sonnet-4.5",
	},
	DefaultModel: "google/gemini-2.5-flash",
	APIKeyURL:    "https://openrouter.ai/keys",
}, {
	ID:              ProviderCustom,
	Label:           "Khác (tương thích OpenAI)",
	Kind:            KindOpenAI,
	Models:          []string{},
	DefaultModel:    "",
	RequiresBaseURL: true,
}}

// FindProvider tra cứu nhà cung cấp theo mã; ok=false nếu không hỗ trợ.
func FindProvider(id string) (Provider, bool) {
	for _, p := range Providers {
		if p.ID == id {
			return p, true
		}
	}
	return Provider{}, false
}

// NormalizeProvider đưa mã nhà cung cấp về giá trị hợp lệ.
// Chuỗi rỗng trả về nhà cung cấp mặc định; mã lạ bị từ chối.
func NormalizeProvider(id string) (string, error) {
	id = strings.ToLower(strings.TrimSpace(id))
	if id == "" {
		return DefaultProvider, nil
	}
	if _, ok := FindProvider(id); !ok {
		return "", fmt.Errorf("nhà cung cấp AI %q không được hỗ trợ", id)
	}
	return id, nil
}

// DefaultModelFor trả về model mặc định của một nhà cung cấp.
func DefaultModelFor(provider string) string {
	if p, ok := FindProvider(provider); ok {
		return p.DefaultModel
	}
	return Providers[0].DefaultModel
}

// Ping gọi một yêu cầu rất ngắn để kiểm tra khoá API và tên model có dùng được
// không. Dùng cho nút "Kiểm tra kết nối" ở màn hình cấu hình của admin.
func (c *Client) Ping(ctx context.Context) error {
	if !c.Enabled() {
		return ErrDisabled
	}
	_, err := c.Generate(ctx, Request{
		System:   "Bạn là công cụ kiểm tra kết nối. Luôn trả lời đúng một từ.",
		Contents: []Content{UserText("Trả lời đúng một từ: OK")},
	})
	return err
}
