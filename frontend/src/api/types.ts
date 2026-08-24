// Các kiểu dữ liệu khớp với JSON do backend Go trả về.

export type Role = 'admin' | 'trainer' | 'supervisor' | 'student'
export type NodeKind = 'folder' | 'lesson' | 'assignment'
export type ContentType = 'video' | 'slide' | 'document' | 'pdf' | 'link' | 'richtext' | 'materials'
export type QuestionType = 'single_choice' | 'multi_choice' | 'essay'
export type ProgramStatus = 'draft' | 'published' | 'archived'

export interface User {
  id: string
  email: string
  fullName: string
  avatarUrl: string
  role: Role
  isActive: boolean
  mustChangePassword: boolean
  hasPassword: boolean
  hasGoogle: boolean
  lastLoginAt: string | null
  createdAt: string
}

export interface Program {
  id: string
  code: string
  /** Dùng cho URL thân thiện (VD: /hoc/attp-2026), lấy từ mã lúc tạo. */
  slug: string
  title: string
  description: string
  coverUrl: string
  status: ProgramStatus
  /** Cho phép học viên tự bấm ghi danh thay vì chờ admin thêm vào. */
  allowSelfEnroll: boolean
  /** Tự động hiện trong "Khoá học của tôi" của mọi người dùng, không cần ghi danh. Chỉ admin đặt được. */
  isDefaultCourse: boolean
  createdBy: string | null
  createdAt: string
  updatedAt: string
  nodeCount: number
  lessonCount: number
  assignmentCount: number
  enrollmentCount: number
  /** Số bài học người dùng hiện tại đã hoàn thành; chỉ có giá trị ở API của học viên. */
  completedLessonCount: number
  /** Người đang xem đã ghi danh chưa; dùng ở trang khám phá khoá học. */
  enrolled: boolean
}

export interface Lesson {
  contentType: ContentType
  driveFileId: string
  embedUrl: string
  durationMinutes: number
  body: string
  /** Tài liệu tải về kèm bài học; dùng cho loại nội dung "materials". */
  attachments: LessonAttachment[]
}

/** Một tài liệu tải về: tên hiển thị và link tải. */
export interface LessonAttachment {
  name: string
  url: string
}

export interface Assignment {
  instructions: string
  timeLimitMinutes: number
  maxAttempts: number
  passScore: number
  shuffleQuestions: boolean
  dueAt: string | null
  questionCount: number
  questions?: Question[]
}

export interface TreeNode {
  id: string
  programId: string
  parentId: string | null
  kind: NodeKind
  /** Duy nhất trong phạm vi chương trình, dùng cho URL thay vì UUID. */
  slug: string
  title: string
  description: string
  position: number
  isPublished: boolean
  isLocked: boolean
  createdAt: string
  updatedAt: string
  lesson?: Lesson
  assignment?: Assignment
  completed?: boolean
  children: TreeNode[]
}

export interface QuestionOption {
  id: string
  content: string
  position: number
  isCorrect: boolean
}

export interface Question {
  id: string
  /** Mã ngắn cố định, không đổi khi sắp xếp lại thứ tự câu hỏi. */
  code: string
  type: QuestionType
  prompt: string
  points: number
  position: number
  explanation: string
  options: QuestionOption[]
}

export interface Enrollment {
  id: string
  programId: string
  userId: string
  role: 'student' | 'trainer'
  enrolledAt: string
  email: string
  fullName: string
}

export interface SubmissionAnswer {
  id: string
  questionId: string
  selectedOptionIds: string[]
  essayText: string
  isCorrect: boolean | null
  score: number
  comment: string
}

export interface Submission {
  id: string
  assignmentId: string
  userId: string
  attemptNo: number
  status: 'submitted' | 'graded'
  autoScore: number
  manualScore: number | null
  maxScore: number
  feedback: string
  gradedBy: string | null
  gradedAt: string | null
  submittedAt: string
  studentName?: string
  studentEmail?: string
  assignmentTitle?: string
  programTitle?: string
  needsGrading: boolean
  answers?: SubmissionAnswer[]
}

/** Lượt làm bài đang mở; đồng hồ đếm ngược dựa trên expiresAt do máy chủ cấp. */
export interface AttemptSession {
  id: string
  startedAt: string
  expiresAt: string | null
}

export interface AttemptView {
  node: TreeNode
  attemptsUsed: number
  maxAttempts: number
  submissions: Submission[]
  session: AttemptSession | null
}

export interface SubmissionDetail {
  submission: Submission
  questions: Question[]
  canGrade: boolean
}

/** Thống kê một câu hỏi trên toàn bộ bài đã nộp. */
export interface QuestionStat {
  questionId: string
  code: string
  type: QuestionType
  prompt: string
  points: number
  position: number
  answerCount: number
  correctCount: number
  blankCount: number
  averageScore: number
  needsGrading: number
}

export interface AssignmentResults {
  submissionCount: number
  studentCount: number
  pendingCount: number
  maxScore: number
  averageScore: number
  passScore: number
  passedCount: number
  questions: QuestionStat[]
}

/** Số liệu tổng quan toàn hệ thống cho trang chủ quản trị. */
export interface DashboardStats {
  programsTotal: number
  programsDraft: number
  programsPublished: number
  programsArchived: number
  usersTotal: number
  adminCount: number
  trainerCount: number
  supervisorCount: number
  studentCount: number
  enrollmentsTotal: number
  submissionsTotal: number
  submissionsPending: number
  lessonsCompleted: number
  topPrograms: DashboardProgram[]
  recentSignups: DashboardUser[]
}

export interface DashboardProgram {
  id: string
  slug: string
  title: string
  code: string
  status: ProgramStatus
  enrollmentCount: number
}

export interface DashboardUser {
  id: string
  fullName: string
  email: string
  role: Role
  createdAt: string
}

export interface GoogleSettings {
  enabled: boolean
  clientId: string
  hasSecret: boolean
  source: 'database' | 'env' | 'none'
  redirectUrl: string
  allowedDomains: string
  autoProvisionRole: '' | 'student' | 'trainer'
}

// Nhãn tiếng Việt dùng chung cho giao diện.
export const ROLE_LABEL: Record<Role, string> = {
  admin: 'Quản trị viên',
  trainer: 'Giảng viên',
  supervisor: 'Giám sát',
  student: 'Học viên',
}

export const STATUS_LABEL: Record<ProgramStatus, string> = {
  draft: 'Bản nháp',
  published: 'Đã xuất bản',
  archived: 'Lưu trữ',
}

export const CONTENT_TYPE_LABEL: Record<ContentType, string> = {
  video: 'Video',
  slide: 'Slide trình chiếu',
  document: 'Tài liệu',
  pdf: 'PDF',
  link: 'Liên kết ngoài',
  richtext: 'Bài đọc tự soạn',
  materials: 'Tài liệu tải về',
}

export const QUESTION_TYPE_LABEL: Record<QuestionType, string> = {
  single_choice: 'Trắc nghiệm — một đáp án',
  multi_choice: 'Trắc nghiệm — nhiều đáp án',
  essay: 'Tự luận',
}
