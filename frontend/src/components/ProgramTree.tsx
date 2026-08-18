import { useState } from 'react'

import type { TreeNode } from '../api/types'
import { IconChevron } from './icons'
import { NodeIcon } from './ui'

/** Vị trí thả so với nút đích. */
type DropZone = 'before' | 'inside' | 'after'

interface Props {
  nodes: TreeNode[]
  selectedId: string | null
  onSelect: (node: TreeNode) => void
  /** Bỏ trống để hiển thị cây ở chế độ chỉ đọc. */
  onMove?: (nodeId: string, parentId: string | null, position: number) => void
  /** Nội dung phụ hiển thị bên phải mỗi dòng (ví dụ dấu tích hoàn thành). */
  renderMeta?: (node: TreeNode) => React.ReactNode
}

export function ProgramTree({ nodes, selectedId, onSelect, onMove, renderMeta }: Props) {
  // Mặc định mở toàn bộ thư mục để người dùng thấy ngay cấu trúc.
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set())
  const [dragId, setDragId] = useState<string | null>(null)
  const [dropTarget, setDropTarget] = useState<{ id: string; zone: DropZone } | null>(null)

  const toggle = (id: string) => {
    setCollapsed((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  /** Xác định thả vào trong hay chèn trước/sau dựa trên vị trí con trỏ trong dòng. */
  const zoneFor = (e: React.DragEvent, node: TreeNode): DropZone => {
    const rect = e.currentTarget.getBoundingClientRect()
    const ratio = (e.clientY - rect.top) / rect.height
    if (node.kind === 'folder') {
      if (ratio < 0.25) return 'before'
      if (ratio > 0.75) return 'after'
      return 'inside'
    }
    return ratio < 0.5 ? 'before' : 'after'
  }

  const handleDrop = (e: React.DragEvent, target: TreeNode, siblings: TreeNode[]) => {
    e.preventDefault()
    e.stopPropagation()
    setDropTarget(null)
    const sourceId = dragId
    setDragId(null)
    if (!onMove || !sourceId || sourceId === target.id) return

    const zone = zoneFor(e, target)
    if (zone === 'inside') {
      onMove(sourceId, target.id, target.children.length)
      return
    }
    const index = siblings.findIndex((n) => n.id === target.id)
    onMove(sourceId, target.parentId, zone === 'before' ? index : index + 1)
  }

  const renderNodes = (list: TreeNode[], depth: number) => (
    <>
      {list.map((node) => {
        const isOpen = !collapsed.has(node.id)
        const hasChildren = node.children.length > 0
        const drop = dropTarget?.id === node.id ? dropTarget.zone : null

        return (
          <div key={node.id}>
            <div
              className={[
                'tree-row',
                selectedId === node.id ? 'selected' : '',
                node.isPublished ? '' : 'unpublished',
                dragId === node.id ? 'dragging' : '',
                drop === 'inside' ? 'drop-inside' : '',
                drop === 'before' ? 'drop-before' : '',
                drop === 'after' ? 'drop-after' : '',
              ].filter(Boolean).join(' ')}
              onClick={() => onSelect(node)}
              draggable={!!onMove}
              onDragStart={(e) => {
                e.stopPropagation()
                setDragId(node.id)
                e.dataTransfer.effectAllowed = 'move'
              }}
              onDragEnd={() => { setDragId(null); setDropTarget(null) }}
              onDragOver={(e) => {
                if (!onMove || !dragId || dragId === node.id) return
                e.preventDefault()
                e.stopPropagation()
                e.dataTransfer.dropEffect = 'move'
                setDropTarget({ id: node.id, zone: zoneFor(e, node) })
              }}
              onDragLeave={() => setDropTarget((cur) => (cur?.id === node.id ? null : cur))}
              onDrop={(e) => handleDrop(e, node, list)}
              title={node.title}
            >
              <span
                className={`tree-caret ${hasChildren ? '' : 'leaf'} ${isOpen ? 'open' : ''}`}
                onClick={(e) => { e.stopPropagation(); toggle(node.id) }}
              >
                <IconChevron />
              </span>
              <span className="tree-icon">
                <NodeIcon kind={node.kind} contentType={node.lesson?.contentType} />
              </span>
              <span className="tree-title">{node.title}</span>
              {renderMeta?.(node)}
              {node.kind === 'assignment' && node.assignment && (
                <span className="tree-meta">{node.assignment.questionCount} câu</span>
              )}
              {!node.isPublished && <span className="tree-meta">ẩn</span>}
            </div>

            {hasChildren && isOpen && (
              <div className="tree-children">{renderNodes(node.children, depth + 1)}</div>
            )}
          </div>
        )
      })}
    </>
  )

  return (
    <div
      className="tree-scroll"
      // Thả vào vùng trống = đưa nút ra cấp gốc, cuối danh sách.
      onDragOver={(e) => { if (onMove && dragId) e.preventDefault() }}
      onDrop={(e) => {
        e.preventDefault()
        if (onMove && dragId) onMove(dragId, null, nodes.length)
        setDragId(null)
        setDropTarget(null)
      }}
    >
      {nodes.length === 0
        ? <p className="muted tiny" style={{ padding: 12 }}>Chưa có nội dung nào. Bấm “Thêm” để bắt đầu.</p>
        : renderNodes(nodes, 0)}
    </div>
  )
}
