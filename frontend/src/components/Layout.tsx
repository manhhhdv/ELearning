import { useEffect, useState } from 'react'
import { Outlet } from 'react-router-dom'

import { api } from '../api/client'
import { canManageContent, useAuth } from '../auth'
import { TopNav } from './TopNav'
import { IconHelp } from './icons'

export function Layout() {
  const { user } = useAuth()
  const [pending, setPending] = useState(0)
  const [helpOpen, setHelpOpen] = useState(false)

  // Chuông báo số bài của chính mình còn chờ giảng viên chấm.
  useEffect(() => {
    if (!user) return
    api.mySubmissions()
      .then((items) => setPending(items.filter((s) => s.status === 'submitted').length))
      .catch(() => setPending(0))
  }, [user])

  if (!user) return null
  const manage = canManageContent(user)

  return (
    <div className="shell">
      <TopNav pendingCount={pending} onToggleHelp={() => setHelpOpen((v) => !v)} />

      <main className="main">
        <Outlet />
      </main>

      {helpOpen && (
        <div className="help-panel">
          <div className="spread" style={{ marginBottom: 10 }}>
            <h3>Hướng dẫn nhanh</h3>
            <button className="btn btn-ghost btn-sm" onClick={() => setHelpOpen(false)} aria-label="Đóng">✕</button>
          </div>
          <ul style={{ margin: 0, paddingLeft: 18 }}>
            <li>Vào <b>Khoá học của tôi</b> để bắt đầu học.</li>
            <li>Bấm <b>Đánh dấu hoàn thành</b> ở cuối mỗi bài để cập nhật tiến độ.</li>
            <li>Bài tập trắc nghiệm có điểm ngay; bài tự luận chờ giảng viên chấm.</li>
            <li>Xem điểm ở mục <b>Kết quả</b>.</li>
            {manage && <li>Soạn nội dung ở mục <b>Quản lý</b>.</li>}
          </ul>
          <p className="hint" style={{ marginTop: 12, marginBottom: 0 }}>
            Cần hỗ trợ thêm, liên hệ quản trị viên của đơn vị bạn.
          </p>
        </div>
      )}

      <button className="help-fab" onClick={() => setHelpOpen((v) => !v)} aria-label="Trợ giúp">
        <IconHelp size={26} />
      </button>
    </div>
  )
}

export function PageHeader({
  title, subtitle, actions,
}: {
  title: string
  subtitle?: string
  actions?: React.ReactNode
}) {
  return (
    <header className="page-header">
      <div>
        <h1>{title}</h1>
        {subtitle && <div className="sub">{subtitle}</div>}
      </div>
      {actions && <div className="wrap-gap">{actions}</div>}
    </header>
  )
}
