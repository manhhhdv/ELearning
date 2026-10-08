package docx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

// buildDocx dựng một tệp .docx tối giản trong bộ nhớ: chỉ cần word/document.xml
// vì đó là phần duy nhất ExtractText đọc tới.
func buildDocx(t *testing.T, documentXML string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(documentXMLPath)
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	if _, err := w.Write([]byte(documentXML)); err != nil {
		t.Fatalf("write zip entry: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func TestExtractTextParagraphsAndTabs(t *testing.T) {
	xmlDoc := `<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Định luật Ôm</w:t></w:r></w:p>
    <w:p><w:r><w:t>Cường độ</w:t></w:r><w:r><w:tab/></w:r><w:r><w:t>Điện áp</w:t></w:r></w:p>
  </w:body>
</w:document>`

	text, err := ExtractText(buildDocx(t, xmlDoc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(text, "Định luật Ôm") {
		t.Fatalf("missing first paragraph: %q", text)
	}
	if !strings.Contains(text, "Cường độ\tĐiện áp") {
		t.Fatalf("expected tab-separated run text, got %q", text)
	}
	// Hai đoạn văn phải cách nhau bởi dòng trống.
	if !strings.Contains(text, "Định luật Ôm\n\nCường độ") {
		t.Fatalf("expected paragraph break, got %q", text)
	}
}

func TestExtractTextEmptyDocument(t *testing.T) {
	xmlDoc := `<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body><w:p><w:r><w:t>   </w:t></w:r></w:p></w:body>
</w:document>`

	_, err := ExtractText(buildDocx(t, xmlDoc))
	if err != ErrEmpty {
		t.Fatalf("expected ErrEmpty, got %v", err)
	}
}

func TestExtractTextNotDocx(t *testing.T) {
	_, err := ExtractText([]byte("not a zip file at all"))
	if err != ErrNotDocx {
		t.Fatalf("expected ErrNotDocx, got %v", err)
	}
}

func TestExtractTextMissingDocumentXML(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("word/other.xml")
	_, _ = w.Write([]byte("<x/>"))
	_ = zw.Close()

	_, err := ExtractText(buf.Bytes())
	if err != ErrNotDocx {
		t.Fatalf("expected ErrNotDocx, got %v", err)
	}
}
