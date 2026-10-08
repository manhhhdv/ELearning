import { useState } from 'react'
import { api } from '../api/client'
import type { TreeNode } from '../api/types'
import { RichTextEditor } from './RichTextEditor'
import { FilePicker } from './FilePicker'
import { Modal, ErrorAlert } from './ui'
import { IconSparkles } from './icons'

export function GenerateLessonModal({ programId, folders, defaultParentId, onClose, onCreated }: {
  programId: string
  folders: { id: string; label: string }[]
  defaultParentId: string | null
  onClose: () => void
  onCreated: (node: TreeNode) => void
}) {
  const [topic, setTopic] = useState('')
  const [grade, setGrade] = useState('')
  const [objectives, setObjectives] = useState('')
  const [content, setContent] = useState('')
  // Nguồn nội dung tham khảo: dán tay hoặc tải lên tài liệu Word đã soạn sẵn.
  const [source, setSource] = useState<'text' | 'file'>('text')
  const [file, setFile] = useState<File | null>(null)
  const [duration, setDuration] = useState(45)
  const [parentId, setParentId] = useState(defaultParentId ?? '')
  const [draft, setDraft] = useState<{ title: string; body: string } | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const close = () => { if (!busy) onClose() }

  const generate = async () => {
    setBusy(true); setError(null)
    try {
      const fields = { programId, topic, grade, objectives, durationMinutes: duration }
      setDraft(source === 'file' && file
        ? await api.generateLessonFromFile(file, fields)
        : await api.generateLesson({ ...fields, content }))
    }
    catch (err) { setError(err instanceof Error ? err.message : 'Không tạo được bài giảng') }
    finally { setBusy(false) }
  }
  const save = async () => {
    if (!draft) return
    setBusy(true); setError(null)
    try {
      const node = await api.createNode(programId, {
        kind: 'lesson', parentId: parentId || null, title: draft.title.trim(), isPublished: false,
        lesson: { contentType: 'richtext', body: draft.body, durationMinutes: duration, source: '', attachments: [] },
      })
      onCreated(node)
    } catch (err) { setError(err instanceof Error ? err.message : 'Không lưu được bài giảng') }
    finally { setBusy(false) }
  }
  return <Modal title="Tạo bài giảng bằng AI" wide onClose={close} footer={<>
    <button className="btn" onClick={draft ? () => setDraft(null) : close} disabled={busy}>{draft ? 'Chỉnh yêu cầu' : 'Huỷ'}</button>
    {draft ? <button className="btn btn-primary" onClick={save} disabled={busy || !draft.title.trim() || !draft.body.trim()}>{busy ? 'Đang lưu…' : 'Lưu bản nháp vào lớp học'}</button>
      : <button className="btn btn-primary" onClick={generate} disabled={busy || !topic.trim() || duration < 1 || duration > 240 || (source === 'file' && !file)}><IconSparkles />{busy ? 'AI đang soạn…' : 'Tạo bài giảng'}</button>}
  </>}>
    <ErrorAlert message={error} />
    <p className="hint">AI soạn bài đọc cho học sinh. Bạn kiểm tra, chỉnh sửa và xuất bản sau khi lưu; bài mới chưa hiển thị với học sinh.</p>
    <div className="field"><label htmlFor="ai-lesson-parent">Chương / thư mục</label><select id="ai-lesson-parent" value={parentId} disabled={busy} onChange={(e) => setParentId(e.target.value)}><option value="">Cấp gốc lớp học</option>{folders.map((f) => <option key={f.id} value={f.id}>{f.label}</option>)}</select></div>
    <fieldset disabled={busy} style={{ border: 0, padding: 0, margin: 0 }}>
    {draft ? <>
      <div className="field"><label htmlFor="ai-lesson-title">Tiêu đề bài giảng</label><input id="ai-lesson-title" maxLength={300} value={draft.title} onChange={(e) => setDraft({ ...draft, title: e.target.value })} /></div>
      <label htmlFor="ai-lesson-body">Nội dung bài giảng</label>
      <RichTextEditor id="ai-lesson-body" value={draft.body} onChange={(body) => setDraft({ ...draft, body })} />
    </> : <>
      <div className="field"><label htmlFor="ai-lesson-topic">Chủ đề bài giảng</label><input id="ai-lesson-topic" value={topic} onChange={(e) => setTopic(e.target.value)} maxLength={300} placeholder="VD: Định luật Ôm và ứng dụng" autoFocus /></div>
      <div className="row">
        <div className="field"><label htmlFor="ai-lesson-grade">Trình độ / lớp</label><input id="ai-lesson-grade" value={grade} onChange={(e) => setGrade(e.target.value)} maxLength={100} placeholder="VD: Lớp 11" /></div>
        <div className="field"><label htmlFor="ai-lesson-duration">Thời lượng (phút)</label><input id="ai-lesson-duration" type="number" min={1} max={240} value={duration} onChange={(e) => setDuration(Number(e.target.value))} /></div>
      </div>
      <div className="field"><label htmlFor="ai-lesson-objectives">Mục tiêu cần đạt</label><textarea id="ai-lesson-objectives" value={objectives} onChange={(e) => setObjectives(e.target.value)} maxLength={4000} /></div>
      <div className="field">
        <label>Nội dung tham khảo (không bắt buộc)</label>
        <div className="wrap-gap">
          <button type="button" className={`btn btn-sm ${source === 'text' ? 'btn-primary' : ''}`} onClick={() => setSource('text')}>Dán nội dung</button>
          <button type="button" className={`btn btn-sm ${source === 'file' ? 'btn-primary' : ''}`} onClick={() => setSource('file')}>Tải lên tài liệu Word (.docx)</button>
        </div>
      </div>
      {source === 'text'
        ? <div className="field"><textarea id="ai-lesson-content" value={content} onChange={(e) => setContent(e.target.value)} rows={5} maxLength={20000} placeholder="Dán kiến thức nguồn hoặc đề cương cần bám sát…" /></div>
        : <div className="field">
            <FilePicker
              accept=".docx,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              label="Chọn tệp Word"
              onPick={setFile}
              hint="AI đọc nội dung tài liệu và soạn bài giảng bám sát. Tệp chỉ đi qua máy chủ, không được lưu lại. Tối đa 20 MB."
            />
          </div>}
    </>}
    </fieldset>
    {busy && <p role="status" className="hint">Đang xử lý, vui lòng giữ cửa sổ này mở…</p>}
  </Modal>
}
