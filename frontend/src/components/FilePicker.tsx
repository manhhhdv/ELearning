import { useId, useRef, useState } from 'react'

import { IconUpload } from './icons'

/**
 * Ô chọn tệp.
 *
 * Trình duyệt tự vẽ nút của `input[type=file]` với chữ tiếng Anh ("Choose File",
 * "No file chosen") và CSS không đổi được chữ đó. Nên ở đây input thật bị ẩn,
 * còn nút và tên tệp do mình dựng — vừa ra tiếng Việt, vừa dùng đúng kiểu nút
 * chung của hệ thống.
 */
export function FilePicker({
  accept, onPick, disabled = false, label = 'Chọn tệp', hint,
}: {
  accept?: string
  onPick: (file: File) => void
  disabled?: boolean
  label?: string
  hint?: string
}) {
  const inputId = useId()
  const inputRef = useRef<HTMLInputElement>(null)
  const [name, setName] = useState<string | null>(null)

  return (
    <div className="file-picker">
      <input
        ref={inputRef}
        id={inputId}
        type="file"
        accept={accept}
        disabled={disabled}
        className="file-picker-input"
        onChange={(e) => {
          const file = e.target.files?.[0]
          if (!file) return
          setName(file.name)
          onPick(file)
          // Xoá giá trị để chọn lại đúng tệp vừa rồi vẫn kích hoạt onChange.
          e.target.value = ''
        }}
      />
      {/* Nhãn <label> gắn với input nên bấm vào là mở hộp thoại chọn tệp,
          không cần JavaScript và bàn phím vẫn dùng được. */}
      <label htmlFor={inputId} className={`btn btn-sm ${disabled ? 'is-disabled' : ''}`}>
        <IconUpload /> {label}
      </label>
      <span className="file-picker-name">
        {name ?? <span className="muted">Chưa chọn tệp nào</span>}
      </span>
      {hint && <div className="hint">{hint}</div>}
    </div>
  )
}
