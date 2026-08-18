import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { Submission } from '../api/types'
import { PageHeader } from '../components/Layout'
import { Loading, formatDate, formatScore } from '../components/ui'

export function MyResultsPage() {
  const [items, setItems] = useState<Submission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.mySubmissions()
      .then(setItems)
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được kết quả'))
      .finally(() => setLoading(false))
  }, [])

  const graded = items.filter((s) => s.status === 'graded')
  const pending = items.length - graded.length

  return (
    <>
      <PageHeader title="Kết quả bài tập" subtitle="Toàn bộ bài bạn đã nộp" />
      <div className="page-body learner">
        {error && <div className="callout warn">{error}</div>}

        {loading ? <Loading /> : items.length === 0 ? (
          <div className="panel">
            <div className="blank">
              <h3>Bạn chưa nộp bài tập nào</h3>
              <p>Kết quả sẽ hiện ở đây sau khi bạn hoàn thành bài tập đầu tiên.</p>
            </div>
          </div>
        ) : (
          <>
            <div className="quiz-head" style={{ marginTop: 0 }}>
              <div className="quiz-stats">
                <div><span>Tổng số bài đã nộp</span><b>{items.length}</b></div>
                <div><span>Đã có điểm</span><b>{graded.length}</b></div>
                <div><span>Đang chờ chấm</span><b>{pending}</b></div>
              </div>
            </div>

            <div className="panel">
              <table className="ltable">
                <thead>
                  <tr>
                    <th>Bài tập</th><th>Khoá học</th><th>Lượt</th>
                    <th>Nộp lúc</th><th>Điểm</th><th>Trạng thái</th><th />
                  </tr>
                </thead>
                <tbody>
                  {items.map((s) => {
                    const total = s.autoScore + (s.manualScore ?? 0)
                    const percent = s.maxScore > 0 ? Math.round((total / s.maxScore) * 100) : 0
                    return (
                      <tr key={s.id}>
                        <td><b>{s.assignmentTitle}</b></td>
                        <td style={{ color: 'var(--c-muted)', fontSize: 13.5 }}>{s.programTitle}</td>
                        <td>#{s.attemptNo}</td>
                        <td style={{ color: 'var(--c-muted)', fontSize: 13.5 }}>{formatDate(s.submittedAt)}</td>
                        <td>
                          <b>{formatScore(total)}</b>
                          <span style={{ color: 'var(--c-muted)' }}> / {formatScore(s.maxScore)}</span>
                          {s.status === 'graded' && (
                            <div style={{ fontSize: 12.5, color: 'var(--c-muted)' }}>{percent}%</div>
                          )}
                        </td>
                        <td>
                          <span className={`pill ${s.status === 'graded' ? 'pill-green' : 'pill-amber'}`}>
                            {s.status === 'graded' ? 'Đã chấm' : 'Chờ chấm'}
                          </span>
                        </td>
                        <td style={{ textAlign: 'right' }}>
                          <Link className="cbtn cbtn-plain cbtn-sm" to={`/bai-nop/${s.id}`}>Xem lại</Link>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          </>
        )}
      </div>
    </>
  )
}
