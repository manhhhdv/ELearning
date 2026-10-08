import { useEffect, useRef, useState } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'

import { useAuth } from '../auth'
import { ROLE_LABEL } from '../api/types'
import { IconBell, IconHelp, IconPanelClose, IconPanelOpen, IconSearch } from './icons'

interface Props {
  pendingCount?: number
  onToggleHelp: () => void
  onToggleSidebar: () => void
  sidebarCollapsed: boolean
}

export function TopBar({ pendingCount = 0, onToggleHelp, onToggleSidebar, sidebarCollapsed }: Props) {
  const { user, signOut } = useAuth()
  const navigate = useNavigate()
  const [term, setTerm] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!menuOpen) return
    const close = (e: MouseEvent) => {
      const target = e.target as Node
      if (menuRef.current && !menuRef.current.contains(target)) setMenuOpen(false)
    }
    document.addEventListener('mousedown', close)
    return () => document.removeEventListener('mousedown', close)
  }, [menuOpen])

  if (!user) return null

  const initials = (user.fullName || user.email)
    .split(' ').filter(Boolean).slice(-2).map((p) => p[0]?.toUpperCase()).join('')

  const search = (e: React.FormEvent) => {
    e.preventDefault()
    const q = term.trim()
    navigate(q ? `/hoc?q=${encodeURIComponent(q)}` : '/hoc')
  }

  return (
    <header className="topbar">
      <button
        className="icon-btn topbar-collapse"
        title={sidebarCollapsed ? 'Mở rộng menu' : 'Thu gọn menu'}
        onClick={onToggleSidebar}
      >
        {sidebarCollapsed ? <IconPanelOpen size={19} /> : <IconPanelClose size={19} />}
      </button>

      <form className="topbar-search" onSubmit={search} role="search">
        <IconSearch size={16} />
        <input
          type="text"
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          placeholder="Tìm kiếm..."
          aria-label="Tìm kiếm"
        />
      </form>

      <div className="topbar-actions">
        <NavLink to="/ket-qua" className="icon-btn" title="Bài đang chờ chấm">
          <IconBell />
          {pendingCount > 0 && <span className="dot">{pendingCount}</span>}
        </NavLink>

        <button className="icon-btn" title="Trợ giúp" onClick={onToggleHelp}>
          <IconHelp />
        </button>

        <div className="user-menu" ref={menuRef}>
          <button className="topbar-user" onClick={() => setMenuOpen((v) => !v)} aria-haspopup="menu" aria-expanded={menuOpen}>
            <span className="avatar">
              {user.avatarUrl ? <img src={user.avatarUrl} alt="" /> : initials || '?'}
            </span>
            <span className="topbar-user-info">
              <b>{user.fullName || 'Chưa đặt tên'}</b>
              <span>{ROLE_LABEL[user.role]}</span>
            </span>
          </button>

          {menuOpen && (
            <div className="user-pop" role="menu">
              <div className="who">
                <b>{user.fullName || 'Chưa đặt tên'}</b>
                <span>{user.email}</span>
                <span className="badge badge-primary" style={{ marginTop: 6 }}>{ROLE_LABEL[user.role]}</span>
              </div>
              <NavLink to="/tac-vu-ai" onClick={() => setMenuOpen(false)}>Tác vụ AI gần đây</NavLink>
              <NavLink to="/doi-mat-khau" onClick={() => setMenuOpen(false)}>Đổi mật khẩu</NavLink>
              <button onClick={() => { signOut(); navigate('/dang-nhap') }}>Đăng xuất</button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
