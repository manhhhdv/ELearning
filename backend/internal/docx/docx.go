// Package docx trích xuất văn bản thuần từ tệp Word (.docx) để đưa cho lớp
// Service AI đọc, tương tự cách hệ thống đưa thẳng PDF cho Gemini.
//
// Gemini không đọc hiểu trực tiếp định dạng .docx như PDF, nên hướng đi ở
// đây là tự bóc tách: .docx thực chất là một tệp zip (chuẩn OOXML) chứa các
// tệp XML, nội dung chính nằm ở word/document.xml. Bóc xong, văn bản thuần
// được đưa vào cùng luồng "nội dung tham khảo" mà giáo viên có thể dán tay,
// không cần thêm phụ thuộc thư viện ngoài.
package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNotDocx báo tệp không phải .docx hợp lệ (không mở được như zip OOXML,
// hoặc thiếu word/document.xml).
var ErrNotDocx = errors.New("tệp không đúng định dạng Word (.docx)")

// ErrEmpty báo tài liệu mở được nhưng không có nội dung văn bản nào.
var ErrEmpty = errors.New("tài liệu không có nội dung văn bản")

// documentXMLPath là đường dẫn cố định của nội dung chính trong một tệp .docx.
const documentXMLPath = "word/document.xml"

// ExtractText đọc toàn bộ tệp .docx trong bộ nhớ và trả về văn bản thuần,
// giữ mỗi đoạn văn (w:p) trên các dòng cách nhau bởi một dòng trống.
func ExtractText(raw []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return "", ErrNotDocx
	}

	var docFile *zip.File
	for _, f := range zr.File {
		if f.Name == documentXMLPath {
			docFile = f
			break
		}
	}
	if docFile == nil {
		return "", ErrNotDocx
	}

	rc, err := docFile.Open()
	if err != nil {
		return "", fmt.Errorf("đọc tài liệu Word: %w", err)
	}
	defer rc.Close()

	xmlBytes, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("đọc tài liệu Word: %w", err)
	}

	text, err := textFromDocumentXML(xmlBytes)
	if err != nil {
		return "", fmt.Errorf("phân tích tài liệu Word: %w", err)
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ErrEmpty
	}
	return text, nil
}

// textFromDocumentXML duyệt tuần tự các token XML của word/document.xml.
// Chỉ quan tâm bốn phần tử: w:t (đoạn chữ), w:tab (dấu tab), w:br/w:cr
// (ngắt dòng trong đoạn) và kết thúc w:p (ngắt đoạn văn). Dùng .Local thay vì
// so cả tên có tiền tố nên không phụ thuộc tiền tố namespace cụ thể ("w").
func textFromDocumentXML(raw []byte) (string, error) {
	dec := xml.NewDecoder(bytes.NewReader(raw))
	var sb strings.Builder
	inText := false

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "t":
				inText = true
			case "tab":
				sb.WriteByte('\t')
			case "br", "cr":
				sb.WriteByte('\n')
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "t":
				inText = false
			case "p":
				sb.WriteString("\n\n")
			}
		case xml.CharData:
			if inText {
				sb.Write(t)
			}
		}
	}
	return sb.String(), nil
}
