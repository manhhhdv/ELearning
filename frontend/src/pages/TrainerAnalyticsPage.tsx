import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { DashboardStats } from '../api/types'
import { STATUS_LABEL } from '../api/types'
import { PageHeader } from '../components/Layout'
import { EmptyState, ErrorAlert, Loading } from '../components/ui'

/**
 * Phân tích học tập dành cho giáo viên: cùng dữ liệu như Bảng điều khiển của
 * admin nhưng chỉ tính trong phạm vi các lớp học họ phụ trách — không lộ số
 * liệu của lớp học khác.
 */
export function TrainerAnalyticsPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.trainerAnalytics()
      .then(setStats)
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được số liệu'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <>
      <PageHeader
        title="Phân tích học tập"
        subtitle="Số liệu trong phạm vi các lớp học bạn phụ trách"
      />

      <div className="page-body">
        <ErrorAlert message={error} />

        {loading ? <Loading /> : !stats ? null : (
          <>
            <div className="stat-row" style={{ marginBottom: 26 }}>
              <StatCard label="Lớp học" value={stats.programsTotal}
                detail={`${stats.programsPublished} xuất bản · ${stats.programsDraft} nháp · ${stats.programsArchived} lưu trữ`} />
              <StatCard label="Học sinh" value={stats.enrollmentsTotal} />
              <StatCard label="Bài nộp" value={stats.submissionsTotal}
                detail={`${stats.submissionsPending} đang chờ chấm`}
                tone={stats.submissionsPending > 0 ? 'warning' : undefined} />
              <StatCard label="Lượt hoàn thành bài học" value={stats.lessonsCompleted} />
            </div>

            <h3 style={{ marginBottom: 10 }}>Lớp học nhiều học sinh nhất</h3>
            <div className="card">
              {stats.topPrograms.length === 0 ? (
                <EmptyState title="Bạn chưa phụ trách lớp học nào" />
              ) : (
                <table>
                  <thead><tr><th>Lớp học</th><th>Trạng thái</th><th>Học sinh</th></tr></thead>
                  <tbody>
                    {stats.topPrograms.map((p) => (
                      <tr key={p.id}>
                        <td>
                          <Link to={`/quan-tri/chuong-trinh/${p.slug}`}><b>{p.title}</b></Link>
                          <div className="tiny muted mono">{p.code}</div>
                        </td>
                        <td>
                          <span className={`badge ${p.status === 'published' ? 'badge-success' : ''}`}>
                            {STATUS_LABEL[p.status]}
                          </span>
                        </td>
                        <td>{p.enrollmentCount}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          </>
        )}
      </div>
    </>
  )
}

function StatCard({
  label, value, detail, tone,
}: {
  label: string
  value: number
  detail?: string
  tone?: 'warning'
}) {
  return (
    <div className="stat">
      <div>
        <div className="label">{label}</div>
        <div className="value" style={tone === 'warning' ? { color: 'var(--warning)' } : undefined}>
          {value.toLocaleString('vi-VN')}
        </div>
        {detail && <div className="tiny muted" style={{ marginTop: 2 }}>{detail}</div>}
      </div>
    </div>
  )
}
