import { useState } from 'react'

import { api } from '../api/client'
import { useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { ErrorAlert, SuccessAlert } from '../components/ui'

export function ChangePasswordPage() {
  const { user, refresh } = useAuth()
  const [current, setCurrent] = useState('')
  const [next, setNext] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [done, setDone] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setDone(null)
    if (next !== confirm) {
      setError('Mật khẩu xác nhận không khớp')
      return
    }
    setBusy(true)
    try {
      await api.changePassword(current, next)
      setDone('Đã đổi mật khẩu thành công')
      setCurrent(''); setNext(''); setConfirm('')
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không đổi được mật khẩu')
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <PageHeader title="Đổi mật khẩu" subtitle={user?.email} />
      <div className="page-body">
        <div className="card card-pad" style={{ maxWidth: 460 }}>
          {user?.mustChangePassword && (
            <div className="alert alert-info">
              Tài khoản đang dùng mật khẩu do quản trị viên cấp. Hãy đổi sang mật khẩu riêng của bạn.
            </div>
          )}
          <ErrorAlert message={error} />
          <SuccessAlert message={done} />

          <form onSubmit={submit}>
            <div className="field">
              <label htmlFor="cur">Mật khẩu hiện tại</label>
              <input
                id="cur" type="password" autoComplete="current-password"
                value={current} onChange={(e) => setCurrent(e.target.value)}
                required={user?.hasPassword}
              />
              {!user?.hasPassword && (
                <div className="hint">Tài khoản đăng nhập bằng Google chưa có mật khẩu — để trống ô này.</div>
              )}
            </div>
            <div className="field">
              <label htmlFor="new">Mật khẩu mới</label>
              <input
                id="new" type="password" autoComplete="new-password"
                value={next} onChange={(e) => setNext(e.target.value)} required
              />
              <div className="hint">Tối thiểu 8 ký tự, gồm cả chữ và số.</div>
            </div>
            <div className="field">
              <label htmlFor="cfm">Nhập lại mật khẩu mới</label>
              <input
                id="cfm" type="password" autoComplete="new-password"
                value={confirm} onChange={(e) => setConfirm(e.target.value)} required
              />
            </div>
            <button className="btn btn-primary" disabled={busy}>
              {busy ? 'Đang lưu…' : 'Cập nhật mật khẩu'}
            </button>
          </form>
        </div>
      </div>
    </>
  )
}
