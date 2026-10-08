import { useCallback, useEffect, useState } from 'react'

import { api } from '../api/client'
import type { PasswordResetRequest } from '../api/types'
import { formatDate } from './ui'

/**
 * Danh sách yêu cầu quên mật khẩu do người dùng gửi từ trang đăng nhập.
 *
 * Hệ thống chưa gửi email tự động nên admin là người cấp lại mật khẩu: bấm
 * "Đặt lại mật khẩu" ở dòng tài khoản tương ứng trong bảng bên dưới, yêu cầu
 * sẽ tự đóng.
 */
export function PasswordResetPanel({ onPick }: {
  /** Gọi khi admin bấm vào một yêu cầu, để lọc nhanh tới tài khoản đó. */
  onPick: (email: string) => void
}) {
  const [items, setItems] = useState<PasswordResetRequest[]>([])
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState<string | null>(null)

  const load = useCallback(async () => {
    try {
      setItems(await api.listPasswordResets(true))
    } catch {
      // Không chặn cả trang Người dùng chỉ vì mục phụ này lỗi.
      setItems([])
    }
  }, [])

  useEffect(() => { load() }, [load])

  const resolve = async (id: string, status: 'done' | 'rejected') => {
    setBusy(id)
    setError(null)
    try {
      await api.resolvePasswordReset(id, status)
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không cập nhật được yêu cầu')
    } finally {
      setBusy(null)
    }
  }

  if (items.length === 0) return null

  return (
    <div className="panel" style={{ marginBottom: 16, borderLeft: '3px solid var(--warning)' }}>
      <div className="spread" style={{ marginBottom: 10 }}>
        <h3 style={{ margin: 0 }}>Yêu cầu đặt lại mật khẩu</h3>
        <span className="badge badge-warning">{items.length} đang chờ</span>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      <table className="ltable">
        <thead>
          <tr><th>Email</th><th>Lời nhắn</th><th>Gửi lúc</th><th /></tr>
        </thead>
        <tbody>
          {items.map((q) => (
            <tr key={q.id}>
              <td>
                {q.userId ? (
                  <button className="linkish" onClick={() => onPick(q.email)}>{q.email}</button>
                ) : (
                  <>
                    {q.email}
                    {/* Email không khớp tài khoản nào: có thể người dùng gõ nhầm. */}
                    <div className="tiny muted">không có tài khoản nào dùng email này</div>
                  </>
                )}
                {q.fullName && <div className="tiny muted">{q.fullName}</div>}
              </td>
              <td className="tiny">{q.note || '—'}</td>
              <td className="tiny muted">{formatDate(q.createdAt)}</td>
              <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                {q.userId && (
                  <>
                    <button className="btn btn-sm" onClick={() => onPick(q.email)}>Tìm tài khoản</button>{' '}
                  </>
                )}
                <button
                  className="btn btn-ghost btn-sm"
                  disabled={busy === q.id}
                  onClick={() => resolve(q.id, 'rejected')}
                  title="Bỏ qua yêu cầu này"
                >
                  Bỏ qua
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="hint" style={{ marginTop: 8 }}>
        Bấm <b>Tìm tài khoản</b> rồi dùng <b>Đặt lại mật khẩu</b> ở dòng tương ứng. Yêu cầu tự đóng
        sau khi bạn cấp mật khẩu mới.
      </div>
    </div>
  )
}
