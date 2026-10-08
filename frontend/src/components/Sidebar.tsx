import { NavLink } from 'react-router-dom'

import { useAIStatus } from '../api/useAIStatus'
import { canAccessAdminArea, useAuth } from '../auth'
import { AI_CHAT_PATH } from './aiChatPath'
import {
  IconBook, IconChevron, IconDashboard, IconGrade, IconSearch,
  IconSettings, IconSparkles, IconUsers,
} from './icons'
import { Logo } from './Logo'

interface NavItem {
  to: string
  label: string
  icon: (size: number) => React.ReactNode
  end?: boolean
}

export function Sidebar({ collapsed }: { collapsed: boolean }) {
  const { user } = useAuth()
  const ai = useAIStatus()

  if (!user) return null
  const manage = canAccessAdminArea(user)

  const items: NavItem[] = [
    { to: '/tong-quan', label: 'Tổng quan', icon: (s) => <IconDashboard size={s} />, end: true },
    ...(manage ? [] : [
      { to: '/hoc', label: 'Lớp học của tôi', icon: (s: number) => <IconBook size={s} />, end: true },
      { to: '/kham-pha', label: 'Khám phá', icon: (s: number) => <IconSearch size={s} /> },
    ]),
    { to: '/tai-lieu', label: 'Tài liệu', icon: (s) => <IconGrade size={s} /> },
    ...(manage ? [] : [
      { to: '/ket-qua', label: 'Kết quả', icon: (s: number) => <IconChevron size={s} /> },
    ]),
  ]

  // Các mục quản lý hiện phẳng ngay trong menu chính, không còn gom vào
  // nhóm "Quản lý" phải bấm mở ra mới thấy.
  const manageItems: NavItem[] = [
    ...(user.role === 'admin' ? [
      { to: '/quan-tri', label: 'Bảng điều khiển', icon: (s: number) => <IconDashboard size={s} />, end: true },
    ] : [
      { to: '/quan-tri/phan-tich', label: 'Phân tích học tập', icon: (s: number) => <IconDashboard size={s} />, end: true },
    ]),
    { to: '/quan-tri/chuong-trinh', label: 'Lớp học', icon: (s) => <IconBook size={s} /> },
    { to: '/quan-tri/cham-bai', label: 'Chấm bài', icon: (s) => <IconGrade size={s} /> },
  ]

  // Ba tính năng AI gom chung một nhóm "Trợ lý AI" để menu chính bớt rối;
  // nhãn bỏ hậu tố "bằng AI" vì tiêu đề nhóm đã nói rõ.
  const aiItems: NavItem[] = manage ? [
    { to: '/quan-tri/giao-an', label: 'Soạn giáo án', icon: (s: number) => <IconSparkles size={s} /> },
    { to: '/quan-tri/tao-de', label: 'Tạo đề', icon: (s: number) => <IconSparkles size={s} /> },
  ] : []

  // Nhóm riêng cho các mục chỉ admin đụng tới, gom dưới nhãn "Quản trị".
  const adminItems: NavItem[] = user.role === 'admin' ? [
    { to: '/quan-tri/nguoi-dung', label: 'Người dùng', icon: (s: number) => <IconUsers size={s} /> },
    { to: '/quan-tri/cai-dat/google', label: 'Đăng nhập Google', icon: (s: number) => <IconSettings size={s} /> },
    { to: '/quan-tri/cai-dat/ai', label: 'Cấu hình AI', icon: (s: number) => <IconSettings size={s} /> },
    { to: '/quan-tri/cai-dat/dang-ky', label: 'Đăng ký tài khoản', icon: (s: number) => <IconSettings size={s} /> },
  ] : []

  return (
    <aside className={`sidebar ${collapsed ? 'collapsed' : ''}`}>
      <NavLink to="/tong-quan" className="sidebar-brand" aria-label="Đào Tạo — trang chủ">
        <Logo size={42} />
        {!collapsed && <span>EduFlow</span>}
      </NavLink>

      <nav className="sidebar-nav">
        {items.map((it) => (
          <NavLink key={it.to} to={it.to} end={it.end} className="sidebar-link" title={it.label}>
            <span className="sidebar-link-icon">{it.icon(18)}</span>
            {!collapsed && <span className="sidebar-link-label">{it.label}</span>}
          </NavLink>
        ))}

        {manage && (
          <>
            {!collapsed && <div className="sidebar-sep" />}
            {manageItems.map((it) => (
              <NavLink key={it.to} to={it.to} end={it.end} className="sidebar-link" title={it.label}>
                <span className="sidebar-link-icon">{it.icon(18)}</span>
                {!collapsed && <span className="sidebar-link-label">{it.label}</span>}
              </NavLink>
            ))}
          </>
        )}

        {(ai?.enabled || aiItems.length > 0) && (
          <>
            <div className="sidebar-sep" />
            {!collapsed && <div className="sidebar-group-label">Trợ lý AI</div>}
            {ai?.enabled && (
              <NavLink to={AI_CHAT_PATH} className="sidebar-link" title="Hỏi trợ lý AI">
                <span className="sidebar-link-icon"><IconSparkles size={18} /></span>
                {!collapsed && <span className="sidebar-link-label">Hỏi trợ lý AI</span>}
              </NavLink>
            )}
            {aiItems.map((it) => (
              <NavLink key={it.to} to={it.to} end={it.end} className="sidebar-link" title={it.label}>
                <span className="sidebar-link-icon">{it.icon(18)}</span>
                {!collapsed && <span className="sidebar-link-label">{it.label}</span>}
              </NavLink>
            ))}
          </>
        )}

        {adminItems.length > 0 && (
          <>
            <div className="sidebar-sep" />
            {!collapsed && <div className="sidebar-group-label">Quản trị</div>}
            {adminItems.map((it) => (
              <NavLink key={it.to} to={it.to} end={it.end} className="sidebar-link" title={it.label}>
                <span className="sidebar-link-icon">{it.icon(18)}</span>
                {!collapsed && <span className="sidebar-link-label">{it.label}</span>}
              </NavLink>
            ))}
          </>
        )}
      </nav>
    </aside>
  )
}
