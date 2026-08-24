import { useEffect, useState } from 'react'

import { api } from '../api/client'
import type { AssignmentInput, LessonAttachmentInput, LessonInput } from '../api/client'
import type { ContentType, TreeNode } from '../api/types'
import { CONTENT_TYPE_LABEL } from '../api/types'
import { AssignmentResultsPanel } from './AssignmentResultsPanel'
import { IconChart, IconPlus, IconTrash } from './icons'
import { QuestionEditor } from './QuestionEditor'
import { RichTextEditor } from './RichTextEditor'
import { ErrorAlert, SuccessAlert } from './ui'

const KIND_LABEL: Record<string, string> = {
  folder: 'Thư mục',
  lesson: 'Bài học',
  assignment: 'Bài tập',
}

interface Props {
  node: TreeNode
  onSaved: (node: TreeNode) => void
  onDeleted: (nodeId: string) => void
  /** Vai trò Giám sát: hiển thị mọi thứ nhưng khoá toàn bộ form, ẩn nút Xoá/Lưu. */
  readOnly?: boolean
}

export function NodeEditor({ node, onSaved, onDeleted, readOnly = false }: Props) {
  const [title, setTitle] = useState(node.title)
  const [description, setDescription] = useState(node.description)
  const [isPublished, setIsPublished] = useState(node.isPublished)
  const [isLocked, setIsLocked] = useState(node.isLocked)

  const [contentType, setContentType] = useState<ContentType>(node.lesson?.contentType ?? 'video')
  // Ô nhập nhận link chia sẻ Drive; khi mở lại bài cũ ta hiển thị chính URL nhúng đã lưu.
  const [source, setSource] = useState(node.lesson?.embedUrl ?? '')
  const [duration, setDuration] = useState(node.lesson?.durationMinutes ?? 0)
  const [body, setBody] = useState(node.lesson?.body ?? '')
  const [attachments, setAttachments] = useState<LessonAttachmentInput[]>(node.lesson?.attachments ?? [])

  const [instructions, setInstructions] = useState(node.assignment?.instructions ?? '')
  const [timeLimit, setTimeLimit] = useState(node.assignment?.timeLimitMinutes ?? 0)
  const [maxAttempts, setMaxAttempts] = useState(node.assignment?.maxAttempts ?? 0)
  const [passScore, setPassScore] = useState(node.assignment?.passScore ?? 0)
  const [dueAt, setDueAt] = useState(toLocalInput(node.assignment?.dueAt ?? null))
  const [shuffle, setShuffle] = useState(node.assignment?.shuffleQuestions ?? false)

  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [tab, setTab] = useState<'questions' | 'results'>('questions')

  // Bài đọc tự soạn thay hẳn khung nhúng Google Drive bằng nội dung viết trong hệ thống.
  const isRichText = contentType === 'richtext'
  // Bài tài liệu không nhúng gì cả: chỉ là danh sách file cho học viên tải về.
  const isMaterials = contentType === 'materials'

  // Nạp lại form khi người dùng chọn một nút khác trên cây.
  useEffect(() => {
    setTitle(node.title)
    setDescription(node.description)
    setIsPublished(node.isPublished)
    setIsLocked(node.isLocked)
    setContentType(node.lesson?.contentType ?? 'video')
    setSource(node.lesson?.embedUrl ?? '')
    setDuration(node.lesson?.durationMinutes ?? 0)
    setBody(node.lesson?.body ?? '')
    setAttachments(node.lesson?.attachments ?? [])
    setInstructions(node.assignment?.instructions ?? '')
    setTimeLimit(node.assignment?.timeLimitMinutes ?? 0)
    setMaxAttempts(node.assignment?.maxAttempts ?? 0)
    setPassScore(node.assignment?.passScore ?? 0)
    setDueAt(toLocalInput(node.assignment?.dueAt ?? null))
    setShuffle(node.assignment?.shuffleQuestions ?? false)
    setError(null)
    setSaved(null)
    setTab('questions')
  }, [node])

  const save = async (e: React.FormEvent) => {
    e.preventDefault()
    setBusy(true)
    setError(null)
    setSaved(null)

    const payload: {
      title: string; description: string; isPublished: boolean; isLocked: boolean
      lesson?: LessonInput; assignment?: AssignmentInput
    } = { title: title.trim(), description, isPublished, isLocked }

    if (node.kind === 'lesson') {
      payload.lesson = {
        contentType,
        source: source.trim(),
        durationMinutes: Number(duration) || 0,
        body,
        attachments: attachments
          .map((a) => ({ name: a.name.trim(), url: a.url.trim() }))
          .filter((a) => a.name !== '' || a.url !== ''),
      }
    }
    if (node.kind === 'assignment') {
      payload.assignment = {
        instructions,
        timeLimitMinutes: Number(timeLimit) || 0,
        maxAttempts: Number(maxAttempts) || 0,
        passScore: Number(passScore) || 0,
        shuffleQuestions: shuffle,
        dueAt: dueAt ? new Date(dueAt).toISOString() : null,
      }
    }

    try {
      const updated = await api.updateNode(node.id, payload)
      onSaved(updated)
      setSaved('Đã lưu thay đổi')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được')
    } finally {
      setBusy(false)
    }
  }

  const setAttachment = (index: number, patch: Partial<LessonAttachmentInput>) => {
    setAttachments(attachments.map((a, i) => (i === index ? { ...a, ...patch } : a)))
  }

  const moveAttachment = (index: number, delta: number) => {
    const target = index + delta
    if (target < 0 || target >= attachments.length) return
    const next = [...attachments]
    const [item] = next.splice(index, 1)
    next.splice(target, 0, item)
    setAttachments(next)
  }

  const remove = async () => {
    const extra = node.children.length > 0 ? '\n\nToàn bộ nội dung bên trong cũng sẽ bị xoá.' : ''
    if (!confirm(`Xoá “${node.title}”?${extra}`)) return
    try {
      await api.deleteNode(node.id)
      onDeleted(node.id)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được')
    }
  }

  return (
    <div className="card">
      <div className="editor-head">
        <div>
          <span className="badge badge-primary">{KIND_LABEL[node.kind]}</span>
          {readOnly && <span className="badge" style={{ marginLeft: 6 }}>Chỉ xem</span>}
          <h2 style={{ marginTop: 8 }}>{node.title}</h2>
        </div>
        {!readOnly && <button className="btn btn-danger btn-sm" onClick={remove}>Xoá</button>}
      </div>

      <form onSubmit={save}>
        <fieldset disabled={readOnly} style={{ border: 'none', padding: 0, margin: 0 }}>
        <div className="editor-body">
          <ErrorAlert message={error} />
          <SuccessAlert message={saved} />

          <div className="field">
            <label htmlFor="n-title">Tiêu đề</label>
            <input id="n-title" type="text" value={title} onChange={(e) => setTitle(e.target.value)} required />
          </div>

          <div className="field">
            <label htmlFor="n-desc">Mô tả ngắn</label>
            <textarea id="n-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>

          {node.kind === 'lesson' && (
            <>
              <div className="row">
                <div className="field">
                  <label htmlFor="n-type">Loại nội dung</label>
                  <select
                    id="n-type" value={contentType}
                    onChange={(e) => setContentType(e.target.value as ContentType)}
                  >
                    {Object.entries(CONTENT_TYPE_LABEL).map(([value, label]) => (
                      <option key={value} value={value}>{label}</option>
                    ))}
                  </select>
                </div>
                <div className="field">
                  <label htmlFor="n-dur">Thời lượng (phút)</label>
                  <input
                    id="n-dur" type="number" min={0} value={duration}
                    onChange={(e) => setDuration(Number(e.target.value))}
                  />
                </div>
              </div>

              {isRichText ? (
                <div className="hint" style={{ marginTop: -4 }}>
                  Bài đọc tự soạn không cần file bên ngoài: nội dung viết thẳng bên dưới bằng
                  markdown, có công thức LaTeX và media nhúng.
                </div>
              ) : isMaterials ? (
                <div className="field">
                  <label>Tài liệu tải về</label>
                  {attachments.map((a, i) => (
                    <div className="attach-row" key={i}>
                      <input
                        type="text"
                        value={a.name}
                        onChange={(e) => setAttachment(i, { name: e.target.value })}
                        placeholder={`Tên tài liệu ${i + 1}`}
                      />
                      <input
                        type="text"
                        value={a.url}
                        onChange={(e) => setAttachment(i, { url: e.target.value })}
                        placeholder="https://… hoặc link Google Drive"
                      />
                      <button
                        type="button" className="btn btn-ghost btn-sm" title="Lên"
                        onClick={() => moveAttachment(i, -1)} disabled={i === 0}
                      >↑</button>
                      <button
                        type="button" className="btn btn-ghost btn-sm" title="Xuống"
                        onClick={() => moveAttachment(i, 1)} disabled={i === attachments.length - 1}
                      >↓</button>
                      <button
                        type="button" className="btn btn-ghost btn-sm" title="Xoá tài liệu"
                        onClick={() => setAttachments(attachments.filter((_, idx) => idx !== i))}
                      ><IconTrash /></button>
                    </div>
                  ))}
                  <button
                    type="button" className="btn btn-sm"
                    onClick={() => setAttachments([...attachments, { name: '', url: '' }])}
                  >
                    <IconPlus /> Thêm tài liệu
                  </button>
                  <div className="hint">
                    Mỗi dòng là một file học viên tải về: đặt tên hiển thị và dán link tải
                    (Google Drive hoặc URL bất kỳ). Với file Drive nhớ đặt quyền chia sẻ
                    “Bất kỳ ai có đường liên kết”. Bỏ trống tên thì hệ thống hiển thị chính link.
                  </div>
                </div>
              ) : (
                <>
                  <div className="field">
                    <label htmlFor="n-src">Link Google Drive</label>
                    <input
                      id="n-src" type="text" value={source}
                      onChange={(e) => setSource(e.target.value)}
                      placeholder="https://drive.google.com/file/d/…/view  hoặc  ID file"
                    />
                    <div className="hint">
                      Dán link chia sẻ từ Google Drive / Google Slides — hệ thống tự chuyển thành liên kết nhúng.
                      Nhớ đặt quyền chia sẻ file thành “Bất kỳ ai có đường liên kết”.
                    </div>
                  </div>

                  {node.lesson?.embedUrl && (
                    <div className="field">
                      <label>Xem trước</label>
                      <iframe
                        className={`embed-frame ${contentType === 'video' ? '' : 'doc'}`}
                        src={node.lesson.embedUrl}
                        title={node.title}
                        allow="autoplay; fullscreen"
                        allowFullScreen
                      />
                      <div className="hint mono" style={{ wordBreak: 'break-all' }}>{node.lesson.embedUrl}</div>
                    </div>
                  )}
                </>
              )}

              <div className="field">
                <label htmlFor="n-body">{isRichText ? 'Nội dung bài học' : 'Ghi chú cho học viên'}</label>
                <RichTextEditor
                  id="n-body"
                  value={body}
                  onChange={setBody}
                  disabled={readOnly}
                  rows={isRichText ? 20 : 8}
                  placeholder={isRichText
                    ? 'Viết nội dung bài học ở đây — dùng thanh công cụ hoặc gõ markdown trực tiếp.'
                    : 'Ghi chú thêm cho học viên (không bắt buộc).'}
                />
              </div>
            </>
          )}

          {node.kind === 'assignment' && (
            <>
              <div className="field">
                <label htmlFor="n-inst">Hướng dẫn làm bài</label>
                <textarea id="n-inst" value={instructions} onChange={(e) => setInstructions(e.target.value)} />
              </div>
              <div className="row">
                <div className="field">
                  <label htmlFor="n-time">Thời gian làm bài (phút)</label>
                  <input
                    id="n-time" type="number" min={0} value={timeLimit}
                    onChange={(e) => setTimeLimit(Number(e.target.value))}
                  />
                  <div className="hint">0 = không giới hạn</div>
                </div>
                <div className="field">
                  <label htmlFor="n-att">Số lượt làm tối đa</label>
                  <input
                    id="n-att" type="number" min={0} value={maxAttempts}
                    onChange={(e) => setMaxAttempts(Number(e.target.value))}
                  />
                  <div className="hint">0 = không giới hạn</div>
                </div>
              </div>
              <div className="row">
                <div className="field">
                  <label htmlFor="n-pass">Điểm đạt</label>
                  <input
                    id="n-pass" type="number" min={0} step="0.5" value={passScore}
                    onChange={(e) => setPassScore(Number(e.target.value))}
                  />
                </div>
                <div className="field">
                  <label htmlFor="n-due">Hạn nộp</label>
                  <input
                    id="n-due" type="datetime-local" value={dueAt}
                    onChange={(e) => setDueAt(e.target.value)}
                  />
                </div>
              </div>

              <div className="field">
                <label className="checkbox">
                  <input
                    type="checkbox" checked={shuffle}
                    onChange={(e) => setShuffle(e.target.checked)}
                  />
                  Trộn thứ tự câu hỏi và phương án
                </label>
                <div className="hint">
                  Mỗi lượt làm bài có một thứ tự riêng, giữ nguyên khi học viên tải lại trang.
                </div>
              </div>
            </>
          )}

          <label className="checkbox">
            <input type="checkbox" checked={isPublished} onChange={(e) => setIsPublished(e.target.checked)} />
            Hiển thị với học viên
          </label>
          <label className="checkbox">
            <input type="checkbox" checked={isLocked} onChange={(e) => setIsLocked(e.target.checked)} />
            Khoá nội dung với học viên
          </label>
        </div>
        </fieldset>

        {!readOnly && (
          <div className="editor-actions">
            <button className="btn btn-primary" disabled={busy}>{busy ? 'Đang lưu…' : 'Lưu thay đổi'}</button>
          </div>
        )}
      </form>

      {node.kind === 'assignment' && (
        <div className="editor-body" style={{ borderTop: '1px solid var(--border)' }}>
          <div className="toolbar">
            <button
              className={`btn btn-sm ${tab === 'questions' ? 'btn-primary' : ''}`}
              onClick={() => setTab('questions')}
            >
              Câu hỏi
            </button>
            <button
              className={`btn btn-sm ${tab === 'results' ? 'btn-primary' : ''}`}
              onClick={() => setTab('results')}
            >
              <IconChart /> Kết quả học viên
            </button>
          </div>

          {tab === 'questions'
            ? <QuestionEditor node={node} onChanged={onSaved} readOnly={readOnly} />
            : <AssignmentResultsPanel node={node} />}
        </div>
      )}
    </div>
  )
}

/** ISO 8601 → giá trị cho input[type=datetime-local] theo giờ địa phương. */
function toLocalInput(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
