import { useEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'

import { api } from '../api/client'
import type { QuestionInput } from '../api/client'
import { useAIStatus } from '../api/useAIStatus'
import type { Question, QuestionLevel, QuestionType, TreeNode } from '../api/types'
import { QUESTION_LEVELS, QUESTION_TYPE_LABEL, QUESTION_TYPES } from '../api/types'
import { ExportQuestionsButtons } from './ExportButtons'
import { GenerateQuestionsModal } from './GenerateQuestionsModal'
import { ImportQuestionsModal } from './ImportQuestionsModal'
import { IconPlus, IconSparkles, IconTrash, IconUpload } from './icons'
import { ErrorAlert } from './ui'

/** Bản nháp câu hỏi đang chỉnh trong form. */
interface Draft {
  code: string
  type: QuestionType
  prompt: string
  points: number
  level: QuestionLevel
  explanation: string
  sampleAnswer: string
  rubric: string
  options: { content: string; isCorrect: boolean }[]
}

const blankOptions = () => [{ content: '', isCorrect: false }, { content: '', isCorrect: false }]

const emptyDraft = (): Draft => ({
  code: '',
  type: 'single_choice',
  prompt: '',
  points: 1,
  level: '',
  explanation: '',
  sampleAnswer: '',
  rubric: '',
  options: blankOptions(),
})

const toDraft = (q: Question): Draft => ({
  code: q.code,
  type: q.type,
  prompt: q.prompt,
  points: q.points,
  level: q.level,
  explanation: q.explanation,
  sampleAnswer: q.sampleAnswer,
  rubric: q.rubric,
  options: q.options.length
    ? q.options.map((o) => ({ content: o.content, isCorrect: o.isCorrect }))
    : blankOptions(),
})

/** Nhãn của danh sách phương án, khác nhau theo từng dạng câu hỏi. */
const optionsLabel = (type: QuestionType) => {
  if (type === 'true_false') return 'Các phát biểu'
  if (type === 'fill_blank') return 'Các đáp án được chấp nhận'
  return 'Phương án trả lời'
}

export function QuestionEditor({
  node, onChanged, readOnly = false, toolbarTarget, onShowQuestions,
}: {
  node: TreeNode
  onChanged: (n: TreeNode) => void
  readOnly?: boolean
  /** Phần tử trên đầu khung soạn nơi đặt thanh công cụ (xuất Word, tạo AI, thêm câu…). */
  toolbarTarget: HTMLElement | null
  /** Chuyển về tab câu hỏi khi thao tác từ thanh công cụ cần hiện danh sách câu. */
  onShowQuestions?: () => void
}) {
  const questions = node.assignment?.questions ?? []
  const [editingId, setEditingId] = useState<string | null>(null)
  const [draft, setDraft] = useState<Draft | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [importing, setImporting] = useState(false)
  const [generating, setGenerating] = useState(false)
  // Nút tạo đề bằng AI chỉ hiện khi admin đã cấu hình nhà cung cấp AI.
  const aiEnabled = useAIStatus()?.enabled ?? false

  const reload = async () => onChanged(await api.getNode(node.id))

  const newFormRef = useRef<HTMLDivElement>(null)

  const startNew = () => { setEditingId('new'); setDraft(emptyDraft()); setError(null); onShowQuestions?.() }

  // Nút "Thêm câu hỏi" nằm trên đầu khung soạn, còn form nằm cuối danh sách: cuộn tới cho thấy ngay.
  useEffect(() => {
    if (editingId === 'new') newFormRef.current?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }, [editingId])
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
      level: draft.level,
      explanation: draft.explanation,
      sampleAnswer: draft.type === 'essay' ? draft.sampleAnswer.trim() : '',
      rubric: draft.type === 'essay' ? draft.rubric.trim() : '',
      // Câu điền khuyết: mọi dòng đều là một đáp án được chấp nhận.
      options: draft.type === 'essay'
        ? []
        : draft.options
          .filter((o) => o.content.trim() !== '')
          .map((o) => ({ ...o, isCorrect: draft.type === 'fill_blank' ? true : o.isCorrect })),
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
      <div style={{ marginBottom: 14 }}>
        <h3>Câu hỏi</h3>
        <div className="tiny muted">
          {questions.length} câu · tổng {Number(totalPoints.toFixed(2))} điểm
        </div>
      </div>

      {toolbarTarget && createPortal(
        <>
          {questions.length > 0 && (
            <ExportQuestionsButtons
              nodeId={node.id}
              onError={(message) => { setError(message); if (message) onShowQuestions?.() }}
            />
          )}
          {!readOnly && (
            <>
              {aiEnabled && (
                <button className="btn btn-sm btn-ai" onClick={() => setGenerating(true)}>
                  <IconSparkles /> Tạo bằng AI
                </button>
              )}
              <button className="btn btn-sm" onClick={() => setImporting(true)}>
                <IconUpload /> Nhập hàng loạt
              </button>
              <button className="btn btn-sm" onClick={startNew} disabled={editingId === 'new'}>
                <IconPlus /> Thêm câu hỏi
              </button>
            </>
          )}
        </>,
        toolbarTarget,
      )}

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
              {q.level && <span className="badge">{q.level}</span>}
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

            {q.type === 'true_false' && (
              <ul style={{ margin: '9px 0 0', paddingLeft: 20 }}>
                {q.options.map((o) => (
                  <li key={o.id} style={{ fontSize: 13.5 }}>
                    {o.content}{' '}
                    <span className={o.isCorrect ? 'badge badge-success' : 'badge'}>
                      {o.isCorrect ? 'Đúng' : 'Sai'}
                    </span>
                  </li>
                ))}
              </ul>
            )}

            {q.type === 'fill_blank' && (
              <div className="hint">
                Đáp án được chấp nhận: {q.options.map((o) => o.content).join(' · ')}
              </div>
            )}

            {(q.type === 'single_choice' || q.type === 'multi_choice') && (
              <ul style={{ margin: '9px 0 0', paddingLeft: 20 }}>
                {q.options.map((o) => (
                  <li key={o.id} className={o.isCorrect ? '' : 'muted'} style={{ fontSize: 13.5 }}>
                    {o.content} {o.isCorrect && <span className="badge badge-success">đúng</span>}
                  </li>
                ))}
              </ul>
            )}

            {q.type === 'essay' && (
              <>
                {q.sampleAnswer && <div className="hint">Đáp án gợi ý: {q.sampleAnswer}</div>}
                {q.rubric && <div className="hint">Tiêu chí chấm: {q.rubric}</div>}
              </>
            )}

            {q.explanation && <div className="hint">Giải thích: {q.explanation}</div>}
          </div>
        )
      ))}

      {!readOnly && editingId === 'new' && draft && (
        <div ref={newFormRef}>
          <QuestionForm draft={draft} setDraft={setDraft} onSubmit={save} onCancel={cancel} busy={busy} />
        </div>
      )}

      {questions.length === 0 && (readOnly || editingId !== 'new') && (
        <p className="muted tiny">Bài tập chưa có câu hỏi nào.</p>
      )}

      {generating && (
        <GenerateQuestionsModal
          nodeId={node.id}
          defaultTopic={node.title}
          onClose={() => setGenerating(false)}
          onImported={async (count) => {
            setGenerating(false)
            await reload()
            setError(null)
            alert(`Đã thêm ${count} câu hỏi do AI tạo.`)
          }}
        />
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

  const isChoice = draft.type === 'single_choice' || draft.type === 'multi_choice'
  const hasOptions = draft.type !== 'essay'

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
            {QUESTION_TYPES.map((value) => (
              <option key={value} value={value}>{QUESTION_TYPE_LABEL[value]}</option>
            ))}
            {/* Chỉ hiện dạng cũ khi câu hỏi đang ở dạng đó, để không soạn thêm câu mới. */}
            {draft.type === 'multi_choice' && (
              <option value="multi_choice">{QUESTION_TYPE_LABEL.multi_choice}</option>
            )}
          </select>
        </div>
        <div className="field" style={{ maxWidth: 160 }}>
          <label>Mức độ</label>
          <select value={draft.level} onChange={(e) => set('level', e.target.value as QuestionLevel)}>
            <option value="">Chưa phân loại</option>
            {QUESTION_LEVELS.map((lv) => <option key={lv} value={lv}>{lv}</option>)}
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
          placeholder={draft.type === 'fill_blank'
            ? 'Dùng ___ để đánh dấu chỗ trống. VD: Đơn vị của điện trở là ___.'
            : 'Nhập đề bài…'}
          required
        />
        {draft.type === 'fill_blank' && (
          <div className="hint">Đánh dấu chỗ cần điền bằng ba dấu gạch dưới: <b>___</b></div>
        )}
        {draft.type === 'true_false' && (
          <div className="hint">Câu dẫn chung, các phát biểu Đúng/Sai nhập ở phần bên dưới.</div>
        )}
      </div>

      {hasOptions && (
        <div className="field">
          <label>{optionsLabel(draft.type)}</label>
          {draft.options.map((o, i) => (
            <div className="option-row" key={i}>
              {/* Câu điền khuyết không có ô đánh dấu: mọi dòng đều là đáp án đúng. */}
              {draft.type !== 'fill_blank' && (
                <label
                  className="checkbox mark"
                  title={draft.type === 'true_false' ? 'Phát biểu này Đúng' : 'Đánh dấu là đáp án đúng'}
                >
                  <input
                    type={draft.type === 'single_choice' ? 'radio' : 'checkbox'}
                    name="correct-option"
                    checked={o.isCorrect}
                    onChange={(e) => setOption(i, { isCorrect: e.target.checked })}
                  />
                </label>
              )}
              <input
                type="text"
                value={o.content}
                onChange={(e) => setOption(i, { content: e.target.value })}
                placeholder={
                  draft.type === 'true_false' ? `Phát biểu ${i + 1}`
                    : draft.type === 'fill_blank' ? `Đáp án được chấp nhận ${i + 1}`
                      : `Phương án ${i + 1}`
                }
              />
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => setDraft({ ...draft, options: draft.options.filter((_, idx) => idx !== i) })}
                disabled={draft.options.length <= (draft.type === 'fill_blank' ? 1 : 2)}
                title="Xoá dòng"
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
            <IconPlus />{' '}
            {draft.type === 'true_false' ? 'Thêm phát biểu'
              : draft.type === 'fill_blank' ? 'Thêm đáp án' : 'Thêm phương án'}
          </button>
          <div className="hint">
            {isChoice && 'Tick vào ô bên trái để đánh dấu đáp án đúng.'}
            {draft.type === 'true_false' && 'Tick vào ô bên trái nếu phát biểu đó Đúng; để trống nghĩa là Sai. Học sinh được chấm điểm theo từng phát biểu.'}
            {draft.type === 'fill_blank' && 'Mỗi dòng là một cách trả lời được chấp nhận. Khi chấm, hệ thống bỏ qua khác biệt hoa/thường, khoảng trắng thừa và dấu câu ở cuối.'}
          </div>
        </div>
      )}

      {draft.type === 'essay' && (
        <>
          <div className="field">
            <label>
              Đáp án gợi ý<span style={{ color: 'var(--c-danger, #dc2626)' }}> *</span>
            </label>
            <textarea
              value={draft.sampleAnswer}
              onChange={(e) => set('sampleAnswer', e.target.value)}
              placeholder="Những ý chính bài làm cần có…"
              required
            />
          </div>
          <div className="field">
            <label>
              Tiêu chí chấm và thang điểm<span style={{ color: 'var(--c-danger, #dc2626)' }}> *</span>
            </label>
            <textarea
              value={draft.rubric}
              onChange={(e) => set('rubric', e.target.value)}
              placeholder="VD: Nêu đúng công thức (1 điểm); trình bày đủ các bước (2 điểm)…"
              required
            />
            <div className="hint">AI dựa vào đáp án gợi ý và tiêu chí này để đề xuất điểm cho bạn duyệt.</div>
          </div>
        </>
      )}

      <div className="field">
        <label>
          Giải thích (hiện sau khi chấm xong)
          <span style={{ color: 'var(--c-danger, #dc2626)' }}> *</span>
        </label>
        <textarea value={draft.explanation} onChange={(e) => set('explanation', e.target.value)} required />
      </div>

      <div className="wrap-gap">
        <button className="btn btn-primary btn-sm" disabled={busy}>{busy ? 'Đang lưu…' : 'Lưu câu hỏi'}</button>
        <button type="button" className="btn btn-sm" onClick={onCancel}>Huỷ</button>
      </div>
    </form>
  )
}
