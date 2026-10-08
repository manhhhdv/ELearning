package docx

// Chuyển nội dung bài giảng tự soạn (Markdown GFM + công thức LaTeX + media
// nhúng, xem frontend/src/components/richtext.ts) thành đoạn văn Word.
//
// Công thức được dựng thành phương trình Word (xem math.go); công thức dùng
// lệnh chưa hỗ trợ thì giữ nguyên dạng $…$ để chuyển tay bằng MathType.

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var markdownParser = goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser()

// Markdown thêm nội dung Markdown vào vùng chứa. Giống trình hiển thị ở
// frontend, một lần xuống dòng đơn cũng là một lần ngắt dòng.
func (c *Container) Markdown(src string) {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	masked, math := extractMath(src)
	source := []byte(masked)
	w := &mdWriter{src: source, math: math}
	w.blocks(c, markdownParser.Parse(text.NewReader(source)), mdBlockCtx{level: -1})
}

// ---------------------------------------------------------------------------
// Tách công thức
// ---------------------------------------------------------------------------

// Công thức được thay bằng ký tự đánh dấu thuộc vùng riêng (private use area)
// trước khi phân tích Markdown, để các ký hiệu _ * \ bên trong TeX không bị
// hiểu thành in nghiêng hay ký tự thoát.
const (
	mathOpen  = ''
	mathClose = ''
)

var mathToken = regexp.MustCompile("(\\d+)")

type mathPart struct {
	// raw là nguyên văn công thức kể cả dấu phân cách; tex là phần bên trong.
	raw     string
	tex     string
	display bool
}

func mathPlaceholder(i int) string {
	return string(mathOpen) + strconv.Itoa(i) + string(mathClose)
}

// extractMath tách công thức ($$…$$, \[…\], \(…\), $…$) ra khỏi Markdown, bỏ
// qua khối mã và mã trong dòng. Quy tắc nhận diện $…$ giống richtext.ts: hai
// đầu không phải khoảng trắng, không xuống dòng, và $ đóng không đứng liền
// trước chữ số — để "giá $5 và $7" không thành công thức.
func extractMath(src string) (string, []mathPart) {
	var out strings.Builder
	var math []mathPart
	var chunk strings.Builder
	flush := func() {
		out.WriteString(extractInlineMath(chunk.String(), &math))
		chunk.Reset()
	}

	lines := strings.SplitAfter(src, "\n")
	fence := ""
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if fence != "" {
			out.WriteString(line)
			if strings.TrimRight(trimmed, " \t\n") == fence {
				fence = ""
			}
			continue
		}
		if f := fenceMarker(trimmed); f != "" {
			flush()
			fence = f
			out.WriteString(line)
			continue
		}
		chunk.WriteString(line)
	}
	flush()
	return out.String(), math
}

// fenceMarker trả về chuỗi ``` hoặc ~~~ mở đầu một khối mã, rỗng nếu không phải.
func fenceMarker(line string) string {
	for _, ch := range []byte{'`', '~'} {
		n := 0
		for n < len(line) && line[n] == ch {
			n++
		}
		if n >= 3 {
			return line[:n]
		}
	}
	return ""
}

func extractInlineMath(s string, math *[]mathPart) string {
	var b strings.Builder
	// keep ghi lại công thức; delim là độ dài dấu phân cách mỗi đầu ($ = 1, $$ \[ \( = 2).
	keep := func(raw string, display bool, delim int) {
		*math = append(*math, mathPart{raw: raw, tex: raw[delim : len(raw)-delim], display: display})
		b.WriteString(mathPlaceholder(len(*math) - 1))
	}

	for i := 0; i < len(s); {
		switch {
		case s[i] == '\\' && i+1 < len(s):
			closer := ""
			switch s[i+1] {
			case '[':
				closer = `\]`
			case '(':
				closer = `\)`
			}
			if closer != "" {
				if j := strings.Index(s[i+2:], closer); j > 0 {
					end := i + 2 + j + 2
					keep(s[i:end], closer == `\]`, 2)
					i = end
					continue
				}
			}
			// Ký tự thoát (VD: \$) giữ nguyên cho Markdown xử lý.
			b.WriteString(s[i : i+2])
			i += 2

		case s[i] == '`':
			n := 0
			for i+n < len(s) && s[i+n] == '`' {
				n++
			}
			if end := closingBackticks(s, i+n, n); end > 0 {
				b.WriteString(s[i:end])
				i = end
			} else {
				b.WriteString(s[i : i+n])
				i += n
			}

		case strings.HasPrefix(s[i:], "$$"):
			if j := strings.Index(s[i+2:], "$$"); j > 0 {
				end := i + 2 + j + 2
				keep(s[i:end], true, 2)
				i = end
				continue
			}
			b.WriteString("$$")
			i += 2

		case s[i] == '$':
			if end := inlineMathEnd(s, i); end > 0 {
				keep(s[i:end], false, 1)
				i = end
				continue
			}
			b.WriteByte('$')
			i++

		default:
			b.WriteByte(s[i])
			i++
		}
	}
	return b.String()
}

// closingBackticks tìm dãy đúng n dấu ` đóng mã trong dòng, trả về vị trí ngay
// sau dãy đó hoặc -1 nếu không có.
func closingBackticks(s string, from, n int) int {
	for i := from; i < len(s); {
		if s[i] != '`' {
			i++
			continue
		}
		run := 0
		for i+run < len(s) && s[i+run] == '`' {
			run++
		}
		if run == n {
			return i + run
		}
		i += run
	}
	return -1
}

// inlineMathEnd trả về vị trí ngay sau dấu $ đóng của công thức bắt đầu ở
// start, hoặc -1 nếu đây không phải công thức.
func inlineMathEnd(s string, start int) int {
	i := start + 1
	if i >= len(s) || isSpace(s[i]) {
		return -1
	}
	for i < len(s) {
		switch s[i] {
		case '\n':
			return -1
		case '\\':
			if i+1 >= len(s) || s[i+1] == '\n' {
				return -1
			}
			i += 2
		case '$':
			if i == start+1 || isSpace(s[i-1]) || s[i-1] == '\\' {
				return -1
			}
			if i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9' {
				return -1
			}
			return i + 1
		default:
			i++
		}
	}
	return -1
}

func isSpace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

// ---------------------------------------------------------------------------
// Duyệt cây Markdown
// ---------------------------------------------------------------------------

type mdWriter struct {
	src  []byte
	math []mathPart
}

// mdBlockCtx là ngữ cảnh khối hiện tại: đang trong danh sách nào, cấp mấy,
// có phải trích dẫn không.
type mdBlockCtx struct {
	list  int
	level int
	// Đoạn đầu tiên của một mục danh sách mang số thứ tự; các đoạn sau chỉ lùi lề.
	numbered *bool
	quote    bool
}

func (w *mdWriter) blocks(c *Container, parent ast.Node, ctx mdBlockCtx) {
	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		w.block(c, n, ctx)
	}
}

func (w *mdWriter) block(c *Container, n ast.Node, ctx mdBlockCtx) {
	switch n := n.(type) {
	case *ast.Heading:
		style := StyleHeading3
		switch n.Level {
		case 1:
			style = StyleHeading1
		case 2:
			style = StyleHeading2
		}
		c.Add(Para{Style: style, Runs: w.inlines(n)})

	case *ast.Paragraph, *ast.TextBlock:
		c.Add(w.para(n, ctx))

	case *ast.List:
		list := c.doc.NewList(n.IsOrdered(), n.Start)
		for item := n.FirstChild(); item != nil; item = item.NextSibling() {
			numbered := false
			w.blocks(c, item, mdBlockCtx{list: list, level: ctx.level + 1, numbered: &numbered, quote: ctx.quote})
		}

	case *ast.Blockquote:
		ctx.quote = true
		// "> [!NOTE]" là hộp nhấn mạnh giống ở frontend: thay nhãn bằng tiêu đề in đậm.
		if first, ok := n.FirstChild().(*ast.Paragraph); ok {
			if label, rest, ok := splitAlert(w.inlines(first)); ok {
				c.Add(Para{Style: StyleQuote, Indent: w.indent(ctx), KeepNext: true, Runs: []Run{{Text: label, Bold: true}}})
				if len(rest) > 0 {
					c.Add(Para{Style: StyleQuote, Indent: w.indent(ctx), Runs: rest})
				}
				for m := first.NextSibling(); m != nil; m = m.NextSibling() {
					w.block(c, m, ctx)
				}
				return
			}
		}
		w.blocks(c, n, ctx)

	case *ast.FencedCodeBlock, *ast.CodeBlock:
		var code strings.Builder
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			code.Write(seg.Value(w.src))
		}
		c.Add(Para{Style: StyleCode, Indent: w.indent(ctx), Runs: []Run{{Text: w.restoreRaw(strings.TrimRight(code.String(), "\n"))}}})

	case *ast.ThematicBreak:
		c.Add(Para{Style: StyleRule})

	case *east.Table:
		w.table(c, n)

	case *ast.HTMLBlock:
		var raw strings.Builder
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			raw.Write(seg.Value(w.src))
		}
		if txt := strings.TrimSpace(htmlTag.ReplaceAllString(raw.String(), "")); txt != "" {
			c.Add(Para{Indent: w.indent(ctx), Runs: []Run{{Text: w.restoreMath(txt)}}})
		}

	case *ast.LinkReferenceDefinition:
		// Định nghĩa liên kết không hiển thị thành chữ.

	default:
		w.blocks(c, n, ctx)
	}
}

// indent là lề trái của nội dung nằm trong một mục danh sách, thẳng hàng với
// chữ (không phải ký hiệu đầu dòng) của mục đó.
func (w *mdWriter) indent(ctx mdBlockCtx) int {
	if ctx.list == 0 {
		return 0
	}
	return 567 * (ctx.level + 1)
}

func (w *mdWriter) para(n ast.Node, ctx mdBlockCtx) Para {
	p := Para{Runs: w.inlines(n)}
	if ctx.quote {
		p.Style = StyleQuote
	}
	switch {
	case ctx.list > 0 && ctx.numbered != nil && !*ctx.numbered:
		p.List, p.Level = ctx.list, ctx.level
		*ctx.numbered = true
	case ctx.list > 0:
		p.Indent = w.indent(ctx)
	}
	if w.isDisplayMath(n) {
		p.Align = AlignCenter
	}
	return p
}

// isDisplayMath cho biết đoạn văn chỉ gồm đúng một công thức khối.
func (w *mdWriter) isDisplayMath(n ast.Node) bool {
	t, ok := n.FirstChild().(*ast.Text)
	if !ok || t.NextSibling() != nil {
		return false
	}
	m := mathToken.FindSubmatch(t.Value(w.src))
	if m == nil || len(m[0]) != len(strings.TrimSpace(string(t.Value(w.src)))) {
		return false
	}
	i, _ := strconv.Atoi(string(m[1]))
	return i < len(w.math) && w.math[i].display
}

func (w *mdWriter) table(c *Container, n *east.Table) {
	cols := len(n.Alignments)
	if cols == 0 {
		return
	}
	weights := make([]int, cols)
	for i := range weights {
		weights[i] = 1
	}
	t := c.NewTable(weights...)
	for row := n.FirstChild(); row != nil; row = row.NextSibling() {
		_, header := row.(*east.TableHeader)
		if header {
			t.HeaderRows++
		}
		cells := t.Row()
		i := 0
		for cell := row.FirstChild(); cell != nil && i < cols; cell = cell.NextSibling() {
			runs := w.inlines(cell)
			if header {
				for j := range runs {
					runs[j].Bold = true
				}
			}
			align := AlignLeft
			if tc, ok := cell.(*east.TableCell); ok {
				switch tc.Alignment {
				case east.AlignCenter:
					align = AlignCenter
				case east.AlignRight:
					align = AlignRight
				}
			}
			cells[i].Add(Para{Align: align, Runs: runs})
			i++
		}
	}
	c.AddTable(t)
}

// ---------------------------------------------------------------------------
// Định dạng trong dòng
// ---------------------------------------------------------------------------

var htmlTag = regexp.MustCompile(`<[^>]*>`)

func (w *mdWriter) inlines(n ast.Node) []Run {
	var runs []Run
	w.inlineChildren(n, Run{}, &runs)
	return mergeRuns(runs)
}

func (w *mdWriter) inlineChildren(parent ast.Node, style Run, out *[]Run) {
	// Thẻ HTML trong dòng (<u>, <sup>…) bật/tắt định dạng cho các nút anh em
	// phía sau nó, nên kiểu chữ được mang theo qua vòng lặp.
	cur := style
	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		if raw, ok := n.(*ast.RawHTML); ok {
			cur = w.applyTag(raw, cur, out)
			continue
		}
		w.inline(n, cur, out)
	}
}

func (w *mdWriter) emit(out *[]Run, style Run, s string) {
	if s == "" {
		return
	}
	r := style
	r.Text = s
	*out = append(*out, r)
}

// emitText ghi chữ thường, tách các ký tự đánh dấu công thức thành run phương trình.
func (w *mdWriter) emitText(out *[]Run, style Run, s string) {
	for s != "" {
		loc := mathToken.FindStringSubmatchIndex(s)
		if loc == nil {
			w.emit(out, style, s)
			return
		}
		w.emit(out, style, s[:loc[0]])
		if i, err := strconv.Atoi(s[loc[2]:loc[3]]); err == nil && i < len(w.math) {
			m := w.math[i]
			*out = append(*out, Run{Math: m.tex, Display: m.display})
		}
		s = s[loc[1]:]
	}
}

func (w *mdWriter) inline(n ast.Node, style Run, out *[]Run) {
	switch n := n.(type) {
	case *ast.Text:
		v := n.Value(w.src)
		if !n.IsRaw() {
			v = util.UnescapePunctuations(util.ResolveNumericReferences(util.ResolveEntityNames(v)))
		}
		w.emitText(out, style, string(v))
		if n.SoftLineBreak() || n.HardLineBreak() {
			w.emit(out, Run{}, "\n")
		}

	case *ast.String:
		w.emitText(out, style, string(n.Value))

	case *ast.CodeSpan:
		var code strings.Builder
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				code.Write(t.Value(w.src))
			}
		}
		style.Code = true
		w.emit(out, style, w.restoreRaw(strings.ReplaceAll(code.String(), "\n", " ")))

	case *ast.Emphasis:
		if n.Level >= 2 {
			style.Bold = true
		} else {
			style.Italic = true
		}
		w.inlineChildren(n, style, out)

	case *east.Strikethrough:
		style.Strike = true
		w.inlineChildren(n, style, out)

	case *ast.Link:
		style.Link = w.restoreRaw(string(n.Destination))
		w.inlineChildren(n, style, out)

	case *ast.AutoLink:
		url := string(n.URL(w.src))
		if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(url), "mailto:") {
			url = "mailto:" + url
		}
		style.Link = url
		w.emit(out, style, string(n.Label(w.src)))

	case *ast.Image:
		// Ảnh, video, âm thanh đều dùng cú pháp ảnh của Markdown (xem richtext.ts).
		// Tệp Word chỉ ghi lại liên kết kèm nhãn để giáo viên mở xem hoặc tự chèn.
		var alt []Run
		w.inlineChildren(n, Run{}, &alt)
		label := mediaLabel(string(n.Destination))
		if caption := strings.TrimSpace(runsText(alt)); caption != "" {
			label += ": " + caption
		}
		style.Italic = true
		style.Link = w.restoreRaw(string(n.Destination))
		w.emit(out, style, "["+label+"]")

	case *east.TaskCheckBox:
		box := "☐ "
		if n.IsChecked {
			box = "☑ "
		}
		w.emit(out, style, box)

	default:
		w.inlineChildren(n, style, out)
	}
}

// applyTag xử lý thẻ HTML trong dòng: <br> ngắt dòng, các thẻ định dạng quen
// thuộc bật/tắt kiểu chữ, còn lại bỏ qua.
func (w *mdWriter) applyTag(n *ast.RawHTML, style Run, out *[]Run) Run {
	var raw strings.Builder
	for i := 0; i < n.Segments.Len(); i++ {
		seg := n.Segments.At(i)
		raw.Write(seg.Value(w.src))
	}
	tag := strings.ToLower(strings.Trim(raw.String(), "<>/ \t\n"))
	if sp := strings.IndexAny(tag, " \t\n"); sp >= 0 {
		tag = tag[:sp]
	}
	closing := strings.HasPrefix(strings.TrimSpace(raw.String()), "</")

	switch tag {
	case "br":
		w.emit(out, Run{}, "\n")
	case "b", "strong":
		style.Bold = !closing
	case "i", "em":
		style.Italic = !closing
	case "u", "ins":
		style.Underline = !closing
	case "s", "del", "strike":
		style.Strike = !closing
	case "sup", "sub":
		style.Script = ""
		if !closing {
			style.Script = map[string]string{"sup": "superscript", "sub": "subscript"}[tag]
		}
	}
	return style
}

// restoreMath thay ký tự đánh dấu bằng công thức, gộp các dòng của công thức
// khối thành một dòng cho dễ chuyển đổi bằng MathType.
func (w *mdWriter) restoreMath(s string) string {
	return w.restore(s, func(p mathPart) string { return strings.Join(strings.Fields(p.raw), " ") })
}

// restoreRaw trả lại nguyên văn công thức (dùng trong khối mã, đường dẫn).
func (w *mdWriter) restoreRaw(s string) string {
	return w.restore(s, func(p mathPart) string { return p.raw })
}

func (w *mdWriter) restore(s string, render func(mathPart) string) string {
	if !strings.ContainsRune(s, mathOpen) {
		return s
	}
	return mathToken.ReplaceAllStringFunc(s, func(tok string) string {
		i, err := strconv.Atoi(tok[len(string(mathOpen)) : len(tok)-len(string(mathClose))])
		if err != nil || i >= len(w.math) {
			return ""
		}
		return render(w.math[i])
	})
}

// mediaLabel đoán loại media từ đường dẫn, giống detectMedia ở frontend.
func mediaLabel(url string) string {
	u := strings.ToLower(url)
	if i := strings.IndexAny(u, "?#"); i >= 0 {
		u = u[:i]
	}
	switch {
	case strings.Contains(u, "youtube.com/") || strings.Contains(u, "youtu.be/") || strings.Contains(u, "vimeo.com/"):
		return "Video"
	case strings.Contains(u, "drive.google.com/") || strings.Contains(u, "docs.google.com/"):
		return "Tài liệu nhúng"
	}
	for _, ext := range []string{".mp4", ".webm", ".ogv", ".mov", ".m4v"} {
		if strings.HasSuffix(u, ext) {
			return "Video"
		}
	}
	for _, ext := range []string{".mp3", ".wav", ".ogg", ".oga", ".m4a", ".aac", ".flac"} {
		if strings.HasSuffix(u, ext) {
			return "Âm thanh"
		}
	}
	return "Hình ảnh"
}

func runsText(runs []Run) string {
	var b strings.Builder
	for _, r := range runs {
		b.WriteString(r.Text)
	}
	return b.String()
}

// mergeRuns gộp các run chữ liền nhau cùng định dạng cho XML gọn hơn.
func mergeRuns(runs []Run) []Run {
	out := runs[:0]
	for _, r := range runs {
		if n := len(out); n > 0 && r.Math == "" && out[n-1].Math == "" {
			last := &out[n-1]
			probe := r
			probe.Text = last.Text
			if probe == *last {
				last.Text += r.Text
				continue
			}
		}
		out = append(out, r)
	}
	return out
}

var alertMark = regexp.MustCompile(`(?i)^\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s*`)

var alertLabel = map[string]string{
	"NOTE": "Ghi chú", "TIP": "Mẹo", "IMPORTANT": "Quan trọng", "WARNING": "Lưu ý", "CAUTION": "Cảnh báo",
}

// splitAlert nhận diện đoạn mở đầu bằng [!NOTE]…, trả về nhãn tiếng Việt và
// phần chữ còn lại sau nhãn.
func splitAlert(runs []Run) (string, []Run, bool) {
	m := alertMark.FindStringSubmatch(runsText(runs))
	if m == nil {
		return "", nil, false
	}
	rest := append([]Run(nil), runs...)
	for n := len(m[0]); n > 0 && len(rest) > 0 && rest[0].Math == ""; {
		if len(rest[0].Text) <= n {
			n -= len(rest[0].Text)
			rest = rest[1:]
			continue
		}
		rest[0].Text = rest[0].Text[n:]
		n = 0
	}
	return alertLabel[strings.ToUpper(m[1])], rest, true
}
