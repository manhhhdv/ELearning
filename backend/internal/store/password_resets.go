package store

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// CreatePasswordResetRequest ghi nhận một yêu cầu đặt lại mật khẩu.
//
// Email không có trong hệ thống vẫn được ghi (user_id để trống) để trang đăng
// nhập không lộ email nào tồn tại; admin nhìn danh sách sẽ thấy ngay dòng nào
// không khớp tài khoản nào.
//
// Bấm lại khi đang có yêu cầu chờ xử lý thì không tạo thêm dòng mới.
func (s *Store) CreatePasswordResetRequest(ctx context.Context, email, note string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := s.pool.Exec(ctx, `
		INSERT INTO password_reset_requests (email, user_id, note)
		VALUES ($1, (SELECT id FROM users WHERE lower(email) = $1), $2)
		ON CONFLICT (lower(email)) WHERE status = 'pending' DO NOTHING`,
		email, note)
	return translate(err, "ghi nhận yêu cầu đặt lại mật khẩu")
}

// ListPasswordResetRequests trả về các yêu cầu, mới nhất trước.
// pendingOnly=true chỉ lấy những yêu cầu admin chưa xử lý.
func (s *Store) ListPasswordResetRequests(ctx context.Context, pendingOnly bool) ([]*models.PasswordResetRequest, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT r.id, r.email, r.user_id, r.note, r.status, r.handled_by, r.handled_at, r.created_at,
		       COALESCE(u.full_name, ''), COALESCE(h.full_name, '')
		FROM password_reset_requests r
		LEFT JOIN users u ON u.id = r.user_id
		LEFT JOIN users h ON h.id = r.handled_by
		WHERE ($1 = false OR r.status = 'pending')
		ORDER BY r.created_at DESC
		LIMIT 200`, pendingOnly)
	if err != nil {
		return nil, translate(err, "liệt kê yêu cầu đặt lại mật khẩu")
	}
	defer rows.Close()

	out := []*models.PasswordResetRequest{}
	for rows.Next() {
		var q models.PasswordResetRequest
		if err := rows.Scan(&q.ID, &q.Email, &q.UserID, &q.Note, &q.Status,
			&q.HandledBy, &q.HandledAt, &q.CreatedAt, &q.FullName, &q.HandledByName); err != nil {
			return nil, translate(err, "đọc yêu cầu đặt lại mật khẩu")
		}
		out = append(out, &q)
	}
	return out, rows.Err()
}

// CountPendingPasswordResets đếm số yêu cầu đang chờ, dùng cho chỉ báo ở giao diện admin.
func (s *Store) CountPendingPasswordResets(ctx context.Context) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM password_reset_requests WHERE status = 'pending'`).Scan(&n)
	return n, translate(err, "đếm yêu cầu đặt lại mật khẩu")
}

// ResolvePasswordResetRequest đánh dấu một yêu cầu đã xử lý xong hoặc bị từ chối.
func (s *Store) ResolvePasswordResetRequest(ctx context.Context, id, adminID uuid.UUID, status string) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE password_reset_requests
		SET status = $3, handled_by = $2, handled_at = now()
		WHERE id = $1 AND status = 'pending'`, id, adminID, status)
	if err != nil {
		return translate(err, "cập nhật yêu cầu đặt lại mật khẩu")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ResolvePasswordResetsForUser đóng mọi yêu cầu đang chờ của một tài khoản.
// Gọi sau khi admin cấp mật khẩu mới, để danh sách chờ tự sạch mà admin không
// phải bấm thêm một lần nữa.
func (s *Store) ResolvePasswordResetsForUser(ctx context.Context, userID, adminID uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE password_reset_requests
		SET status = 'done', handled_by = $2, handled_at = now()
		WHERE user_id = $1 AND status = 'pending'`, userID, adminID)
	return translate(err, "đóng yêu cầu đặt lại mật khẩu")
}
