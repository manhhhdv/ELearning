package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/util"
)

const programColumns = `
	p.id, p.code, p.slug, p.title, p.description, p.cover_url, p.status,
	p.allow_self_enroll, p.is_default_course,
	p.created_by, p.created_at, p.updated_at`

// programStats đếm số nút, bài học, bài tập và lượt ghi danh của từng chương trình.
// Tham số cuối là ID người dùng đang xem, dùng để đếm số bài họ đã hoàn thành (NULL = bỏ qua).
// Khoá học mặc định coi như "đã ghi danh" với mọi người xem, dù không có dòng enrollments nào.
const programStats = `
	COALESCE((SELECT count(*) FROM nodes n WHERE n.program_id = p.id), 0),
	COALESCE((SELECT count(*) FROM nodes n WHERE n.program_id = p.id AND n.kind = 'lesson'), 0),
	COALESCE((SELECT count(*) FROM nodes n WHERE n.program_id = p.id AND n.kind = 'assignment'), 0),
	COALESCE((SELECT count(*) FROM enrollments e WHERE e.program_id = p.id), 0),
	COALESCE((SELECT count(*) FROM lesson_progress lp
	          JOIN nodes n ON n.id = lp.node_id
	          WHERE n.program_id = p.id AND lp.user_id = $%d::uuid), 0),
	(p.is_default_course OR EXISTS (
		SELECT 1 FROM enrollments e WHERE e.program_id = p.id AND e.user_id = $%d::uuid))`

type programScanner interface {
	Scan(dest ...any) error
}

func scanProgram(row programScanner) (*models.Program, error) {
	var p models.Program
	err := row.Scan(&p.ID, &p.Code, &p.Slug, &p.Title, &p.Description, &p.CoverURL, &p.Status,
		&p.AllowSelfEnroll, &p.IsDefaultCourse,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.NodeCount, &p.LessonCount, &p.AssignmentCount, &p.EnrollmentCount,
		&p.CompletedLessonCount, &p.Enrolled)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

type CreateProgramParams struct {
	Code            string
	Title           string
	Description     string
	CoverURL        string
	Status          string
	AllowSelfEnroll bool
	IsDefaultCourse bool
	CreatedBy       uuid.UUID
}

func (s *Store) CreateProgram(ctx context.Context, p CreateProgramParams) (*models.Program, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	code := strings.TrimSpace(p.Code)
	slug, err := uniqueProgramSlug(ctx, tx, util.Slugify(code))
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO programs (code, slug, title, description, cover_url, status, allow_self_enroll, is_default_course, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		code, slug, p.Title, p.Description, p.CoverURL, p.Status, p.AllowSelfEnroll, p.IsDefaultCourse, p.CreatedBy).Scan(&id)
	if err != nil {
		return nil, translate(err, "tạo chương trình")
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetProgram(ctx, id, uuid.Nil)
}

// uniqueProgramSlug thêm hậu tố -2, -3… nếu slug cơ sở đã tồn tại.
func uniqueProgramSlug(ctx context.Context, tx pgx.Tx, base string) (string, error) {
	return uniqueProgramSlugExcluding(ctx, tx, base, uuid.Nil)
}

// uniqueProgramSlugExcluding giống uniqueProgramSlug nhưng bỏ qua một ID cụ thể —
// dùng khi đổi mã của chính chương trình đó, để không tự đụng độ với slug cũ của mình.
func uniqueProgramSlugExcluding(ctx context.Context, tx pgx.Tx, base string, excluding uuid.UUID) (string, error) {
	candidate := base
	for i := 2; ; i++ {
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM programs WHERE slug = $1 AND id <> $2)`,
			candidate, excluding).Scan(&exists); err != nil {
			return "", translate(err, "kiểm tra slug")
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// GetProgramBySlug đọc chương trình theo slug trên URL thay vì UUID nội bộ.
func (s *Store) GetProgramBySlug(ctx context.Context, slug string, viewerID uuid.UUID) (*models.Program, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+programColumns+`, `+fmt.Sprintf(programStats, 2, 2)+` FROM programs p WHERE p.slug = $1`,
		slug, nullableUUID(viewerID))
	p, err := scanProgram(row)
	if err != nil {
		return nil, translate(err, "đọc chương trình")
	}
	return p, nil
}

// GetProgram đọc một chương trình. viewerID khác uuid.Nil thì kèm theo tiến độ học của người đó.
func (s *Store) GetProgram(ctx context.Context, id uuid.UUID, viewerID uuid.UUID) (*models.Program, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+programColumns+`, `+fmt.Sprintf(programStats, 2, 2)+` FROM programs p WHERE p.id = $1`,
		id, nullableUUID(viewerID))
	p, err := scanProgram(row)
	if err != nil {
		return nil, translate(err, "đọc chương trình")
	}
	return p, nil
}

type ListProgramsFilter struct {
	Search string
	Status string
	// Khi khác uuid.Nil, chỉ trả về chương trình mà người dùng này được ghi danh.
	EnrolledUserID uuid.UUID
	// Khi khác uuid.Nil, chỉ trả về chương trình do người dùng này tạo hoặc được ghi danh.
	// Dùng cho giảng viên: không cho thấy chương trình của người khác.
	VisibleToUserID uuid.UUID
	// Người đang xem, dùng để tính số bài đã hoàn thành trả kèm mỗi chương trình.
	ViewerID uuid.UUID
	// Chỉ lấy chương trình đang mở cho học viên tự ghi danh.
	OnlySelfEnrollable bool
}

func (s *Store) ListPrograms(ctx context.Context, f ListProgramsFilter) ([]*models.Program, error) {
	search := "%" + strings.ToLower(strings.TrimSpace(f.Search)) + "%"

	rows, err := s.pool.Query(ctx, `
		SELECT `+programColumns+`, `+fmt.Sprintf(programStats, 5, 5)+`
		FROM programs p
		WHERE ($1 = '' OR p.status = $1)
		  AND ($2 = '%%' OR lower(p.title) LIKE $2 OR lower(p.code) LIKE $2)
		  AND ($3::uuid IS NULL OR p.is_default_course OR EXISTS (
		        SELECT 1 FROM enrollments e WHERE e.program_id = p.id AND e.user_id = $3::uuid))
		  AND ($4::uuid IS NULL OR p.created_by = $4::uuid OR EXISTS (
		        SELECT 1 FROM enrollments e WHERE e.program_id = p.id AND e.user_id = $4::uuid))
		  AND ($6 = false OR p.allow_self_enroll)
		ORDER BY p.updated_at DESC
	`, f.Status, search, nullableUUID(f.EnrolledUserID), nullableUUID(f.VisibleToUserID), nullableUUID(f.ViewerID), f.OnlySelfEnrollable)
	if err != nil {
		return nil, translate(err, "liệt kê chương trình")
	}
	defer rows.Close()

	out := []*models.Program{}
	for rows.Next() {
		p, err := scanProgram(rows)
		if err != nil {
			return nil, translate(err, "đọc chương trình")
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

type UpdateProgramParams struct {
	Code            *string
	Title           *string
	Description     *string
	CoverURL        *string
	Status          *string
	AllowSelfEnroll *bool
	IsDefaultCourse *bool
}

func (s *Store) UpdateProgram(ctx context.Context, id uuid.UUID, p UpdateProgramParams) (*models.Program, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Đổi mã thì slug đổi theo (loại trừ chính bản ghi này khi kiểm tra trùng),
	// để URL luôn khớp với mã hiện tại thay vì giữ mã cũ đã bỏ.
	var slug *string
	if p.Code != nil {
		s, err := uniqueProgramSlugExcluding(ctx, tx, util.Slugify(*p.Code), id)
		if err != nil {
			return nil, err
		}
		slug = &s
	}

	tag, err := tx.Exec(ctx, `
		UPDATE programs SET
			code        = COALESCE($2, code),
			slug        = COALESCE($3, slug),
			title       = COALESCE($4, title),
			description = COALESCE($5, description),
			cover_url   = COALESCE($6, cover_url),
			status      = COALESCE($7, status),
			allow_self_enroll = COALESCE($8, allow_self_enroll),
			is_default_course = COALESCE($9, is_default_course)
		WHERE id = $1`, id, p.Code, slug, p.Title, p.Description, p.CoverURL, p.Status,
		p.AllowSelfEnroll, p.IsDefaultCourse)
	if err != nil {
		return nil, translate(err, "cập nhật chương trình")
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetProgram(ctx, id, uuid.Nil)
}

func (s *Store) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM programs WHERE id = $1`, id)
	if err != nil {
		return translate(err, "xoá chương trình")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// IsEnrolled kiểm tra một người dùng có được ghi danh vào chương trình hay không.
func (s *Store) IsEnrolled(ctx context.Context, programID, userID uuid.UUID) (bool, error) {
	var ok bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM enrollments WHERE program_id = $1 AND user_id = $2)`,
		programID, userID).Scan(&ok)
	if err != nil {
		return false, translate(err, "kiểm tra ghi danh")
	}
	return ok, nil
}

func nullableUUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}
