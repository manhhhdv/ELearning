import { Navigate, Route, Routes } from 'react-router-dom'

import { AI_CHAT_PATH } from './components/aiChatPath'
import { Layout } from './components/Layout'
import { Loading } from './components/ui'
import { canAccessAdminArea, useAuth } from './auth'
import { AdminDashboardPage } from './pages/AdminDashboardPage'
import { AIChatPage } from './pages/AIChatPage'
import { AISettingsPage } from './pages/AISettingsPage'
import { AITasksPage } from './pages/AITasksPage'
import { ChangePasswordPage } from './pages/ChangePasswordPage'
import { CatalogPage } from './pages/CatalogPage'
import { CoursePage } from './pages/CoursePage'
import { GoogleSettingsPage } from './pages/GoogleSettingsPage'
import { GradingPage } from './pages/GradingPage'
import { LessonPlanPage } from './pages/LessonPlanPage'
import { LoginPage } from './pages/LoginPage'
import { MaterialsPage } from './pages/MaterialsPage'
import { MyProgramsPage } from './pages/MyProgramsPage'
import { MyResultsPage } from './pages/MyResultsPage'
import { PreviewPage } from './pages/PreviewPage'
import { ProgramBuilderPage } from './pages/ProgramBuilderPage'
import { ProgramsPage } from './pages/ProgramsPage'
import { QuestionGenPage } from './pages/QuestionGenPage'
import { SignupSettingsPage } from './pages/SignupSettingsPage'
import { SubmissionPage } from './pages/SubmissionPage'
import { TrainerAnalyticsPage } from './pages/TrainerAnalyticsPage'
import { UsersPage } from './pages/UsersPage'
import { WorkspacePage } from './pages/WorkspacePage'

export default function App() {
  const { user, loading } = useAuth()

  if (loading) return <Loading label="Đang khôi phục phiên đăng nhập…" />

  if (!user) {
    return (
      <Routes>
        <Route path="/dang-nhap" element={<LoginPage />} />
        <Route path="*" element={<Navigate to="/dang-nhap" replace />} />
      </Routes>
    )
  }

  // Xem được khu quản lý: admin, giáo viên (sửa được) và Giám sát (chỉ xem).
  const manage = canAccessAdminArea(user)
  // Trang chính là "không gian làm việc AI" theo vai trò; khu quản lý và khu
  // học tập đều mở từ đây hoặc từ thanh điều hướng.
  const home = '/tong-quan'

  return (
    <Routes>
      <Route path="/dang-nhap" element={<Navigate to={home} replace />} />

      {/* Trình học một lớp chiếm trọn màn hình: thanh nội dung lớp học thay cho sidebar chung.
          Mã bài + mã lớp nằm trên URL để tải lại trang vẫn giữ nguyên vị trí đang học. */}
      <Route path="/hoc/:programSlug" element={<CoursePage />} />
      <Route path="/hoc/:programSlug/:nodeSlug" element={<CoursePage />} />

      {/* Xem trước như học sinh ngay trong lúc soạn — không tính điểm/tiến độ thật. */}
      {manage && (
        <>
          <Route path="/xem-truoc/:programSlug" element={<PreviewPage />} />
          <Route path="/xem-truoc/:programSlug/:nodeSlug" element={<PreviewPage />} />
        </>
      )}

      <Route element={<Layout />}>
        <Route path="/" element={<Navigate to={home} replace />} />
        <Route path="/tong-quan" element={<WorkspacePage />} />
        <Route path="/doi-mat-khau" element={<ChangePasswordPage />} />

        {/* Khu vực học tập, mọi vai trò đều dùng được */}
        <Route path="/hoc" element={<MyProgramsPage />} />
        <Route path="/kham-pha" element={<CatalogPage />} />
        <Route path="/ket-qua" element={<MyResultsPage />} />
        <Route path="/tai-lieu" element={<MaterialsPage />} />
        <Route path="/bai-nop/:submissionId" element={<SubmissionPage />} />
        {/* Trợ lý AI và nhật ký tác vụ AI dùng chung cho mọi vai trò. */}
        <Route path={AI_CHAT_PATH} element={<AIChatPage />} />
        <Route path="/tac-vu-ai" element={<AITasksPage />} />

        {/* Khu vực quản lý */}
        {manage && (
          <>
            <Route path="/quan-tri" element={<AdminDashboardPage />} />
            <Route path="/quan-tri/phan-tich" element={<TrainerAnalyticsPage />} />
            <Route path="/quan-tri/chuong-trinh" element={<ProgramsPage />} />
            <Route path="/quan-tri/chuong-trinh/:programSlug" element={<ProgramBuilderPage />} />
            <Route path="/quan-tri/cham-bai" element={<GradingPage />} />
            <Route path="/quan-tri/giao-an" element={<LessonPlanPage />} />
            <Route path="/quan-tri/tao-de" element={<QuestionGenPage />} />
          </>
        )}
        {user.role === 'admin' && (
          <>
            <Route path="/quan-tri/nguoi-dung" element={<UsersPage />} />
            <Route path="/quan-tri/cai-dat/google" element={<GoogleSettingsPage />} />
            <Route path="/quan-tri/cai-dat/ai" element={<AISettingsPage />} />
            <Route path="/quan-tri/cai-dat/dang-ky" element={<SignupSettingsPage />} />
          </>
        )}

        <Route path="*" element={<Navigate to={home} replace />} />
      </Route>
    </Routes>
  )
}
