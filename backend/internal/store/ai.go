package store

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// ---------------------------------------------------------------------------
// Giáo án do AI sinh
// ---------------------------------------------------------------------------

type SaveLessonPlanParams struct {
	OwnerID         uuid.UUID
	ProgramID       *uuid.UUID
	Subject         string
	Grade           string
	Topic           string
	Objectives      string
	DurationMinutes int
	Title           string
	Content         json.RawMessage
}

const lessonPlanColumns = `lp.id, lp.owner_id, lp.program_id, lp.subject, lp.grade, lp.topic,
	lp.objectives, lp.duration_minutes, lp.title, lp.content, lp.created_at, lp.updated_at`

type rowScanner interface{ Scan(dest ...any) error }

func scanLessonPlan(row rowScanner, withOwner bool) (*models.LessonPlan, error) {
	var p models.LessonPlan
	dest := []any{&p.ID, &p.OwnerID, &p.ProgramID, &p.Subject, &p.Grade, &p.Topic,
		&p.Objectives, &p.DurationMinutes, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt}
	if withOwner {
		dest = append(dest, &p.OwnerName)
	}
	if err := row.Scan(dest...); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) CreateLessonPlan(ctx context.Context, p SaveLessonPlanParams) (*models.LessonPlan, error) {
	if len(p.Content) == 0 {
		p.Content = json.RawMessage("{}")
	}
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		INSERT INTO lesson_plans (owner_id, program_id, subject, grade, topic, objectives,
		                          duration_minutes, title, content)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		p.OwnerID, p.ProgramID, p.Subject, p.Grade, p.Topic, p.Objectives,
		p.DurationMinutes, p.Title, []byte(p.Content)).Scan(&id)
	if err != nil {
		return nil, translate(err, "lưu giáo án")
	}
	return s.GetLessonPlan(ctx, id)
}

// UpdateLessonPlan lưu bản giáo viên đã chỉnh sửa. Chỉ chủ sở hữu sửa được.
func (s *Store) UpdateLessonPlan(ctx context.Context, id, ownerID uuid.UUID, title string, content json.RawMessage) (*models.LessonPlan, error) {
	if len(content) == 0 {
		content = json.RawMessage("{}")
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE lesson_plans SET title = $3, content = $4
		WHERE id = $1 AND owner_id = $2`, id, ownerID, title, []byte(content))
	if err != nil {
		return nil, translate(err, "cập nhật giáo án")
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return s.GetLessonPlan(ctx, id)
}

func (s *Store) GetLessonPlan(ctx context.Context, id uuid.UUID) (*models.LessonPlan, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+lessonPlanColumns+`, u.full_name
		FROM lesson_plans lp JOIN users u ON u.id = lp.owner_id
		WHERE lp.id = $1`, id)
	p, err := scanLessonPlan(row, true)
	if err != nil {
		return nil, translate(err, "đọc giáo án")
	}
	return p, nil
}

func (s *Store) ListLessonPlans(ctx context.Context, ownerID uuid.UUID) ([]*models.LessonPlan, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+lessonPlanColumns+`, u.full_name
		FROM lesson_plans lp JOIN users u ON u.id = lp.owner_id
		WHERE lp.owner_id = $1
		ORDER BY lp.created_at DESC
		LIMIT 200`, ownerID)
	if err != nil {
		return nil, translate(err, "liệt kê giáo án")
	}
	defer rows.Close()

	out := []*models.LessonPlan{}
	for rows.Next() {
		p, err := scanLessonPlan(rows, true)
		if err != nil {
			return nil, translate(err, "đọc giáo án")
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) DeleteLessonPlan(ctx context.Context, id, ownerID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM lesson_plans WHERE id = $1 AND owner_id = $2`, id, ownerID)
	if err != nil {
		return translate(err, "xoá giáo án")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// Nhật ký tác vụ AI
// ---------------------------------------------------------------------------

type LogAITaskParams struct {
	UserID   uuid.UUID
	Kind     string
	Title    string
	Status   string
	Detail   string
	Attempts int
}

// LogAITask ghi lại một lần gọi AI. Lỗi khi ghi nhật ký không được làm hỏng
// kết quả đã sinh ra, nên hàm gọi chỉ cần log lại chứ không trả lỗi cho người dùng.
func (s *Store) LogAITask(ctx context.Context, p LogAITaskParams) error {
	if p.Status == "" {
		p.Status = "ok"
	}
	if p.Attempts <= 0 {
		p.Attempts = 1
	}
	// Thông báo lỗi có thể rất dài; cắt bớt để bảng nhật ký không phình.
	if len(p.Detail) > 2000 {
		p.Detail = p.Detail[:2000]
	}
	_, err := s.pool.Exec(ctx, `
		INSERT INTO ai_tasks (user_id, kind, title, status, detail, attempts)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		p.UserID, p.Kind, p.Title, p.Status, p.Detail, p.Attempts)
	return translate(err, "ghi nhật ký tác vụ AI")
}

func (s *Store) ListAITasks(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AITask, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, kind, title, status, detail, attempts, created_at
		FROM ai_tasks WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, translate(err, "đọc nhật ký tác vụ AI")
	}
	defer rows.Close()

	out := []*models.AITask{}
	for rows.Next() {
		var t models.AITask
		if err := rows.Scan(&t.ID, &t.UserID, &t.Kind, &t.Title, &t.Status,
			&t.Detail, &t.Attempts, &t.CreatedAt); err != nil {
			return nil, translate(err, "đọc nhật ký tác vụ AI")
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------------------
// Hội thoại trợ lý AI
// ---------------------------------------------------------------------------

func (s *Store) ListConversations(ctx context.Context, userID uuid.UUID) ([]*models.AIConversation, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, title, created_at, updated_at, program_id, node_id
		FROM ai_conversations WHERE user_id = $1 ORDER BY updated_at DESC LIMIT 100`, userID)
	if err != nil {
		return nil, translate(err, "liệt kê hội thoại")
	}
	defer rows.Close()

	out := []*models.AIConversation{}
	for rows.Next() {
		var c models.AIConversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt, &c.ProgramID, &c.NodeID); err != nil {
			return nil, translate(err, "đọc hội thoại")
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

// GetConversation đọc một hội thoại kèm toàn bộ lượt trao đổi.
// Trả về ErrNotFound nếu hội thoại không thuộc về người dùng này.
func (s *Store) GetConversation(ctx context.Context, id, userID uuid.UUID) (*models.AIConversation, error) {
	var c models.AIConversation
	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, title, created_at, updated_at, program_id, node_id
		FROM ai_conversations WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt, &c.ProgramID, &c.NodeID)
	if err != nil {
		return nil, translate(err, "đọc hội thoại")
	}

	rows, err := s.pool.Query(ctx, `
		SELECT id, role, content, created_at FROM ai_messages
		WHERE conversation_id = $1 ORDER BY created_at, id`, id)
	if err != nil {
		return nil, translate(err, "đọc tin nhắn")
	}
	defer rows.Close()

	c.Messages = []*models.AIMessage{}
	for rows.Next() {
		var m models.AIMessage
		if err := rows.Scan(&m.ID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, translate(err, "đọc tin nhắn")
		}
		c.Messages = append(c.Messages, &m)
	}
	return &c, rows.Err()
}

// CreateConversation mở một cuộc trò chuyện mới, lấy câu hỏi đầu làm tiêu đề.
func (s *Store) CreateConversation(ctx context.Context, userID uuid.UUID, firstMessage string, programID, nodeID *uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		INSERT INTO ai_conversations (user_id, title, program_id, node_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, conversationTitle(firstMessage), programID, nodeID).Scan(&id)
	if err != nil {
		return uuid.Nil, translate(err, "tạo hội thoại")
	}
	return id, nil
}

// AppendMessages ghi cặp câu hỏi – câu trả lời vào hội thoại và đẩy hội thoại
// lên đầu danh sách (trigger cập nhật updated_at).
func (s *Store) AppendMessages(ctx context.Context, conversationID uuid.UUID, question, answer string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range []struct{ role, content string }{{"user", question}, {"model", answer}} {
		if _, err := tx.Exec(ctx, `
			INSERT INTO ai_messages (conversation_id, role, content) VALUES ($1, $2, $3)`,
			conversationID, m.role, m.content); err != nil {
			return translate(err, "lưu tin nhắn")
		}
	}
	if _, err := tx.Exec(ctx, `
		UPDATE ai_conversations SET updated_at = now() WHERE id = $1`, conversationID); err != nil {
		return translate(err, "cập nhật hội thoại")
	}
	return tx.Commit(ctx)
}

func (s *Store) DeleteConversation(ctx context.Context, id, userID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM ai_conversations WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return translate(err, "xoá hội thoại")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// conversationTitle rút gọn câu hỏi đầu tiên thành tiêu đề hiển thị ở danh sách.
func conversationTitle(msg string) string {
	t := strings.Join(strings.Fields(msg), " ")
	runes := []rune(t)
	if len(runes) > 60 {
		return string(runes[:60]) + "…"
	}
	if t == "" {
		return "Cuộc trò chuyện mới"
	}
	return t
}

// ---------------------------------------------------------------------------
// Điểm tự luận do AI gợi ý
// ---------------------------------------------------------------------------

// EssayAnswerToGrade gom đủ dữ liệu để lớp Service AI chấm một câu tự luận.
type EssayAnswerToGrade struct {
	AnswerID     uuid.UUID
	QuestionID   uuid.UUID
	Prompt       string
	SampleAnswer string
	Rubric       string
	Points       float64
	EssayText    string
}

// ListEssayAnswers trả về các câu tự luận của một bài nộp để AI gợi ý điểm.
func (s *Store) ListEssayAnswers(ctx context.Context, submissionID uuid.UUID) ([]EssayAnswerToGrade, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT sa.id, q.id, q.prompt, q.sample_answer, q.rubric, q.points, sa.essay_text
		FROM submission_answers sa
		JOIN questions q ON q.id = sa.question_id
		WHERE sa.submission_id = $1 AND q.type = 'essay'
		ORDER BY q.position, q.created_at`, submissionID)
	if err != nil {
		return nil, translate(err, "đọc câu tự luận")
	}
	defer rows.Close()

	out := []EssayAnswerToGrade{}
	for rows.Next() {
		var a EssayAnswerToGrade
		if err := rows.Scan(&a.AnswerID, &a.QuestionID, &a.Prompt, &a.SampleAnswer,
			&a.Rubric, &a.Points, &a.EssayText); err != nil {
			return nil, translate(err, "đọc câu tự luận")
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// SaveAISuggestion lưu điểm và nhận xét AI gợi ý. Cột score không bị đụng tới:
// điểm chính thức chỉ thay đổi khi giáo viên bấm lưu ở màn hình chấm bài.
func (s *Store) SaveAISuggestion(ctx context.Context, answerID uuid.UUID, score float64, comment string) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE submission_answers SET ai_score = $2, ai_comment = $3 WHERE id = $1`,
		answerID, score, comment)
	if err != nil {
		return translate(err, "lưu gợi ý chấm điểm")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
