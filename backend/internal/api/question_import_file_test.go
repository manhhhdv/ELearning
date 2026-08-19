package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestParseImportFileUnsupportedExtension(t *testing.T) {
	_, errMsg := parseImportFile(strings.NewReader("x"), "cauhoi.docx")
	if errMsg == "" {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestParseCSVRowsCommaDelimited(t *testing.T) {
	csvText := "Mã,Nội dung,Điểm,PA1,PA2,PA3,PA4,Đáp án đúng,Giải thích\n" +
		`C10,"Câu hỏi có, dấu phẩy",2,0-5,10-15,,,1,Giải thích A` + "\n"

	rows, errMsg := parseCSVRows(strings.NewReader(csvText))
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d: %v", len(rows), rows)
	}
	if rows[1][1] != "Câu hỏi có, dấu phẩy" {
		t.Fatalf("quoted comma not preserved: %q", rows[1][1])
	}
}

func TestParseCSVRowsSemicolonDelimited(t *testing.T) {
	// Kiểu Việt Nam: dấu phẩy dùng làm ký tự thập phân nên CSV xuất ra dùng ';'.
	csvText := "Mã;Nội dung;Điểm;PA1;PA2;PA3;PA4;Đáp án đúng;Giải thích\n" +
		"C10;Nhiệt độ bảo quản;2,5;0-5;10-15;;;1;Theo QCVN\n"

	rows, errMsg := parseCSVRows(strings.NewReader(csvText))
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if len(rows) != 2 || rows[1][2] != "2,5" {
		t.Fatalf("unexpected rows: %v", rows)
	}
}

func TestParseCSVRowsStripsBOM(t *testing.T) {
	csvText := "\uFEFFMã;Nội dung;Giải thích\nC10;Câu hỏi;Giải thích A\n"

	rows, errMsg := parseCSVRows(strings.NewReader(csvText))
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if rows[0][0] != "Mã" {
		t.Fatalf("BOM not stripped from first cell: %q", rows[0][0])
	}
}

func TestParseXLSXRowsRoundTrip(t *testing.T) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	_ = f.SetCellValue(sheet, "A1", "Mã")
	_ = f.SetCellValue(sheet, "B1", "Nội dung")
	_ = f.SetCellValue(sheet, "A2", "C10")
	_ = f.SetCellValue(sheet, "B2", "Nhiệt độ bảo quản lạnh là?")

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		t.Fatalf("failed to build fixture xlsx: %v", err)
	}

	rows, errMsg := parseXLSXRows(bytes.NewReader(buf.Bytes()))
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}
	if len(rows) != 2 || rows[1][0] != "C10" || rows[1][1] != "Nhiệt độ bảo quản lạnh là?" {
		t.Fatalf("unexpected rows: %v", rows)
	}
}

func TestParseXLSXRowsInvalidFile(t *testing.T) {
	_, errMsg := parseXLSXRows(strings.NewReader("not an xlsx file"))
	if errMsg == "" {
		t.Fatal("expected error for invalid xlsx content")
	}
}

func TestDetectCSVDelimiter(t *testing.T) {
	if got := detectCSVDelimiter("a;b;c\n1;2;3"); got != ';' {
		t.Fatalf("expected ';', got %q", got)
	}
	if got := detectCSVDelimiter("a,b,c\n1,2,3"); got != ',' {
		t.Fatalf("expected ',', got %q", got)
	}
}

func TestRowsToTabText(t *testing.T) {
	rows := [][]string{
		{"Mã", "Nội dung"},
		{"", ""},
		{"C10", "Câu hỏi 1"},
		{"C11", "Có\ttab lẫn vào"},
	}
	got := rowsToTabText(rows)
	want := "Mã\tNội dung\nC10\tCâu hỏi 1\nC11\tCó tab lẫn vào"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
