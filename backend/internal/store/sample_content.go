package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// sampleLesson là một bài giảng mẫu (nội dung văn bản, không phải giáo án).
type sampleLesson struct {
	Title           string
	DurationMinutes int
	Body            string
}

// sampleOption là một phương án của câu hỏi trắc nghiệm mẫu.
type sampleOption struct {
	Content   string
	IsCorrect bool
}

// sampleQuestion là một câu hỏi trắc nghiệm mẫu (chỉ dùng single_choice cho đơn giản).
type sampleQuestion struct {
	Prompt      string
	Level       string
	Points      float64
	Explanation string
	Options     []sampleOption
}

// sampleQuiz là bài tập trắc nghiệm mẫu đi kèm các bài giảng của một môn.
type sampleQuiz struct {
	Title        string
	Instructions string
	Questions    []sampleQuestion
}

// sampleChapter là một chương của lớp học mẫu, tạo thành một thư mục chứa các bài giảng.
type sampleChapter struct {
	Title   string
	Lessons []sampleLesson
}

// sampleCourse là toàn bộ một lớp học mẫu: thông tin lớp, các chương bài giảng
// và một bài tập trắc nghiệm ôn tập cuối lớp.
type sampleCourse struct {
	Code        string // mã lớp, VD "TOAN1", "NGUVAN10", "TIENGANH12"
	Title       string // VD "Toán lớp 1"
	Description string
	Cover       string // tên ảnh bìa trong frontend/public/covers (math, literature, language…)
	Chapters    []sampleChapter
	Quiz        sampleQuiz
}

// sampleCourses gom các lớp học mẫu từ lớp 1 đến lớp 12 theo Chương trình
// GDPT 2018: Toán, Tiếng Việt (lớp 1–5) / Ngữ văn (lớp 6–12) và Tiếng Anh.
// Nội dung từng nhóm lớp nằm ở các file sample_*.go.
func sampleCourses() []sampleCourse {
	var all []sampleCourse
	for _, group := range [][]sampleCourse{
		sampleToanPrimary, sampleToanSecondary,
		sampleVanPrimary, sampleVanSecondary,
		sampleAnhPrimary, sampleAnhSecondary,
	} {
		all = append(all, group...)
	}
	return all
}

// ensureSampleContent thêm các chương, bài giảng và bài tập mẫu cho một lớp
// học mẫu nếu lớp đó hiện chưa có nội dung nào. Nếu chương trình đã có ít nhất
// một nút (do seed trước, hoặc do giáo viên tự soạn/xoá) thì bỏ qua hẳn — seed
// không bao giờ ghi đè nội dung đã tồn tại.
func (s *Store) ensureSampleContent(ctx context.Context, programID uuid.UUID, course sampleCourse) error {
	var nodeCount int
	if err := s.pool.QueryRow(ctx, `SELECT count(*) FROM nodes WHERE program_id = $1`, programID).
		Scan(&nodeCount); err != nil {
		return fmt.Errorf("đếm nội dung hiện có: %w", err)
	}
	if nodeCount > 0 {
		return nil
	}

	for _, chapter := range course.Chapters {
		folder, err := s.CreateNode(ctx, SaveNodeParams{
			ProgramID:   programID,
			Kind:        models.KindFolder,
			Title:       chapter.Title,
			IsPublished: true,
		})
		if err != nil {
			return fmt.Errorf("tạo chương mẫu %q: %w", chapter.Title, err)
		}
		for _, lesson := range chapter.Lessons {
			if _, err := s.CreateNode(ctx, SaveNodeParams{
				ProgramID:   programID,
				ParentID:    &folder.ID,
				Kind:        models.KindLesson,
				Title:       lesson.Title,
				IsPublished: true,
				Lesson: &models.Lesson{
					ContentType:     "richtext",
					Body:            lesson.Body,
					DurationMinutes: lesson.DurationMinutes,
					Attachments:     []models.LessonAttachment{},
				},
			}); err != nil {
				return fmt.Errorf("tạo bài giảng mẫu %q: %w", lesson.Title, err)
			}
		}
	}

	// Bài ôn tập nằm trong một chương riêng ở cuối: mục lục lớp học đưa các mục
	// ở cấp gốc lên trước mọi chương, nên để ở gốc thì bài ôn tập sẽ đứng đầu.
	quiz := course.Quiz
	reviewFolder, err := s.CreateNode(ctx, SaveNodeParams{
		ProgramID:   programID,
		Kind:        models.KindFolder,
		Title:       "Ôn tập",
		IsPublished: true,
	})
	if err != nil {
		return fmt.Errorf("tạo chương ôn tập: %w", err)
	}
	assignmentNode, err := s.CreateNode(ctx, SaveNodeParams{
		ProgramID:   programID,
		ParentID:    &reviewFolder.ID,
		Kind:        models.KindAssignment,
		Title:       quiz.Title,
		IsPublished: true,
		Assignment: &models.Assignment{
			Instructions: quiz.Instructions,
			MaxAttempts:  3,
			PassScore:    5,
		},
	})
	if err != nil {
		return fmt.Errorf("tạo bài tập mẫu %q: %w", quiz.Title, err)
	}

	for _, q := range quiz.Questions {
		opts := make([]QuestionOptionInput, len(q.Options))
		for i, o := range q.Options {
			opts[i] = QuestionOptionInput{Content: o.Content, IsCorrect: o.IsCorrect}
		}
		if _, err := s.CreateQuestion(ctx, SaveQuestionParams{
			AssignmentID: assignmentNode.ID,
			Type:         models.QuestionSingleChoice,
			Prompt:       q.Prompt,
			Points:       q.Points,
			Level:        q.Level,
			Explanation:  q.Explanation,
			Options:      opts,
		}); err != nil {
			return fmt.Errorf("tạo câu hỏi mẫu %q: %w", q.Prompt, err)
		}
	}

	return nil
}
