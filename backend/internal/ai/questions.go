package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Tên dạng câu hỏi dùng trong prompt và trong JSON mà AI trả về.
// Cố tình khác với hằng số của package models để prompt đọc tự nhiên hơn
// với model; việc quy đổi nằm ở hàm ToQuestionType.
const (
	KindMultipleChoice = "multiple_choice"
	KindTrueFalse      = "true_false"
	KindFillBlank      = "fill_blank"
	KindEssay          = "essay"
)

// Ba mức độ nhận thức được dùng trong đề (mục 3.4.2).
var Levels = []string{"Nhận biết", "Thông hiểu", "Vận dụng"}

// QuestionSpec là yêu cầu số lượng câu cho một dạng câu hỏi.
type QuestionSpec struct {
	Kind  string `json:"kind"`
	Count int    `json:"count"`
}

// GenerateQuestionsParams gom toàn bộ lựa chọn của giáo viên ở màn hình tạo đề.
type GenerateQuestionsParams struct {
	Subject string
	Grade   string
	Topic   string
	// Nội dung bài học do giáo viên nhập; bỏ trống khi sinh câu hỏi từ PDF.
	Content string
	// Tài liệu PDF đã mã hoá base64; có thì AI đọc tài liệu thay cho Content.
	PDFBase64 string
	PDFName   string
	// Số câu của từng dạng.
	Specs []QuestionSpec
	// Các mức độ nhận thức được phép xuất hiện; rỗng = dùng cả ba.
	Levels []string
}

// GeneratedQuestion là một câu hỏi sau khi đã kiểm tra hợp lệ.
type GeneratedQuestion struct {
	Type        string  `json:"type"`
	Level       string  `json:"level"`
	Question    string  `json:"question"`
	Explanation string  `json:"explanation"`
	Points      float64 `json:"points"`
	// multiple_choice
	Options []string `json:"options,omitempty"`
	Answer  *int     `json:"answer,omitempty"`
	// true_false
	Statements []Statement `json:"statements,omitempty"`
	// fill_blank
	Answers []string `json:"answers,omitempty"`
	// essay
	SampleAnswer string `json:"sampleAnswer,omitempty"`
	Rubric       string `json:"rubric,omitempty"`
}

type Statement struct {
	Text   string `json:"text"`
	Answer *bool  `json:"answer"`
}

// GenerateQuestionsResult kèm số lần gọi API để hiển thị ở nhật ký tác vụ.
type GenerateQuestionsResult struct {
	Questions []GeneratedQuestion `json:"questions"`
	Attempts  int                 `json:"attempts"`
}

// maxGenerateAttempts: sai cấu trúc thì yêu cầu AI sinh lại, nhưng có trần để
// một yêu cầu hỏng không kéo dài vô hạn (mục 4.7).
const maxGenerateAttempts = 3

// GenerateQuestions sinh câu hỏi rồi kiểm tra tự động; không đạt thì sinh lại.
func (c *Client) GenerateQuestions(ctx context.Context, p GenerateQuestionsParams) (*GenerateQuestionsResult, error) {
	if !c.Enabled() {
		return nil, ErrDisabled
	}
	specs, total, err := normalizeSpecs(p.Specs)
	if err != nil {
		return nil, err
	}
	p.Specs = specs

	prompt := buildQuestionPrompt(p)
	var lastErr error

	for attempt := 1; attempt <= maxGenerateAttempts; attempt++ {
		text := prompt
		if lastErr != nil {
			// Nói rõ lỗi của lần trước để lần sinh lại không lặp lại đúng lỗi đó.
			text = prompt + "\n\nLƯU Ý: Kết quả lần trước bị loại vì: " +
				lastErr.Error() + "\nHãy sinh lại toàn bộ cho đúng yêu cầu."
		}

		var content Content
		if p.PDFBase64 != "" {
			content = UserWithPDF(text, p.PDFBase64)
		} else {
			content = UserText(text)
		}

		raw, err := c.Generate(ctx, Request{
			System:      questionSystemPrompt,
			Contents:    []Content{content},
			JSON:        true,
			Temperature: 0.4,
		})
		if err != nil {
			return nil, err
		}

		var payload struct {
			Questions []GeneratedQuestion `json:"questions"`
		}
		if err := extractJSON(raw, &payload); err != nil {
			lastErr = err
			continue
		}
		questions, err := validateQuestions(payload.Questions, specs, total)
		if err != nil {
			lastErr = err
			continue
		}
		return &GenerateQuestionsResult{Questions: questions, Attempts: attempt}, nil
	}

	return nil, fmt.Errorf("AI đã sinh lại %d lần nhưng kết quả vẫn không đạt yêu cầu (%v)", maxGenerateAttempts, lastErr)
}

const questionSystemPrompt = `Bạn là một trợ lý ra đề cho giáo viên phổ thông Việt Nam.
Bạn luôn trả lời bằng tiếng Việt và luôn trả về đúng một đối tượng JSON, không kèm lời dẫn hay chú thích.`

func buildQuestionPrompt(p GenerateQuestionsParams) string {
	b := &promptBuilder{
		role: fmt.Sprintf("Bạn là giáo viên %s giàu kinh nghiệm, đang soạn đề kiểm tra cho học sinh %s.",
			orDefault(p.Subject, "bộ môn"), orDefault(p.Grade, "trung học phổ thông")),
	}

	b.addContext("Môn học", p.Subject)
	b.addContext("Lớp", p.Grade)
	b.addContext("Chủ đề", p.Topic)
	if p.PDFBase64 != "" {
		b.addContext("Tài liệu đính kèm", orDefault(p.PDFName, "tài liệu PDF"))
		b.context = append(b.context, "- Toàn bộ câu hỏi phải bám sát nội dung tài liệu PDF đính kèm ở trên.")
	} else {
		b.addContext("Nội dung bài học", p.Content)
	}

	levels := p.Levels
	if len(levels) == 0 {
		levels = Levels
	}
	for _, spec := range p.Specs {
		b.addTask("Soạn %d câu dạng %s (%s).", spec.Count, vietnameseKind(spec.Kind), spec.Kind)
	}
	b.addTask("Phân bổ các câu theo mức độ nhận thức: %s.", strings.Join(levels, ", "))

	b.constraints = []string{
		"Nội dung phải chính xác về kiến thức và bám sát chương trình giáo dục phổ thông Việt Nam.",
		"Viết toàn bộ bằng tiếng Việt, diễn đạt phù hợp với lứa tuổi học sinh.",
		"Không lặp lại câu hỏi, không ra hai câu hỏi cùng một ý.",
		`Trường "level" chỉ nhận đúng một trong các giá trị: ` + strings.Join(levels, " / ") + ".",
		"Mỗi câu đều phải có trường \"explanation\" giải thích ngắn gọn vì sao đáp án đó đúng.",
		"Câu trắc nghiệm phải có đúng 4 phương án và đúng một phương án đúng.",
		"Câu đúng/sai phải có từ 2 đến 4 phát biểu, mỗi phát biểu có giá trị đúng hoặc sai rõ ràng.",
		"Câu điền khuyết dùng dấu ___ ở vị trí cần điền và liệt kê mọi cách diễn đạt được chấp nhận.",
		"Câu tự luận phải có đáp án gợi ý và tiêu chí chấm chi tiết.",
	}
	b.format = questionJSONFormat
	return b.String()
}

const questionJSONFormat = `Chỉ trả về JSON theo đúng cấu trúc sau, không thêm bất kỳ ký tự nào khác:
{
  "questions": [
    {
      "type": "multiple_choice",
      "level": "Thông hiểu",
      "question": "Đơn vị của điện trở là gì?",
      "options": ["Vôn (V)", "Ampe (A)", "Ôm (Ω)", "Oát (W)"],
      "answer": 2,
      "explanation": "Điện trở đo bằng ôm, ký hiệu Ω."
    },
    {
      "type": "true_false",
      "level": "Nhận biết",
      "question": "Xét các phát biểu sau về dòng điện:",
      "statements": [
        { "text": "Dòng điện là dòng chuyển dời có hướng của các điện tích.", "answer": true },
        { "text": "Đơn vị của cường độ dòng điện là vôn.", "answer": false }
      ],
      "explanation": "Cường độ dòng điện đo bằng ampe, không phải vôn."
    },
    {
      "type": "fill_blank",
      "level": "Thông hiểu",
      "question": "Định luật Ôm cho biết cường độ dòng điện tỉ lệ ___ với hiệu điện thế.",
      "answers": ["thuận", "tỉ lệ thuận"],
      "explanation": "I = U/R nên I tỉ lệ thuận với U khi R không đổi."
    },
    {
      "type": "essay",
      "level": "Vận dụng",
      "question": "Trình bày cách xác định điện trở của một dây dẫn bằng vôn kế và ampe kế.",
      "sampleAnswer": "Mắc vôn kế song song, ampe kế nối tiếp, đo U và I rồi tính R = U/I.",
      "rubric": "Nêu đúng sơ đồ mạch (1 điểm); nêu đúng công thức R = U/I (1 điểm); trình bày rõ các bước đo (1 điểm).",
      "points": 3,
      "explanation": "Câu hỏi kiểm tra khả năng vận dụng định luật Ôm vào thực hành."
    }
  ]
}
Trường "answer" của câu trắc nghiệm là chỉ số của phương án đúng trong mảng "options", đếm từ 0.`

// normalizeSpecs bỏ các dạng có số câu bằng 0 và chặn yêu cầu quá lớn.
func normalizeSpecs(specs []QuestionSpec) ([]QuestionSpec, int, error) {
	out := make([]QuestionSpec, 0, len(specs))
	total := 0
	seen := map[string]bool{}
	for _, s := range specs {
		kind := strings.TrimSpace(s.Kind)
		if s.Count <= 0 {
			continue
		}
		if !validKind(kind) {
			return nil, 0, fmt.Errorf("dạng câu hỏi %q không được hỗ trợ", s.Kind)
		}
		if seen[kind] {
			return nil, 0, fmt.Errorf("dạng câu hỏi %q bị khai báo hai lần", vietnameseKind(kind))
		}
		seen[kind] = true
		out = append(out, QuestionSpec{Kind: kind, Count: s.Count})
		total += s.Count
	}
	if total == 0 {
		return nil, 0, errors.New("vui lòng chọn số câu cho ít nhất một dạng câu hỏi")
	}
	if total > 50 {
		return nil, 0, errors.New("mỗi lần chỉ tạo tối đa 50 câu hỏi")
	}
	return out, total, nil
}

// validateQuestions là lớp kiểm tra tự động ở mục 4.7: đủ số câu, đúng số câu
// từng dạng và mỗi câu có đáp án hợp lệ. Chỉ cần một câu sai là loại cả mẻ và
// yêu cầu AI sinh lại, vì đề thiếu câu không dùng được.
func validateQuestions(items []GeneratedQuestion, specs []QuestionSpec, total int) ([]GeneratedQuestion, error) {
	if len(items) != total {
		return nil, fmt.Errorf("cần %d câu nhưng AI trả về %d câu", total, len(items))
	}

	want := map[string]int{}
	for _, s := range specs {
		want[s.Kind] = s.Count
	}
	got := map[string]int{}
	seenPrompt := map[string]bool{}

	out := make([]GeneratedQuestion, 0, len(items))
	for i := range items {
		q := items[i]
		q.Type = strings.TrimSpace(q.Type)
		q.Question = strings.TrimSpace(q.Question)
		q.Explanation = strings.TrimSpace(q.Explanation)
		q.Level = normalizeLevel(q.Level)

		pos := i + 1
		if q.Question == "" {
			return nil, fmt.Errorf("câu %d thiếu nội dung câu hỏi", pos)
		}
		if q.Explanation == "" {
			return nil, fmt.Errorf("câu %d thiếu phần giải thích", pos)
		}
		key := strings.ToLower(q.Question)
		if seenPrompt[key] {
			return nil, fmt.Errorf("câu %d bị trùng nội dung với một câu trước đó", pos)
		}
		seenPrompt[key] = true

		if err := validateOne(&q, pos); err != nil {
			return nil, err
		}
		if q.Points <= 0 {
			q.Points = 1
		}
		got[q.Type]++
		out = append(out, q)
	}

	for kind, n := range want {
		if got[kind] != n {
			return nil, fmt.Errorf("cần %d câu dạng %s nhưng nhận được %d câu",
				n, vietnameseKind(kind), got[kind])
		}
	}
	return out, nil
}

// validateOne kiểm tra phần riêng của từng dạng câu hỏi.
func validateOne(q *GeneratedQuestion, pos int) error {
	switch q.Type {
	case KindMultipleChoice:
		opts := trimAll(q.Options)
		if len(opts) != 4 {
			return fmt.Errorf("câu %d: câu trắc nghiệm phải có đúng 4 phương án, đang có %d", pos, len(opts))
		}
		if dup := firstDuplicate(opts); dup != "" {
			return fmt.Errorf("câu %d: hai phương án trùng nhau (%q)", pos, dup)
		}
		if q.Answer == nil || *q.Answer < 0 || *q.Answer >= len(opts) {
			return fmt.Errorf("câu %d: chỉ số đáp án đúng không hợp lệ", pos)
		}
		q.Options = opts

	case KindTrueFalse:
		sts := make([]Statement, 0, len(q.Statements))
		for _, s := range q.Statements {
			text := strings.TrimSpace(s.Text)
			if text == "" {
				continue
			}
			if s.Answer == nil {
				return fmt.Errorf("câu %d: phát biểu %q chưa có giá trị Đúng/Sai", pos, text)
			}
			sts = append(sts, Statement{Text: text, Answer: s.Answer})
		}
		if len(sts) < 2 {
			return fmt.Errorf("câu %d: câu đúng/sai cần ít nhất 2 phát biểu", pos)
		}
		q.Statements = sts

	case KindFillBlank:
		answers := trimAll(q.Answers)
		if len(answers) == 0 {
			return fmt.Errorf("câu %d: câu điền khuyết phải có ít nhất một đáp án", pos)
		}
		if !strings.Contains(q.Question, "___") {
			return fmt.Errorf("câu %d: câu điền khuyết phải có chỗ trống đánh dấu bằng ___", pos)
		}
		q.Answers = answers

	case KindEssay:
		q.SampleAnswer = strings.TrimSpace(q.SampleAnswer)
		q.Rubric = strings.TrimSpace(q.Rubric)
		if q.SampleAnswer == "" {
			return fmt.Errorf("câu %d: câu tự luận thiếu đáp án gợi ý", pos)
		}
		if q.Rubric == "" {
			return fmt.Errorf("câu %d: câu tự luận thiếu tiêu chí chấm", pos)
		}

	default:
		return fmt.Errorf("câu %d: dạng câu hỏi %q không được hỗ trợ", pos, q.Type)
	}
	return nil
}

func validKind(kind string) bool {
	switch kind {
	case KindMultipleChoice, KindTrueFalse, KindFillBlank, KindEssay:
		return true
	}
	return false
}

func vietnameseKind(kind string) string {
	switch kind {
	case KindMultipleChoice:
		return "trắc nghiệm nhiều lựa chọn"
	case KindTrueFalse:
		return "đúng/sai"
	case KindFillBlank:
		return "điền khuyết"
	case KindEssay:
		return "tự luận"
	}
	return kind
}

// normalizeLevel đưa mức độ về đúng một trong ba giá trị chuẩn; không nhận ra
// thì để rỗng chứ không loại cả câu, vì mức độ chỉ dùng để phân loại.
func normalizeLevel(level string) string {
	l := strings.ToLower(strings.TrimSpace(level))
	for _, want := range Levels {
		if l == strings.ToLower(want) {
			return want
		}
	}
	switch {
	case strings.Contains(l, "nhận biết"), l == "remember", l == "knowledge":
		return "Nhận biết"
	case strings.Contains(l, "thông hiểu"), l == "understand", l == "comprehension":
		return "Thông hiểu"
	case strings.Contains(l, "vận dụng"), l == "apply", l == "application":
		return "Vận dụng"
	}
	return ""
}

func trimAll(in []string) []string {
	out := make([]string, 0, len(in))
	for _, s := range in {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func firstDuplicate(in []string) string {
	seen := map[string]bool{}
	for _, s := range in {
		k := strings.ToLower(s)
		if seen[k] {
			return s
		}
		seen[k] = true
	}
	return ""
}

func orDefault(v, fallback string) string {
	if v = strings.TrimSpace(v); v != "" {
		return v
	}
	return fallback
}
