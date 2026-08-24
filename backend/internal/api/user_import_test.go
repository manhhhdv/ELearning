package api

import (
	"testing"

	"github.com/manhnv/elearning/backend/internal/models"
)

func TestCleanUserItemNormalizes(t *testing.T) {
	got, msg := cleanUserItem(userItemInput{Email: "  An.Nguyen@CongTy.VN ", FullName: " Nguyễn Văn An "})
	if msg != "" {
		t.Fatalf("unexpected error: %s", msg)
	}
	if got.Email != "an.nguyen@congty.vn" {
		t.Fatalf("email not normalized: %q", got.Email)
	}
	if got.FullName != "Nguyễn Văn An" {
		t.Fatalf("full name not trimmed: %q", got.FullName)
	}
	if got.Role != models.RoleStudent {
		t.Fatalf("expected default role student, got %q", got.Role)
	}
}

func TestCleanUserItemRejectsBadInput(t *testing.T) {
	cases := []struct {
		name string
		in   userItemInput
	}{
		{"thiếu email", userItemInput{FullName: "Nguyễn Văn An"}},
		{"email sai định dạng", userItemInput{Email: "an.nguyen(a)congty"}},
		{"email kèm tên hiển thị", userItemInput{Email: "An <an@congty.vn>"}},
		{"email không có tên miền đủ", userItemInput{Email: "an@congty"}},
		{"vai trò lạ", userItemInput{Email: "an@congty.vn", Role: "giám đốc"}},
		{"mật khẩu quá ngắn", userItemInput{Email: "an@congty.vn", Password: "abc1"}},
		{"mật khẩu thiếu số", userItemInput{Email: "an@congty.vn", Password: "matkhaudai"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, msg := cleanUserItem(tc.in); msg == "" {
				t.Fatalf("expected error for %+v", tc.in)
			}
		})
	}
}

func TestCleanUserItemAcceptsValidRolesAndPassword(t *testing.T) {
	for _, role := range []string{models.RoleAdmin, models.RoleTrainer, models.RoleSupervisor, models.RoleStudent} {
		got, msg := cleanUserItem(userItemInput{Email: "an@congty.vn", Role: " " + role + " ", Password: "MatKhau123"})
		if msg != "" {
			t.Fatalf("unexpected error for role %q: %s", role, msg)
		}
		if got.Role != role || got.Password != "MatKhau123" {
			t.Fatalf("unexpected item: %+v", got)
		}
	}
}
