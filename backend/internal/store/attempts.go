package store

import (
	"context"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/manhnv/elearning/backend/internal/models"
)

// Cho phép trễ mạng khi học viên bấm nộp đúng lúc chuông reo.
const submitGracePeriod = 45 * time.Second

// AttemptSession là một lượt làm bài đang mở của học viên.
type AttemptSession struct {
	ID        uuid.UUID  `json:"id"`
	StartedAt time.Time  `json:"startedAt"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

// Expired cho biết phiên đã quá hạn nộp (đã tính khoảng trễ cho phép).
func (s *AttemptSession) Expired(now time.Time) bool {
	return s.ExpiresAt != nil && now.After(s.ExpiresAt.Add(submitGracePeriod))
}

// StartAttempt mở một lượt làm bài mới, hoặc trả lại lượt đang dang dở.
// Đồng hồ tiếp tục chạy khi học viên tải lại trang: mốc bắt đầu chỉ ghi một lần.
func (s *Store) StartAttempt(ctx context.Context, assignmentID, userID uuid.UUID) (*AttemptSession, error) {
	var timeLimit, maxAttempts int
	err := s.pool.QueryRow(ctx, `
		SELECT time_limit_minutes, max_attempts FROM assignments WHERE node_id = $1`, assignmentID).
		Scan(&timeLimit, &maxAttempts)
	if err != nil {
		return nil, translate(err, "đọc cấu hình bài tập")
	}

	if existing, err := s.OpenAttempt(ctx, assignmentID, userID); err != nil {
		return nil, err
	} else if existing != nil && !existing.Expired(time.Now()) {
		return existing, nil
	}

	// Số lượt tính theo bài đã nộp, nên phiên bỏ dở không làm mất lượt của học viên.
	submitted, err := s.CountAttempts(ctx, assignmentID, userID)
	if err != nil {
		return nil, err
	}
	if maxAttempts > 0 && submitted >= maxAttempts {
		return nil, Invalidf("Bạn đã dùng hết %d lượt làm bài cho phép", maxAttempts)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Dọn phiên cũ đã quá hạn để không vướng ràng buộc "mỗi bài một phiên đang mở".
	if _, err := tx.Exec(ctx, `
		DELETE FROM attempt_sessions
		WHERE assignment_id = $1 AND user_id = $2 AND submitted_at IS NULL`, assignmentID, userID); err != nil {
		return nil, translate(err, "dọn phiên làm bài cũ")
	}

	var session AttemptSession
	err = tx.QueryRow(ctx, `
		INSERT INTO attempt_sessions (assignment_id, user_id, expires_at)
		VALUES ($1, $2, CASE WHEN $3 > 0 THEN now() + make_interval(mins => $3) ELSE NULL END)
		RETURNING id, started_at, expires_at`, assignmentID, userID, timeLimit).
		Scan(&session.ID, &session.StartedAt, &session.ExpiresAt)
	if err != nil {
		return nil, translate(err, "mở lượt làm bài")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &session, nil
}

// OpenAttempt trả về lượt làm bài chưa nộp, hoặc nil nếu không có.
func (s *Store) OpenAttempt(ctx context.Context, assignmentID, userID uuid.UUID) (*AttemptSession, error) {
	var session AttemptSession
	err := s.pool.QueryRow(ctx, `
		SELECT id, started_at, expires_at FROM attempt_sessions
		WHERE assignment_id = $1 AND user_id = $2 AND submitted_at IS NULL`, assignmentID, userID).
		Scan(&session.ID, &session.StartedAt, &session.ExpiresAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, translate(err, "đọc lượt làm bài")
	}
	return &session, nil
}

// ShuffleForSession xáo thứ tự câu hỏi và phương án theo một hạt giống cố định
// lấy từ ID phiên, để học viên tải lại trang vẫn thấy đúng thứ tự cũ.
func ShuffleForSession(questions []*models.Question, sessionID uuid.UUID) {
	seed := int64(binary.BigEndian.Uint64(sessionID[:8]))
	rng := rand.New(rand.NewSource(seed))

	rng.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	for _, q := range questions {
		if q.Type == models.QuestionEssay {
			continue
		}
		rng.Shuffle(len(q.Options), func(i, j int) {
			q.Options[i], q.Options[j] = q.Options[j], q.Options[i]
		})
	}
}
