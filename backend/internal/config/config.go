package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config gom toàn bộ tham số vận hành, đọc từ biến môi trường / file .env.
type Config struct {
	Port        string
	DatabaseURL string
	Env         string

	JWTSecret    string
	JWTTTL       time.Duration
	AllowOrigins []string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	// Chỉ cho phép đăng nhập Google với email thuộc các domain này (rỗng = không giới hạn).
	GoogleAllowedDomains []string
	// Email đăng nhập Google lần đầu mà chưa có trong hệ thống sẽ được tạo tự động với vai trò này.
	// Để trống nghĩa là chỉ những tài khoản admin đã tạo sẵn mới đăng nhập được.
	GoogleAutoProvisionRole string

	FrontendURL string

	SeedAdminEmail    string
	SeedAdminPassword string
	SeedAdminName     string
}

func Load() (*Config, error) {
	// Không lỗi nếu thiếu .env — môi trường production thường dùng biến môi trường thật.
	_ = godotenv.Load(".env", "../.env")

	cfg := &Config{
		Port:                    env("PORT", "8082"),
		DatabaseURL:             env("DATABASE_URL", "postgres://elearning:elearning@localhost:5433/elearning?sslmode=disable"),
		Env:                     env("APP_ENV", "development"),
		JWTSecret:               env("JWT_SECRET", ""),
		JWTTTL:                  envDuration("JWT_TTL", 24*time.Hour),
		AllowOrigins:            envList("ALLOW_ORIGINS", "http://localhost:3006"),
		GoogleClientID:          env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:      env("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:       env("GOOGLE_REDIRECT_URL", "http://localhost:8082/api/auth/google/callback"),
		GoogleAllowedDomains:    envList("GOOGLE_ALLOWED_DOMAINS", ""),
		GoogleAutoProvisionRole: env("GOOGLE_AUTO_PROVISION_ROLE", ""),
		FrontendURL:             env("FRONTEND_URL", "http://localhost:3006"),
		SeedAdminEmail:          env("SEED_ADMIN_EMAIL", "admin@elearning.local"),
		SeedAdminPassword:       env("SEED_ADMIN_PASSWORD", "Admin@12345"),
		SeedAdminName:           env("SEED_ADMIN_NAME", "Quản trị hệ thống"),
	}

	if cfg.JWTSecret == "" {
		if cfg.Env == "production" {
			return nil, fmt.Errorf("JWT_SECRET là bắt buộc khi APP_ENV=production")
		}
		cfg.JWTSecret = "dev-secret-doi-truoc-khi-len-production"
	}
	if r := cfg.GoogleAutoProvisionRole; r != "" && r != "student" && r != "trainer" {
		return nil, fmt.Errorf("GOOGLE_AUTO_PROVISION_ROLE chỉ nhận giá trị student hoặc trainer, nhận được %q", r)
	}
	return cfg, nil
}

// GoogleEnabled cho biết đã cấu hình đủ credentials để bật nút đăng nhập Google chưa.
func (c *Config) GoogleEnabled() bool {
	return c.GoogleClientID != "" && c.GoogleClientSecret != ""
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envList(key, fallback string) []string {
	raw := env(key, fallback)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envDuration(key string, fallback time.Duration) time.Duration {
	raw := env(key, "")
	if raw == "" {
		return fallback
	}
	if d, err := time.ParseDuration(raw); err == nil {
		return d
	}
	if hours, err := strconv.Atoi(raw); err == nil {
		return time.Duration(hours) * time.Hour
	}
	return fallback
}
