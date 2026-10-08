import { useEffect, useRef } from 'react'
import { Navigate } from 'react-router-dom'

import { useAIStatus } from '../api/useAIStatus'
import { useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { RichContent } from '../components/RichContent'
import { chatSuggestions, useAIChat } from '../components/useAIChat'
import { IconPlus, IconTrash } from '../components/icons'
import { Loading, formatDate } from '../components/ui'

/**
 * Trang "Hỏi trợ lý AI" — bản toàn màn hình của bảng trò chuyện nổi.
 *
 * Ở đây không gắn ngữ cảnh bài học nào: người dùng hỏi chung, còn khi đang học
 * một bài thì bảng nổi trong trình học vẫn gửi kèm mã lớp và mã bài.
 */
export function AIChatPage() {
  const { user } = useAuth()
  const aiStatus = useAIStatus()
  const {
    history, conversationId, messages, input, setInput, busy, error, sources,
    loadHistory, startNew, openConversation, removeConversation, send, cancelPending,
  } = useAIChat({})
  const inputRef = useRef<HTMLTextAreaElement>(null)
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    loadHistory()
    inputRef.current?.focus()
  }, [loadHistory])

  useEffect(() => cancelPending, [cancelPending])

  // Luôn cuộn tới lượt trả lời mới nhất.
  useEffect(() => {
    if (messages.length === 0) return
    bottomRef.current?.scrollIntoView({
      behavior: window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'instant' : 'smooth',
      block: 'end',
    })
  }, [messages, busy])

  if (!user) return null
  if (aiStatus === null) return <Loading label="Đang kiểm tra cấu hình AI…" />
  // Máy chủ tắt AI thì trang này không có gì để hiển thị.
  if (!aiStatus.enabled) return <Navigate to="/tong-quan" replace />

  const isStudent = user.role === 'student'
  const suggestions = chatSuggestions(isStudent)

  const askNow = (text: string) => {
    setInput(text)
    inputRef.current?.focus()
  }

  return (
    <>
      <PageHeader
        title={isStudent ? 'Trợ lý học tập AI' : 'Trợ lý AI cho giáo viên'}
        subtitle={isStudent
          ? 'Hỏi bài, xin giải thích từng bước. Đây là công cụ hỗ trợ, khi cần chốt kiến thức em hãy hỏi thêm thầy cô.'
          : 'Ý tưởng bài dạy, hoạt động trên lớp và câu hỏi kiểm tra.'}
        actions={
          <button className="btn btn-sm" onClick={startNew} disabled={busy || (!conversationId && messages.length === 0)}>
            <IconPlus /> Cuộc trò chuyện mới
          </button>
        }
      />

      <div className="page-body">
        <div className="ai-page">
          <section className="ai-page-main card">
            <div className="ai-page-log" role="log" aria-live="polite" aria-busy={busy}>
              {messages.length === 0 && (
                <div className="ai-page-empty">
                  <p className="hint" style={{ marginTop: 0 }}>
                    {isStudent
                      ? 'Trợ lý sẽ gợi ý và giải thích từng bước thay vì đưa ngay đáp án.'
                      : 'Trợ lý hỗ trợ tìm ý tưởng bài học, thiết kế hoạt động và gợi ý câu hỏi kiểm tra.'}
                  </p>
                  <div className="wrap-gap">
                    {suggestions.map((sg) => (
                      <button key={sg} type="button" className="btn btn-sm" onClick={() => askNow(sg)}>{sg}</button>
                    ))}
                  </div>
                </div>
              )}

              {messages.map((m) => (
                <div className={`ai-msg ai-msg-${m.role}`} key={m.id}>
                  {m.role === 'model' ? <RichContent value={m.content} /> : m.content}
                </div>
              ))}

              {sources.length > 0 && (
                <details className="ai-sources">
                  <summary>Nguồn tham khảo đã gửi cho AI ({sources.length})</summary>
                  {sources.map((source) => (
                    <a key={source.nodeId} href={source.url} target="_blank" rel="noreferrer">
                      {source.title}
                      {!source.contentAvailable && <small>Chỉ có mô tả; chưa đọc nội dung tệp/video</small>}
                    </a>
                  ))}
                </details>
              )}
              {busy && <div className="ai-msg ai-msg-model muted">Đang soạn câu trả lời…</div>}
              {error && <div className="alert alert-error">{error}</div>}
              <div ref={bottomRef} />
            </div>

            <form className="ai-page-foot" onSubmit={send}>
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
                rows={3}
              />
              <button className="btn btn-primary" disabled={busy || !input.trim()}>Gửi</button>
            </form>
          </section>

          <aside className="ai-page-side card">
            <h3 className="ai-page-side-title">Cuộc trò chuyện gần đây</h3>
            {history.length === 0 ? (
              <p className="muted tiny" style={{ margin: 0 }}>Chưa có cuộc trò chuyện nào.</p>
            ) : (
              history.map((c) => (
                <div className={`ai-history-row${c.id === conversationId ? ' active' : ''}`} key={c.id}>
                  <button className="ai-history-item" disabled={busy} onClick={() => openConversation(c.id)}>
                    <span className="ai-history-title">{c.title}</span>
                    <small className="muted">{formatDate(c.updatedAt)}</small>
                  </button>
                  <button
                    className="btn btn-ghost btn-sm"
                    onClick={() => removeConversation(c.id)}
                    title="Xoá" aria-label={`Xoá cuộc trò chuyện ${c.title}`} disabled={busy}
                  >
                    <IconTrash />
                  </button>
                </div>
              ))
            )}
          </aside>
        </div>
      </div>
    </>
  )
}
