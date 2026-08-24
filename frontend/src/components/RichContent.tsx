import { useMemo } from 'react'

import { renderRichText } from './richtext'

/**
 * Hiển thị nội dung tự soạn đã kết xuất (markdown + LaTeX + media).
 * HTML luôn đi qua bộ lọc trong richtext.ts nên an toàn để gắn thẳng vào trang.
 */
export function RichContent({ value, className }: { value: string; className?: string }) {
  const html = useMemo(() => renderRichText(value), [value])
  if (!html) return null
  return <div className={className ? `rich ${className}` : 'rich'} dangerouslySetInnerHTML={{ __html: html }} />
}
