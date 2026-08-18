import { useCallback, useEffect, useState } from 'react'

import { api } from '../api/client'
import type { IssuedCredentials } from '../api/client'
import type { Role, User } from '../api/types'
import { useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { IconPlus } from '../components/icons'
import { EmptyState, ErrorAlert, Loading, Modal, formatDate } from '../components/ui'

export function UsersPage() {
  const { user: me } = useAuth()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [creating, setCreating] = useState(false)
  // Mật khẩu thô chỉ hiện đúng một lần ngay sau khi tạo hoặc đặt lại.
  const [issued, setIssued] = useState<IssuedCredentials | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const res = await api.listUsers({ search, role: roleFilter })
      setUsers(res.items)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được danh sách')
    } finally {
      setLoading(false)
    }
  }, [search, roleFilter])

  useEffect(() => {
    const timer = setTimeout(load, search ? 250 : 0)
    return () => clearTimeout(timer)
  }, [load, search])

  const patch = async (id: string, body: { role?: string; isActive?: boolean }) => {
    try {
      await api.updateUser(id, body)
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không cập nhật được')
    }
  }

  const resetPassword = async (u: User) => {
    if (!confirm(`Cấp mật khẩu mới cho ${u.email}?`)) return
    try {
      setIssued(await api.resetPassword(u.id))
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không đặt lại được mật khẩu')
    }
  }

  const remove = async (u: User) => {
    if (!confirm(`Xoá tài khoản ${u.email}? Thao tác này không hoàn tác được.`)) return
    try {
      await api.deleteUser(u.id)
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được tài khoản')
    }
  }

  return (
    <>
      <PageHeader
        title="Người dùng"
        subtitle="Cấp tài khoản, phân vai trò và đặt lại mật khẩu"
        actions={
          <button className="btn btn-primary" onClick={() => setCreating(true)}>
            <IconPlus /> Cấp tài khoản
          </button>
        }
      />

      <div className="page-body">
        <div className="toolbar">
          <input
            className="grow" type="text" placeholder="Tìm theo tên hoặc email…"
            value={search} onChange={(e) => setSearch(e.target.value)}
          />
          <select value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)} style={{ width: 180 }}>
            <option value="">Tất cả vai trò</option>
            <option value="admin">Quản trị viên</option>
            <option value="trainer">Giảng viên</option>
            <option value="supervisor">Giám sát</option>
            <option value="student">Học viên</option>
          </select>
        </div>

        <ErrorAlert message={error} />

        {loading ? <Loading /> : (
          <div className="card">
            {users.length === 0 ? (
              <EmptyState title="Không tìm thấy người dùng nào" />
            ) : (
              <table>
                <thead>
                  <tr>
                    <th>Họ tên</th><th>Email</th><th>Vai trò</th>
                    <th>Đăng nhập</th><th>Trạng thái</th><th>Lần cuối</th><th></th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.id}>
                      <td>
                        <b>{u.fullName || '—'}</b>
                        {u.id === me?.id && <span className="badge badge-primary" style={{ marginLeft: 6 }}>bạn</span>}
                      </td>
                      <td className="muted">{u.email}</td>
                      <td>
                        <select
                          value={u.role}
                          disabled={u.id === me?.id}
                          onChange={(e) => patch(u.id, { role: e.target.value as Role })}
                          style={{ width: 140 }}
                        >
                          <option value="admin">Quản trị viên</option>
                          <option value="trainer">Giảng viên</option>
                          <option value="supervisor">Giám sát</option>
                          <option value="student">Học viên</option>
                        </select>
                      </td>
                      <td className="tiny muted">
                        {[u.hasPassword && 'Mật khẩu', u.hasGoogle && 'Google'].filter(Boolean).join(' + ') || '—'}
                        {u.mustChangePassword && <div className="badge badge-warning">chờ đổi mật khẩu</div>}
                      </td>
                      <td>
                        <span className={`badge ${u.isActive ? 'badge-success' : 'badge-danger'}`}>
                          {u.isActive ? 'Đang hoạt động' : 'Đã khoá'}
                        </span>
                      </td>
                      <td className="tiny muted">{formatDate(u.lastLoginAt)}</td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        <button className="btn btn-sm" onClick={() => resetPassword(u)}>Mật khẩu mới</button>{' '}
                        <button
                          className="btn btn-sm"
                          disabled={u.id === me?.id}
                          onClick={() => patch(u.id, { isActive: !u.isActive })}
                        >
                          {u.isActive ? 'Khoá' : 'Mở khoá'}
                        </button>{' '}
                        <button
                          className="btn btn-sm btn-danger"
                          disabled={u.id === me?.id}
                          onClick={() => remove(u)}
                        >
                          Xoá
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}
      </div>

      {creating && (
        <CreateUserModal
          onClose={() => setCreating(false)}
          onCreated={async (result) => { setCreating(false); setIssued(result); await load() }}
        />
      )}

      {issued && (
        <Modal
          title="Thông tin đăng nhập"
          onClose={() => setIssued(null)}
          footer={<button className="btn btn-primary" onClick={() => setIssued(null)}>Đã ghi lại</button>}
        >
          <div className="alert alert-info">
            Mật khẩu chỉ hiển thị một lần. Hãy gửi cho người dùng ngay bây giờ.
          </div>
          <div className="field">
            <label>Email</label>
            <input type="text" readOnly value={issued.user.email} />
          </div>
          <div className="field">
            <label>Mật khẩu</label>
            <input className="mono" type="text" readOnly value={issued.password} onFocus={(e) => e.target.select()} />
          </div>
          <p className="hint">Người dùng sẽ được nhắc đổi mật khẩu ở lần đăng nhập đầu tiên.</p>
        </Modal>
      )}
    </>
  )
}

function CreateUserModal({
  onClose, onCreated,
}: {
  onClose: () => void
  onCreated: (r: IssuedCredentials) => void
}) {
  const [email, setEmail] = useState('')
  const [fullName, setFullName] = useState('')
  const [role, setRole] = useState<Role>('student')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      onCreated(await api.createUser({ email: email.trim(), fullName: fullName.trim(), role, password }))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tạo được tài khoản')
      setBusy(false)
    }
  }

  return (
    <Modal
      title="Cấp tài khoản mới"
      onClose={onClose}
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" form="create-user" disabled={busy}>
            {busy ? 'Đang tạo…' : 'Cấp tài khoản'}
          </button>
        </>
      }
    >
      <ErrorAlert message={error} />
      <form id="create-user" onSubmit={submit}>
        <div className="field">
          <label htmlFor="u-email">Email</label>
          <input
            id="u-email" type="email" value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="ten@congty.vn" required autoFocus
          />
          <div className="hint">Nếu dùng email Google, người dùng có thể đăng nhập thẳng bằng Google.</div>
        </div>
        <div className="field">
          <label htmlFor="u-name">Họ và tên</label>
          <input id="u-name" type="text" value={fullName} onChange={(e) => setFullName(e.target.value)} />
        </div>
        <div className="field">
          <label htmlFor="u-role">Vai trò</label>
          <select id="u-role" value={role} onChange={(e) => setRole(e.target.value as Role)}>
            <option value="student">Học viên</option>
            <option value="trainer">Giảng viên</option>
            <option value="supervisor">Giám sát — chỉ xem, không sửa được nội dung</option>
            <option value="admin">Quản trị viên</option>
          </select>
        </div>
        <div className="field">
          <label htmlFor="u-pass">Mật khẩu</label>
          <input
            id="u-pass" type="text" value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="Để trống để hệ thống sinh ngẫu nhiên"
          />
        </div>
      </form>
    </Modal>
  )
}
