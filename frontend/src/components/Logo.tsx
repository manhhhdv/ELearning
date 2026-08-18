/**
 * Logo Tập Huấn: mũ tốt nghiệp trên nền xanh bo góc.
 * Giữ đồng bộ với public/favicon.svg — sửa hình ở đây thì sửa cả file kia.
 */

interface Props {
  size?: number
  /** Kèm chữ "Tập Huấn" bên phải hình. */
  withWordmark?: boolean
  /** Dùng trên nền tối: chữ chuyển sang trắng. */
  onDark?: boolean
}

export function LogoMark({ size = 36 }: { size?: number }) {
  // id gradient phải khác nhau giữa các bản sao, nếu không bản sau sẽ ăn theo bản đầu.
  const gradientId = `logo-tile-${size}`
  return (
    <svg width={size} height={size} viewBox="0 0 48 48" aria-hidden focusable="false">
      <defs>
        <linearGradient id={gradientId} x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stopColor="#3d94f6" />
          <stop offset="1" stopColor="#0c3f9e" />
        </linearGradient>
      </defs>
      <rect width="48" height="48" rx="11" fill={`url(#${gradientId})`} />
      <path d="M24 10.5 42 18.2 24 25.9 6 18.2Z" fill="#fff" />
      <path
        d="M14.4 21.3v6.1c0 .72.35 1.4.95 1.82C17.9 30.94 20.86 31.8 24 31.8s6.1-.86 8.65-2.57c.6-.42.95-1.1.95-1.82v-6.1L24 27.6Z"
        fill="#cfe4ff"
      />
      <path d="M39.4 19.3v8.4" stroke="#fff" strokeWidth="1.7" strokeLinecap="round" />
      <circle cx="39.4" cy="30.2" r="2.4" fill="#9ecbff" />
    </svg>
  )
}

export function Logo({ size = 36, withWordmark = true, onDark = false }: Props) {
  return (
    <span className="logo">
      <LogoMark size={size} />
      {withWordmark && (
        <span className="logo-text" style={onDark ? { color: '#fff' } : undefined}>
          Tập Huấn
        </span>
      )}
    </span>
  )
}
