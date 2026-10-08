package docx

// Chuyển công thức LaTeX sang phương trình Word (OMML — Office Math Markup
// Language) để giáo viên mở tệp là thấy công thức dựng sẵn, sửa tiếp được
// bằng trình soạn phương trình của Word, thay vì phải đọc mã LaTeX.
//
// Chỉ hỗ trợ tập lệnh hay gặp trong bài giảng phổ thông: phân số, căn, mũ và
// chỉ số, chữ Hy Lạp, ký hiệu quan hệ/tập hợp, dấu ngoặc co giãn, tổng/tích
// phân, hệ phương trình (cases), ma trận, dấu mũ vectơ… Gặp lệnh lạ, cả công
// thức được giữ nguyên dạng $…$ — MathType ("Toggle TeX") vẫn chuyển được.

import (
	"strconv"
	"strings"
	"unicode"
)

const nsM = "http://schemas.openxmlformats.org/officeDocument/2006/math"

// texToOMML trả về phần thân bên trong <m:oMath> của công thức, ok=false nếu
// công thức dùng cú pháp chưa hỗ trợ.
func texToOMML(tex string) (xml string, ok bool) {
	p := &texParser{toks: tokenizeTeX(tex)}
	out := p.list(func(texTok) bool { return false })
	if p.failed || p.pos < len(p.toks) || strings.TrimSpace(out) == "" {
		return "", false
	}
	return out, true
}

// ---------------------------------------------------------------------------
// Tách token
// ---------------------------------------------------------------------------

type texTokKind int

const (
	tkChar  texTokKind = iota // một ký tự thường
	tkCmd                     // \lệnh (val không kèm dấu \)
	tkOpen                    // {
	tkClose                   // }
	tkSup                     // ^
	tkSub                     // _
	tkAmp                     // & (ngăn cột trong môi trường)
	tkSpace                   // khoảng trắng — trong chế độ toán TeX bỏ qua
)

type texTok struct {
	kind texTokKind
	val  string
}

func tokenizeTeX(s string) []texTok {
	var toks []texTok
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		switch {
		case r == '\\':
			j := i + 1
			switch {
			case j >= len(rs):
				toks = append(toks, texTok{tkChar, `\`})
			case isASCIILetter(rs[j]):
				for j < len(rs) && isASCIILetter(rs[j]) {
					j++
				}
				toks = append(toks, texTok{tkCmd, string(rs[i+1 : j])})
				i = j - 1
			default:
				toks = append(toks, texTok{tkCmd, string(rs[j])})
				i = j
			}
		case r == '{':
			toks = append(toks, texTok{tkOpen, "{"})
		case r == '}':
			toks = append(toks, texTok{tkClose, "}"})
		case r == '^':
			toks = append(toks, texTok{tkSup, "^"})
		case r == '_':
			toks = append(toks, texTok{tkSub, "_"})
		case r == '&':
			toks = append(toks, texTok{tkAmp, "&"})
		case r == '%':
			for i < len(rs) && rs[i] != '\n' {
				i++
			}
		case unicode.IsSpace(r):
			if n := len(toks); n == 0 || toks[n-1].kind != tkSpace {
				toks = append(toks, texTok{tkSpace, " "})
			}
		default:
			toks = append(toks, texTok{tkChar, string(r)})
		}
	}
	return toks
}

func isASCIILetter(r rune) bool { return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' }

// ---------------------------------------------------------------------------
// Phân tích và dựng OMML
// ---------------------------------------------------------------------------

type texParser struct {
	toks   []texTok
	pos    int
	failed bool
}

// texAtom là một phần tử toán học trước khi gắn chỉ số trên/dưới.
type texAtom struct {
	xml string
	// text là chữ thường (ký tự, số, ký hiệu) chưa bọc thành run: các atom
	// chữ liền nhau được gộp chung một run, giống cách Word tự lưu phương trình.
	text string
	// nary là ký hiệu toán tử lớn (∑, ∫…): chỉ số trên/dưới thành cận.
	nary string
	// limit đánh dấu lim, max, min…: chỉ số dưới đặt ngay bên dưới tên hàm.
	limit bool
}

// render trả về XML của atom, bọc chữ thường thành run khi cần.
func (a texAtom) render() string {
	if a.xml == "" {
		return mathRun(a.text, "")
	}
	return a.xml
}

func (p *texParser) fail() { p.failed = true }

func (p *texParser) peek() (texTok, bool) {
	if p.pos >= len(p.toks) {
		return texTok{}, false
	}
	return p.toks[p.pos], true
}

func (p *texParser) next() (texTok, bool) {
	t, ok := p.peek()
	if ok {
		p.pos++
	}
	return t, ok
}

func (p *texParser) skipSpaces() {
	for p.pos < len(p.toks) && p.toks[p.pos].kind == tkSpace {
		p.pos++
	}
}

func isKind(k texTokKind) func(texTok) bool {
	return func(t texTok) bool { return t.kind == k }
}

func isCmd(names ...string) func(texTok) bool {
	return func(t texTok) bool {
		if t.kind != tkCmd {
			return false
		}
		for _, n := range names {
			if t.val == n {
				return true
			}
		}
		return false
	}
}

// list đọc một dãy phần tử tới khi gặp token dừng (không tiêu thụ token đó).
func (p *texParser) list(stop func(texTok) bool) string {
	var out, plain strings.Builder
	flush := func() {
		out.WriteString(mathRun(plain.String(), ""))
		plain.Reset()
	}
	for !p.failed {
		p.skipSpaces()
		t, ok := p.peek()
		if !ok || stop(t) {
			break
		}
		if e := p.element(stop); e.xml == "" {
			plain.WriteString(e.text)
		} else {
			flush()
			out.WriteString(e.xml)
		}
	}
	flush()
	return out.String()
}

// element đọc một phần tử kèm chỉ số trên/dưới của nó. Chữ thường không có
// chỉ số được giữ ở dạng text để list gộp run.
func (p *texParser) element(stop func(texTok) bool) texAtom {
	a := p.atom()
	sub, sup, hasSub, hasSup := p.scripts()
	switch {
	case a.nary != "":
		body := ""
		p.skipSpaces()
		if t, ok := p.peek(); ok && !stop(t) && t.kind != tkAmp && !isCmd(`\`, "end", "right")(t) {
			body = p.element(stop).render()
		}
		return texAtom{xml: naryXML(a.nary, sub, sup, hasSub, hasSup, body)}
	case a.limit && hasSub && !hasSup:
		return texAtom{xml: `<m:limLow><m:e>` + a.render() + `</m:e><m:lim>` + sub + `</m:lim></m:limLow>`}
	case !hasSub && !hasSup:
		return a
	}
	return texAtom{xml: scriptsXML(a.render(), sub, sup, hasSub, hasSup)}
}

func (p *texParser) scripts() (sub, sup string, hasSub, hasSup bool) {
	for i := 0; i < 2; i++ {
		save := p.pos
		p.skipSpaces()
		t, ok := p.peek()
		switch {
		case ok && t.kind == tkSup && !hasSup:
			p.pos++
			sup, hasSup = p.arg(), true
		case ok && t.kind == tkSub && !hasSub:
			p.pos++
			sub, hasSub = p.arg(), true
		default:
			p.pos = save
			return
		}
	}
	return
}

// arg đọc một đối số của lệnh: một nhóm {…}, một lệnh, hoặc đúng một ký tự
// (TeX hiểu \frac12 là \frac{1}{2}).
func (p *texParser) arg() string {
	p.skipSpaces()
	t, ok := p.next()
	if !ok {
		p.fail()
		return ""
	}
	switch t.kind {
	case tkOpen:
		inner := p.list(isKind(tkClose))
		p.expect(tkClose)
		return inner
	case tkChar:
		return mathRun(charText(t.val), "")
	case tkCmd:
		return p.command(t.val).render()
	}
	p.fail()
	return ""
}

func (p *texParser) expect(kind texTokKind) {
	if t, ok := p.next(); !ok || t.kind != kind {
		p.fail()
	}
}

// rawArg đọc nguyên văn một nhóm {…} (dùng cho \text, tên môi trường…).
func (p *texParser) rawArg() string {
	p.skipSpaces()
	if t, ok := p.next(); !ok || t.kind != tkOpen {
		p.fail()
		return ""
	}
	var b strings.Builder
	depth := 1
	for {
		t, ok := p.next()
		if !ok {
			p.fail()
			return b.String()
		}
		switch t.kind {
		case tkOpen:
			depth++
			continue
		case tkClose:
			if depth--; depth == 0 {
				return b.String()
			}
			continue
		case tkCmd:
			if sym, ok := texSymbols[t.val]; ok {
				b.WriteString(sym)
			} else if sp, ok := texSpaces[t.val]; ok {
				b.WriteString(sp)
			} else {
				b.WriteString(t.val)
			}
			continue
		}
		b.WriteString(t.val)
	}
}

// textArg đọc đối số dạng chữ: một nhóm {…} lấy nguyên văn, hoặc một ký tự.
func (p *texParser) textArg() string {
	p.skipSpaces()
	if t, ok := p.peek(); ok && t.kind == tkChar {
		p.pos++
		return t.val
	}
	return p.rawArg()
}

// plainArg đọc đối số chỉ gồm chữ thường (không lệnh, không chỉ số); nếu không
// phải vậy thì trả ok=false và không tiêu thụ token nào.
func (p *texParser) plainArg() (string, bool) {
	save := p.pos
	p.skipSpaces()
	t, ok := p.next()
	if ok && t.kind == tkChar {
		return t.val, true
	}
	if !ok || t.kind != tkOpen {
		p.pos = save
		return "", false
	}
	var b strings.Builder
	for {
		t, ok := p.next()
		switch {
		case ok && (t.kind == tkChar || t.kind == tkSpace):
			b.WriteString(t.val)
			continue
		case ok && t.kind == tkClose:
			return b.String(), true
		}
		p.pos = save
		return "", false
	}
}

func (p *texParser) atom() texAtom {
	t, ok := p.next()
	if !ok {
		p.fail()
		return texAtom{}
	}
	switch t.kind {
	case tkOpen:
		inner := p.list(isKind(tkClose))
		p.expect(tkClose)
		return texAtom{xml: inner}
	case tkSup, tkSub:
		// Chỉ số không có cơ sở (VD: ^{14}C): để scripts() đọc.
		p.pos--
		return texAtom{}
	case tkAmp:
		return texAtom{text: "\u2003"}
	case tkChar:
		if isDigitStr(t.val) {
			num := t.val
			for {
				n, ok := p.peek()
				if ok && n.kind == tkChar && isDigitStr(n.val) {
					num += n.val
					p.pos++
					continue
				}
				// Dấu thập phân chỉ thuộc về số khi ngay sau nó là chữ số.
				if ok && n.kind == tkChar && (n.val == "." || n.val == ",") && p.pos+1 < len(p.toks) &&
					p.toks[p.pos+1].kind == tkChar && isDigitStr(p.toks[p.pos+1].val) {
					num += n.val
					p.pos++
					continue
				}
				break
			}
			return texAtom{text: num}
		}
		return texAtom{text: charText(t.val)}
	case tkCmd:
		return p.command(t.val)
	}
	p.fail()
	return texAtom{}
}

func (p *texParser) command(name string) texAtom {
	switch name {
	case "frac", "dfrac", "tfrac", "cfrac":
		num := p.arg()
		den := p.arg()
		return texAtom{xml: `<m:f><m:num>` + num + `</m:num><m:den>` + den + `</m:den></m:f>`}

	case "binom", "dbinom", "tbinom":
		top := p.arg()
		bot := p.arg()
		return texAtom{xml: delimXML("(", ")", `<m:f><m:fPr><m:type m:val="noBar"/></m:fPr><m:num>`+top+`</m:num><m:den>`+bot+`</m:den></m:f>`)}

	case "sqrt":
		deg := ""
		p.skipSpaces()
		if t, ok := p.peek(); ok && t.kind == tkChar && t.val == "[" {
			p.pos++
			deg = p.list(func(t texTok) bool { return t.kind == tkChar && t.val == "]" })
			if t, ok := p.next(); !ok || t.val != "]" {
				p.fail()
			}
		}
		body := p.arg()
		if deg == "" {
			return texAtom{xml: `<m:rad><m:radPr><m:degHide m:val="1"/></m:radPr><m:deg/><m:e>` + body + `</m:e></m:rad>`}
		}
		return texAtom{xml: `<m:rad><m:deg>` + deg + `</m:deg><m:e>` + body + `</m:e></m:rad>`}

	case "left":
		open := p.delimiter()
		inner := p.list(isCmd("right"))
		if t, ok := p.next(); !ok || t.val != "right" {
			p.fail()
		}
		return texAtom{xml: delimXML(open, p.delimiter(), inner)}

	case "text", "textrm", "textnormal", "mbox", "mathrm", "operatorname", "textit", "mathit", "textsf", "mathsf", "texttt", "mathtt":
		return texAtom{xml: mathRun(p.textArg(), "p")}

	case "mathbf", "textbf", "boldsymbol", "bm", "pmb":
		// Chữ in đậm (VD: vectơ lực F); nhóm có công thức lồng bên trong thì dựng bình thường.
		if text, ok := p.plainArg(); ok {
			return texAtom{xml: mathRun(text, "b")}
		}
		return texAtom{xml: p.arg()}

	case "mathbb":
		return texAtom{xml: mathRun(mapLetters(p.textArg(), doubleStruck), "p")}

	case "mathcal", "mathscr":
		return texAtom{xml: mathRun(mapLetters(p.textArg(), script), "p")}

	case "vec", "overrightarrow":
		return texAtom{xml: accentXML("⃗", p.arg())}
	case "hat", "widehat":
		return texAtom{xml: accentXML("̂", p.arg())}
	case "tilde", "widetilde":
		return texAtom{xml: accentXML("̃", p.arg())}
	case "dot":
		return texAtom{xml: accentXML("̇", p.arg())}
	case "ddot":
		return texAtom{xml: accentXML("̈", p.arg())}
	case "check":
		return texAtom{xml: accentXML("̌", p.arg())}
	case "breve":
		return texAtom{xml: accentXML("̆", p.arg())}
	case "overleftarrow":
		return texAtom{xml: accentXML("⃖", p.arg())}
	case "bar", "overline":
		return texAtom{xml: `<m:bar><m:barPr><m:pos m:val="top"/></m:barPr><m:e>` + p.arg() + `</m:e></m:bar>`}
	case "underline":
		return texAtom{xml: `<m:bar><m:barPr><m:pos m:val="bot"/></m:barPr><m:e>` + p.arg() + `</m:e></m:bar>`}

	case "overset", "stackrel":
		top := p.arg()
		base := p.arg()
		return texAtom{xml: `<m:limUpp><m:e>` + base + `</m:e><m:lim>` + top + `</m:lim></m:limUpp>`}
	case "underset":
		bot := p.arg()
		base := p.arg()
		return texAtom{xml: `<m:limLow><m:e>` + base + `</m:e><m:lim>` + bot + `</m:lim></m:limLow>`}

	case "begin":
		return texAtom{xml: p.environment(p.rawArg())}

	case "not":
		p.skipSpaces()
		t, ok := p.next()
		if !ok {
			p.fail()
			return texAtom{}
		}
		if neg, ok := negated[t.val]; ok {
			return texAtom{text: neg}
		}
		p.fail()
		return texAtom{}

	case "pmod":
		return texAtom{xml: delimXML("(", ")", mathRun("mod", "p")+mathRun(" ", "")+p.arg())}
	case "bmod", "mod":
		return texAtom{xml: mathRun(" mod ", "p")}

	case "color":
		p.rawArg()
		return texAtom{}
	case "textcolor":
		p.rawArg()
		return texAtom{xml: p.arg()}

	case "displaystyle", "textstyle", "scriptstyle", "scriptscriptstyle", "limits", "nolimits",
		"big", "Big", "bigg", "Bigg", "bigl", "bigr", "Bigl", "Bigr", "biggl", "biggr", "Biggl", "Biggr",
		"bigm", "Bigm", "nonumber", "notag", "hline", "middle":
		return texAtom{}

	case `\`:
		// Xuống dòng ngoài môi trường nhiều dòng: coi như khoảng trắng.
		return texAtom{text: " "}
	}

	if sp, ok := texSpaces[name]; ok {
		return texAtom{text: sp}
	}
	if sym, ok := texSymbols[name]; ok {
		return texAtom{text: sym}
	}
	if sym, ok := texNary[name]; ok {
		return texAtom{nary: sym}
	}
	if texFunctions[name] {
		return texAtom{xml: mathRun(name, "p"), limit: texLimits[name]}
	}
	p.fail()
	return texAtom{}
}

// delimiter đọc ký hiệu ngoặc đi sau \left / \right; "." là ngoặc rỗng.
func (p *texParser) delimiter() string {
	p.skipSpaces()
	t, ok := p.next()
	if !ok {
		p.fail()
		return ""
	}
	if t.kind == tkChar {
		if t.val == "." {
			return ""
		}
		return t.val
	}
	if t.kind == tkCmd {
		if d, ok := texDelims[t.val]; ok {
			return d
		}
	}
	p.fail()
	return ""
}

// environment dựng \begin{…}…\end{…}: hệ phương trình, các dòng căn thẳng, ma trận.
func (p *texParser) environment(name string) string {
	base := strings.TrimSuffix(name, "*")
	if base == "array" || base == "alignat" {
		p.rawArg() // bỏ qua mô tả cột / số cột
	}
	rows := p.envRows(name)
	if p.failed {
		return ""
	}

	switch base {
	case "cases", "dcases":
		return delimXML("{", "", eqArrXML(rows, " "))
	case "rcases":
		return delimXML("", "}", eqArrXML(rows, " "))
	case "aligned", "align", "alignat", "gathered", "gather", "split", "eqnarray", "multline":
		return eqArrXML(rows, "")
	case "matrix", "smallmatrix", "array":
		return matrixXML(rows)
	case "pmatrix":
		return delimXML("(", ")", matrixXML(rows))
	case "bmatrix":
		return delimXML("[", "]", matrixXML(rows))
	case "Bmatrix":
		return delimXML("{", "}", matrixXML(rows))
	case "vmatrix":
		return delimXML("|", "|", matrixXML(rows))
	case "Vmatrix":
		return delimXML("‖", "‖", matrixXML(rows))
	}
	p.fail()
	return ""
}

// envRows đọc các dòng (ngăn bởi \\) và cột (ngăn bởi &) tới \end{name}.
func (p *texParser) envRows(name string) [][]string {
	var rows [][]string
	var row []string
	stop := func(t texTok) bool { return t.kind == tkAmp || isCmd(`\`, "end")(t) }
	for !p.failed {
		row = append(row, p.list(stop))
		t, ok := p.next()
		if !ok {
			p.fail()
			break
		}
		if t.kind == tkAmp {
			continue
		}
		if t.val == `\` {
			rows = append(rows, row)
			row = nil
			p.skipSpaces()
			// Bỏ qua khoảng cách dòng tuỳ chọn kiểu \\[2pt].
			if n, ok := p.peek(); ok && n.kind == tkChar && n.val == "[" {
				for t, ok := p.next(); ok && t.val != "]"; t, ok = p.next() {
				}
			}
			continue
		}
		// \end{…}
		if p.rawArg() != name {
			p.fail()
		}
		if !(len(row) == 1 && strings.TrimSpace(row[0]) == "") || len(rows) == 0 {
			rows = append(rows, row)
		}
		break
	}
	return rows
}

// ---------------------------------------------------------------------------
// Mảnh XML
// ---------------------------------------------------------------------------

func mathRun(text, sty string) string {
	if text == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString("<m:r>")
	if sty != "" {
		b.WriteString(`<m:rPr><m:sty m:val="` + sty + `"/></m:rPr>`)
	}
	b.WriteString(`<w:rPr><w:rFonts w:ascii="Cambria Math" w:hAnsi="Cambria Math"/></w:rPr>`)
	b.WriteString(`<m:t xml:space="preserve">` + escapeXML(text) + `</m:t></m:r>`)
	return b.String()
}

func scriptsXML(base, sub, sup string, hasSub, hasSup bool) string {
	switch {
	case hasSub && hasSup:
		return `<m:sSubSup><m:e>` + base + `</m:e><m:sub>` + sub + `</m:sub><m:sup>` + sup + `</m:sup></m:sSubSup>`
	case hasSup:
		return `<m:sSup><m:e>` + base + `</m:e><m:sup>` + sup + `</m:sup></m:sSup>`
	case hasSub:
		return `<m:sSub><m:e>` + base + `</m:e><m:sub>` + sub + `</m:sub></m:sSub>`
	}
	return base
}

func naryXML(chr, sub, sup string, hasSub, hasSup bool, body string) string {
	loc := "undOvr"
	if strings.ContainsAny(chr, "∫∬∭∮") {
		loc = "subSup"
	}
	var b strings.Builder
	b.WriteString(`<m:nary><m:naryPr><m:chr m:val="` + chr + `"/><m:limLoc m:val="` + loc + `"/>`)
	if !hasSub {
		b.WriteString(`<m:subHide m:val="1"/>`)
	}
	if !hasSup {
		b.WriteString(`<m:supHide m:val="1"/>`)
	}
	b.WriteString(`</m:naryPr><m:sub>` + sub + `</m:sub><m:sup>` + sup + `</m:sup><m:e>` + body + `</m:e></m:nary>`)
	return b.String()
}

func delimXML(open, close, inner string) string {
	return `<m:d><m:dPr><m:begChr m:val="` + escapeXML(open) + `"/><m:endChr m:val="` + escapeXML(close) +
		`"/></m:dPr><m:e>` + inner + `</m:e></m:d>`
}

func accentXML(chr, body string) string {
	return `<m:acc><m:accPr><m:chr m:val="` + chr + `"/></m:accPr><m:e>` + body + `</m:e></m:acc>`
}

// eqArrXML dựng các dòng xếp chồng; sep chèn vào chỗ dấu & (VD: giữa biểu
// thức và điều kiện trong hệ cases).
func eqArrXML(rows [][]string, sep string) string {
	var b strings.Builder
	b.WriteString("<m:eqArr>")
	for _, row := range rows {
		b.WriteString("<m:e>" + strings.Join(row, mathRun(sep, "")) + "</m:e>")
	}
	b.WriteString("</m:eqArr>")
	return b.String()
}

func matrixXML(rows [][]string) string {
	cols := 1
	for _, row := range rows {
		cols = max(cols, len(row))
	}
	var b strings.Builder
	b.WriteString(`<m:m><m:mPr><m:mcs><m:mc><m:mcPr><m:count m:val="`)
	b.WriteString(strconv.Itoa(cols))
	b.WriteString(`"/><m:mcJc m:val="center"/></m:mcPr></m:mc></m:mcs></m:mPr>`)
	for _, row := range rows {
		b.WriteString("<m:mr>")
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			b.WriteString("<m:e>" + cell + "</m:e>")
		}
		b.WriteString("</m:mr>")
	}
	b.WriteString("</m:m>")
	return b.String()
}

func isDigitStr(s string) bool { return len(s) == 1 && s[0] >= '0' && s[0] <= '9' }

// charText đổi vài ký tự ASCII sang dạng ký hiệu toán học tương ứng.
func charText(s string) string {
	switch s {
	case "-":
		return "−"
	case "*":
		return "∗"
	case "'":
		return "′"
	case "~":
		return " "
	}
	return s
}

func mapLetters(s string, table func(rune) rune) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(table(r))
	}
	return b.String()
}

// doubleStruck: R → ℝ, N → ℕ… (khối Mathematical Double-Struck, trừ các chữ đã có sẵn ở BMP).
func doubleStruck(r rune) rune {
	if bmp, ok := map[rune]rune{'C': 'ℂ', 'H': 'ℍ', 'N': 'ℕ', 'P': 'ℙ', 'Q': 'ℚ', 'R': 'ℝ', 'Z': 'ℤ'}[r]; ok {
		return bmp
	}
	switch {
	case r >= 'A' && r <= 'Z':
		return 0x1D538 + (r - 'A')
	case r >= 'a' && r <= 'z':
		return 0x1D552 + (r - 'a')
	case r >= '0' && r <= '9':
		return 0x1D7D8 + (r - '0')
	}
	return r
}

// script: E → ℰ, F → ℱ… (khối Mathematical Script, trừ các chữ đã có sẵn ở BMP).
func script(r rune) rune {
	if bmp, ok := map[rune]rune{
		'B': 'ℬ', 'E': 'ℰ', 'F': 'ℱ', 'H': 'ℋ', 'I': 'ℐ', 'L': 'ℒ', 'M': 'ℳ', 'R': 'ℛ',
		'e': 'ℯ', 'g': 'ℊ', 'o': 'ℴ',
	}[r]; ok {
		return bmp
	}
	switch {
	case r >= 'A' && r <= 'Z':
		return 0x1D49C + (r - 'A')
	case r >= 'a' && r <= 'z':
		return 0x1D4B6 + (r - 'a')
	}
	return r
}

// ---------------------------------------------------------------------------
// Bảng ký hiệu
// ---------------------------------------------------------------------------

var texSpaces = map[string]string{
	",": " ", ":": " ", ">": " ", ";": " ", "!": "", " ": " ",
	"quad": " ", "qquad": "  ", "enspace": " ", "thinspace": " ",
}

var texDelims = map[string]string{
	"{": "{", "}": "}", "lbrace": "{", "rbrace": "}", "langle": "⟨", "rangle": "⟩",
	"|": "‖", "vert": "|", "lvert": "|", "rvert": "|", "Vert": "‖", "lVert": "‖", "rVert": "‖",
	"lfloor": "⌊", "rfloor": "⌋", "lceil": "⌈", "rceil": "⌉", "lbrack": "[", "rbrack": "]",
}

var negated = map[string]string{
	"=": "≠", "in": "∉", "subset": "⊄", "subseteq": "⊈", "supset": "⊅", "equiv": "≢",
	"parallel": "∦", "exists": "∄", "<": "≮", ">": "≯", "le": "≰", "leq": "≰", "ge": "≱", "geq": "≱",
}

var texNary = map[string]string{
	"sum": "∑", "prod": "∏", "coprod": "∐", "int": "∫", "iint": "∬", "iiint": "∭", "oint": "∮",
	"bigcup": "⋃", "bigcap": "⋂", "bigoplus": "⨁", "bigotimes": "⨂", "bigvee": "⋁", "bigwedge": "⋀",
}

var texFunctions = map[string]bool{
	"sin": true, "cos": true, "tan": true, "cot": true, "sec": true, "csc": true,
	"arcsin": true, "arccos": true, "arctan": true, "sinh": true, "cosh": true, "tanh": true, "coth": true,
	"log": true, "ln": true, "lg": true, "exp": true, "lim": true, "max": true, "min": true, "sup": true, "inf": true,
	"det": true, "gcd": true, "deg": true, "dim": true, "ker": true, "arg": true, "sgn": true, "tg": true, "cotg": true,
}

var texLimits = map[string]bool{"lim": true, "max": true, "min": true, "sup": true, "inf": true}

var texSymbols = map[string]string{
	// Chữ Hy Lạp
	"alpha": "α", "beta": "β", "gamma": "γ", "delta": "δ", "epsilon": "ϵ", "varepsilon": "ε", "zeta": "ζ",
	"eta": "η", "theta": "θ", "vartheta": "ϑ", "iota": "ι", "kappa": "κ", "lambda": "λ", "mu": "μ", "nu": "ν",
	"xi": "ξ", "omicron": "ο", "pi": "π", "varpi": "ϖ", "rho": "ρ", "varrho": "ϱ", "sigma": "σ", "varsigma": "ς",
	"tau": "τ", "upsilon": "υ", "phi": "ϕ", "varphi": "φ", "chi": "χ", "psi": "ψ", "omega": "ω",
	"Gamma": "Γ", "Delta": "Δ", "Theta": "Θ", "Lambda": "Λ", "Xi": "Ξ", "Pi": "Π", "Sigma": "Σ",
	"Upsilon": "Υ", "Phi": "Φ", "Psi": "Ψ", "Omega": "Ω",

	// Phép toán và quan hệ
	"pm": "±", "mp": "∓", "times": "×", "cdot": "⋅", "div": "÷", "ast": "∗", "star": "⋆", "circ": "∘",
	"bullet": "∙", "oplus": "⊕", "otimes": "⊗", "neq": "≠", "ne": "≠", "le": "≤", "leq": "≤", "ge": "≥",
	"geq": "≥", "leqslant": "⩽", "geqslant": "⩾", "ll": "≪", "gg": "≫", "approx": "≈", "equiv": "≡",
	"sim": "∼", "simeq": "≃", "cong": "≅", "propto": "∝", "perp": "⊥", "parallel": "∥", "mid": "∣",
	"nmid": "∤", "models": "⊨", "vdash": "⊢",

	// Tập hợp và logic
	"in": "∈", "notin": "∉", "ni": "∋", "subset": "⊂", "subseteq": "⊆", "supset": "⊃", "supseteq": "⊇",
	"subsetneq": "⊊", "cup": "∪", "cap": "∩", "setminus": "∖", "emptyset": "∅", "varnothing": "∅",
	"forall": "∀", "exists": "∃", "nexists": "∄", "neg": "¬", "lnot": "¬", "land": "∧", "lor": "∨",
	"wedge": "∧", "vee": "∨",

	// Mũi tên
	"to": "→", "rightarrow": "→", "leftarrow": "←", "gets": "←", "leftrightarrow": "↔",
	"Rightarrow": "⇒", "Leftarrow": "⇐", "Leftrightarrow": "⇔", "implies": "⟹", "impliedby": "⟸",
	"iff": "⟺", "longrightarrow": "⟶", "longleftarrow": "⟵", "Longrightarrow": "⟹", "Longleftrightarrow": "⟺",
	"mapsto": "↦", "uparrow": "↑", "downarrow": "↓", "rightleftharpoons": "⇌", "nearrow": "↗", "searrow": "↘",

	// Hình học và khác
	"angle": "∠", "measuredangle": "∡", "triangle": "△", "square": "□", "degree": "°", "infty": "∞",
	"partial": "∂", "nabla": "∇", "hbar": "ℏ", "ell": "ℓ", "Re": "ℜ", "Im": "ℑ", "aleph": "ℵ",
	"ldots": "…", "dots": "…", "cdots": "⋯", "vdots": "⋮", "ddots": "⋱", "prime": "′",
	"therefore": "∴", "because": "∵", "backslash": `\`, "colon": ":",
	"langle": "⟨", "rangle": "⟩", "lfloor": "⌊", "rfloor": "⌋", "lceil": "⌈", "rceil": "⌉",
	"lbrace": "{", "rbrace": "}", "vert": "|", "Vert": "‖",
	"{": "{", "}": "}", "|": "‖", "%": "%", "$": "$", "&": "&", "#": "#", "_": "_",
}
