import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

import { api } from '../api/client'
import { useAuth } from '../auth'
import { ErrorAlert } from '../components/ui'
import { IconGoogle } from '../components/icons'
import { Logo } from '../components/Logo'

export function LoginPage() {
  const { signIn, signInWithToken } = useAuth()
  const [params, setParams] = useSearchParams()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [googleEnabled, setGoogleEnabled] = useState(false)

  useEffect(() => {
    api.authConfig().then((c) => setGoogleEnabled(c.googleEnabled)).catch(() => setGoogleEnabled(false))
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

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setBusy(true)
    try {
      await signIn(email.trim(), password)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Đăng nhập thất bại')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="login-wrap">
      <div className="card login-card">
        <div style={{ marginBottom: 20 }}>
          <Logo size={44} />
        </div>
        <h1>Đăng nhập</h1>
        <p className="lead">Hệ thống tập huấn và đào tạo trực tuyến</p>

        <ErrorAlert message={error} />

        <form onSubmit={submit}>
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
          </div>
          <div className="field">
            <label htmlFor="password">Mật khẩu</label>
            <input
              id="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button className="btn btn-primary btn-block" disabled={busy}>
            {busy ? 'Đang xử lý…' : 'Đăng nhập'}
          </button>
        </form>

        {googleEnabled && (
          <>
            <div className="divider">hoặc</div>
            {/* Điều hướng cả trang: luồng OAuth phải rời khỏi ứng dụng. */}
            <a className="btn btn-block btn-google" href="/api/auth/google/start">
              <IconGoogle /> Đăng nhập bằng Google
            </a>
          </>
        )}

        <p className="hint center" style={{ marginTop: 18 }}>
          Chưa có tài khoản? Liên hệ quản trị viên để được cấp.
        </p>
      </div>
    </div>
  )
}
