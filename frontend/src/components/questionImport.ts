import type { QuestionInput } from '../api/client'
import type { QuestionType } from '../api/types'

export type ImportFormat = 'text' | 'table'

export interface ParseIssue {
  line: number
  message: string
}

export interface ParseResult {
  questions: QuestionInput[]
  issues: ParseIssue[]
}

/** Suy ra loại câu hỏi từ số phương án và số đáp án đúng. */
function inferType(optionCount: number, correctCount: number): QuestionType {
  if (optionCount === 0) return 'essay'
  return correctCount > 1 ? 'multi_choice' : 'single_choice'
}

/** Tách phần "[2đ]" hoặc "[2 điểm]" ở cuối dòng đề bài. */
function extractPoints(line: string): { prompt: string; points: number } {
  const m = line.match(/\[\s*([\d.,]+)\s*(?:đ|điểm|d|point|pts?)?\s*\]\s*$/i)
  if (!m) return { prompt: line.trim(), points: 1 }
  const points = Number(m[1].replace(',', '.'))
  return {
    prompt: line.slice(0, m.index).trim(),
    points: Number.isFinite(points) && points > 0 ? points : 1,
  }
}

/**
 * Định dạng văn bản: mỗi câu là một khối, ngăn cách bằng dòng trống.
 *
 *   [MÃ] Nội dung câu hỏi [2đ]
 *   * Phương án đúng
 *   - Phương án sai
 *   > Giải thích
 *
 * Không có dòng phương án nào thì câu đó là tự luận.
 */
function parseText(raw: string): ParseResult {
  const questions: QuestionInput[] = []
  const issues: ParseIssue[] = []

  const lines = raw.replace(/\r\n?/g, '\n').split('\n')
  let block: { line: number; text: string }[] = []

  const flush = () => {
    if (block.length === 0) return
    const startLine = block[0].line

    const head = block[0].text
    const codeMatch = head.match(/^\[([^\]]{1,20})\]\s*/)
    const code = codeMatch ? codeMatch[1].trim() : ''
    const { prompt, points } = extractPoints(codeMatch ? head.slice(codeMatch[0].length) : head)

    const options: { content: string; isCorrect: boolean }[] = []
    let explanation = ''

    for (const { text } of block.slice(1)) {
      if (text.startsWith('>')) {
        explanation = (explanation ? explanation + '\n' : '') + text.slice(1).trim()
      } else if (text.startsWith('*') || text.startsWith('+')) {
        options.push({ content: text.slice(1).trim(), isCorrect: true })
      } else if (text.startsWith('-')) {
        options.push({ content: text.slice(1).trim(), isCorrect: false })
      } else {
        issues.push({
          line: startLine,
          message: `Dòng "${text.slice(0, 30)}" không rõ là phương án hay giải thích — cần bắt đầu bằng *, - hoặc >`,
        })
      }
    }

    block = []

    if (!prompt) {
      issues.push({ line: startLine, message: 'Thiếu nội dung câu hỏi' })
      return
    }

    const filled = options.filter((o) => o.content !== '')
    const correct = filled.filter((o) => o.isCorrect).length

    if (filled.length > 0) {
      if (filled.length < 2) {
        issues.push({ line: startLine, message: 'Câu trắc nghiệm cần ít nhất 2 phương án' })
        return
      }
      if (correct === 0) {
        issues.push({ line: startLine, message: 'Chưa đánh dấu phương án đúng (dùng dấu * ở đầu dòng)' })
        return
      }
    }

    questions.push({
      code,
      type: inferType(filled.length, correct),
      prompt,
      points,
      explanation,
      options: filled,
    })
  }

  lines.forEach((rawLine, i) => {
    const text = rawLine.trim()
    if (text === '') flush()
    else block.push({ line: i + 1, text })
  })
  flush()

  return { questions, issues }
}

/**
 * Định dạng bảng dán từ Excel/Google Sheets, các cột cách nhau bằng Tab:
 *
 *   Mã | Nội dung | Điểm | PA1 | PA2 | PA3 | PA4 | Đáp án đúng | Giải thích
 *
 * Cột "Đáp án đúng" ghi số thứ tự phương án (1,3) hoặc chữ cái (A,C).
 * Bỏ trống hết các cột phương án thì câu đó là tự luận.
 */
function parseTable(raw: string): ParseResult {
  const questions: QuestionInput[] = []
  const issues: ParseIssue[] = []

  const rows = raw.replace(/\r\n?/g, '\n').split('\n').filter((r) => r.trim() !== '')

  rows.forEach((row, index) => {
    const line = index + 1
    // Chấp nhận Tab (dán từ bảng tính) hoặc dấu chấm phẩy (CSV kiểu Việt Nam).
    const cells = (row.includes('\t') ? row.split('\t') : row.split(';')).map((c) => c.trim())

    // Bỏ qua dòng tiêu đề.
    if (index === 0 && /^(mã|ma|code|stt)$/i.test(cells[0] ?? '')) return

    const [code = '', prompt = '', pointsRaw = '', o1 = '', o2 = '', o3 = '', o4 = '', answerRaw = '', explanation = ''] = cells

    if (!prompt) {
      issues.push({ line, message: 'Thiếu nội dung câu hỏi ở cột thứ 2' })
      return
    }

    const points = Number((pointsRaw || '1').replace(',', '.'))
    const optionTexts = [o1, o2, o3, o4].filter((o) => o !== '')

    if (optionTexts.length === 0) {
      questions.push({
        code, type: 'essay', prompt,
        points: Number.isFinite(points) && points > 0 ? points : 1,
        explanation, options: [],
      })
      return
    }

    if (optionTexts.length < 2) {
      issues.push({ line, message: 'Câu trắc nghiệm cần ít nhất 2 phương án' })
      return
    }

    // "1,3" hoặc "A,C" hoặc "A C" đều hiểu được.
    const marks = answerRaw.split(/[,;\s]+/).map((t) => t.trim()).filter(Boolean)
    const correctIdx = new Set<number>()
    for (const mark of marks) {
      const asNumber = Number(mark)
      if (Number.isFinite(asNumber) && asNumber >= 1) correctIdx.add(asNumber - 1)
      else if (/^[a-z]$/i.test(mark)) correctIdx.add(mark.toUpperCase().charCodeAt(0) - 65)
    }

    if (correctIdx.size === 0) {
      issues.push({ line, message: 'Cột đáp án đúng để trống hoặc không đọc được (ghi 1,3 hoặc A,C)' })
      return
    }
    if ([...correctIdx].some((i) => i >= optionTexts.length)) {
      issues.push({ line, message: 'Đáp án đúng trỏ tới phương án không tồn tại' })
      return
    }

    questions.push({
      code,
      type: inferType(optionTexts.length, correctIdx.size),
      prompt,
      points: Number.isFinite(points) && points > 0 ? points : 1,
      explanation,
      options: optionTexts.map((content, i) => ({ content, isCorrect: correctIdx.has(i) })),
    })
  })

  return { questions, issues }
}

export function parseQuestions(raw: string, format: ImportFormat): ParseResult {
  if (raw.trim() === '') return { questions: [], issues: [] }
  return format === 'text' ? parseText(raw) : parseTable(raw)
}

export const SAMPLE_TEXT = `[C10] Nhiệt độ bảo quản lạnh thực phẩm tươi sống là? [2đ]
* 0–5°C
- 10–15°C
- 20–25°C
> Theo QCVN, ngăn mát giữ 0–5°C.

Chọn các bước rửa tay đúng quy trình: [3đ]
* Làm ướt tay
* Xoa xà phòng 20 giây
- Lau tay vào áo

Trình bày quy trình xử lý khi phát hiện thực phẩm hết hạn. [3đ]`

export const SAMPLE_TABLE = `Mã\tNội dung câu hỏi\tĐiểm\tPhương án 1\tPhương án 2\tPhương án 3\tPhương án 4\tĐáp án đúng\tGiải thích
C10\tNhiệt độ bảo quản lạnh là?\t2\t0–5°C\t10–15°C\t20–25°C\t\t1\tTheo QCVN
C11\tChọn các bước rửa tay đúng\t3\tLàm ướt tay\tXoa xà phòng\tLau vào áo\t\t1,2\t
C12\tNêu 3 nguyên tắc bảo quản thực phẩm\t3\t\t\t\t\t\t`
