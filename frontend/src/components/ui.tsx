import type { ReactNode } from 'react'

import type { NodeKind, ContentType } from '../api/types'
import { IconDoc, IconFolder, IconSlide, IconTask, IconVideo } from './icons'

export function Loading({ label = 'Đang tải…' }: { label?: string }) {
  return <div className="loading">{label}</div>
}

export function ErrorAlert({ message }: { message: string | null }) {
  if (!message) return null
  return <div className="alert alert-error">{message}</div>
}

export function SuccessAlert({ message }: { message: string | null }) {
  if (!message) return null
  return <div className="alert alert-success">{message}</div>
}

export function EmptyState({ title, children }: { title: string; children?: ReactNode }) {
  return (
    <div className="empty">
      <h3>{title}</h3>
      {children}
    </div>
  )
}

export function Modal({
  title, onClose, children, footer, wide,
}: {
  title: string
  onClose: () => void
  children: ReactNode
  footer?: ReactNode
  wide?: boolean
}) {
  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div
        className="modal"
        style={wide ? { maxWidth: 640 } : undefined}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
      >
        <div className="modal-head">
          <h2>{title}</h2>
          <button className="btn btn-ghost btn-sm" onClick={onClose} aria-label="Đóng">✕</button>
        </div>
        <div className="modal-body">{children}</div>
        {footer && <div className="modal-foot">{footer}</div>}
      </div>
    </div>
  )
}

/** Biểu tượng tương ứng loại nút trên cây; bài học đổi theo loại nội dung. */
export function NodeIcon({ kind, contentType }: { kind: NodeKind; contentType?: ContentType }) {
  if (kind === 'folder') return <IconFolder />
  if (kind === 'assignment') return <IconTask />
  switch (contentType) {
    case 'slide': return <IconSlide />
    case 'document':
    case 'pdf': return <IconDoc />
    default: return <IconVideo />
  }
}

export function formatDate(value: string | null | undefined) {
  if (!value) return '—'
  return new Date(value).toLocaleString('vi-VN', {
    day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}

/** Bỏ phần thập phân thừa: 7.50 → 7,5 */
export function formatScore(value: number) {
  return Number(value.toFixed(2)).toLocaleString('vi-VN')
}
