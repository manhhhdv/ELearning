import { useState } from 'react'

import { api } from '../api/client'
import { IconDownload } from './icons'

/**
 * Hai nút tải bộ câu hỏi của một bài tập về dạng Word (.docx): đề cho học sinh,
 * và đề kèm trang đáp án, lời giải, hướng dẫn chấm. Tệp lấy theo dữ liệu đã lưu.
 */
export function ExportQuestionsButtons({ nodeId, onError }: {
  nodeId: string
  onError: (message: string | null) => void
}) {
  const [busy, setBusy] = useState<'plain' | 'answers' | null>(null)

  const run = async (withAnswers: boolean) => {
    setBusy(withAnswers ? 'answers' : 'plain')
    onError(null)
    try {
      await api.exportNode(nodeId, withAnswers)
    } catch (err) {
      onError(err instanceof Error ? err.message : 'Không xuất được tệp Word')
    } finally {
      setBusy(null)
    }
  }

  return (
    <>
      <button
        type="button" className="btn btn-sm" onClick={() => run(false)} disabled={busy !== null}
        title="Tải đề dạng Word (.docx), không kèm đáp án"
      >
        <IconDownload /> {busy === 'plain' ? 'Đang xuất…' : 'Xuất đề Word'}
      </button>
      <button
        type="button" className="btn btn-sm" onClick={() => run(true)} disabled={busy !== null}
        title="Tải đề kèm trang đáp án, lời giải và hướng dẫn chấm (.docx)"
      >
        <IconDownload /> {busy === 'answers' ? 'Đang xuất…' : 'Đề + đáp án'}
      </button>
    </>
  )
}
