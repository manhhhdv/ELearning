import type {
  AIChannelInput, AICourseSource, AIConversation, AIGradeSuggestion, AISettings, AIStatus, AITask, AssignmentResults, AttemptView,
  DashboardStats, Enrollment, GenerateQuestionsResult, GoogleSettings, LessonPlan,
  LessonPlanContent, Program, ProgramGradebook, Question, QuestionLevel, QuestionSpec,
  Material, MaterialKind, MaterialList, PasswordResetRequest, QuestionType, SignupSettings,
  Submission, SubmissionDetail, TreeNode, User, Workspace,
} from './types'

const TOKEN_KEY = 'elearning.token'

/**
 * Token nằm ở localStorage khi người dùng chọn "Ghi nhớ đăng nhập" (còn sau khi
 * đóng trình duyệt), ngược lại ở sessionStorage (mất khi đóng tab).
 *
 * Mọi lần đọc/ghi đều bọc try/catch: chế độ duyệt web riêng tư hoặc cấu hình
 * chặn cookie có thể làm các kho này ném lỗi.
 */
export function getToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_KEY) ?? sessionStorage.getItem(TOKEN_KEY)
  } catch {
    return null
  }
}

export function setToken(token: string | null, remember = true) {
  try {
    // Luôn xoá ở cả hai kho trước, tránh còn sót token cũ ở kho không dùng tới.
    localStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(TOKEN_KEY)
    if (token) {
      (remember ? localStorage : sessionStorage).setItem(TOKEN_KEY, token)
    }
  } catch {
    // Không lưu được thì phiên chỉ sống trong lần tải trang này — vẫn dùng được.
  }
}

/** Lỗi kèm mã HTTP để giao diện phân biệt hết phiên với lỗi nghiệp vụ. */
export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers)
  headers.set('Accept', 'application/json')
  if (init.body && !(init.body instanceof FormData)) headers.set('Content-Type', 'application/json')

  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const res = await fetch(`/api${path}`, { ...init, headers })

  if (res.status === 204) return undefined as T
  const text = await res.text()
  const data = text ? JSON.parse(text) : null

  if (!res.ok) {
    throw new ApiError(res.status, data?.error ?? 'Không kết nối được máy chủ')
  }
  return data as T
}

/**
 * Tải một tệp (VD: Word .docx) từ API về máy người dùng. Không dùng được thẻ
 * <a href> trực tiếp vì API cần header Authorization, nên tải về dạng Blob rồi
 * mới lưu. Tên tệp lấy từ Content-Disposition, thiếu thì dùng fallbackName.
 */
async function download(path: string, fallbackName: string, init: RequestInit = {}): Promise<void> {
  const headers = new Headers(init.headers)
  if (init.body) headers.set('Content-Type', 'application/json')
  const token = getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const res = await fetch(`/api${path}`, { ...init, headers })
  if (!res.ok) {
    let message = 'Không tải được tệp'
    try {
      message = JSON.parse(await res.text())?.error ?? message
    } catch {
      // Phản hồi lỗi không phải JSON — giữ thông báo mặc định.
    }
    throw new ApiError(res.status, message)
  }

  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileNameFrom(res.headers.get('Content-Disposition')) ?? fallbackName
  document.body.appendChild(a)
  a.click()
  a.remove()
  // Thu hồi sau một nhịp để trình duyệt kịp bắt đầu tải.
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}

function fileNameFrom(disposition: string | null): string | null {
  const match = disposition?.match(/filename="?([^";]+)"?/i)
  return match ? match[1] : null
}

const get = <T>(path: string) => request<T>(path)
// Không set Content-Type thủ công: trình duyệt tự thêm boundary cho multipart/form-data.
const postFile = <T>(path: string, file: File) => {
  const body = new FormData()
  body.append('file', file)
  return request<T>(path, { method: 'POST', body })
}
const post = <T>(path: string, body?: unknown) =>
  request<T>(path, { method: 'POST', body: body === undefined ? undefined : JSON.stringify(body) })
const patch = <T>(path: string, body: unknown) =>
  request<T>(path, { method: 'PATCH', body: JSON.stringify(body) })
const put = <T>(path: string, body: unknown) =>
  request<T>(path, { method: 'PUT', body: JSON.stringify(body) })
const del = (path: string) => request<void>(path, { method: 'DELETE' })

// --- Kiểu dữ liệu gửi lên ---

export interface LessonInput {
  contentType: string
  source: string
  durationMinutes: number
  body: string
  attachments: LessonAttachmentInput[]
}

export interface LessonAttachmentInput {
  name: string
  url: string
}

export interface StructureItemInput {
  level: number
  kind: string
  title: string
  description: string
  contentType: string
  source: string
  durationMinutes: number
}

export interface AssignmentInput {
  instructions: string
  timeLimitMinutes: number
  maxAttempts: number
  passScore: number
  shuffleQuestions: boolean
  dueAt: string | null
}

/** Kết quả cấp tài khoản / đặt lại mật khẩu; mật khẩu thô chỉ trả về đúng một lần. */
export interface IssuedCredentials {
  user: User
  password: string
}

export interface UserImportItem {
  email: string
  fullName: string
  role: string
  /** Bỏ trống để máy chủ sinh mật khẩu ngẫu nhiên cho dòng này. */
  password: string
}

/** Kết quả của một dòng trong lô nhập tài khoản. */
export interface UserImportResult {
  row: number
  email: string
  fullName: string
  role: string
  status: 'created' | 'skipped' | 'failed'
  /** Chỉ có ở dòng tạo mới thành công, và chỉ trả về đúng một lần. */
  password: string
  message: string
}

export interface UserImportReport {
  created: number
  skipped: number
  failed: number
  results: UserImportResult[]
}

/** Dữ liệu một tài liệu dùng chung gửi lên máy chủ. */
export interface MaterialInput {
  title: string
  description: string
  category: string
  url: string
  kind: MaterialKind
  isPublished: boolean
}

export interface QuestionInput {
  /** Bỏ trống khi tạo mới thì máy chủ tự sinh mã C01, C02… */
  code?: string
  type: QuestionType
  prompt: string
  points: number
  level?: QuestionLevel
  explanation: string
  /** Chỉ dùng cho câu tự luận. */
  sampleAnswer?: string
  rubric?: string
  /** Phương án / phát biểu / đáp án được chấp nhận, tuỳ theo `type`. */
  options: { content: string; isCorrect: boolean }[]
}

/** Thông tin bài học giáo viên nhập ở màn hình tạo giáo án bằng AI. */
export interface LessonPlanInput {
  subject: string
  grade: string
  topic: string
  content: string
  objectives: string
  durationMinutes: number
}

/** Yêu cầu tạo câu hỏi từ chủ đề / nội dung bài học. */
export interface GenerateQuestionsInput {
  programId?: string
  sourceNodeId?: string
  subject: string
  grade: string
  topic: string
  content: string
  levels: string[]
  specs: QuestionSpec[]
}

export const api = {
  authConfig: () => get<{
    googleEnabled: boolean
    signupEnabled: boolean
    signupAllowedDomains: string
  }>('/auth/config'),
  /** Tự đăng ký tài khoản; trả về token dùng được ngay như khi đăng nhập. */
  register: (body: { email: string; fullName: string; password: string }) =>
    post<{ token: string; user: User }>('/auth/register', body),
  /** Gửi yêu cầu đặt lại mật khẩu để quản trị viên xử lý. */
  forgotPassword: (email: string, note: string) =>
    post<{ message: string }>('/auth/forgot-password', { email, note }),
  login: (email: string, password: string) =>
    post<{ token: string; user: User }>('/auth/login', { email, password }),
  me: () => get<User>('/auth/me'),
  changePassword: (currentPassword: string, newPassword: string) =>
    post<{ message: string }>('/auth/change-password', { currentPassword, newPassword }),

  // Người dùng
  listUsers: (params: { search?: string; role?: string } = {}) => {
    const q = new URLSearchParams()
    if (params.search) q.set('search', params.search)
    if (params.role) q.set('role', params.role)
    return get<{ items: User[]; total: number }>(`/users?${q}`)
  },
  createUser: (body: { email: string; fullName: string; role: string; password?: string }) =>
    post<IssuedCredentials>('/users', body),
  importUsers: (items: UserImportItem[]) => post<UserImportReport>('/users/import', { items }),
  importUsersFile: (file: File) => postFile<{ text: string }>('/users/import-file', file),
  updateUser: (id: string, body: { fullName?: string; role?: string; isActive?: boolean }) =>
    patch<User>(`/users/${id}`, body),
  resetPassword: (id: string, password?: string) =>
    post<IssuedCredentials>(`/users/${id}/password`, { password: password ?? '' }),
  deleteUser: (id: string) => del(`/users/${id}`),

  // Chương trình
  listPrograms: (params: { search?: string; status?: string } = {}) => {
    const q = new URLSearchParams()
    if (params.search) q.set('search', params.search)
    if (params.status) q.set('status', params.status)
    return get<Program[]>(`/programs?${q}`)
  },
  createProgram: (body: {
    code: string; title: string; description: string; status: string
    allowSelfEnroll: boolean; isDefaultCourse?: boolean
  }) =>
    post<Program>('/programs', body),
  getProgram: (id: string) => get<Program>(`/programs/${id}`),
  getProgramBySlug: (slug: string) => get<Program>(`/programs/slug/${slug}`),
  updateProgram: (id: string, body: Partial<{
    code: string; title: string; description: string; coverUrl: string
    status: string; allowSelfEnroll: boolean; isDefaultCourse: boolean
  }>) =>
    patch<Program>(`/programs/${id}`, body),
  deleteProgram: (id: string) => del(`/programs/${id}`),
  getTree: (id: string) => get<TreeNode[]>(`/programs/${id}/tree`),

  // Cây nội dung
  createNode: (programId: string, body: {
    parentId: string | null
    kind: string
    title: string
    description?: string
    isPublished?: boolean
    isLocked?: boolean
    lesson?: LessonInput
    assignment?: AssignmentInput
  }) => post<TreeNode>(`/programs/${programId}/nodes`, body),
  getNode: (nodeId: string) => get<TreeNode>(`/nodes/${nodeId}`),
  /** Xuất toàn bộ lớp học thành tệp .zip: mỗi bài học một tệp Word, mỗi bài tập một tệp đề và một tệp đề kèm đáp án. */
  exportProgram: (programId: string) => download(`/programs/${programId}/export`, 'lop-hoc.zip'),
  /** Xuất bài học (bài giảng) hoặc bộ câu hỏi của bài tập ra tệp Word; withAnswers in thêm đáp án. */
  exportNode: (nodeId: string, withAnswers = false) =>
    download(`/nodes/${nodeId}/export${withAnswers ? '?answers=1' : ''}`, 'noi-dung.docx'),
  importStructure: (programId: string, parentId: string | null, items: StructureItemInput[]) =>
    post<{ imported: number }>(`/programs/${programId}/structure/import`, { parentId, items }),
  importStructureFile: (programId: string, file: File) =>
    postFile<{ text: string }>(`/programs/${programId}/structure/import-file`, file),
  updateNode: (nodeId: string, body: {
    title?: string
    description?: string
    isPublished?: boolean
    isLocked?: boolean
    lesson?: LessonInput
    assignment?: AssignmentInput
  }) => patch<TreeNode>(`/nodes/${nodeId}`, body),
  deleteNode: (nodeId: string) => del(`/nodes/${nodeId}`),
  moveNode: (nodeId: string, parentId: string | null, position: number) =>
    post<{ message: string }>(`/nodes/${nodeId}/move`, { parentId, position }),

  // Câu hỏi
  createQuestion: (nodeId: string, body: QuestionInput) =>
    post<Question>(`/nodes/${nodeId}/questions`, body),
  updateQuestion: (questionId: string, body: QuestionInput) =>
    patch<Question>(`/questions/${questionId}`, body),
  deleteQuestion: (questionId: string) => del(`/questions/${questionId}`),
  reorderQuestions: (nodeId: string, questionIds: string[]) =>
    post<{ message: string }>(`/nodes/${nodeId}/questions/reorder`, { questionIds }),
  importQuestions: (nodeId: string, questions: QuestionInput[]) =>
    post<{ imported: number; questions: Question[] }>(`/nodes/${nodeId}/questions/import`, { questions }),
  importQuestionsFile: (nodeId: string, file: File) =>
    postFile<{ text: string }>(`/nodes/${nodeId}/questions/import-file`, file),
  assignmentResults: (nodeId: string) => get<AssignmentResults>(`/nodes/${nodeId}/results`),

  // Ghi danh
  listEnrollments: (programId: string) => get<Enrollment[]>(`/programs/${programId}/enrollments`),
  enroll: (programId: string, userIds: string[], role: 'student' | 'trainer') =>
    post<Enrollment[]>(`/programs/${programId}/enrollments`, { userIds, role }),
  unenroll: (programId: string, userId: string) => del(`/programs/${programId}/enrollments/${userId}`),

  // Học sinh
  myPrograms: () => get<Program[]>('/my/programs'),
  catalog: (search = '') => get<Program[]>(`/catalog?search=${encodeURIComponent(search)}`),
  selfEnroll: (programId: string) => post<Program>(`/programs/${programId}/self-enroll`),
  selfUnenroll: (programId: string) => del(`/programs/${programId}/self-enroll`),
  mySubmissions: () => get<Submission[]>('/my/submissions'),
  /** Dữ liệu trang tổng quan: lịch học, tiến độ, cảnh báo học tập. */
  workspace: () => get<Workspace>('/my/workspace'),
  markComplete: (nodeId: string, completed: boolean) =>
    post<{ completed: boolean }>(`/nodes/${nodeId}/complete`, { completed }),
  getAttempt: (nodeId: string) => get<AttemptView>(`/nodes/${nodeId}/attempt`),
  startAttempt: (nodeId: string) => post<AttemptView>(`/nodes/${nodeId}/attempt/start`),
  submitAssignment: (nodeId: string, answers: {
    questionId: string
    selectedOptionIds: string[]
    essayText: string
  }[]) => post<Submission>(`/nodes/${nodeId}/submit`, { answers }),

  // Chấm bài
  programGradebook: (programId: string) => get<ProgramGradebook>(`/programs/${programId}/gradebook`),
  listProgramSubmissions: (programId: string, onlyPending = false) =>
    get<Submission[]>(`/programs/${programId}/submissions${onlyPending ? '?pending=true' : ''}`),
  getSubmission: (id: string) => get<SubmissionDetail>(`/submissions/${id}`),
  gradeSubmission: (id: string, feedback: string, answers: { answerId: string; score: number; comment: string }[]) =>
    post<Submission>(`/submissions/${id}/grade`, { feedback, answers }),

  // --- Các chức năng AI ---
  aiStatus: () => get<AIStatus>('/ai/status'),

  /** Sinh bản nháp giáo án; máy chủ tự kiểm tra và yêu cầu AI soạn lại nếu chưa đạt. */
  generateLessonPlan: (body: LessonPlanInput) =>
    post<{ plan: LessonPlanContent; attempts: number }>('/ai/lesson-plan', body),
  listLessonPlans: () => get<LessonPlan[]>('/ai/lesson-plans'),
  getLessonPlan: (id: string) => get<LessonPlan>(`/ai/lesson-plans/${id}`),
  saveLessonPlan: (body: {
    programId?: string | null
    subject: string
    grade: string
    topic: string
    objectives: string
    durationMinutes: number
    title: string
    content: LessonPlanContent
  }) => post<LessonPlan>('/ai/lesson-plans', body),
  updateLessonPlan: (id: string, title: string, content: LessonPlanContent) =>
    patch<LessonPlan>(`/ai/lesson-plans/${id}`, { title, content }),
  deleteLessonPlan: (id: string) => del(`/ai/lesson-plans/${id}`),
  /** Xuất giáo án đang soạn trên màn hình (kể cả chưa lưu) ra tệp Word. */
  exportLessonPlanDraft: (body: { subject: string; grade: string; durationMinutes: number; content: LessonPlanContent }) =>
    download('/ai/lesson-plans/export', 'giao-an.docx', { method: 'POST', body: JSON.stringify(body) }),
  /** Xuất một giáo án đã lưu ra tệp Word. */
  exportLessonPlan: (id: string) => download(`/ai/lesson-plans/${id}/export`, 'giao-an.docx'),

  generateLesson: (body: { programId: string; topic: string; grade: string; objectives: string; content: string; durationMinutes: number }) =>
    post<{ title: string; body: string }>('/ai/lesson', body),
  /** Sinh bài giảng bám sát một tài liệu Word (.docx) tải lên thay vì dán tay nội dung. Tệp chỉ đi qua máy chủ, không được lưu lại. */
  generateLessonFromFile: (file: File, fields: { programId: string; topic: string; grade: string; objectives: string; durationMinutes: number }) => {
    const body = new FormData()
    body.append('file', file)
    body.append('programId', fields.programId)
    body.append('topic', fields.topic)
    body.append('grade', fields.grade)
    body.append('objectives', fields.objectives)
    body.append('durationMinutes', String(fields.durationMinutes))
    return request<{ title: string; body: string }>('/ai/lesson/file', { method: 'POST', body })
  },

  generateQuestions: (body: GenerateQuestionsInput) =>
    post<GenerateQuestionsResult>('/ai/questions', body),
  /** Sinh câu hỏi bám sát một tài liệu PDF hoặc Word (.docx) tải lên. Tệp chỉ đi qua máy chủ, không được lưu lại. */
  generateQuestionsFromFile: (file: File, fields: Omit<GenerateQuestionsInput, 'content'>) => {
    const body = new FormData()
    body.append('file', file)
    body.append('subject', fields.subject)
    body.append('grade', fields.grade)
    body.append('topic', fields.topic)
    body.append('levels', fields.levels.join(','))
    body.append('specs', JSON.stringify(fields.specs))
    return request<GenerateQuestionsResult>('/ai/questions/file', { method: 'POST', body })
  },

  /** AI gợi ý điểm cho mọi câu tự luận của một bài nộp; giáo viên vẫn phải xác nhận. */
  aiGradeSubmission: (submissionId: string) =>
    post<{ suggestions: AIGradeSuggestion[] }>(`/submissions/${submissionId}/ai-grade`),

  aiChat: (body: { conversationId?: string; message: string; subject?: string; grade?: string; programId?: string; nodeId?: string }) =>
    post<{ conversationId: string; answer: string; sources: AICourseSource[] }>('/ai/chat', body),
  listConversations: () => get<AIConversation[]>('/ai/conversations'),
  getConversation: (id: string) => get<AIConversation>(`/ai/conversations/${id}`),
  deleteConversation: (id: string) => del(`/ai/conversations/${id}`),
  listAITasks: (limit = 20) => get<AITask[]>(`/ai/tasks?limit=${limit}`),

  // Quản trị hệ thống
  dashboard: () => get<DashboardStats>('/admin/dashboard'),
  trainerAnalytics: () => get<DashboardStats>('/admin/analytics'),
  getGoogleSettings: () => get<GoogleSettings>('/admin/settings/google'),
  saveGoogleSettings: (body: {
    enabled: boolean
    clientId?: string
    clientSecret?: string
    allowedDomains?: string
    autoProvisionRole?: string
  }) => put<GoogleSettings>('/admin/settings/google', body),

  // Kho tài liệu dùng chung. Ai đăng nhập cũng xem được; chỉ giáo viên và
  // quản trị viên mới thêm/sửa/xoá.
  listMaterials: (search = '', category = '') =>
    get<MaterialList>(`/materials/?search=${encodeURIComponent(search)}&category=${encodeURIComponent(category)}`),
  createMaterial: (body: MaterialInput) => post<Material>('/materials/', body),
  updateMaterial: (id: string, body: MaterialInput) => patch<Material>(`/materials/${id}`, body),
  deleteMaterial: (id: string) => del(`/materials/${id}`),

  // Yêu cầu quên mật khẩu (chỉ admin).
  listPasswordResets: (pendingOnly = false) =>
    get<PasswordResetRequest[]>(`/users/password-resets${pendingOnly ? '?pending=true' : ''}`),
  resolvePasswordReset: (id: string, status: 'done' | 'rejected') =>
    post<void>(`/users/password-resets/${id}`, { status }),

  // Cho phép tự đăng ký tài khoản (chỉ admin).
  getSignupSettings: () => get<SignupSettings>('/admin/settings/signup'),
  saveSignupSettings: (body: SignupSettings) =>
    put<SignupSettings>('/admin/settings/signup', body),

  // Cấu hình nhà cung cấp AI (chỉ admin). Khoá API không bao giờ được trả về.
  getAISettings: () => get<AISettings>('/admin/settings/ai'),
  saveAISettings: (body: { enabled: boolean; channels: AIChannelInput[] }) =>
    put<AISettings>('/admin/settings/ai', body),
  /** Gọi thử nhà cung cấp AI. Truyền khoá mới để kiểm tra trước khi lưu. */
  testAISettings: (body: Partial<AIChannelInput> & { index?: number }) =>
    post<{ ok: boolean; message: string }>('/admin/settings/ai/test', body),
}
