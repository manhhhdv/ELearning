package util

import "testing"

func TestSlugify(t *testing.T) {
	cases := []struct{ in, want string }{
		{"An toàn thực phẩm 2026", "an-toan-thuc-pham-2026"},
		{"Chương 1 — Kiến thức cơ bản", "chuong-1-kien-thuc-co-ban"},
		{"Đăng ký & Đào tạo", "dang-ky-dao-tao"},
		{"ATTP-2026", "attp-2026"},
		{"   ", "muc"},
		{"!!!", "muc"},
		{"Bài 1: Giới thiệu", "bai-1-gioi-thieu"},
	}
	for _, tc := range cases {
		if got := Slugify(tc.in); got != tc.want {
			t.Errorf("Slugify(%q) = %q, muốn %q", tc.in, got, tc.want)
		}
	}
}
