import { useCallback, useRef, useState } from 'react'

import { api } from '../api/client'
import type { AICourseSource, AIConversation, AIMessage } from '../api/types'

export interface AIChatScope {
  subject?: string
  grade?: string
  programId?: string
  nodeId?: string
}

/**
 * Phần lõi của trợ lý AI: lịch sử, tin nhắn và việc gửi câu hỏi.
 *
 * Dùng chung cho bảng trò chuyện nổi và trang "Hỏi trợ lý AI" nên chỉ giữ dữ
 * liệu, không dính gì tới cách bày trí.
 */
export function useAIChat({ subject, grade, programId, nodeId }: AIChatScope) {
  const [history, setHistory] = useState<AIConversation[]>([])
  const [conversationId, setConversationId] = useState<string | undefined>()
  const [messages, setMessages] = useState<AIMessage[]>([])
  const [input, setInput] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [sources, setSources] = useState<AICourseSource[]>([])

  // Mỗi lần đổi cuộc trò chuyện là một "phiên"; câu trả lời về muộn của phiên cũ bị bỏ.
  const requestVersion = useRef(0)

  // Lịch sử chỉ hiện các cuộc trò chuyện cùng ngữ cảnh bài học đang mở.
  const scopedHistory = useCallback(
    (items: AIConversation[]) => items.filter((c) =>
      (c.programId ?? undefined) === programId && (c.nodeId ?? undefined) === nodeId),
    [programId, nodeId],
  )

  const loadHistory = useCallback(() => {
    api.listConversations()
      .then((items) => setHistory(scopedHistory(items)))
      .catch(() => setHistory([]))
  }, [scopedHistory])

  const startNew = useCallback(() => {
    requestVersion.current++
    setSources([])
    setConversationId(undefined)
    setMessages([])
    setError(null)
  }, [])

  const openConversation = useCallback(async (id: string) => {
    const version = ++requestVersion.current
    setBusy(true)
    setError(null)
    try {
      const conv = await api.getConversation(id)
      if (version !== requestVersion.current) return
      if ((conv.programId ?? undefined) !== programId || (conv.nodeId ?? undefined) !== nodeId) return
      setSources([])
      setConversationId(conv.id)
      setMessages(conv.messages ?? [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được cuộc trò chuyện')
    } finally {
      if (version === requestVersion.current) setBusy(false)
    }
  }, [programId, nodeId])

  const removeConversation = useCallback(async (id: string) => {
    if (!confirm('Xoá cuộc trò chuyện này?')) return
    try {
      await api.deleteConversation(id)
      setHistory((prev) => prev.filter((c) => c.id !== id))
      if (conversationId === id) startNew()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không xoá được')
    }
  }, [conversationId, startNew])

  const send = useCallback(async (e?: { preventDefault?: () => void }) => {
    e?.preventDefault?.()
    const text = input.trim()
    if (!text || busy) return

    const version = ++requestVersion.current
    // Hiện câu hỏi ngay để khung chat không bị "đứng" trong lúc chờ AI.
    const pending: AIMessage = {
      id: `tmp-${Date.now()}`, role: 'user', content: text, createdAt: new Date().toISOString(),
    }
    setMessages((prev) => [...prev, pending])
    setInput('')
    setBusy(true)
    setError(null)

    try {
      const res = await api.aiChat({ conversationId, message: text, subject, grade, programId, nodeId })
      if (version !== requestVersion.current) return
      setSources(res.sources ?? [])
      setConversationId(res.conversationId)
      setMessages((prev) => [...prev, {
        id: `a-${Date.now()}`, role: 'model', content: res.answer, createdAt: new Date().toISOString(),
      }])
      api.listConversations().then((items) => setHistory(scopedHistory(items))).catch(() => undefined)
    } catch (err) {
      if (version !== requestVersion.current) return
      setError(err instanceof Error ? err.message : 'Trợ lý AI chưa trả lời được')
      // Bỏ câu hỏi vừa thêm để người dùng gửi lại mà không bị trùng.
      setMessages((prev) => prev.filter((m) => m.id !== pending.id))
      setInput(text)
    } finally {
      if (version === requestVersion.current) setBusy(false)
    }
  }, [busy, conversationId, grade, input, nodeId, programId, scopedHistory, subject])

  const cancelPending = useCallback(() => { requestVersion.current++ }, [])

  return {
    history, conversationId, messages, input, setInput, busy, error, sources,
    loadHistory, startNew, openConversation, removeConversation, send, cancelPending,
  }
}

/** Gợi ý câu hỏi mồi, khác nhau theo vai trò và theo việc có đang mở bài học hay không. */
export function chatSuggestions(isStudent: boolean, programId?: string): string[] {
  if (!isStudent) {
    return ['Gợi ý hoạt động mở bài', 'Thiết kế 3 câu hỏi vận dụng', 'Cách giải thích khái niệm này cho dễ hiểu']
  }
  return [
    programId ? 'Tóm tắt các ý chính trong bài đang học' : 'Giải thích giúp em khái niệm này',
    'Em nên bắt đầu ôn từ đâu?',
    'Cho em một ví dụ dễ hiểu',
  ]
}
