package docx

// Phần ghi: dựng tệp Word (.docx) từ đầu để giáo viên tải giáo án, bộ câu hỏi
// và bài giảng về máy, chỉnh sửa tiếp hoặc in ra.
//
// Giống phần đọc, ở đây tự ghép các tệp XML của chuẩn OOXML rồi nén zip thay
// vì kéo thêm thư viện ngoài: chỉ cần một tập nhỏ tính năng (đoạn văn, tiêu đề,
// danh sách, bảng, liên kết, số trang) nên tự viết vẫn gọn và kiểm soát được
// định dạng theo thói quen văn bản của giáo viên Việt Nam (khổ A4, Times New
// Roman, lề trái 3 cm).

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Khổ giấy A4 và lề trang theo đơn vị twip (1/20 point, 1 cm ≈ 567 twip).
const (
	pageWidth   = 11906
	pageHeight  = 16838
	marginLeft  = 1701 // 3 cm
	marginRight = 1134 // 2 cm
	marginTop   = 1134
	marginBot   = 1134

	// ContentWidth là bề ngang vùng soạn thảo, dùng để chia độ rộng cột bảng.
	ContentWidth = pageWidth - marginLeft - marginRight

	// cellPadding là lề trái + phải bên trong mỗi ô bảng.
	cellPadding = 2 * 100
)

// Align là căn lề của đoạn văn.
type Align string

const (
	AlignLeft    Align = ""
	AlignCenter  Align = "center"
	AlignRight   Align = "right"
	AlignJustify Align = "both"
)

// Các kiểu đoạn văn khai báo trong styles.xml.
const (
	StyleTitle    = "Title"
	StyleSubtitle = "Subtitle"
	StyleHeading1 = "Heading1"
	StyleHeading2 = "Heading2"
	StyleHeading3 = "Heading3"
	StyleQuote    = "Quote"
	StyleCode     = "Code"
	StyleRule     = "Rule"
	styleTable    = "TableText"
)

// Run là một đoạn chữ cùng định dạng. Ký tự "\n" trong Text thành ngắt dòng,
// "\t" thành dấu tab.
type Run struct {
	Text      string
	Bold      bool
	Italic    bool
	Underline bool
	Strike    bool
	Code      bool
	// "superscript" hoặc "subscript"; rỗng = bình thường.
	Script string
	// Link là địa chỉ http(s)/mailto; địa chỉ khác bị bỏ qua, chữ vẫn giữ.
	Link string
	// Math là công thức LaTeX (không kèm dấu $), ghi thành phương trình Word;
	// Text bị bỏ qua. Display đánh dấu công thức khối ($$…$$).
	Math    string
	Display bool
}

// Para mô tả một đoạn văn.
type Para struct {
	// Style là một trong các hằng Style*; rỗng = đoạn văn thường.
	Style string
	Align Align
	// List là mã danh sách lấy từ Document.NewList; 0 = không thuộc danh sách.
	List  int
	Level int
	// Indent là lùi đầu dòng (twip) cho đoạn không đánh số.
	Indent int
	// SpaceBefore là khoảng cách phía trên đoạn (twip); 0 = theo kiểu đoạn.
	SpaceBefore int
	// Tabs là các điểm dừng tab, dùng cho dòng kẻ chấm kiểu "Họ và tên: ……".
	Tabs     []TabStop
	KeepNext bool
	Runs     []Run
}

// TabStop là một điểm dừng tab, tính từ lề trái (twip). DotLeader điền dấu
// chấm vào khoảng trống trước điểm dừng.
type TabStop struct {
	Pos       int
	Right     bool
	DotLeader bool
}

// Document là một tài liệu Word đang dựng. Nội dung được thêm vào Body();
// gọi Bytes() để lấy tệp .docx hoàn chỉnh.
type Document struct {
	// Title và Author ghi vào thuộc tính tệp (File → Info trong Word).
	Title  string
	Author string

	body  *Container
	links []string
	lists []listDef
}

type listDef struct {
	ordered bool
	start   int
}

// New tạo một tài liệu trống.
func New(title string) *Document {
	d := &Document{Title: title}
	d.body = &Container{doc: d, width: ContentWidth}
	return d
}

// Body là vùng nội dung chính của tài liệu.
func (d *Document) Body() *Container { return d.body }

// NewList cấp mã cho một danh sách mới. Mỗi danh sách có mã riêng để số thứ
// tự bắt đầu lại từ start thay vì nối tiếp danh sách trước đó.
func (d *Document) NewList(ordered bool, start int) int {
	if start < 1 {
		start = 1
	}
	d.lists = append(d.lists, listDef{ordered: ordered, start: start})
	return len(d.lists)
}

func (d *Document) linkID(url string) string {
	d.links = append(d.links, url)
	return fmt.Sprintf("rIdL%d", len(d.links))
}

// Container gom các khối nội dung (đoạn văn, bảng): thân tài liệu hoặc một ô bảng.
type Container struct {
	doc   *Document
	buf   strings.Builder
	width int
	cell  bool
	// Ô bảng bắt buộc kết thúc bằng một đoạn văn, nên cần biết khối cuối có phải
	// bảng; đoạn ngay sau bảng cũng được giãn cách để không dính vào bảng.
	endsWithTable bool
}

// Width là bề ngang khả dụng (twip) của vùng chứa.
func (c *Container) Width() int { return c.width }

// Empty cho biết vùng chứa chưa có nội dung nào.
func (c *Container) Empty() bool { return c.buf.Len() == 0 }

// Add thêm một đoạn văn.
func (c *Container) Add(p Para) {
	if p.Style == "" && c.cell {
		p.Style = styleTable
	}
	if c.endsWithTable && p.SpaceBefore == 0 {
		p.SpaceBefore = 160
	}
	b := &c.buf
	b.WriteString("<w:p>")
	writeParaProps(b, p)
	if m, ok := soleDisplayMath(p.Runs); ok {
		if xml, ok := texToOMML(m.Math); ok {
			// Công thức khối đứng riêng một đoạn: dùng oMathPara để Word căn giữa như phương trình.
			b.WriteString("<m:oMathPara><m:oMath>" + xml + "</m:oMath></m:oMathPara></w:p>")
			c.endsWithTable = false
			return
		}
	}
	for _, r := range p.Runs {
		c.writeRun(r)
	}
	b.WriteString("</w:p>")
	c.endsWithTable = false
}

// soleDisplayMath trả về công thức khối nếu đoạn văn chỉ gồm đúng công thức đó.
func soleDisplayMath(runs []Run) (Run, bool) {
	var found *Run
	for i := range runs {
		switch {
		case runs[i].Math != "":
			if found != nil || !runs[i].Display {
				return Run{}, false
			}
			found = &runs[i]
		case strings.TrimSpace(runs[i].Text) != "":
			return Run{}, false
		}
	}
	if found == nil {
		return Run{}, false
	}
	return *found, true
}

// Text thêm một đoạn văn chỉ gồm một đoạn chữ thường.
func (c *Container) Text(style, text string) {
	c.Add(Para{Style: style, Runs: []Run{{Text: text}}})
}

// PageBreak ngắt sang trang mới.
func (c *Container) PageBreak() {
	c.buf.WriteString(`<w:p><w:r><w:br w:type="page"/></w:r></w:p>`)
	c.endsWithTable = false
}

func writeParaProps(b *strings.Builder, p Para) {
	var props strings.Builder
	if p.Style != "" {
		fmt.Fprintf(&props, `<w:pStyle w:val="%s"/>`, p.Style)
	}
	if p.KeepNext {
		props.WriteString("<w:keepNext/>")
	}
	if p.List > 0 {
		fmt.Fprintf(&props, `<w:numPr><w:ilvl w:val="%d"/><w:numId w:val="%d"/></w:numPr>`, clampLevel(p.Level), p.List)
	}
	if len(p.Tabs) > 0 {
		props.WriteString("<w:tabs>")
		for _, t := range p.Tabs {
			align, leader := "left", "none"
			if t.Right {
				align = "right"
			}
			if t.DotLeader {
				leader = "dot"
			}
			fmt.Fprintf(&props, `<w:tab w:val="%s" w:leader="%s" w:pos="%d"/>`, align, leader, t.Pos)
		}
		props.WriteString("</w:tabs>")
	}
	if p.SpaceBefore > 0 {
		fmt.Fprintf(&props, `<w:spacing w:before="%d"/>`, p.SpaceBefore)
	}
	if p.List == 0 && p.Indent > 0 {
		fmt.Fprintf(&props, `<w:ind w:left="%d"/>`, p.Indent)
	}
	if p.Align != AlignLeft {
		fmt.Fprintf(&props, `<w:jc w:val="%s"/>`, p.Align)
	}
	if props.Len() > 0 {
		b.WriteString("<w:pPr>")
		b.WriteString(props.String())
		b.WriteString("</w:pPr>")
	}
}

func clampLevel(level int) int {
	return max(0, min(level, 8))
}

func (c *Container) writeRun(r Run) {
	if r.Math != "" {
		if xml, ok := texToOMML(r.Math); ok {
			c.buf.WriteString("<m:oMath>" + xml + "</m:oMath>")
			return
		}
		// Cú pháp chưa hỗ trợ: giữ nguyên LaTeX để MathType / Word chuyển tay.
		tex := strings.Join(strings.Fields(r.Math), " ")
		if r.Display {
			r.Text = "$$ " + tex + " $$"
		} else {
			r.Text = "$" + tex + "$"
		}
		r.Math = ""
	}
	if r.Text == "" {
		return
	}
	b := &c.buf
	link := safeLink(r.Link)
	if link != "" {
		fmt.Fprintf(b, `<w:hyperlink r:id="%s">`, c.doc.linkID(link))
	}

	b.WriteString("<w:r>")
	var props strings.Builder
	switch {
	case link != "":
		props.WriteString(`<w:rStyle w:val="Hyperlink"/>`)
	case r.Code:
		props.WriteString(`<w:rStyle w:val="CodeChar"/>`)
	}
	if r.Bold {
		props.WriteString("<w:b/><w:bCs/>")
	}
	if r.Italic {
		props.WriteString("<w:i/><w:iCs/>")
	}
	if r.Strike {
		props.WriteString("<w:strike/>")
	}
	if r.Underline && link == "" {
		props.WriteString(`<w:u w:val="single"/>`)
	}
	if r.Script == "superscript" || r.Script == "subscript" {
		fmt.Fprintf(&props, `<w:vertAlign w:val="%s"/>`, r.Script)
	}
	if props.Len() > 0 {
		b.WriteString("<w:rPr>")
		b.WriteString(props.String())
		b.WriteString("</w:rPr>")
	}
	writeRunText(b, r.Text)
	b.WriteString("</w:r>")

	if link != "" {
		b.WriteString("</w:hyperlink>")
	}
}

// writeRunText ghi chữ của một run, đổi "\n" thành ngắt dòng và "\t" thành tab.
func writeRunText(b *strings.Builder, text string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	start := 0
	flush := func(end int) {
		if end > start {
			b.WriteString(`<w:t xml:space="preserve">`)
			b.WriteString(escapeXML(text[start:end]))
			b.WriteString("</w:t>")
		}
	}
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\n':
			flush(i)
			b.WriteString("<w:br/>")
			start = i + 1
		case '\t':
			flush(i)
			b.WriteString("<w:tab/>")
			start = i + 1
		}
	}
	flush(len(text))
}

// safeLink chỉ giữ lại liên kết http(s) và mailto.
func safeLink(url string) string {
	url = strings.TrimSpace(url)
	lower := strings.ToLower(url)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "mailto:") {
		return url
	}
	return ""
}

// escapeXML thoát ký tự đặc biệt và bỏ các ký tự điều khiển mà XML 1.0 không
// cho phép — chỉ một ký tự lạ trong nội dung AI sinh ra cũng đủ làm Word báo
// tệp hỏng.
func escapeXML(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == '&':
			b.WriteString("&amp;")
		case r == '<':
			b.WriteString("&lt;")
		case r == '>':
			b.WriteString("&gt;")
		case r == '"':
			b.WriteString("&quot;")
		case r == '\t' || r == '\n' || r == '\r',
			r >= 0x20 && r <= 0xD7FF,
			r >= 0xE000 && r <= 0xFFFD,
			r >= 0x10000 && r <= 0x10FFFF:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ---------------------------------------------------------------------------
// Bảng
// ---------------------------------------------------------------------------

// Table là một bảng đang dựng. Tạo bằng Container.NewTable, thêm dòng bằng
// Row rồi đưa vào tài liệu bằng Container.AddTable.
type Table struct {
	doc    *Document
	widths []int
	// HeaderRows là số dòng đầu được tô nền và lặp lại ở đầu mỗi trang.
	HeaderRows int
	rows       [][]*Container
}

// NewTable tạo bảng với độ rộng tương đối của từng cột; tổng độ rộng được co
// giãn cho vừa bề ngang của vùng chứa.
func (c *Container) NewTable(weights ...int) *Table {
	total := 0
	for _, w := range weights {
		total += max(w, 1)
	}
	widths := make([]int, len(weights))
	for i, w := range weights {
		widths[i] = c.width * max(w, 1) / total
	}
	return &Table{doc: c.doc, widths: widths}
}

// Row thêm một dòng mới và trả về các ô của dòng đó để điền nội dung.
func (t *Table) Row() []*Container {
	cells := make([]*Container, len(t.widths))
	for i, w := range t.widths {
		cells[i] = &Container{doc: t.doc, width: max(w-cellPadding, 567), cell: true}
	}
	t.rows = append(t.rows, cells)
	return cells
}

// AddTable ghi bảng vào vùng chứa.
func (c *Container) AddTable(t *Table) {
	if len(t.rows) == 0 {
		return
	}
	b := &c.buf
	total := 0
	for _, w := range t.widths {
		total += w
	}
	fmt.Fprintf(b, `<w:tbl><w:tblPr><w:tblW w:w="%d" w:type="dxa"/>`, total)
	b.WriteString(`<w:tblBorders>`)
	for _, side := range []string{"top", "left", "bottom", "right", "insideH", "insideV"} {
		fmt.Fprintf(b, `<w:%s w:val="single" w:sz="4" w:space="0" w:color="808080"/>`, side)
	}
	b.WriteString(`</w:tblBorders><w:tblLayout w:type="fixed"/>`)
	b.WriteString(`<w:tblCellMar><w:top w:w="40" w:type="dxa"/><w:left w:w="100" w:type="dxa"/>` +
		`<w:bottom w:w="40" w:type="dxa"/><w:right w:w="100" w:type="dxa"/></w:tblCellMar>`)
	b.WriteString(`<w:tblLook w:val="04A0" w:firstRow="1" w:lastRow="0" w:firstColumn="0" w:lastColumn="0" w:noHBand="1" w:noVBand="1"/></w:tblPr>`)

	b.WriteString("<w:tblGrid>")
	for _, w := range t.widths {
		fmt.Fprintf(b, `<w:gridCol w:w="%d"/>`, w)
	}
	b.WriteString("</w:tblGrid>")

	for i, row := range t.rows {
		header := i < t.HeaderRows
		b.WriteString("<w:tr>")
		if header {
			b.WriteString("<w:trPr><w:cantSplit/><w:tblHeader/></w:trPr>")
		}
		for j, cell := range row {
			fmt.Fprintf(b, `<w:tc><w:tcPr><w:tcW w:w="%d" w:type="dxa"/>`, t.widths[j])
			if header {
				b.WriteString(`<w:shd w:val="clear" w:color="auto" w:fill="E8EEF6"/>`)
			}
			b.WriteString("</w:tcPr>")
			b.WriteString(cell.buf.String())
			if cell.Empty() || cell.endsWithTable {
				b.WriteString(`<w:p><w:pPr><w:pStyle w:val="TableText"/></w:pPr></w:p>`)
			}
			b.WriteString("</w:tc>")
		}
		b.WriteString("</w:tr>")
	}
	b.WriteString("</w:tbl>")
	c.endsWithTable = true
}

// ---------------------------------------------------------------------------
// Đóng gói
// ---------------------------------------------------------------------------

// Bytes đóng gói tài liệu thành tệp .docx.
func (d *Document) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	parts := []struct{ name, content string }{
		{"[Content_Types].xml", contentTypesXML},
		{"_rels/.rels", rootRelsXML},
		{"docProps/core.xml", d.coreXML()},
		{"word/_rels/document.xml.rels", d.documentRelsXML()},
		{"word/document.xml", d.documentXML()},
		{"word/styles.xml", stylesXML},
		{"word/numbering.xml", d.numberingXML()},
		{"word/settings.xml", settingsXML},
		{"word/footer1.xml", footerXML},
	}
	for _, p := range parts {
		w, err := zw.Create(p.name)
		if err != nil {
			return nil, fmt.Errorf("tạo %s: %w", p.name, err)
		}
		if _, err := w.Write([]byte(p.content)); err != nil {
			return nil, fmt.Errorf("ghi %s: %w", p.name, err)
		}
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("đóng gói tệp Word: %w", err)
	}
	return buf.Bytes(), nil
}

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

const (
	nsW = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
	nsR = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
)

func (d *Document) documentXML() string {
	body := d.body
	if body.endsWithTable {
		// Word cần một đoạn văn giữa bảng cuối cùng và phần thiết lập trang.
		body.buf.WriteString("<w:p/>")
		body.endsWithTable = false
	}
	return xmlHeader +
		`<w:document xmlns:w="` + nsW + `" xmlns:r="` + nsR + `" xmlns:m="` + nsM + `"><w:body>` +
		body.buf.String() +
		`<w:sectPr><w:footerReference w:type="default" r:id="rIdFooter"/>` +
		fmt.Sprintf(`<w:pgSz w:w="%d" w:h="%d"/>`, pageWidth, pageHeight) +
		fmt.Sprintf(`<w:pgMar w:top="%d" w:right="%d" w:bottom="%d" w:left="%d" w:header="567" w:footer="567" w:gutter="0"/>`,
			marginTop, marginRight, marginBot, marginLeft) +
		`</w:sectPr></w:body></w:document>`
}

func (d *Document) documentRelsXML() string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	b.WriteString(`<Relationship Id="rIdStyles" Type="` + nsR + `/styles" Target="styles.xml"/>`)
	b.WriteString(`<Relationship Id="rIdNumbering" Type="` + nsR + `/numbering" Target="numbering.xml"/>`)
	b.WriteString(`<Relationship Id="rIdSettings" Type="` + nsR + `/settings" Target="settings.xml"/>`)
	b.WriteString(`<Relationship Id="rIdFooter" Type="` + nsR + `/footer" Target="footer1.xml"/>`)
	for i, url := range d.links {
		fmt.Fprintf(&b, `<Relationship Id="rIdL%d" Type="%s/hyperlink" Target="%s" TargetMode="External"/>`,
			i+1, nsR, escapeXML(url))
	}
	b.WriteString(`</Relationships>`)
	return b.String()
}

func (d *Document) coreXML() string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"` +
		` xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/"` +
		` xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`)
	if d.Title != "" {
		b.WriteString("<dc:title>" + escapeXML(d.Title) + "</dc:title>")
	}
	if d.Author != "" {
		b.WriteString("<dc:creator>" + escapeXML(d.Author) + "</dc:creator>")
	}
	now := time.Now().UTC().Format(time.RFC3339)
	b.WriteString(`<dcterms:created xsi:type="dcterms:W3CDTF">` + now + `</dcterms:created>`)
	b.WriteString(`<dcterms:modified xsi:type="dcterms:W3CDTF">` + now + `</dcterms:modified>`)
	b.WriteString(`</cp:coreProperties>`)
	return b.String()
}

// Ký hiệu đầu dòng theo từng cấp: "-" rồi "+" như văn bản hành chính quen thuộc.
var bulletSymbols = []string{"-", "+", "•"}

// Kiểu đánh số theo từng cấp của danh sách có thứ tự: 1. → a) → i.
var orderedFormats = []struct{ format, text string }{
	{"decimal", "%%%d."},
	{"lowerLetter", "%%%d)"},
	{"lowerRoman", "%%%d."},
}

func (d *Document) numberingXML() string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`<w:numbering xmlns:w="` + nsW + `">`)
	for id, ordered := range []bool{false, true} {
		fmt.Fprintf(&b, `<w:abstractNum w:abstractNumId="%d"><w:multiLevelType w:val="hybridMultilevel"/>`, id)
		for lvl := 0; lvl <= 8; lvl++ {
			fmt.Fprintf(&b, `<w:lvl w:ilvl="%d"><w:start w:val="1"/>`, lvl)
			if ordered {
				f := orderedFormats[lvl%len(orderedFormats)]
				fmt.Fprintf(&b, `<w:numFmt w:val="%s"/><w:lvlText w:val="%s"/>`, f.format, fmt.Sprintf(f.text, lvl+1))
			} else {
				fmt.Fprintf(&b, `<w:numFmt w:val="bullet"/><w:lvlText w:val="%s"/>`, bulletSymbols[lvl%len(bulletSymbols)])
			}
			fmt.Fprintf(&b, `<w:lvlJc w:val="left"/><w:pPr><w:ind w:left="%d" w:hanging="283"/></w:pPr></w:lvl>`, 567*(lvl+1))
		}
		b.WriteString(`</w:abstractNum>`)
	}
	for i, l := range d.lists {
		abstract := 0
		if l.ordered {
			abstract = 1
		}
		fmt.Fprintf(&b, `<w:num w:numId="%d"><w:abstractNumId w:val="%d"/>`, i+1, abstract)
		fmt.Fprintf(&b, `<w:lvlOverride w:ilvl="0"><w:startOverride w:val="%d"/></w:lvlOverride></w:num>`, l.start)
	}
	b.WriteString(`</w:numbering>`)
	return b.String()
}

const contentTypesXML = xmlHeader +
	`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
	`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
	`<Default Extension="xml" ContentType="application/xml"/>` +
	`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
	`<Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>` +
	`<Override PartName="/word/numbering.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.numbering+xml"/>` +
	`<Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/>` +
	`<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>` +
	`<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>` +
	`</Types>`

const rootRelsXML = xmlHeader +
	`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
	`<Relationship Id="rId1" Type="` + nsR + `/officeDocument" Target="word/document.xml"/>` +
	`<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>` +
	`</Relationships>`

// settingsXML đặt chế độ tương thích Word 2013+ để tệp không mở ở "Compatibility Mode".
const settingsXML = xmlHeader +
	`<w:settings xmlns:w="` + nsW + `">` +
	`<w:defaultTabStop w:val="567"/>` +
	`<w:characterSpacingControl w:val="doNotCompress"/>` +
	`<w:compat><w:compatSetting w:name="compatibilityMode" w:uri="http://schemas.microsoft.com/office/word" w:val="15"/></w:compat>` +
	`</w:settings>`

// footerXML đánh số trang ở giữa chân trang.
const footerXML = xmlHeader +
	`<w:ftr xmlns:w="` + nsW + `" xmlns:r="` + nsR + `">` +
	`<w:p><w:pPr><w:jc w:val="center"/></w:pPr>` +
	`<w:fldSimple w:instr=" PAGE "><w:r><w:t>1</w:t></w:r></w:fldSimple></w:p>` +
	`</w:ftr>`

// stylesXML: Times New Roman 13 pt, giãn dòng 1,15 cho toàn văn bản; tiêu đề
// các cấp có outlineLvl để hiện trong ngăn điều hướng của Word.
const stylesXML = xmlHeader +
	`<w:styles xmlns:w="` + nsW + `">` +
	`<w:docDefaults>` +
	`<w:rPrDefault><w:rPr><w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman" w:cs="Times New Roman"/>` +
	`<w:sz w:val="26"/><w:szCs w:val="26"/><w:lang w:val="vi-VN" w:eastAsia="en-US" w:bidi="ar-SA"/></w:rPr></w:rPrDefault>` +
	`<w:pPrDefault><w:pPr><w:spacing w:after="80" w:line="276" w:lineRule="auto"/></w:pPr></w:pPrDefault>` +
	`</w:docDefaults>` +

	`<w:style w:type="paragraph" w:default="1" w:styleId="Normal"><w:name w:val="Normal"/><w:qFormat/></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Title"><w:name w:val="Title"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:keepNext/><w:spacing w:before="0" w:after="60"/><w:jc w:val="center"/></w:pPr>` +
	`<w:rPr><w:b/><w:bCs/><w:sz w:val="32"/><w:szCs w:val="32"/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Subtitle"><w:name w:val="Subtitle"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:spacing w:after="240"/><w:jc w:val="center"/></w:pPr>` +
	`<w:rPr><w:i/><w:iCs/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Heading1"><w:name w:val="heading 1"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:keepNext/><w:keepLines/><w:spacing w:before="240" w:after="120"/><w:outlineLvl w:val="0"/></w:pPr>` +
	`<w:rPr><w:b/><w:bCs/><w:sz w:val="28"/><w:szCs w:val="28"/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Heading2"><w:name w:val="heading 2"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:keepNext/><w:keepLines/><w:spacing w:before="200" w:after="80"/><w:outlineLvl w:val="1"/></w:pPr>` +
	`<w:rPr><w:b/><w:bCs/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Heading3"><w:name w:val="heading 3"/><w:basedOn w:val="Normal"/><w:next w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:keepNext/><w:keepLines/><w:spacing w:before="160" w:after="60"/><w:outlineLvl w:val="2"/></w:pPr>` +
	`<w:rPr><w:b/><w:bCs/><w:i/><w:iCs/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Quote"><w:name w:val="Quote"/><w:basedOn w:val="Normal"/><w:qFormat/>` +
	`<w:pPr><w:pBdr><w:left w:val="single" w:sz="12" w:space="8" w:color="A0A0A0"/></w:pBdr><w:ind w:left="567"/></w:pPr>` +
	`<w:rPr><w:i/><w:iCs/><w:color w:val="404040"/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Code"><w:name w:val="Code"/><w:basedOn w:val="Normal"/>` +
	`<w:pPr><w:shd w:val="clear" w:color="auto" w:fill="F2F2F2"/><w:spacing w:after="120" w:line="240" w:lineRule="auto"/></w:pPr>` +
	`<w:rPr><w:rFonts w:ascii="Consolas" w:hAnsi="Consolas" w:cs="Consolas"/><w:sz w:val="20"/><w:szCs w:val="20"/></w:rPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="Rule"><w:name w:val="Horizontal Rule"/><w:basedOn w:val="Normal"/>` +
	`<w:pPr><w:pBdr><w:bottom w:val="single" w:sz="6" w:space="1" w:color="A0A0A0"/></w:pBdr><w:spacing w:after="160"/></w:pPr></w:style>` +

	`<w:style w:type="paragraph" w:styleId="TableText"><w:name w:val="Table Text"/><w:basedOn w:val="Normal"/>` +
	`<w:pPr><w:spacing w:before="20" w:after="40" w:line="264" w:lineRule="auto"/></w:pPr></w:style>` +

	`<w:style w:type="character" w:styleId="Hyperlink"><w:name w:val="Hyperlink"/>` +
	`<w:rPr><w:color w:val="0563C1"/><w:u w:val="single"/></w:rPr></w:style>` +

	`<w:style w:type="character" w:styleId="CodeChar"><w:name w:val="Code Char"/>` +
	`<w:rPr><w:rFonts w:ascii="Consolas" w:hAnsi="Consolas" w:cs="Consolas"/><w:sz w:val="22"/><w:szCs w:val="22"/>` +
	`<w:shd w:val="clear" w:color="auto" w:fill="F2F2F2"/></w:rPr></w:style>` +
	`</w:styles>`
