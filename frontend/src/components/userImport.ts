/** Phân tích bảng danh sách tài khoản (dán từ Excel hoặc tải file lên). */

import type { UserImportItem } from '../api/client'
import type { Role } from '../api/types'

export interface UserImportIssue {
  line: number
  message: string
}

export interface UserImportParseResult {
  items: UserImportItem[]
  issues: UserImportIssue[]
}

/** Bỏ dấu tiếng Việt để nhận diện tên vai trò người dùng gõ tự do. */
function normalize(text: string): string {
  return text.trim().toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '').replace(/đ/g, 'd')
}

function parseRole(cell: string): Role | null {
  const v = normalize(cell)
  if (v === '') return 'student'
  if (['student', 'hoc vien', 'hoc sinh', 'nhan vien', 'hv'].includes(v)) return 'student'
  if (['trainer', 'giang vien', 'giao vien', 'gv'].includes(v)) return 'trainer'
  if (['supervisor', 'giam sat', 'giam sat vien'].includes(v)) return 'supervisor'
  if (['admin', 'quan tri', 'quan tri vien', 'qtv'].includes(v)) return 'admin'
  return null
}

// Kiểm tra tối thiểu ở phía giao diện; máy chủ vẫn kiểm lại trước khi tạo tài khoản.
const EMAIL_RE = /^[^\s@,;]+@[^\s@,;]+\.[^\s@,;.]+$/

/** Cùng quy tắc với auth.ValidatePassword ở máy chủ. */
function passwordIssue(password: string): string | null {
  if (password.length < 8) return 'mật khẩu phải có ít nhất 8 ký tự'
  if (!/[a-zA-Z]/.test(password) || !/\d/.test(password)) return 'mật khẩu phải gồm cả chữ và số'
  return null
}

/** Dòng đầu chỉ là tiêu đề cột khi ô email của nó không phải một địa chỉ hợp lệ. */
function isHeaderRow(cells: string[]): boolean {
  const joined = normalize(cells.join(' '))
  return !EMAIL_RE.test(cells[0] ?? '') && (joined.includes('email') || joined.includes('ho ten') || joined.includes('vai tro'))
}

/**
 * Thứ tự cột: Email · Họ tên · Vai trò · Mật khẩu.
 * Bỏ trống vai trò thì mặc định là Học viên; bỏ trống mật khẩu thì máy chủ tự sinh.
 */
export function parseUsers(raw: string): UserImportParseResult {
  const items: UserImportItem[] = []
  const issues: UserImportIssue[] = []
  // Email đã gặp -> số dòng, để báo trùng ngay trong file trước khi gửi lên.
  const seen = new Map<string, number>()

  raw.replace(/\r\n?/g, '\n').split('\n').forEach((line, idx) => {
    const lineNo = idx + 1
    if (line.trim() === '') return

    const cells = line.split('\t').map((c) => c.trim())
    if (idx === 0 && isHeaderRow(cells)) return

    const [emailCell = '', nameCell = '', roleCell = '', passwordCell = ''] = cells

    const email = emailCell.toLowerCase()
    if (email === '') {
      issues.push({ line: lineNo, message: 'thiếu email' })
      return
    }
    if (!EMAIL_RE.test(email)) {
      issues.push({ line: lineNo, message: `email “${emailCell}” không hợp lệ` })
      return
    }
    const duplicateOf = seen.get(email)
    if (duplicateOf !== undefined) {
      issues.push({ line: lineNo, message: `email “${email}” đã xuất hiện ở dòng ${duplicateOf}` })
      return
    }

    const role = parseRole(roleCell)
    if (role === null) {
      issues.push({
        line: lineNo,
        message: `không hiểu vai trò “${roleCell}” (dùng Học viên / Giảng viên / Giám sát / Quản trị viên)`,
      })
      return
    }

    if (passwordCell !== '') {
      const problem = passwordIssue(passwordCell)
      if (problem !== null) {
        issues.push({ line: lineNo, message: problem })
        return
      }
    }

    seen.set(email, lineNo)
    items.push({ email, fullName: nameCell, role, password: passwordCell })
  })

  return { items, issues }
}

export const SAMPLE_USERS = `Email\tHọ và tên\tVai trò\tMật khẩu
an.nguyen@congty.vn\tNguyễn Văn An\tHọc viên\t
binh.tran@congty.vn\tTrần Thị Bình\tGiảng viên\t
cuong.le@congty.vn\tLê Mạnh Cường\tGiám sát\tMatKhau123
dung.pham@congty.vn\tPhạm Thị Dung\t\t`
