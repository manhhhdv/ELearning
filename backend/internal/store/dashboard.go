package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DashboardStats tổng hợp số liệu toàn hệ thống cho trang chủ quản trị.
type DashboardStats struct {
	ProgramsTotal     int `json:"programsTotal"`
	ProgramsDraft     int `json:"programsDraft"`
	ProgramsPublished int `json:"programsPublished"`
	ProgramsArchived  int `json:"programsArchived"`

	UsersTotal      int `json:"usersTotal"`
	AdminCount      int `json:"adminCount"`
	TrainerCount    int `json:"trainerCount"`
	SupervisorCount int `json:"supervisorCount"`
	StudentCount    int `json:"studentCount"`

	EnrollmentsTotal   int `json:"enrollmentsTotal"`
	SubmissionsTotal   int `json:"submissionsTotal"`
	SubmissionsPending int `json:"submissionsPending"`
	LessonsCompleted   int `json:"lessonsCompleted"`

	TopPrograms   []*DashboardProgram `json:"topPrograms"`
	RecentSignups []*DashboardUser    `json:"recentSignups"`
}

type DashboardProgram struct {
	ID              uuid.UUID `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Code            string    `json:"code"`
	Status          string    `json:"status"`
	EnrollmentCount int       `json:"enrollmentCount"`
}

type DashboardUser struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// DashboardStats gom số liệu tổng quan bằng vài truy vấn đếm nhanh — dữ liệu quy mô
// nhỏ (vài nghìn bản ghi) nên không cần bảng tổng hợp riêng.
func (s *Store) DashboardStats(ctx context.Context) (*DashboardStats, error) {
	out := &DashboardStats{TopPrograms: []*DashboardProgram{}, RecentSignups: []*DashboardUser{}}

	err := s.pool.QueryRow(ctx, `
		SELECT
			count(*),
			count(*) FILTER (WHERE status = 'draft'),
			count(*) FILTER (WHERE status = 'published'),
			count(*) FILTER (WHERE status = 'archived')
		FROM programs`).Scan(&out.ProgramsTotal, &out.ProgramsDraft, &out.ProgramsPublished, &out.ProgramsArchived)
	if err != nil {
		return nil, translate(err, "đếm chương trình")
	}

	err = s.pool.QueryRow(ctx, `
		SELECT
			count(*),
			count(*) FILTER (WHERE role = 'admin'),
			count(*) FILTER (WHERE role = 'trainer'),
			count(*) FILTER (WHERE role = 'supervisor'),
			count(*) FILTER (WHERE role = 'student')
		FROM users`).Scan(&out.UsersTotal, &out.AdminCount, &out.TrainerCount, &out.SupervisorCount, &out.StudentCount)
	if err != nil {
		return nil, translate(err, "đếm người dùng")
	}

	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM enrollments`).Scan(&out.EnrollmentsTotal); err != nil {
		return nil, translate(err, "đếm ghi danh")
	}
	err = s.pool.QueryRow(ctx, `
		SELECT count(*), count(*) FILTER (WHERE status = 'submitted') FROM submissions`).
		Scan(&out.SubmissionsTotal, &out.SubmissionsPending)
	if err != nil {
		return nil, translate(err, "đếm bài nộp")
	}
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM lesson_progress`).Scan(&out.LessonsCompleted); err != nil {
		return nil, translate(err, "đếm tiến độ học")
	}

	progRows, err := s.pool.Query(ctx, `
		SELECT p.id, p.slug, p.title, p.code, p.status,
		       (SELECT count(*) FROM enrollments e WHERE e.program_id = p.id) AS ec
		FROM programs p
		ORDER BY ec DESC, p.updated_at DESC
		LIMIT 5`)
	if err != nil {
		return nil, translate(err, "xếp hạng chương trình")
	}
	for progRows.Next() {
		var dp DashboardProgram
		if err := progRows.Scan(&dp.ID, &dp.Slug, &dp.Title, &dp.Code, &dp.Status, &dp.EnrollmentCount); err != nil {
			progRows.Close()
			return nil, translate(err, "xếp hạng chương trình")
		}
		out.TopPrograms = append(out.TopPrograms, &dp)
	}
	progRows.Close()
	if err := progRows.Err(); err != nil {
		return nil, err
	}

	userRows, err := s.pool.Query(ctx, `
		SELECT id, full_name, email, role, created_at
		FROM users ORDER BY created_at DESC LIMIT 5`)
	if err != nil {
		return nil, translate(err, "danh sách người dùng mới")
	}
	defer userRows.Close()
	for userRows.Next() {
		var du DashboardUser
		if err := userRows.Scan(&du.ID, &du.FullName, &du.Email, &du.Role, &du.CreatedAt); err != nil {
			return nil, translate(err, "danh sách người dùng mới")
		}
		out.RecentSignups = append(out.RecentSignups, &du)
	}
	return out, userRows.Err()
}
