package api

import (
	"testing"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

func essayQuestion() *models.Question {
	return &models.Question{
		ID: uuid.New(), Type: models.QuestionEssay, Prompt: "Trình bày…",
		Explanation:  "giải thích",
		SampleAnswer: "đáp án gợi ý",
		Rubric:       "sơ đồ 1đ; công thức 1đ",
	}
}

func fillBlankQuestion() *models.Question {
	return &models.Question{
		ID: uuid.New(), Type: models.QuestionFillBlank, Prompt: "Điền ___ vào",
		Explanation: "giải thích",
		Options: []*models.QuestionOption{
			{ID: uuid.New(), Content: "thuận", IsCorrect: true},
			{ID: uuid.New(), Content: "tỉ lệ thuận", IsCorrect: true},
		},
	}
}

// Với câu điền khuyết, chính nội dung phương án là đáp án: tắt cờ IsCorrect
// thôi thì chưa đủ, phải gỡ hẳn danh sách.
func TestStripQuestionAnswersRemovesFillBlankOptions(t *testing.T) {
	q := fillBlankQuestion()
	stripQuestionAnswers([]*models.Question{q})

	if len(q.Options) != 0 {
		t.Errorf("đáp án điền khuyết phải bị gỡ hết, còn lại %d dòng", len(q.Options))
	}
	if q.Explanation != "" {
		t.Error("giải thích phải bị gỡ")
	}
}

func TestStripQuestionAnswersRemovesEssayGuidance(t *testing.T) {
	q := essayQuestion()
	stripQuestionAnswers([]*models.Question{q})

	if q.SampleAnswer != "" {
		t.Error("đáp án gợi ý phải bị gỡ khi học viên đang làm bài")
	}
	if q.Rubric != "" {
		t.Error("tiêu chí chấm phải bị gỡ khi học viên đang làm bài")
	}
}

// Câu trắc nghiệm giữ lại nội dung phương án (học viên cần đọc để chọn),
// chỉ tắt cờ đáp án đúng.
func TestStripQuestionAnswersKeepsChoiceContent(t *testing.T) {
	q := &models.Question{
		ID: uuid.New(), Type: models.QuestionSingleChoice, Prompt: "Chọn…",
		Explanation: "giải thích",
		Options: []*models.QuestionOption{
			{ID: uuid.New(), Content: "Vôn", IsCorrect: false},
			{ID: uuid.New(), Content: "Ôm", IsCorrect: true},
		},
	}
	stripQuestionAnswers([]*models.Question{q})

	if len(q.Options) != 2 {
		t.Fatalf("phương án trắc nghiệm phải được giữ, còn %d", len(q.Options))
	}
	for _, o := range q.Options {
		if o.IsCorrect {
			t.Error("cờ đáp án đúng phải bị tắt")
		}
		if o.Content == "" {
			t.Error("nội dung phương án phải được giữ để học viên đọc")
		}
	}
}
