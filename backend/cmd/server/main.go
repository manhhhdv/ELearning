// Command server khởi động API của hệ thống đào tạo trực tuyến.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/manhnv/elearning/backend/internal/api"
	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/config"
	"github.com/manhnv/elearning/backend/internal/database"
	"github.com/manhnv/elearning/backend/internal/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	if err := run(); err != nil {
		slog.Error("máy chủ dừng do lỗi", "lỗi", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		return err
	}

	st := store.New(pool)
	if err := seedAdmin(ctx, st, cfg); err != nil {
		return err
	}
	if cfg.SeedSampleSubjects {
		if err := seedSampleSubjects(ctx, st, cfg); err != nil {
			return err
		}
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.NewServer(cfg, st).Router(),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("API đang lắng nghe", "địa chỉ", srv.Addr, "google_login", cfg.GoogleEnabled())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("nhận tín hiệu dừng, đang tắt máy chủ")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// seedAdmin tạo tài khoản quản trị đầu tiên khi database còn trống.
func seedAdmin(ctx context.Context, st *store.Store, cfg *config.Config) error {
	hash, err := auth.HashPassword(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}
	created, err := st.EnsureSeedAdmin(ctx, cfg.SeedAdminEmail, cfg.SeedAdminName, hash)
	if err != nil {
		return err
	}
	if created {
		slog.Warn("đã tạo tài khoản quản trị đầu tiên — hãy đổi mật khẩu ngay sau khi đăng nhập",
			"email", cfg.SeedAdminEmail)
	}
	return nil
}

// seedSampleSubjects tạo sẵn các lớp học mẫu lớp 1–12 (Toán, Tiếng Việt/Ngữ văn, Tiếng Anh)
// khi cfg.SeedSampleSubjects bật, để có ngay dữ liệu thử nghiệm cho các chức
// năng tạo bài giảng / ra đề / ghi danh thay vì phải tạo tay từng lớp. Bật
// mặc định ở mọi môi trường trừ production; đặt SEED_SAMPLE_SUBJECTS=true để
// chủ động bật ở production khi cần seed dữ liệu demo lên đó.
func seedSampleSubjects(ctx context.Context, st *store.Store, cfg *config.Config) error {
	admin, err := st.GetUserByEmail(ctx, cfg.SeedAdminEmail)
	if err != nil {
		// Không có tài khoản quản trị nào để gán làm người tạo thì bỏ qua,
		// không chặn khởi động máy chủ vì việc này.
		slog.Warn("không tìm được tài khoản quản trị để gán lớp học mẫu, bỏ qua", "lỗi", err)
		return nil
	}
	return st.EnsureSampleSubjects(ctx, admin.ID)
}
