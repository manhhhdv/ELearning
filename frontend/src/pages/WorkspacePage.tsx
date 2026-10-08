import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import { api } from '../api/client'
import { useAIStatus } from '../api/useAIStatus'
import type { AITask, Workspace } from '../api/types'
import { canManageContent, useAuth } from '../auth'
import { PageHeader } from '../components/Layout'
import { RichContent } from '../components/RichContent'
import { IconBook, IconCheck, IconDoc, IconPlus, IconSparkles } from '../components/icons'
import { ErrorAlert, Loading, formatDate, formatScore } from '../components/ui'

/** Gợi ý nhanh cho ô nhập yêu cầu AI, khác nhau theo vai trò. */
const STUDENT_PROMPTS = [
  'Giải thích giúp em khái niệm này',
  'Lập kế hoạch ôn tập cho tuần này',
  'Cho em một ví dụ dễ hiểu',
  'Em nên bắt đầu ôn từ đâu?',
]

const TEACHER_PROMPTS = [
  'Gợi ý hoạt động mở bài',
  'Thiết kế 3 câu hỏi vận dụng',
  'Cách giải thích khái niệm này cho dễ hiểu',
  'Đề xuất tiêu chí chấm bài tự luận',
]

const TASK_LABEL: Record<AITask['kind'], string> = {
  lesson: 'Tạo bài giảng',
  lesson_plan: 'Soạn giáo án',
  questions: 'Tạo câu hỏi',
  grade_essay: 'Chấm tự luận',
  chat: 'Trợ lý AI',
}

/**
 * Trang tổng quan — "không gian làm việc AI" (mục 4.2).
 *
 * Cùng một bố cục cho mọi vai trò: ô nhập yêu cầu AI kèm gợi ý nhanh, các lối
 * tắt tới chức năng chính, danh sách tác vụ AI gần đây và lịch học. Người dạy
 * thấy thêm bài chờ chấm và cảnh báo học tập.
 */
export function WorkspacePage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const ai = useAIStatus()

  const [data, setData] = useState<Workspace | null>(null)
  const [tasks, setTasks] = useState<AITask[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [prompt, setPrompt] = useState('')
  const [answer, setAnswer] = useState<string | null>(null)
  const [asking, setAsking] = useState(false)

  const manage = user ? canManageContent(user) : false

  useEffect(() => {
    Promise.all([
      api.workspace().then(setData),
      api.listAITasks(6).then(setTasks).catch(() => setTasks([])),
    ])
      .catch((err) => setError(err instanceof Error ? err.message : 'Không tải được trang tổng quan'))
      .finally(() => setLoading(false))
  }, [])

  /** Hỏi nhanh trợ lý AI ngay tại trang tổng quan. */
  const ask = async (e: React.FormEvent) => {
    e.preventDefault()
    const text = prompt.trim()
    if (!text || asking) return
    setAsking(true)
    setAnswer(null)
    setError(null)
    try {
      const res = await api.aiChat({ message: text })
      setAnswer(res.answer)
      setPrompt('')
      api.listAITasks(6).then(setTasks).catch(() => undefined)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Trợ lý AI chưa trả lời được')
    } finally {
      setAsking(false)
    }
  }

  if (loading) return <Loading label="Đang tải trang tổng quan…" />

  const suggestions = manage ? TEACHER_PROMPTS : STUDENT_PROMPTS
  const overdue = data?.schedule.filter((s) => s.overdue) ?? []
  const upcoming = data?.schedule.filter((s) => !s.overdue) ?? []

  return (
    <>
      <PageHeader
        title={`Xin chào, ${user?.fullName ?? ''}`}
        subtitle={manage
          ? 'Không gian làm việc: soạn nội dung với AI, theo dõi lớp và chấm bài'
          : 'Không gian học tập: hỏi trợ lý AI, xem bài cần làm và tiến độ của bạn'}
      />

      <div className="page-body">
        <ErrorAlert message={error} />

        {/* --- Ô nhập yêu cầu cho AI --- */}
        {ai?.enabled && (
          <div className="panel ai-workspace">
            <form onSubmit={ask}>
              <div className="field" style={{ marginBottom: 10 }}>
                <label htmlFor="ai-ask">
                  <IconSparkles /> Hỏi trợ lý AI
                </label>
                <textarea
                  id="ai-ask"
                  rows={2}
                  value={prompt}
                  onChange={(e) => setPrompt(e.target.value)}
                  onKeyDown={(e) => {
                    // Enter gửi, Shift+Enter xuống dòng.
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault()
                      void ask(e)
                    }
                  }}
                  placeholder={manage ? 'Bạn cần hỗ trợ gì cho bài dạy?' : 'Em muốn hỏi điều gì?'}
                />
              </div>
              <div className="wrap-gap">
                {suggestions.map((sg) => (
                  <button key={sg} type="button" className="btn btn-sm" onClick={() => setPrompt(sg)}>
                    {sg}
                  </button>
                ))}
                <span className="grow" />
                <button className="btn btn-primary btn-sm" disabled={asking || !prompt.trim()}>
                  {asking ? 'Đang soạn câu trả lời…' : 'Gửi'}
                </button>
              </div>
            </form>

            {answer && (
              <div className="ai-msg ai-msg-model" style={{ marginTop: 14, marginRight: 0 }}>
                <RichContent value={answer} />
              </div>
            )}
          </div>
        )}

        {/* --- Lối tắt tới các chức năng chính --- */}
        <div className="shortcut-grid">
          {manage ? (
            <>
              <ShortcutCard
                color="violet"
                icon={<IconPlus />} title="Tạo đề bằng AI"
                desc="Sinh câu hỏi 4 dạng từ chủ đề hoặc tài liệu PDF"
                onClick={() => navigate('/quan-tri/chuong-trinh')}
              />
              <ShortcutCard
                color="blue"
                icon={<IconSparkles />} title="Soạn giáo án bằng AI"
                desc="Bản nháp giáo án theo từng hoạt động dạy học"
                onClick={() => navigate('/quan-tri/giao-an')}
              />
              <ShortcutCard
                color="emerald"
                icon={<IconCheck />} title="Chấm bài"
                desc={data?.pendingGrading ? `${data.pendingGrading} bài đang chờ chấm` : 'Không có bài nào chờ chấm'}
                badge={data?.pendingGrading || undefined}
                onClick={() => navigate('/quan-tri/cham-bai')}
              />
              <ShortcutCard
                color="amber"
                icon={<IconDoc />} title="Kho tài liệu"
                desc="Tài liệu dùng chung của trường"
                onClick={() => navigate('/tai-lieu')}
              />
            </>
          ) : (
            <>
              <ShortcutCard
                color="violet"
                icon={<IconBook />} title="Lớp học của tôi"
                desc="Tiếp tục bài đang học"
                onClick={() => navigate('/hoc')}
              />
              <ShortcutCard
                color="blue"
                icon={<IconPlus />} title="Khám phá lớp học"
                desc="Tìm và ghi danh lớp học mới"
                onClick={() => navigate('/kham-pha')}
              />
              <ShortcutCard
                color="emerald"
                icon={<IconCheck />} title="Kết quả học tập"
                desc="Điểm, đáp án và nhận xét từng bài"
                onClick={() => navigate('/ket-qua')}
              />
              <ShortcutCard
                color="amber"
                icon={<IconDoc />} title="Kho tài liệu"
                desc="Tài liệu dùng chung của trường"
                onClick={() => navigate('/tai-lieu')}
              />
            </>
          )}
        </div>

        <div className="workspace-cols">
          <div className="workspace-main">
          {manage ? (
            <>
              {/* --- Chấm bài đang chờ + cảnh báo học tập, thay cho các khối chỉ dành cho học sinh --- */}
              <section className="panel">
                <div className="spread" style={{ marginBottom: 12 }}>
                  <h3 style={{ margin: 0 }}>Bài đang chờ chấm</h3>
                  {!!data?.pendingGrading && (
                    <span className="badge badge-danger">{data.pendingGrading} bài</span>
                  )}
                </div>
                {!data?.pendingGrading ? (
                  <EmptyLine
                    icon={<IconCheck size={18} />}
                    title="Không có bài nào chờ chấm"
                    hint="Bài học sinh nộp sẽ hiện ở đây để bạn chấm điểm."
                  />
                ) : (
                  <EmptyLine
                    icon={<IconCheck size={18} />}
                    title={`${data.pendingGrading} bài đang chờ bạn chấm`}
                    hint="Vào trang Chấm bài để xem và chấm điểm chi tiết."
                  />
                )}
                {!!data?.pendingGrading && (
                  <button
                    type="button"
                    className="btn btn-sm"
                    style={{ marginTop: 12 }}
                    onClick={() => navigate('/quan-tri/cham-bai')}
                  >
                    Đi tới trang Chấm bài
                  </button>
                )}
              </section>

              <section className="panel">
                <h3 style={{ marginTop: 0 }}>Cảnh báo học tập</h3>
                {!data || data.alerts.length === 0 ? (
                  <p className="muted tiny">Chưa có học sinh nào cần chú ý.</p>
                ) : (
                  <>
                    <table className="ltable">
                      <thead>
                        <tr><th>Học sinh</th><th>Lớp học</th><th>Vấn đề</th></tr>
                      </thead>
                      <tbody>
                        {data.alerts.map((a) => (
                          <tr key={`${a.userId}-${a.programId}-${a.kind}`}>
                            <td>
                              {a.fullName}
                              <div className="tiny muted">{a.email}</div>
                            </td>
                            <td className="tiny">{a.programName}</td>
                            <td>
                              {a.kind === 'low_score' ? (
                                <span className="badge badge-danger">
                                  Điểm trung bình {a.value}%
                                </span>
                              ) : (
                                <span className="badge badge-warning">
                                  {a.value} bài quá hạn chưa nộp
                                </span>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    <div className="hint" style={{ marginTop: 10 }}>
                      Cảnh báo dựa trên quy tắc cố định: điểm trung bình dưới 50% hoặc có bài quá hạn
                      chưa nộp. Việc dùng AI phân tích sâu dữ liệu học tập là hướng phát triển tiếp theo.
                    </div>
                  </>
                )}
              </section>
            </>
          ) : (
            <>
              {/* --- Lịch học / việc cần làm --- */}
              <section className="panel">
                <div className="spread" style={{ marginBottom: 12 }}>
                  <h3 style={{ margin: 0 }}>Bài cần làm</h3>
                  {overdue.length > 0 && (
                    <span className="badge badge-danger">{overdue.length} bài quá hạn</span>
                  )}
                </div>
                {data && data.schedule.length === 0 ? (
                  <EmptyLine
                    icon={<IconCheck size={18} />}
                    title="Không có bài nào sắp tới hạn"
                    hint="Bài tập giáo viên giao sẽ hiện ở đây kèm hạn nộp."
                  />
                ) : (
                  <>
                    <ul className="task-list">
                      {[...overdue, ...upcoming].slice(0, 8).map((it) => (
                        <li key={it.nodeId} className={it.overdue ? 'overdue' : ''}>
                          <Link to={`/hoc/${it.programSlug}/${it.nodeSlug}`}>{it.title}</Link>
                          <div className="tiny muted">
                            {it.programName}
                            {it.dueAt && <> · hạn {formatDate(it.dueAt)}</>}
                            {it.overdue && <> · <b>quá hạn</b></>}
                          </div>
                        </li>
                      ))}
                    </ul>
                  </>
                )}
              </section>

              {/* --- Tiến độ học tập (mục 3.5) --- */}
              {data && data.progress.length > 0 && (
                <section className="panel">
                  <h3 style={{ marginTop: 0 }}>Tiến độ học tập</h3>
                  <div className="progress-list">
                    {data.progress.map((p) => {
                      const pct = p.total > 0 ? Math.round((p.completed / p.total) * 100) : 0
                      return (
                        <div className="progress-row" key={p.programId}>
                          <div className="progress-row-head">
                            <Link to={`/hoc/${p.programSlug}`}>{p.title}</Link>
                            <span className="tiny muted">
                              {p.averageScore === null
                                ? 'chưa có điểm'
                                : `điểm TB ${formatScore(p.averageScore)}`}
                            </span>
                          </div>
                          <div className="progress-line">
                            <div className="progress-fill" style={{ width: `${pct}%` }} />
                          </div>
                          <span className="tiny muted">{p.completed}/{p.total} bài · {pct}%</span>
                        </div>
                      )
                    })}
                  </div>
                </section>
              )}
            </>
          )}
          </div>

          <div className="workspace-side">
          {/* --- Tác vụ AI gần đây, chỉ khi chức năng AI đang bật --- */}
          {ai?.enabled && (
          <section className="panel">
            <div className="spread" style={{ marginBottom: 10 }}>
              <h3 style={{ margin: 0 }}>Tác vụ AI gần đây</h3>
              <Link className="tiny" to="/tac-vu-ai">Xem tất cả</Link>
            </div>
            {tasks.length === 0 ? (
              <p className="muted tiny">Bạn chưa dùng chức năng AI nào.</p>
            ) : (
              <ul className="task-list">
                {tasks.map((t) => (
                  <li key={t.id}>
                    <span className="badge badge-primary">{TASK_LABEL[t.kind] ?? t.kind}</span>{' '}
                    {t.title || '—'}
                    <div className="tiny muted">
                      {formatDate(t.createdAt)}
                      {t.status === 'error' && <> · <b>lỗi</b></>}
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </section>
          )}
          </div>
        </div>
      </div>
    </>
  )
}

/** Trạng thái rỗng gọn cho một khối nhỏ, thay cho một dòng chữ xám trơ trọi. */
function EmptyLine({ icon, title, hint }: {
  icon: React.ReactNode
  title: string
  hint: string
}) {
  return (
    <div className="empty-line">
      <span className="empty-line-icon">{icon}</span>
      <div>
        <b>{title}</b>
        <div className="tiny muted">{hint}</div>
      </div>
    </div>
  )
}

function ShortcutCard({
  icon, title, desc, onClick, badge, color,
}: {
  icon: React.ReactNode
  title: string
  desc: string
  onClick: () => void
  badge?: number
  color?: 'violet' | 'blue' | 'emerald' | 'amber'
}) {
  return (
    <button type="button" className={`shortcut-card${color ? ` c-${color}` : ''}`} onClick={onClick}>
      <span className="shortcut-icon">{icon}</span>
      <span className="shortcut-text">
        <b>{title}</b>
        <span className="tiny muted">{desc}</span>
      </span>
      {badge ? <span className="badge badge-danger">{badge}</span> : null}
    </button>
  )
}
