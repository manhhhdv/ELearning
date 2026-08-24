import type {
  AssignmentResults, AttemptView, DashboardStats, Enrollment, GoogleSettings,
  Program, Question, QuestionType, Submission, SubmissionDetail, TreeNode, User,
} from './types'

const TOKEN_KEY = 'elearning.token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string | null) {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
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

export interface QuestionInput {
  /** Bỏ trống khi tạo mới thì máy chủ tự sinh mã C01, C02… */
  code?: string
  type: QuestionType
  prompt: string
  points: number
  explanation: string
  options: { content: string; isCorrect: boolean }[]
}

export const api = {
  authConfig: () => get<{ googleEnabled: boolean }>('/auth/config'),
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
    lesson?: LessonInput
    assignment?: AssignmentInput
  }) => post<TreeNode>(`/programs/${programId}/nodes`, body),
  getNode: (nodeId: string) => get<TreeNode>(`/nodes/${nodeId}`),
  importStructure: (programId: string, parentId: string | null, items: StructureItemInput[]) =>
    post<{ imported: number }>(`/programs/${programId}/structure/import`, { parentId, items }),
  importStructureFile: (programId: string, file: File) =>
    postFile<{ text: string }>(`/programs/${programId}/structure/import-file`, file),
  updateNode: (nodeId: string, body: {
    title?: string
    description?: string
    isPublished?: boolean
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

  // Học viên
  myPrograms: () => get<Program[]>('/my/programs'),
  catalog: (search = '') => get<Program[]>(`/catalog?search=${encodeURIComponent(search)}`),
  selfEnroll: (programId: string) => post<Program>(`/programs/${programId}/self-enroll`),
  selfUnenroll: (programId: string) => del(`/programs/${programId}/self-enroll`),
  mySubmissions: () => get<Submission[]>('/my/submissions'),
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
  listProgramSubmissions: (programId: string, onlyPending = false) =>
    get<Submission[]>(`/programs/${programId}/submissions${onlyPending ? '?pending=true' : ''}`),
  getSubmission: (id: string) => get<SubmissionDetail>(`/submissions/${id}`),
  gradeSubmission: (id: string, feedback: string, answers: { answerId: string; score: number; comment: string }[]) =>
    post<Submission>(`/submissions/${id}/grade`, { feedback, answers }),

  // Quản trị hệ thống
  dashboard: () => get<DashboardStats>('/admin/dashboard'),
  getGoogleSettings: () => get<GoogleSettings>('/admin/settings/google'),
  saveGoogleSettings: (body: {
    enabled: boolean
    clientId?: string
    clientSecret?: string
    allowedDomains?: string
    autoProvisionRole?: string
  }) => put<GoogleSettings>('/admin/settings/google', body),
}
