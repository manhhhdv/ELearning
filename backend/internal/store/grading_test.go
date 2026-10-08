package store

import (
	"testing"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

func opt(content string, correct bool) *models.QuestionOption {
	return &models.QuestionOption{ID: uuid.New(), Content: content, IsCorrect: correct}
}

func TestScoreAnswerEssayLeftForTeacher(t *testing.T) {
	q := &models.Question{Type: models.QuestionEssay, Points: 3}
	score, isCorrect, auto := scoreAnswer(q, AnswerInput{EssayText: "bài làm"})
	if auto {
		t.Fatal("câu tự luận không được chấm tự động")
	}
	if score != 0 || isCorrect != nil {
		t.Fatalf("câu tự luận phải để trống điểm, nhận được score=%v isCorrect=%v", score, isCorrect)
	}
}

func TestScoreTrueFalsePartialCredit(t *testing.T) {
	a, b, c, d := opt("p1", true), opt("p2", false), opt("p3", true), opt("p4", false)
	q := &models.Question{Type: models.QuestionTrueFalse, Points: 4, Options: []*models.QuestionOption{a, b, c, d}}

	cases := []struct {
		name     string
		selected []uuid.UUID
		want     float64
		correct  bool
	}{
		{"đúng hết", []uuid.UUID{a.ID, c.ID}, 4, true},
		{"sai hết", []uuid.UUID{b.ID, d.ID}, 0, false},
		{"đúng một nửa", []uuid.UUID{a.ID, b.ID}, 2, false},
		{"không chọn gì", nil, 2, false}, // hai phát biểu Sai được tính đúng
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			score, isCorrect, auto := scoreAnswer(q, AnswerInput{SelectedOptionIDs: tc.selected})
			if !auto {
				t.Fatal("câu đúng/sai phải được chấm tự động")
			}
			if score != tc.want {
				t.Errorf("điểm = %v, mong đợi %v", score, tc.want)
			}
			if isCorrect == nil || *isCorrect != tc.correct {
				t.Errorf("isCorrect = %v, mong đợi %v", isCorrect, tc.correct)
			}
		})
	}
}

func TestScoreFillBlankNormalizesAnswer(t *testing.T) {
	q := &models.Question{
		Type:    models.QuestionFillBlank,
		Points:  2,
		Options: []*models.QuestionOption{opt("tỉ lệ thuận", true), opt("thuận", true)},
	}

	cases := []struct {
		text string
		want float64
	}{
		{"tỉ lệ thuận", 2},
		{"  Tỉ Lệ   Thuận  ", 2}, // thừa khoảng trắng và khác hoa thường
		{"thuận.", 2},            // dấu câu ở cuối
		{"THUẬN", 2},
		{"nghịch", 0},
		{"", 0},
	}
	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			score, _, auto := scoreAnswer(q, AnswerInput{EssayText: tc.text})
			if !auto {
				t.Fatal("câu điền khuyết phải được chấm tự động")
			}
			if score != tc.want {
				t.Errorf("điểm cho %q = %v, mong đợi %v", tc.text, score, tc.want)
			}
		})
	}
}

// Tiếng Việt gõ kiểu tổ hợp (NFD) phải khớp với đáp án lưu kiểu dựng sẵn (NFC).
func TestScoreFillBlankUnicodeNormalization(t *testing.T) {
	q := &models.Question{
		Type:    models.QuestionFillBlank,
		Points:  1,
		Options: []*models.QuestionOption{opt("điện trở", true)},
	}
	decomposed := "điện trở"
	if score, _, _ := scoreAnswer(q, AnswerInput{EssayText: decomposed}); score != 1 {
		t.Errorf("câu trả lời gõ tổ hợp phải được chấp nhận, điểm = %v", score)
	}
}

func TestScoreSingleChoiceExactMatch(t *testing.T) {
	a, b := opt("A", true), opt("B", false)
	q := &models.Question{Type: models.QuestionSingleChoice, Points: 1, Options: []*models.QuestionOption{a, b}}

	if score, _, _ := scoreAnswer(q, AnswerInput{SelectedOptionIDs: []uuid.UUID{a.ID}}); score != 1 {
		t.Errorf("chọn đúng phải được 1 điểm, nhận được %v", score)
	}
	if score, _, _ := scoreAnswer(q, AnswerInput{SelectedOptionIDs: []uuid.UUID{a.ID, b.ID}}); score != 0 {
		t.Errorf("chọn thừa phương án phải bị 0 điểm, nhận được %v", score)
	}
}
