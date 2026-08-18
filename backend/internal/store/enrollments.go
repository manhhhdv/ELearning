package store

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// Enroll ghi danh một người dùng vào chương trình; gọi lại nhiều lần chỉ cập nhật vai trò.
func (s *Store) Enroll(ctx context.Context, programID, userID uuid.UUID, role string) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO enrollments (program_id, user_id, role) VALUES ($1, $2, $3)
		ON CONFLICT (program_id, user_id) DO UPDATE SET role = EXCLUDED.role`,
		programID, userID, role)
	return translate(err, "ghi danh")
}

func (s *Store) Unenroll(ctx context.Context, programID, userID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM enrollments WHERE program_id = $1 AND user_id = $2`, programID, userID)
	if err != nil {
		return translate(err, "huỷ ghi danh")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ListProgramEnrollments liệt kê học viên / giảng viên của một chương trình.
func (s *Store) ListProgramEnrollments(ctx context.Context, programID uuid.UUID) ([]*models.Enrollment, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT e.id, e.program_id, e.user_id, e.role, e.enrolled_at, u.email, u.full_name
		FROM enrollments e JOIN users u ON u.id = e.user_id
		WHERE e.program_id = $1
		ORDER BY u.full_name, u.email`, programID)
	if err != nil {
		return nil, translate(err, "liệt kê ghi danh")
	}
	defer rows.Close()

	out := []*models.Enrollment{}
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.ID, &e.ProgramID, &e.UserID, &e.Role, &e.EnrolledAt, &e.Email, &e.FullName); err != nil {
			return nil, translate(err, "đọc ghi danh")
		}
		out = append(out, &e)
	}
	return out, rows.Err()
}

// EnrollmentRole trả về vai trò của người dùng trong chương trình ("" nếu chưa ghi danh).
func (s *Store) EnrollmentRole(ctx context.Context, programID, userID uuid.UUID) (string, error) {
	var role string
	err := s.pool.QueryRow(ctx, `
		SELECT role FROM enrollments WHERE program_id = $1 AND user_id = $2`, programID, userID).Scan(&role)
	if err != nil {
		if err := translate(err, "đọc vai trò ghi danh"); !errors.Is(err, ErrNotFound) {
			return "", err
		}
		return "", nil
	}
	return role, nil
}

// MarkLessonComplete đánh dấu học viên đã hoàn thành một bài học.
func (s *Store) MarkLessonComplete(ctx context.Context, userID, nodeID uuid.UUID, done bool) error {
	if !done {
		_, err := s.pool.Exec(ctx, `DELETE FROM lesson_progress WHERE user_id = $1 AND node_id = $2`, userID, nodeID)
		return translate(err, "bỏ đánh dấu hoàn thành")
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO lesson_progress (user_id, node_id) VALUES ($1, $2)
		ON CONFLICT (user_id, node_id) DO NOTHING`, userID, nodeID)
	return translate(err, "đánh dấu hoàn thành")
}

// CompletedNodeIDs trả về các bài học người dùng đã hoàn thành trong một chương trình.
func (s *Store) CompletedNodeIDs(ctx context.Context, userID, programID uuid.UUID) (map[uuid.UUID]bool, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT lp.node_id FROM lesson_progress lp
		JOIN nodes n ON n.id = lp.node_id
		WHERE lp.user_id = $1 AND n.program_id = $2`, userID, programID)
	if err != nil {
		return nil, translate(err, "đọc tiến độ học")
	}
	defer rows.Close()

	done := map[uuid.UUID]bool{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, translate(err, "đọc tiến độ học")
		}
		done[id] = true
	}
	return done, rows.Err()
}
