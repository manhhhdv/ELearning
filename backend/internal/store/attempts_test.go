package store

import (
	"testing"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

func sampleQuestions() []*models.Question {
	questions := make([]*models.Question, 0, 6)
	for i := 0; i < 6; i++ {
		q := &models.Question{
			ID:   uuid.New(),
			Type: models.QuestionSingleChoice,
		}
		for j := 0; j < 4; j++ {
			q.Options = append(q.Options, &models.QuestionOption{ID: uuid.New()})
		}
		questions = append(questions, q)
	}
	return questions
}

func order(questions []*models.Question) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(questions))
	for _, q := range questions {
		ids = append(ids, q.ID)
		for _, o := range q.Options {
			ids = append(ids, o.ID)
		}
	}
	return ids
}

func equal(a, b []uuid.UUID) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Cùng một phiên phải cho ra đúng thứ tự cũ, để học viên tải lại trang không bị xáo đề.
func TestShuffleForSessionOnDinh(t *testing.T) {
	sessionID := uuid.New()

	first := sampleQuestions()
	second := make([]*models.Question, len(first))
	copy(second, first)
	// Nhân bản sâu để hai lần trộn không dùng chung slice phương án.
	for i, q := range first {
		clone := *q
		clone.Options = append([]*models.QuestionOption(nil), q.Options...)
		second[i] = &clone
	}

	ShuffleForSession(first, sessionID)
	ShuffleForSession(second, sessionID)

	if !equal(order(first), order(second)) {
		t.Error("cùng một phiên nhưng thứ tự trộn lại khác nhau")
	}
}

// Hai lượt làm bài khác nhau thì thứ tự phải khác nhau.
func TestShuffleForSessionKhacPhien(t *testing.T) {
	base := sampleQuestions()

	a := make([]*models.Question, len(base))
	b := make([]*models.Question, len(base))
	for i, q := range base {
		ca, cb := *q, *q
		ca.Options = append([]*models.QuestionOption(nil), q.Options...)
		cb.Options = append([]*models.QuestionOption(nil), q.Options...)
		a[i], b[i] = &ca, &cb
	}

	ShuffleForSession(a, uuid.New())
	ShuffleForSession(b, uuid.New())

	if equal(order(a), order(b)) {
		t.Error("hai phiên khác nhau lại cho cùng một thứ tự")
	}
}

// Trộn không được làm mất hay nhân bản câu hỏi.
func TestShuffleGiuNguyenSoLuong(t *testing.T) {
	questions := sampleQuestions()
	before := map[uuid.UUID]int{}
	for _, q := range questions {
		before[q.ID] = len(q.Options)
	}

	ShuffleForSession(questions, uuid.New())

	if len(questions) != len(before) {
		t.Fatalf("số câu hỏi đổi từ %d thành %d", len(before), len(questions))
	}
	for _, q := range questions {
		want, ok := before[q.ID]
		if !ok {
			t.Fatalf("xuất hiện câu hỏi lạ %s", q.ID)
		}
		if len(q.Options) != want {
			t.Errorf("câu %s có %d phương án, mong đợi %d", q.ID, len(q.Options), want)
		}
	}
}
