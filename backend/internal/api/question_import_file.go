package api

import (
	"encoding/csv"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"

	"github.com/manhnv/elearning/backend/internal/models"
)

// maxImportFileSize giới hạn dung lượng file .xlsx/.csv tải lên để nhập câu hỏi hàng loạt.
const maxImportFileSize = 5 << 20 // 5MB

type importQuestionsFileResponse struct {
	Text string `json:"text"`
}

// handleImportQuestionsFile nhận file .xlsx hoặc .csv, đọc thành văn bản dạng bảng
// (mỗi ô cách nhau bằng Tab) để tái dùng đúng logic phân tích và xem trước ở
// handleImportQuestions/parseTable phía giao diện — chỉ khác chỗ lấy dữ liệu.
func (s *Server) handleImportQuestionsFile(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, _, ok := s.loadNode(w, r, nodeID, true)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment {
		writeError(w, http.StatusBadRequest, "Chỉ có thể nhập câu hỏi vào một bài tập")
		return
	}

	rows, errMsg := readUploadedTable(w, r)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	writeJSON(w, http.StatusOK, importQuestionsFileResponse{Text: rowsToTabText(rows)})
}

// readUploadedTable nhận file .xlsx/.csv từ form multipart và trả về các hàng dữ liệu.
func readUploadedTable(w http.ResponseWriter, r *http.Request) ([][]string, string) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImportFileSize)
	if err := r.ParseMultipartForm(maxImportFileSize); err != nil {
		return nil, "File quá lớn hoặc không đọc được (tối đa 5MB)"
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, "Vui lòng chọn một file để tải lên"
	}
	defer file.Close()

	rows, errMsg := parseImportFile(file, header.Filename)
	if errMsg != "" {
		return nil, errMsg
	}
	if len(rows) == 0 {
		return nil, "File không có dữ liệu"
	}
	return rows, ""
}

// parseImportFile đọc file theo phần mở rộng, trả về các hàng dữ liệu dạng bảng.
func parseImportFile(file io.Reader, filename string) ([][]string, string) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".xlsx":
		return parseXLSXRows(file)
	case ".csv":
		return parseCSVRows(file)
	default:
		return nil, "Chỉ hỗ trợ file .xlsx hoặc .csv"
	}
}

func parseXLSXRows(file io.Reader) ([][]string, string) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, "Không đọc được file .xlsx. Vui lòng kiểm tra lại định dạng."
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, "File .xlsx không có trang tính nào"
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, "Không đọc được dữ liệu trong file .xlsx"
	}
	return rows, ""
}

func parseCSVRows(file io.Reader) ([][]string, string) {
	raw, err := io.ReadAll(file)
	if err != nil {
		return nil, "Không đọc được file .csv"
	}
	text := strings.TrimPrefix(string(raw), "\uFEFF") // bỏ BOM nếu có
	if !utf8.ValidString(text) {
		return nil, "File .csv cần được lưu ở định dạng UTF-8"
	}

	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = detectCSVDelimiter(text)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, "Không đọc được dữ liệu trong file .csv"
	}
	return rows, ""
}

// detectCSVDelimiter đoán dấu phân cách: nhiều file CSV xuất từ Excel/Sheets ở
// vùng Việt Nam dùng dấu chấm phẩy vì dấu phẩy là ký tự thập phân.
func detectCSVDelimiter(text string) rune {
	firstLine := text
	if idx := strings.IndexAny(text, "\r\n"); idx >= 0 {
		firstLine = text[:idx]
	}
	if strings.Count(firstLine, ";") > strings.Count(firstLine, ",") {
		return ';'
	}
	return ','
}

// rowsToTabText nối các hàng thành văn bản Tab-separated, bỏ hàng trống hoàn toàn.
func rowsToTabText(rows [][]string) string {
	var b strings.Builder
	for _, row := range rows {
		empty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}
		for i, cell := range row {
			if i > 0 {
				b.WriteByte('\t')
			}
			b.WriteString(strings.ReplaceAll(strings.TrimSpace(cell), "\t", " "))
		}
		b.WriteByte('\n')
	}
	return strings.TrimSuffix(b.String(), "\n")
}
