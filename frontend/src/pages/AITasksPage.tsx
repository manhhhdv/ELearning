import { useEffect, useState } from 'react'

import { api } from '../api/client'
import type { AITask, AITaskKind } from '../api/types'
import { PageHeader } from '../components/Layout'
import { ErrorAlert, Loading, formatDate } from '../components/ui'

const KIND_LABEL: Record<AITaskKind, string> = {
  lesson: 'Tạo bài giảng',
  lesson_plan: 'Soạn giáo án',
  questions: 'Tạo câu hỏi',
  grade_essay: 'Chấm tự luận',
  chat: 'Trợ lý AI',
}

/** Nhật ký "tác vụ AI gần đây" của người đang đăng nhập (mục 4.2). */
export function AITasksPage() {
  const [tasks, setTasks] = useState<AITask[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.listAITasks(100)
      .then(setTasks)
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được nhật ký'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <Loading label="Đang tải nhật ký…" />

  return (
    <>
      <PageHeader
        title="Tác vụ AI gần đây"
        subtitle="Nhật ký các lần hệ thống gọi AI từ tài khoản của bạn."
      />
      <div className="page-body">
        <ErrorAlert message={error} />

        <div className="panel">
        {tasks.length === 0 ? (
          <p className="muted tiny">Bạn chưa dùng chức năng AI nào.</p>
        ) : (
          <table className="ltable">
            <thead>
              <tr><th>Thời điểm</th><th>Loại tác vụ</th><th>Nội dung</th><th>Kết quả</th></tr>
            </thead>
            <tbody>
              {tasks.map((t) => (
                <tr key={t.id}>
                  <td className="tiny muted">{formatDate(t.createdAt)}</td>
                  <td><span className="badge badge-primary">{KIND_LABEL[t.kind] ?? t.kind}</span></td>
                  <td>
                    {t.title || '—'}
                    {t.status === 'error' && t.detail && (
                      <div className="hint" style={{ marginTop: 2 }}>{t.detail}</div>
                    )}
                  </td>
                  <td>
                    {t.status === 'ok' ? (
                      <span className="badge badge-success">
                        Thành công{t.attempts > 1 && ` (sinh lại ${t.attempts - 1} lần)`}
                      </span>
                    ) : (
                      <span className="badge badge-danger">Lỗi</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        </div>
      </div>
    </>
  )
}
