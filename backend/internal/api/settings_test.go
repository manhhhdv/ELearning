package api

import "testing"

func TestMaskAPIKey(t *testing.T) {
	cases := map[string]string{
		"":                 "",
		"abcd":             "••••",
		"abc":              "••••",
		"AIzaSyD-xxxx1234": "••••1234",
	}
	for in, want := range cases {
		if got := maskAPIKey(in); got != want {
			t.Errorf("maskAPIKey(%q) = %q, mong đợi %q", in, got, want)
		}
	}
}

// Khoá API không bao giờ được trả nguyên văn ra ngoài, kể cả khi đủ ngắn.
func TestMaskAPIKeyNeverLeaksWholeKey(t *testing.T) {
	for _, key := range []string{"abcde", "AIzaSyD-verylongsecretkey"} {
		if got := maskAPIKey(key); got == key {
			t.Errorf("maskAPIKey(%q) trả về nguyên khoá", key)
		}
	}
}
