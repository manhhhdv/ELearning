package docx

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

// wellFormed bọc phần OMML trong một phần tử gốc có khai báo namespace rồi
// kiểm tra XML hợp lệ.
func wellFormed(t *testing.T, tex, omml string) {
	t.Helper()
	doc := `<root xmlns:m="` + nsM + `" xmlns:w="` + nsW + `"><m:oMath>` + omml + `</m:oMath></root>`
	dec := xml.NewDecoder(strings.NewReader(doc))
	for {
		if _, err := dec.Token(); err == io.EOF {
			return
		} else if err != nil {
			t.Fatalf("%s → invalid XML: %v\n%s", tex, err, omml)
		}
	}
}

func TestTeXToOMML(t *testing.T) {
	cases := []struct {
		tex  string
		want []string
	}{
		{`x^2 + y_1`, []string{`<m:sSup><m:e>`, `<m:sSub><m:e>`}},
		{`a_i^2`, []string{`<m:sSubSup>`}},
		{`\frac{a+b}{2}`, []string{`<m:f><m:num>`, `<m:den>`}},
		{`\dfrac12`, []string{`<m:num><m:r><w:rPr><w:rFonts w:ascii="Cambria Math" w:hAnsi="Cambria Math"/></w:rPr><m:t xml:space="preserve">1</m:t>`}},
		{`\sqrt{x}`, []string{`<m:degHide m:val="1"/>`}},
		{`\sqrt[3]{x}`, []string{`<m:deg><m:r>`}},
		{`\left(-\dfrac{b}{2a};\ -\dfrac{\Delta}{4a}\right)`, []string{`<m:begChr m:val="("/><m:endChr m:val=")"/>`, "Δ", "−"}},
		{`\left.\frac{x}{y}\right|`, []string{`<m:begChr m:val=""/><m:endChr m:val="|"/>`}},
		{`\begin{cases} x + y = 3 \\ x - y = 1 \end{cases}`, []string{`<m:begChr m:val="{"/><m:endChr m:val=""/>`, `<m:eqArr><m:e>`}},
		{`\begin{pmatrix} 1 & 2 \\ 3 & 4 \end{pmatrix}`, []string{`<m:count m:val="2"/>`, `<m:mr><m:e>`}},
		{`\sum_{i=1}^{n} i^2`, []string{`<m:chr m:val="∑"/><m:limLoc m:val="undOvr"/>`, `<m:e><m:sSup>`}},
		{`\int_0^1 x\,dx`, []string{`<m:chr m:val="∫"/><m:limLoc m:val="subSup"/>`}},
		{`\lim_{x \to 0} \frac{\sin x}{x}`, []string{`<m:limLow><m:e>`, "→", `<m:sty m:val="p"/>`}},
		{`x \in \mathbb{R}, \mathcal{E}`, []string{"∈", "ℝ", "ℰ"}},
		{`\vec{F} = m\vec{a}`, []string{`<m:acc><m:accPr><m:chr m:val="` + "⃗" + `"/>`}},
		{`\overline{AB}`, []string{`<m:pos m:val="top"/>`}},
		{`v = 10\text{ m/s}^2`, []string{`<m:sty m:val="p"/></m:rPr>`, " m/s"}},
		{`30^\circ, 2{,}5 \le x \ne 3.14`, []string{"∘", ">5≤x≠3.14<"}},
		// Chữ liền nhau gộp một run; chỉ số chỉ gắn vào ký tự ngay trước nó.
		{`ax^2 + bx`, []string{`>a</m:t></m:r><m:sSup><m:e>`, `>x</m:t>`, `>+bx<`}},
		{`\sin x`, []string{`<m:sty m:val="p"/></m:rPr>`, `>sin</m:t></m:r><m:r>`}},
		{`\mathbf{F}_{net}`, []string{`<m:sty m:val="b"/>`}},
		{`\not\in \emptyset`, []string{"∉", "∅"}},
		{`^{14}_{6}C`, []string{`<m:sSubSup><m:e></m:e>`}},
		{`a < b \quad \% & x`, []string{"&lt;", " ", "%"}},
	}
	for _, c := range cases {
		omml, ok := texToOMML(c.tex)
		if !ok {
			t.Errorf("%s: conversion failed", c.tex)
			continue
		}
		wellFormed(t, c.tex, omml)
		for _, w := range c.want {
			if !strings.Contains(omml, w) {
				t.Errorf("%s: missing %q in\n%s", c.tex, w, omml)
			}
		}
	}
}

func TestTeXToOMMLRejectsUnsupported(t *testing.T) {
	for _, tex := range []string{
		`\foo{x}`,                    // lệnh lạ
		`\frac{a}{b`,                 // thiếu ngoặc đóng
		`a}`,                         // thừa ngoặc đóng
		`\left( x`,                   // thiếu \right
		`\begin{tikzcd}\end{tikzcd}`, // môi trường lạ
		`\begin{cases} x \end{array}`,
		`   `,
	} {
		if omml, ok := texToOMML(tex); ok {
			t.Errorf("%q should be rejected, got %s", tex, omml)
		}
	}
}

func TestMathFallbackKeepsTeX(t *testing.T) {
	parts, text := build(t, func(d *Document) {
		d.Body().Add(Para{Runs: []Run{{Text: "Xem "}, {Math: `\unknown{x}`}}})
		d.Body().Add(Para{Runs: []Run{{Math: `\weird`, Display: true}}})
	})
	if !strings.Contains(text, `Xem $\unknown{x}$`) || !strings.Contains(text, `$$ \weird $$`) {
		t.Errorf("unsupported formula should stay as TeX: %q", text)
	}
	if strings.Contains(parts["word/document.xml"], "<m:oMath") {
		t.Errorf("no equation expected for unsupported formula")
	}
}
