package config

import "testing"

// TestSeedSampleSubjectsDefault kiểm tra cờ tạo dữ liệu mẫu (lớp học mẫu lớp 1–12, thử
// nghiệm): mặc định bật ở mọi môi trường trừ production, và SEED_SAMPLE_SUBJECTS
// luôn ghi đè tường minh được — kể cả để chủ động bật ở production khi cần
// seed dữ liệu demo lên đó.
func TestSeedSampleSubjectsDefault(t *testing.T) {
	cases := []struct {
		name     string
		env      string
		override string // rỗng = không đặt SEED_SAMPLE_SUBJECTS
		want     bool
	}{
		{name: "development mặc định bật", env: "development", want: true},
		{name: "production mặc định tắt", env: "production", want: false},
		{name: "production bật tường minh", env: "production", override: "true", want: true},
		{name: "development tắt tường minh", env: "development", override: "false", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tc.env)
			t.Setenv("JWT_SECRET", "kiem-tra-khong-dung-that")
			if tc.override == "" {
				t.Setenv("SEED_SAMPLE_SUBJECTS", "")
			} else {
				t.Setenv("SEED_SAMPLE_SUBJECTS", tc.override)
			}

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() lỗi: %v", err)
			}
			if cfg.SeedSampleSubjects != tc.want {
				t.Fatalf("SeedSampleSubjects = %v, muốn %v", cfg.SeedSampleSubjects, tc.want)
			}
		})
	}
}
