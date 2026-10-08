import { useEffect, useRef, useState } from 'react'
import { useLocation } from 'react-router-dom'

import { useAIStatus } from '../api/useAIStatus'
import { useAuth } from '../auth'
import { AI_CHAT_PATH } from './aiChatPath'
import { RichContent } from './RichContent'
import { chatSuggestions, useAIChat } from './useAIChat'
import { IconChat, IconPlus, IconTrash } from './icons'

/** Context is loaded by the server using the course and lesson IDs. */
export function AIAssistant({ subject, grade, programId, nodeId, lessonTitle }: {
  subject?: string
  grade?: string
  programId?: string
  nodeId?: string
  lessonTitle?: string
} = {}) {
  const { user } = useAuth()
  const aiStatus = useAIStatus()
  const location = useLocation()
  const [open, setOpen] = useState(false)
  const chat = useAIChat({ subject, grade, programId, nodeId })
  const {
    history, messages, input, setInput, busy, error, sources,
    loadHistory, startNew, openConversation, removeConversation, send, cancelPending,
  } = chat
  const inputRef = useRef<HTMLTextAreaElement>(null)
  const triggerRef = useRef<HTMLButtonElement>(null)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => cancelPending, [cancelPending])

  useEffect(() => {
    if (!open) return
    loadHistory()
    inputRef.current?.focus()
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') { setOpen(false); triggerRef.current?.focus() }
    }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [open, loadHistory])

  // Luôn cuộn tới lượt trả lời mới nhất.
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'instant' : 'smooth' })
  }, [messages, busy])

  if (!user || !aiStatus?.enabled) return null
  // Trên trang trợ lý AI đã có sẵn khung trò chuyện toàn màn hình, nút nổi chỉ gây rối.
  if (location.pathname === AI_CHAT_PATH) return null

  const isStudent = user.role === 'student'
  const suggestions = chatSuggestions(isStudent, programId)

  return (
    <>
      <button
        ref={triggerRef}
        className="ai-fab"
        aria-expanded={open}
        aria-controls="ai-assistant-panel"
        onClick={() => setOpen((v) => !v)}
        aria-label="Trợ lý AI"
        title="Trợ lý AI"
      >
        <IconChat size={22} />
      </button>

      {open && (
        <div id="ai-assistant-panel" className="ai-panel" role="dialog" aria-label="Trợ lý AI">
          <div className="ai-head">
            <div>
              <b>{isStudent ? 'Trợ lý học tập AI' : 'Trợ lý AI cho giáo viên'}</b>
              <div className="tiny muted">
                {isStudent
                  ? 'Hỏi bài, xin giải thích từng bước'
                  : 'Ý tưởng bài dạy, hoạt động, câu hỏi'}
              </div>
            </div>
            <button className="btn btn-ghost btn-sm" onClick={startNew} disabled={busy} title="Cuộc trò chuyện mới" aria-label="Cuộc trò chuyện mới">
              <IconPlus />
            </button>
            <button className="btn btn-ghost btn-sm" onClick={() => setOpen(false)} aria-label="Đóng">✕</button>
          </div>

          {programId && <div className="ai-context"><b>{subject}</b><span>{lessonTitle ?? 'Nội dung lớp học'}</span><small>Ưu tiên bài đang mở và các bài liên quan bạn được phép xem.</small></div>}
          <div className="ai-body" role="log" aria-live="polite" aria-busy={busy}>
            {messages.length === 0 && (
              <>
                <p className="hint" style={{ marginTop: 0 }}>
                  {isStudent
                    ? 'Trợ lý sẽ gợi ý và giải thích từng bước thay vì đưa ngay đáp án. Đây là công cụ hỗ trợ, khi cần chốt kiến thức em hãy hỏi thêm thầy cô.'
                    : 'Trợ lý hỗ trợ tìm ý tưởng bài học, thiết kế hoạt động và gợi ý câu hỏi kiểm tra.'}
                </p>
                <div className="wrap-gap" style={{ marginBottom: 14 }}>
                  {suggestions.map((sg) => (
                    <button key={sg} className="btn btn-sm" onClick={() => setInput(sg)}>{sg}</button>
                  ))}
                </div>

                {history.length > 0 && (
                  <>
                    <div className="tiny muted" style={{ marginBottom: 6 }}>Cuộc trò chuyện gần đây</div>
                    {history.slice(0, 8).map((c) => (
                      <div className="ai-history-row" key={c.id}>
                        <button className="ai-history-item" disabled={busy} onClick={() => openConversation(c.id)}>
                          {c.title}
                        </button>
                        <button
                          className="btn btn-ghost btn-sm"
                          onClick={() => removeConversation(c.id)}
                          title="Xoá" aria-label={`Xoá cuộc trò chuyện ${c.title}`} disabled={busy}
                        >
                          <IconTrash />
                        </button>
                      </div>
                    ))}
                  </>
                )}
              </>
            )}

            {messages.map((m) => (
              <div className={`ai-msg ai-msg-${m.role}`} key={m.id}>{m.role === 'model' ? <RichContent value={m.content} /> : m.content}</div>
            ))}
            {sources.length > 0 && <details className="ai-sources">
              <summary>Nguồn tham khảo đã gửi cho AI ({sources.length})</summary>
              {sources.map((source) => <a key={source.nodeId} href={source.url} target="_blank" rel="noreferrer">
                {source.title}{!source.contentAvailable && <small>Chỉ có mô tả; chưa đọc nội dung tệp/video</small>}
              </a>)}
            </details>}
            {busy && <div className="ai-msg ai-msg-model muted">Đang soạn câu trả lời…</div>}
            {error && <div className="alert alert-error">{error}</div>}
            <div ref={bottomRef} />
          </div>

          <form className="ai-foot" onSubmit={send}>
            <textarea
              ref={inputRef}
              aria-label="Câu hỏi cho trợ lý AI"
              maxLength={4000}
              disabled={busy}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                // Enter gửi, Shift+Enter xuống dòng.
                if (e.key === 'Enter' && !e.shiftKey && !e.nativeEvent.isComposing) {
                  e.preventDefault()
                  void send()
                }
              }}
              placeholder={isStudent ? 'Em muốn hỏi điều gì?' : 'Bạn cần hỗ trợ gì?'}
              rows={2}
            />
            <button className="btn btn-primary btn-sm" disabled={busy || !input.trim()}>Gửi</button>
          </form>
        </div>
      )}
    </>
  )
}
