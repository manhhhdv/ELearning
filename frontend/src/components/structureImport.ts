/** Phân tích bảng cấu trúc chương trình (dán từ Excel hoặc tải file lên). */

export interface StructureItem {
  level: number
  kind: 'folder' | 'lesson' | 'assignment'
  title: string
  description: string
  contentType: string
  source: string
  durationMinutes: number
}

export interface StructureIssue {
  line: number
  message: string
}

export interface StructureParseResult {
  items: StructureItem[]
  issues: StructureIssue[]
}

export const STRUCTURE_KIND_LABEL: Record<StructureItem['kind'], string> = {
  folder: 'Thư mục',
  lesson: 'Bài học',
  assignment: 'Bài tập',
}

const CONTENT_TYPES = ['video', 'slide', 'document', 'pdf', 'link', 'richtext']

/** Bỏ dấu tiếng Việt để nhận diện tên loại người dùng gõ tự do. */
function normalize(text: string): string {
  return text.trim().toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '').replace(/đ/g, 'd')
}

function parseKind(cell: string): StructureItem['kind'] | null {
  const v = normalize(cell)
  if (v === '') return null
  if (['folder', 'thu muc', 'chuong', 'phan', 'muc'].includes(v)) return 'folder'
  if (['lesson', 'bai hoc', 'bai giang', 'bai'].includes(v)) return 'lesson'
  if (['assignment', 'bai tap', 'bai kiem tra', 'kiem tra', 'quiz', 'de thi'].includes(v)) return 'assignment'
  return null
}

function parseContentType(cell: string): string | null {
  const v = normalize(cell)
  if (v === '') return 'video'
  if (['video', 'phim'].includes(v)) return 'video'
  if (['slide', 'trinh chieu', 'ppt', 'pptx'].includes(v)) return 'slide'
  if (['document', 'tai lieu', 'doc', 'docx', 'van ban'].includes(v)) return 'document'
  if (['pdf'].includes(v)) return 'pdf'
  if (['link', 'lien ket', 'url'].includes(v)) return 'link'
  if (['richtext', 'rich text', 'tu soan', 'bai doc', 'bai doc tu soan', 'markdown'].includes(v)) return 'richtext'
  return CONTENT_TYPES.includes(v) ? v : null
}

/** Dòng đầu chỉ là tiêu đề cột khi ô "Loại" của nó không đọc được thành loại nội dung. */
function isHeaderRow(cells: string[]): boolean {
  const joined = normalize(cells.join(' '))
  return parseKind(cells[1] ?? '') === null && (joined.includes('tieu de') || joined.includes('loai') || joined.includes('cap'))
}

/**
 * Thứ tự cột: Cấp · Loại · Tiêu đề · Loại nội dung · Link/Nguồn · Thời lượng · Mô tả.
 * Cấp 1 nằm ngay dưới thư mục đích; mỗi mục bám vào thư mục gần nhất ở cấp trên nó.
 */
export function parseStructure(raw: string): StructureParseResult {
  const items: StructureItem[] = []
  const issues: StructureIssue[] = []
  // Cấp sâu nhất được phép cho dòng kế tiếp: chỉ thư mục mới mở thêm một cấp.
  let maxLevel = 1

  raw.replace(/\r\n?/g, '\n').split('\n').forEach((line, idx) => {
    const lineNo = idx + 1
    if (line.trim() === '') return

    const cells = line.split('\t').map((c) => c.trim())
    if (idx === 0 && isHeaderRow(cells)) return

    const [levelCell = '', kindCell = '', titleCell = '', contentCell = '', sourceCell = '', durationCell = '', descCell = ''] = cells

    const kind = parseKind(kindCell)
    if (kind === null) {
      issues.push({ line: lineNo, message: `không hiểu loại “${kindCell}” (dùng Thư mục / Bài học / Bài tập)` })
      return
    }
    const title = titleCell
    if (title === '') {
      issues.push({ line: lineNo, message: 'thiếu tiêu đề' })
      return
    }

    let level = levelCell === '' ? 1 : Number(levelCell.replace(',', '.'))
    if (!Number.isInteger(level) || level < 1) {
      issues.push({ line: lineNo, message: `cấp “${levelCell}” phải là số nguyên từ 1 trở lên` })
      return
    }
    if (level > maxLevel) {
      issues.push({
        line: lineNo,
        message: `cấp ${level} bị nhảy bậc — dòng này chỉ có thể ở cấp ${maxLevel} trở xuống`,
      })
      return
    }

    const contentType = kind === 'lesson' ? parseContentType(contentCell) : 'video'
    if (contentType === null) {
      issues.push({ line: lineNo, message: `loại nội dung “${contentCell}” không hợp lệ (video / slide / document / pdf / link / richtext)` })
      return
    }

    let duration = 0
    if (kind === 'lesson' && durationCell !== '') {
      duration = Number(durationCell.replace(',', '.'))
      if (!Number.isFinite(duration) || duration < 0) {
        issues.push({ line: lineNo, message: `thời lượng “${durationCell}” không hợp lệ` })
        return
      }
      duration = Math.round(duration)
    }

    items.push({
      level,
      kind,
      title,
      description: descCell,
      contentType,
      source: kind === 'lesson' ? sourceCell : '',
      durationMinutes: duration,
    })
    // Chỉ thư mục mới nhận con, nên chỉ thư mục mới mở thêm một cấp cho dòng sau.
    maxLevel = kind === 'folder' ? level + 1 : level
  })

  return { items, issues }
}

export const SAMPLE_STRUCTURE = `Cấp\tLoại\tTiêu đề\tLoại nội dung\tLink/Nguồn\tThời lượng\tMô tả
1\tThư mục\tChương 1. Nhập môn Big Data\t\t\t\tTổng quan khái niệm
2\tBài học\tBài 1. Big Data là gì\tvideo\thttps://drive.google.com/file/d/ABC123/view\t25\t
2\tBài học\tBài 2. Hệ sinh thái Hadoop\tslide\thttps://drive.google.com/file/d/DEF456/view\t30\t
2\tBài tập\tKiểm tra chương 1\t\t\t\tNhập câu hỏi sau khi tạo
1\tThư mục\tChương 2. Xử lý dữ liệu\t\t\t\t
2\tBài học\tBài 3. Apache Spark\tvideo\t\t40\t`
