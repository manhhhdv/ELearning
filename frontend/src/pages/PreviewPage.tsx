import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Program, Question, TreeNode } from '../api/types'
import { QUESTION_TYPE_LABEL } from '../api/types'
import { CourseRail, itemMeta, learningSequence } from '../components/CourseRail'
import { IconChevron, IconEye } from '../components/icons'
import { RichContent } from '../components/RichContent'
import { Loading } from '../components/ui'

/**
 * Xem trước nội dung đúng như học viên sẽ thấy, nhưng dùng dữ liệu của người
 * quản lý (thấy cả bài chưa xuất bản, đáp án đúng) và không ghi lại bất kỳ
 * hành động nào — không có nút hoàn thành, không nộp bài thật.
 */
export function PreviewPage() {
  const { programSlug = '', nodeSlug } = useParams()
  const navigate = useNavigate()
  const [program, setProgram] = useState<Program | null>(null)
  const [tree, setTree] = useState<TreeNode[]>([])
  const [current, setCurrent] = useState<TreeNode | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    api.getProgramBySlug(programSlug)
      .then(async (p) => {
        setProgram(p)
        const nodes = await api.getTree(p.id)
        setTree(nodes)
        const sequence = learningSequence(nodes)
        const wanted = nodeSlug ? sequence.find((n) => n.slug === nodeSlug) : undefined
        const target = wanted ?? sequence[0] ?? null
        setCurrent(target && target.kind !== 'folder' ? await api.getNode(target.id) : target)
        setError(null)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không mở được nội dung'))
      .finally(() => setLoading(false))
  }, [programSlug, nodeSlug])

  const sequence = useMemo(() => learningSequence(tree), [tree])
  const index = current ? sequence.findIndex((n) => n.id === current.id) : -1
  const prev = index > 0 ? sequence[index - 1] : null
  const next = index >= 0 && index < sequence.length - 1 ? sequence[index + 1] : null

  const open = async (node: TreeNode) => {
    try {
      setCurrent(node.kind === 'folder' ? node : await api.getNode(node.id))
      navigate(`/xem-truoc/${programSlug}/${node.slug}`, { replace: true })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được nội dung')
    }
  }

  if (loading) return <Loading label="Đang mở bản xem trước…" />
  if (!program) {
    return (
      <div className="learner learner-shell">
        <div className="stage"><div className="callout warn">{error ?? 'Không tìm thấy nội dung'}</div></div>
      </div>
    )
  }

  return (
    <div className="learner learner-shell">
      <header className="course-topbar">
        <button className="cbtn cbtn-plain cbtn-sm" onClick={() => window.close()} title="Đóng bản xem trước">
          <span style={{ transform: 'rotate(180deg)', display: 'inline-flex' }}><IconChevron size={13} /></span>
          Đóng
        </button>
        <div className="course-name">
          <span className="pill pill-amber" style={{ marginRight: 8 }}><IconEye size={12} /> Xem trước</span>
          {program.title}
        </div>
        <div className="who">
          <Link className="cbtn cbtn-plain cbtn-sm" to={`/quan-tri/chuong-trinh/${program.slug}`}>
            Về trang soạn thảo
          </Link>
        </div>
      </header>

      <div className="callout warn" style={{ margin: '10px 16px 0' }}>
        Đây là bản xem trước dành cho người soạn — không tính tiến độ, không lưu bài làm, kể cả nội
        dung chưa xuất bản cũng hiển thị ở đây.
      </div>

      <div className="course-body">
        <CourseRail tree={tree} currentId={current?.id ?? null} onOpen={open} />

        <main className="stage">
          {error && <div className="callout warn">{error}</div>}

          {!current ? (
            <div className="blank">
              <h3>Chưa có nội dung</h3>
              <p>Thêm bài học hoặc bài tập ở trang soạn thảo để xem trước tại đây.</p>
            </div>
          ) : current.kind === 'folder' ? (
            <>
              <div className="crumb">{program.title}</div>
              <h1>{current.title}</h1>
              {current.description && <p className="prose prose-muted" style={{ marginTop: 12 }}>{current.description}</p>}
            </>
          ) : current.kind === 'lesson' ? (
            <>
              <div className="crumb">{program.title}</div>
              <h1>{current.title}</h1>
              <div className="stage-meta">
                <span>{itemMeta(current)}</span>
                {!current.isPublished && <span className="pill pill-amber">Chưa xuất bản</span>}
              </div>
              {current.description && (
                <p className="prose prose-muted" style={{ marginTop: 14 }}>{current.description}</p>
              )}
              {current.lesson?.contentType === 'richtext' ? (
                current.lesson.body.trim() ? (
                  <RichContent value={current.lesson.body} className="lesson-rich" />
                ) : (
                  <div className="callout warn">Bài học này chưa có nội dung.</div>
                )
              ) : (
                <>
                  {current.lesson?.embedUrl ? (
                    <iframe
                      className={`player ${current.lesson.contentType === 'video' ? '' : 'doc'}`}
                      src={current.lesson.embedUrl}
                      title={current.title}
                      allow="autoplay; fullscreen"
                      allowFullScreen
                    />
                  ) : (
                    <div className="callout warn">Bài học này chưa được gắn nội dung.</div>
                  )}

                  {current.lesson?.body && (
                    <div className="note-box">
                      <h3>Ghi chú bài học</h3>
                      <RichContent value={current.lesson.body} />
                    </div>
                  )}
                </>
              )}
            </>
          ) : (
            <AssignmentPreview node={current} programTitle={program.title} />
          )}

          <div className="stage-nav">
            <button className="cbtn cbtn-plain cbtn-sm" disabled={!prev} onClick={() => prev && open(prev)}>
              <span style={{ transform: 'rotate(180deg)', display: 'inline-flex' }}><IconChevron size={13} /></span>
              Trước
            </button>
            <div className="spacer" />
            <button className="cbtn cbtn-plain cbtn-sm" disabled={!next} onClick={() => next && open(next)}>
              Tiếp theo <IconChevron size={13} />
            </button>
          </div>
        </main>
      </div>
    </div>
  )
}

/** Xem đề bài tập kèm đáp án đúng — không có ô nhập, không nộp được. */
function AssignmentPreview({ node, programTitle }: { node: TreeNode; programTitle: string }) {
  const assignment = node.assignment
  const questions = assignment?.questions ?? []

  return (
    <>
      <div className="crumb">{programTitle}</div>
      <h1>{node.title}</h1>
      <div className="stage-meta">
        <span className="pill pill-blue">Bài tập</span>
        <span>{questions.length} câu hỏi</span>
        {!node.isPublished && <span className="pill pill-amber">Chưa xuất bản</span>}
      </div>

      {assignment?.instructions && <div className="callout">{assignment.instructions}</div>}

      {questions.length === 0 ? (
        <p className="prose prose-muted" style={{ marginTop: 16 }}>Bài tập này chưa có câu hỏi.</p>
      ) : (
        questions.map((q, i) => <PreviewQuestion key={q.id} question={q} index={i} />)
      )}
    </>
  )
}

function PreviewQuestion({ question, index }: { question: Question; index: number }) {
  return (
    <article className="qcard">
      <div className="qcard-top">
        <span className="no">Câu {index + 1}</span>
        <span className="pill mono">{question.code}</span>
        <span className="grow" />
        <span className="pill">{question.points} điểm</span>
        <span className="pill pill-blue">{QUESTION_TYPE_LABEL[question.type]}</span>
      </div>
      <div className="stem">{question.prompt}</div>

      {question.type === 'essay' ? (
        <div className="essay-answer" style={{ color: 'var(--c-muted)', fontStyle: 'italic' }}>
          Học viên sẽ nhập câu trả lời tự luận tại đây.
        </div>
      ) : (
        question.options.map((o) => (
          <div className={`choice review ${o.isCorrect ? 'right' : ''}`} key={o.id}>
            <span style={{ width: 17, textAlign: 'center', marginTop: 2 }}>{o.isCorrect ? '●' : '○'}</span>
            <span className="body">{o.content}</span>
            {o.isCorrect && <span className="pill pill-green">đáp án đúng</span>}
          </div>
        ))
      )}

      {question.explanation && (
        <div className="callout" style={{ margin: '14px 0 0' }}>
          <b>Giải thích</b>
          <div style={{ marginTop: 4 }}>{question.explanation}</div>
        </div>
      )}
    </article>
  )
}
