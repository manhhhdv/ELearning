import { useCallback, useEffect, useRef, useState } from 'react'
import { Link } from 'react-router-dom'

import { api } from '../api/client'
import type { AttemptView, Question, TreeNode } from '../api/types'
import { CountdownTimer } from './CountdownTimer'
import { IconCheck, IconChevron } from './icons'
import { Loading, formatDate, formatScore } from './ui'

/** Câu trả lời đang soạn: danh sách ID phương án cho trắc nghiệm, chuỗi cho tự luận. */
type Draft = Record<string, { options: string[]; essay: string }>

interface Props {
  node: TreeNode
  programTitle: string
  onSubmitted: () => void
  /** Điều hướng lớp học, chỉ hiện khi học sinh chưa vào màn hình làm bài. */
  prevNode: TreeNode | null
  nextNode: TreeNode | null
  onNavigate: (node: TreeNode) => void
}

export function AssignmentPlayer({
  node, programTitle, onSubmitted, prevNode, nextNode, onNavigate,
}: Props) {
  const [view, setView] = useState<AttemptView | null>(null)
  const [draft, setDraft] = useState<Draft>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [taking, setTaking] = useState(false)
  // Chặn nộp trùng khi đồng hồ về 0 đúng lúc học sinh bấm nút.
  const submittingRef = useRef(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await api.getAttempt(node.id)
      setView(data)
      setDraft((prev) => Object.fromEntries(
        (data.node.assignment?.questions ?? []).map(
          (q) => [q.id, prev[q.id] ?? { options: [], essay: '' }]),
      ))
      // Học sinh tải lại trang giữa chừng: quay về đúng màn hình làm bài, đồng hồ vẫn chạy tiếp.
      setTaking(data.session !== null)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được bài tập')
    } finally {
      setLoading(false)
    }
  }, [node.id])

  useEffect(() => { load() }, [load])

  const send = useCallback(async () => {
    setBusy(true)
    setError(null)
    try {
      const items = view?.node.assignment?.questions ?? []
      await api.submitAssignment(node.id, items.map((q) => ({
        questionId: q.id,
        selectedOptionIds: draft[q.id]?.options ?? [],
        essayText: draft[q.id]?.essay ?? '',
      })))
      setTaking(false)
      await load()
      onSubmitted()
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không nộp được bài')
      // Máy chủ từ chối vì hết giờ: rời màn hình làm bài và nạp lại trạng thái thật.
      setTaking(false)
      await load()
    } finally {
      setBusy(false)
    }
  }, [node.id, view, draft, load, onSubmitted])

  const submit = async () => {
    const missing = questions.length - answeredCount
    if (missing > 0 && !confirm(`Còn ${missing} câu chưa trả lời. Vẫn nộp bài?`)) return
    await send()
  }

  // Hết giờ thì nộp luôn phần đã làm, không hỏi lại.
  const handleExpire = useCallback(() => {
    if (submittingRef.current) return
    submittingRef.current = true
    void send().finally(() => { submittingRef.current = false })
  }, [send])

  if (loading) return <Loading label="Đang tải bài tập…" />
  if (!view) return <div className="callout warn">{error ?? 'Không mở được bài tập'}</div>

  const assignment = view.node.assignment
  const questions = assignment?.questions ?? []
  const totalPoints = questions.reduce((sum, q) => sum + q.points, 0)
  const attemptsLeft = view.maxAttempts > 0 ? view.maxAttempts - view.attemptsUsed : null
  const outOfAttempts = attemptsLeft !== null && attemptsLeft <= 0
  const best = view.submissions.reduce<number | null>(
    (max, s) => Math.max(max ?? 0, s.autoScore + (s.manualScore ?? 0)), null)
  // Bài còn câu tự luận chưa chấm chỉ mới có điểm phần trắc nghiệm. Hiển thị
  // con số đó như "điểm cao nhất" khiến học sinh tưởng mình đã mất điểm, nên
  // phải nói rõ đây là điểm tạm tính.
  const hasPending = view.submissions.some((s) => s.status !== 'graded')

  const begin = async () => {
    setBusy(true)
    setError(null)
    try {
      const data = await api.startAttempt(node.id)
      setView(data)
      setDraft(Object.fromEntries(
        (data.node.assignment?.questions ?? []).map((q) => [q.id, { options: [], essay: '' }]),
      ))
      setTaking(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được lượt làm bài')
    } finally {
      setBusy(false)
    }
  }

  /**
   * Ghi lựa chọn của học sinh.
   * Với câu đúng/sai, `options` chứa ID của các phát biểu được đánh dấu "Đúng";
   * phát biểu không nằm trong danh sách được hiểu là "Sai".
   */
  const pick = (q: Question, optionId: string, checked: boolean) => {
    setDraft((prev) => {
      const cur = prev[q.id] ?? { options: [], essay: '' }
      const options = q.type === 'single_choice'
        ? [optionId]
        : checked
          ? [...cur.options.filter((id) => id !== optionId), optionId]
          : cur.options.filter((id) => id !== optionId)
      return { ...prev, [q.id]: { ...cur, options } }
    })
  }

  const setEssay = (questionId: string, essay: string) => {
    setDraft((prev) => ({
      ...prev,
      [questionId]: { options: prev[questionId]?.options ?? [], essay },
    }))
  }

  const answeredCount = questions.filter((q) => {
    const a = draft[q.id]
    // Câu đúng/sai luôn coi như đã trả lời: không đánh dấu nghĩa là chọn "Sai" cho mọi phát biểu.
    if (q.type === 'true_false') return true
    if (q.type === 'essay' || q.type === 'fill_blank') return !!a?.essay.trim()
    return !!a?.options.length
  }).length


  return (
    <>
      <div className="crumb">{programTitle}</div>
      <h1>{node.title}</h1>
      <div className="stage-meta">
        <span className="pill pill-blue">Bài tập</span>
        <span>{questions.length} câu hỏi</span>
        <span className="dot">·</span>
        <span>{formatScore(totalPoints)} điểm</span>
        {assignment && assignment.timeLimitMinutes > 0 && (
          <><span className="dot">·</span><span>{assignment.timeLimitMinutes} phút</span></>
        )}
      </div>

      {error && <div className="callout warn">{error}</div>}

      {!taking && (
        <>
          <div className="quiz-head">
            <div className="quiz-stats">
              <div>
                <span>{hasPending ? 'Điểm tạm tính' : 'Điểm cao nhất'}</span>
                <b>{best === null ? '—' : `${formatScore(best)} / ${formatScore(totalPoints)}`}</b>
                {hasPending && (
                  <small className="stat-note">Chưa gồm điểm tự luận đang chờ chấm</small>
                )}
              </div>
              <div>
                <span>Lượt làm bài</span>
                <b>{view.attemptsUsed}{view.maxAttempts > 0 ? ` / ${view.maxAttempts}` : ''}</b>
              </div>
              {assignment && assignment.passScore > 0 && (
                <div>
                  <span>Điểm đạt</span>
                  <b>{formatScore(assignment.passScore)}</b>
                </div>
              )}
              {assignment?.dueAt && (
                <div>
                  <span>Hạn nộp</span>
                  <b style={{ fontSize: 14.5 }}>{formatDate(assignment.dueAt)}</b>
                </div>
              )}
            </div>

            <div style={{ marginTop: 18 }}>
              {questions.length === 0 ? (
                <span style={{ color: 'var(--c-muted)' }}>Bài tập này chưa có câu hỏi.</span>
              ) : outOfAttempts ? (
                <span style={{ color: 'var(--c-muted)' }}>
                  Bạn đã dùng hết {view.maxAttempts} lượt làm bài cho phép.
                </span>
              ) : (
                <button className="cbtn cbtn-fill" onClick={begin} disabled={busy}>
                  {busy ? 'Đang mở đề…' : view.attemptsUsed > 0 ? 'Làm lại bài' : 'Bắt đầu làm bài'}
                </button>
              )}
            </div>
          </div>

          {assignment?.instructions && (
            <div className="callout">{assignment.instructions}</div>
          )}

          {view.submissions.length > 0 && (
            <>
              <h2 style={{ fontSize: 17, marginBottom: 12 }}>Lịch sử làm bài</h2>
              <div className="panel">
                <table className="ltable">
                  <thead>
                    <tr><th>Lượt</th><th>Nộp lúc</th><th>Điểm</th><th>Trạng thái</th><th /></tr>
                  </thead>
                  <tbody>
                    {view.submissions.map((s) => (
                      <tr key={s.id}>
                        <td>#{s.attemptNo}</td>
                        <td style={{ color: 'var(--c-muted)', fontSize: 13.5 }}>{formatDate(s.submittedAt)}</td>
                        <td>
                          <b>{formatScore(s.autoScore + (s.manualScore ?? 0))}</b>
                          <span style={{ color: 'var(--c-muted)' }}> / {formatScore(s.maxScore)}</span>
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
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}

          <div className="stage-nav">
            <button
              className="cbtn cbtn-plain cbtn-sm"
              disabled={!prevNode}
              onClick={() => prevNode && onNavigate(prevNode)}
            >
              <span style={{ transform: 'rotate(180deg)', display: 'inline-flex' }}><IconChevron size={13} /></span>
              Bài trước
            </button>
            <div className="spacer" />
            {nextNode && (
              <button className="cbtn cbtn-fill cbtn-sm" onClick={() => onNavigate(nextNode)}>
                Bài tiếp theo <IconChevron size={13} />
              </button>
            )}
          </div>
        </>
      )}

      {taking && (
        <>
          <div className="quiz-bar">
            <b style={{ fontSize: 15 }}>Đang làm bài</b>
            <span className="pill">{questions.length} câu · {formatScore(totalPoints)} điểm</span>
            {view.session?.expiresAt && (
              <CountdownTimer expiresAt={view.session.expiresAt} onExpire={handleExpire} />
            )}
          </div>

          {view.session?.expiresAt && (
            <div className="callout warn">
              Bài có giới hạn thời gian. Hết giờ hệ thống sẽ tự nộp phần bạn đã làm.
            </div>
          )}

          {assignment?.instructions && <div className="callout">{assignment.instructions}</div>}

          {questions.map((q, i) => (
            <article className="qcard" key={q.id}>
              <div className="qcard-top">
                <span className="no">Câu {i + 1}</span>
                <span className="grow" />
                <span className="pill">{formatScore(q.points)} điểm</span>
                {q.type === 'multi_choice' && <span className="pill pill-blue">Chọn nhiều đáp án</span>}
                {q.type === 'true_false' && <span className="pill pill-blue">Đúng / Sai</span>}
                {q.type === 'fill_blank' && <span className="pill pill-blue">Điền khuyết</span>}
                {q.type === 'essay' && <span className="pill pill-blue">Tự luận</span>}
              </div>

              <div className="stem">{q.prompt}</div>

              {/* Tự luận và điền khuyết đều nhập bằng bàn phím, khác nhau ở kích thước ô nhập. */}
              {q.type === 'essay' && (
                <textarea
                  className="essay-box"
                  value={draft[q.id]?.essay ?? ''}
                  onChange={(e) => setEssay(q.id, e.target.value)}
                  placeholder="Nhập câu trả lời của bạn…"
                />
              )}

              {q.type === 'fill_blank' && (
                <input
                  className="blank-box"
                  type="text"
                  value={draft[q.id]?.essay ?? ''}
                  onChange={(e) => setEssay(q.id, e.target.value)}
                  placeholder="Điền vào chỗ trống…"
                />
              )}

              {q.type === 'true_false' && (
                <div className="tf-list">
                  {q.options.map((o) => {
                    const on = draft[q.id]?.options.includes(o.id) ?? false
                    return (
                      <div className="tf-row" key={o.id}>
                        <span className="tf-text">{o.content}</span>
                        <div className="tf-pick">
                          {/* Không chọn gì cũng được chấm: khi đó phát biểu tính là "Sai". */}
                          <label className={`tf-opt ${on ? 'on' : ''}`}>
                            <input
                              type="radio"
                              name={`tf-${o.id}`}
                              checked={on}
                              onChange={() => pick(q, o.id, true)}
                            />
                            Đúng
                          </label>
                          <label className={`tf-opt ${!on ? 'on' : ''}`}>
                            <input
                              type="radio"
                              name={`tf-${o.id}`}
                              checked={!on}
                              onChange={() => pick(q, o.id, false)}
                            />
                            Sai
                          </label>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}

              {(q.type === 'single_choice' || q.type === 'multi_choice') && (
                q.options.map((o) => {
                  const on = draft[q.id]?.options.includes(o.id) ?? false
                  return (
                    <label className={`choice ${on ? 'on' : ''}`} key={o.id}>
                      <input
                        type={q.type === 'single_choice' ? 'radio' : 'checkbox'}
                        name={`q-${q.id}`}
                        checked={on}
                        onChange={(e) => pick(q, o.id, e.target.checked)}
                      />
                      <span className="body">{o.content}</span>
                    </label>
                  )
                })
              )}
            </article>
          ))}

          <div className="stage-nav">
            <span style={{ color: 'var(--c-muted)', fontSize: 14 }}>
              Đã trả lời {answeredCount}/{questions.length} câu
            </span>
            <div className="spacer" />
            <button className="cbtn cbtn-plain cbtn-sm" onClick={() => setTaking(false)} disabled={busy}>
              Thoát
            </button>
            <button className="cbtn cbtn-fill cbtn-sm" onClick={submit} disabled={busy}>
              <IconCheck size={14} /> {busy ? 'Đang nộp…' : 'Nộp bài'}
            </button>
          </div>
        </>
      )}
    </>
  )
}
