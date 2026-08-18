import { useEffect, useState } from 'react'

import type { TreeNode } from '../api/types'
import { CONTENT_TYPE_LABEL } from '../api/types'
import { IconCheck, IconChevron, IconDoc, IconSlide, IconTask, IconVideo } from './icons'

/** Một mục học được: bài học hoặc bài tập, kèm độ sâu để thụt lề. */
export interface CourseItem {
  node: TreeNode
  depth: number
}

/** Một chương: thư mục cấp một, chứa các mục bên trong đã được làm phẳng. */
export interface CourseModule {
  folder: TreeNode | null
  items: CourseItem[]
}

/** Làm phẳng nhánh con thành danh sách mục học, giữ nguyên thứ tự hiển thị. */
function flattenItems(nodes: TreeNode[], depth = 0): CourseItem[] {
  return nodes.flatMap((n) =>
    n.kind === 'folder'
      ? flattenItems(n.children, depth + 1)
      : [{ node: n, depth }],
  )
}

/**
 * Gom cây thành các chương: mỗi thư mục cấp một là một chương,
 * các mục nằm thẳng ở cấp gốc được gộp vào một chương không tên đứng đầu.
 */
export function buildModules(tree: TreeNode[]): CourseModule[] {
  const modules: CourseModule[] = []
  const loose = tree.filter((n) => n.kind !== 'folder')
  if (loose.length > 0) {
    modules.push({ folder: null, items: loose.map((node) => ({ node, depth: 0 })) })
  }
  for (const node of tree) {
    if (node.kind === 'folder') {
      modules.push({ folder: node, items: flattenItems(node.children) })
    }
  }
  return modules
}

/** Danh sách phẳng theo đúng thứ tự học, dùng cho nút Bài trước / Bài tiếp theo. */
export function learningSequence(tree: TreeNode[]): TreeNode[] {
  return buildModules(tree).flatMap((m) => m.items.map((i) => i.node))
}

/**
 * Tiêu đề chương hiển thị trên thanh bên. Chỉ tự thêm tiền tố "Chương N" khi
 * người soạn chưa tự đánh số, tránh ra chuỗi kiểu "Chương 1: Chương 1 — ...".
 */
export function moduleHeading(title: string, order: number): string {
  const numbered = /^\s*(chương|chuong|module|phần|phan|tuần|tuan|unit|bài|bai|part)\b/i.test(title)
    || /^\s*\d+\s*[.):\-–—]/.test(title)
  return numbered ? title : `Chương ${order}: ${title}`
}

export function itemIcon(node: TreeNode) {
  if (node.kind === 'assignment') return <IconTask size={16} />
  switch (node.lesson?.contentType) {
    case 'slide': return <IconSlide size={16} />
    case 'document':
    case 'pdf': return <IconDoc size={16} />
    default: return <IconVideo size={16} />
  }
}

/** Dòng mô tả ngắn dưới tiêu đề mục: loại nội dung và thời lượng / số câu hỏi. */
export function itemMeta(node: TreeNode): string {
  if (node.kind === 'assignment') {
    const count = node.assignment?.questionCount ?? 0
    return count > 0 ? `Bài tập · ${count} câu hỏi` : 'Bài tập'
  }
  const kindLabel = CONTENT_TYPE_LABEL[node.lesson?.contentType ?? 'video']
  const minutes = node.lesson?.durationMinutes ?? 0
  return minutes > 0 ? `${kindLabel} · ${minutes} phút` : kindLabel
}

interface Props {
  tree: TreeNode[]
  currentId: string | null
  onOpen: (node: TreeNode) => void
  /** Chỉ có tác dụng ở màn hình hẹp, nơi mục lục được gập lại mặc định. */
  open?: boolean
}

export function CourseRail({ tree, currentId, onOpen, open = false }: Props) {
  const modules = buildModules(tree)
  const [closed, setClosed] = useState<Set<string>>(new Set())

  // Luôn mở sẵn chương chứa bài đang xem, kể cả khi người dùng vừa đóng nó.
  useEffect(() => {
    if (!currentId) return
    const owner = modules.find((m) => m.items.some((i) => i.node.id === currentId))
    if (owner?.folder) {
      setClosed((prev) => {
        if (!prev.has(owner.folder!.id)) return prev
        const next = new Set(prev)
        next.delete(owner.folder!.id)
        return next
      })
    }
  }, [currentId]) // eslint-disable-line react-hooks/exhaustive-deps

  const lessons = modules.flatMap((m) => m.items).filter((i) => i.node.kind === 'lesson')
  const done = lessons.filter((i) => i.node.completed).length
  const percent = lessons.length ? Math.round((done / lessons.length) * 100) : 0

  const toggle = (id: string) => {
    setClosed((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  return (
    <nav className={`rail ${open ? 'open' : ''}`} aria-label="Nội dung khoá học">
      <div className="rail-head">
        <h2>Nội dung khoá học</h2>
        <div className="rail-progress">
          <div className="bar"><i style={{ width: `${percent}%` }} /></div>
          <span className="pct">{percent}%</span>
        </div>
        <div style={{ fontSize: 12.5, color: 'var(--c-muted)', marginTop: 6 }}>
          Đã hoàn thành {done}/{lessons.length} bài học
        </div>
      </div>

      {modules.map((m, index) => {
        const id = m.folder?.id ?? '__loose__'
        const isOpen = !closed.has(id)
        const modLessons = m.items.filter((i) => i.node.kind === 'lesson')
        const modDone = modLessons.filter((i) => i.node.completed).length

        return (
          <section className="module" key={id}>
            {m.folder && (
              <button className="module-head" onClick={() => toggle(id)} aria-expanded={isOpen}>
                <span className={`caret ${isOpen ? 'open' : ''}`}><IconChevron size={13} /></span>
                <span className="label">
                  <b>{moduleHeading(m.folder.title, index + (modules[0]?.folder ? 1 : 0))}</b>
                  <span>
                    {m.items.length} mục
                    {modLessons.length > 0 && ` · hoàn thành ${modDone}/${modLessons.length}`}
                  </span>
                </span>
              </button>
            )}

            {isOpen && (
              <div className="module-items">
                {m.items.length === 0 && (
                  <p style={{ padding: '4px 30px 8px', fontSize: 13, color: 'var(--c-faint)' }}>
                    Chương này chưa có nội dung.
                  </p>
                )}
                {m.items.map(({ node, depth }) => (
                  <button
                    key={node.id}
                    className={`item-row ${currentId === node.id ? 'active' : ''}`}
                    style={depth > 1 ? { paddingLeft: 30 + (depth - 1) * 14 } : undefined}
                    onClick={() => onOpen(node)}
                  >
                    <span className={`tick ${node.completed ? 'done' : ''}`}>
                      <IconCheck size={12} />
                    </span>
                    <span className="txt">
                      <b>{node.title}</b>
                      <span>{itemIcon(node)} {itemMeta(node)}</span>
                    </span>
                  </button>
                ))}
              </div>
            )}
          </section>
        )
      })}
    </nav>
  )
}
