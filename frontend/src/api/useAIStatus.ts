import { useEffect, useState } from 'react'

import { api } from './client'
import type { AIStatus } from './types'

/**
 * Trạng thái chức năng AI, dùng để ẩn/hiện các nút gọi AI.
 *
 * Cấu hình AI sửa được ở giao diện admin nên không thể quyết định lúc build;
 * mỗi màn hình cần biết đều hỏi máy chủ. Kết quả được nhớ lại trong phiên làm
 * việc để nhiều thành phần cùng dùng chỉ tốn một lần gọi.
 */
let cached: Promise<AIStatus> | null = null

function fetchStatus(): Promise<AIStatus> {
  if (!cached) {
    cached = api.aiStatus().catch((err) => {
      // Lỗi mạng không được "đóng băng" trạng thái tắt cho cả phiên.
      cached = null
      throw err
    })
  }
  return cached
}

/** Xoá bộ nhớ đệm sau khi admin lưu cấu hình mới. */
export function invalidateAIStatus() {
  cached = null
}

/** `null` nghĩa là chưa biết (đang hỏi máy chủ). */
export function useAIStatus(): AIStatus | null {
  const [status, setStatus] = useState<AIStatus | null>(null)

  useEffect(() => {
    let alive = true
    fetchStatus()
      .then((st) => { if (alive) setStatus(st) })
      .catch(() => {
        if (alive) setStatus({ enabled: false, provider: '', model: '', levels: [] })
      })
    return () => { alive = false }
  }, [])

  return status
}
