// Bộ biểu tượng SVG nhỏ dùng chung, kế thừa màu chữ hiện hành.

interface IconProps {
  size?: number
  className?: string
}

function svgProps(size: number, className?: string) {
  return {
    width: size,
    height: size,
    viewBox: '0 0 24 24',
    fill: 'none',
    stroke: 'currentColor',
    strokeWidth: 1.9,
    strokeLinecap: 'round' as const,
    strokeLinejoin: 'round' as const,
    className,
    'aria-hidden': true,
  }
}

export const IconFolder = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2Z" />
  </svg>
)

export const IconVideo = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <rect x="2" y="5" width="14" height="14" rx="2" />
    <path d="m16 11 6-3v8l-6-3z" />
  </svg>
)

export const IconSlide = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <rect x="3" y="4" width="18" height="12" rx="2" />
    <path d="M12 16v4M9 20h6" />
  </svg>
)

export const IconDoc = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8Z" />
    <path d="M14 3v5h5M9 13h6M9 17h4" />
  </svg>
)

export const IconTask = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2" />
    <rect x="9" y="3" width="6" height="4" rx="1" />
    <path d="m9 14 2 2 4-4" />
  </svg>
)

export const IconPlus = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}><path d="M12 5v14M5 12h14" /></svg>
)

export const IconTrash = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M3 6h18M8 6V4a1 1 0 0 1 1-1h6a1 1 0 0 1 1 1v2M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6" />
  </svg>
)

export const IconChevron = ({ size = 12, className }: IconProps) => (
  <svg {...svgProps(size, className)}><path d="m9 6 6 6-6 6" /></svg>
)

export const IconCheck = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}><path d="m5 13 4 4L19 7" /></svg>
)

export const IconUsers = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M16 20v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
    <circle cx="9" cy="7" r="3" />
    <path d="M22 20v-2a4 4 0 0 0-3-3.9M16 4.1a4 4 0 0 1 0 5.8" />
  </svg>
)

export const IconBook = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M4 5a2 2 0 0 1 2-2h13v16H6a2 2 0 0 0-2 2Z" />
    <path d="M4 19a2 2 0 0 1 2-2h13" />
  </svg>
)

export const IconGrade = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M12 3 3 8l9 5 9-5-9-5Z" />
    <path d="M7 10.5V16c0 1.1 2.2 2.5 5 2.5s5-1.4 5-2.5v-5.5" />
  </svg>
)

export const IconLogout = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
  </svg>
)

export const IconSearch = ({ size = 18, className }: IconProps) => (
  <svg {...svgProps(size, className)}><circle cx="11" cy="11" r="7" /><path d="m20 20-3.5-3.5" /></svg>
)

export const IconBell = ({ size = 20, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
    <path d="M13.7 21a2 2 0 0 1-3.4 0" />
  </svg>
)

export const IconHelp = ({ size = 20, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <circle cx="12" cy="12" r="9" />
    <path d="M9.1 9.2a3 3 0 0 1 5.8 1c0 2-3 2.5-3 4" />
    <path d="M12 17.5h.01" />
  </svg>
)

export const IconArrowRight = ({ size = 16, className }: IconProps) => (
  <svg {...svgProps(size, className)}><path d="M5 12h14M13 6l6 6-6 6" /></svg>
)

/* Bộ icon cho thẻ thống kê trên trang Tổng quan: tô màu đặc, không dùng emoji
   vì emoji hiển thị khác nhau tuỳ hệ điều hành và hay bị bạc màu ở kích thước lớn. */

export const IconEnrolled = ({ size = 28 }: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" aria-hidden>
    <rect x="4" y="3" width="13" height="18" rx="2" fill="#2f7bf0" />
    <rect x="7" y="1.6" width="7" height="3.4" rx="1.2" fill="#1b4f9e" />
    <path d="M7.5 10h6M7.5 13.5h6M7.5 17h3.5" stroke="#fff" strokeWidth="1.6" strokeLinecap="round" />
    <circle cx="17.5" cy="17" r="4.5" fill="#f5a524" />
    <path d="m15.6 17 1.4 1.4 2.6-2.8" stroke="#fff" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
)

export const IconFinished = ({ size = 28 }: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" aria-hidden>
    <path d="M6 3v18" stroke="#0f7a45" strokeWidth="2.2" strokeLinecap="round" />
    <path d="M7.6 4.2h11.2l-2.6 4 2.6 4H7.6z" fill="#22b573" />
  </svg>
)

export const IconSubmitted = ({ size = 28 }: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" aria-hidden>
    <path d="M5 4a2 2 0 0 1 2-2h7l5 5v13a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2z" fill="#f0a428" />
    <path d="M14 2v5h5" fill="#c9821a" />
    <path d="M8.5 12.5h7M8.5 16h4.5" stroke="#fff" strokeWidth="1.7" strokeLinecap="round" />
  </svg>
)

export const IconUpload = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
    <path d="M7 9l5-5 5 5M12 4v12" />
  </svg>
)

export const IconChart = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M3 3v18h18" />
    <path d="M7 15v-4M12 15V7M17 15v-6" />
  </svg>
)

export const IconPanelClose = ({ size = 18, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <rect x="3" y="4" width="18" height="16" rx="2" />
    <path d="M9 4v16" />
    <path d="m15 10-2 2 2 2" />
  </svg>
)

export const IconPanelOpen = ({ size = 18, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <rect x="3" y="4" width="18" height="16" rx="2" />
    <path d="M9 4v16" />
    <path d="m13 10 2 2-2 2" />
  </svg>
)

export const IconFullscreen = ({ size = 18, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M4 9V5a1 1 0 0 1 1-1h4M15 4h4a1 1 0 0 1 1 1v4M20 15v4a1 1 0 0 1-1 1h-4M9 20H5a1 1 0 0 1-1-1v-4" />
  </svg>
)

export const IconExitFullscreen = ({ size = 18, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M9 4v4a1 1 0 0 1-1 1H4M20 9h-4a1 1 0 0 1-1-1V4M15 20v-4a1 1 0 0 1 1-1h4M4 15h4a1 1 0 0 1 1 1v4" />
  </svg>
)

export const IconEye = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
    <circle cx="12" cy="12" r="3" />
  </svg>
)

export const IconDashboard = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <rect x="3" y="3" width="7" height="9" rx="1.5" />
    <rect x="14" y="3" width="7" height="5" rx="1.5" />
    <rect x="14" y="12" width="7" height="9" rx="1.5" />
    <rect x="3" y="16" width="7" height="5" rx="1.5" />
  </svg>
)

export const IconSettings = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <circle cx="12" cy="12" r="3" />
    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.6 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.6a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1Z" />
  </svg>
)

export const IconPin = ({ size = 15, className }: IconProps) => (
  <svg {...svgProps(size, className)}>
    <path d="M12 17v5" />
    <path d="M9 3h6l1 6 3 3v2H5v-2l3-3Z" />
  </svg>
)

export const IconGoogle = ({ size = 17 }: IconProps) => (
  <svg width={size} height={size} viewBox="0 0 48 48" aria-hidden>
    <path fill="#4285F4" d="M45.1 24.5c0-1.6-.1-3.2-.4-4.7H24v8.9h11.8c-.5 2.7-2 5-4.3 6.6v5.5h7c4.1-3.8 6.6-9.4 6.6-16.3z" />
    <path fill="#34A853" d="M24 46c5.8 0 10.7-1.9 14.3-5.2l-7-5.5c-1.9 1.3-4.4 2.1-7.3 2.1-5.6 0-10.4-3.8-12.1-8.9H4.7v5.6C8.3 41.3 15.6 46 24 46z" />
    <path fill="#FBBC05" d="M11.9 28.5c-.4-1.3-.7-2.7-.7-4.5s.3-3.2.7-4.5v-5.6H4.7C3.1 17 2.2 20.4 2.2 24s.9 7 2.5 10.1l7.2-5.6z" />
    <path fill="#EA4335" d="M24 10.6c3.2 0 6 1.1 8.2 3.2l6.2-6.2C34.7 4.1 29.8 2 24 2 15.6 2 8.3 6.7 4.7 13.9l7.2 5.6c1.7-5.1 6.5-8.9 12.1-8.9z" />
  </svg>
)
