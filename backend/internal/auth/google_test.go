package auth

import (
	"strings"
	"testing"
	"time"
)

func newTestAuthenticator() *GoogleAuthenticator {
	return NewGoogleAuthenticator("client-id", "client-secret",
		"http://localhost:8081/api/auth/google/callback", "khoa-bi-mat-de-test")
}

func TestStateHopLe(t *testing.T) {
	g := newTestAuthenticator()

	state, err := g.NewState()
	if err != nil {
		t.Fatalf("không sinh được state: %v", err)
	}
	if err := g.VerifyState(state); err != nil {
		t.Errorf("state vừa sinh ra lại bị từ chối: %v", err)
	}
}

func TestStateMoiLanMotKhac(t *testing.T) {
	g := newTestAuthenticator()

	first, _ := g.NewState()
	second, _ := g.NewState()
	if first == second {
		t.Error("hai lần sinh state cho ra cùng một chuỗi, nonce không ngẫu nhiên")
	}
}

func TestStateBiSuaThiTuChoi(t *testing.T) {
	g := newTestAuthenticator()
	state, _ := g.NewState()
	parts := strings.Split(state, ".")

	cases := map[string]string{
		"đổi nonce":        "deadbeef." + parts[1] + "." + parts[2],
		"kéo dài hạn":      parts[0] + ".9999999999." + parts[2],
		"chữ ký sai":       parts[0] + "." + parts[1] + ".chu-ky-gia",
		"thiếu thành phần": parts[0] + "." + parts[1],
		"chuỗi rỗng":       "",
	}

	for name, tampered := range cases {
		t.Run(name, func(t *testing.T) {
			if err := g.VerifyState(tampered); err == nil {
				t.Error("state bị sửa nhưng vẫn được chấp nhận")
			}
		})
	}
}

func TestStateKhacKhoaThiTuChoi(t *testing.T) {
	state, _ := newTestAuthenticator().NewState()

	other := NewGoogleAuthenticator("client-id", "client-secret",
		"http://localhost:8081/api/auth/google/callback", "mot-khoa-hoan-toan-khac")
	if err := other.VerifyState(state); err == nil {
		t.Error("state ký bằng khoá khác nhưng vẫn được chấp nhận")
	}
}

func TestStateHetHan(t *testing.T) {
	g := newTestAuthenticator()

	// Dựng thủ công một state có hạn đã trôi qua, ký bằng đúng khoá thật.
	payload := "abc123." + itoa(time.Now().Add(-time.Minute).Unix())
	expired := payload + "." + g.sign(payload)

	err := g.VerifyState(expired)
	if err == nil {
		t.Fatal("state hết hạn nhưng vẫn được chấp nhận")
	}
	if !strings.Contains(err.Error(), "hết hạn") {
		t.Errorf("thông báo lỗi chưa nói rõ lý do hết hạn: %v", err)
	}
}

func TestAuthCodeURLChuaDuThamSo(t *testing.T) {
	g := newTestAuthenticator()
	url := g.AuthCodeURL("state-mau")

	for _, want := range []string{
		"accounts.google.com",
		"client_id=client-id",
		"state=state-mau",
		"prompt=select_account",
		"scope=",
	} {
		if !strings.Contains(url, want) {
			t.Errorf("URL đăng nhập thiếu %q\nURL: %s", want, url)
		}
	}
}

func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf []byte
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	if neg {
		return "-" + string(buf)
	}
	return string(buf)
}
