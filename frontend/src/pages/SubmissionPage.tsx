import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Question, SubmissionDetail } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconCheck } from '../components/icons'
import { Loading, formatDate, formatScore } from '../components/ui'

/** Điểm và nhận xét người chấm đang nhập cho từng câu tự luận. */
type Grades = Record<string, { score: number; comment: string }>

export function SubmissionPage() {
  const { submissionId = '' } = useParams()
  const navigate = useNavigate()
  const [detail, setDetail] = useState<SubmissionDetail | null>(null)
  const [grades, setGrades] = useState<Grades>({})
  const [feedback, setFeedback] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const d = await api.getSubmission(submissionId)
      setDetail(d)
      setFeedback(d.submission.feedback)
      setGrades(Object.fromEntries(
        (d.submission.answers ?? []).map((a) => [a.id, { score: a.score, comment: a.comment }]),
      ))
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được bài nộp')
    } finally {
      setLoading(false)
    }
  }, [submissionId])

  useEffect(() => { load() }, [load])

  if (loading) return <Loading />
  if (!detail) {
    return <div className="page-body learner"><div className="callout warn">{error ?? 'Không tìm thấy bài nộp'}</div></div>
  }

  const { submission, questions, canGrade } = detail
  const byQuestion = new Map(questions.map((q) => [q.id, q]))
  const answers = submission.answers ?? []
  const essayAnswers = answers.filter((a) => byQuestion.get(a.questionId)?.type === 'essay')
  const graded = submission.status === 'graded'
  const total = submission.autoScore + (submission.manualScore ?? 0)
  const percent = submission.maxScore > 0 ? Math.round((total / submission.maxScore) * 100) : 0
  const correctCount = answers.filter((a) => a.isCorrect === true).length
  const gradedCount = answers.filter((a) => byQuestion.get(a.questionId)?.type !== 'essay').length

  const save = async () => {
    setBusy(true)
    setError(null)
    setSaved(null)
    try {
      await api.gradeSubmission(
        submissionId,
        feedback,
        essayAnswers.map((a) => ({
          answerId: a.id,
          score: Number(grades[a.id]?.score) || 0,
          comment: grades[a.id]?.comment ?? '',
        })),
      )
      setSaved('Đã lưu điểm và chốt kết quả bài nộp')
      await load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không lưu được điểm')
    } finally {
      setBusy(false)
    }
  }

  return (
    <>
      <PageHeader
        title={submission.assignmentTitle ?? 'Bài nộp'}
        subtitle={`${submission.studentName || submission.studentEmail} · lượt #${submission.attemptNo} · nộp ${formatDate(submission.submittedAt)}`}
        actions={<button className="btn" onClick={() => navigate(-1)}>← Quay lại</button>}
      />

      <div className="page-body learner">
        {error && <div className="callout warn">{error}</div>}
        {saved && <div className="callout ok">{saved}</div>}

        <div className="result-head">
          <div
            className="score-ring"
            style={{
              // Bài chưa chấm xong chỉ hiển thị phần điểm đã có, tô màu chờ.
              '--pct': graded ? percent : 0,
              '--ring': graded ? 'var(--c-green)' : 'var(--c-amber)',
            } as React.CSSProperties}
          >
            <div>
              {graded ? (
                <><b>{percent}%</b><span>điểm số</span></>
              ) : (
                <><b style={{ fontSize: 15 }}>Chờ</b><span>chấm bài</span></>
              )}
            </div>
          </div>

          <div style={{ flex: 1, minWidth: 220 }}>
            <div style={{ fontSize: 13, color: 'var(--c-muted)' }}>Tổng điểm</div>
            <div style={{ fontSize: 30, fontWeight: 700, lineHeight: 1.2 }}>
              {formatScore(total)}
              <span style={{ color: 'var(--c-muted)', fontSize: 20 }}> / {formatScore(submission.maxScore)}</span>
            </div>
            <div className="quiz-stats" style={{ marginTop: 12, gap: 22 }}>
              <div>
                <span>Trắc nghiệm</span>
                <b style={{ fontSize: 15 }}>{correctCount}/{gradedCount} câu đúng · {formatScore(submission.autoScore)} điểm</b>
              </div>
              {submission.manualScore !== null && (
                <div>
                  <span>Tự luận</span>
                  <b style={{ fontSize: 15 }}>{formatScore(submission.manualScore)} điểm</b>
                </div>
              )}
            </div>
          </div>

          <span className={`pill ${graded ? 'pill-green' : 'pill-amber'}`}>
            {graded ? 'Đã chấm' : 'Chờ giảng viên chấm'}
          </span>
        </div>

        {submission.feedback && !canGrade && (
          <div className="callout">
            <b>Nhận xét của giảng viên</b>
            <div style={{ marginTop: 4 }}>{submission.feedback}</div>
          </div>
        )}

        {!graded && !canGrade && (
          <div className="callout warn">
            Bài có câu tự luận nên đang chờ giảng viên chấm. Đáp án đúng sẽ hiển thị sau khi có điểm.
          </div>
        )}

        {answers.map((a, i) => {
          const q = byQuestion.get(a.questionId)
          if (!q) return null
          const isEssay = q.type === 'essay'

          return (
            <article className="qcard" key={a.id}>
              <div className="qcard-top">
                <span className="no">Câu {i + 1}</span>
                {a.isCorrect === true && <span className="pill pill-green"><IconCheck size={12} /> Đúng</span>}
                {a.isCorrect === false && <span className="pill pill-red">Sai</span>}
                <span className="grow" />
                <span className="pill">{formatScore(a.score)} / {formatScore(q.points)} điểm</span>
              </div>

              <div className="stem">{q.prompt}</div>

              {isEssay ? (
                <div className="essay-answer">
                  {a.essayText || <span style={{ color: 'var(--c-faint)' }}>— Không trả lời —</span>}
                </div>
              ) : (
                <ChoiceReview question={q} selected={a.selectedOptionIds} />
              )}

              {q.explanation && (
                <div className="callout" style={{ margin: '14px 0 0' }}>
                  <b>Giải thích</b>
                  <div style={{ marginTop: 4 }}>{q.explanation}</div>
                </div>
              )}

              {isEssay && canGrade && (
                <div style={{ display: 'flex', gap: 14, marginTop: 16, flexWrap: 'wrap' }}>
                  <label style={{ width: 130 }}>
                    <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 5 }}>
                      Điểm (tối đa {formatScore(q.points)})
                    </div>
                    <input
                      type="number" min={0} max={q.points} step={0.5}
                      value={grades[a.id]?.score ?? 0}
                      onChange={(e) => setGrades((prev) => ({
                        ...prev,
                        [a.id]: { score: Number(e.target.value), comment: prev[a.id]?.comment ?? '' },
                      }))}
                    />
                  </label>
                  <label style={{ flex: 1, minWidth: 200 }}>
                    <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 5 }}>Nhận xét cho câu này</div>
                    <input
                      type="text"
                      value={grades[a.id]?.comment ?? ''}
                      onChange={(e) => setGrades((prev) => ({
                        ...prev,
                        [a.id]: { score: prev[a.id]?.score ?? 0, comment: e.target.value },
                      }))}
                    />
                  </label>
                </div>
              )}

              {isEssay && !canGrade && a.comment && (
                <div className="callout" style={{ margin: '14px 0 0' }}>
                  <b>Nhận xét</b>
                  <div style={{ marginTop: 4 }}>{a.comment}</div>
                </div>
              )}
            </article>
          )
        })}

        {canGrade && (
          <div className="qcard">
            <h3 style={{ fontSize: 16, marginBottom: 10 }}>Nhận xét chung</h3>
            <textarea
              className="essay-box"
              style={{ minHeight: 110 }}
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              placeholder="Nhận xét gửi tới học viên…"
            />
            <div style={{ marginTop: 14 }}>
              <button className="cbtn cbtn-fill" onClick={save} disabled={busy}>
                {busy ? 'Đang lưu…' : graded ? 'Cập nhật điểm' : 'Lưu điểm và chốt kết quả'}
              </button>
              {essayAnswers.length === 0 && (
                <p style={{ fontSize: 13, color: 'var(--c-muted)', marginTop: 8 }}>
                  Bài này chỉ có câu trắc nghiệm nên đã được chấm tự động.
                </p>
              )}
            </div>
          </div>
        )}
      </div>
    </>
  )
}

/** Hiển thị các phương án kèm lựa chọn của học viên và đáp án đúng (nếu được phép xem). */
function ChoiceReview({ question, selected }: { question: Question; selected: string[] }) {
  const picked = new Set(selected)
  // Backend chỉ trả cờ đáp án đúng khi người xem được phép biết.
  const revealed = question.options.some((o) => o.isCorrect)

  return (
    <>
      {question.options.map((o) => {
        const chosen = picked.has(o.id)
        let cls = 'choice review'
        if (revealed && o.isCorrect) cls += ' right'
        else if (chosen && revealed) cls += ' wrong'
        else if (chosen) cls += ' on'

        return (
          <div className={cls} key={o.id}>
            <span style={{ width: 17, textAlign: 'center', marginTop: 2 }}>{chosen ? '●' : '○'}</span>
            <span className="body">{o.content}</span>
            {revealed && o.isCorrect && <span className="pill pill-green">đáp án đúng</span>}
            {revealed && chosen && !o.isCorrect && <span className="pill pill-red">bạn đã chọn</span>}
          </div>
        )
      })}
    </>
  )
}
