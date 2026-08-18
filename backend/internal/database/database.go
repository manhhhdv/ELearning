package database

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/manhnv/elearning/backend/internal/migrations"
)

// Connect mở pool kết nối và chờ database sẵn sàng.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("chuỗi kết nối không hợp lệ: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("không tạo được pool: %w", err)
	}

	// Postgres trong Docker có thể chưa nhận kết nối ngay khi tiến trình khởi động.
	var lastErr error
	for attempt := 1; attempt <= 10; attempt++ {
		if lastErr = pool.Ping(ctx); lastErr == nil {
			return pool, nil
		}
		slog.Warn("chờ database sẵn sàng", "lần", attempt, "lỗi", lastErr)
		select {
		case <-ctx.Done():
			pool.Close()
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	pool.Close()
	return nil, fmt.Errorf("không kết nối được database: %w", lastErr)
}

// Migrate chạy lần lượt các file .sql chưa được áp dụng, mỗi file trong một transaction.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version     text PRIMARY KEY,
			applied_at  timestamptz NOT NULL DEFAULT now()
		)`)
	if err != nil {
		return fmt.Errorf("tạo bảng schema_migrations: %w", err)
	}

	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("đọc schema_migrations: %w", err)
	}
	applied := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return err
		}
		applied[v] = true
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("đọc thư mục migrations: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		if applied[name] {
			continue
		}
		body, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("đọc migration %s: %w", name, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(body)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("chạy migration %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		slog.Info("đã áp dụng migration", "file", name)
	}
	return nil
}
