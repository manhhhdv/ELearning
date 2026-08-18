// Package models chứa các kiểu dữ liệu dùng chung giữa tầng lưu trữ và tầng HTTP.
package models

import (
	"time"

	"github.com/google/uuid"
)

// Vai trò người dùng trong hệ thống.
const (
	RoleAdmin      = "admin"
	RoleTrainer    = "trainer"
	RoleSupervisor = "supervisor"
	RoleStudent    = "student"
)

// Loại nút trên cây chương trình đào tạo.
const (
	KindFolder     = "folder"
	KindLesson     = "lesson"
	KindAssignment = "assignment"
)

// Loại câu hỏi trong bài tập.
const (
	QuestionSingleChoice = "single_choice"
	QuestionMultiChoice  = "multi_choice"
	QuestionEssay        = "essay"
)

type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	FullName           string     `json:"fullName"`
	AvatarURL          string     `json:"avatarUrl"`
	Role               string     `json:"role"`
	IsActive           bool       `json:"isActive"`
	MustChangePassword bool       `json:"mustChangePassword"`
	HasPassword        bool       `json:"hasPassword"`
	HasGoogle          bool       `json:"hasGoogle"`
	LastLoginAt        *time.Time `json:"lastLoginAt"`
	CreatedAt          time.Time  `json:"createdAt"`
}

type Program struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	// Dùng để tạo URL thân thiện (VD: /hoc/attp-2026), lấy từ Code lúc tạo.
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	Status      string `json:"status"`
	// Cho phép học viên tự bấm ghi danh thay vì chờ admin thêm vào.
	AllowSelfEnroll bool       `json:"allowSelfEnroll"`
	CreatedBy       *uuid.UUID `json:"createdBy"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`

	// Thống kê phụ trợ cho màn hình danh sách.
	NodeCount       int `json:"nodeCount"`
	LessonCount     int `json:"lessonCount"`
	AssignmentCount int `json:"assignmentCount"`
	EnrollmentCount int `json:"enrollmentCount"`
	// Số bài học người dùng hiện tại đã hoàn thành; chỉ có giá trị ở API của học viên.
	CompletedLessonCount int `json:"completedLessonCount"`
	// Người đang xem đã ghi danh chưa; dùng cho trang khám phá khoá học.
	Enrolled bool `json:"enrolled"`
}

// Node là một nút bất kỳ trên cây: thư mục, bài học hoặc bài tập.
type Node struct {
	ID        uuid.UUID  `json:"id"`
	ProgramID uuid.UUID  `json:"programId"`
	ParentID  *uuid.UUID `json:"parentId"`
	Kind      string     `json:"kind"`
	// Duy nhất trong phạm vi chương trình, sinh từ Title lúc tạo.
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Lesson     *Lesson     `json:"lesson,omitempty"`
	Assignment *Assignment `json:"assignment,omitempty"`

	// Trạng thái hoàn thành của người dùng đang gọi API cây nội dung.
	Completed bool `json:"completed,omitempty"`

	Children []*Node `json:"children"`
}

type Lesson struct {
	ContentType     string `json:"contentType"`
	DriveFileID     string `json:"driveFileId"`
	EmbedURL        string `json:"embedUrl"`
	DurationMinutes int    `json:"durationMinutes"`
	Body            string `json:"body"`
}

type Assignment struct {
	Instructions     string     `json:"instructions"`
	TimeLimitMinutes int        `json:"timeLimitMinutes"`
	MaxAttempts      int        `json:"maxAttempts"`
	PassScore        float64    `json:"passScore"`
	ShuffleQuestions bool       `json:"shuffleQuestions"`
	DueAt            *time.Time `json:"dueAt"`

	QuestionCount int         `json:"questionCount"`
	Questions     []*Question `json:"questions,omitempty"`
}

type Question struct {
	ID uuid.UUID `json:"id"`
	// Mã ngắn cố định do người soạn đặt, không đổi khi sắp xếp lại thứ tự.
	Code        string            `json:"code"`
	Type        string            `json:"type"`
	Prompt      string            `json:"prompt"`
	Points      float64           `json:"points"`
	Position    int               `json:"position"`
	Explanation string            `json:"explanation"`
	Options     []*QuestionOption `json:"options"`
}

type QuestionOption struct {
	ID       uuid.UUID `json:"id"`
	Content  string    `json:"content"`
	Position int       `json:"position"`
	// Bị lược bỏ khi trả về cho học viên đang làm bài.
	IsCorrect bool `json:"isCorrect"`
}

type Enrollment struct {
	ID         uuid.UUID `json:"id"`
	ProgramID  uuid.UUID `json:"programId"`
	UserID     uuid.UUID `json:"userId"`
	Role       string    `json:"role"`
	EnrolledAt time.Time `json:"enrolledAt"`

	Email        string `json:"email"`
	FullName     string `json:"fullName"`
	ProgramTitle string `json:"programTitle"`
	ProgramCode  string `json:"programCode"`
}

type Submission struct {
	ID           uuid.UUID  `json:"id"`
	AssignmentID uuid.UUID  `json:"assignmentId"`
	UserID       uuid.UUID  `json:"userId"`
	AttemptNo    int        `json:"attemptNo"`
	Status       string     `json:"status"`
	AutoScore    float64    `json:"autoScore"`
	ManualScore  *float64   `json:"manualScore"`
	MaxScore     float64    `json:"maxScore"`
	Feedback     string     `json:"feedback"`
	GradedBy     *uuid.UUID `json:"gradedBy"`
	GradedAt     *time.Time `json:"gradedAt"`
	SubmittedAt  time.Time  `json:"submittedAt"`

	// Thông tin hiển thị kèm ở màn hình chấm bài.
	StudentName     string `json:"studentName,omitempty"`
	StudentEmail    string `json:"studentEmail,omitempty"`
	AssignmentTitle string `json:"assignmentTitle,omitempty"`
	ProgramTitle    string `json:"programTitle,omitempty"`
	NeedsGrading    bool   `json:"needsGrading"`

	Answers []*SubmissionAnswer `json:"answers,omitempty"`
}

// TotalScore là tổng điểm cuối cùng: phần tự động cộng phần chấm tay (nếu đã chấm).
func (s *Submission) TotalScore() float64 {
	if s.ManualScore == nil {
		return s.AutoScore
	}
	return s.AutoScore + *s.ManualScore
}

type SubmissionAnswer struct {
	ID                uuid.UUID   `json:"id"`
	QuestionID        uuid.UUID   `json:"questionId"`
	SelectedOptionIDs []uuid.UUID `json:"selectedOptionIds"`
	EssayText         string      `json:"essayText"`
	IsCorrect         *bool       `json:"isCorrect"`
	Score             float64     `json:"score"`
	Comment           string      `json:"comment"`

	Question *Question `json:"question,omitempty"`
}
