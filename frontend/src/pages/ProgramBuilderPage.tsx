import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Enrollment, Program, ProgramStatus, TreeNode, User } from '../api/types'
import { ROLE_LABEL, STATUS_LABEL } from '../api/types'
import { isReadOnlyViewer, useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { NodeEditor } from '../components/NodeEditor'
import { ProgramTree } from '../components/ProgramTree'
import { IconEye, IconPlus } from '../components/icons'
import { EmptyState, ErrorAlert, Loading, Modal, formatDate } from '../components/ui'

type Tab = 'content' | 'learners' | 'settings'

export function ProgramBuilderPage() {
  const { programSlug = '' } = useParams()
  const { user } = useAuth()
  // Giám sát xem được mọi thứ trong trang này nhưng không có quyền sửa/xoá/tạo mới.
  const readOnly = isReadOnlyViewer(user)
  const [tab, setTab] = useState<Tab>('content')
  const [program, setProgram] = useState<Program | null>(null)
  const [tree, setTree] = useState<TreeNode[]>([])
  const [selected, setSelected] = useState<TreeNode | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [adding, setAdding] = useState(false)

  // ID nội bộ, tra một lần từ slug trên URL rồi dùng lại cho mọi lời gọi API.
  const programId = program?.id ?? ''

  const loadTree = useCallback(async (keepSelectedId?: string | null) => {
    if (!programId) return []
    const nodes = await api.getTree(programId)
    setTree(nodes)
    if (keepSelectedId) {
      const found = findNode(nodes, keepSelectedId)
      // Nút đang mở có thể vừa bị xoá — khi đó đóng khung soạn thảo lại.
      setSelected(found ? await api.getNode(found.id) : null)
    }
    return nodes
  }, [programId])

  useEffect(() => {
    setLoading(true)
    api.getProgramBySlug(programSlug)
      .then(async (p) => {
        setProgram(p)
        setTree(await api.getTree(p.id))
        setError(null)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được chương trình'))
      .finally(() => setLoading(false))
  }, [programSlug])

  const selectNode = async (node: TreeNode) => {
    try {
      // Cây chỉ mang dữ liệu tóm tắt; lấy bản đầy đủ (kèm câu hỏi) khi mở ra sửa.
      setSelected(await api.getNode(node.id))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được nội dung')
    }
  }

  const handleMove = async (nodeId: string, parentId: string | null, position: number) => {
    try {
      await api.moveNode(nodeId, parentId, position)
      await loadTree(selected?.id ?? null)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không di chuyển được')
    }
  }

  if (loading) return <Loading />
  if (!program) return <div className="page-body"><ErrorAlert message={error ?? 'Không tìm thấy chương trình'} /></div>

  return (
    <>
      <PageHeader
        title={program.title}
        subtitle={`${program.code} · ${STATUS_LABEL[program.status]} · ${program.lessonCount} bài học · ${program.assignmentCount} bài tập${readOnly ? ' · Chế độ chỉ xem' : ''}`}
        actions={
          <>
            <Link className="btn" to="/quan-tri/chuong-trinh">← Danh sách</Link>
            <Link className="btn" to={`/xem-truoc/${program.slug}`} target="_blank" rel="noreferrer">
              <IconEye /> Xem trước
            </Link>
            {!readOnly && (
              <button className="btn btn-primary" onClick={() => setAdding(true)}>
                <IconPlus /> Thêm nội dung
              </button>
            )}
          </>
        }
      />

      <div className="page-body">
        <div className="toolbar">
          <button className={`btn ${tab === 'content' ? 'btn-primary' : ''}`} onClick={() => setTab('content')}>Nội dung</button>
          <button className={`btn ${tab === 'learners' ? 'btn-primary' : ''}`} onClick={() => setTab('learners')}>Học viên</button>
          <button className={`btn ${tab === 'settings' ? 'btn-primary' : ''}`} onClick={() => setTab('settings')}>Cài đặt</button>
        </div>

        <ErrorAlert message={error} />

        {tab === 'content' && (
          <div className="builder">
            <div className="card tree-panel">
              <div className="tree-panel-head">
                <b>Cấu trúc chương trình</b>
                {!readOnly && <span className="tiny faint">kéo để sắp xếp</span>}
              </div>
              <ProgramTree
                nodes={tree}
                selectedId={selected?.id ?? null}
                onSelect={selectNode}
                onMove={readOnly ? undefined : handleMove}
              />
            </div>

            {selected ? (
              <NodeEditor
                node={selected}
                readOnly={readOnly}
                onSaved={async (updated) => {
                  setSelected(updated)
                  setTree(await api.getTree(programId))
                }}
                onDeleted={async () => {
                  setSelected(null)
                  setTree(await api.getTree(programId))
                }}
              />
            ) : (
              <div className="card">
                <EmptyState title="Chọn một mục để xem">
                  <p>Bấm vào một mục trên cây bên trái{!readOnly && ', hoặc tạo mục mới bằng nút “Thêm nội dung”'}.</p>
                </EmptyState>
              </div>
            )}
          </div>
        )}

        {tab === 'learners' && <LearnersTab programId={programId} readOnly={readOnly} />}
        {tab === 'settings' && <SettingsTab program={program} onUpdated={setProgram} readOnly={readOnly} />}
      </div>

      {adding && (
        <AddNodeModal
          programId={programId}
          tree={tree}
          defaultParentId={selected?.kind === 'folder' ? selected.id : selected?.parentId ?? null}
          onClose={() => setAdding(false)}
          onCreated={async (node) => {
            setAdding(false)
            setTree(await api.getTree(programId))
            setSelected(await api.getNode(node.id))
          }}
        />
      )}
    </>
  )
}

function findNode(nodes: TreeNode[], id: string): TreeNode | null {
  for (const n of nodes) {
    if (n.id === id) return n
    const hit = findNode(n.children, id)
    if (hit) return hit
  }
  return null
}

/** Danh sách thư mục phẳng để chọn nơi đặt nội dung mới. */
function folderOptions(nodes: TreeNode[], depth = 0): { id: string; label: string }[] {
  return nodes.flatMap((n) =>
    n.kind === 'folder'
      ? [{ id: n.id, label: `${'— '.repeat(depth)}${n.title}` }, ...folderOptions(n.children, depth + 1)]
      : [],
  )
}

function AddNodeModal({
  programId, tree, defaultParentId, onClose, onCreated,
}: {
  programId: string
  tree: TreeNode[]
  defaultParentId: string | null
  onClose: () => void
  onCreated: (node: TreeNode) => void
}) {
  const [kind, setKind] = useState<'folder' | 'lesson' | 'assignment'>('folder')
  const [title, setTitle] = useState('')
  const [parentId, setParentId] = useState<string | null>(defaultParentId)
  const [source, setSource] = useState('')
  const [contentType, setContentType] = useState('video')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const folders = folderOptions(tree)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      const node = await api.createNode(programId, {
        parentId,
        kind,
        title: title.trim(),
        ...(kind === 'lesson'
          ? { lesson: { contentType, source: source.trim(), durationMinutes: 0, body: '' } }
          : {}),
        ...(kind === 'assignment'
          ? {
              assignment: {
                instructions: '', timeLimitMinutes: 0, maxAttempts: 0,
                passScore: 0, shuffleQuestions: false, dueAt: null,
              },
            }
          : {}),
      })
      onCreated(node)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tạo được nội dung')
      setBusy(false)
    }
  }

  return (
    <Modal
      title="Thêm nội dung"
      onClose={onClose}
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" form="add-node" disabled={busy}>
            {busy ? 'Đang tạo…' : 'Tạo'}
          </button>
        </>
      }
    >
      <ErrorAlert message={error} />
      <form id="add-node" onSubmit={submit}>
        <div className="field">
          <label>Loại</label>
          <div className="wrap-gap">
            {([
              ['folder', 'Thư mục / Chương'],
              ['lesson', 'Bài học'],
              ['assignment', 'Bài tập'],
            ] as const).map(([value, label]) => (
              <button
                type="button"
                key={value}
                className={`btn btn-sm ${kind === value ? 'btn-primary' : ''}`}
                onClick={() => setKind(value)}
              >
                {label}
              </button>
            ))}
          </div>
        </div>

        <div className="field">
          <label htmlFor="a-title">Tiêu đề</label>
          <input
            id="a-title" type="text" value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder={kind === 'folder' ? 'VD: Chương 1 — Kiến thức cơ bản' : 'VD: Bài 1 — Giới thiệu'}
            required autoFocus
          />
        </div>

        <div className="field">
          <label htmlFor="a-parent">Đặt trong</label>
          <select
            id="a-parent"
            value={parentId ?? ''}
            onChange={(e) => setParentId(e.target.value || null)}
          >
            <option value="">— Cấp gốc —</option>
            {folders.map((f) => <option key={f.id} value={f.id}>{f.label}</option>)}
          </select>
          <div className="hint">Chỉ thư mục mới chứa được nội dung con.</div>
        </div>

        {kind === 'lesson' && (
          <>
            <div className="field">
              <label htmlFor="a-ct">Loại nội dung</label>
              <select id="a-ct" value={contentType} onChange={(e) => setContentType(e.target.value)}>
                <option value="video">Video</option>
                <option value="slide">Slide trình chiếu</option>
                <option value="document">Tài liệu</option>
                <option value="pdf">PDF</option>
                <option value="link">Liên kết ngoài</option>
              </select>
            </div>
            <div className="field">
              <label htmlFor="a-src">Link Google Drive</label>
              <input
                id="a-src" type="text" value={source}
                onChange={(e) => setSource(e.target.value)}
                placeholder="Có thể dán sau khi tạo"
              />
            </div>
          </>
        )}
      </form>
    </Modal>
  )
}

function LearnersTab({ programId, readOnly }: { programId: string; readOnly: boolean }) {
  const [enrollments, setEnrollments] = useState<Enrollment[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [picking, setPicking] = useState(false)
  const enrolledIds = useMemo(() => new Set(enrollments.map((e) => e.userId)), [enrollments])

  const load = useCallback(async () => {
    setLoading(true)
    try {
      setEnrollments(await api.listEnrollments(programId))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được danh sách')
    } finally {
      setLoading(false)
    }
  }, [programId])

  useEffect(() => { load() }, [load])

  const remove = async (userId: string, name: string) => {
    if (!confirm(`Gỡ ${name} khỏi chương trình?`)) return
    await api.unenroll(programId, userId)
    await load()
  }

  if (loading) return <Loading />

  return (
    <>
      {!readOnly && (
        <div className="toolbar">
          <button className="btn btn-primary" onClick={() => setPicking(true)}>
            <IconPlus /> Ghi danh người dùng
          </button>
        </div>
      )}
      <ErrorAlert message={error} />

      <div className="card">
        {enrollments.length === 0 ? (
          <EmptyState title="Chưa có ai được ghi danh">
            <p>Ghi danh học viên để họ nhìn thấy chương trình này.</p>
          </EmptyState>
        ) : (
          <table>
            <thead>
              <tr><th>Họ tên</th><th>Email</th><th>Vai trò</th><th>Ngày ghi danh</th>{!readOnly && <th></th>}</tr>
            </thead>
            <tbody>
              {enrollments.map((e) => (
                <tr key={e.id}>
                  <td><b>{e.fullName || '—'}</b></td>
                  <td className="muted">{e.email}</td>
                  <td>
                    <span className={`badge ${e.role === 'trainer' ? 'badge-primary' : ''}`}>
                      {e.role === 'trainer' ? 'Giảng viên' : 'Học viên'}
                    </span>
                  </td>
                  <td className="tiny muted">{formatDate(e.enrolledAt)}</td>
                  {!readOnly && (
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-sm btn-danger" onClick={() => remove(e.userId, e.fullName || e.email)}>
                        Gỡ
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {picking && (
        <EnrollModal
          programId={programId}
          enrolled={enrolledIds}
          onClose={() => setPicking(false)}
          onDone={async () => { setPicking(false); await load() }}
        />
      )}
    </>
  )
}

function EnrollModal({
  programId, enrolled, onClose, onDone,
}: {
  programId: string
  enrolled: Set<string>
  onClose: () => void
  onDone: () => void
}) {
  const [users, setUsers] = useState<User[]>([])
  const [search, setSearch] = useState('')
  const [picked, setPicked] = useState<Set<string>>(new Set())
  const [role, setRole] = useState<'student' | 'trainer'>('student')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    const timer = setTimeout(() => {
      api.listUsers({ search })
        .then((r) => setUsers(r.items.filter((u) => !enrolled.has(u.id))))
        .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được người dùng'))
    }, search ? 250 : 0)
    return () => clearTimeout(timer)
  }, [search, enrolled])

  const toggle = (id: string) => {
    setPicked((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  const submit = async () => {
    setBusy(true)
    try {
      await api.enroll(programId, [...picked], role)
      onDone()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không ghi danh được')
      setBusy(false)
    }
  }

  return (
    <Modal
      title="Ghi danh vào chương trình"
      onClose={onClose}
      wide
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" onClick={submit} disabled={busy || picked.size === 0}>
            Ghi danh {picked.size > 0 && `(${picked.size})`}
          </button>
        </>
      }
    >
      <ErrorAlert message={error} />
      <div className="field">
        <label>Vai trò trong chương trình</label>
        <select value={role} onChange={(e) => setRole(e.target.value as 'student' | 'trainer')}>
          <option value="student">Học viên — chỉ học và làm bài</option>
          <option value="trainer">Giảng viên — sửa nội dung và chấm bài</option>
        </select>
      </div>
      <div className="field">
        <label>Tìm người dùng</label>
        <input type="text" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Tên hoặc email…" />
      </div>

      <div style={{ maxHeight: 280, overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 6 }}>
        {users.length === 0 ? (
          <p className="muted tiny" style={{ padding: 14, margin: 0 }}>Không còn người dùng nào để ghi danh.</p>
        ) : users.map((u) => (
          <label
            key={u.id}
            className="checkbox"
            style={{ padding: '9px 13px', borderBottom: '1px solid var(--border)' }}
          >
            <input type="checkbox" checked={picked.has(u.id)} onChange={() => toggle(u.id)} />
            <span className="grow">
              <b>{u.fullName || u.email}</b>
              <span className="tiny muted"> · {u.email} · {ROLE_LABEL[u.role]}</span>
            </span>
          </label>
        ))}
      </div>
    </Modal>
  )
}

function SettingsTab({
  program, onUpdated, readOnly,
}: {
  program: Program
  onUpdated: (p: Program) => void
  readOnly: boolean
}) {
  const [title, setTitle] = useState(program.title)
  const [code, setCode] = useState(program.code)
  const [description, setDescription] = useState(program.description)
  const [coverUrl, setCoverUrl] = useState(program.coverUrl)
  const [status, setStatus] = useState<ProgramStatus>(program.status)
  const [allowSelfEnroll, setAllowSelfEnroll] = useState(program.allowSelfEnroll)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState(false)
  const [busy, setBusy] = useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setSaved(false)
    try {
      onUpdated(await api.updateProgram(program.id, {
        title: title.trim(), code: code.trim(), description,
        coverUrl: coverUrl.trim(), status, allowSelfEnroll,
      }))
      setSaved(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được')
    } finally {
      setBusy(false)
    }
  }

  return (
    <form className="card card-pad" style={{ maxWidth: 560 }} onSubmit={submit}>
      <ErrorAlert message={error} />
      {saved && <div className="alert alert-success">Đã lưu thay đổi</div>}
      {readOnly && (
        <div className="alert alert-info">Bạn đang xem ở chế độ Giám sát — không sửa được cấu hình này.</div>
      )}

      <fieldset disabled={readOnly} style={{ border: 'none', padding: 0, margin: 0 }}>
      <div className="row">
        <div className="field">
          <label htmlFor="s-code">Mã</label>
          <input id="s-code" type="text" value={code} onChange={(e) => setCode(e.target.value)} required />
        </div>
        <div className="field">
          <label htmlFor="s-status">Trạng thái</label>
          <select id="s-status" value={status} onChange={(e) => setStatus(e.target.value as ProgramStatus)}>
            <option value="draft">Bản nháp</option>
            <option value="published">Đã xuất bản</option>
            <option value="archived">Lưu trữ</option>
          </select>
        </div>
      </div>
      <div className="field">
        <label htmlFor="s-title">Tên chương trình</label>
        <input id="s-title" type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
      </div>
      <div className="field">
        <label htmlFor="s-desc">Mô tả</label>
        <textarea id="s-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
      </div>
      <div className="field">
        <label htmlFor="s-cover">Ảnh bìa</label>
        <input
          id="s-cover" type="text" value={coverUrl}
          onChange={(e) => setCoverUrl(e.target.value)}
          placeholder="https://…/anh-bia.jpg"
        />
        <div className="hint">
          Dán link ảnh hiển thị trên thẻ khoá học. Để trống thì hệ thống dùng dải màu suy từ mã khoá.
        </div>
        {coverUrl.trim() && (
          <img
            src={coverUrl}
            alt="Xem trước ảnh bìa"
            style={{
              marginTop: 10, width: '100%', maxWidth: 320, aspectRatio: '16 / 9',
              objectFit: 'cover', borderRadius: 8, border: '1px solid var(--border)',
            }}
          />
        )}
      </div>
      <div className="field">
        <label className="checkbox">
          <input
            type="checkbox" checked={allowSelfEnroll}
            onChange={(e) => setAllowSelfEnroll(e.target.checked)}
          />
          Cho học viên tự ghi danh
        </label>
        <div className="hint">
          Bật thì khoá xuất hiện ở mục “Khám phá khoá học” để học viên tự đăng ký, và họ cũng tự rời
          được. Tắt thì chỉ quản trị viên hoặc giảng viên mới thêm được người vào khoá.
        </div>
      </div>

      </fieldset>

      {!readOnly && (
        <button className="btn btn-primary" disabled={busy}>{busy ? 'Đang lưu…' : 'Lưu thay đổi'}</button>
      )}
    </form>
  )
}
