import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Role } from '../api/types'
import { useAuth } from '../auth'
import { ErrorAlert, SuccessAlert } from '../components/ui'
import { IconGoogle } from '../components/icons'
import { Logo } from '../components/Logo'

/**
 * Vai trò người dùng chọn ở trang đăng nhập. Đây chỉ là lựa chọn để hệ thống
 * biết đưa người dùng vào khu vực nào; vai trò thật nằm ở tài khoản và do
 * quản trị viên cấp, nên hệ thống vẫn đối chiếu sau khi đăng nhập thành công.
 */
const ROLE_TABS: { value: Role; label: string }[] = [
  { value: 'student', label: 'Học sinh' },
  { value: 'trainer', label: 'Giáo viên' },
  { value: 'admin', label: 'Quản trị viên' },
]

/** Vai trò nào được phép vào khu vực của vai trò đã chọn ở tab. */
function roleMatches(picked: Role, actual: Role): boolean {
  if (picked === actual) return true
  // Admin vào được mọi khu vực; Giám sát dùng chung khu quản lý với giáo viên.
  if (actual === 'admin') return true
  if (actual === 'supervisor') return picked === 'trainer'
  return false
}

type Mode = 'login' | 'register' | 'forgot'

export function LoginPage() {
  const { signIn, signInWithToken, register } = useAuth()
  const [params, setParams] = useSearchParams()

  const [mode, setMode] = useState<Mode>('login')
  const [role, setRole] = useState<Role>('student')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [fullName, setFullName] = useState('')
  const [note, setNote] = useState('')
  const [remember, setRemember] = useState(true)

  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [config, setConfig] = useState({
    googleEnabled: false, signupEnabled: false, signupAllowedDomains: '',
  })

  useEffect(() => {
    api.authConfig().then(setConfig).catch(() => undefined)
  }, [])

  // Backend chuyển hướng về đây kèm token sau khi đăng nhập Google thành công.
  useEffect(() => {
    const token = params.get('token')
    const err = params.get('error')
    if (err) {
      setError(err)
      setParams({}, { replace: true })
      return
    }
    if (token) {
      setBusy(true)
      signInWithToken(token)
        .catch(() => setError('Không xác thực được phiên đăng nhập Google'))
        .finally(() => { setBusy(false); setParams({}, { replace: true }) })
    }
  }, [params, setParams, signInWithToken])

  const switchMode = (next: Mode) => {
    setMode(next)
    setError(null)
    setNotice(null)
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setNotice(null)
    setBusy(true)

    try {
      if (mode === 'forgot') {
        const res = await api.forgotPassword(email.trim(), note.trim())
        setNotice(res.message)
        setNote('')
        return
      }

      const user = mode === 'register'
        ? await register({ email: email.trim(), fullName: fullName.trim(), password }, remember)
        : await signIn(email.trim(), password, remember)

      // Đăng nhập đã thành công; chỉ nhắc nếu vai trò thật khác tab đã chọn.
      if (!roleMatches(role, user.role)) {
        setNotice(`Tài khoản này có vai trò ${ROLE_TABS.find((t) => t.value === user.role)?.label ?? user.role}. `
          + 'Hệ thống sẽ mở đúng khu vực tương ứng.')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thực hiện được, vui lòng thử lại')
    } finally {
      setBusy(false)
    }
  }

  const title = mode === 'register' ? 'Đăng ký tài khoản'
    : mode === 'forgot' ? 'Quên mật khẩu' : 'Đăng nhập'

  return (
    <div className="login-wrap">
      <div className="login-hero">
        <div className="login-hero-badge">
          <span className="login-hero-crest">TP</span>
          <div>
            <div className="school-tag">Trường THPT</div>
            <div className="school-name">Trần Phú - Hoàn Kiếm</div>
          </div>
        </div>
        <h1>Hệ thống <span>đào tạo trực tuyến</span> của nhà trường</h1>
        <p className="tagline">
          Nơi thầy cô giao bài, học sinh học tập và ôn luyện mọi lúc, mọi nơi.
        </p>
        <div className="login-hero-features">
          <div className="login-hero-feature">
            <span className="ico">📘</span>
            <div><b>Học tập cá nhân hoá</b><span>Theo tiến độ từng học sinh</span></div>
          </div>
          <div className="login-hero-feature">
            <span className="ico">📝</span>
            <div><b>Bài tập &amp; kiểm tra</b><span>Giao bài, chấm điểm trực tuyến</span></div>
          </div>
          <div className="login-hero-feature">
            <span className="ico">📊</span>
            <div><b>Theo dõi kết quả</b><span>Thống kê tiến độ rõ ràng</span></div>
          </div>
          <div className="login-hero-feature">
            <span className="ico">🏫</span>
            <div><b>Kết nối nhà trường</b><span>Giáo viên - học sinh - phụ huynh</span></div>
          </div>
        </div>
      </div>

      <div className="login-panel">
      <div className="card login-card">
        <div style={{ marginBottom: 20 }}>
          <Logo size={44} />
        </div>
        <h1>{title}</h1>
        <p className="lead">Hệ thống đào tạo trực tuyến</p>

        {mode === 'login' && (
          <div className="role-tabs" role="tablist" aria-label="Chọn vai trò">
            {ROLE_TABS.map((t) => (
              <button
                key={t.value}
                type="button"
                role="tab"
                aria-selected={role === t.value}
                className={`role-tab ${role === t.value ? 'on' : ''}`}
                onClick={() => setRole(t.value)}
              >
                {t.label}
              </button>
            ))}
          </div>
        )}

        <ErrorAlert message={error} />
        <SuccessAlert message={notice} />

        <form onSubmit={submit}>
          {mode === 'register' && (
            <div className="field">
              <label htmlFor="fullName">Họ và tên</label>
              <input
                id="fullName" type="text" autoComplete="name" value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                placeholder="Nguyễn Văn A" required
              />
            </div>
          )}

          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              type="email"
              autoComplete="username"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="ten@congty.vn"
              required
            />
            {mode === 'register' && config.signupAllowedDomains && (
              <div className="hint">Chỉ nhận email thuộc: {config.signupAllowedDomains}</div>
            )}
          </div>

          {mode !== 'forgot' && (
            <div className="field">
              <label htmlFor="password">Mật khẩu</label>
              <input
                id="password"
                type="password"
                autoComplete={mode === 'register' ? 'new-password' : 'current-password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
              {mode === 'register' && (
                <div className="hint">Ít nhất 8 ký tự, nên có cả chữ và số.</div>
              )}
            </div>
          )}

          {mode === 'forgot' && (
            <div className="field">
              <label htmlFor="note">Lời nhắn cho quản trị viên</label>
              <textarea
                id="note" value={note} rows={2}
                onChange={(e) => setNote(e.target.value)}
                placeholder="VD: em là học sinh lớp 11A2, quên mật khẩu ạ"
              />
              <div className="hint">
                Hệ thống chưa gửi email tự động. Yêu cầu của bạn sẽ được chuyển tới quản trị viên
                để cấp lại mật khẩu.
              </div>
            </div>
          )}

          {mode !== 'forgot' && (
            <label className="checkbox" style={{ marginBottom: 14 }}>
              <input
                type="checkbox"
                checked={remember}
                onChange={(e) => setRemember(e.target.checked)}
              />
              Ghi nhớ đăng nhập trên thiết bị này
            </label>
          )}

          <button className="btn btn-primary btn-block" disabled={busy}>
            {busy ? 'Đang xử lý…'
              : mode === 'register' ? 'Tạo tài khoản'
                : mode === 'forgot' ? 'Gửi yêu cầu' : 'Đăng nhập'}
          </button>
        </form>

        {mode === 'login' && config.googleEnabled && (
          <>
            <div className="divider">hoặc</div>
            {/* Điều hướng cả trang: luồng OAuth phải rời khỏi ứng dụng. */}
            <a className="btn btn-block btn-google" href="/api/auth/google/start">
              <IconGoogle /> Đăng nhập bằng Google
            </a>
          </>
        )}

        <div className="login-links">
          {mode === 'login' ? (
            <>
              <button type="button" className="linkish" onClick={() => switchMode('forgot')}>
                Quên mật khẩu?
              </button>
              {config.signupEnabled ? (
                <button type="button" className="linkish" onClick={() => switchMode('register')}>
                  Đăng ký tài khoản mới
                </button>
              ) : (
                <span className="hint">Chưa có tài khoản? Liên hệ quản trị viên để được cấp.</span>
              )}
            </>
          ) : (
            <button type="button" className="linkish" onClick={() => switchMode('login')}>
              ← Quay lại đăng nhập
            </button>
          )}
        </div>
      </div>
      </div>
    </div>
  )
}
