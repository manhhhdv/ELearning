package api

import "testing"

func TestDomainAllowed(t *testing.T) {
	cases := []struct {
		name    string
		domains []string
		email   string
		want    bool
	}{
		{"không giới hạn thì nhận mọi email", nil, "ai.do@gmail.com", true},
		{"đúng domain cho phép", []string{"congty.vn"}, "an@congty.vn", true},
		{"không phân biệt hoa thường", []string{"CongTy.vn"}, "an@CONGTY.VN", true},
		{"chấp nhận khai báo có @ đứng đầu", []string{"@congty.vn"}, "an@congty.vn", true},
		{"một trong nhiều domain", []string{"a.vn", "congty.vn"}, "an@congty.vn", true},
		{"sai domain thì chặn", []string{"congty.vn"}, "an@gmail.com", false},
		{"tên miền con không tự động được nhận", []string{"congty.vn"}, "an@chinhanh.congty.vn", false},
		{"chuỗi không phải email", []string{"congty.vn"}, "khong-co-a-cong", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := domainAllowed(tc.email, tc.domains); got != tc.want {
				t.Errorf("domainAllowed(%q, %v) = %v, mong đợi %v", tc.email, tc.domains, got, tc.want)
			}
		})
	}
}

func TestSplitDomains(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"congty.vn", []string{"congty.vn"}},
		{"congty.vn, chinhanh.congty.vn", []string{"congty.vn", "chinhanh.congty.vn"}},
		{"  a.vn ,, b.vn  ", []string{"a.vn", "b.vn"}},
	}
	for _, tc := range cases {
		got := splitDomains(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("splitDomains(%q) = %v, mong đợi %v", tc.in, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("splitDomains(%q) = %v, mong đợi %v", tc.in, got, tc.want)
				break
			}
		}
	}
}
