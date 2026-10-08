import { useEffect, useState } from 'react'

import { api } from '../api/client'
import type { QuestionInput } from '../api/client'
import type {
  TreeNode, AIQuestionKind, GeneratedQuestion, QuestionLevel, QuestionSpec, QuestionType,
} from '../api/types'
import { AI_QUESTION_KINDS, QUESTION_LEVELS, QUESTION_TYPE_LABEL } from '../api/types'
import { IconSparkles, IconTrash } from './icons'
import { FilePicker } from './FilePicker'
import { Modal } from './ui'

interface Props {
  nodeId: string
  /** Tên bài tập, dùng làm chủ đề mặc định cho đỡ phải gõ lại. */
  defaultTopic?: string
  onClose: () => void
  onImported: (count: number) => void
}

/** Dạng câu hỏi của lớp Service AI quy về dạng lưu trong ngân hàng câu hỏi. */
const toQuestionType = (kind: AIQuestionKind): QuestionType =>
  kind === 'multiple_choice' ? 'single_choice' : kind

/**
 * Chuyển một câu do AI sinh sang đúng dữ liệu lưu xuống máy chủ.
 * Ý nghĩa của `options` thay đổi theo từng dạng — xem chú thích ở `Question`.
 */
function toQuestionInput(q: GeneratedQuestion): QuestionInput {
  const base = {
    type: toQuestionType(q.type),
    prompt: q.question,
    points: q.points || 1,
    level: q.level,
    explanation: q.explanation,
  }

  switch (q.type) {
    case 'multiple_choice':
      return {
        ...base,
        options: (q.options ?? []).map((content, i) => ({ content, isCorrect: i === q.answer })),
      }
    case 'true_false':
      return {
        ...base,
        options: (q.statements ?? []).map((st) => ({ content: st.text, isCorrect: st.answer })),
      }
    case 'fill_blank':
      return {
        ...base,
        options: (q.answers ?? []).map((content) => ({ content, isCorrect: true })),
      }
    case 'essay':
      return { ...base, sampleAnswer: q.sampleAnswer ?? '', rubric: q.rubric ?? '', options: [] }
  }
}

export function GenerateQuestionsModal({ nodeId, defaultTopic = '', onClose, onImported }: Props) {
  // Nguồn nội dung: nhập chủ đề hay tải lên tài liệu PDF / Word (mục 3.4.2).
  const [source, setSource] = useState<'course' | 'topic' | 'file'>('course')
  const [programId, setProgramId] = useState('')
  const [lessons, setLessons] = useState<TreeNode[]>([])
  const [sourceNodeId, setSourceNodeId] = useState('')
  const [subject, setSubject] = useState('')
  const [grade, setGrade] = useState('')
  const [topic, setTopic] = useState(defaultTopic)
  const [content, setContent] = useState('')
  const [file, setFile] = useState<File | null>(null)
  const [levels, setLevels] = useState<string[]>([...QUESTION_LEVELS])
  const [counts, setCounts] = useState<Record<AIQuestionKind, number>>({
    multiple_choice: 5, true_false: 0, fill_blank: 0, essay: 0,
  })

  // Câu hỏi AI đã sinh, chờ giáo viên duyệt trước khi lưu (mục 4.7).
  const [drafts, setDrafts] = useState<GeneratedQuestion[] | null>(null)
  const [attempts, setAttempts] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  useEffect(() => {
    let active = true
    const flatten = (nodes: TreeNode[]): TreeNode[] => nodes.flatMap((n) => [n, ...flatten(n.children)])
    api.getNode(nodeId).then(async (node) => {
      const tree = await api.getTree(node.programId)
      if (!active) return
      setProgramId(node.programId)
      setLessons(flatten(tree).filter((n) => n.kind === 'lesson' && n.lesson?.body.trim()))
    }).catch((err) => { if (active) setError(err.message) })
    return () => { active = false }
  }, [nodeId])

  const total = Object.values(counts).reduce((a, b) => a + b, 0)
  const specs: QuestionSpec[] = AI_QUESTION_KINDS
    .filter(({ kind }) => counts[kind] > 0)
    .map(({ kind }) => ({ kind, count: counts[kind] }))

  const toggleLevel = (lv: string) =>
    setLevels((prev) => (prev.includes(lv) ? prev.filter((x) => x !== lv) : [...prev, lv]))

  const generate = async () => {
    setBusy(true)
    setError(null)
    try {
      const fields = { subject, grade, topic, levels, specs }
      const res = source === 'file' && file
        ? await api.generateQuestionsFromFile(file, fields)
        : await api.generateQuestions({ ...fields, content, ...(source === 'course' ? { programId, sourceNodeId: sourceNodeId || undefined } : {}) })
      setDrafts(res.questions)
      setAttempts(res.attempts)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tạo được câu hỏi')
    } finally {
      setBusy(false)
    }
  }

  const save = async () => {
    if (!drafts?.length) return
    setBusy(true)
    setError(null)
    try {
      const res = await api.importQuestions(nodeId, drafts.map(toQuestionInput))
      onImported(res.imported)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được câu hỏi')
      setBusy(false)
    }
  }

  const canGenerate = !busy
    && total > 0
    && levels.length > 0
    && (source === 'course' ? !!programId && lessons.length > 0 : source === 'file' ? file !== null : topic.trim() !== '' || content.trim() !== '')

  return (
    <Modal
      title="Tạo câu hỏi bằng AI"
      onClose={() => { if (!busy) onClose() }}
      wide
      footer={drafts ? (
        <>
          <button className="btn" onClick={() => setDrafts(null)} disabled={busy}>Quay lại</button>
          <button className="btn btn-primary" onClick={save} disabled={busy || drafts.length === 0}>
            {busy ? 'Đang lưu…' : `Lưu ${drafts.length} câu vào bài tập`}
          </button>
        </>
      ) : (
        <>
          <button className="btn" onClick={onClose}>Huỷ</button>
          <button className="btn btn-primary" onClick={generate} disabled={!canGenerate}>
            <IconSparkles /> {busy ? 'AI đang soạn…' : `Tạo ${total} câu hỏi`}
          </button>
        </>
      )}
    >
      {error && <div className="alert alert-error">{error}</div>}

      {drafts ? (
        <>
          <p className="hint" style={{ marginTop: 0 }}>
            Hệ thống đã kiểm tra cấu trúc và số lượng câu
            {attempts > 1 && <> (AI phải sinh lại {attempts - 1} lần mới đạt)</>}.
            {' '}Bạn hãy đọc lại nội dung, bỏ những câu chưa phù hợp rồi lưu. Sau khi lưu vẫn sửa được từng câu.
          </p>
          {drafts.map((q, i) => (
            <div className="question-card" key={i}>
              <div className="question-card-head">
                <span className="badge">Câu {i + 1}</span>
                <span className="badge badge-primary">{QUESTION_TYPE_LABEL[toQuestionType(q.type)]}</span>
                {q.level && <span className="badge">{q.level}</span>}
                <span className="tiny muted grow">{q.points || 1} điểm</span>
                <button
                  className="btn btn-ghost btn-sm"
                  title="Bỏ câu này"
                  onClick={() => setDrafts(drafts.filter((_, idx) => idx !== i))}
                >
                  <IconTrash />
                </button>
              </div>
              <div style={{ whiteSpace: 'pre-wrap' }}>{q.question}</div>

              {q.type === 'multiple_choice' && (
                <ul style={{ margin: '9px 0 0', paddingLeft: 20 }}>
                  {q.options?.map((o, oi) => (
                    <li key={oi} className={oi === q.answer ? '' : 'muted'} style={{ fontSize: 13.5 }}>
                      {o} {oi === q.answer && <span className="badge badge-success">đúng</span>}
                    </li>
                  ))}
                </ul>
              )}
              {q.type === 'true_false' && (
                <ul style={{ margin: '9px 0 0', paddingLeft: 20 }}>
                  {q.statements?.map((st, si) => (
                    <li key={si} style={{ fontSize: 13.5 }}>
                      {st.text}{' '}
                      <span className={st.answer ? 'badge badge-success' : 'badge'}>
                        {st.answer ? 'Đúng' : 'Sai'}
                      </span>
                    </li>
                  ))}
                </ul>
              )}
              {q.type === 'fill_blank' && (
                <div className="hint">Đáp án được chấp nhận: {q.answers?.join(' · ')}</div>
              )}
              {q.type === 'essay' && (
                <>
                  <div className="hint">Đáp án gợi ý: {q.sampleAnswer}</div>
                  <div className="hint">Tiêu chí chấm: {q.rubric}</div>
                </>
              )}
              <div className="hint">Giải thích: {q.explanation}</div>
            </div>
          ))}
          {drafts.length === 0 && <p className="muted tiny">Bạn đã bỏ hết câu hỏi. Bấm Quay lại để tạo mẻ mới.</p>}
        </>
      ) : (
        <>
          <div className="field">
            <label>Nguồn nội dung</label>
            <div className="wrap-gap">
              <button type="button" className={`btn btn-sm ${source === 'course' ? 'btn-primary' : ''}`} onClick={() => setSource('course')}>Bài giảng trong lớp học</button>
              <button
                type="button"
                className={`btn btn-sm ${source === 'topic' ? 'btn-primary' : ''}`}
                onClick={() => setSource('topic')}
              >
                Nhập chủ đề / nội dung
              </button>
              <button
                type="button"
                className={`btn btn-sm ${source === 'file' ? 'btn-primary' : ''}`}
                onClick={() => setSource('file')}
              >
                Tải lên tài liệu (PDF / Word)
              </button>
            </div>
          </div>

          <div className="row">
            <div className="field">
              <label>Lớp học</label>
              <input value={subject} onChange={(e) => setSubject(e.target.value)} placeholder="VD: Vật lí" />
            </div>
            <div className="field">
              <label>Lớp</label>
              <input value={grade} onChange={(e) => setGrade(e.target.value)} placeholder="VD: lớp 11" />
            </div>
          </div>

          <div className="field">
            <label>Chủ đề</label>
            <input value={topic} onChange={(e) => setTopic(e.target.value)} placeholder="VD: Định luật Ôm" />
          </div>

          {source === 'course' ? (
            <div className="field">
              <label htmlFor="question-source-lesson">Bài giảng nguồn</label>
              <select id="question-source-lesson" value={sourceNodeId} onChange={(e) => setSourceNodeId(e.target.value)}>
                <option value="">Các bài liên quan trong lớp học</option>
                {lessons.map((n) => <option key={n.id} value={n.id}>{n.title}</option>)}
              </select>
              <p className="hint">{lessons.length ? 'AI nhận nội dung văn bản đã lưu trong lớp học. Bài được chọn sẽ được ưu tiên làm nguồn.' : 'Chưa có bài giảng văn bản. Hãy tạo bài giảng trước hoặc chọn nhập nội dung / PDF.'}</p>
            </div>
          ) : source === 'topic' ? (
            <div className="field">
              <label>Nội dung bài học</label>
              <textarea
                value={content}
                onChange={(e) => setContent(e.target.value)}
                placeholder="Dán nội dung bài học để AI bám sát khi ra đề (không bắt buộc nếu đã có chủ đề)…"
                rows={5}
              />
            </div>
          ) : (
            <div className="field">
              <label>Tài liệu PDF hoặc Word (.docx)</label>
              <FilePicker
                accept="application/pdf,.pdf,.docx,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
                label="Chọn tệp"
                onPick={setFile}
                hint="AI đọc trực tiếp tài liệu và ra đề bám sát nội dung. Tệp chỉ đi qua máy chủ, không được lưu lại. Tối đa 12 MB với PDF, 20 MB với Word."
              />
            </div>
          )}

          <div className="field">
            <label>Số câu theo từng dạng</label>
            {AI_QUESTION_KINDS.map(({ kind, label }) => (
              <div className="option-row" key={kind}>
                <input
                  type="number"
                  min={0}
                  max={50}
                  style={{ maxWidth: 90 }}
                  value={counts[kind]}
                  onChange={(e) => setCounts({ ...counts, [kind]: Math.max(0, Number(e.target.value) || 0) })}
                />
                <span style={{ fontSize: 13.5 }}>{label}</span>
              </div>
            ))}
            <div className="hint">Tổng {total} câu. Mỗi lần tạo tối đa 50 câu.</div>
          </div>

          <div className="field">
            <label>Mức độ nhận thức</label>
            <div className="wrap-gap">
              {QUESTION_LEVELS.map((lv) => (
                <button
                  type="button"
                  key={lv}
                  className={`btn btn-sm ${levels.includes(lv) ? 'btn-primary' : ''}`}
                  onClick={() => toggleLevel(lv)}
                >
                  {lv}
                </button>
              ))}
            </div>
            {levels.length === 0 && <div className="hint">Hãy chọn ít nhất một mức độ.</div>}
          </div>

          <p className="hint">
            Câu hỏi do AI tạo là <b>bản nháp</b>. Hệ thống tự kiểm tra cấu trúc và số lượng, nhưng
            tính đúng đắn về kiến thức vẫn cần bạn duyệt trước khi giao cho học sinh.
          </p>
        </>
      )}
    </Modal>
  )
}

// `QuestionLevel` chỉ dùng cho kiểu dữ liệu; khai báo lại để tránh cảnh báo import thừa.
export type { QuestionLevel }
