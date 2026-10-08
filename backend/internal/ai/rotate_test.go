package ai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// geminiOK trả về một phản hồi Gemini hợp lệ chứa text.
func geminiOK(text string) string {
	return `{"candidates":[{"content":{"parts":[{"text":"` + text + `"}]}}]}`
}

func TestGenerateChuyenKenhKhiRateLimit(t *testing.T) {
	var hits int32
	limited := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":{"code":429,"message":"Quota exceeded","status":"RESOURCE_EXHAUSTED"}}`))
	}))
	defer limited.Close()
	healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(geminiOK("xin chao")))
	}))
	defer healthy.Close()

	c := NewRotatingClient([]Credential{
		{Label: "chinh", Provider: ProviderGemini, APIKey: "k1", Model: "m", BaseURL: limited.URL},
		{Label: "du phong", Provider: ProviderGemini, APIKey: "k2", Model: "m", BaseURL: healthy.URL},
	})

	got, err := c.Generate(context.Background(), Request{Contents: []Content{UserText("hi")}})
	if err != nil || got != "xin chao" {
		t.Fatalf("phải rơi sang kênh dự phòng, nhận được %q (lỗi %v)", got, err)
	}
	// Kênh bị rate limit phải bị treo: lần gọi sau không đụng vào nó nữa.
	before := atomic.LoadInt32(&hits)
	if _, err := c.Generate(context.Background(), Request{Contents: []Content{UserText("hi")}}); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&hits) != before {
		t.Error("kênh đang bị treo vẫn được gọi lại")
	}
}

func TestGenerateKhongChuyenKenhVoiLoiCauHinh(t *testing.T) {
	var second int32
	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":400,"message":"API key not valid"}}`))
	}))
	defer bad.Close()
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&second, 1)
		w.Write([]byte(geminiOK("khong nen goi")))
	}))
	defer other.Close()

	c := NewRotatingClient([]Credential{
		{Provider: ProviderGemini, APIKey: "k1", Model: "m", BaseURL: bad.URL},
		{Provider: ProviderGemini, APIKey: "k2", Model: "m", BaseURL: other.URL},
	})
	if _, err := c.Generate(context.Background(), Request{Contents: []Content{UserText("hi")}}); err == nil {
		t.Fatal("khoá sai phải trả lỗi ngay")
	}
	if atomic.LoadInt32(&second) != 0 {
		t.Error("lỗi cấu hình không được kéo theo việc đốt kênh khác")
	}
}

func TestGenerateQuaKenhOpenAI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			t.Errorf("đường dẫn sai: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer k1" {
			t.Errorf("thiếu Bearer token, nhận được %q", got)
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"ok nhe"}}]}`))
	}))
	defer srv.Close()

	c := NewRotatingClient([]Credential{{Provider: ProviderOpenAI, APIKey: "k1", Model: "m", BaseURL: srv.URL}})
	got, err := c.Generate(context.Background(), Request{System: "s", Contents: []Content{UserText("hi")}})
	if err != nil || got != "ok nhe" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestGenerateQuaKenhAnthropic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "k1" || r.Header.Get("anthropic-version") == "" {
			t.Error("thiếu header xác thực của Anthropic")
		}
		w.Write([]byte(`{"content":[{"type":"text","text":"chao ban"}]}`))
	}))
	defer srv.Close()

	c := NewRotatingClient([]Credential{{Provider: ProviderAnthropic, APIKey: "k1", Model: "m", BaseURL: srv.URL}})
	got, err := c.Generate(context.Background(), Request{Contents: []Content{UserText("hi")}})
	if err != nil || got != "chao ban" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestRotatingClientBoQuaKenhKhongDungDuoc(t *testing.T) {
	c := NewRotatingClient([]Credential{
		{Provider: ProviderGemini, APIKey: "   ", Model: "m"},
		{Provider: ProviderGemini, APIKey: "k", Model: "m", Disabled: true},
		{Provider: "khong-ton-tai", APIKey: "k"},
		{Provider: ProviderCustom, APIKey: "k", Model: "m"}, // thiếu baseUrl
		{Provider: ProviderGemini, APIKey: "k"},
	})
	if c.Channels() != 1 {
		t.Fatalf("chỉ một kênh hợp lệ, nhận được %d", c.Channels())
	}
	if c.Model() != DefaultModelFor(ProviderGemini) {
		t.Errorf("model rỗng phải rơi về mặc định, nhận được %q", c.Model())
	}
}
