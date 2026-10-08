// Package export dựng tệp Word (.docx) cho giáo viên tải về: giáo án, bộ câu
// hỏi của bài tập (kèm hoặc không kèm đáp án) và bài giảng.
//
// Bố cục bám theo văn bản giáo viên Việt Nam quen dùng — giáo án theo khung kế
// hoạch bài dạy (Mục tiêu → Thiết bị → Tiến trình với bảng hoạt động của giáo
// viên / học sinh), đề kiểm tra có dòng họ tên và phần đáp án tách trang riêng
// — để tải về là in hoặc chỉnh tiếp được ngay.
package export

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/manhnv/elearning/backend/internal/ai"
	"github.com/manhnv/elearning/backend/internal/docx"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/util"
)

// ContentType là kiểu MIME của tệp .docx.
const ContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

// FileName tạo tên tệp không dấu, VD: FileName("giao an", "Định luật Ôm") → "giao-an-dinh-luat-om.docx".
func FileName(prefix, title string) string {
	return util.Slugify(prefix) + "-" + slugPart(title) + ".docx"
}

// ---------------------------------------------------------------------------
// Giáo án
// ---------------------------------------------------------------------------

// LessonPlanInfo là thông tin chung in ở đầu giáo án.
type LessonPlanInfo struct {
	Subject         string
	Grade           string
	DurationMinutes int
	Author          string
}

// LessonPlan dựng tệp Word cho một giáo án.
func LessonPlan(info LessonPlanInfo, plan *ai.LessonPlan) ([]byte, error) {
	title := strings.TrimSpace(plan.Title)
	if title == "" {
		title = "Giáo án"
	}
	d := docx.New(title)
	d.Author = info.Author
	body := d.Body()

	body.Text(docx.StyleTitle, "KẾ HOẠCH BÀI DẠY")
	body.Text(docx.StyleTitle, title)
	body.Text(docx.StyleSubtitle, joinNonEmpty(" · ",
		labelled("Môn học", info.Subject),
		labelled("Lớp", info.Grade),
		labelled("Thời lượng", minutes(info.DurationMinutes)),
	))

	section := sectionCounter()
	if items := nonEmpty(plan.Objectives); len(items) > 0 {
		body.Text(docx.StyleHeading1, section("YÊU CẦU CẦN ĐẠT"))
		bulletList(d, body, items)
	}
	if items := nonEmpty(plan.Materials); len(items) > 0 {
		body.Text(docx.StyleHeading1, section("THIẾT BỊ DẠY HỌC VÀ HỌC LIỆU"))
		bulletList(d, body, items)
	}

	if len(plan.Activities) > 0 {
		body.Text(docx.StyleHeading1, section("TIẾN TRÌNH DẠY HỌC"))
		for i, a := range plan.Activities {
			name := strings.TrimSpace(a.Name)
			if name == "" {
				name = fmt.Sprintf("Hoạt động %d", i+1)
			}
			if a.DurationMinutes > 0 {
				name += " (" + minutes(a.DurationMinutes) + ")"
			}
			body.Text(docx.StyleHeading2, name)

			part := partCounter()
			if s := strings.TrimSpace(a.Objective); s != "" {
				body.Add(labelPara(part("Mục tiêu")+": ", s))
			}
			if s := strings.TrimSpace(a.Output); s != "" {
				body.Add(labelPara(part("Sản phẩm")+": ", s))
			}
			body.Add(docx.Para{KeepNext: true, Runs: []docx.Run{{Text: part("Tổ chức thực hiện") + ":", Bold: true}}})
			t := body.NewTable(1, 1)
			t.HeaderRows = 1
			head := t.Row()
			head[0].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: "Hoạt động của giáo viên", Bold: true}}})
			head[1].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: "Hoạt động của học sinh", Bold: true}}})
			row := t.Row()
			row[0].Markdown(a.TeacherActions)
			row[1].Markdown(a.StudentActions)
			body.AddTable(t)
		}
	}

	if s := strings.TrimSpace(plan.Homework); s != "" {
		body.Text(docx.StyleHeading1, section("BÀI TẬP VỀ NHÀ"))
		body.Markdown(s)
	}
	if s := strings.TrimSpace(plan.Notes); s != "" {
		body.Text(docx.StyleHeading1, section("LƯU Ý KHI TỔ CHỨC DẠY HỌC"))
		body.Markdown(s)
	}
	return d.Bytes()
}

// ---------------------------------------------------------------------------
// Bộ câu hỏi
// ---------------------------------------------------------------------------

// QuestionSetInfo là thông tin chung in ở đầu đề.
type QuestionSetInfo struct {
	ProgramTitle     string
	Title            string
	Instructions     string
	TimeLimitMinutes int
	Author           string
}

var blankMark = regexp.MustCompile(`_{3,}`)

const dottedBlank = "…………………"

// Questions dựng tệp Word cho bộ câu hỏi của một bài tập. Khi withAnswers,
// phần đáp án và hướng dẫn chấm được in ở trang riêng sau đề, để giáo viên
// vẫn in được đề từ cùng một tệp.
func Questions(info QuestionSetInfo, questions []*models.Question, withAnswers bool) ([]byte, error) {
	title := strings.TrimSpace(info.Title)
	if title == "" {
		title = "Đề bài"
	}
	d := docx.New(title)
	d.Author = info.Author
	body := d.Body()

	total := 0.0
	for _, q := range questions {
		total += q.Points
	}
	body.Text(docx.StyleTitle, title)
	body.Text(docx.StyleSubtitle, joinNonEmpty(" · ",
		labelled("Lớp học", info.ProgramTitle),
		labelled("Thời gian làm bài", minutes(info.TimeLimitMinutes)),
		fmt.Sprintf("%d câu", len(questions)),
		"tổng "+points(total)+" điểm",
	))
	body.Add(docx.Para{
		Tabs: []docx.TabStop{{Pos: docx.ContentWidth * 7 / 10, DotLeader: true}, {Pos: docx.ContentWidth, Right: true, DotLeader: true}},
		Runs: []docx.Run{{Text: "Họ và tên:\t  Lớp:\t"}},
	})
	if s := strings.TrimSpace(info.Instructions); s != "" {
		body.Add(docx.Para{Runs: []docx.Run{{Text: "Hướng dẫn: ", Bold: true, Italic: true}, {Text: s, Italic: true}}})
	}
	body.Add(docx.Para{Style: docx.StyleRule})

	// Chỉ ghi điểm từng câu khi các câu có điểm khác nhau; đều nhau thì tổng
	// điểm ở đầu đề là đủ.
	showPoints := false
	for _, q := range questions {
		if q.Points != questions[0].Points {
			showPoints = true
			break
		}
	}

	for i, q := range questions {
		writeQuestion(body, i+1, q, showPoints)
	}

	if withAnswers && len(questions) > 0 {
		body.PageBreak()
		body.Text(docx.StyleTitle, "ĐÁP ÁN VÀ HƯỚNG DẪN CHẤM")
		body.Text(docx.StyleSubtitle, title)
		writeAnswerGrid(body, questions)
		for i, q := range questions {
			writeAnswer(body, i+1, q)
		}
	}
	return d.Bytes()
}

func writeQuestion(body *docx.Container, num int, q *models.Question, showPoints bool) {
	prompt := strings.TrimSpace(q.Prompt)
	if q.Type == models.QuestionFillBlank {
		prompt = blankMark.ReplaceAllString(prompt, dottedBlank)
	}
	runs := []docx.Run{{Text: fmt.Sprintf("Câu %d. ", num), Bold: true}, {Text: prompt}}
	if q.Type == models.QuestionMultiChoice {
		runs = append(runs, docx.Run{Text: " (Chọn tất cả phương án đúng)", Italic: true})
	}
	if showPoints {
		runs = append(runs, docx.Run{Text: " (" + points(q.Points) + " điểm)", Italic: true})
	}
	body.Add(docx.Para{SpaceBefore: 160, KeepNext: len(q.Options) > 0, Runs: runs})

	switch q.Type {
	case models.QuestionSingleChoice, models.QuestionMultiChoice:
		for i, o := range q.Options {
			body.Add(docx.Para{
				Indent:   567,
				KeepNext: i < len(q.Options)-1,
				Runs:     []docx.Run{{Text: optionLetter(i) + ". ", Bold: true}, {Text: strings.TrimSpace(o.Content)}},
			})
		}

	case models.QuestionTrueFalse:
		t := body.NewTable(8, 1, 1)
		t.HeaderRows = 1
		head := t.Row()
		head[0].Add(docx.Para{Runs: []docx.Run{{Text: "Phát biểu", Bold: true}}})
		head[1].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: "Đúng", Bold: true}}})
		head[2].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: "Sai", Bold: true}}})
		for i, o := range q.Options {
			row := t.Row()
			row[0].Add(docx.Para{Runs: []docx.Run{{Text: statementLetter(i) + ") ", Bold: true}, {Text: strings.TrimSpace(o.Content)}}})
		}
		body.AddTable(t)

	case models.QuestionFillBlank:
		if !blankMark.MatchString(q.Prompt) {
			body.Add(docx.Para{Indent: 567, Runs: []docx.Run{{Text: "Trả lời: " + dottedBlank + dottedBlank}}})
		}
	}
}

// writeAnswerGrid in bảng đáp án gọn cho các câu chọn phương án, mỗi hàng 10
// câu, để giáo viên dò nhanh khi chấm tay.
func writeAnswerGrid(body *docx.Container, questions []*models.Question) {
	type cell struct{ num, answer string }
	var cells []cell
	for i, q := range questions {
		if q.Type == models.QuestionSingleChoice || q.Type == models.QuestionMultiChoice {
			cells = append(cells, cell{strconv.Itoa(i + 1), strings.Join(correctLetters(q), ", ")})
		}
	}
	if len(cells) < 2 {
		return
	}

	const perRow = 10
	body.Add(docx.Para{Style: docx.StyleHeading2, Runs: []docx.Run{{Text: "Bảng đáp án phần chọn phương án"}}})
	weights := make([]int, perRow+1)
	for i := range weights {
		weights[i] = 2
	}
	weights[0] = 3
	t := body.NewTable(weights...)
	for start := 0; start < len(cells); start += perRow {
		nums, answers := t.Row(), t.Row()
		nums[0].Add(docx.Para{Runs: []docx.Run{{Text: "Câu", Bold: true}}})
		answers[0].Add(docx.Para{Runs: []docx.Run{{Text: "Đáp án", Bold: true}}})
		for j := 0; j < perRow && start+j < len(cells); j++ {
			c := cells[start+j]
			nums[j+1].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: c.num, Bold: true}}})
			answers[j+1].Add(docx.Para{Align: docx.AlignCenter, Runs: []docx.Run{{Text: c.answer}}})
		}
	}
	body.AddTable(t)
	body.Text(docx.StyleHeading2, "Đáp án chi tiết")
}

func writeAnswer(body *docx.Container, num int, q *models.Question) {
	head := []docx.Run{{Text: fmt.Sprintf("Câu %d", num), Bold: true}}
	if meta := joinNonEmpty(" · ", q.Level, points(q.Points)+" điểm"); meta != "" {
		head = append(head, docx.Run{Text: " (" + meta + ")", Italic: true})
	}
	body.Add(docx.Para{SpaceBefore: 160, KeepNext: true, Runs: head})

	switch q.Type {
	case models.QuestionSingleChoice, models.QuestionMultiChoice:
		answer := strings.Join(correctLetters(q), ", ")
		if answer == "" {
			answer = "(chưa đánh dấu phương án đúng)"
		}
		body.Add(labelPara("Đáp án: ", answer))

	case models.QuestionTrueFalse:
		parts := make([]string, len(q.Options))
		for i, o := range q.Options {
			v := "Sai"
			if o.IsCorrect {
				v = "Đúng"
			}
			parts[i] = statementLetter(i) + ") " + v
		}
		body.Add(labelPara("Đáp án: ", strings.Join(parts, "; ")))

	case models.QuestionFillBlank:
		accepted := make([]string, 0, len(q.Options))
		for _, o := range q.Options {
			if s := strings.TrimSpace(o.Content); s != "" {
				accepted = append(accepted, s)
			}
		}
		body.Add(labelPara("Đáp án được chấp nhận: ", strings.Join(accepted, "; ")))

	case models.QuestionEssay:
		if s := strings.TrimSpace(q.SampleAnswer); s != "" {
			body.Add(labelPara("Đáp án gợi ý: ", s))
		}
		if s := strings.TrimSpace(q.Rubric); s != "" {
			body.Add(labelPara("Hướng dẫn chấm: ", s))
		}
	}

	if s := strings.TrimSpace(q.Explanation); s != "" {
		body.Add(docx.Para{Runs: []docx.Run{{Text: "Lời giải: ", Bold: true, Italic: true}, {Text: s, Italic: true}}})
	}
}

func correctLetters(q *models.Question) []string {
	var letters []string
	for i, o := range q.Options {
		if o.IsCorrect {
			letters = append(letters, optionLetter(i))
		}
	}
	return letters
}

// optionLetter: 0 → A, 1 → B… (quá 26 phương án thì dùng số).
func optionLetter(i int) string {
	if i < 26 {
		return string(rune('A' + i))
	}
	return strconv.Itoa(i + 1)
}

// statementLetter: 0 → a, 1 → b… cho các phát biểu của câu đúng/sai.
func statementLetter(i int) string {
	if i < 26 {
		return string(rune('a' + i))
	}
	return strconv.Itoa(i + 1)
}

// ---------------------------------------------------------------------------
// Bài giảng
// ---------------------------------------------------------------------------

var contentTypeLabel = map[string]string{
	"video":     "Video",
	"slide":     "Slide trình chiếu",
	"document":  "Tài liệu",
	"pdf":       "PDF",
	"link":      "Liên kết ngoài",
	"materials": "Tài liệu tải về",
}

// Lesson dựng tệp Word cho một bài học. Bài đọc tự soạn được chuyển trọn nội
// dung; các loại nhúng (video, slide, PDF…) chỉ ghi lại liên kết và ghi chú.
func Lesson(programTitle, author string, node *models.Node) ([]byte, error) {
	title := strings.TrimSpace(node.Title)
	if title == "" {
		title = "Bài giảng"
	}
	d := docx.New(title)
	d.Author = author
	body := d.Body()

	lesson := node.Lesson
	if lesson == nil {
		lesson = &models.Lesson{}
	}
	body.Text(docx.StyleTitle, title)
	if meta := joinNonEmpty(" · ", labelled("Lớp học", programTitle), labelled("Thời lượng", minutes(lesson.DurationMinutes))); meta != "" {
		body.Text(docx.StyleSubtitle, meta)
	}
	if s := strings.TrimSpace(node.Description); s != "" {
		body.Add(docx.Para{Runs: []docx.Run{{Text: s, Italic: true}}})
	}

	if lesson.ContentType != "richtext" && lesson.ContentType != "materials" {
		url := strings.TrimSpace(lesson.EmbedURL)
		if url == "" && lesson.DriveFileID != "" {
			url = "https://drive.google.com/file/d/" + lesson.DriveFileID + "/view"
		}
		if url != "" {
			label := contentTypeLabel[lesson.ContentType]
			if label == "" {
				label = "Nội dung"
			}
			body.Add(docx.Para{Runs: []docx.Run{{Text: label + ": ", Bold: true}, {Text: url, Link: url}}})
		}
	}

	if s := strings.TrimSpace(lesson.Body); s != "" {
		body.Markdown(s)
	}

	if len(lesson.Attachments) > 0 {
		body.Text(docx.StyleHeading2, "Tài liệu kèm theo")
		list := d.NewList(false, 1)
		for _, a := range lesson.Attachments {
			name := strings.TrimSpace(a.Name)
			if name == "" {
				name = a.URL
			}
			body.Add(docx.Para{List: list, Runs: []docx.Run{{Text: name, Link: a.URL}}})
		}
	}
	return d.Bytes()
}

// ---------------------------------------------------------------------------
// Tiện ích chung
// ---------------------------------------------------------------------------

// sectionCounter đánh số mục lớn bằng số La Mã: I. …, II. …
func sectionCounter() func(string) string {
	roman := []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	n := 0
	return func(title string) string {
		n++
		if n <= len(roman) {
			return roman[n-1] + ". " + title
		}
		return strconv.Itoa(n) + ". " + title
	}
}

// partCounter đánh số ý nhỏ: a) …, b) …
func partCounter() func(string) string {
	n := 0
	return func(title string) string {
		n++
		return statementLetter(n-1) + ") " + title
	}
}

func bulletList(d *docx.Document, c *docx.Container, items []string) {
	list := d.NewList(false, 1)
	for _, s := range items {
		c.Add(docx.Para{List: list, Runs: []docx.Run{{Text: s}}})
	}
}

// labelPara là đoạn văn dạng "Nhãn in đậm: nội dung".
func labelPara(label, text string) docx.Para {
	return docx.Para{Runs: []docx.Run{{Text: label, Bold: true}, {Text: text}}}
}

func labelled(label, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return label + ": " + value
}

func minutes(n int) string {
	if n <= 0 {
		return ""
	}
	return fmt.Sprintf("%d phút", n)
}

// points in điểm theo cách viết số thập phân của Việt Nam: 0.5 → "0,5".
func points(p float64) string {
	return strings.ReplaceAll(strconv.FormatFloat(p, 'f', -1, 64), ".", ",")
}

func joinNonEmpty(sep string, parts ...string) string {
	return strings.Join(nonEmpty(parts), sep)
}

func nonEmpty(items []string) []string {
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
