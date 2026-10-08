package export

import (
	"strings"
	"testing"

	"github.com/manhnv/elearning/backend/internal/ai"
	"github.com/manhnv/elearning/backend/internal/docx"
	"github.com/manhnv/elearning/backend/internal/models"
)

// extract nhận thẳng kết quả của hàm dựng tệp: extract(t)(Questions(...)).
func extract(t *testing.T) func([]byte, error) string {
	return func(raw []byte, err error) string {
		t.Helper()
		if err != nil {
			t.Fatalf("build: %v", err)
		}
		text, err := docx.ExtractText(raw)
		if err != nil {
			t.Fatalf("ExtractText: %v", err)
		}
		return text
	}
}

func assertOrder(t *testing.T, text string, parts ...string) {
	t.Helper()
	pos := 0
	for _, p := range parts {
		i := strings.Index(text[pos:], p)
		if i < 0 {
			t.Fatalf("%q not found (in order) in:\n%s", p, text)
		}
		pos += i + len(p)
	}
}

func TestLessonPlan(t *testing.T) {
	plan := &ai.LessonPlan{
		Title:      "Định luật Ôm",
		Objectives: []string{"Phát biểu được định luật Ôm", "  "},
		Activities: []ai.LessonActivity{
			{Name: "Hoạt động 1: Khởi động", DurationMinutes: 5, Objective: "Gợi tò mò",
				TeacherActions: "- Chiếu video\n- Đặt câu hỏi **vì sao**", StudentActions: "Quan sát, trả lời", Output: "Câu trả lời"},
			{Name: "Hoạt động 2: Hình thành kiến thức", DurationMinutes: 40,
				TeacherActions: "Hướng dẫn thí nghiệm", StudentActions: "Làm thí nghiệm"},
		},
		Homework: "Làm bài 1, 2 SGK",
	}
	text := extract(t)(LessonPlan(LessonPlanInfo{Subject: "Vật lí", Grade: "11", DurationMinutes: 45}, plan))

	assertOrder(t, text,
		"KẾ HOẠCH BÀI DẠY", "Định luật Ôm", "Môn học: Vật lí · Lớp: 11 · Thời lượng: 45 phút",
		"I. YÊU CẦU CẦN ĐẠT", "Phát biểu được định luật Ôm",
		// Không có thiết bị nên mục kế tiếp được đánh số II.
		"II. TIẾN TRÌNH DẠY HỌC",
		"Hoạt động 1: Khởi động (5 phút)", "a) Mục tiêu: ", "Gợi tò mò", "b) Sản phẩm: ", "c) Tổ chức thực hiện:",
		"Hoạt động của giáo viên", "Hoạt động của học sinh", "Chiếu video", "vì sao", "Quan sát, trả lời",
		"Hoạt động 2: Hình thành kiến thức (40 phút)", "a) Tổ chức thực hiện:",
		"III. BÀI TẬP VỀ NHÀ", "Làm bài 1, 2 SGK",
	)
	if strings.Contains(text, "**") {
		t.Errorf("markdown markers leaked: %q", text)
	}
	if strings.Contains(text, "LƯU Ý") {
		t.Errorf("empty notes section should be omitted")
	}
}

func sampleQuestions() []*models.Question {
	return []*models.Question{
		{Type: models.QuestionSingleChoice, Prompt: "Đơn vị điện trở?", Points: 0.5, Level: models.LevelRemember,
			Explanation: "Ôm là đơn vị điện trở.",
			Options:     []*models.QuestionOption{{Content: "Vôn"}, {Content: "Ôm", IsCorrect: true}, {Content: "Ampe"}}},
		{Type: models.QuestionMultiChoice, Prompt: "Đại lượng vô hướng?", Points: 0.5,
			Options: []*models.QuestionOption{{Content: "Khối lượng", IsCorrect: true}, {Content: "Lực"}, {Content: "Nhiệt độ", IsCorrect: true}}},
		{Type: models.QuestionTrueFalse, Prompt: "Xét mạch điện:", Points: 1,
			Options: []*models.QuestionOption{{Content: "I tỉ lệ với U", IsCorrect: true}, {Content: "R phụ thuộc U"}}},
		{Type: models.QuestionFillBlank, Prompt: "Công thức: I = ___", Points: 1,
			Options: []*models.QuestionOption{{Content: "U/R", IsCorrect: true}, {Content: "U : R", IsCorrect: true}}},
		{Type: models.QuestionEssay, Prompt: "Trình bày cách đo điện trở.", Points: 2,
			SampleAnswer: "Dùng vôn kế và ampe kế.", Rubric: "Nêu đúng sơ đồ (1đ)."},
	}
}

func TestQuestionsWithoutAnswers(t *testing.T) {
	info := QuestionSetInfo{ProgramTitle: "Vật lí 11", Title: "Kiểm tra 15 phút", Instructions: "Làm bài nghiêm túc", TimeLimitMinutes: 15}
	text := extract(t)(Questions(info, sampleQuestions(), false))

	assertOrder(t, text,
		"Kiểm tra 15 phút", "Lớp học: Vật lí 11 · Thời gian làm bài: 15 phút · 5 câu · tổng 5 điểm",
		"Họ và tên:", "Hướng dẫn: ", "Làm bài nghiêm túc",
		"Câu 1. ", "Đơn vị điện trở?", "(0,5 điểm)", "A. ", "Vôn", "B. ", "Ôm", "C. ", "Ampe",
		"Câu 2. ", "(Chọn tất cả phương án đúng)",
		"Câu 3. ", "Phát biểu", "Đúng", "Sai", "a) ", "I tỉ lệ với U", "b) ", "R phụ thuộc U",
		"Câu 4. ", "Công thức: I = "+dottedBlank,
		"Câu 5. ", "(2 điểm)",
	)
	for _, leaked := range []string{"ĐÁP ÁN", "Lời giải", "U/R", "Dùng vôn kế", "___"} {
		if strings.Contains(text, leaked) {
			t.Errorf("student copy must not contain %q", leaked)
		}
	}
}

func TestQuestionsWithAnswers(t *testing.T) {
	text := extract(t)(Questions(QuestionSetInfo{Title: "Kiểm tra"}, sampleQuestions(), true))

	assertOrder(t, text,
		"Câu 5. ", "ĐÁP ÁN VÀ HƯỚNG DẪN CHẤM",
		"Bảng đáp án", "Câu", "1", "2", "Đáp án", "B", "A, C",
		"Câu 1", "(Nhận biết · 0,5 điểm)", "Đáp án: ", "B", "Lời giải: ", "Ôm là đơn vị điện trở.",
		"Câu 2", "Đáp án: ", "A, C",
		"Câu 3", "Đáp án: ", "a) Đúng; b) Sai",
		"Câu 4", "Đáp án được chấp nhận: ", "U/R; U : R",
		"Câu 5", "Đáp án gợi ý: ", "Dùng vôn kế và ampe kế.", "Hướng dẫn chấm: ", "Nêu đúng sơ đồ (1đ).",
	)
}

func TestQuestionsEqualPointsHidePerQuestionPoints(t *testing.T) {
	qs := []*models.Question{
		{Type: models.QuestionEssay, Prompt: "Câu một", Points: 1},
		{Type: models.QuestionEssay, Prompt: "Câu hai", Points: 1},
	}
	text := extract(t)(Questions(QuestionSetInfo{Title: "Đề"}, qs, false))
	if strings.Contains(text, "(1 điểm)") {
		t.Errorf("per-question points should be hidden when all equal: %q", text)
	}
}

func TestLessonRichText(t *testing.T) {
	node := &models.Node{
		Title:       "Định luật Ôm",
		Description: "Bài đọc trước giờ học",
		Lesson: &models.Lesson{
			ContentType:     "richtext",
			DurationMinutes: 20,
			Body:            "## Nội dung\n\nCông thức $I = \\frac{U}{R}$.",
			Attachments:     []models.LessonAttachment{{Name: "Phiếu học tập", URL: "https://example.com/phieu.pdf"}},
		},
	}
	text := extract(t)(Lesson("Vật lí 11", "", node))
	assertOrder(t, text,
		"Định luật Ôm", "Lớp học: Vật lí 11 · Thời lượng: 20 phút", "Bài đọc trước giờ học",
		"Nội dung", "Công thức ", "Tài liệu kèm theo", "Phiếu học tập",
	)
	if strings.Contains(text, `\frac`) {
		t.Errorf("formula should be converted to a Word equation: %q", text)
	}
}

func TestLessonEmbedLinksSource(t *testing.T) {
	node := &models.Node{
		Title:  "Video thí nghiệm",
		Lesson: &models.Lesson{ContentType: "video", EmbedURL: "https://youtu.be/abc", Body: "Xem trước khi đến lớp."},
	}
	text := extract(t)(Lesson("", "", node))
	assertOrder(t, text, "Video thí nghiệm", "Video: ", "https://youtu.be/abc", "Xem trước khi đến lớp.")
}

func TestFileName(t *testing.T) {
	if got := FileName("giao an", "Định luật Ôm – lớp 11"); got != "giao-an-dinh-luat-om-lop-11.docx" {
		t.Errorf("FileName = %q", got)
	}
	long := FileName("de", strings.Repeat("rất dài ", 30))
	if len(long) > 85 || strings.Contains(long, "-.docx") {
		t.Errorf("long file name not trimmed cleanly: %q", long)
	}
}
