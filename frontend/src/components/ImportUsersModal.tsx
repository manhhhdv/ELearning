import { useMemo, useState } from 'react'

import { api } from '../api/client'
import type { UserImportReport, UserImportResult } from '../api/client'
import { ROLE_LABEL } from '../api/types'
import type { Role } from '../api/types'
import { SAMPLE_USERS, parseUsers } from './userImport'
import { Modal } from './ui'

interface Props {
  onClose: () => void
  /** Gọi ngay khi nhập xong để trang danh sách tải lại. */
  onImported: (report: UserImportReport) => void
}

const STATUS_LABEL: Record<UserImportResult['status'], string> = {
  created: 'Đã tạo',
  skipped: 'Bỏ qua',
  failed: 'Lỗi',
}

const STATUS_CLASS: Record<UserImportResult['status'], string> = {
  created: 'badge-success',
  skipped: 'badge-warning',
  failed: 'badge-danger',
}

export function ImportUsersModal({ onClose, onImported }: Props) {
  const [raw, setRaw] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [uploading, setUploading] = useState(false)
  // Có báo cáo nghĩa là đã nhập xong: đổi hẳn nội dung modal sang phần kết quả.
  const [report, setReport] = useState<UserImportReport | null>(null)

  // Phân tích lại mỗi lần gõ để người dùng thấy ngay danh sách sắp được tạo.
  const { items, issues } = useMemo(() => parseUsers(raw), [raw])

  const handleFile = async (file: File) => {
    setUploading(true)
    setError(null)
    try {
      const res = await api.importUsersFile(file)
      setRaw(res.text)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không đọc được file')
    } finally {
      setUploading(false)
    }
  }

  const submit = async () => {
    setBusy(true)
    setError(null)
    try {
      const res = await api.importUsers(items)
      setReport(res)
      onImported(res)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không nhập được danh sách')
    } finally {
      setBusy(false)
    }
  }

  if (report) {
    return (
      <ImportUsersReport report={report} onClose={onClose} />
    )
  }

  const counts = items.reduce(
    (acc, it) => ({ ...acc, [it.role]: (acc[it.role] ?? 0) + 1 }),
    {} as Record<string, number>,
  )

  return (
    <Modal
      title="Nhập tài khoản hàng loạt"
      onClose={onClose}
      wide
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button
            className="btn btn-primary"
            onClick={submit}
            disabled={busy || items.length === 0 || issues.length > 0}
          >
            {busy ? 'Đang nhập…' : `Nhập ${items.length} tài khoản`}
          </button>
        </>
      }
    >
      {error && <div className="alert alert-error">{error}</div>}

      <details style={{ marginBottom: 14 }}>
        <summary style={{ cursor: 'pointer', fontWeight: 600, fontSize: 13.5 }}>Cách viết bảng</summary>
        <ul className="hint" style={{ paddingLeft: 18, marginTop: 8 }}>
          <li>Thứ tự cột: <b>Email · Họ và tên · Vai trò · Mật khẩu</b>.</li>
          <li><b>Vai trò</b> ghi <code>Học viên</code>, <code>Giảng viên</code>, <code>Giám sát</code> hoặc <code>Quản trị viên</code>; để trống thì mặc định là Học viên.</li>
          <li><b>Mật khẩu</b> để trống thì hệ thống tự sinh ngẫu nhiên; nếu tự đặt thì cần ít nhất 8 ký tự, gồm cả chữ và số.</li>
          <li>Email đã có tài khoản sẽ được bỏ qua, không ghi đè dữ liệu cũ.</li>
          <li>Mọi tài khoản mới đều được yêu cầu đổi mật khẩu ở lần đăng nhập đầu tiên.</li>
          <li>Dòng tiêu đề được bỏ qua tự động; mỗi lần nhập tối đa 500 dòng.</li>
        </ul>
        <button
          type="button"
          className="btn btn-sm"
          style={{ marginTop: 8 }}
          onClick={() => setRaw(SAMPLE_USERS)}
        >
          Điền dữ liệu mẫu
        </button>
      </details>

      <div className="field">
        <label htmlFor="import-users-file">Tải file lên (.xlsx, .csv)</label>
        <input
          id="import-users-file"
          type="file"
          accept=".xlsx,.csv"
          disabled={uploading}
          onChange={(e) => {
            const file = e.target.files?.[0]
            e.target.value = ''
            if (file) void handleFile(file)
          }}
        />
        {uploading && <span className="tiny muted" style={{ marginLeft: 8 }}>Đang đọc file…</span>}
      </div>

      <div className="field">
        <label htmlFor="import-users-src">Hoặc dán từ Excel / Google Sheets</label>
        <textarea
          id="import-users-src"
          value={raw}
          onChange={(e) => setRaw(e.target.value)}
          placeholder="Dán vùng bảng đã sao chép… hoặc tải file lên ở trên"
          style={{ minHeight: 170, fontFamily: 'ui-monospace, monospace' }}
        />
      </div>

      {issues.length > 0 && (
        <div className="alert alert-error">
          <b>Cần sửa {issues.length} chỗ trước khi nhập:</b>
          <ul style={{ margin: '6px 0 0', paddingLeft: 18 }}>
            {issues.slice(0, 8).map((it, i) => (
              <li key={i}>Dòng {it.line}: {it.message}</li>
            ))}
            {issues.length > 8 && <li>… và {issues.length - 8} lỗi khác</li>}
          </ul>
        </div>
      )}

      {items.length > 0 && (
        <>
          <div className="spread" style={{ marginBottom: 8 }}>
            <b>Xem trước</b>
            <span className="tiny muted">
              {(['student', 'trainer', 'supervisor', 'admin'] as Role[])
                .filter((role) => counts[role])
                .map((role) => `${counts[role]} ${ROLE_LABEL[role].toLowerCase()}`)
                .join(' · ')}
            </span>
          </div>
          <div style={{ maxHeight: 260, overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 8 }}>
            {items.map((it, i) => (
              <div key={i} style={{ padding: '8px 13px', borderBottom: '1px solid var(--border)' }}>
                <div className="wrap-gap" style={{ alignItems: 'center' }}>
                  <span className="badge badge-primary">{ROLE_LABEL[it.role as Role]}</span>
                  <span style={{ fontSize: 13.5 }}>{it.fullName || '—'}</span>
                  <span className="tiny muted">{it.email}</span>
                  {it.password && <span className="tiny muted">· mật khẩu tự đặt</span>}
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </Modal>
  )
}

/** Bảng kết quả kèm mật khẩu — mật khẩu chỉ hiển thị đúng một lần nên cần tải về ngay. */
function ImportUsersReport({ report, onClose }: { report: UserImportReport; onClose: () => void }) {
  const created = report.results.filter((r) => r.status === 'created')

  return (
    <Modal
      title="Kết quả nhập tài khoản"
      onClose={onClose}
      wide
      footer={
        <>
          {created.length > 0 && (
            <button className="btn" onClick={() => downloadCredentials(created)}>
              Tải file mật khẩu (.csv)
            </button>
          )}
          <button className="btn btn-primary" onClick={onClose}>Đã ghi lại</button>
        </>
      }
    >
      <div className={`alert ${report.failed > 0 ? 'alert-error' : 'alert-info'}`}>
        Đã tạo <b>{report.created}</b> tài khoản
        {report.skipped > 0 && <> · bỏ qua <b>{report.skipped}</b></>}
        {report.failed > 0 && <> · lỗi <b>{report.failed}</b></>}.
        {created.length > 0 && ' Mật khẩu chỉ hiển thị một lần — hãy tải file về hoặc sao chép ngay.'}
      </div>

      <div style={{ maxHeight: 340, overflowY: 'auto' }}>
        <table>
          <thead>
            <tr>
              <th>Email</th><th>Họ tên</th><th>Vai trò</th><th>Mật khẩu</th><th>Kết quả</th>
            </tr>
          </thead>
          <tbody>
            {report.results.map((r) => (
              <tr key={r.row}>
                <td className="tiny">{r.email || '—'}</td>
                <td className="tiny">{r.fullName || '—'}</td>
                <td className="tiny muted">{ROLE_LABEL[r.role as Role] ?? r.role}</td>
                <td className="tiny mono">{r.password || '—'}</td>
                <td className="tiny">
                  <span className={`badge ${STATUS_CLASS[r.status]}`}>{STATUS_LABEL[r.status]}</span>
                  {r.message && <div className="muted">{r.message}</div>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Modal>
  )
}

/** Xuất file CSV thông tin đăng nhập; dùng dấu ; và BOM để Excel mở đúng tiếng Việt. */
function downloadCredentials(created: UserImportResult[]) {
  const cell = (v: string) => (/[";\n]/.test(v) ? `"${v.replace(/"/g, '""')}"` : v)
  const lines = [
    ['Email', 'Họ và tên', 'Vai trò', 'Mật khẩu'].join(';'),
    ...created.map((r) => [r.email, r.fullName, ROLE_LABEL[r.role as Role] ?? r.role, r.password].map(cell).join(';')),
  ]

  const blob = new Blob(['\uFEFF' + lines.join('\r\n')], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'tai-khoan-moi.csv'
  a.click()
  URL.revokeObjectURL(url)
}
