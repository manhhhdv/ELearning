// Các kiểu dữ liệu khớp với JSON do backend Go trả về.

export type Role = 'admin' | 'trainer' | 'supervisor' | 'student'
export type NodeKind = 'folder' | 'lesson' | 'assignment'
export type ContentType = 'video' | 'slide' | 'document' | 'pdf' | 'link' | 'richtext' | 'materials'
/**
 * Bốn dạng câu hỏi chuẩn của hệ thống. `multi_choice` được giữ lại để hiển thị
 * đúng các câu đã soạn trước khi có bốn dạng chuẩn.
 */
export type QuestionType = 'single_choice' | 'true_false' | 'fill_blank' | 'essay' | 'multi_choice'

/** Mức độ nhận thức của câu hỏi. */
export type QuestionLevel = 'Nhận biết' | 'Thông hiểu' | 'Vận dụng' | ''

export const QUESTION_LEVELS: Exclude<QuestionLevel, ''>[] = ['Nhận biết', 'Thông hiểu', 'Vận dụng']

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
  /** Cho phép học sinh tự bấm ghi danh thay vì chờ admin thêm vào. */
  allowSelfEnroll: boolean
  /** Tự động hiện trong "Lớp học của tôi" của mọi người dùng, không cần ghi danh. Chỉ admin đặt được. */
  isDefaultCourse: boolean
  createdBy: string | null
  createdAt: string
  updatedAt: string
  nodeCount: number
  lessonCount: number
  assignmentCount: number
  enrollmentCount: number
  /** Số bài học người dùng hiện tại đã hoàn thành; chỉ có giá trị ở API của học sinh. */
  completedLessonCount: number
  /** Người đang xem đã ghi danh chưa; dùng ở trang khám phá lớp học. */
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
  /** Mức độ nhận thức; rỗng nghĩa là chưa phân loại. */
  level: QuestionLevel
  position: number
  explanation: string
  /** Chỉ dùng cho câu tự luận. */
  sampleAnswer: string
  rubric: string
  /**
   * Ý nghĩa thay đổi theo `type`:
   * - `single_choice` / `multi_choice`: các phương án, `isCorrect` đánh dấu phương án đúng
   * - `true_false`: mỗi dòng là một phát biểu, `isCorrect` là giá trị Đúng/Sai của phát biểu
   * - `fill_blank`: mỗi dòng là một đáp án được chấp nhận
   * - `essay`: luôn rỗng
   *
   * Khi học sinh đang làm bài, backend gỡ bỏ đáp án: `isCorrect` luôn false và
   * riêng câu điền khuyết thì danh sách này rỗng.
   */
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
  /** Với câu đúng/sai, đây là các phát biểu học sinh đánh dấu "Đúng". */
  selectedOptionIds: string[]
  /** Dùng cho cả câu tự luận lẫn câu điền khuyết. */
  essayText: string
  isCorrect: boolean | null
  score: number
  comment: string
  /** Điểm và nhận xét AI gợi ý; chỉ tham khảo, điểm chính thức là `score`. */
  aiScore: number | null
  aiComment: string
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

/** Một cột trên bảng điểm: một bài tập của chương trình. */
export interface GradebookAssignment {
  nodeId: string
  slug: string
  title: string
  questionCount: number
  /** Tổng điểm các câu hỏi hiện có của bài tập. */
  maxScore: number
  /** 0 = không đặt ngưỡng đạt. */
  passScore: number
  submittedCount: number
  pendingCount: number
  passedCount: number
  averageScore: number
}

/** Điểm của một học sinh ở một bài tập — lấy lượt có điểm cao nhất. */
export interface GradebookCell {
  assignmentId: string
  submissionId: string
  score: number
  /** Thang điểm lúc nộp, có thể khác thang hiện tại nếu câu hỏi đã đổi. */
  maxScore: number
  attempts: number
  /** Số lượt còn chờ chấm tay. */
  pending: number
  submittedAt: string
}

export interface GradebookRow {
  userId: string
  fullName: string
  email: string
  /** Đã bị gỡ ghi danh nhưng vẫn còn bài nộp cũ thì false. */
  enrolled: boolean
  /** Khoá là ID bài tập; bài chưa nộp thì không có khoá tương ứng. */
  cells: Record<string, GradebookCell>
  doneCount: number
  pendingCount: number
  totalScore: number
}

/** Bảng điểm bài tập của cả chương trình. */
export interface ProgramGradebook {
  assignments: GradebookAssignment[]
  rows: GradebookRow[]
  totalMaxScore: number
  averageScore: number
  pendingCount: number
  submittedCells: number
  totalCells: number
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
  trainer: 'Giáo viên',
  supervisor: 'Giám sát',
  student: 'Học sinh',
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
  single_choice: 'Trắc nghiệm nhiều lựa chọn',
  true_false: 'Đúng / Sai',
  fill_blank: 'Điền khuyết',
  essay: 'Tự luận',
  multi_choice: 'Trắc nghiệm — nhiều đáp án',
}

/** Bốn dạng chuẩn hiện ở trình soạn thảo, theo đúng thứ tự trong Bảng 3.1. */
export const QUESTION_TYPES: QuestionType[] = ['single_choice', 'true_false', 'fill_blank', 'essay']


// ---------------------------------------------------------------------------
// Các chức năng AI
// ---------------------------------------------------------------------------

/** Dạng câu hỏi theo cách gọi của lớp Service AI. */
export type AIQuestionKind = 'multiple_choice' | 'true_false' | 'fill_blank' | 'essay'

export const AI_QUESTION_KINDS: { kind: AIQuestionKind; label: string }[] = [
  { kind: 'multiple_choice', label: 'Trắc nghiệm nhiều lựa chọn' },
  { kind: 'true_false', label: 'Đúng / Sai' },
  { kind: 'fill_blank', label: 'Điền khuyết' },
  { kind: 'essay', label: 'Tự luận' },
]

// ---------------------------------------------------------------------------
// Trang tổng quan ("không gian làm việc AI")
// ---------------------------------------------------------------------------

/** Một bài tập được giao, dùng cho mục "lịch học" ở trang tổng quan. */
export interface ScheduleItem {
  nodeId: string
  nodeSlug: string
  title: string
  programId: string
  programSlug: string
  programName: string
  dueAt: string | null
  submitted: boolean
  overdue: boolean
}

/** Tiến độ học một lớp: số bài đã hoàn thành trên tổng số bài. */
export interface CourseProgress {
  programId: string
  programSlug: string
  title: string
  completed: number
  total: number
  /** Điểm trung bình các bài đã nộp; null khi chưa nộp bài nào. */
  averageScore: number | null
}

/** Một học sinh cần giáo viên chú ý. Luật cố định, chưa dùng AI phân tích. */
export interface LearningAlert {
  userId: string
  fullName: string
  email: string
  programId: string
  programSlug: string
  programName: string
  kind: 'low_score' | 'overdue'
  /** Điểm trung bình (%) với low_score, số bài quá hạn với overdue. */
  value: number
}

export interface Workspace {
  schedule: ScheduleItem[]
  progress: CourseProgress[]
  pendingGrading: number
  alerts: LearningAlert[]
}

/** Loại tài liệu trong kho dùng chung. */
export type MaterialKind = 'pdf' | 'slide' | 'document' | 'video' | 'link'

export const MATERIAL_KIND_LABEL: Record<MaterialKind, string> = {
  pdf: 'PDF',
  slide: 'Slide trình chiếu',
  document: 'Văn bản',
  video: 'Video',
  link: 'Liên kết',
}

/** Một tài liệu trong kho dùng chung toàn hệ thống, không thuộc lớp nào. */
export interface Material {
  id: string
  title: string
  description: string
  category: string
  url: string
  /** Rỗng khi là link ngoài, không phải Google Drive. */
  driveFileId: string
  kind: MaterialKind
  isPublished: boolean
  createdBy: string | null
  createdAt: string
  updatedAt: string
  createdByName: string
}

export interface MaterialList {
  items: Material[]
  /** Các phân loại đang được dùng, để dựng ô lọc. */
  categories: string[]
  canManage: boolean
}

export interface PasswordResetRequest {
  id: string
  email: string
  /** null khi email không khớp tài khoản nào trong hệ thống. */
  userId: string | null
  note: string
  status: 'pending' | 'done' | 'rejected'
  handledBy: string | null
  handledAt: string | null
  createdAt: string
  fullName: string
  handledByName: string
}

export interface SignupSettings {
  enabled: boolean
  allowedDomains: string
}

export interface AIStatus {
  enabled: boolean
  provider: string
  model: string
  levels: string[]
}

/** Một nhà cung cấp AI được hỗ trợ, dùng cho ô chọn ở màn hình cấu hình. */
export interface AIProvider {
  id: string
  label: string
  /** Giao thức REST: gemini | openai | anthropic. */
  kind: string
  endpoint: string
  models: string[]
  defaultModel: string
  apiKeyUrl: string
  /** Nhà cung cấp không có endpoint cố định, admin phải tự nhập. */
  requiresBaseUrl: boolean
}

/** Một kênh gọi AI đã lưu. Khoá thật không bao giờ được gửi về trình duyệt. */
export interface AIChannel {
  label: string
  provider: string
  model: string
  baseUrl: string
  disabled: boolean
  hasApiKey: boolean
  /** 4 ký tự cuối của khoá. */
  apiKeyHint: string
}

/** Một kênh gửi lên máy chủ; apiKey rỗng nghĩa là giữ nguyên khoá đã lưu. */
export interface AIChannelInput {
  label: string
  provider: string
  apiKey: string
  model: string
  baseUrl: string
  disabled: boolean
}

export interface AISettings {
  enabled: boolean
  source: 'database' | 'env' | 'none'
  /** Các kênh theo đúng thứ tự xoay vòng khi gọi AI. */
  channels: AIChannel[]
  /** Thông tin kênh đầu tiên, để các chỗ chỉ cần biết "đang dùng model gì". */
  provider: string
  model: string
  hasApiKey: boolean
  apiKeyHint: string
  providers: AIProvider[]
}

export interface QuestionSpec {
  kind: AIQuestionKind
  count: number
}

/** Một câu hỏi do AI sinh, đã qua bước kiểm tra tự động ở máy chủ. */
export interface GeneratedQuestion {
  type: AIQuestionKind
  level: QuestionLevel
  question: string
  explanation: string
  points: number
  options?: string[]
  answer?: number
  statements?: { text: string; answer: boolean }[]
  answers?: string[]
  sampleAnswer?: string
  rubric?: string
}

export interface GenerateQuestionsResult {
  questions: GeneratedQuestion[]
  /** Số lần phải gọi Gemini cho tới khi kết quả đạt yêu cầu. */
  attempts: number
}

export interface LessonActivity {
  name: string
  durationMinutes: number
  objective: string
  teacherActions: string
  studentActions: string
  output: string
}

export interface LessonPlanContent {
  title: string
  objectives: string[]
  materials: string[]
  activities: LessonActivity[]
  homework: string
  notes: string
}

export interface LessonPlan {
  id: string
  ownerId: string
  programId: string | null
  subject: string
  grade: string
  topic: string
  objectives: string
  durationMinutes: number
  title: string
  content: LessonPlanContent
  createdAt: string
  updatedAt: string
  ownerName?: string
}

export type AITaskKind = 'lesson' | 'lesson_plan' | 'questions' | 'grade_essay' | 'chat'

export interface AITask {
  id: string
  userId: string
  kind: AITaskKind
  title: string
  status: 'ok' | 'error'
  detail: string
  attempts: number
  createdAt: string
}

export interface AIMessage {
  id: string
  role: 'user' | 'model'
  content: string
  createdAt: string
}

export interface AICourseSource {
  nodeId: string
  title: string
  url: string
  contentAvailable: boolean
}

export interface AIConversation {
  programId: string | null
  nodeId: string | null
  id: string
  userId: string
  title: string
  createdAt: string
  updatedAt: string
  messages?: AIMessage[]
}

/** Gợi ý chấm điểm của AI cho một câu tự luận. */
export interface AIGradeSuggestion {
  answerId: string
  questionId: string
  score: number
  maxPoints: number
  comment: string
  strengths: string[]
  missing: string[]
}
