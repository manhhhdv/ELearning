import { Navigate, Route, Routes } from 'react-router-dom'

import { Layout } from './components/Layout'
import { Loading } from './components/ui'
import { canAccessAdminArea, useAuth } from './auth'
import { AdminDashboardPage } from './pages/AdminDashboardPage'
import { ChangePasswordPage } from './pages/ChangePasswordPage'
import { CatalogPage } from './pages/CatalogPage'
import { CoursePage } from './pages/CoursePage'
import { GoogleSettingsPage } from './pages/GoogleSettingsPage'
import { GradingPage } from './pages/GradingPage'
import { LoginPage } from './pages/LoginPage'
import { MyProgramsPage } from './pages/MyProgramsPage'
import { MyResultsPage } from './pages/MyResultsPage'
import { PreviewPage } from './pages/PreviewPage'
import { ProgramBuilderPage } from './pages/ProgramBuilderPage'
import { ProgramsPage } from './pages/ProgramsPage'
import { SubmissionPage } from './pages/SubmissionPage'
import { UsersPage } from './pages/UsersPage'

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

  // Xem được khu quản lý: admin, giảng viên (sửa được) và Giám sát (chỉ xem).
  const manage = canAccessAdminArea(user)
  // Trang chính là danh sách khoá học kèm thống kê; khu quản lý nằm trên thanh điều hướng.
  const home = '/hoc'

  return (
    <Routes>
      <Route path="/dang-nhap" element={<Navigate to={home} replace />} />

      {/* Trình học một khoá chiếm trọn màn hình: thanh nội dung khoá học thay cho sidebar chung.
          Mã bài + mã khoá nằm trên URL để tải lại trang vẫn giữ nguyên vị trí đang học. */}
      <Route path="/hoc/:programSlug" element={<CoursePage />} />
      <Route path="/hoc/:programSlug/:nodeSlug" element={<CoursePage />} />

      {/* Xem trước như học viên ngay trong lúc soạn — không tính điểm/tiến độ thật. */}
      {manage && (
        <>
          <Route path="/xem-truoc/:programSlug" element={<PreviewPage />} />
          <Route path="/xem-truoc/:programSlug/:nodeSlug" element={<PreviewPage />} />
        </>
      )}

      <Route element={<Layout />}>
        <Route path="/" element={<Navigate to={home} replace />} />
        {/* Giữ đường dẫn cũ để link đã chia sẻ không bị hỏng. */}
        <Route path="/tong-quan" element={<Navigate to="/hoc" replace />} />
        <Route path="/doi-mat-khau" element={<ChangePasswordPage />} />

        {/* Khu vực học tập, mọi vai trò đều dùng được */}
        <Route path="/hoc" element={<MyProgramsPage />} />
        <Route path="/kham-pha" element={<CatalogPage />} />
        <Route path="/ket-qua" element={<MyResultsPage />} />
        <Route path="/bai-nop/:submissionId" element={<SubmissionPage />} />

        {/* Khu vực quản lý */}
        {manage && (
          <>
            <Route path="/quan-tri" element={<AdminDashboardPage />} />
            <Route path="/quan-tri/chuong-trinh" element={<ProgramsPage />} />
            <Route path="/quan-tri/chuong-trinh/:programSlug" element={<ProgramBuilderPage />} />
            <Route path="/quan-tri/cham-bai" element={<GradingPage />} />
          </>
        )}
        {user.role === 'admin' && (
          <>
            <Route path="/quan-tri/nguoi-dung" element={<UsersPage />} />
            <Route path="/quan-tri/cai-dat/google" element={<GoogleSettingsPage />} />
          </>
        )}

        <Route path="*" element={<Navigate to={home} replace />} />
      </Route>
    </Routes>
  )
}
