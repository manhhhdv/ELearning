package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/manhnv/elearning/backend/internal/models"
)

// getAssignment đọc cấu hình bài tập kèm toàn bộ câu hỏi và đáp án.
func (s *Store) getAssignment(ctx context.Context, nodeID uuid.UUID) (*models.Assignment, error) {
	var a models.Assignment
	err := s.pool.QueryRow(ctx, `
		SELECT instructions, time_limit_minutes, max_attempts, pass_score, shuffle_questions, due_at
		FROM assignments WHERE node_id = $1`, nodeID).
		Scan(&a.Instructions, &a.TimeLimitMinutes, &a.MaxAttempts, &a.PassScore, &a.ShuffleQuestions, &a.DueAt)
	if err != nil {
		return nil, translate(err, "đọc bài tập")
	}

	questions, err := s.ListQuestions(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	a.Questions = questions
	a.QuestionCount = len(questions)
	return &a, nil
}

// ListQuestions trả về câu hỏi của một bài tập theo đúng thứ tự đã sắp.
func (s *Store) ListQuestions(ctx context.Context, assignmentID uuid.UUID) ([]*models.Question, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, code, type, prompt, points, position, explanation
		FROM questions WHERE assignment_id = $1 ORDER BY position, created_at`, assignmentID)
	if err != nil {
		return nil, translate(err, "liệt kê câu hỏi")
	}
	defer rows.Close()

	questions := []*models.Question{}
	byID := map[uuid.UUID]*models.Question{}
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.Code, &q.Type, &q.Prompt, &q.Points, &q.Position, &q.Explanation); err != nil {
			return nil, translate(err, "đọc câu hỏi")
		}
		q.Options = []*models.QuestionOption{}
		questions = append(questions, &q)
		byID[q.ID] = &q
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return questions, nil
	}

	optRows, err := s.pool.Query(ctx, `
		SELECT o.id, o.question_id, o.content, o.is_correct, o.position
		FROM question_options o
		JOIN questions q ON q.id = o.question_id
		WHERE q.assignment_id = $1
		ORDER BY o.position, o.id`, assignmentID)
	if err != nil {
		return nil, translate(err, "đọc đáp án")
	}
	defer optRows.Close()

	for optRows.Next() {
		var questionID uuid.UUID
		var o models.QuestionOption
		if err := optRows.Scan(&o.ID, &questionID, &o.Content, &o.IsCorrect, &o.Position); err != nil {
			return nil, translate(err, "đọc đáp án")
		}
		if q, ok := byID[questionID]; ok {
			q.Options = append(q.Options, &o)
		}
	}
	return questions, optRows.Err()
}

// SaveQuestionParams mô tả một câu hỏi được lưu từ trình soạn thảo.
// ID rỗng nghĩa là thêm mới; danh sách đáp án luôn được ghi đè toàn bộ.
type SaveQuestionParams struct {
	AssignmentID uuid.UUID
	// Để trống thì hệ thống tự sinh mã kế tiếp dạng C01, C02…
	Code        string
	Type        string
	Prompt      string
	Points      float64
	Explanation string
	Options     []QuestionOptionInput
}

type QuestionOptionInput struct {
	Content   string
	IsCorrect bool
}

func (s *Store) CreateQuestion(ctx context.Context, p SaveQuestionParams) (*models.Question, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var position int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(max(position) + 1, 0) FROM questions WHERE assignment_id = $1`,
		p.AssignmentID).Scan(&position); err != nil {
		return nil, translate(err, "tính thứ tự câu hỏi")
	}

	code := p.Code
	if code == "" {
		if code, err = nextQuestionCode(ctx, tx, p.AssignmentID); err != nil {
			return nil, err
		}
	}

	var id uuid.UUID
	if err := tx.QueryRow(ctx, `
		INSERT INTO questions (assignment_id, code, type, prompt, points, position, explanation)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		p.AssignmentID, code, p.Type, p.Prompt, p.Points, position, p.Explanation).Scan(&id); err != nil {
		return nil, translate(err, "tạo câu hỏi")
	}
	if err := replaceOptions(ctx, tx, id, p.Type, p.Options); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetQuestion(ctx, id)
}

func (s *Store) UpdateQuestion(ctx context.Context, id uuid.UUID, p SaveQuestionParams) (*models.Question, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Mã để trống nghĩa là giữ nguyên mã cũ, không tự sinh lại.
	tag, err := tx.Exec(ctx, `
		UPDATE questions
		SET type = $2, prompt = $3, points = $4, explanation = $5,
		    code = CASE WHEN $6 = '' THEN code ELSE $6 END
		WHERE id = $1`,
		id, p.Type, p.Prompt, p.Points, p.Explanation, p.Code)
	if err != nil {
		return nil, translate(err, "cập nhật câu hỏi")
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	if err := replaceOptions(ctx, tx, id, p.Type, p.Options); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetQuestion(ctx, id)
}

// replaceOptions ghi đè toàn bộ đáp án của một câu hỏi.
// Câu tự luận không có đáp án nên chỉ cần xoá sạch danh sách cũ.
func replaceOptions(ctx context.Context, tx pgx.Tx, questionID uuid.UUID, qType string, opts []QuestionOptionInput) error {
	if _, err := tx.Exec(ctx, `DELETE FROM question_options WHERE question_id = $1`, questionID); err != nil {
		return translate(err, "xoá đáp án cũ")
	}
	if qType == models.QuestionEssay {
		return nil
	}
	for i, o := range opts {
		if _, err := tx.Exec(ctx, `
			INSERT INTO question_options (question_id, content, is_correct, position)
			VALUES ($1, $2, $3, $4)`, questionID, o.Content, o.IsCorrect, i); err != nil {
			return translate(err, "lưu đáp án")
		}
	}
	return nil
}

func (s *Store) GetQuestion(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	var q models.Question
	err := s.pool.QueryRow(ctx, `
		SELECT id, code, type, prompt, points, position, explanation FROM questions WHERE id = $1`, id).
		Scan(&q.ID, &q.Code, &q.Type, &q.Prompt, &q.Points, &q.Position, &q.Explanation)
	if err != nil {
		return nil, translate(err, "đọc câu hỏi")
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, content, is_correct, position FROM question_options
		WHERE question_id = $1 ORDER BY position, id`, id)
	if err != nil {
		return nil, translate(err, "đọc đáp án")
	}
	defer rows.Close()

	q.Options = []*models.QuestionOption{}
	for rows.Next() {
		var o models.QuestionOption
		if err := rows.Scan(&o.ID, &o.Content, &o.IsCorrect, &o.Position); err != nil {
			return nil, translate(err, "đọc đáp án")
		}
		q.Options = append(q.Options, &o)
	}
	return &q, rows.Err()
}

func (s *Store) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM questions WHERE id = $1`, id)
	if err != nil {
		return translate(err, "xoá câu hỏi")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ReorderQuestions ghi lại thứ tự câu hỏi theo đúng danh sách ID nhận được.
func (s *Store) ReorderQuestions(ctx context.Context, assignmentID uuid.UUID, ids []uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, id := range ids {
		if _, err := tx.Exec(ctx, `
			UPDATE questions SET position = $3 WHERE id = $1 AND assignment_id = $2`, id, assignmentID, i); err != nil {
			return translate(err, "sắp xếp câu hỏi")
		}
	}
	return tx.Commit(ctx)
}

// nextQuestionCode sinh mã kế tiếp cho một bài tập: C01, C02, … bỏ qua các mã đã dùng.
func nextQuestionCode(ctx context.Context, tx pgx.Tx, assignmentID uuid.UUID) (string, error) {
	used := map[string]bool{}
	rows, err := tx.Query(ctx, `SELECT lower(code) FROM questions WHERE assignment_id = $1`, assignmentID)
	if err != nil {
		return "", translate(err, "đọc mã câu hỏi")
	}
	defer rows.Close()
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return "", translate(err, "đọc mã câu hỏi")
		}
		used[c] = true
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("C%02d", i)
		if !used[strings.ToLower(candidate)] {
			return candidate, nil
		}
	}
}

// ImportQuestions thêm nhiều câu hỏi trong một transaction: hỏng một câu thì huỷ toàn bộ,
// tránh để bài tập rơi vào trạng thái nhập dở.
func (s *Store) ImportQuestions(ctx context.Context, assignmentID uuid.UUID, items []SaveQuestionParams) (int, error) {
	if len(items) == 0 {
		return 0, Invalidf("Không có câu hỏi nào để nhập")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var position int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(max(position) + 1, 0) FROM questions WHERE assignment_id = $1`,
		assignmentID).Scan(&position); err != nil {
		return 0, translate(err, "tính thứ tự câu hỏi")
	}

	for i, item := range items {
		code := item.Code
		if code == "" {
			if code, err = nextQuestionCode(ctx, tx, assignmentID); err != nil {
				return 0, err
			}
		}

		var id uuid.UUID
		err := tx.QueryRow(ctx, `
			INSERT INTO questions (assignment_id, code, type, prompt, points, position, explanation)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			assignmentID, code, item.Type, item.Prompt, item.Points, position+i, item.Explanation).Scan(&id)
		if err != nil {
			if e := translate(err, "nhập câu hỏi"); errors.Is(e, ErrConflict) {
				return 0, Invalidf("Mã câu hỏi %q bị trùng, vui lòng sửa lại trước khi nhập", code)
			}
			return 0, translate(err, "nhập câu hỏi")
		}
		if err := replaceOptions(ctx, tx, id, item.Type, item.Options); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(items), nil
}

// QuestionAssignmentID trả về bài tập chứa câu hỏi, dùng để kiểm tra quyền.
func (s *Store) QuestionAssignmentID(ctx context.Context, questionID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT assignment_id FROM questions WHERE id = $1`, questionID).Scan(&id)
	if err != nil {
		return uuid.Nil, translate(err, "đọc bài tập của câu hỏi")
	}
	return id, nil
}
