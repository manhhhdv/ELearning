import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { DashboardStats } from '../api/types'
import { ROLE_LABEL, STATUS_LABEL } from '../api/types'
import { useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { EmptyState, ErrorAlert, Loading, formatDate } from '../components/ui'

export function AdminDashboardPage() {
  const { user } = useAuth()
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.dashboard()
      .then(setStats)
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được số liệu'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <>
      <PageHeader
        title="Bảng điều khiển"
        subtitle={user?.role === 'supervisor'
          ? 'Số liệu tổng quan toàn hệ thống (chế độ Giám sát — chỉ xem)'
          : 'Số liệu tổng quan toàn hệ thống'}
      />

      <div className="page-body">
        <ErrorAlert message={error} />

        {loading ? <Loading /> : !stats ? null : (
          <>
            <div className="stat-row" style={{ marginBottom: 26 }}>
              <StatCard label="Chương trình" value={stats.programsTotal}
                detail={`${stats.programsPublished} xuất bản · ${stats.programsDraft} nháp · ${stats.programsArchived} lưu trữ`} />
              <StatCard label="Người dùng" value={stats.usersTotal}
                detail={`${stats.studentCount} học viên · ${stats.trainerCount} giảng viên`} />
              <StatCard label="Lượt ghi danh" value={stats.enrollmentsTotal} />
              <StatCard label="Bài nộp" value={stats.submissionsTotal}
                detail={`${stats.submissionsPending} đang chờ chấm`}
                tone={stats.submissionsPending > 0 ? 'warning' : undefined} />
            </div>

            <div className="row" style={{ alignItems: 'flex-start', gap: 20 }}>
              <div style={{ flex: 1, minWidth: 0 }}>
                <h3 style={{ marginBottom: 10 }}>Người dùng theo vai trò</h3>
                <div className="card card-pad">
                  <RoleBar label={ROLE_LABEL.admin} count={stats.adminCount} total={stats.usersTotal} />
                  <RoleBar label={ROLE_LABEL.trainer} count={stats.trainerCount} total={stats.usersTotal} />
                  <RoleBar label={ROLE_LABEL.supervisor} count={stats.supervisorCount} total={stats.usersTotal} />
                  <RoleBar label={ROLE_LABEL.student} count={stats.studentCount} total={stats.usersTotal} last />
                </div>

                <h3 style={{ margin: '22px 0 10px' }}>Người dùng mới nhất</h3>
                <div className="card">
                  {stats.recentSignups.length === 0 ? (
                    <EmptyState title="Chưa có ai" />
                  ) : (
                    <table>
                      <thead><tr><th>Họ tên</th><th>Email</th><th>Vai trò</th><th>Ngày tạo</th></tr></thead>
                      <tbody>
                        {stats.recentSignups.map((u) => (
                          <tr key={u.id}>
                            <td><b>{u.fullName || '—'}</b></td>
                            <td className="muted tiny">{u.email}</td>
                            <td><span className="badge badge-primary">{ROLE_LABEL[u.role]}</span></td>
                            <td className="tiny muted">{formatDate(u.createdAt)}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                </div>
              </div>

              <div style={{ flex: 1, minWidth: 0 }}>
                <h3 style={{ marginBottom: 10 }}>Chương trình nhiều học viên nhất</h3>
                <div className="card">
                  {stats.topPrograms.length === 0 ? (
                    <EmptyState title="Chưa có chương trình nào" />
                  ) : (
                    <table>
                      <thead><tr><th>Chương trình</th><th>Trạng thái</th><th>Học viên</th></tr></thead>
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
              </div>
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

function RoleBar({
  label, count, total, last = false,
}: {
  label: string
  count: number
  total: number
  last?: boolean
}) {
  const percent = total > 0 ? Math.round((count / total) * 100) : 0
  return (
    <div style={{ marginBottom: last ? 0 : 12 }}>
      <div className="spread tiny" style={{ marginBottom: 4 }}>
        <span>{label}</span>
        <span className="muted">{count}</span>
      </div>
      <div className="progress-bar"><i style={{ width: `${percent}%` }} /></div>
    </div>
  )
}
