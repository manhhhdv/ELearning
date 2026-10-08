package ai

import (
	"context"
	"fmt"
	"strings"
)

// LessonPlanParams là thông tin bài học giáo viên nhập ở màn hình tạo giáo án.
type LessonPlanParams struct {
	Subject         string
	Grade           string
	Topic           string
	Content         string
	Objectives      string
	DurationMinutes int
}

// LessonPlan là bản nháp giáo án, tổ chức theo từng hoạt động dạy học.
type LessonPlan struct {
	Title      string           `json:"title"`
	Objectives []string         `json:"objectives"`
	Materials  []string         `json:"materials"`
	Activities []LessonActivity `json:"activities"`
	Homework   string           `json:"homework"`
	Notes      string           `json:"notes"`
}

type LessonActivity struct {
	Name            string `json:"name"`
	DurationMinutes int    `json:"durationMinutes"`
	Objective       string `json:"objective"`
	TeacherActions  string `json:"teacherActions"`
	StudentActions  string `json:"studentActions"`
	Output          string `json:"output"`
}

const lessonPlanSystemPrompt = `Bạn là chuyên gia thiết kế bài dạy theo Chương trình giáo dục phổ thông 2018 của Việt Nam.
Bạn luôn trả lời bằng tiếng Việt và luôn trả về đúng một đối tượng JSON, không kèm lời dẫn hay chú thích.`

// GenerateLessonPlan gọi Gemini và kiểm tra bản nháp giáo án trước khi trả về.
func (c *Client) GenerateLessonPlan(ctx context.Context, p LessonPlanParams) (*LessonPlan, int, error) {
	if !c.Enabled() {
		return nil, 0, ErrDisabled
	}
	if strings.TrimSpace(p.Topic) == "" {
		return nil, 0, fmt.Errorf("vui lòng nhập chủ đề bài học")
	}
	if p.DurationMinutes <= 0 {
		p.DurationMinutes = 45
	}

	prompt := buildLessonPlanPrompt(p)
	var lastErr error

	for attempt := 1; attempt <= maxGenerateAttempts; attempt++ {
		text := prompt
		if lastErr != nil {
			text = prompt + "\n\nLƯU Ý: Kết quả lần trước bị loại vì: " + lastErr.Error() + "\nHãy soạn lại cho đúng yêu cầu."
		}
		raw, err := c.Generate(ctx, Request{
			System:      lessonPlanSystemPrompt,
			Contents:    []Content{UserText(text)},
			JSON:        true,
			Temperature: 0.6,
		})
		if err != nil {
			return nil, 0, err
		}

		var plan LessonPlan
		if err := extractJSON(raw, &plan); err != nil {
			lastErr = err
			continue
		}
		if err := normalizeLessonPlan(&plan, p); err != nil {
			lastErr = err
			continue
		}
		return &plan, attempt, nil
	}
	return nil, 0, fmt.Errorf("AI đã soạn lại %d lần nhưng giáo án vẫn không đạt yêu cầu (%v)", maxGenerateAttempts, lastErr)
}

func buildLessonPlanPrompt(p LessonPlanParams) string {
	b := &promptBuilder{
		role: fmt.Sprintf("Bạn là giáo viên %s giàu kinh nghiệm, đang soạn giáo án cho học sinh %s.",
			orDefault(p.Subject, "bộ môn"), orDefault(p.Grade, "trung học phổ thông")),
	}
	b.addContext("Môn học", p.Subject)
	b.addContext("Lớp", p.Grade)
	b.addContext("Chủ đề bài học", p.Topic)
	b.addContext("Nội dung trọng tâm", p.Content)
	b.addContext("Yêu cầu cần đạt", p.Objectives)
	b.addContext("Thời lượng", fmt.Sprintf("%d phút", p.DurationMinutes))

	b.addTask("Soạn một giáo án hoàn chỉnh cho bài học trên.")
	b.addTask("Chia bài dạy thành các hoạt động: Khởi động, Hình thành kiến thức, Luyện tập, Vận dụng.")
	b.addTask("Với mỗi hoạt động, nêu rõ mục tiêu, hoạt động của giáo viên, hoạt động của học sinh và sản phẩm cần đạt.")
	b.addTask("Đề xuất bài tập về nhà và những lưu ý sư phạm khi dạy bài này.")

	b.constraints = []string{
		"Viết toàn bộ bằng tiếng Việt, dùng thuật ngữ sư phạm phổ thông.",
		fmt.Sprintf("Tổng thời lượng các hoạt động phải bằng đúng %d phút.", p.DurationMinutes),
		"Nội dung phải chính xác về kiến thức và bám sát Chương trình giáo dục phổ thông 2018.",
		"Mỗi hoạt động mô tả cụ thể, giáo viên đọc là dạy được ngay, không nói chung chung.",
	}
	b.format = `Chỉ trả về JSON theo đúng cấu trúc sau, không thêm ký tự nào khác:
{
  "title": "Tên bài dạy",
  "objectives": ["Yêu cầu cần đạt 1", "Yêu cầu cần đạt 2"],
  "materials": ["Thiết bị, học liệu cần chuẩn bị"],
  "activities": [
    {
      "name": "Hoạt động 1: Khởi động",
      "durationMinutes": 5,
      "objective": "Mục tiêu của hoạt động",
      "teacherActions": "Giáo viên làm gì",
      "studentActions": "Học sinh làm gì",
      "output": "Sản phẩm học tập cần đạt"
    }
  ],
  "homework": "Bài tập về nhà",
  "notes": "Lưu ý khi tổ chức dạy học"
}`
	return b.String()
}

// normalizeLessonPlan là lớp kiểm tra tự động cho giáo án: phải có tiêu đề,
// yêu cầu cần đạt và các hoạt động mô tả đầy đủ.
func normalizeLessonPlan(plan *LessonPlan, p LessonPlanParams) error {
	plan.Title = strings.TrimSpace(plan.Title)
	if plan.Title == "" {
		plan.Title = strings.TrimSpace(p.Topic)
	}
	plan.Objectives = trimAll(plan.Objectives)
	plan.Materials = trimAll(plan.Materials)
	plan.Homework = strings.TrimSpace(plan.Homework)
	plan.Notes = strings.TrimSpace(plan.Notes)

	if len(plan.Objectives) == 0 {
		return fmt.Errorf("giáo án thiếu phần yêu cầu cần đạt")
	}

	acts := make([]LessonActivity, 0, len(plan.Activities))
	for _, a := range plan.Activities {
		a.Name = strings.TrimSpace(a.Name)
		a.Objective = strings.TrimSpace(a.Objective)
		a.TeacherActions = strings.TrimSpace(a.TeacherActions)
		a.StudentActions = strings.TrimSpace(a.StudentActions)
		a.Output = strings.TrimSpace(a.Output)
		if a.Name == "" {
			continue
		}
		if a.TeacherActions == "" || a.StudentActions == "" {
			return fmt.Errorf("hoạt động %q thiếu mô tả việc của giáo viên hoặc của học sinh", a.Name)
		}
		if a.DurationMinutes < 0 {
			a.DurationMinutes = 0
		}
		acts = append(acts, a)
	}
	if len(acts) < 2 {
		return fmt.Errorf("giáo án cần ít nhất 2 hoạt động dạy học")
	}
	plan.Activities = acts
	return nil
}
