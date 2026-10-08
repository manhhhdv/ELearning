package ai

import "testing"

func mc(q string) GeneratedQuestion {
	answer := 2
	return GeneratedQuestion{
		Type: KindMultipleChoice, Level: "Thông hiểu", Question: q,
		Options: []string{"A", "B", "C", "D"}, Answer: &answer, Explanation: "vì thế",
	}
}

func specsOf(kind string, n int) []QuestionSpec { return []QuestionSpec{{Kind: kind, Count: n}} }

func TestValidateQuestionsAcceptsValidBatch(t *testing.T) {
	items := []GeneratedQuestion{mc("câu 1"), mc("câu 2")}
	out, err := validateQuestions(items, specsOf(KindMultipleChoice, 2), 2)
	if err != nil {
		t.Fatalf("mẻ hợp lệ bị loại: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("mong đợi 2 câu, nhận được %d", len(out))
	}
	if out[0].Points != 1 {
		t.Errorf("điểm mặc định phải là 1, nhận được %v", out[0].Points)
	}
}

func TestValidateQuestionsRejectsBadBatches(t *testing.T) {
	yes, no := true, false
	bad3 := mc("thiếu phương án")
	bad3.Options = []string{"A", "B", "C"}
	noAnswer := mc("thiếu đáp án")
	noAnswer.Answer = nil
	outOfRange := mc("đáp án ngoài phạm vi")
	idx := 9
	outOfRange.Answer = &idx
	noExplain := mc("thiếu giải thích")
	noExplain.Explanation = ""

	cases := []struct {
		name  string
		items []GeneratedQuestion
		specs []QuestionSpec
		total int
	}{
		{"thiếu câu", []GeneratedQuestion{mc("một")}, specsOf(KindMultipleChoice, 2), 2},
		{"trùng câu hỏi", []GeneratedQuestion{mc("x"), mc("x")}, specsOf(KindMultipleChoice, 2), 2},
		{"không đủ 4 phương án", []GeneratedQuestion{bad3}, specsOf(KindMultipleChoice, 1), 1},
		{"thiếu đáp án", []GeneratedQuestion{noAnswer}, specsOf(KindMultipleChoice, 1), 1},
		{"đáp án ngoài phạm vi", []GeneratedQuestion{outOfRange}, specsOf(KindMultipleChoice, 1), 1},
		{"thiếu giải thích", []GeneratedQuestion{noExplain}, specsOf(KindMultipleChoice, 1), 1},
		{"đúng/sai chỉ có 1 phát biểu", []GeneratedQuestion{{
			Type: KindTrueFalse, Question: "q", Explanation: "e",
			Statements: []Statement{{Text: "p1", Answer: &yes}},
		}}, specsOf(KindTrueFalse, 1), 1},
		{"đúng/sai thiếu giá trị", []GeneratedQuestion{{
			Type: KindTrueFalse, Question: "q", Explanation: "e",
			Statements: []Statement{{Text: "p1", Answer: &yes}, {Text: "p2"}},
		}}, specsOf(KindTrueFalse, 1), 1},
		{"điền khuyết không có chỗ trống", []GeneratedQuestion{{
			Type: KindFillBlank, Question: "không có chỗ trống", Explanation: "e",
			Answers: []string{"x"},
		}}, specsOf(KindFillBlank, 1), 1},
		{"điền khuyết không có đáp án", []GeneratedQuestion{{
			Type: KindFillBlank, Question: "điền ___ vào", Explanation: "e",
		}}, specsOf(KindFillBlank, 1), 1},
		{"tự luận thiếu tiêu chí chấm", []GeneratedQuestion{{
			Type: KindEssay, Question: "q", Explanation: "e", SampleAnswer: "đáp án",
		}}, specsOf(KindEssay, 1), 1},
		{"sai dạng câu hỏi", []GeneratedQuestion{{
			Type: "matching", Question: "q", Explanation: "e",
		}}, specsOf(KindEssay, 1), 1},
		{"lệch số câu từng dạng", []GeneratedQuestion{mc("a"), {
			Type: KindTrueFalse, Question: "q", Explanation: "e",
			Statements: []Statement{{Text: "p1", Answer: &yes}, {Text: "p2", Answer: &no}},
		}}, []QuestionSpec{{Kind: KindMultipleChoice, Count: 2}}, 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := validateQuestions(tc.items, tc.specs, tc.total); err == nil {
				t.Fatal("mẻ hỏng đáng lẽ phải bị loại để AI sinh lại")
			}
		})
	}
}

func TestNormalizeSpecs(t *testing.T) {
	specs, total, err := normalizeSpecs([]QuestionSpec{
		{Kind: KindMultipleChoice, Count: 3},
		{Kind: KindEssay, Count: 0}, // dạng không chọn thì bỏ qua
	})
	if err != nil {
		t.Fatalf("lỗi không mong đợi: %v", err)
	}
	if total != 3 || len(specs) != 1 {
		t.Fatalf("mong đợi 1 dạng / 3 câu, nhận được %d dạng / %d câu", len(specs), total)
	}

	if _, _, err := normalizeSpecs(nil); err == nil {
		t.Error("không chọn dạng nào thì phải báo lỗi")
	}
	if _, _, err := normalizeSpecs([]QuestionSpec{{Kind: "matching", Count: 1}}); err == nil {
		t.Error("dạng không hỗ trợ phải bị từ chối")
	}
	if _, _, err := normalizeSpecs([]QuestionSpec{{Kind: KindEssay, Count: 99}}); err == nil {
		t.Error("yêu cầu quá 50 câu phải bị từ chối")
	}
}

func TestNormalizeLevel(t *testing.T) {
	cases := map[string]string{
		"nhận biết":  "Nhận biết",
		"Thông hiểu": "Thông hiểu",
		"apply":      "Vận dụng",
		"linh tinh":  "",
	}
	for in, want := range cases {
		if got := normalizeLevel(in); got != want {
			t.Errorf("normalizeLevel(%q) = %q, mong đợi %q", in, got, want)
		}
	}
}

func TestExtractJSONUnwrapsCodeFence(t *testing.T) {
	var out struct {
		A int `json:"a"`
	}
	if err := extractJSON("```json\n{\"a\": 5}\n```", &out); err != nil {
		t.Fatalf("lỗi không mong đợi: %v", err)
	}
	if out.A != 5 {
		t.Errorf("a = %d, mong đợi 5", out.A)
	}
	if err := extractJSON("không phải JSON", &out); err == nil {
		t.Error("chuỗi không phải JSON phải báo lỗi")
	}
}
