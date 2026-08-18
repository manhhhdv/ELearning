import { useState } from 'react'

import { api } from '../api/client'
import type { QuestionInput } from '../api/client'
import type { Question, QuestionType, TreeNode } from '../api/types'
import { QUESTION_TYPE_LABEL } from '../api/types'
import { ImportQuestionsModal } from './ImportQuestionsModal'
import { IconPlus, IconTrash, IconUpload } from './icons'
import { ErrorAlert } from './ui'

/** Bản nháp câu hỏi đang chỉnh trong form. */
interface Draft {
  code: string
  type: QuestionType
  prompt: string
  points: number
  explanation: string
  options: { content: string; isCorrect: boolean }[]
}

const emptyDraft = (): Draft => ({
  code: '',
  type: 'single_choice',
  prompt: '',
  points: 1,
  explanation: '',
  options: [{ content: '', isCorrect: false }, { content: '', isCorrect: false }],
})

const toDraft = (q: Question): Draft => ({
  code: q.code,
  type: q.type,
  prompt: q.prompt,
  points: q.points,
  explanation: q.explanation,
  options: q.options.length
    ? q.options.map((o) => ({ content: o.content, isCorrect: o.isCorrect }))
    : [{ content: '', isCorrect: false }, { content: '', isCorrect: false }],
})

export function QuestionEditor({
  node, onChanged, readOnly = false,
}: {
  node: TreeNode
  onChanged: (n: TreeNode) => void
  readOnly?: boolean
}) {
  const questions = node.assignment?.questions ?? []
  const [editingId, setEditingId] = useState<string | null>(null)
  const [draft, setDraft] = useState<Draft | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [importing, setImporting] = useState(false)

  const reload = async () => onChanged(await api.getNode(node.id))

  const startNew = () => { setEditingId('new'); setDraft(emptyDraft()); setError(null) }
  const startEdit = (q: Question) => { setEditingId(q.id); setDraft(toDraft(q)); setError(null) }
  const cancel = () => { setEditingId(null); setDraft(null); setError(null) }

  const save = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!draft || !editingId) return
    setBusy(true)
    setError(null)

    const payload: QuestionInput = {
      code: draft.code.trim(),
      type: draft.type,
      prompt: draft.prompt.trim(),
      points: Number(draft.points) || 1,
      explanation: draft.explanation,
      options: draft.type === 'essay' ? [] : draft.options.filter((o) => o.content.trim() !== ''),
    }

    try {
      if (editingId === 'new') await api.createQuestion(node.id, payload)
      else await api.updateQuestion(editingId, payload)
      await reload()
      cancel()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được câu hỏi')
    } finally {
      setBusy(false)
    }
  }

  const remove = async (q: Question) => {
    if (!confirm('Xoá câu hỏi này?')) return
    try {
      await api.deleteQuestion(q.id)
      await reload()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được câu hỏi')
    }
  }

  /** Đổi chỗ hai câu liền kề rồi gửi lại toàn bộ thứ tự. */
  const move = async (index: number, delta: number) => {
    const next = [...questions]
    const target = index + delta
    if (target < 0 || target >= next.length) return
    ;[next[index], next[target]] = [next[target], next[index]]
    try {
      await api.reorderQuestions(node.id, next.map((q) => q.id))
      await reload()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không sắp xếp được')
    }
  }

  const totalPoints = questions.reduce((sum, q) => sum + q.points, 0)

  return (
    <>
      <div className="spread" style={{ marginBottom: 14 }}>
        <div>
          <h3>Câu hỏi</h3>
          <div className="tiny muted">
            {questions.length} câu · tổng {Number(totalPoints.toFixed(2))} điểm
          </div>
        </div>
        {!readOnly && (
          <div className="wrap-gap">
            <button className="btn btn-sm" onClick={() => setImporting(true)}>
              <IconUpload /> Nhập hàng loạt
            </button>
            <button className="btn btn-sm" onClick={startNew} disabled={editingId === 'new'}>
              <IconPlus /> Thêm câu hỏi
            </button>
          </div>
        )}
      </div>

      <ErrorAlert message={error} />

      {questions.map((q, i) => (
        !readOnly && editingId === q.id && draft ? (
          <QuestionForm
            key={q.id}
            draft={draft} setDraft={setDraft} onSubmit={save} onCancel={cancel} busy={busy}
          />
        ) : (
          <div className="question-card" key={q.id}>
            <div className="question-card-head">
              <span className="badge">Câu {i + 1}</span>
              <span className="badge mono" title="Mã cố định của câu hỏi">{q.code}</span>
              <span className="badge badge-primary">{QUESTION_TYPE_LABEL[q.type]}</span>
              <span className="tiny muted grow">{Number(q.points)} điểm</span>
              {!readOnly && (
                <>
                  <button className="btn btn-ghost btn-sm" onClick={() => move(i, -1)} disabled={i === 0} title="Lên">↑</button>
                  <button className="btn btn-ghost btn-sm" onClick={() => move(i, 1)} disabled={i === questions.length - 1} title="Xuống">↓</button>
                  <button className="btn btn-sm" onClick={() => startEdit(q)}>Sửa</button>
                  <button className="btn btn-ghost btn-sm" onClick={() => remove(q)} title="Xoá"><IconTrash /></button>
                </>
              )}
            </div>
            <div style={{ whiteSpace: 'pre-wrap' }}>{q.prompt}</div>
            {q.type !== 'essay' && (
              <ul style={{ margin: '9px 0 0', paddingLeft: 20 }}>
                {q.options.map((o) => (
                  <li key={o.id} className={o.isCorrect ? '' : 'muted'} style={{ fontSize: 13.5 }}>
                    {o.content} {o.isCorrect && <span className="badge badge-success">đúng</span>}
                  </li>
                ))}
              </ul>
            )}
            {q.explanation && <div className="hint">Giải thích: {q.explanation}</div>}
          </div>
        )
      ))}

      {!readOnly && editingId === 'new' && draft && (
        <QuestionForm draft={draft} setDraft={setDraft} onSubmit={save} onCancel={cancel} busy={busy} />
      )}

      {questions.length === 0 && (readOnly || editingId !== 'new') && (
        <p className="muted tiny">Bài tập chưa có câu hỏi nào.</p>
      )}

      {importing && (
        <ImportQuestionsModal
          nodeId={node.id}
          onClose={() => setImporting(false)}
          onImported={async (count) => {
            setImporting(false)
            await reload()
            setError(null)
            alert(`Đã nhập ${count} câu hỏi.`)
          }}
        />
      )}
    </>
  )
}

function QuestionForm({
  draft, setDraft, onSubmit, onCancel, busy,
}: {
  draft: Draft
  setDraft: (d: Draft) => void
  onSubmit: (e: React.FormEvent) => void
  onCancel: () => void
  busy: boolean
}) {
  const set = <K extends keyof Draft>(key: K, value: Draft[K]) => setDraft({ ...draft, [key]: value })

  const setOption = (index: number, patch: Partial<Draft['options'][number]>) => {
    const options = draft.options.map((o, i) => (i === index ? { ...o, ...patch } : o))
    // Câu một đáp án: đánh dấu đúng ở đâu thì bỏ dấu ở mọi phương án còn lại.
    if (patch.isCorrect && draft.type === 'single_choice') {
      options.forEach((o, i) => { if (i !== index) o.isCorrect = false })
    }
    setDraft({ ...draft, options })
  }

  const changeType = (type: QuestionType) => {
    // Chuyển sang một đáp án thì chỉ giữ lại dấu đúng đầu tiên.
    if (type === 'single_choice') {
      let seen = false
      const options = draft.options.map((o) => {
        const keep = o.isCorrect && !seen
        if (o.isCorrect) seen = true
        return { ...o, isCorrect: keep }
      })
      setDraft({ ...draft, type, options })
      return
    }
    setDraft({ ...draft, type })
  }

  return (
    <form className="question-card" onSubmit={onSubmit} style={{ borderColor: 'var(--primary)' }}>
      <div className="row">
        <div className="field" style={{ maxWidth: 130 }}>
          <label>Mã câu hỏi</label>
          <input
            type="text" className="mono" value={draft.code}
            onChange={(e) => set('code', e.target.value)}
            placeholder="tự sinh"
          />
        </div>
        <div className="field">
          <label>Loại câu hỏi</label>
          <select value={draft.type} onChange={(e) => changeType(e.target.value as QuestionType)}>
            {Object.entries(QUESTION_TYPE_LABEL).map(([value, label]) => (
              <option key={value} value={value}>{label}</option>
            ))}
          </select>
        </div>
        <div className="field" style={{ maxWidth: 120 }}>
          <label>Điểm</label>
          <input
            type="number" min={0.5} step={0.5} value={draft.points}
            onChange={(e) => set('points', Number(e.target.value))}
          />
        </div>
      </div>

      <div className="field">
        <label>Nội dung câu hỏi</label>
        <textarea
          value={draft.prompt}
          onChange={(e) => set('prompt', e.target.value)}
          placeholder="Nhập đề bài…"
          required
        />
      </div>

      {draft.type !== 'essay' && (
        <div className="field">
          <label>Phương án trả lời</label>
          {draft.options.map((o, i) => (
            <div className="option-row" key={i}>
              <label className="checkbox mark" title="Đánh dấu là đáp án đúng">
                <input
                  type={draft.type === 'single_choice' ? 'radio' : 'checkbox'}
                  name="correct-option"
                  checked={o.isCorrect}
                  onChange={(e) => setOption(i, { isCorrect: e.target.checked })}
                />
              </label>
              <input
                type="text"
                value={o.content}
                onChange={(e) => setOption(i, { content: e.target.value })}
                placeholder={`Phương án ${i + 1}`}
              />
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => setDraft({ ...draft, options: draft.options.filter((_, idx) => idx !== i) })}
                disabled={draft.options.length <= 2}
                title="Xoá phương án"
              >
                <IconTrash />
              </button>
            </div>
          ))}
          <button
            type="button"
            className="btn btn-sm"
            onClick={() => setDraft({ ...draft, options: [...draft.options, { content: '', isCorrect: false }] })}
          >
            <IconPlus /> Thêm phương án
          </button>
          <div className="hint">Tick vào ô bên trái để đánh dấu đáp án đúng.</div>
        </div>
      )}

      <div className="field">
        <label>Giải thích {draft.type === 'essay' ? '/ đáp án gợi ý cho người chấm' : '(hiện sau khi chấm xong)'}</label>
        <textarea value={draft.explanation} onChange={(e) => set('explanation', e.target.value)} />
      </div>

      <div className="wrap-gap">
        <button className="btn btn-primary btn-sm" disabled={busy}>{busy ? 'Đang lưu…' : 'Lưu câu hỏi'}</button>
        <button type="button" className="btn btn-sm" onClick={onCancel}>Huỷ</button>
      </div>
    </form>
  )
}
