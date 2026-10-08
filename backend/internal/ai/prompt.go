package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Mẫu prompt của hệ thống gồm năm phần (mục 3.4.6): Vai trò, Bối cảnh,
// Nhiệm vụ, Ràng buộc, Định dạng đầu ra. Người dùng chỉ chọn thông tin trong
// ngoặc vuông, hệ thống tự ghép thành prompt hoàn chỉnh.
type promptBuilder struct {
	role        string
	context     []string
	task        []string
	constraints []string
	format      string
}

func (b *promptBuilder) addContext(label, value string) {
	if v := strings.TrimSpace(value); v != "" {
		b.context = append(b.context, fmt.Sprintf("- %s: %s", label, v))
	}
}

func (b *promptBuilder) addTask(format string, args ...any) {
	b.task = append(b.task, "- "+fmt.Sprintf(format, args...))
}

func (b *promptBuilder) String() string {
	var sb strings.Builder
	sb.WriteString("VAI TRÒ:\n")
	sb.WriteString(b.role)
	if len(b.context) > 0 {
		sb.WriteString("\n\nBỐI CẢNH:\n")
		sb.WriteString(strings.Join(b.context, "\n"))
	}
	if len(b.task) > 0 {
		sb.WriteString("\n\nNHIỆM VỤ:\n")
		sb.WriteString(strings.Join(b.task, "\n"))
	}
	if len(b.constraints) > 0 {
		sb.WriteString("\n\nRÀNG BUỘC:\n")
		for _, c := range b.constraints {
			sb.WriteString("- " + c + "\n")
		}
	}
	if b.format != "" {
		sb.WriteString("\nĐỊNH DẠNG ĐẦU RA:\n")
		sb.WriteString(b.format)
	}
	return strings.TrimSpace(sb.String())
}

// extractJSON lấy khối JSON đầu tiên trong câu trả lời. Dù đã yêu cầu
// responseMimeType=application/json, model vẫn có lúc bọc kết quả trong ```json.
func extractJSON(raw string, dst any) error {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```json")
		s = strings.TrimPrefix(s, "```")
		if i := strings.LastIndex(s, "```"); i >= 0 {
			s = s[:i]
		}
		s = strings.TrimSpace(s)
	}
	// Cắt về đúng cặp ngoặc ngoài cùng nếu model kèm thêm lời dẫn.
	if start := strings.IndexAny(s, "{["); start > 0 {
		s = s[start:]
	}
	if end := strings.LastIndexAny(s, "}]"); end >= 0 && end < len(s)-1 {
		s = s[:end+1]
	}
	if err := json.Unmarshal([]byte(s), dst); err != nil {
		return fmt.Errorf("kết quả AI không đúng cấu trúc JSON: %w", err)
	}
	return nil
}
