/**
 * Kết xuất nội dung tự soạn: markdown (GFM) + công thức LaTeX + media nhúng.
 *
 * Thứ tự xử lý quan trọng:
 *   1. Che các đoạn mã (``` và `inline`) để ký hiệu $ bên trong không bị coi là công thức.
 *   2. Tách công thức LaTeX ra ngoài, thay bằng ký tự đánh dấu — markdown không được đụng vào TeX.
 *   3. Chuyển markdown sang HTML rồi lọc sạch bằng DOMPurify.
 *   4. Chỉ khi HTML đã sạch mới đắp KaTeX vào các ký tự đánh dấu, và chỉ ở vị trí văn bản
 *      (không bao giờ ở trong thuộc tính) nên HTML của KaTeX không thể phá khung thẻ.
 */
import DOMPurify from 'dompurify'
import katex from 'katex'
import { Marked, type Tokens } from 'marked'

/** Ký tự vùng riêng (private use area) — không xuất hiện trong nội dung người dùng gõ. */
const MARK = '\uE000'
const MATH_TOKEN = /\uE000m(\d+)\uE000/
const MATH_TOKEN_ALL = /\uE000m(\d+)\uE000/g

interface MathPart {
  tex: string
  display: boolean
}

// ---------------------------------------------------------------------------
// Nhận diện media từ đường dẫn
// ---------------------------------------------------------------------------

const VIDEO_EXT = /\.(mp4|webm|ogv|mov|m4v)(\?|#|$)/i
const AUDIO_EXT = /\.(mp3|wav|ogg|oga|m4a|aac|flac)(\?|#|$)/i

/** Các host được phép nhúng iframe; ngoài danh sách này iframe sẽ bị loại bỏ. */
const EMBED_HOSTS = new Set([
  'www.youtube.com',
  'youtube.com',
  'www.youtube-nocookie.com',
  'player.vimeo.com',
  'drive.google.com',
  'docs.google.com',
])

export type MediaKind = 'image' | 'video' | 'audio' | 'embed'

/** Đổi link chia sẻ quen thuộc (YouTube, Vimeo, Drive) thành URL nhúng được. */
export function detectMedia(rawUrl: string): { kind: MediaKind; src: string } {
  const url = rawUrl.trim()
  let parsed: URL | null = null
  try {
    parsed = new URL(url, window.location.origin)
  } catch {
    parsed = null
  }

  if (parsed && (parsed.protocol === 'http:' || parsed.protocol === 'https:')) {
    const host = parsed.hostname.replace(/^m\./, '')
    const path = parsed.pathname

    if (host === 'youtu.be' && path.length > 1) {
      return { kind: 'embed', src: `https://www.youtube-nocookie.com/embed/${path.slice(1)}` }
    }
    if (host.endsWith('youtube.com') || host.endsWith('youtube-nocookie.com')) {
      const id = parsed.searchParams.get('v') ?? path.match(/\/(?:embed|shorts|live)\/([^/?#]+)/)?.[1]
      if (id) return { kind: 'embed', src: `https://www.youtube-nocookie.com/embed/${id}` }
    }
    if (host.endsWith('vimeo.com')) {
      const id = path.match(/\/(\d+)/)?.[1]
      if (id) return { kind: 'embed', src: `https://player.vimeo.com/video/${id}` }
    }
    if (host === 'drive.google.com') {
      const id = path.match(/\/file\/d\/([^/?#]+)/)?.[1] ?? parsed.searchParams.get('id')
      // Ảnh trên Drive hiển thị đẹp hơn qua thumbnail; file khác dùng trình xem của Drive.
      if (id) return { kind: 'embed', src: `https://drive.google.com/file/d/${id}/preview` }
    }
    if (host === 'docs.google.com') {
      const m = path.match(/\/(presentation|document|spreadsheets)\/d\/([^/?#]+)/)
      if (m) {
        const suffix = m[1] === 'presentation' ? '/embed' : '/preview'
        return { kind: 'embed', src: `https://docs.google.com/${m[1]}/d/${m[2]}${suffix}` }
      }
    }
  }

  if (VIDEO_EXT.test(url)) return { kind: 'video', src: url }
  if (AUDIO_EXT.test(url)) return { kind: 'audio', src: url }
  return { kind: 'image', src: url }
}

// ---------------------------------------------------------------------------
// Markdown
// ---------------------------------------------------------------------------

const escapeAttr = (value: string) =>
  value.replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

/** Chỉ giữ link http(s) hoặc đường dẫn nội bộ; phần còn lại (javascript:, data:…) bị bỏ. */
function safeUrl(url: string): string {
  const value = url.trim()
  if (/^(https?:)?\/\//i.test(value) || value.startsWith('/')) return value
  return ''
}

/** Cú pháp ảnh của markdown dùng chung cho mọi loại media, tự đoán theo đuôi file / host. */
function renderMedia(href: string, alt: string, title: string): string {
  const { kind, src } = detectMedia(href)
  const url = safeUrl(src)
  if (!url) return escapeAttr(alt || href)

  const label = escapeAttr(alt)
  const caption = title ? `<figcaption>${escapeAttr(title)}</figcaption>` : ''

  switch (kind) {
    case 'embed':
      return (
        `<figure class="rich-media rich-embed">` +
        `<iframe src="${escapeAttr(url)}" title="${label || 'Nội dung nhúng'}" loading="lazy"` +
        ` allow="accelerometer; autoplay; clipboard-write; encrypted-media; picture-in-picture; fullscreen"` +
        ` allowfullscreen></iframe>${caption}</figure>`
      )
    case 'video':
      return (
        `<figure class="rich-media"><video src="${escapeAttr(url)}" controls preload="metadata"` +
        ` playsinline${label ? ` aria-label="${label}"` : ''}></video>${caption}</figure>`
      )
    case 'audio':
      return (
        `<figure class="rich-media"><audio src="${escapeAttr(url)}" controls preload="metadata"` +
        `${label ? ` aria-label="${label}"` : ''}></audio>${caption}</figure>`
      )
    default:
      return (
        `<figure class="rich-media"><img src="${escapeAttr(url)}" alt="${label}" />` +
        `${caption}</figure>`
      )
  }
}

const markdown = new Marked({ gfm: true, breaks: true, async: false })
markdown.use({
  renderer: {
    image({ href, title, text }: Tokens.Image) {
      return renderMedia(href, text ?? '', title ?? '')
    },
  },
})

// ---------------------------------------------------------------------------
// Che đoạn mã và tách công thức
// ---------------------------------------------------------------------------

const FENCED = /(^|\n)[ \t]*(`{3,}|~{3,})[\s\S]*?(?:\n[ \t]*\2[ \t]*(?=\n|$)|$)/g
const INLINE_CODE = /(`+)[\s\S]*?\1/g

function maskCode(source: string): { text: string; codes: string[] } {
  const codes: string[] = []
  const keep = (match: string, lead = '') => {
    codes.push(match.slice(lead.length))
    return `${lead}${MARK}c${codes.length - 1}${MARK}`
  }
  return {
    text: source
      .replace(FENCED, (match, lead: string) => keep(match, lead))
      .replace(INLINE_CODE, (match) => keep(match)),
    codes,
  }
}

const restoreCode = (text: string, codes: string[]) =>
  text.replace(/\uE000c(\d+)\uE000/g, (_, i: string) => codes[Number(i)] ?? '')

// $$…$$ và \[…\] là công thức khối; $…$ và \(…\) là công thức trong dòng.
// Với $…$ phải chặn hai đầu bằng ký tự không phải khoảng trắng để "giá $5 và $7" không thành công thức.
const MATH_PATTERNS: { re: RegExp; display: boolean }[] = [
  { re: /(?<!\\)\$\$([\s\S]+?)\$\$/g, display: true },
  { re: /(?<!\\)\\\[([\s\S]+?)\\\]/g, display: true },
  { re: /(?<!\\)\\\(([\s\S]+?)\\\)/g, display: false },
  { re: /(?<!\\)\$(?!\s)((?:[^$\\\n]|\\.)+?)(?<![\s\\])\$(?!\d)/g, display: false },
]

function extractMath(source: string): { text: string; math: MathPart[] } {
  const math: MathPart[] = []
  let text = source
  for (const { re, display } of MATH_PATTERNS) {
    text = text.replace(re, (_, tex: string) => {
      math.push({ tex: tex.trim(), display })
      return `${MARK}m${math.length - 1}${MARK}`
    })
  }
  return { text, math }
}

// ---------------------------------------------------------------------------
// Lọc HTML
// ---------------------------------------------------------------------------

let hooksReady = false

function installHooks() {
  if (hooksReady) return
  hooksReady = true

  DOMPurify.addHook('afterSanitizeAttributes', (node) => {
    if (!(node instanceof Element)) return

    // Danh sách việc cần làm của markdown sinh ra ô tích; chỉ giữ đúng dạng đó và luôn khoá lại.
    if (node.tagName === 'INPUT') {
      if (node.getAttribute('type') !== 'checkbox') {
        node.remove()
        return
      }
      node.setAttribute('disabled', '')
      node.removeAttribute('name')
      node.removeAttribute('value')
      return
    }

    if (node.tagName === 'IFRAME') {
      let host = ''
      try {
        host = new URL(node.getAttribute('src') ?? '', window.location.origin).hostname
      } catch {
        host = ''
      }
      if (!EMBED_HOSTS.has(host)) node.remove()
      return
    }

    // Link ra ngoài luôn mở tab mới và không rò rỉ phiên làm việc.
    if (node.tagName === 'A' && /^https?:/i.test(node.getAttribute('href') ?? '')) {
      node.setAttribute('target', '_blank')
      node.setAttribute('rel', 'noopener noreferrer')
    }
  })
}

const SANITIZE_CONFIG = {
  ADD_TAGS: ['iframe'],
  ADD_ATTR: ['allow', 'allowfullscreen', 'loading', 'controls', 'playsinline', 'preload', 'target'],
  FORBID_TAGS: ['style', 'form', 'textarea', 'select', 'button'],
}

// ---------------------------------------------------------------------------
// Đắp KaTeX và tinh chỉnh HTML sau khi đã lọc
// ---------------------------------------------------------------------------

function paintMath(root: DocumentFragment, math: MathPart[]) {
  if (!math.length) return

  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const hits: Text[] = []
  for (let node = walker.nextNode(); node; node = walker.nextNode()) {
    if (MATH_TOKEN.test(node.nodeValue ?? '')) hits.push(node as Text)
  }

  for (const node of hits) {
    const pieces = (node.nodeValue ?? '').split(MATH_TOKEN_ALL)
    const frag = document.createDocumentFragment()
    pieces.forEach((piece, index) => {
      // split với nhóm bắt: vị trí lẻ là chỉ số công thức, vị trí chẵn là văn bản thường.
      if (index % 2 === 0) {
        if (piece) frag.append(piece)
        return
      }
      const part = math[Number(piece)]
      if (!part) return
      const holder = document.createElement('span')
      holder.className = part.display ? 'rich-math rich-math-block' : 'rich-math'
      try {
        katex.render(part.tex, holder, {
          displayMode: part.display,
          throwOnError: false,
          strict: 'ignore',
          output: 'htmlAndMathml',
          trust: false,
        })
      } catch {
        holder.classList.add('rich-math-error')
        holder.textContent = part.display ? `$$${part.tex}$$` : `$${part.tex}$`
      }
      frag.append(holder)
    })
    node.replaceWith(frag)
  }
}

const ALERT_LABEL: Record<string, string> = {
  NOTE: 'Ghi chú',
  TIP: 'Mẹo',
  IMPORTANT: 'Quan trọng',
  WARNING: 'Lưu ý',
  CAUTION: 'Cảnh báo',
}

/** Bảng rộng cần khung cuộn riêng, và trích dẫn dạng `> [!NOTE]` trở thành hộp nhấn mạnh. */
function enhance(root: DocumentFragment) {
  for (const table of Array.from(root.querySelectorAll('table'))) {
    const scroller = document.createElement('div')
    scroller.className = 'rich-scroll'
    table.replaceWith(scroller)
    scroller.append(table)
  }

  for (const quote of Array.from(root.querySelectorAll('blockquote'))) {
    const first = quote.firstElementChild
    const match = first?.textContent?.match(/^\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*/i)
    if (!first || !match) continue

    const key = match[1].toUpperCase()
    quote.classList.add('rich-alert', `rich-alert-${key.toLowerCase()}`)
    // Cắt đúng phần "[!NOTE]" ở đầu, giữ nguyên phần định dạng còn lại của dòng.
    trimLeading(first, match[0].length)
    // Nhãn nằm riêng một dòng nên phần còn lại mở đầu bằng <br>; bỏ đi cho khỏi hụt một dòng trống.
    while (first.firstChild && (first.firstChild.nodeName === 'BR'
      || (first.firstChild.nodeType === Node.TEXT_NODE && !first.firstChild.nodeValue?.trim()))) {
      first.firstChild.remove()
    }
    if (!first.textContent?.trim() && !first.querySelector('img, video, audio, iframe')) first.remove()
    const heading = document.createElement('p')
    heading.className = 'rich-alert-title'
    heading.textContent = ALERT_LABEL[key]
    quote.prepend(heading)
  }
}

/** Bỏ `count` ký tự đầu tiên của một phần tử, đi xuyên qua các thẻ con nếu cần. */
function trimLeading(element: Element, count: number) {
  const walker = document.createTreeWalker(element, NodeFilter.SHOW_TEXT)
  let left = count
  while (left > 0) {
    const node = walker.nextNode() as Text | null
    if (!node) return
    const take = Math.min(left, node.data.length)
    node.deleteData(0, take)
    left -= take
  }
}

// ---------------------------------------------------------------------------

/** Chuyển nội dung tự soạn thành HTML an toàn để gắn vào trang. */
export function renderRichText(source: string): string {
  if (!source || !source.trim()) return ''
  installHooks()

  const masked = maskCode(source)
  const extracted = extractMath(masked.text)
  const html = markdown.parse(restoreCode(extracted.text, masked.codes)) as string

  const template = document.createElement('template')
  template.innerHTML = DOMPurify.sanitize(html, SANITIZE_CONFIG)
  paintMath(template.content, extracted.math)
  enhance(template.content)
  return template.innerHTML
}
