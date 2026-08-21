import { useMemo, useState } from 'react'

import { api } from '../api/client'
import { SAMPLE_STRUCTURE, STRUCTURE_KIND_LABEL, parseStructure } from './structureImport'
import { Modal } from './ui'

interface Props {
  programId: string
  /** Danh sách thư mục phẳng để chọn nơi đặt nội dung nhập vào. */
  folders: { id: string; label: string }[]
  defaultParentId: string | null
  onClose: () => void
  onImported: (count: number) => void
}

export function ImportStructureModal({ programId, folders, defaultParentId, onClose, onImported }: Props) {
  const [parentId, setParentId] = useState<string | null>(defaultParentId)
  const [raw, setRaw] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [uploading, setUploading] = useState(false)

  // Phân tích lại mỗi lần gõ để người dùng thấy ngay cây sẽ được tạo ra.
  const { items, issues } = useMemo(() => parseStructure(raw), [raw])

  const handleFile = async (file: File) => {
    setUploading(true)
    setError(null)
    try {
      const res = await api.importStructureFile(programId, file)
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
      const res = await api.importStructure(programId, parentId, items)
      onImported(res.imported)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không nhập được cấu trúc')
      setBusy(false)
    }
  }

  const counts = items.reduce(
    (acc, it) => ({ ...acc, [it.kind]: (acc[it.kind] ?? 0) + 1 }),
    {} as Record<string, number>,
  )

  return (
    <Modal
      title="Nhập cấu trúc chương trình"
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
            {busy ? 'Đang nhập…' : `Nhập ${items.length} mục`}
          </button>
        </>
      }
    >
      {error && <div className="alert alert-error">{error}</div>}

      <div className="field">
        <label htmlFor="import-parent">Đặt vào</label>
        <select
          id="import-parent"
          value={parentId ?? ''}
          onChange={(e) => setParentId(e.target.value || null)}
        >
          <option value="">Gốc chương trình</option>
          {folders.map((f) => (
            <option key={f.id} value={f.id}>{f.label}</option>
          ))}
        </select>
      </div>

      <details style={{ marginBottom: 14 }}>
        <summary style={{ cursor: 'pointer', fontWeight: 600, fontSize: 13.5 }}>Cách viết bảng</summary>
        <ul className="hint" style={{ paddingLeft: 18, marginTop: 8 }}>
          <li>Thứ tự cột: <b>Cấp · Loại · Tiêu đề · Loại nội dung · Link/Nguồn · Thời lượng · Mô tả</b>.</li>
          <li><b>Cấp</b> là số bậc trong cây: 1 nằm ngay trong mục đã chọn ở trên, 2 nằm trong thư mục cấp 1 gần nhất phía trên.</li>
          <li><b>Loại</b> ghi <code>Thư mục</code>, <code>Bài học</code> hoặc <code>Bài tập</code>. Chỉ thư mục mới chứa được mục con.</li>
          <li>Bốn cột cuối chỉ dùng cho bài học: loại nội dung (<code>video/slide/document/pdf/link</code>), link chia sẻ Google Drive, số phút.</li>
          <li>Bài tập được tạo rỗng — mở từng bài để nhập câu hỏi sau.</li>
          <li>Dòng tiêu đề được bỏ qua tự động; mọi mục đều ở trạng thái đã xuất bản.</li>
        </ul>
        <button
          type="button"
          className="btn btn-sm"
          style={{ marginTop: 8 }}
          onClick={() => setRaw(SAMPLE_STRUCTURE)}
        >
          Điền dữ liệu mẫu
        </button>
      </details>

      <div className="field">
        <label htmlFor="import-structure-file">Tải file lên (.xlsx, .csv)</label>
        <input
          id="import-structure-file"
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
        <label htmlFor="import-structure-src">Hoặc dán từ Excel / Google Sheets</label>
        <textarea
          id="import-structure-src"
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
              {counts.folder ?? 0} thư mục · {counts.lesson ?? 0} bài học · {counts.assignment ?? 0} bài tập
            </span>
          </div>
          <div style={{ maxHeight: 260, overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 8 }}>
            {items.map((it, i) => (
              <div
                key={i}
                style={{ padding: '8px 13px', borderBottom: '1px solid var(--border)', paddingLeft: 13 + (it.level - 1) * 20 }}
              >
                <div className="wrap-gap" style={{ alignItems: 'center' }}>
                  <span className="badge badge-primary">{STRUCTURE_KIND_LABEL[it.kind]}</span>
                  <span style={{ fontSize: 13.5 }}>{it.title}</span>
                  {it.kind === 'lesson' && (
                    <span className="tiny muted">
                      {it.contentType}{it.durationMinutes > 0 ? ` · ${it.durationMinutes} phút` : ''}{it.source ? ' · có link' : ''}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </Modal>
  )
}
