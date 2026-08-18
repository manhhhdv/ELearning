package store

import (
	"context"

	"github.com/google/uuid"
)

// QuestionStat là kết quả tổng hợp của một câu hỏi trên toàn bộ bài đã nộp.
type QuestionStat struct {
	QuestionID uuid.UUID `json:"questionId"`
	Code       string    `json:"code"`
	Type       string    `json:"type"`
	Prompt     string    `json:"prompt"`
	Points     float64   `json:"points"`
	Position   int       `json:"position"`

	AnswerCount  int     `json:"answerCount"`  // số lượt trả lời câu này
	CorrectCount int     `json:"correctCount"` // số lượt đúng, chỉ có với câu trắc nghiệm
	BlankCount   int     `json:"blankCount"`   // số lượt bỏ trống
	AverageScore float64 `json:"averageScore"`
	// Câu tự luận chưa chấm thì không có tỉ lệ đúng/sai để hiển thị.
	NeedsGrading int `json:"needsGrading"`
}

// AssignmentResults là số liệu tổng hợp một bài tập, phục vụ màn hình kết quả của admin.
type AssignmentResults struct {
	SubmissionCount int             `json:"submissionCount"`
	StudentCount    int             `json:"studentCount"`
	PendingCount    int             `json:"pendingCount"`
	MaxScore        float64         `json:"maxScore"`
	AverageScore    float64         `json:"averageScore"`
	PassScore       float64         `json:"passScore"`
	PassedCount     int             `json:"passedCount"`
	Questions       []*QuestionStat `json:"questions"`
}

// AssignmentResults tổng hợp kết quả của một bài tập theo từng câu hỏi.
func (s *Store) AssignmentResults(ctx context.Context, assignmentID uuid.UUID) (*AssignmentResults, error) {
	out := &AssignmentResults{Questions: []*QuestionStat{}}

	err := s.pool.QueryRow(ctx, `
		SELECT
			count(*),
			count(DISTINCT s.user_id),
			count(*) FILTER (WHERE s.status = 'submitted'),
			COALESCE(max(s.max_score), 0),
			COALESCE(avg(s.auto_score + COALESCE(s.manual_score, 0)), 0)
		FROM submissions s WHERE s.assignment_id = $1`, assignmentID).
		Scan(&out.SubmissionCount, &out.StudentCount, &out.PendingCount, &out.MaxScore, &out.AverageScore)
	if err != nil {
		return nil, translate(err, "tổng hợp bài nộp")
	}

	// Điểm đạt lấy từ cấu hình bài tập; 0 nghĩa là không đặt ngưỡng.
	if err := s.pool.QueryRow(ctx, `
		SELECT pass_score FROM assignments WHERE node_id = $1`, assignmentID).Scan(&out.PassScore); err != nil {
		return nil, translate(err, "đọc cấu hình bài tập")
	}
	if out.PassScore > 0 {
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*) FROM submissions
			WHERE assignment_id = $1 AND status = 'graded'
			  AND auto_score + COALESCE(manual_score, 0) >= $2`,
			assignmentID, out.PassScore).Scan(&out.PassedCount); err != nil {
			return nil, translate(err, "đếm bài đạt")
		}
	}

	rows, err := s.pool.Query(ctx, `
		SELECT
			q.id, q.code, q.type, q.prompt, q.points, q.position,
			count(sa.id),
			count(*) FILTER (WHERE sa.is_correct IS TRUE),
			count(*) FILTER (WHERE sa.selected_option_ids = '{}' AND btrim(sa.essay_text) = ''),
			COALESCE(avg(sa.score), 0),
			count(*) FILTER (WHERE q.type = 'essay' AND sub.status = 'submitted')
		FROM questions q
		LEFT JOIN submission_answers sa ON sa.question_id = q.id
		LEFT JOIN submissions sub ON sub.id = sa.submission_id
		WHERE q.assignment_id = $1
		GROUP BY q.id, q.code, q.type, q.prompt, q.points, q.position
		ORDER BY q.position, q.code`, assignmentID)
	if err != nil {
		return nil, translate(err, "tổng hợp theo câu hỏi")
	}
	defer rows.Close()

	for rows.Next() {
		var st QuestionStat
		if err := rows.Scan(&st.QuestionID, &st.Code, &st.Type, &st.Prompt, &st.Points, &st.Position,
			&st.AnswerCount, &st.CorrectCount, &st.BlankCount, &st.AverageScore, &st.NeedsGrading); err != nil {
			return nil, translate(err, "đọc thống kê câu hỏi")
		}
		out.Questions = append(out.Questions, &st)
	}
	return out, rows.Err()
}
