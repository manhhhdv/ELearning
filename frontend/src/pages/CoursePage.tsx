import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'

import { api } from '../api/client'
import type { Program, TreeNode } from '../api/types'
import { useAuth } from '../auth'
import { AssignmentPlayer } from '../components/AssignmentPlayer'
import { CourseRail, itemMeta, learningSequence } from '../components/CourseRail'
import {
  IconCheck, IconChevron, IconExitFullscreen, IconFullscreen, IconPanelClose, IconPanelOpen,
} from '../components/icons'
import { LogoMark } from '../components/Logo'
import { RichContent } from '../components/RichContent'
import { Loading } from '../components/ui'

export function CoursePage() {
  const { programSlug = '', nodeSlug } = useParams()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [program, setProgram] = useState<Program | null>(null)
  const [tree, setTree] = useState<TreeNode[]>([])
  const [current, setCurrent] = useState<TreeNode | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  // Chỉ dùng ở màn hình hẹp; trên desktop mục lục luôn hiển thị.
  const [navOpen, setNavOpen] = useState(false)
  // Thu gọn mục lục trên desktop để dành chỗ cho nội dung.
  const [railHidden, setRailHidden] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)
  // Dự phòng khi trình duyệt chặn Fullscreen API (trang nhúng, chính sách của tổ chức):
  // phủ kín khung nhìn bằng CSS thay vì bỏ hẳn tính năng.
  const [immersive, setImmersive] = useState(false)
  const shellRef = useRef<HTMLDivElement>(null)

  // ID nội bộ của chương trình, tra một lần từ slug rồi dùng lại cho mọi lời gọi API.
  const programId = program?.id ?? ''

  const refreshTree = useCallback(async () => {
    if (!programId) return []
    const nodes = await api.getTree(programId)
    setTree(nodes)
    return nodes
  }, [programId])

  useEffect(() => {
    setLoading(true)
    api.getProgramBySlug(programSlug)
      .then(async (p) => {
        setProgram(p)
        const nodes = await api.getTree(p.id)
        setTree(nodes)
        const sequence = learningSequence(nodes)
        // Ưu tiên bài chỉ đích danh trên URL, nếu không có thì mở bài đầu tiên chưa hoàn thành.
        const wanted = nodeSlug ? sequence.find((n) => n.slug === nodeSlug) : undefined
        const target = wanted ?? sequence.find((n) => !n.completed) ?? sequence[0] ?? null
        if (target && target.kind !== 'folder') {
          api.getNode(target.id).then(setCurrent).catch(() => setCurrent(target))
        } else {
          setCurrent(target)
        }
        setError(null)
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Không mở được khoá học'))
      .finally(() => setLoading(false))
    // Cố ý bỏ nodeSlug khỏi phụ thuộc: chỉ dùng nó cho lần mở đầu tiên,
    // các lần đổi bài sau do hàm open ghi thẳng lên URL.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [programSlug])

  // Trình duyệt có thể rời toàn màn hình bằng phím Esc, phải theo dõi sự kiện để nút hiển thị đúng.
  useEffect(() => {
    const sync = () => setIsFullscreen(document.fullscreenElement !== null)
    document.addEventListener('fullscreenchange', sync)
    return () => document.removeEventListener('fullscreenchange', sync)
  }, [])

  const toggleFullscreen = async () => {
    if (document.fullscreenElement) {
      await document.exitFullscreen().catch(() => undefined)
      return
    }
    if (immersive) {
      setImmersive(false)
      return
    }
    try {
      await shellRef.current?.requestFullscreen()
    } catch {
      setImmersive(true)
    }
  }

  // Esc thoát chế độ phủ kín, giống hành vi quen thuộc của toàn màn hình thật.
  useEffect(() => {
    if (!immersive) return
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') setImmersive(false) }
    document.addEventListener('keydown', onKey)
    return () => document.removeEventListener('keydown', onKey)
  }, [immersive])

  const sequence = useMemo(() => learningSequence(tree), [tree])
  const index = current ? sequence.findIndex((n) => n.id === current.id) : -1
  const prev = index > 0 ? sequence[index - 1] : null
  const next = index >= 0 && index < sequence.length - 1 ? sequence[index + 1] : null

  const open = useCallback(async (node: TreeNode) => {
    try {
      // Cây chỉ có dữ liệu tóm tắt, lấy bản đầy đủ khi mở nội dung.
      setCurrent(node.kind === 'folder' ? node : await api.getNode(node.id))
      setNavOpen(false)
      // replace: chuyển bài không tạo thêm mục lịch sử, nút Back vẫn quay về danh sách khoá học.
      navigate(`/hoc/${programSlug}/${node.slug}`, { replace: true })
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không mở được nội dung')
    }
  }, [programSlug, navigate])

  const completed = current ? isDone(tree, current.id) : false

  const toggleComplete = async () => {
    if (!current) return
    await api.markComplete(current.id, !completed)
    await refreshTree()
  }

  /** Đánh dấu xong rồi chuyển luôn sang bài kế tiếp — thao tác thường dùng nhất. */
  const completeAndContinue = async () => {
    if (!current) return
    if (!completed) await api.markComplete(current.id, true)
    const nodes = await refreshTree()
    if (next) {
      const fresh = learningSequence(nodes).find((n) => n.id === next.id)
      if (fresh) await open(fresh)
    }
  }

  if (loading) return <Loading label="Đang mở khoá học…" />

  if (!program) {
    return (
      <div className="learner learner-shell">
        <div className="stage">
          <div className="callout warn">{error ?? 'Không tìm thấy khoá học'}</div>
          <Link className="cbtn" to="/hoc">Về danh sách khoá học</Link>
        </div>
      </div>
    )
  }

  const initials = (user?.fullName || user?.email || '?')
    .split(' ').filter(Boolean).slice(-2).map((p) => p[0]?.toUpperCase()).join('')

  return (
    <div
      className={[
        'learner learner-shell',
        railHidden ? 'rail-hidden' : '',
        immersive ? 'immersive' : '',
      ].filter(Boolean).join(' ')}
      ref={shellRef}
    >
      <header className="course-topbar">
        <Link to="/hoc" aria-label="Tập Huấn — trang chủ" style={{ display: 'inline-flex' }}>
          <LogoMark size={32} />
        </Link>
        <Link className="back" to="/hoc">
          <span style={{ transform: 'rotate(180deg)', display: 'inline-flex' }}><IconChevron size={13} /></span>
          Khoá học của tôi
        </Link>
        <button
          className="cbtn cbtn-plain cbtn-sm rail-toggle"
          onClick={() => setNavOpen((v) => !v)}
          aria-expanded={navOpen}
        >
          Nội dung
        </button>

        {/* Chỉ hiện trên desktop: thu gọn hẳn cột mục lục. */}
        <button
          className="icon-square rail-collapse"
          onClick={() => setRailHidden((v) => !v)}
          title={railHidden ? 'Hiện mục lục khoá học' : 'Thu gọn mục lục khoá học'}
          aria-label={railHidden ? 'Hiện mục lục khoá học' : 'Thu gọn mục lục khoá học'}
        >
          {railHidden ? <IconPanelOpen /> : <IconPanelClose />}
        </button>
        <div className="course-name">{program.title}</div>
        <div className="who">
          <button
            className="icon-square"
            onClick={toggleFullscreen}
            title={isFullscreen || immersive ? 'Thoát toàn màn hình' : 'Toàn màn hình'}
            aria-label={isFullscreen || immersive ? 'Thoát toàn màn hình' : 'Toàn màn hình'}
          >
            {isFullscreen || immersive ? <IconExitFullscreen /> : <IconFullscreen />}
          </button>
          <span className="pill">{program.code}</span>
          <div className="avatar" title={user?.fullName || user?.email}>{initials}</div>
        </div>
      </header>

      <div className="course-body">
        <CourseRail tree={tree} currentId={current?.id ?? null} onOpen={open} open={navOpen} />

        <main className="stage">
          {error && <div className="callout warn">{error}</div>}

          {!current ? (
            <div className="blank">
              <h3>Khoá học chưa có nội dung</h3>
              <p>Giảng viên chưa đăng bài học nào cho khoá này.</p>
            </div>
          ) : current.kind === 'folder' ? (
            <>
              <div className="crumb">{program.title}</div>
              <h1>{current.title}</h1>
              {current.description && <p className="prose prose-muted" style={{ marginTop: 12 }}>{current.description}</p>}
              <p className="prose prose-muted" style={{ marginTop: 12 }}>
                Chọn một mục trong chương này ở thanh bên trái để bắt đầu.
              </p>
            </>
          ) : current.kind === 'lesson' ? (
            <>
              <div className="crumb">{program.title}</div>
              <h1>{current.title}</h1>
              <div className="stage-meta">
                <span>{itemMeta(current)}</span>
                {completed && (
                  <>
                    <span className="dot">·</span>
                    <span className="pill pill-green"><IconCheck size={12} /> Đã hoàn thành</span>
                  </>
                )}
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
            <AssignmentPlayer
              key={current.id}
              node={current}
              onSubmitted={refreshTree}
              programTitle={program.title}
              prevNode={prev}
              nextNode={next}
              onNavigate={open}
            />
          )}

          {/* Bài tập tự dựng thanh đáy riêng vì lúc đang làm bài cần nút nộp thay cho điều hướng. */}
          {current?.kind === 'lesson' && (
            <div className="stage-nav">
              <button className="cbtn cbtn-plain cbtn-sm" disabled={!prev} onClick={() => prev && open(prev)}>
                <span style={{ transform: 'rotate(180deg)', display: 'inline-flex' }}><IconChevron size={13} /></span>
                Bài trước
              </button>
              <div className="spacer" />

              {completed ? (
                <button className="cbtn cbtn-done cbtn-sm" onClick={toggleComplete}>
                  <IconCheck size={14} /> Đã hoàn thành
                </button>
              ) : (
                <button className="cbtn cbtn-sm" onClick={toggleComplete}>Đánh dấu hoàn thành</button>
              )}

              {next ? (
                <button className="cbtn cbtn-fill cbtn-sm" onClick={completeAndContinue}>
                  Bài tiếp theo <IconChevron size={13} />
                </button>
              ) : (
                !completed && (
                  <button className="cbtn cbtn-fill cbtn-sm" onClick={toggleComplete}>
                    Hoàn thành khoá học
                  </button>
                )
              )}
            </div>
          )}

        </main>
      </div>
    </div>
  )
}

/** Cờ hoàn thành chỉ có trên dữ liệu cây, không có ở API lấy chi tiết một nút. */
function isDone(nodes: TreeNode[], id: string): boolean {
  for (const n of nodes) {
    if (n.id === id) return n.completed ?? false
    const hit = isDone(n.children, id)
    if (hit) return true
  }
  return false
}
