package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/manhnv/elearning/backend/internal/models"
)

// TestSampleCoursesWellFormed kiểm tra dữ liệu mẫu lớp 1–12 đủ và đúng khuôn:
// mỗi môn có đủ 12 lớp, mã không trùng, mỗi câu hỏi có đúng một đáp án đúng,
// mức độ hợp lệ và tổng điểm bài ôn tập là 10.
func TestSampleCoursesWellFormed(t *testing.T) {
	courses := sampleCourses()

	want := map[string]bool{}
	for g := 1; g <= 12; g++ {
		want[fmt.Sprintf("TOAN%d", g)] = true
		want[fmt.Sprintf("TIENGANH%d", g)] = true
		if g <= 5 {
			want[fmt.Sprintf("TIENGVIET%d", g)] = true
		} else {
			want[fmt.Sprintf("NGUVAN%d", g)] = true
		}
	}

	levels := map[string]bool{models.LevelRemember: true, models.LevelUnderstand: true, models.LevelApply: true}
	seen := map[string]bool{}
	for _, c := range courses {
		if seen[c.Code] {
			t.Errorf("mã lớp trùng: %s", c.Code)
		}
		seen[c.Code] = true
		if !want[c.Code] {
			t.Errorf("mã lớp ngoài danh sách: %s", c.Code)
		}
		if c.Title == "" || c.Description == "" {
			t.Errorf("%s: thiếu tên hoặc mô tả", c.Code)
		}
		if _, err := os.Stat(filepath.Join("..", "..", "..", "frontend", "public", "covers", c.Cover+".svg")); err != nil {
			t.Errorf("%s: không có ảnh bìa %q", c.Code, c.Cover)
		}
		if len(c.Chapters) == 0 {
			t.Errorf("%s: không có chương nào", c.Code)
		}
		for _, ch := range c.Chapters {
			if ch.Title == "" || len(ch.Lessons) == 0 {
				t.Errorf("%s: chương %q rỗng", c.Code, ch.Title)
			}
			for _, l := range ch.Lessons {
				if l.Title == "" || strings.TrimSpace(l.Body) == "" || l.DurationMinutes <= 0 {
					t.Errorf("%s: bài giảng %q thiếu nội dung", c.Code, l.Title)
				}
			}
		}

		if c.Quiz.Title == "" || len(c.Quiz.Questions) == 0 {
			t.Errorf("%s: thiếu bài ôn tập", c.Code)
		}
		var total float64
		for _, q := range c.Quiz.Questions {
			total += q.Points
			if !levels[q.Level] {
				t.Errorf("%s: câu %q có mức độ không hợp lệ %q", c.Code, q.Prompt, q.Level)
			}
			if q.Explanation == "" {
				t.Errorf("%s: câu %q thiếu lời giải thích", c.Code, q.Prompt)
			}
			correct := 0
			for _, o := range q.Options {
				if o.IsCorrect {
					correct++
				}
			}
			if len(q.Options) < 2 || correct != 1 {
				t.Errorf("%s: câu %q cần ≥ 2 phương án và đúng 1 đáp án đúng (có %d)", c.Code, q.Prompt, correct)
			}
		}
		if total != 10 {
			t.Errorf("%s: tổng điểm bài ôn tập là %v, cần 10", c.Code, total)
		}
	}

	for code := range want {
		if !seen[code] {
			t.Errorf("thiếu lớp %s", code)
		}
	}
}
