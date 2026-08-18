import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { api } from '../api/client'
import type { Program, ProgramStatus } from '../api/types'
import { STATUS_LABEL } from '../api/types'
import { isReadOnlyViewer, useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { IconPlus } from '../components/icons'
import { EmptyState, ErrorAlert, Loading, Modal, formatDate } from '../components/ui'

const STATUS_BADGE: Record<ProgramStatus, string> = {
  draft: 'badge',
  published: 'badge badge-success',
  archived: 'badge badge-warning',
}

export function ProgramsPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const readOnly = isReadOnlyViewer(user)
  const [programs, setPrograms] = useState<Program[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [creating, setCreating] = useState(false)

  const load = useCallback(async (term: string) => {
    setLoading(true)
    try {
      setPrograms(await api.listPrograms({ search: term }))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được danh sách')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    const timer = setTimeout(() => load(search), search ? 250 : 0)
    return () => clearTimeout(timer)
  }, [search, load])

  return (
    <>
      <PageHeader
        title="Chương trình đào tạo"
        subtitle={readOnly ? 'Xem danh sách chương trình đào tạo (chế độ Giám sát)' : 'Tạo và sắp xếp nội dung đào tạo theo cấu trúc cây'}
        actions={!readOnly && (
          <button className="btn btn-primary" onClick={() => setCreating(true)}>
            <IconPlus /> Chương trình mới
          </button>
        )}
      />

      <div className="page-body">
        <div className="toolbar">
          <input
            className="grow"
            type="text"
            placeholder="Tìm theo tên hoặc mã chương trình…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <ErrorAlert message={error} />

        {loading ? <Loading /> : programs.length === 0 ? (
          <div className="card">
            <EmptyState title="Chưa có chương trình nào">
              <p>Bấm “Chương trình mới” để tạo chương trình đào tạo đầu tiên.</p>
            </EmptyState>
          </div>
        ) : (
          <div className="card">
            <table>
              <thead>
                <tr>
                  <th>Mã</th>
                  <th>Tên chương trình</th>
                  <th>Trạng thái</th>
                  <th>Nội dung</th>
                  <th>Học viên</th>
                  <th>Cập nhật</th>
                </tr>
              </thead>
              <tbody>
                {programs.map((p) => (
                  <tr
                    key={p.id}
                    style={{ cursor: 'pointer' }}
                    onClick={() => navigate(`/quan-tri/chuong-trinh/${p.slug}`)}
                  >
                    <td className="mono tiny">{p.code}</td>
                    <td>
                      <b>{p.title}</b>
                      {p.description && <div className="tiny muted">{p.description}</div>}
                    </td>
                    <td><span className={STATUS_BADGE[p.status]}>{STATUS_LABEL[p.status]}</span></td>
                    <td className="tiny muted">{p.lessonCount} bài học · {p.assignmentCount} bài tập</td>
                    <td className="tiny muted">{p.enrollmentCount}</td>
                    <td className="tiny muted">{formatDate(p.updatedAt)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {creating && (
        <CreateProgramModal
          onClose={() => setCreating(false)}
          onCreated={(p) => navigate(`/quan-tri/chuong-trinh/${p.slug}`)}
        />
      )}
    </>
  )
}

function CreateProgramModal({
  onClose, onCreated,
}: {
  onClose: () => void
  onCreated: (p: Program) => void
}) {
  const [code, setCode] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [status, setStatus] = useState<ProgramStatus>('draft')
  const [allowSelfEnroll, setAllowSelfEnroll] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      onCreated(await api.createProgram({
        code: code.trim(), title: title.trim(), description, status, allowSelfEnroll,
      }))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tạo được chương trình')
      setBusy(false)
    }
  }

  return (
    <Modal
      title="Chương trình đào tạo mới"
      onClose={onClose}
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" form="create-program" disabled={busy}>
            {busy ? 'Đang tạo…' : 'Tạo chương trình'}
          </button>
        </>
      }
    >
      <ErrorAlert message={error} />
      <form id="create-program" onSubmit={submit}>
        <div className="field">
          <label htmlFor="p-code">Mã chương trình</label>
          <input
            id="p-code" type="text" value={code}
            onChange={(e) => setCode(e.target.value)}
            placeholder="VD: ATTP-2026" required
          />
          <div className="hint">Mã không trùng nhau, dùng để tra cứu nhanh.</div>
        </div>
        <div className="field">
          <label htmlFor="p-title">Tên chương trình</label>
          <input
            id="p-title" type="text" value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="VD: An toàn thực phẩm 2026" required
          />
        </div>
        <div className="field">
          <label htmlFor="p-desc">Mô tả</label>
          <textarea id="p-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
        </div>
        <div className="field">
          <label htmlFor="p-status">Trạng thái</label>
          <select id="p-status" value={status} onChange={(e) => setStatus(e.target.value as ProgramStatus)}>
            <option value="draft">Bản nháp — chỉ người quản lý thấy</option>
            <option value="published">Đã xuất bản — học viên được ghi danh sẽ thấy</option>
            <option value="archived">Lưu trữ</option>
          </select>
        </div>
        <label className="checkbox">
          <input
            type="checkbox" checked={allowSelfEnroll}
            onChange={(e) => setAllowSelfEnroll(e.target.checked)}
          />
          Cho học viên tự ghi danh
        </label>
        <div className="hint">
          Bật thì khoá học xuất hiện ở mục “Khám phá khoá học”, học viên tự bấm đăng ký.
        </div>
      </form>
    </Modal>
  )
}
