package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

// unzipParts mở tệp .docx và trả về nội dung từng phần, đồng thời kiểm tra
// mọi phần XML đều hợp lệ — một thẻ lệch là Word báo tệp hỏng.
func unzipParts(t *testing.T, raw []byte) map[string]string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	parts := map[string]string{}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		b, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read %s: %v", f.Name, err)
		}
		dec := xml.NewDecoder(bytes.NewReader(b))
		for {
			if _, err := dec.Token(); err == io.EOF {
				break
			} else if err != nil {
				t.Fatalf("%s is not well-formed XML: %v", f.Name, err)
			}
		}
		parts[f.Name] = string(b)
	}
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "word/document.xml", "word/styles.xml", "word/numbering.xml"} {
		if _, ok := parts[name]; !ok {
			t.Fatalf("missing part %s", name)
		}
	}
	return parts
}

func build(t *testing.T, fill func(d *Document)) (map[string]string, string) {
	t.Helper()
	d := New("Thử nghiệm")
	fill(d)
	raw, err := d.Bytes()
	if err != nil {
		t.Fatalf("Bytes: %v", err)
	}
	parts := unzipParts(t, raw)
	text, err := ExtractText(raw)
	if err != nil {
		t.Fatalf("ExtractText: %v", err)
	}
	return parts, text
}

func TestWriterRoundTripsText(t *testing.T) {
	_, text := build(t, func(d *Document) {
		d.Body().Text(StyleTitle, "Định luật Ôm")
		d.Body().Add(Para{Runs: []Run{{Text: "I = U/R", Bold: true}, {Text: " & <R> \"điện trở\""}}})
	})
	for _, want := range []string{"Định luật Ôm", "I = U/R", `& <R> "điện trở"`} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q missing %q", text, want)
		}
	}
}

func TestWriterDropsInvalidXMLChars(t *testing.T) {
	_, text := build(t, func(d *Document) {
		d.Body().Text("", "a\x00b\x0bc\x1fd")
	})
	if !strings.Contains(text, "abcd") {
		t.Errorf("control characters should be dropped, got %q", text)
	}
}

func TestWriterLineBreaksAndTabs(t *testing.T) {
	parts, _ := build(t, func(d *Document) {
		d.Body().Text("", "dòng 1\ndòng 2\tcột")
	})
	doc := parts["word/document.xml"]
	if !strings.Contains(doc, "<w:br/>") || !strings.Contains(doc, "<w:tab/>") {
		t.Errorf("expected line break and tab in %s", doc)
	}
}

func TestWriterHyperlinksOnlyForSafeSchemes(t *testing.T) {
	parts, _ := build(t, func(d *Document) {
		d.Body().Add(Para{Runs: []Run{
			{Text: "an toàn", Link: "https://example.com/?a=1&b=2"},
			{Text: "nguy hiểm", Link: "javascript:alert(1)"},
		}})
	})
	rels := parts["word/_rels/document.xml.rels"]
	if !strings.Contains(rels, `Target="https://example.com/?a=1&amp;b=2" TargetMode="External"`) {
		t.Errorf("missing https relationship: %s", rels)
	}
	if strings.Contains(rels, "javascript") {
		t.Errorf("javascript: link must not be written: %s", rels)
	}
	if n := strings.Count(parts["word/document.xml"], "<w:hyperlink"); n != 1 {
		t.Errorf("want 1 hyperlink, got %d", n)
	}
}

func TestWriterListsRestartNumbering(t *testing.T) {
	parts, _ := build(t, func(d *Document) {
		a := d.NewList(true, 1)
		b := d.NewList(true, 3)
		d.Body().Add(Para{List: a, Runs: []Run{{Text: "một"}}})
		d.Body().Add(Para{List: b, Runs: []Run{{Text: "ba"}}})
	})
	num := parts["word/numbering.xml"]
	for _, want := range []string{
		`<w:num w:numId="1"><w:abstractNumId w:val="1"/><w:lvlOverride w:ilvl="0"><w:startOverride w:val="1"/>`,
		`<w:num w:numId="2"><w:abstractNumId w:val="1"/><w:lvlOverride w:ilvl="0"><w:startOverride w:val="3"/>`,
	} {
		if !strings.Contains(num, want) {
			t.Errorf("numbering.xml missing %s", want)
		}
	}
}

func TestWriterTableCellsAlwaysEndWithParagraph(t *testing.T) {
	parts, text := build(t, func(d *Document) {
		outer := d.Body().NewTable(1, 1)
		row := outer.Row()
		row[0].Text("", "ô trái")
		inner := row[1].NewTable(1)
		inner.Row()[0].Text("", "bảng lồng")
		row[1].AddTable(inner)
		d.Body().AddTable(outer)
	})
	doc := parts["word/document.xml"]
	if !strings.Contains(doc, "</w:tbl><w:p><w:pPr><w:pStyle w:val=\"TableText\"/></w:pPr></w:p></w:tc>") {
		t.Errorf("cell ending with a nested table needs a trailing paragraph: %s", doc)
	}
	if !strings.Contains(doc, "</w:tbl><w:p/><w:sectPr>") {
		t.Errorf("body ending with a table needs a paragraph before sectPr")
	}
	if !strings.Contains(text, "ô trái") || !strings.Contains(text, "bảng lồng") {
		t.Errorf("table text lost: %q", text)
	}
}

func TestMarkdownBlocks(t *testing.T) {
	src := "# Bài 1\n\nĐoạn **đậm** và *nghiêng*, [liên kết](https://example.com).\n\n" +
		"- ý một\n- ý hai\n  - ý con\n\n1. bước một\n2. bước hai\n\n" +
		"> trích dẫn\n\n```\ncode $x$ **không đổi**\n```\n\n| A | B |\n|---|--:|\n| 1 | 2 |\n\n---\n\n- [x] đã xong\n"
	parts, text := build(t, func(d *Document) { d.Body().Markdown(src) })
	doc := parts["word/document.xml"]

	checks := []string{
		`<w:pStyle w:val="Heading1"/>`,
		`<w:b/><w:bCs/></w:rPr><w:t xml:space="preserve">đậm</w:t>`,
		`<w:i/><w:iCs/></w:rPr><w:t xml:space="preserve">nghiêng</w:t>`,
		`<w:numPr><w:ilvl w:val="1"/>`,
		`<w:pStyle w:val="Quote"/>`,
		`<w:pStyle w:val="Code"/>`,
		`code $x$ **không đổi**`,
		`<w:tbl>`,
		`<w:jc w:val="right"/>`,
		`<w:pStyle w:val="Rule"/>`,
		`☑ đã xong`,
	}
	for _, want := range checks {
		if !strings.Contains(doc, want) {
			t.Errorf("document.xml missing %s", want)
		}
	}
	if strings.Contains(text, "**") && !strings.Contains(text, "**không đổi**") {
		t.Errorf("markdown markers leaked into text: %q", text)
	}
	if !strings.Contains(parts["word/_rels/document.xml.rels"], "https://example.com") {
		t.Errorf("link relationship missing")
	}
}

func TestMarkdownMathBecomesEquations(t *testing.T) {
	src := "Ta có $a_1 * b_2$ và giá $5 và $7.\n\n$$\n\\frac{a}{b}\n  = \\sqrt{x}\n$$\n\n" +
		"`$không phải công thức$` và \\$ thoát. Lệnh lạ $\\foo{x}$ giữ nguyên."
	parts, text := build(t, func(d *Document) { d.Body().Markdown(src) })
	doc := parts["word/document.xml"]

	for _, want := range []string{
		// Công thức trong dòng thành phương trình, chỉ số dưới dựng bằng sSub.
		`<m:oMath><m:sSub><m:e>`,
		// Công thức khối đứng riêng một đoạn thành oMathPara.
		`<m:oMathPara><m:oMath><m:f><m:num>`,
		`<m:rad><m:radPr><m:degHide m:val="1"/></m:radPr>`,
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("document.xml missing %s", want)
		}
	}
	for _, want := range []string{"giá $5 và $7", "$không phải công thức$", "$ thoát", `$\foo{x}$ giữ nguyên`} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q missing %q", text, want)
		}
	}
	if strings.ContainsRune(text, mathOpen) || strings.ContainsRune(text, mathClose) {
		t.Errorf("math placeholder leaked: %q", text)
	}
}

func TestMarkdownAlert(t *testing.T) {
	parts, text := build(t, func(d *Document) {
		d.Body().Markdown("> [!NOTE]\n> Dấu của $\\Delta$ quyết định số nghiệm.")
	})
	if strings.Contains(text, "[!NOTE]") || !strings.Contains(text, "Ghi chú") || !strings.Contains(text, "quyết định số nghiệm") {
		t.Errorf("alert not converted: %q", text)
	}
	if !strings.Contains(parts["word/document.xml"], "<m:oMath>") {
		t.Errorf("math inside alert should stay an equation")
	}
}

func TestMarkdownMediaBecomesLabelledLink(t *testing.T) {
	src := "![Sơ đồ mạch](https://example.com/mach.png)\n\n![](https://youtu.be/abc)"
	parts, text := build(t, func(d *Document) { d.Body().Markdown(src) })
	if !strings.Contains(text, "[Hình ảnh: Sơ đồ mạch]") || !strings.Contains(text, "[Video]") {
		t.Errorf("media labels missing: %q", text)
	}
	if !strings.Contains(parts["word/_rels/document.xml.rels"], "https://youtu.be/abc") {
		t.Errorf("media link missing")
	}
}

func TestMarkdownInlineHTML(t *testing.T) {
	parts, text := build(t, func(d *Document) {
		d.Body().Markdown("H<sub>2</sub>O và x<sup>2</sup>, <u>gạch chân</u><br>dòng mới")
	})
	doc := parts["word/document.xml"]
	for _, want := range []string{`<w:vertAlign w:val="subscript"/>`, `<w:vertAlign w:val="superscript"/>`, `<w:u w:val="single"/>`} {
		if !strings.Contains(doc, want) {
			t.Errorf("document.xml missing %s", want)
		}
	}
	if strings.Contains(text, "<u>") || !strings.Contains(text, "gạch chân") {
		t.Errorf("inline html not converted: %q", text)
	}
}
