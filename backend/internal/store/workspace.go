package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Workspace là dữ liệu cho trang tổng quan ("không gian làm việc AI").
// Phần nào không thuộc vai trò người đang xem thì để rỗng.
type Workspace struct {
	// Việc cần làm hôm nay: bài được giao sắp tới hạn hoặc đã quá hạn.
	Schedule []*ScheduleItem `json:"schedule"`
	// Tiến độ học theo từng lớp: số bài đã hoàn thành trên tổng số bài.
	Progress []*CourseProgress `json:"progress"`

	// --- Chỉ dành cho người dạy và quản trị viên ---
	// Bài nộp đang chờ chấm.
	PendingGrading int `json:"pendingGrading"`
	// Cảnh báo học tập: học viên cần chú ý (điểm thấp, quá hạn chưa nộp).
	Alerts []*LearningAlert `json:"alerts"`
}

// ScheduleItem là một bài tập được giao, sắp xếp theo hạn nộp.
type ScheduleItem struct {
	NodeID      uuid.UUID  `json:"nodeId"`
	NodeSlug    string     `json:"nodeSlug"`
	Title       string     `json:"title"`
	ProgramID   uuid.UUID  `json:"programId"`
	ProgramSlug string     `json:"programSlug"`
	ProgramName string     `json:"programName"`
	DueAt       *time.Time `json:"dueAt"`
	Submitted   bool       `json:"submitted"`
	Overdue     bool       `json:"overdue"`
}

// CourseProgress là tiến độ học một lớp của người đang xem.
type CourseProgress struct {
	ProgramID   uuid.UUID `json:"programId"`
	ProgramSlug string    `json:"programSlug"`
	Title       string    `json:"title"`
	Completed   int       `json:"completed"`
	Total       int       `json:"total"`
	// Điểm trung bình các bài đã nộp của lớp này; nil khi chưa nộp bài nào.
	AverageScore *float64 `json:"averageScore"`
}

// LearningAlert là một học viên cần giáo viên chú ý.
type LearningAlert struct {
	UserID      uuid.UUID `json:"userId"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	ProgramID   uuid.UUID `json:"programId"`
	ProgramSlug string    `json:"programSlug"`
	ProgramName string    `json:"programName"`
	// "low_score" điểm trung bình dưới ngưỡng đạt; "overdue" có bài quá hạn chưa nộp.
	Kind string `json:"kind"`
	// Số liệu kèm theo: điểm trung bình (low_score) hoặc số bài quá hạn (overdue).
	Value float64 `json:"value"`
}

// scheduleHorizon là khoảng nhìn về phía trước của "lịch học trong ngày".
// Lấy rộng hơn một ngày để học viên thấy trước bài sắp tới hạn, thay vì chỉ
// biết đúng hôm bài hết hạn.
const scheduleHorizon = 7 * 24 * time.Hour

// StudentWorkspace gom lịch học và tiến độ của một học viên.
func (s *Store) StudentWorkspace(ctx context.Context, userID uuid.UUID) (*Workspace, error) {
	ws := &Workspace{Schedule: []*ScheduleItem{}, Progress: []*CourseProgress{}, Alerts: []*LearningAlert{}}

	// Bài tập đã xuất bản trong các lớp người dùng đang theo học, còn hạn hoặc
	// đã quá hạn mà chưa nộp. Bài không đặt hạn cũng hiện nếu chưa làm.
	rows, err := s.pool.Query(ctx, `
		SELECT n.id, n.slug, n.title, p.id, p.slug, p.title, a.due_at,
		       EXISTS (SELECT 1 FROM submissions sub
		               WHERE sub.assignment_id = n.id AND sub.user_id = $1) AS submitted
		FROM nodes n
		JOIN assignments a ON a.node_id = n.id
		JOIN programs p ON p.id = n.program_id
		WHERE n.kind = 'assignment'
		  AND n.is_published = true
		  AND p.status = 'published'
		  AND (
		    EXISTS (SELECT 1 FROM enrollments e WHERE e.program_id = p.id AND e.user_id = $1)
		    OR p.is_default_course = true
		  )
		  AND (a.due_at IS NULL OR a.due_at <= now() + $2::interval)
		ORDER BY a.due_at NULLS LAST, n.title
		LIMIT 50`, userID, scheduleHorizon.String())
	if err != nil {
		return nil, translate(err, "đọc lịch học")
	}
	defer rows.Close()

	now := time.Now()
	for rows.Next() {
		var it ScheduleItem
		if err := rows.Scan(&it.NodeID, &it.NodeSlug, &it.Title, &it.ProgramID,
			&it.ProgramSlug, &it.ProgramName, &it.DueAt, &it.Submitted); err != nil {
			return nil, translate(err, "đọc lịch học")
		}
		// Bài đã nộp rồi thì không còn là việc cần làm.
		if it.Submitted {
			continue
		}
		it.Overdue = it.DueAt != nil && it.DueAt.Before(now)
		ws.Schedule = append(ws.Schedule, &it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	progress, err := s.courseProgress(ctx, userID)
	if err != nil {
		return nil, err
	}
	ws.Progress = progress
	return ws, nil
}

// courseProgress tính số bài đã hoàn thành trên tổng số bài của từng lớp,
// kèm điểm trung bình các bài đã nộp (mục 3.5).
func (s *Store) courseProgress(ctx context.Context, userID uuid.UUID) ([]*CourseProgress, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT p.id, p.slug, p.title,
		       COALESCE((SELECT count(*) FROM lesson_progress lp
		                 JOIN nodes ln ON ln.id = lp.node_id
		                 WHERE ln.program_id = p.id AND lp.user_id = $1), 0) AS completed,
		       COALESCE((SELECT count(*) FROM nodes ln
		                 WHERE ln.program_id = p.id AND ln.kind = 'lesson' AND ln.is_published = true), 0) AS total,
		       (SELECT avg(sub.auto_score + COALESCE(sub.manual_score, 0))
		        FROM submissions sub
		        JOIN nodes an ON an.id = sub.assignment_id
		        WHERE an.program_id = p.id AND sub.user_id = $1) AS avg_score
		FROM programs p
		WHERE p.status = 'published'
		  AND (
		    EXISTS (SELECT 1 FROM enrollments e WHERE e.program_id = p.id AND e.user_id = $1)
		    OR p.is_default_course = true
		  )
		ORDER BY p.title
		LIMIT 50`, userID)
	if err != nil {
		return nil, translate(err, "đọc tiến độ học tập")
	}
	defer rows.Close()

	out := []*CourseProgress{}
	for rows.Next() {
		var c CourseProgress
		if err := rows.Scan(&c.ProgramID, &c.ProgramSlug, &c.Title,
			&c.Completed, &c.Total, &c.AverageScore); err != nil {
			return nil, translate(err, "đọc tiến độ học tập")
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

// TrainerWorkspace gom số bài chờ chấm và các cảnh báo học tập trong phạm vi
// những lớp mà người dạy phụ trách. Admin thấy toàn hệ thống.
func (s *Store) TrainerWorkspace(ctx context.Context, userID uuid.UUID, allPrograms bool) (*Workspace, error) {
	ws := &Workspace{Schedule: []*ScheduleItem{}, Progress: []*CourseProgress{}, Alerts: []*LearningAlert{}}

	if err := s.pool.QueryRow(ctx, `
		SELECT count(*)
		FROM submissions sub
		JOIN nodes n ON n.id = sub.assignment_id
		WHERE sub.status = 'submitted'
		  AND ($2 OR EXISTS (SELECT 1 FROM enrollments e
		                     WHERE e.program_id = n.program_id AND e.user_id = $1 AND e.role = 'trainer'))`,
		userID, allPrograms).Scan(&ws.PendingGrading); err != nil {
		return nil, translate(err, "đếm bài chờ chấm")
	}

	alerts, err := s.learningAlerts(ctx, userID, allPrograms)
	if err != nil {
		return nil, err
	}
	ws.Alerts = alerts
	return ws, nil
}

// learningAlerts tìm học viên cần chú ý theo hai quy tắc đơn giản:
// điểm trung bình dưới điểm đạt của bài, hoặc có bài quá hạn chưa nộp.
//
// Đây là luật cố định, không phải phân tích bằng AI — việc dùng AI để cảnh báo
// sớm là hướng phát triển tiếp theo.
func (s *Store) learningAlerts(ctx context.Context, userID uuid.UUID, allPrograms bool) ([]*LearningAlert, error) {
	rows, err := s.pool.Query(ctx, `
		WITH scope AS (
		    SELECT p.id, p.slug, p.title
		    FROM programs p
		    WHERE p.status = 'published'
		      AND ($2 OR EXISTS (SELECT 1 FROM enrollments e
		                         WHERE e.program_id = p.id AND e.user_id = $1 AND e.role = 'trainer'))
		),
		learners AS (
		    SELECT e.user_id, u.full_name, u.email, s.id AS program_id, s.slug, s.title
		    FROM enrollments e
		    JOIN scope s ON s.id = e.program_id
		    JOIN users u ON u.id = e.user_id
		    WHERE e.role = 'student' AND u.is_active = true
		),
		-- Quá hạn: bài đã hết hạn mà học viên chưa có bài nộp nào.
		overdue AS (
		    SELECT l.user_id, l.full_name, l.email, l.program_id, l.slug, l.title,
		           'overdue' AS kind, count(*)::numeric AS value
		    FROM learners l
		    JOIN nodes n ON n.program_id = l.program_id AND n.kind = 'assignment' AND n.is_published = true
		    JOIN assignments a ON a.node_id = n.id
		    WHERE a.due_at IS NOT NULL AND a.due_at < now()
		      AND NOT EXISTS (SELECT 1 FROM submissions sub
		                      WHERE sub.assignment_id = n.id AND sub.user_id = l.user_id)
		    GROUP BY 1, 2, 3, 4, 5, 6
		),
		-- Điểm thấp: trung bình phần trăm điểm các bài đã nộp dưới 50%.
		low_score AS (
		    SELECT l.user_id, l.full_name, l.email, l.program_id, l.slug, l.title,
		           'low_score' AS kind,
		           round(avg(100.0 * (sub.auto_score + COALESCE(sub.manual_score, 0))
		                     / NULLIF(sub.max_score, 0)), 1) AS value
		    FROM learners l
		    JOIN nodes n ON n.program_id = l.program_id
		    JOIN submissions sub ON sub.assignment_id = n.id AND sub.user_id = l.user_id
		    WHERE sub.status = 'graded' AND sub.max_score > 0
		    GROUP BY 1, 2, 3, 4, 5, 6
		    HAVING avg(100.0 * (sub.auto_score + COALESCE(sub.manual_score, 0))
		               / NULLIF(sub.max_score, 0)) < 50
		)
		SELECT * FROM overdue
		UNION ALL
		SELECT * FROM low_score
		ORDER BY kind, value DESC
		LIMIT 50`, userID, allPrograms)
	if err != nil {
		return nil, translate(err, "đọc cảnh báo học tập")
	}
	defer rows.Close()

	out := []*LearningAlert{}
	for rows.Next() {
		var a LearningAlert
		if err := rows.Scan(&a.UserID, &a.FullName, &a.Email, &a.ProgramID,
			&a.ProgramSlug, &a.ProgramName, &a.Kind, &a.Value); err != nil {
			return nil, translate(err, "đọc cảnh báo học tập")
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}
