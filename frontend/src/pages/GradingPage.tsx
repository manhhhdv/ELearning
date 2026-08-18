import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { Program, Submission } from '../api/types'
import { PageHeader } from '../components/Layout'
import { EmptyState, ErrorAlert, Loading, formatDate, formatScore } from '../components/ui'

export function GradingPage() {
  const [programs, setPrograms] = useState<Program[]>([])
  const [programId, setProgramId] = useState('')
  const [onlyPending, setOnlyPending] = useState(true)
  const [items, setItems] = useState<Submission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api.listPrograms()
      .then((list) => {
        setPrograms(list)
        setProgramId((cur) => cur || list[0]?.id || '')
        if (list.length === 0) setLoading(false)
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : 'Không tải được chương trình')
        setLoading(false)
      })
  }, [])

  const load = useCallback(async () => {
    if (!programId) return
    setLoading(true)
    try {
      setItems(await api.listProgramSubmissions(programId, onlyPending))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được bài nộp')
    } finally {
      setLoading(false)
    }
  }, [programId, onlyPending])

  useEffect(() => { load() }, [load])

  return (
    <>
      <PageHeader title="Chấm bài" subtitle="Bài nộp của học viên trong các chương trình bạn phụ trách" />
      <div className="page-body">
        <div className="toolbar">
          <select
            value={programId}
            onChange={(e) => setProgramId(e.target.value)}
            className="grow"
            style={{ maxWidth: 380 }}
          >
            {programs.length === 0 && <option value="">Chưa có chương trình nào</option>}
            {programs.map((p) => <option key={p.id} value={p.id}>{p.title}</option>)}
          </select>
          <label className="checkbox">
            <input type="checkbox" checked={onlyPending} onChange={(e) => setOnlyPending(e.target.checked)} />
            Chỉ hiện bài chờ chấm
          </label>
        </div>

        <ErrorAlert message={error} />

        {loading ? <Loading /> : (
          <div className="card">
            {items.length === 0 ? (
              <EmptyState title={onlyPending ? 'Không còn bài nào chờ chấm' : 'Chưa có bài nộp nào'} />
            ) : (
              <table>
                <thead>
                  <tr><th>Học viên</th><th>Bài tập</th><th>Lượt</th><th>Nộp lúc</th><th>Điểm</th><th>Trạng thái</th><th></th></tr>
                </thead>
                <tbody>
                  {items.map((s) => (
                    <tr key={s.id}>
                      <td>
                        <b>{s.studentName || '—'}</b>
                        <div className="tiny muted">{s.studentEmail}</div>
                      </td>
                      <td>{s.assignmentTitle}</td>
                      <td>#{s.attemptNo}</td>
                      <td className="tiny muted">{formatDate(s.submittedAt)}</td>
                      <td>
                        <b>{formatScore(s.autoScore + (s.manualScore ?? 0))}</b>
                        <span className="muted"> / {formatScore(s.maxScore)}</span>
                      </td>
                      <td>
                        <span className={`badge ${s.status === 'graded' ? 'badge-success' : 'badge-warning'}`}>
                          {s.status === 'graded' ? 'Đã chấm' : 'Chờ chấm'}
                        </span>
                      </td>
                      <td style={{ textAlign: 'right' }}>
                        <Link className="btn btn-sm btn-primary" to={`/bai-nop/${s.id}`}>
                          {s.needsGrading ? 'Chấm bài' : 'Xem'}
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}
      </div>
    </>
  )
}
