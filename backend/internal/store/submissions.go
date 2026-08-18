package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// AnswerInput là một câu trả lời do học viên gửi lên.
type AnswerInput struct {
	QuestionID        uuid.UUID
	SelectedOptionIDs []uuid.UUID
	EssayText         string
}

// SubmitAssignment lưu bài làm, tự động chấm phần trắc nghiệm và để phần tự luận chờ giảng viên.
func (s *Store) SubmitAssignment(ctx context.Context, assignmentID, userID uuid.UUID, answers []AnswerInput) (*models.Submission, error) {
	questions, err := s.ListQuestions(ctx, assignmentID)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, Invalidf("Bài tập này chưa có câu hỏi nào")
	}

	var maxAttempts int
	if err := s.pool.QueryRow(ctx, `SELECT max_attempts FROM assignments WHERE node_id = $1`, assignmentID).
		Scan(&maxAttempts); err != nil {
		return nil, translate(err, "đọc cấu hình bài tập")
	}

	// Phiên làm bài là bằng chứng học viên đã bấm bắt đầu và chưa quá giờ.
	session, err := s.OpenAttempt(ctx, assignmentID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, Invalidf("Lượt làm bài đã kết thúc. Hãy bấm Bắt đầu làm bài rồi làm lại.")
	}
	if session.Expired(time.Now()) {
		// Xoá phiên quá hạn để học viên bắt đầu lượt mới mà không bị vướng ràng buộc.
		if _, err := s.pool.Exec(ctx, `DELETE FROM attempt_sessions WHERE id = $1`, session.ID); err != nil {
			return nil, translate(err, "dọn phiên quá hạn")
		}
		return nil, Invalidf("Đã hết thời gian làm bài, bài nộp không được ghi nhận.")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var attemptNo int
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(max(attempt_no), 0) + 1 FROM submissions
		WHERE assignment_id = $1 AND user_id = $2`, assignmentID, userID).Scan(&attemptNo); err != nil {
		return nil, translate(err, "đếm lượt làm bài")
	}
	if maxAttempts > 0 && attemptNo > maxAttempts {
		return nil, Invalidf("Bạn đã dùng hết %d lượt làm bài cho phép", maxAttempts)
	}

	answerByQuestion := make(map[uuid.UUID]AnswerInput, len(answers))
	for _, a := range answers {
		answerByQuestion[a.QuestionID] = a
	}

	var maxScore, autoScore float64
	needsGrading := false
	graded := make([]gradedAnswer, 0, len(questions))

	for _, q := range questions {
		maxScore += q.Points
		answer := answerByQuestion[q.ID]

		if q.Type == models.QuestionEssay {
			needsGrading = true
			graded = append(graded, gradedAnswer{questionID: q.ID, essayText: answer.EssayText})
			continue
		}

		correct := isChoiceCorrect(q, answer.SelectedOptionIDs)
		score := 0.0
		if correct {
			score = q.Points
			autoScore += q.Points
		}
		isCorrect := correct
		graded = append(graded, gradedAnswer{
			questionID: q.ID,
			selected:   answer.SelectedOptionIDs,
			isCorrect:  &isCorrect,
			score:      score,
		})
	}

	// Bài chỉ có trắc nghiệm được chấm xong ngay, bài có tự luận chờ giảng viên vào điểm.
	status := "graded"
	var manualScore *float64
	if needsGrading {
		status = "submitted"
	} else {
		zero := 0.0
		manualScore = &zero
	}

	var submissionID uuid.UUID
	if err := tx.QueryRow(ctx, `
		INSERT INTO submissions (assignment_id, user_id, attempt_no, status, auto_score, manual_score, max_score, graded_at, session_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CASE WHEN $4 = 'graded' THEN now() ELSE NULL END, $8)
		RETURNING id`,
		assignmentID, userID, attemptNo, status, autoScore, manualScore, maxScore, session.ID).Scan(&submissionID); err != nil {
		return nil, translate(err, "lưu bài nộp")
	}

	if _, err := tx.Exec(ctx, `
		UPDATE attempt_sessions SET submitted_at = now() WHERE id = $1`, session.ID); err != nil {
		return nil, translate(err, "đóng lượt làm bài")
	}

	for _, g := range graded {
		// Cột selected_option_ids là NOT NULL: slice rỗng phải khác nil để pgx ghi '{}' thay vì NULL.
		if g.selected == nil {
			g.selected = []uuid.UUID{}
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO submission_answers (submission_id, question_id, selected_option_ids, essay_text, is_correct, score)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			submissionID, g.questionID, g.selected, g.essayText, g.isCorrect, g.score); err != nil {
			return nil, translate(err, "lưu câu trả lời")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetSubmission(ctx, submissionID)
}

type gradedAnswer struct {
	questionID uuid.UUID
	selected   []uuid.UUID
	essayText  string
	isCorrect  *bool
	score      float64
}

// isChoiceCorrect chấm trắc nghiệm theo nguyên tắc đúng trọn vẹn:
// tập đáp án chọn phải trùng khít tập đáp án đúng.
func isChoiceCorrect(q *models.Question, selected []uuid.UUID) bool {
	correct := map[uuid.UUID]bool{}
	for _, o := range q.Options {
		if o.IsCorrect {
			correct[o.ID] = true
		}
	}
	if len(correct) == 0 {
		return false
	}

	chosen := map[uuid.UUID]bool{}
	for _, id := range selected {
		chosen[id] = true
	}
	if len(chosen) != len(correct) {
		return false
	}
	for id := range correct {
		if !chosen[id] {
			return false
		}
	}
	return true
}

const submissionColumns = `
	s.id, s.assignment_id, s.user_id, s.attempt_no, s.status, s.auto_score, s.manual_score,
	s.max_score, s.feedback, s.graded_by, s.graded_at, s.submitted_at`

type submissionScanner interface {
	Scan(dest ...any) error
}

func scanSubmission(row submissionScanner, withUser bool) (*models.Submission, error) {
	var s models.Submission
	dest := []any{&s.ID, &s.AssignmentID, &s.UserID, &s.AttemptNo, &s.Status, &s.AutoScore,
		&s.ManualScore, &s.MaxScore, &s.Feedback, &s.GradedBy, &s.GradedAt, &s.SubmittedAt}
	if withUser {
		dest = append(dest, &s.StudentName, &s.StudentEmail, &s.AssignmentTitle, &s.ProgramTitle)
	}
	if err := row.Scan(dest...); err != nil {
		return nil, err
	}
	s.NeedsGrading = s.Status == "submitted"
	return &s, nil
}

func (s *Store) GetSubmission(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+submissionColumns+`, u.full_name, u.email, n.title, p.title
		FROM submissions s
		JOIN users u ON u.id = s.user_id
		JOIN nodes n ON n.id = s.assignment_id
		JOIN programs p ON p.id = n.program_id
		WHERE s.id = $1`, id)
	sub, err := scanSubmission(row, true)
	if err != nil {
		return nil, translate(err, "đọc bài nộp")
	}

	answers, err := s.listAnswers(ctx, id)
	if err != nil {
		return nil, err
	}
	sub.Answers = answers
	return sub, nil
}

func (s *Store) listAnswers(ctx context.Context, submissionID uuid.UUID) ([]*models.SubmissionAnswer, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT sa.id, sa.question_id, sa.selected_option_ids, sa.essay_text, sa.is_correct, sa.score, sa.comment
		FROM submission_answers sa
		JOIN questions q ON q.id = sa.question_id
		WHERE sa.submission_id = $1
		ORDER BY q.position, q.created_at`, submissionID)
	if err != nil {
		return nil, translate(err, "đọc câu trả lời")
	}
	defer rows.Close()

	out := []*models.SubmissionAnswer{}
	for rows.Next() {
		var a models.SubmissionAnswer
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.SelectedOptionIDs, &a.EssayText,
			&a.IsCorrect, &a.Score, &a.Comment); err != nil {
			return nil, translate(err, "đọc câu trả lời")
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

type ListSubmissionsFilter struct {
	AssignmentID uuid.UUID
	ProgramID    uuid.UUID
	UserID       uuid.UUID
	OnlyPending  bool
}

func (s *Store) ListSubmissions(ctx context.Context, f ListSubmissionsFilter) ([]*models.Submission, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+submissionColumns+`, u.full_name, u.email, n.title, p.title
		FROM submissions s
		JOIN users u ON u.id = s.user_id
		JOIN nodes n ON n.id = s.assignment_id
		JOIN programs p ON p.id = n.program_id
		WHERE ($1::uuid IS NULL OR s.assignment_id = $1::uuid)
		  AND ($2::uuid IS NULL OR n.program_id = $2::uuid)
		  AND ($3::uuid IS NULL OR s.user_id = $3::uuid)
		  AND ($4 = false OR s.status = 'submitted')
		ORDER BY s.submitted_at DESC
		LIMIT 200`,
		nullableUUID(f.AssignmentID), nullableUUID(f.ProgramID), nullableUUID(f.UserID), f.OnlyPending)
	if err != nil {
		return nil, translate(err, "liệt kê bài nộp")
	}
	defer rows.Close()

	out := []*models.Submission{}
	for rows.Next() {
		sub, err := scanSubmission(rows, true)
		if err != nil {
			return nil, translate(err, "đọc bài nộp")
		}
		out = append(out, sub)
	}
	return out, rows.Err()
}

// GradeAnswerInput là điểm và nhận xét giảng viên chấm cho một câu tự luận.
type GradeAnswerInput struct {
	AnswerID uuid.UUID
	Score    float64
	Comment  string
}

// GradeSubmission ghi điểm phần tự luận và chốt trạng thái bài nộp.
func (s *Store) GradeSubmission(ctx context.Context, submissionID, graderID uuid.UUID, feedback string, answers []GradeAnswerInput) (*models.Submission, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	for _, a := range answers {
		if _, err := tx.Exec(ctx, `
			UPDATE submission_answers SET score = $3, comment = $4
			WHERE id = $1 AND submission_id = $2`, a.AnswerID, submissionID, a.Score, a.Comment); err != nil {
			return nil, translate(err, "ghi điểm câu trả lời")
		}
	}

	// Điểm chấm tay là tổng điểm của các câu tự luận (câu trắc nghiệm đã tính vào auto_score).
	var manualScore float64
	if err := tx.QueryRow(ctx, `
		SELECT COALESCE(sum(sa.score), 0)
		FROM submission_answers sa
		JOIN questions q ON q.id = sa.question_id
		WHERE sa.submission_id = $1 AND q.type = 'essay'`, submissionID).Scan(&manualScore); err != nil {
		return nil, translate(err, "tính điểm tự luận")
	}

	tag, err := tx.Exec(ctx, `
		UPDATE submissions
		SET manual_score = $2, feedback = $3, status = 'graded', graded_by = $4, graded_at = now()
		WHERE id = $1`, submissionID, manualScore, feedback, graderID)
	if err != nil {
		return nil, translate(err, "chốt điểm bài nộp")
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetSubmission(ctx, submissionID)
}

// CountAttempts đếm số lần một học viên đã nộp một bài tập.
func (s *Store) CountAttempts(ctx context.Context, assignmentID, userID uuid.UUID) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM submissions WHERE assignment_id = $1 AND user_id = $2`,
		assignmentID, userID).Scan(&n)
	if err != nil {
		return 0, translate(err, "đếm lượt làm bài")
	}
	return n, nil
}
