import { useEffect, useRef, useState } from 'react'
import { NavLink, useLocation, useNavigate } from 'react-router-dom'

import { canAccessAdminArea, useAuth } from '../auth'
import { ROLE_LABEL } from '../api/types'
import { IconBell, IconChevron, IconHelp, IconSearch } from './icons'
import { Logo } from './Logo'

interface Props {
  /** Số bài nộp đang chờ chấm, hiện thành chấm đỏ trên chuông. */
  pendingCount?: number
  onToggleHelp: () => void
}

export function TopNav({ pendingCount = 0, onToggleHelp }: Props) {
  const { user, signOut } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [term, setTerm] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [manageOpen, setManageOpen] = useState(false)
  // Màn hình hẹp không đủ chỗ cho ô tìm kiếm, thu về một nút bấm mở ra.
  const [searchOpen, setSearchOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const manageRef = useRef<HTMLDivElement>(null)

  // Bấm ra ngoài thì đóng các menu đang mở.
  useEffect(() => {
    if (!menuOpen && !manageOpen) return
    const close = (e: MouseEvent) => {
      const target = e.target as Node
      if (menuRef.current && !menuRef.current.contains(target)) setMenuOpen(false)
      if (manageRef.current && !manageRef.current.contains(target)) setManageOpen(false)
    }
    document.addEventListener('mousedown', close)
    return () => document.removeEventListener('mousedown', close)
  }, [menuOpen, manageOpen])

  if (!user) return null

  const manage = canAccessAdminArea(user)
  const initials = (user.fullName || user.email)
    .split(' ').filter(Boolean).slice(-2).map((p) => p[0]?.toUpperCase()).join('')

  const search = (e: React.FormEvent) => {
    e.preventDefault()
    const q = term.trim()
    navigate(q ? `/hoc?q=${encodeURIComponent(q)}` : '/hoc')
  }

  return (
    <header className="topnav">
      <NavLink to="/hoc" className="brand" aria-label="Tập Huấn — trang chủ">
        <Logo size={36} />
      </NavLink>

      <nav className="topnav-links">
        <NavLink to="/hoc" end>Khoá học của tôi</NavLink>
        <NavLink to="/kham-pha">Khám phá</NavLink>
        <NavLink to="/ket-qua">Kết quả</NavLink>

        {manage && (
          <div className="nav-drop" ref={manageRef}>
            <button
              className={location.pathname.startsWith('/quan-tri') ? 'active' : undefined}
              onClick={() => setManageOpen((v) => !v)}
              aria-haspopup="menu"
              aria-expanded={manageOpen}
            >
              Quản lý
              <span style={{ display: 'inline-flex', transform: 'rotate(90deg)' }}>
                <IconChevron size={12} />
              </span>
            </button>
            {manageOpen && (
              <div className="user-pop" role="menu" onClick={() => setManageOpen(false)}>
                <NavLink to="/quan-tri">Bảng điều khiển</NavLink>
                <NavLink to="/quan-tri/chuong-trinh">Chương trình đào tạo</NavLink>
                <NavLink to="/quan-tri/cham-bai">Chấm bài</NavLink>
                {user.role === 'admin' && (
                  <>
                    <NavLink to="/quan-tri/nguoi-dung">Người dùng</NavLink>
                    <div className="pop-sep" />
                    <NavLink to="/quan-tri/cai-dat/google">Đăng nhập Google</NavLink>
                  </>
                )}
              </div>
            )}
          </div>
        )}
      </nav>

      <form className={`topnav-search ${searchOpen ? 'open' : ''}`} onSubmit={search} role="search">
        <input
          type="text"
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          placeholder="Bạn muốn học gì hôm nay?"
          aria-label="Tìm khoá học"
        />
        <button type="submit" aria-label="Tìm kiếm"><IconSearch /></button>
      </form>

      <div className="topnav-actions">
        <button
          className="icon-btn search-toggle"
          title="Tìm kiếm"
          onClick={() => setSearchOpen((v) => !v)}
        >
          <IconSearch size={20} />
        </button>

        <NavLink to="/ket-qua" className="icon-btn" title="Bài đang chờ chấm">
          <IconBell />
          {pendingCount > 0 && <span className="dot">{pendingCount}</span>}
        </NavLink>

        <button className="icon-btn" title="Trợ giúp" onClick={onToggleHelp}>
          <IconHelp />
        </button>

        <div className="user-menu" ref={menuRef}>
          <button onClick={() => setMenuOpen((v) => !v)} aria-haspopup="menu" aria-expanded={menuOpen}>
            <span className="avatar">
              {user.avatarUrl ? <img src={user.avatarUrl} alt="" /> : initials || '?'}
            </span>
          </button>

          {menuOpen && (
            <div className="user-pop" role="menu">
              <div className="who">
                <b>{user.fullName || 'Chưa đặt tên'}</b>
                <span>{user.email}</span>
                <span className="badge badge-primary" style={{ marginTop: 6 }}>{ROLE_LABEL[user.role]}</span>
              </div>
              {/* Trên màn hình hẹp thanh điều hướng bị ẩn, gom các mục vào đây. */}
              <div className="only-narrow">
                <NavLink to="/hoc" end onClick={() => setMenuOpen(false)}>Khoá học của tôi</NavLink>
                <NavLink to="/kham-pha" onClick={() => setMenuOpen(false)}>Khám phá khoá học</NavLink>
                <NavLink to="/ket-qua" onClick={() => setMenuOpen(false)}>Kết quả</NavLink>
                {manage && (
                  <>
                    <NavLink to="/quan-tri" onClick={() => setMenuOpen(false)}>Bảng điều khiển</NavLink>
                    <NavLink to="/quan-tri/chuong-trinh" onClick={() => setMenuOpen(false)}>
                      Chương trình đào tạo
                    </NavLink>
                    <NavLink to="/quan-tri/cham-bai" onClick={() => setMenuOpen(false)}>Chấm bài</NavLink>
                    {user.role === 'admin' && (
                      <>
                        <NavLink to="/quan-tri/nguoi-dung" onClick={() => setMenuOpen(false)}>Người dùng</NavLink>
                        <NavLink to="/quan-tri/cai-dat/google" onClick={() => setMenuOpen(false)}>
                          Đăng nhập Google
                        </NavLink>
                      </>
                    )}
                  </>
                )}
                <div className="pop-sep" />
              </div>

              <NavLink to="/doi-mat-khau" onClick={() => setMenuOpen(false)}>Đổi mật khẩu</NavLink>
              <button onClick={() => { signOut(); navigate('/dang-nhap') }}>Đăng xuất</button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}
