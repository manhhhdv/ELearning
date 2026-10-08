import { useCallback, useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import { api } from '../api/client'
import { useAIStatus } from '../api/useAIStatus'
import type { Question, SubmissionDetail } from '../api/types'
import { PageHeader } from '../components/Layout'
import { IconCheck, IconSparkles } from '../components/icons'
import { Loading, formatDate, formatScore } from '../components/ui'

/** Điểm và nhận xét người chấm đang nhập cho từng câu tự luận. */
type Grades = Record<string, { score: number; comment: string }>

export function SubmissionPage() {
  const { submissionId = '' } = useParams()
  const navigate = useNavigate()
  const [detail, setDetail] = useState<SubmissionDetail | null>(null)
  const [grades, setGrades] = useState<Grades>({})
  // Những câu tự luận người chấm đã chạm tới trong phiên làm việc này.
  const [touched, setTouched] = useState<Set<string>>(new Set())
  const [feedback, setFeedback] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saved, setSaved] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const [aiBusy, setAiBusy] = useState(false)
  // Chỉ hiện nút nhờ AI chấm khi admin đã cấu hình nhà cung cấp AI.
  const aiEnabled = useAIStatus()?.enabled ?? false

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const d = await api.getSubmission(submissionId)
      setDetail(d)
      setFeedback(d.submission.feedback)
      setGrades(Object.fromEntries(
        (d.submission.answers ?? []).map((a) => [a.id, { score: a.score, comment: a.comment }]),
      ))
      setTouched(new Set())
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

  /**
   * Nhờ AI gợi ý điểm cho toàn bộ câu tự luận. Gợi ý được điền sẵn vào form
   * để người chấm sửa; điểm chỉ được ghi nhận khi bấm lưu (mục 3.4.3).
   */
  const askAI = async () => {
    setAiBusy(true)
    setError(null)
    setSaved(null)
    try {
      const { suggestions } = await api.aiGradeSubmission(submissionId)
      setGrades((prev) => {
        const next = { ...prev }
        for (const s of suggestions) {
          next[s.answerId] = { score: s.score, comment: s.comment }
        }
        return next
      })
      // Gợi ý đã điền vào form nên các câu đó tính là đã cho điểm.
      setTouched((prev) => {
        const next = new Set(prev)
        for (const s of suggestions) next.add(s.answerId)
        return next
      })
      // Nạp lại để aiScore / aiComment hiện trong khung gợi ý của từng câu.
      await load()
      setSaved(`AI đã gợi ý điểm cho ${suggestions.length} câu. Hãy kiểm tra rồi bấm lưu để chốt.`)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'AI chưa chấm được bài này')
    } finally {
      setAiBusy(false)
    }
  }

  // Tổng điểm đang hiển thị trên thanh chấm: phần tự động cộng phần đang nhập,
  // cập nhật ngay khi gõ chứ không chờ lưu.
  const gradedNow = submission.autoScore
    + essayAnswers.reduce((sum, a) => sum + (Number(grades[a.id]?.score) || 0), 0)
  // "Đã cho điểm" phải tính cả điểm 0 mà người chấm chủ động đặt, nên không thể
  // so sánh với 0. Bài đã chốt thì mọi câu coi như đã chấm; bài đang chấm dở thì
  // dựa vào những câu người chấm đã thực sự chạm tới trong phiên này.
  const scoredCount = graded
    ? essayAnswers.length
    : essayAnswers.filter((a) => touched.has(a.id)).length

  const markTouched = (answerId: string) =>
    setTouched((prev) => (prev.has(answerId) ? prev : new Set(prev).add(answerId)))

  const setScore = (answerId: string, score: number) => {
    markTouched(answerId)
    setGrades((prev) => ({ ...prev, [answerId]: { score, comment: prev[answerId]?.comment ?? '' } }))
  }

  const setComment = (answerId: string, comment: string) => {
    markTouched(answerId)
    setGrades((prev) => ({ ...prev, [answerId]: { score: prev[answerId]?.score ?? 0, comment } }))
  }

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
            {graded ? 'Đã chấm' : 'Chờ giáo viên chấm'}
          </span>
        </div>

        {submission.feedback && !canGrade && (
          <div className="callout teacher-note">
            <b>Nhận xét của giáo viên</b>
            <div style={{ marginTop: 4 }}>{submission.feedback}</div>
          </div>
        )}

        {!graded && !canGrade && (
          <div className="callout warn">
            Bài có câu tự luận nên đang chờ giáo viên chấm. Đáp án đúng sẽ hiển thị sau khi có điểm.
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
                {/* Câu tự luận chưa chấm hiện "Chưa chấm" thay vì 0 điểm —
                    0/3 trông như học sinh bị điểm liệt chứ không phải chưa chấm. */}
                {isEssay && !graded ? (
                  <span className="pill pill-amber">Chưa chấm · tối đa {formatScore(q.points)} điểm</span>
                ) : (
                  <span className="pill">{formatScore(a.score)} / {formatScore(q.points)} điểm</span>
                )}
              </div>

              <div className="stem">{q.prompt}</div>

              {isEssay && (
                <div className="answer-block">
                  <div className="answer-label">Bài làm của học sinh</div>
                  <div className="essay-answer">
                    {a.essayText || <span style={{ color: 'var(--c-faint)' }}>— Không trả lời —</span>}
                  </div>
                </div>
              )}

              {q.type === 'fill_blank' && (
                <FillBlankReview question={q} answer={a.essayText} correct={a.isCorrect} />
              )}

              {q.type === 'true_false' && (
                <TrueFalseReview question={q} selected={a.selectedOptionIds} />
              )}

              {(q.type === 'single_choice' || q.type === 'multi_choice') && (
                <ChoiceReview question={q} selected={a.selectedOptionIds} />
              )}

              {q.explanation && (
                <div className="callout" style={{ margin: '14px 0 0' }}>
                  <b>Giải thích</b>
                  <div style={{ marginTop: 4 }}>{q.explanation}</div>
                </div>
              )}

              {/* Người chấm cần đối chiếu với đáp án gợi ý và tiêu chí chấm do
                  chính mình soạn. Backend chỉ gửi hai trường này cho người có
                  quyền chấm hoặc giám sát, học sinh không nhận được. */}
              {isEssay && canGrade && (q.sampleAnswer || q.rubric) && (
                <div className="grading-guide">
                  <div className="grading-guide-head">Hướng dẫn chấm</div>
                  {q.sampleAnswer && (
                    <div className="grading-guide-row">
                      <span>Đáp án gợi ý</span>
                      <p>{q.sampleAnswer}</p>
                    </div>
                  )}
                  {q.rubric && (
                    <div className="grading-guide-row">
                      <span>Tiêu chí chấm</span>
                      <p>{q.rubric}</p>
                    </div>
                  )}
                </div>
              )}

              {isEssay && canGrade && a.aiScore !== null && (
                <div className="callout ai-suggestion" style={{ margin: '14px 0 0' }}>
                  <b>AI gợi ý: {formatScore(a.aiScore)} / {formatScore(q.points)} điểm</b>
                  {a.aiComment && <div style={{ marginTop: 4 }}>{a.aiComment}</div>}
                  <button
                    className="cbtn cbtn-plain cbtn-sm"
                    style={{ marginTop: 10 }}
                    onClick={() => {
                      setScore(a.id, a.aiScore ?? 0)
                      setComment(a.id, a.aiComment)
                    }}
                  >
                    Dùng gợi ý này
                  </button>
                  <div className="hint" style={{ marginTop: 8 }}>
                    Đây chỉ là đề xuất. Điểm chính thức là điểm bạn nhập và lưu bên dưới.
                  </div>
                </div>
              )}

              {isEssay && canGrade && (
                <div className="grade-form">
                  <div className="grade-score">
                    <label htmlFor={`score-${a.id}`}>Điểm</label>
                    <div className="grade-score-input">
                      <input
                        id={`score-${a.id}`}
                        type="number" min={0} max={q.points} step={0.5}
                        value={grades[a.id]?.score ?? 0}
                        onChange={(e) => setScore(a.id, Number(e.target.value))}
                      />
                      <span>/ {formatScore(q.points)}</span>
                    </div>
                    {/* Chấm nhanh: phần lớn câu tự luận rơi vào đúng / một nửa / sai,
                        gõ tay từng số cho cả lớp rất mất thời gian. */}
                    <div className="grade-quick">
                      {[
                        { label: 'Tối đa', value: q.points },
                        { label: 'Một nửa', value: Math.round((q.points / 2) * 2) / 2 },
                        { label: '0', value: 0 },
                      ].map((o) => (
                        <button
                          key={o.label}
                          type="button"
                          className={`btn btn-sm ${grades[a.id]?.score === o.value ? 'btn-primary' : ''}`}
                          onClick={() => setScore(a.id, o.value)}
                        >
                          {o.label}
                        </button>
                      ))}
                    </div>
                  </div>

                  <label className="grade-comment">
                    <span>Nhận xét cho câu này</span>
                    <textarea
                      rows={3}
                      value={grades[a.id]?.comment ?? ''}
                      onChange={(e) => setComment(a.id, e.target.value)}
                      placeholder="Chỗ em làm được, chỗ cần bổ sung…"
                    />
                  </label>
                </div>
              )}

              {/* Sau khi chấm xong, học sinh xem được đáp án gợi ý để tự đối
                  chiếu. Tiêu chí chấm là công cụ nội bộ nên backend đã gỡ. */}
              {isEssay && !canGrade && graded && q.sampleAnswer && (
                <div className="callout" style={{ margin: '14px 0 0' }}>
                  <b>Đáp án gợi ý</b>
                  <div style={{ marginTop: 4 }}>{q.sampleAnswer}</div>
                </div>
              )}

              {isEssay && !canGrade && a.comment && (
                <div className="callout teacher-note">
                  <b>Nhận xét của giáo viên</b>
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
              placeholder="Nhận xét gửi tới học sinh…"
            />
            {essayAnswers.length === 0 && (
              <p style={{ fontSize: 13.5, color: 'var(--c-muted)', marginTop: 12, marginBottom: 0 }}>
                Bài này chỉ có câu trắc nghiệm nên đã được chấm tự động.
              </p>
            )}
          </div>
        )}
      </div>

      {/* Thanh chấm bài dính đáy: bài dài vài chục câu thì không phải cuộn
          xuống tận cuối mới bấm lưu được, và luôn thấy còn mấy câu chưa chấm. */}
      {canGrade && (
        <div className="grade-bar">
          <div className="grade-bar-inner">
            <div className="grade-bar-info">
              <b>{formatScore(gradedNow)} / {formatScore(submission.maxScore)} điểm</b>
              <span className="tiny muted">
                {essayAnswers.length === 0
                  ? 'Đã chấm tự động toàn bộ'
                  : `Tự luận: đã cho điểm ${scoredCount}/${essayAnswers.length} câu`}
              </span>
            </div>
            <span className="grow" />
            {aiEnabled && essayAnswers.length > 0 && (
              <button className="cbtn cbtn-plain cbtn-sm" onClick={askAI} disabled={busy || aiBusy}>
                <IconSparkles size={14} /> {aiBusy ? 'AI đang chấm…' : 'Nhờ AI gợi ý điểm'}
              </button>
            )}
            <button className="cbtn cbtn-fill cbtn-sm" onClick={save} disabled={busy || aiBusy}>
              {busy ? 'Đang lưu…' : graded ? 'Cập nhật điểm' : 'Lưu điểm và chốt kết quả'}
            </button>
          </div>
        </div>
      )}
    </>
  )
}

/**
 * Câu điền khuyết: khi chưa được phép xem đáp án, backend trả về danh sách
 * phương án rỗng nên chỉ hiện bài làm của học sinh.
 */
function FillBlankReview({
  question, answer, correct,
}: {
  question: Question
  answer: string
  correct: boolean | null
}) {
  const accepted = question.options.map((o) => o.content)
  return (
    <>
      <div className={`choice review ${correct === true ? 'right' : correct === false ? 'wrong' : 'on'}`}>
        <span className="body">
          {answer || <span style={{ color: 'var(--c-faint)' }}>— Không trả lời —</span>}
        </span>
      </div>
      {accepted.length > 0 && (
        <div className="hint">Đáp án được chấp nhận: {accepted.join(' · ')}</div>
      )}
    </>
  )
}

/**
 * Câu đúng/sai: mỗi phát biểu chấm riêng. Phát biểu học sinh không đánh dấu
 * được hiểu là chọn "Sai".
 */
function TrueFalseReview({ question, selected }: { question: Question; selected: string[] }) {
  const picked = new Set(selected)
  // Backend chỉ trả cờ đáp án đúng khi người xem được phép biết. Trước lúc đó
  // mọi phát biểu đều có isCorrect = false nên không suy ra được đáp án.
  const revealed = question.options.some((o) => o.isCorrect)

  return (
    <div className="tf-list">
      {question.options.map((o) => {
        const chose = picked.has(o.id)
        const right = revealed && chose === o.isCorrect
        let cls = 'choice review'
        if (revealed) cls += right ? ' right' : ' wrong'
        return (
          <div className={cls} key={o.id}>
            <span className="body">{o.content}</span>
            <span className="pill">Bạn chọn: {chose ? 'Đúng' : 'Sai'}</span>
            {revealed && <span className="pill pill-green">Đáp án: {o.isCorrect ? 'Đúng' : 'Sai'}</span>}
          </div>
        )
      })}
    </div>
  )
}

/** Hiển thị các phương án kèm lựa chọn của học sinh và đáp án đúng (nếu được phép xem). */
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
