package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GradebookAssignment là một cột trên bảng điểm — một bài tập của chương trình,
// kèm số liệu tổng hợp của cả cột (tính ở Go sau khi ghép ô điểm).
type GradebookAssignment struct {
	NodeID        uuid.UUID `json:"nodeId"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	QuestionCount int       `json:"questionCount"`
	MaxScore      float64   `json:"maxScore"`  // tổng điểm các câu hỏi hiện có
	PassScore     float64   `json:"passScore"` // 0 = không đặt ngưỡng

	SubmittedCount int     `json:"submittedCount"` // số học sinh đã nộp ít nhất một lượt
	PendingCount   int     `json:"pendingCount"`   // số học sinh còn lượt chờ chấm
	PassedCount    int     `json:"passedCount"`
	AverageScore   float64 `json:"averageScore"` // trung bình điểm cao nhất của những người đã nộp
}

// GradebookCell là điểm của một học sinh ở một bài tập: lấy **lượt điểm cao nhất**
// giống cách học sinh nhìn thấy điểm của mình ở màn hình làm bài.
type GradebookCell struct {
	AssignmentID uuid.UUID `json:"assignmentId"`
	SubmissionID uuid.UUID `json:"submissionId"`
	Score        float64   `json:"score"`
	MaxScore     float64   `json:"maxScore"` // thang điểm lúc nộp, có thể khác thang hiện tại
	Attempts     int       `json:"attempts"`
	Pending      int       `json:"pending"` // số lượt còn chờ chấm tay
	SubmittedAt  time.Time `json:"submittedAt"`
}

// GradebookRow là một học sinh trên bảng điểm.
type GradebookRow struct {
	UserID   uuid.UUID `json:"userId"`
	FullName string    `json:"fullName"`
	Email    string    `json:"email"`
	// Đã bị gỡ ghi danh nhưng vẫn còn bài nộp cũ thì enrolled = false.
	Enrolled bool `json:"enrolled"`

	Cells map[string]*GradebookCell `json:"cells"` // khoá là ID bài tập

	DoneCount    int     `json:"doneCount"`    // số bài tập đã nộp
	PendingCount int     `json:"pendingCount"` // số bài tập còn lượt chờ chấm
	TotalScore   float64 `json:"totalScore"`
}

// ProgramGradebook là bảng điểm bài tập của cả chương trình.
type ProgramGradebook struct {
	Assignments []*GradebookAssignment `json:"assignments"`
	Rows        []*GradebookRow        `json:"rows"`

	TotalMaxScore  float64 `json:"totalMaxScore"` // tổng thang điểm của mọi bài tập
	AverageScore   float64 `json:"averageScore"`  // trung bình điểm tổng của học sinh đã nộp ít nhất một bài
	PendingCount   int     `json:"pendingCount"`  // tổng số ô còn lượt chờ chấm
	SubmittedCells int     `json:"submittedCells"`
	TotalCells     int     `json:"totalCells"` // số học sinh × số bài tập
}

// ProgramGradebook tổng hợp điểm bài tập của toàn bộ học sinh trong một chương trình.
func (s *Store) ProgramGradebook(ctx context.Context, programID uuid.UUID) (*ProgramGradebook, error) {
	out := &ProgramGradebook{Assignments: []*GradebookAssignment{}, Rows: []*GradebookRow{}}

	// Cột: các bài tập theo đúng thứ tự trên cây nội dung (duyệt đệ quy theo position).
	rows, err := s.pool.Query(ctx, `
		WITH RECURSIVE walk AS (
			SELECT n.id, n.slug, n.title, ARRAY[n.position] AS path
			FROM nodes n WHERE n.program_id = $1 AND n.parent_id IS NULL
		UNION ALL
			SELECT c.id, c.slug, c.title, w.path || c.position
			FROM nodes c JOIN walk w ON c.parent_id = w.id
		)
		SELECT w.id, w.slug, w.title, a.pass_score,
		       COALESCE(sum(q.points), 0), count(q.id)
		FROM walk w
		JOIN assignments a ON a.node_id = w.id
		LEFT JOIN questions q ON q.assignment_id = w.id
		GROUP BY w.id, w.slug, w.title, w.path, a.pass_score
		ORDER BY w.path`, programID)
	if err != nil {
		return nil, translate(err, "liệt kê bài tập")
	}
	defer rows.Close()

	byAssignment := map[uuid.UUID]*GradebookAssignment{}
	for rows.Next() {
		var a GradebookAssignment
		if err := rows.Scan(&a.NodeID, &a.Slug, &a.Title, &a.PassScore, &a.MaxScore, &a.QuestionCount); err != nil {
			return nil, translate(err, "đọc bài tập")
		}
		out.Assignments = append(out.Assignments, &a)
		byAssignment[a.NodeID] = &a
		out.TotalMaxScore += a.MaxScore
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Dòng: học sinh đang ghi danh, cộng thêm người đã nộp bài nhưng bị gỡ ghi danh
	// (giữ lại điểm cũ thay vì để nó biến mất khỏi bảng). Giáo viên không tính vào đây.
	userRows, err := s.pool.Query(ctx, `
		SELECT u.id, u.full_name, u.email, (e.user_id IS NOT NULL)
		FROM users u
		LEFT JOIN enrollments e
		       ON e.user_id = u.id AND e.program_id = $1 AND e.role = 'student'
		WHERE (e.user_id IS NOT NULL
		       OR EXISTS (SELECT 1 FROM submissions s
		                  JOIN nodes n ON n.id = s.assignment_id
		                  WHERE s.user_id = u.id AND n.program_id = $1))
		  AND NOT EXISTS (SELECT 1 FROM enrollments t
		                  WHERE t.program_id = $1 AND t.user_id = u.id AND t.role = 'trainer')
		ORDER BY u.full_name, u.email`, programID)
	if err != nil {
		return nil, translate(err, "liệt kê học sinh")
	}
	defer userRows.Close()

	byUser := map[uuid.UUID]*GradebookRow{}
	for userRows.Next() {
		row := &GradebookRow{Cells: map[string]*GradebookCell{}}
		if err := userRows.Scan(&row.UserID, &row.FullName, &row.Email, &row.Enrolled); err != nil {
			return nil, translate(err, "đọc học sinh")
		}
		out.Rows = append(out.Rows, row)
		byUser[row.UserID] = row
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	// Ô điểm: lượt cao điểm nhất của mỗi cặp (học sinh, bài tập), kèm số lượt và số lượt chờ chấm.
	cellRows, err := s.pool.Query(ctx, `
		WITH subs AS (
			SELECT s.* FROM submissions s
			JOIN nodes n ON n.id = s.assignment_id
			WHERE n.program_id = $1
		), agg AS (
			SELECT user_id, assignment_id, count(*) AS attempts,
			       count(*) FILTER (WHERE status = 'submitted') AS pending
			FROM subs GROUP BY user_id, assignment_id
		), best AS (
			SELECT DISTINCT ON (user_id, assignment_id)
			       user_id, assignment_id, id,
			       auto_score + COALESCE(manual_score, 0) AS score,
			       max_score, submitted_at
			FROM subs
			ORDER BY user_id, assignment_id,
			         (auto_score + COALESCE(manual_score, 0)) DESC, submitted_at DESC
		)
		SELECT b.user_id, b.assignment_id, b.id, b.score, b.max_score, b.submitted_at,
		       a.attempts, a.pending
		FROM best b JOIN agg a ON a.user_id = b.user_id AND a.assignment_id = b.assignment_id`,
		programID)
	if err != nil {
		return nil, translate(err, "tổng hợp điểm bài tập")
	}
	defer cellRows.Close()

	for cellRows.Next() {
		var userID uuid.UUID
		var c GradebookCell
		if err := cellRows.Scan(&userID, &c.AssignmentID, &c.SubmissionID, &c.Score, &c.MaxScore,
			&c.SubmittedAt, &c.Attempts, &c.Pending); err != nil {
			return nil, translate(err, "đọc điểm bài tập")
		}
		row, ok := byUser[userID]
		if !ok {
			// Người nộp bài nhưng đang là giáo viên của chương trình — không có dòng trên bảng.
			continue
		}
		row.Cells[c.AssignmentID.String()] = &c

		row.DoneCount++
		row.TotalScore += c.Score
		if c.Pending > 0 {
			row.PendingCount++
		}

		col, ok := byAssignment[c.AssignmentID]
		if !ok {
			continue
		}
		col.SubmittedCount++
		col.AverageScore += c.Score // cộng dồn, chia ở dưới
		if c.Pending > 0 {
			col.PendingCount++
		}
		if col.PassScore > 0 && c.Pending == 0 && c.Score >= col.PassScore {
			col.PassedCount++
		}
	}
	if err := cellRows.Err(); err != nil {
		return nil, err
	}

	for _, col := range out.Assignments {
		if col.SubmittedCount > 0 {
			col.AverageScore /= float64(col.SubmittedCount)
		}
		out.SubmittedCells += col.SubmittedCount
		out.PendingCount += col.PendingCount
	}
	out.TotalCells = len(out.Rows) * len(out.Assignments)

	active := 0
	for _, row := range out.Rows {
		if row.DoneCount > 0 {
			active++
			out.AverageScore += row.TotalScore
		}
	}
	if active > 0 {
		out.AverageScore /= float64(active)
	}
	return out, nil
}
