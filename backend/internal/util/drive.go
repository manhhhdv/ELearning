// Package util chứa các hàm tiện ích dùng chung.
package util

import (
	"net/url"
	"regexp"
	"strings"
)

// Bắt ID file trong các dạng link Google Drive / Docs phổ biến.
var driveIDPatterns = []*regexp.Regexp{
	regexp.MustCompile(`/file/d/([a-zA-Z0-9_-]{10,})`),
	regexp.MustCompile(`/presentation/d/([a-zA-Z0-9_-]{10,})`),
	regexp.MustCompile(`/document/d/([a-zA-Z0-9_-]{10,})`),
	regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9_-]{10,})`),
	regexp.MustCompile(`/d/([a-zA-Z0-9_-]{10,})`),
	regexp.MustCompile(`[?&]id=([a-zA-Z0-9_-]{10,})`),
}

var bareIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{10,}$`)

// ExtractDriveID lấy ID file từ một link chia sẻ Google Drive, hoặc trả lại chính chuỗi
// nếu người dùng đã dán sẵn ID.
func ExtractDriveID(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if bareIDPattern.MatchString(input) && !strings.Contains(input, "/") {
		return input
	}
	for _, re := range driveIDPatterns {
		if m := re.FindStringSubmatch(input); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}

// BuildEmbedURL dựng URL nhúng iframe phù hợp với loại nội dung.
// Nếu người dùng dán một link không thuộc Google Drive, link đó được giữ nguyên.
func BuildEmbedURL(contentType, input string) (driveID, embedURL string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ""
	}

	id := ExtractDriveID(input)
	if id == "" {
		// Không phải link Drive: chấp nhận link http(s) bất kỳ (YouTube, Vimeo, trang nội bộ...).
		if u, err := url.Parse(input); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
			return "", input
		}
		return "", ""
	}

	switch contentType {
	case "slide":
		if strings.Contains(input, "docs.google.com/presentation") || !strings.Contains(input, "http") {
			return id, "https://docs.google.com/presentation/d/" + id + "/embed?start=false&loop=false&delayms=5000"
		}
		return id, "https://drive.google.com/file/d/" + id + "/preview"
	case "document":
		if strings.Contains(input, "docs.google.com/document") {
			return id, "https://docs.google.com/document/d/" + id + "/preview"
		}
		if strings.Contains(input, "docs.google.com/spreadsheets") {
			return id, "https://docs.google.com/spreadsheets/d/" + id + "/preview"
		}
		return id, "https://drive.google.com/file/d/" + id + "/preview"
	default:
		// video, pdf, link: trình xem sẵn có của Drive xử lý được tất cả.
		return id, "https://drive.google.com/file/d/" + id + "/preview"
	}
}
