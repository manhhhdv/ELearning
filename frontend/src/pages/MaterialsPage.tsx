import { useCallback, useEffect, useState } from 'react'

import { api } from '../api/client'
import type { MaterialInput } from '../api/client'
import type { Material, MaterialKind, MaterialList } from '../api/types'
import { MATERIAL_KIND_LABEL } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconDoc, IconPlus, IconTrash } from '../components/icons'
import { EmptyState, ErrorAlert, Loading, Modal, formatDate } from '../components/ui'

const emptyDraft = (): MaterialInput => ({
  title: '', description: '', category: '', url: '', kind: 'pdf', isPublished: true,
})

const KINDS: MaterialKind[] = ['pdf', 'slide', 'document', 'video', 'link']

/**
 * Kho tài liệu dùng chung toàn hệ thống (mục 4.4).
 *
 * Khác với tài liệu đính kèm trong từng bài học, tài liệu ở đây không thuộc
 * lớp nào: ai đăng nhập cũng tra cứu được. Giáo viên và quản trị viên thêm,
 * sửa, ẩn/hiện tài liệu ngay trên trang này.
 */
export function MaterialsPage() {
  const [data, setData] = useState<MaterialList | null>(null)
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editing, setEditing] = useState<Material | 'new' | null>(null)

  const load = useCallback(async () => {
    try {
      setData(await api.listMaterials(search, category))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được kho tài liệu')
    } finally {
      setLoading(false)
    }
  }, [search, category])

  // Gõ tới đâu lọc tới đó, nhưng chờ một nhịp để không gọi máy chủ mỗi ký tự.
  useEffect(() => {
    const timer = setTimeout(load, 250)
    return () => clearTimeout(timer)
  }, [load])

  const remove = async (m: Material) => {
    if (!confirm(`Xoá tài liệu "${m.title}"?`)) return
    try {
      await api.deleteMaterial(m.id)
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được tài liệu')
    }
  }

  if (loading) return <Loading label="Đang tải kho tài liệu…" />

  const canManage = data?.canManage ?? false
  const items = data?.items ?? []

  return (
    <>
      <PageHeader
        title="Kho tài liệu"
        subtitle="Tài liệu dùng chung cho toàn hệ thống, không thuộc lớp học nào"
        actions={canManage && (
          <button className="btn btn-primary" onClick={() => setEditing('new')}>
            <IconPlus /> Thêm tài liệu
          </button>
        )}
      />

      <div className="page-body learner">
        <ErrorAlert message={error} />

        <div className="toolbar">
          <input
            className="grow" type="search" placeholder="Tìm theo tên hoặc mô tả…"
            value={search} onChange={(e) => setSearch(e.target.value)}
          />
          <select value={category} onChange={(e) => setCategory(e.target.value)} style={{ width: 200 }}>
            <option value="">Tất cả phân loại</option>
            {data?.categories.map((c) => <option key={c} value={c}>{c}</option>)}
          </select>
        </div>

        {items.length === 0 ? (
          <EmptyState title="Chưa có tài liệu nào">
            <p className="hint">
              {canManage
                ? 'Bấm "Thêm tài liệu" để đăng tài liệu dùng chung đầu tiên.'
                : 'Giáo viên chưa đăng tài liệu dùng chung nào.'}
            </p>
          </EmptyState>
        ) : (
          <div className="material-grid">
            {items.map((m) => (
              <article className={`material-card ${m.isPublished ? '' : 'draft'}`} key={m.id}>
                <div className="material-kind">
                  <IconDoc size={16} /> {MATERIAL_KIND_LABEL[m.kind]}
                  {!m.isPublished && <span className="badge badge-warning">Bản nháp</span>}
                </div>

                <h3>
                  {/* Mở tab mới: tài liệu nằm ở Drive hoặc trang ngoài. */}
                  <a href={m.url} target="_blank" rel="noreferrer noopener">{m.title}</a>
                </h3>
                {m.description && <p className="material-desc">{m.description}</p>}

                <div className="material-foot">
                  {m.category && <span className="pill">{m.category}</span>}
                  <span className="tiny muted">{formatDate(m.createdAt)}</span>
                  {canManage && (
                    <>
                      <span className="grow" />
                      <button className="btn btn-sm" onClick={() => setEditing(m)}>Sửa</button>
                      <button className="btn btn-ghost btn-sm" onClick={() => remove(m)} title="Xoá">
                        <IconTrash />
                      </button>
                    </>
                  )}
                </div>
              </article>
            ))}
          </div>
        )}
      </div>

      {editing && (
        <MaterialModal
          material={editing === 'new' ? null : editing}
          categories={data?.categories ?? []}
          onClose={() => setEditing(null)}
          onSaved={async () => { setEditing(null); await load() }}
        />
      )}
    </>
  )
}

function MaterialModal({
  material, categories, onClose, onSaved,
}: {
  material: Material | null
  categories: string[]
  onClose: () => void
  onSaved: () => void
}) {
  const [draft, setDraft] = useState<MaterialInput>(material
    ? {
      title: material.title, description: material.description, category: material.category,
      url: material.url, kind: material.kind, isPublished: material.isPublished,
    }
    : emptyDraft())
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const set = <K extends keyof MaterialInput>(key: K, value: MaterialInput[K]) =>
    setDraft((prev) => ({ ...prev, [key]: value }))

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    try {
      if (material) await api.updateMaterial(material.id, draft)
      else await api.createMaterial(draft)
      onSaved()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được tài liệu')
      setBusy(false)
    }
  }

  return (
    <Modal
      title={material ? 'Sửa tài liệu' : 'Thêm tài liệu'}
      onClose={onClose}
      wide
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" onClick={submit} disabled={busy}>
            {busy ? 'Đang lưu…' : 'Lưu tài liệu'}
          </button>
        </>
      }
    >
      {error && <div className="alert alert-error">{error}</div>}

      <form onSubmit={submit}>
        <div className="field">
          <label htmlFor="m-title">Tên tài liệu</label>
          <input
            id="m-title" value={draft.title} required
            onChange={(e) => set('title', e.target.value)}
            placeholder="VD: Đề cương ôn tập học kỳ I"
          />
        </div>

        <div className="field">
          <label htmlFor="m-url">Đường dẫn</label>
          <input
            id="m-url" value={draft.url} required
            onChange={(e) => set('url', e.target.value)}
            placeholder="Dán link Google Drive hoặc link http(s)"
          />
          <div className="hint">
            Hệ thống lưu đường dẫn chứ không lưu bản sao tệp, nên nhớ đặt quyền chia sẻ cho
            người học xem được.
          </div>
        </div>

        <div className="row">
          <div className="field">
            <label htmlFor="m-kind">Loại tài liệu</label>
            <select id="m-kind" value={draft.kind} onChange={(e) => set('kind', e.target.value as MaterialKind)}>
              {KINDS.map((k) => <option key={k} value={k}>{MATERIAL_KIND_LABEL[k]}</option>)}
            </select>
          </div>
          <div className="field">
            <label htmlFor="m-cat">Phân loại</label>
            <input
              id="m-cat" value={draft.category} list="material-categories"
              onChange={(e) => set('category', e.target.value)}
              placeholder="VD: Vật lí 11"
            />
            <datalist id="material-categories">
              {categories.map((c) => <option key={c} value={c} />)}
            </datalist>
          </div>
        </div>

        <div className="field">
          <label htmlFor="m-desc">Mô tả</label>
          <textarea
            id="m-desc" rows={3} value={draft.description}
            onChange={(e) => set('description', e.target.value)}
            placeholder="Tài liệu này dùng cho việc gì…"
          />
        </div>

        <label className="checkbox">
          <input
            type="checkbox" checked={draft.isPublished}
            onChange={(e) => set('isPublished', e.target.checked)}
          />
          Hiện cho người học
        </label>
        {!draft.isPublished && (
          <div className="hint">Đang tắt: chỉ giáo viên và quản trị viên nhìn thấy tài liệu này.</div>
        )}
      </form>
    </Modal>
  )
}
