// Package models chứa các kiểu dữ liệu dùng chung giữa tầng lưu trữ và tầng HTTP.
package models

import (
	"encoding/json"
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

// Loại câu hỏi trong bài tập. Bốn dạng đầu là bốn dạng chuẩn của hệ thống
// (Bảng 3.1); QuestionMultiChoice được giữ lại cho dữ liệu đã soạn từ trước.
const (
	QuestionSingleChoice = "single_choice"
	QuestionTrueFalse    = "true_false"
	QuestionFillBlank    = "fill_blank"
	QuestionEssay        = "essay"
	QuestionMultiChoice  = "multi_choice"
)

// Mức độ nhận thức của câu hỏi.
const (
	LevelRemember   = "Nhận biết"
	LevelUnderstand = "Thông hiểu"
	LevelApply      = "Vận dụng"
)

// AutoGraded cho biết dạng câu hỏi có được chấm tự động ngay khi nộp bài không.
// Chỉ câu tự luận cần giáo viên chấm (mục 3.4.3).
func AutoGraded(questionType string) bool { return questionType != QuestionEssay }

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
	AllowSelfEnroll bool `json:"allowSelfEnroll"`
	// Tự động hiện trong "Lớp học của tôi" của mọi người dùng, không cần ghi danh.
	IsDefaultCourse bool       `json:"isDefaultCourse"`
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
	// Người đang xem đã ghi danh chưa; dùng cho trang khám phá lớp học.
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
	IsLocked    bool      `json:"isLocked"`
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
	// Danh sách tài liệu tải về kèm theo bài học (dùng cho loại "materials").
	Attachments []LessonAttachment `json:"attachments"`
}

// LessonAttachment là một tài liệu đính kèm: tên hiển thị và link tải.
type LessonAttachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
	Code   string  `json:"code"`
	Type   string  `json:"type"`
	Prompt string  `json:"prompt"`
	Points float64 `json:"points"`
	// Mức độ nhận thức: Nhận biết / Thông hiểu / Vận dụng. Rỗng = chưa phân loại.
	Level       string `json:"level"`
	Position    int    `json:"position"`
	Explanation string `json:"explanation"`
	// Đáp án gợi ý và tiêu chí chấm, chỉ dùng cho câu tự luận.
	SampleAnswer string `json:"sampleAnswer"`
	Rubric       string `json:"rubric"`
	// Ý nghĩa của Options thay đổi theo Type:
	//   single_choice / multi_choice : các phương án, IsCorrect = phương án đúng
	//   true_false                   : mỗi dòng là một phát biểu, IsCorrect = giá trị Đúng/Sai
	//   fill_blank                   : mỗi dòng là một đáp án được chấp nhận (IsCorrect luôn true)
	//   essay                        : luôn rỗng
	Options []*QuestionOption `json:"options"`
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
	ID uuid.UUID `json:"id"`
	// Với true_false, đây là các phát biểu học sinh đánh dấu "Đúng".
	QuestionID        uuid.UUID   `json:"questionId"`
	SelectedOptionIDs []uuid.UUID `json:"selectedOptionIds"`
	// Dùng cho cả câu tự luận lẫn câu điền khuyết.
	EssayText string  `json:"essayText"`
	IsCorrect *bool   `json:"isCorrect"`
	Score     float64 `json:"score"`
	Comment   string  `json:"comment"`
	// Điểm và nhận xét do AI gợi ý cho câu tự luận; chỉ là tham khảo cho
	// giáo viên, điểm chính thức nằm ở Score.
	AIScore   *float64 `json:"aiScore"`
	AIComment string   `json:"aiComment"`

	Question *Question `json:"question,omitempty"`
}

// Material là một tài liệu trong kho dùng chung toàn hệ thống (mục 4.4).
// Khác với tài liệu đính kèm bài học, tài liệu ở đây không thuộc lớp nào.
type Material struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	URL         string    `json:"url"`
	// Rỗng khi tài liệu là link ngoài, không phải Google Drive.
	DriveFileID string     `json:"driveFileId"`
	Kind        string     `json:"kind"`
	IsPublished bool       `json:"isPublished"`
	CreatedBy   *uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`

	CreatedByName string `json:"createdByName"`
}

// Các loại tài liệu trong kho dùng chung.
const (
	MaterialPDF      = "pdf"
	MaterialSlide    = "slide"
	MaterialDocument = "document"
	MaterialVideo    = "video"
	MaterialLink     = "link"
)

// PasswordResetRequest là yêu cầu đặt lại mật khẩu do người dùng gửi từ trang
// đăng nhập. Hệ thống chưa gửi email nên admin là người cấp lại mật khẩu.
type PasswordResetRequest struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	// Rỗng khi email không khớp tài khoản nào trong hệ thống.
	UserID    *uuid.UUID `json:"userId"`
	Note      string     `json:"note"`
	Status    string     `json:"status"`
	HandledBy *uuid.UUID `json:"handledBy"`
	HandledAt *time.Time `json:"handledAt"`
	CreatedAt time.Time  `json:"createdAt"`

	FullName      string `json:"fullName"`
	HandledByName string `json:"handledByName"`
}

// ---------------------------------------------------------------------------
// Lớp chức năng AI
// ---------------------------------------------------------------------------

// LessonPlan là giáo án do AI sinh, đã được giáo viên duyệt và lưu lại.
// Content giữ nguyên cấu trúc JSON do lớp Service AI trả về.
type LessonPlan struct {
	ID              uuid.UUID       `json:"id"`
	OwnerID         uuid.UUID       `json:"ownerId"`
	ProgramID       *uuid.UUID      `json:"programId"`
	Subject         string          `json:"subject"`
	Grade           string          `json:"grade"`
	Topic           string          `json:"topic"`
	Objectives      string          `json:"objectives"`
	DurationMinutes int             `json:"durationMinutes"`
	Title           string          `json:"title"`
	Content         json.RawMessage `json:"content"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`

	OwnerName string `json:"ownerName,omitempty"`
}

// AITask là một dòng nhật ký tác vụ AI, hiển thị ở mục "Tác vụ AI gần đây".
type AITask struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Kind      string    `json:"kind"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Detail    string    `json:"detail"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"createdAt"`
}

// Các giá trị của AITask.Kind.
const (
	AIKindLessonPlan = "lesson_plan"
	AIKindQuestions  = "questions"
	AIKindGradeEssay = "grade_essay"
	AIKindLesson     = "lesson"
	AIKindChat       = "chat"
)

// AIConversation là một cuộc trò chuyện với trợ lý AI.
type AIConversation struct {
	ProgramID *uuid.UUID   `json:"programId"`
	NodeID    *uuid.UUID   `json:"nodeId"`
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"userId"`
	Title     string       `json:"title"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Messages  []*AIMessage `json:"messages,omitempty"`
}

type AIMessage struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
