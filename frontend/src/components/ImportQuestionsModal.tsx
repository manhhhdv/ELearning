import { useMemo, useState } from 'react'

import { api } from '../api/client'
import { QUESTION_TYPE_LABEL } from '../api/types'
import type { ImportFormat } from './questionImport'
import { SAMPLE_TABLE, SAMPLE_TEXT, parseQuestions } from './questionImport'
import { Modal } from './ui'

interface Props {
  nodeId: string
  onClose: () => void
  onImported: (count: number) => void
}

export function ImportQuestionsModal({ nodeId, onClose, onImported }: Props) {
  const [format, setFormat] = useState<ImportFormat>('text')
  const [raw, setRaw] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  // Phân tích lại mỗi lần gõ để người dùng thấy ngay mình nhập đúng chưa.
  const { questions, issues } = useMemo(() => parseQuestions(raw, format), [raw, format])

  const submit = async () => {
    setBusy(true)
    setError(null)
    try {
      const res = await api.importQuestions(nodeId, questions)
      onImported(res.imported)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không nhập được câu hỏi')
      setBusy(false)
    }
  }

  return (
    <Modal
      title="Nhập câu hỏi hàng loạt"
      onClose={onClose}
      wide
      footer={
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button
            className="btn btn-primary"
            onClick={submit}
            disabled={busy || questions.length === 0 || issues.length > 0}
          >
            {busy ? 'Đang nhập…' : `Nhập ${questions.length} câu hỏi`}
          </button>
        </>
      }
    >
      {error && <div className="alert alert-error">{error}</div>}

      <div className="field">
        <label>Định dạng nguồn</label>
        <div className="wrap-gap">
          <button
            type="button"
            className={`btn btn-sm ${format === 'text' ? 'btn-primary' : ''}`}
            onClick={() => setFormat('text')}
          >
            Soạn dạng văn bản
          </button>
          <button
            type="button"
            className={`btn btn-sm ${format === 'table' ? 'btn-primary' : ''}`}
            onClick={() => setFormat('table')}
          >
            Dán từ Excel / Google Sheets
          </button>
        </div>
      </div>

      <details style={{ marginBottom: 14 }}>
        <summary style={{ cursor: 'pointer', fontWeight: 600, fontSize: 13.5 }}>
          Cách viết {format === 'text' ? 'văn bản' : 'bảng'}
        </summary>
        {format === 'text' ? (
          <ul className="hint" style={{ paddingLeft: 18, marginTop: 8 }}>
            <li>Mỗi câu là một khối, các khối cách nhau bằng <b>một dòng trống</b>.</li>
            <li>Dòng đầu là đề bài. Đặt mã riêng bằng <code>[C10]</code> ở đầu, điểm bằng <code>[2đ]</code> ở cuối.</li>
            <li>Dòng bắt đầu bằng <code>*</code> là phương án đúng, <code>-</code> là phương án sai.</li>
            <li>Dòng bắt đầu bằng <code>&gt;</code> là giải thích.</li>
            <li>Khối không có dòng phương án nào sẽ thành <b>câu tự luận</b>.</li>
          </ul>
        ) : (
          <ul className="hint" style={{ paddingLeft: 18, marginTop: 8 }}>
            <li>Chọn vùng trong Excel rồi Ctrl+C, dán thẳng vào ô bên dưới.</li>
            <li>Thứ tự cột: <b>Mã · Nội dung · Điểm · PA1 · PA2 · PA3 · PA4 · Đáp án đúng · Giải thích</b>.</li>
            <li>Cột đáp án đúng ghi số thứ tự (<code>1,3</code>) hoặc chữ cái (<code>A,C</code>).</li>
            <li>Để trống hết các cột phương án thì câu đó là <b>tự luận</b>.</li>
            <li>Dòng tiêu đề được bỏ qua tự động.</li>
          </ul>
        )}
        <button
          type="button"
          className="btn btn-sm"
          style={{ marginTop: 8 }}
          onClick={() => setRaw(format === 'text' ? SAMPLE_TEXT : SAMPLE_TABLE)}
        >
          Điền dữ liệu mẫu
        </button>
      </details>

      <div className="field">
        <label htmlFor="import-src">Dán nội dung</label>
        <textarea
          id="import-src"
          value={raw}
          onChange={(e) => setRaw(e.target.value)}
          placeholder={format === 'text' ? 'Dán hoặc gõ câu hỏi…' : 'Dán vùng bảng đã sao chép…'}
          style={{ minHeight: 190, fontFamily: format === 'table' ? 'ui-monospace, monospace' : undefined }}
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

      {questions.length > 0 && (
        <>
          <div className="spread" style={{ marginBottom: 8 }}>
            <b>Xem trước</b>
            <span className="tiny muted">
              {questions.length} câu · tổng {Number(questions.reduce((s, q) => s + q.points, 0).toFixed(2))} điểm
            </span>
          </div>
          <div style={{ maxHeight: 240, overflowY: 'auto', border: '1px solid var(--border)', borderRadius: 8 }}>
            {questions.map((q, i) => (
              <div key={i} style={{ padding: '10px 13px', borderBottom: '1px solid var(--border)' }}>
                <div className="wrap-gap" style={{ marginBottom: 4 }}>
                  <span className="badge">{q.code || 'mã tự sinh'}</span>
                  <span className="badge badge-primary">{QUESTION_TYPE_LABEL[q.type]}</span>
                  <span className="tiny muted">{q.points} điểm</span>
                </div>
                <div style={{ fontSize: 13.5 }}>{q.prompt}</div>
                {q.options.length > 0 && (
                  <div className="tiny muted" style={{ marginTop: 3 }}>
                    {q.options.map((o) => (o.isCorrect ? `✓ ${o.content}` : o.content)).join('  ·  ')}
                  </div>
                )}
              </div>
            ))}
          </div>
        </>
      )}
    </Modal>
  )
}
