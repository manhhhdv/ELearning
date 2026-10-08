import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { GradebookAssignment, GradebookCell, GradebookRow, ProgramGradebook } from '../api/types'
import { EmptyState, ErrorAlert, Loading, formatDate, formatScore } from './ui'
import { IconDownload } from './icons'

type SortKey = 'name' | 'score' | 'done'

/**
 * Bảng điểm bài tập của cả chương trình: mỗi dòng một học sinh, mỗi cột một bài tập.
 * Ô điểm lấy **lượt cao điểm nhất** — đúng với con số học sinh nhìn thấy ở màn hình bài tập.
 */
export function GradebookPanel({ programId, programCode }: { programId: string; programCode: string }) {
  const [book, setBook] = useState<ProgramGradebook | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [sort, setSort] = useState<SortKey>('name')

  const load = useCallback(async () => {
    setLoading(true)
    try {
      setBook(await api.programGradebook(programId))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không tải được bảng điểm')
    } finally {
      setLoading(false)
    }
  }, [programId])

  useEffect(() => { load() }, [load])

  const rows = useMemo(() => {
    if (!book) return []
    const needle = search.trim().toLowerCase()
    const filtered = needle
      ? book.rows.filter((r) =>
          r.fullName.toLowerCase().includes(needle) || r.email.toLowerCase().includes(needle))
      : book.rows
    const sorted = [...filtered]
    if (sort === 'score') sorted.sort((a, b) => b.totalScore - a.totalScore)
    if (sort === 'done') sorted.sort((a, b) => b.doneCount - a.doneCount)
    return sorted
  }, [book, search, sort])

  if (loading) return <Loading />
  if (!book) return <ErrorAlert message={error ?? 'Không tải được bảng điểm'} />

  if (book.assignments.length === 0) {
    return (
      <div className="card">
        <EmptyState title="Lớp học chưa có bài tập nào">
          <p>Thêm bài tập ở tab “Nội dung” rồi quay lại đây để xem điểm của học sinh.</p>
        </EmptyState>
      </div>
    )
  }

  const submitRate = book.totalCells > 0 ? Math.round((book.submittedCells / book.totalCells) * 100) : 0

  return (
    <>
      <ErrorAlert message={error} />

      <div
        className="stat-row"
        style={{ gap: 14, marginBottom: 24, gridTemplateColumns: 'repeat(auto-fit, minmax(196px, 1fr))' }}
      >
        <Metric label="Học sinh" value={String(book.rows.length)} />
        <Metric label="Bài tập" value={String(book.assignments.length)} />
        <Metric
          label="Tỉ lệ đã nộp"
          value={`${submitRate}%`}
          hint={`${book.submittedCells}/${book.totalCells} lượt bài`}
        />
        <Metric
          label="Điểm tổng trung bình"
          value={book.submittedCells > 0
            ? `${formatScore(book.averageScore)} / ${formatScore(book.totalMaxScore)}`
            : '—'}
          hint="Tính trên học sinh đã nộp ít nhất một bài"
        />
        <Metric
          label="Chờ chấm"
          value={String(book.pendingCount)}
          tone={book.pendingCount > 0 ? 'warn' : undefined}
        />
      </div>

      <div className="toolbar">
        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Tìm học sinh theo tên hoặc email…"
          style={{ maxWidth: 300 }}
        />
        <select value={sort} onChange={(e) => setSort(e.target.value as SortKey)} style={{ maxWidth: 210 }}>
          <option value="name">Sắp xếp: theo tên</option>
          <option value="score">Sắp xếp: điểm tổng cao nhất</option>
          <option value="done">Sắp xếp: nộp nhiều bài nhất</option>
        </select>
        <button className="btn" onClick={() => downloadCsv(book, programCode)}>
          <IconDownload size={16} /> Xuất CSV
        </button>
      </div>

      <div className="card gradebook">
        <table>
          <thead>
            <tr>
              <th className="gb-name">Học sinh</th>
              {book.assignments.map((a) => (
                <th key={a.nodeId} className="gb-col">
                  <div className="gb-col-title" title={a.title}>{a.title}</div>
                  <div className="tiny faint">
                    thang {formatScore(a.maxScore)}
                    {a.passScore > 0 && ` · đạt từ ${formatScore(a.passScore)}`}
                  </div>
                </th>
              ))}
              <th className="gb-total">Điểm tổng</th>
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td colSpan={book.assignments.length + 2} className="muted tiny center" style={{ padding: 26 }}>
                  {book.rows.length === 0
                    ? 'Chưa có học sinh nào được ghi danh vào lớp học.'
                    : 'Không có học sinh nào khớp từ khoá tìm kiếm.'}
                </td>
              </tr>
            ) : rows.map((row) => (
              <tr key={row.userId}>
                <td className="gb-name">
                  <b>{row.fullName || '—'}</b>
                  {!row.enrolled && <span className="badge" style={{ marginLeft: 6 }}>đã gỡ ghi danh</span>}
                  <div className="tiny muted">{row.email}</div>
                </td>
                {book.assignments.map((a) => (
                  <td key={a.nodeId} className="gb-col">
                    <ScoreCell cell={row.cells[a.nodeId]} assignment={a} />
                  </td>
                ))}
                <td className="gb-total">
                  <b>{formatScore(row.totalScore)}</b>
                  <span className="muted"> / {formatScore(book.totalMaxScore)}</span>
                  <div className="tiny muted">
                    đã nộp {row.doneCount}/{book.assignments.length}
                    {row.pendingCount > 0 && ` · ${row.pendingCount} chờ chấm`}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
          <tfoot>
            <tr>
              <td className="gb-name"><b>Trung bình lớp</b></td>
              {book.assignments.map((a) => (
                <td key={a.nodeId} className="gb-col">
                  {a.submittedCount === 0 ? (
                    <span className="muted">—</span>
                  ) : (
                    <>
                      <b>{formatScore(a.averageScore)}</b>
                      <span className="muted"> / {formatScore(a.maxScore)}</span>
                      <div className="tiny muted">
                        {a.submittedCount}/{book.rows.length} đã nộp
                        {a.passScore > 0 && ` · ${a.passedCount} đạt`}
                      </div>
                    </>
                  )}
                </td>
              ))}
              <td className="gb-total">
                {book.submittedCells > 0 ? (
                  <b>{formatScore(book.averageScore)}</b>
                ) : (
                  <span className="muted">—</span>
                )}
              </td>
            </tr>
          </tfoot>
        </table>
      </div>
    </>
  )
}

/** Một ô điểm: điểm cao nhất của học sinh ở bài tập đó, bấm vào mở bài nộp. */
function ScoreCell({ cell, assignment }: { cell?: GradebookCell; assignment: GradebookAssignment }) {
  if (!cell) return <span className="muted">—</span>

  const scale = cell.maxScore > 0 ? cell.maxScore : assignment.maxScore
  const passed = assignment.passScore > 0
    ? cell.score >= assignment.passScore
    : scale > 0 && cell.score / scale >= 0.5

  return (
    <Link className="gb-cell" to={`/bai-nop/${cell.submissionId}`} title={`Nộp lúc ${formatDate(cell.submittedAt)}`}>
      <b className={cell.pending > 0 ? 'gb-pending' : passed ? 'gb-pass' : 'gb-fail'}>
        {formatScore(cell.score)}
      </b>
      <span className="muted"> / {formatScore(scale)}</span>
      <div className="tiny muted">
        {cell.pending > 0 ? 'chờ chấm' : `${cell.attempts} lượt`}
      </div>
    </Link>
  )
}

function Metric({ label, value, hint, tone }: { label: string; value: string; hint?: string; tone?: 'warn' }) {
  return (
    <div className="card card-pad" style={{ padding: '16px 18px' }}>
      <div className="tiny muted">{label}</div>
      <div style={{
        fontSize: 22, fontWeight: 700, lineHeight: 1.3,
        color: tone === 'warn' ? 'var(--warning)' : undefined,
      }}>
        {value}
      </div>
      {hint && <div className="tiny faint">{hint}</div>}
    </div>
  )
}

/** Xuất bảng điểm ra CSV mở được bằng Excel (có BOM để không vỡ dấu tiếng Việt). */
function downloadCsv(book: ProgramGradebook, programCode: string) {
  const esc = (v: string | number) => `"${String(v).replace(/"/g, '""')}"`
  const header = ['Họ và tên', 'Email', ...book.assignments.map((a) => a.title), 'Điểm tổng', 'Thang điểm tổng']
  const lines = [header.map(esc).join(',')]

  for (const row of book.rows) {
    lines.push([
      esc(row.fullName),
      esc(row.email),
      ...book.assignments.map((a) => esc(cellText(row, a))),
      esc(round2(row.totalScore)),
      esc(round2(book.totalMaxScore)),
    ].join(','))
  }

  const blob = new Blob(['﻿' + lines.join('\r\n')], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `diem-bai-tap-${programCode.toLowerCase() || 'chuong-trinh'}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

function cellText(row: GradebookRow, assignment: GradebookAssignment) {
  const cell = row.cells[assignment.nodeId]
  if (!cell) return ''
  return cell.pending > 0 ? `${round2(cell.score)} (chờ chấm)` : round2(cell.score)
}

const round2 = (v: number) => Number(v.toFixed(2))
