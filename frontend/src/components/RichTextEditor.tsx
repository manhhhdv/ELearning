import { useLayoutEffect, useRef, useState } from 'react'

import { RichContent } from './RichContent'

type Mode = 'write' | 'split' | 'preview'

interface Props {
  id?: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  disabled?: boolean
  /** Số dòng của ô soạn thảo; nội dung bài học cần cao hơn ghi chú ngắn. */
  rows?: number
}

const MODE_LABEL: Record<Mode, string> = {
  write: 'Soạn thảo',
  split: 'Chia đôi',
  preview: 'Xem trước',
}

/**
 * Ô soạn thảo markdown có thanh công cụ và khung xem trước trực tiếp.
 * Nội dung lưu xuống vẫn là văn bản thuần (markdown) nên dễ sao chép, so sánh và nhập từ file.
 */
export function RichTextEditor({
  id, value, onChange, placeholder, disabled = false, rows = 14,
}: Props) {
  const area = useRef<HTMLTextAreaElement>(null)
  const [mode, setMode] = useState<Mode>('split')
  const [panel, setPanel] = useState<'link' | 'media' | null>(null)
  const [url, setUrl] = useState('')
  const [label, setLabel] = useState('')
  // Vị trí con trỏ lúc mở bảng chèn — nhấn vào ô URL sẽ làm textarea mất tiêu điểm.
  const anchor = useRef<[number, number]>([0, 0])

  // Vùng chọn cần đặt lại sau khi React vẽ xong giá trị mới, nếu không con trỏ sẽ nhảy về cuối.
  const pending = useRef<[number, number] | null>(null)
  useLayoutEffect(() => {
    if (!pending.current) return
    const [from, to] = pending.current
    pending.current = null
    const el = area.current
    if (!el) return
    el.focus()
    el.setSelectionRange(from, to)
  })

  /** Ghi giá trị mới rồi trả tiêu điểm cùng vùng chọn về ô soạn thảo. */
  const commit = (next: string, from: number, to: number) => {
    pending.current = [from, to]
    onChange(next)
  }

  const selection = (): [number, number] => {
    const el = area.current
    if (!el) return [value.length, value.length]
    return [el.selectionStart, el.selectionEnd]
  }

  /** Bọc vùng đang chọn giữa hai chuỗi; nếu chưa chọn gì thì chèn chữ mẫu. */
  const wrap = (before: string, after: string, sample: string) => {
    const [start, end] = selection()
    const picked = value.slice(start, end) || sample
    const next = `${value.slice(0, start)}${before}${picked}${after}${value.slice(end)}`
    commit(next, start + before.length, start + before.length + picked.length)
  }

  /** Thêm tiền tố vào từng dòng của vùng chọn (danh sách, trích dẫn, tiêu đề…). */
  const prefixLines = (make: (index: number) => string) => {
    const [start, end] = selection()
    const from = value.lastIndexOf('\n', start - 1) + 1
    const to = value.indexOf('\n', end) === -1 ? value.length : value.indexOf('\n', end)
    const block = value.slice(from, to) || ''
    const lines = block.split('\n').map((line, i) => `${make(i)}${line}`)
    const next = `${value.slice(0, from)}${lines.join('\n')}${value.slice(to)}`
    commit(next, from, from + lines.join('\n').length)
  }

  /** Chèn một khối riêng, tự chừa dòng trống trước và sau. */
  const insertBlock = (block: string, cursorOffset = block.length) => {
    const [start, end] = selection()
    const before = value.slice(0, start).replace(/\n*$/, '')
    const after = value.slice(end).replace(/^\n*/, '')
    const head = before ? `${before}\n\n` : ''
    const tail = after ? `\n\n${after}` : '\n'
    const next = `${head}${block}${tail}`
    commit(next, head.length + cursorOffset, head.length + cursorOffset)
  }

  const openPanel = (kind: 'link' | 'media') => {
    const [start, end] = selection()
    anchor.current = [start, end]
    setLabel(value.slice(start, end))
    setUrl('')
    setPanel(kind)
  }

  const insertFromPanel = () => {
    const address = url.trim()
    if (!address) return
    const [start, end] = anchor.current
    const text = label.trim() || (panel === 'media' ? 'Media' : address)
    const snippet = panel === 'media' ? `![${text}](${address})` : `[${text}](${address})`
    const next = `${value.slice(0, start)}${snippet}${value.slice(end)}`
    setPanel(null)
    commit(next, start + snippet.length, start + snippet.length)
  }

  const onKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (!(e.metaKey || e.ctrlKey)) return
    const key = e.key.toLowerCase()
    if (key === 'b') { e.preventDefault(); wrap('**', '**', 'chữ đậm') }
    if (key === 'i') { e.preventDefault(); wrap('*', '*', 'chữ nghiêng') }
    if (key === 'k') { e.preventDefault(); openPanel('link') }
  }

  const tools: { label: string; title: string; run: () => void }[] = [
    { label: 'B', title: 'In đậm (Ctrl/Cmd + B)', run: () => wrap('**', '**', 'chữ đậm') },
    { label: 'I', title: 'In nghiêng (Ctrl/Cmd + I)', run: () => wrap('*', '*', 'chữ nghiêng') },
    { label: 'H2', title: 'Tiêu đề mục', run: () => prefixLines(() => '## ') },
    { label: '• Danh sách', title: 'Danh sách gạch đầu dòng', run: () => prefixLines(() => '- ') },
    { label: '1. Đánh số', title: 'Danh sách đánh số', run: () => prefixLines((i) => `${i + 1}. `) },
    { label: '❝ Trích', title: 'Trích dẫn', run: () => prefixLines(() => '> ') },
    { label: '</>', title: 'Mã trong dòng', run: () => wrap('`', '`', 'mã') },
    { label: 'Khối mã', title: 'Khối mã nhiều dòng', run: () => insertBlock('```\nnội dung mã\n```', 4) },
    { label: 'Liên kết', title: 'Chèn liên kết (Ctrl/Cmd + K)', run: () => openPanel('link') },
    { label: 'Media', title: 'Chèn ảnh, video, audio hoặc link YouTube/Drive', run: () => openPanel('media') },
    { label: '∑ Công thức', title: 'Công thức LaTeX trong dòng', run: () => wrap('$', '$', 'a^2 + b^2 = c^2') },
    {
      label: '∑ Khối',
      title: 'Công thức LaTeX riêng một dòng',
      run: () => insertBlock('$$\n\\int_0^1 x^2\\,dx = \\frac{1}{3}\n$$', 3),
    },
    {
      label: '▦ Bảng',
      title: 'Chèn bảng',
      run: () => insertBlock('| Cột 1 | Cột 2 |\n| --- | --- |\n| Giá trị | Giá trị |', 2),
    },
    { label: '— Ngăn', title: 'Đường ngăn cách', run: () => insertBlock('---') },
  ]

  return (
    <div className={`rte rte-${mode}`}>
      <div className="rte-bar">
        <div className="rte-tools">
          {tools.map((tool) => (
            <button
              key={tool.label}
              type="button"
              className="rte-btn"
              title={tool.title}
              disabled={disabled}
              onClick={tool.run}
            >
              {tool.label}
            </button>
          ))}
        </div>
        <div className="rte-modes">
          {(Object.keys(MODE_LABEL) as Mode[]).map((value_) => (
            <button
              key={value_}
              type="button"
              className={`rte-mode ${mode === value_ ? 'is-on' : ''}`}
              onClick={() => setMode(value_)}
            >
              {MODE_LABEL[value_]}
            </button>
          ))}
        </div>
      </div>

      {panel && (
        <div className="rte-panel">
          <input
            type="text"
            autoFocus
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') { e.preventDefault(); insertFromPanel() }
              if (e.key === 'Escape') setPanel(null)
            }}
            placeholder={panel === 'media'
              ? 'Link ảnh / video / audio / YouTube / Google Drive'
              : 'https://…'}
          />
          <input
            type="text"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') { e.preventDefault(); insertFromPanel() }
              if (e.key === 'Escape') setPanel(null)
            }}
            placeholder={panel === 'media' ? 'Mô tả ngắn (alt)' : 'Chữ hiển thị'}
          />
          <button type="button" className="btn btn-primary btn-sm" onClick={insertFromPanel}>Chèn</button>
          <button type="button" className="btn btn-sm" onClick={() => setPanel(null)}>Huỷ</button>
        </div>
      )}

      <div className="rte-panes">
        {mode !== 'preview' && (
          <textarea
            id={id}
            ref={area}
            className="rte-input"
            rows={rows}
            value={value}
            disabled={disabled}
            placeholder={placeholder}
            onChange={(e) => onChange(e.target.value)}
            onKeyDown={onKeyDown}
            spellCheck={false}
          />
        )}
        {mode !== 'write' && (
          <div className="rte-preview">
            {value.trim()
              ? <RichContent value={value} />
              : <p className="tiny muted">Chưa có nội dung để xem trước.</p>}
          </div>
        )}
      </div>

      <details className="rte-help">
        <summary>Cú pháp nhanh</summary>
        <ul>
          <li><code>**đậm**</code>, <code>*nghiêng*</code>, <code>~~gạch ngang~~</code>, <code>`mã`</code></li>
          <li><code># Tiêu đề</code>, <code>- gạch đầu dòng</code>, <code>1. đánh số</code>, <code>- [ ] việc cần làm</code></li>
          <li><code>[chữ](https://…)</code> — liên kết; <code>![mô tả](link)</code> — ảnh, video, audio, YouTube, Drive</li>
          <li><code>$a^2+b^2$</code> — công thức trong dòng; <code>$$…$$</code> — công thức riêng dòng (LaTeX)</li>
          <li><code>&gt; [!NOTE]</code> mở đầu trích dẫn để tạo hộp nhấn mạnh (NOTE, TIP, IMPORTANT, WARNING, CAUTION)</li>
        </ul>
      </details>
    </div>
  )
}
