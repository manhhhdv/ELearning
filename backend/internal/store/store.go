// Package store chứa toàn bộ truy vấn dữ liệu tới Postgres.
package store

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Các lỗi nghiệp vụ được tầng HTTP ánh xạ sang mã trạng thái phù hợp.
var (
	ErrNotFound = errors.New("không tìm thấy dữ liệu")
	ErrConflict = errors.New("dữ liệu đã tồn tại")
)

// InvalidError là vi phạm quy tắc nghiệp vụ do người dùng gây ra.
// Thông báo được viết cho người dùng cuối nên tầng HTTP trả nguyên văn kèm mã 400.
type InvalidError struct{ Msg string }

func (e *InvalidError) Error() string { return e.Msg }

// Invalidf tạo một InvalidError với thông báo đã định dạng.
func Invalidf(format string, args ...any) error {
	return &InvalidError{Msg: fmt.Sprintf(format, args...)}
}

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store { return &Store{pool: pool} }

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

// translate quy đổi lỗi pgx thành lỗi nghiệp vụ của package.
func translate(err error, action string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}
	return fmt.Errorf("%s: %w", action, err)
}
