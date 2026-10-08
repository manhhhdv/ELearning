package ai

import "testing"

func TestNormalizeProvider(t *testing.T) {
	if got, err := NormalizeProvider(""); err != nil || got != DefaultProvider {
		t.Errorf("chuỗi rỗng phải trả về nhà cung cấp mặc định, nhận được %q (lỗi %v)", got, err)
	}
	if got, err := NormalizeProvider("  GEMINI "); err != nil || got != ProviderGemini {
		t.Errorf("phải bỏ khoảng trắng và không phân biệt hoa thường, nhận được %q (lỗi %v)", got, err)
	}
	if _, err := NormalizeProvider("chatgpt"); err == nil {
		t.Error("nhà cung cấp chưa hỗ trợ phải bị từ chối")
	}
}

func TestDefaultModelFor(t *testing.T) {
	if got := DefaultModelFor(ProviderGemini); got == "" {
		t.Error("nhà cung cấp Gemini phải có model mặc định")
	}
	// Mã lạ không được làm hỏng cấu hình: vẫn trả về một model dùng được.
	if got := DefaultModelFor("khong-ton-tai"); got == "" {
		t.Error("mã lạ vẫn phải trả về một model mặc định")
	}
}

func TestNewClientFillsDefaults(t *testing.T) {
	c := NewClient("", "khoa-test", "")
	if c.Provider() != DefaultProvider {
		t.Errorf("provider = %q, mong đợi %q", c.Provider(), DefaultProvider)
	}
	if c.Model() != DefaultModelFor(DefaultProvider) {
		t.Errorf("model = %q, mong đợi model mặc định", c.Model())
	}
	if !c.Enabled() {
		t.Error("có khoá API thì client phải ở trạng thái bật")
	}
	if NewClient("", "   ", "").Enabled() {
		t.Error("khoá toàn khoảng trắng phải coi như chưa cấu hình")
	}
}
