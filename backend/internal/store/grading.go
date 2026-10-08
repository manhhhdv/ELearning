package store

import (
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"

	"github.com/manhnv/elearning/backend/internal/models"
)

// scoreAnswer chấm một câu trả lời theo đúng dạng câu hỏi (mục 3.4.3).
// Câu tự luận trả về autoGraded = false: điểm do giáo viên quyết định.
func scoreAnswer(q *models.Question, answer AnswerInput) (score float64, isCorrect *bool, autoGraded bool) {
	switch q.Type {
	case models.QuestionEssay:
		return 0, nil, false

	case models.QuestionTrueFalse:
		score = scoreTrueFalse(q, answer.SelectedOptionIDs)

	case models.QuestionFillBlank:
		if matchFillBlank(q, answer.EssayText) {
			score = q.Points
		}

	default: // single_choice, multi_choice
		if isChoiceCorrect(q, answer.SelectedOptionIDs) {
			score = q.Points
		}
	}

	// "Đúng" ở đây nghĩa là đạt điểm tối đa của câu; câu đúng/sai đạt một phần
	// vẫn được tính điểm nhưng không đánh dấu là đúng.
	correct := q.Points > 0 && score >= q.Points
	return score, &correct, true
}

// scoreTrueFalse chấm theo từng phát biểu: mỗi phát biểu đúng được một phần
// điểm bằng nhau, nên học sinh làm đúng một nửa vẫn có điểm tương ứng.
//
// Quy ước: selected chứa ID của các phát biểu học sinh đánh dấu "Đúng";
// phát biểu không có trong danh sách coi như học sinh chọn "Sai".
func scoreTrueFalse(q *models.Question, selected []uuid.UUID) float64 {
	if len(q.Options) == 0 {
		return 0
	}
	marked := make(map[uuid.UUID]bool, len(selected))
	for _, id := range selected {
		marked[id] = true
	}

	right := 0
	for _, o := range q.Options {
		if marked[o.ID] == o.IsCorrect {
			right++
		}
	}
	return q.Points * float64(right) / float64(len(q.Options))
}

// matchFillBlank so câu trả lời với danh sách đáp án được chấp nhận sau khi
// chuẩn hoá, để khác biệt nhỏ về cách gõ không bị chấm sai (mục 3.4.3).
func matchFillBlank(q *models.Question, text string) bool {
	answer := normalizeAnswer(text)
	if answer == "" {
		return false
	}
	for _, o := range q.Options {
		if normalizeAnswer(o.Content) == answer {
			return true
		}
	}
	return false
}

// normalizeAnswer đưa câu trả lời về dạng so sánh được:
// chuẩn hoá Unicode (tiếng Việt gõ tổ hợp hay dựng sẵn đều cho cùng kết quả),
// bỏ hoa/thường, gộp khoảng trắng thừa và bỏ dấu câu ở hai đầu.
func normalizeAnswer(s string) string {
	s = norm.NFC.String(s)
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r)
	})
	return strings.Join(strings.Fields(s), " ")
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
