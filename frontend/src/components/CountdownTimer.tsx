import { useEffect, useState } from 'react'

/** Số giây còn lại tới mốc hết giờ, không âm. */
function secondsLeft(deadline: string): number {
  return Math.max(0, Math.round((new Date(deadline).getTime() - Date.now()) / 1000))
}

function format(total: number): string {
  const m = Math.floor(total / 60)
  const s = total % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

interface Props {
  /** Mốc hết giờ do máy chủ cấp — đồng hồ chỉ hiển thị, việc chặn nộp muộn do máy chủ quyết định. */
  expiresAt: string
  onExpire: () => void
}

export function CountdownTimer({ expiresAt, onExpire }: Props) {
  const [left, setLeft] = useState(() => secondsLeft(expiresAt))

  useEffect(() => {
    setLeft(secondsLeft(expiresAt))
    const timer = setInterval(() => {
      const remaining = secondsLeft(expiresAt)
      setLeft(remaining)
      if (remaining <= 0) {
        clearInterval(timer)
        onExpire()
      }
    }, 1000)
    return () => clearInterval(timer)
  }, [expiresAt, onExpire])

  // Dưới 1 phút chuyển đỏ, dưới 5 phút chuyển hổ phách để học viên chú ý.
  const tone = left <= 60 ? 'pill-red' : left <= 300 ? 'pill-amber' : 'pill-blue'

  return (
    <span className={`pill ${tone} timer`} role="timer" aria-live="off">
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" aria-hidden>
        <circle cx="12" cy="12" r="9" />
        <path d="M12 7v5l3 2" strokeLinecap="round" />
      </svg>
      Còn {format(left)}
    </span>
  )
}
