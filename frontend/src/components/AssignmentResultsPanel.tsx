import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { AssignmentResults, Submission, TreeNode } from '../api/types'
import { QUESTION_TYPE_LABEL } from '../api/types'
import { Loading, formatDate, formatScore } from './ui'

/** Bảng kết quả một bài tập: số liệu tổng, thống kê từng câu và danh sách bài nộp. */
export function AssignmentResultsPanel({ node }: { node: TreeNode }) {
  const [results, setResults] = useState<AssignmentResults | null>(null)
  const [submissions, setSubmissions] = useState<Submission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const [r, subs] = await Promise.all([
        api.assignmentResults(node.id),
        api.listProgramSubmissions(node.programId),
      ])
      setResults(r)
      setSubmissions(subs.filter((s) => s.assignmentId === node.id))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được kết quả')
    } finally {
      setLoading(false)
    }
  }, [node.id, node.programId])

  useEffect(() => { load() }, [load])

  if (loading) return <Loading />
  if (!results) return <div className="alert alert-error">{error ?? 'Không tải được kết quả'}</div>

  const graded = results.submissionCount - results.pendingCount

  return (
    <>
      {error && <div className="alert alert-error">{error}</div>}

      <div className="stat-row" style={{ gap: 14, marginBottom: 24 }}>
        <Metric label="Lượt nộp" value={String(results.submissionCount)} />
        <Metric label="Học viên đã làm" value={String(results.studentCount)} />
        <Metric label="Chờ chấm" value={String(results.pendingCount)} tone={results.pendingCount > 0 ? 'warn' : undefined} />
        <Metric
          label="Điểm trung bình"
          value={results.submissionCount > 0
            ? `${formatScore(results.averageScore)} / ${formatScore(results.maxScore)}`
            : '—'}
        />
        {results.passScore > 0 && (
          <Metric
            label={`Đạt (từ ${formatScore(results.passScore)} điểm)`}
            value={graded > 0 ? `${results.passedCount}/${graded}` : '—'}
          />
        )}
      </div>

      <h3 style={{ marginBottom: 10 }}>Kết quả theo từng câu hỏi</h3>
      <div className="card results-table" style={{ marginBottom: 26 }}>
        <table>
          <thead>
            <tr>
              <th style={{ width: 76 }}>Mã</th>
              <th>Câu hỏi</th>
              <th style={{ width: 150 }}>Loại</th>
              <th style={{ width: 90 }}>Trả lời</th>
              <th style={{ width: 110 }}>Tỉ lệ đúng</th>
              <th style={{ width: 110 }}>Điểm TB</th>
            </tr>
          </thead>
          <tbody>
            {results.questions.map((q) => {
              const rate = q.answerCount > 0 ? Math.round((q.correctCount / q.answerCount) * 100) : null
              return (
                <tr key={q.questionId}>
                  <td className="mono tiny"><b>{q.code}</b></td>
                  <td>
                    {q.prompt.length > 62 ? q.prompt.slice(0, 62) + '…' : q.prompt}
                    {q.blankCount > 0 && (
                      <div className="tiny muted">{q.blankCount} lượt bỏ trống</div>
                    )}
                  </td>
                  <td className="tiny muted">{QUESTION_TYPE_LABEL[q.type]}</td>
                  <td>{q.answerCount}</td>
                  <td>
                    {q.type === 'essay' ? (
                      q.needsGrading > 0
                        ? <span className="badge badge-warning">{q.needsGrading} chờ chấm</span>
                        : <span className="tiny muted">chấm tay</span>
                    ) : rate === null ? (
                      <span className="tiny muted">—</span>
                    ) : (
                      <>
                        <div className="progress-bar" style={{ marginBottom: 3 }}>
                          <i style={{
                            width: `${rate}%`,
                            background: rate >= 70 ? 'var(--success)' : rate >= 40 ? 'var(--warning)' : 'var(--danger)',
                          }} />
                        </div>
                        <span className="tiny muted">{q.correctCount}/{q.answerCount} · {rate}%</span>
                      </>
                    )}
                  </td>
                  <td>{formatScore(q.averageScore)} / {formatScore(q.points)}</td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>

      <h3 style={{ marginBottom: 10 }}>Bài nộp của học viên</h3>
      <div className="card">
        {submissions.length === 0 ? (
          <div className="empty" style={{ padding: '30px 20px' }}>
            <h3>Chưa có học viên nào nộp bài</h3>
          </div>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Học viên</th><th style={{ width: 80 }}>Lượt</th>
                <th style={{ width: 170 }}>Nộp lúc</th><th style={{ width: 130 }}>Điểm</th>
                <th style={{ width: 120 }}>Trạng thái</th><th />
              </tr>
            </thead>
            <tbody>
              {submissions.map((s) => (
                <tr key={s.id}>
                  <td>
                    <b>{s.studentName || '—'}</b>
                    <div className="tiny muted">{s.studentEmail}</div>
                  </td>
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
                    <Link className="btn btn-sm" to={`/bai-nop/${s.id}`}>
                      {s.needsGrading ? 'Chấm bài' : 'Xem'}
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </>
  )
}

function Metric({ label, value, tone }: { label: string; value: string; tone?: 'warn' }) {
  return (
    <div className="card card-pad" style={{ padding: '16px 18px' }}>
      <div className="tiny muted">{label}</div>
      <div style={{
        fontSize: 22, fontWeight: 700, lineHeight: 1.3,
        color: tone === 'warn' ? 'var(--warning)' : undefined,
      }}>
        {value}
      </div>
    </div>
  )
}
