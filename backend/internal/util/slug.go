package util

import (
	"regexp"
	"strings"
)

// Bảng chuyển ký tự có dấu tiếng Việt sang không dấu, dùng cho slug URL.
var vietnameseMap = map[rune]rune{
	'à': 'a', 'á': 'a', 'ạ': 'a', 'ả': 'a', 'ã': 'a',
	'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ậ': 'a', 'ẩ': 'a', 'ẫ': 'a',
	'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ặ': 'a', 'ẳ': 'a', 'ẵ': 'a',
	'è': 'e', 'é': 'e', 'ẹ': 'e', 'ẻ': 'e', 'ẽ': 'e',
	'ê': 'e', 'ề': 'e', 'ế': 'e', 'ệ': 'e', 'ể': 'e', 'ễ': 'e',
	'ì': 'i', 'í': 'i', 'ị': 'i', 'ỉ': 'i', 'ĩ': 'i',
	'ò': 'o', 'ó': 'o', 'ọ': 'o', 'ỏ': 'o', 'õ': 'o',
	'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ộ': 'o', 'ổ': 'o', 'ỗ': 'o',
	'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ợ': 'o', 'ở': 'o', 'ỡ': 'o',
	'ù': 'u', 'ú': 'u', 'ụ': 'u', 'ủ': 'u', 'ũ': 'u',
	'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ự': 'u', 'ử': 'u', 'ữ': 'u',
	'ỳ': 'y', 'ý': 'y', 'ỵ': 'y', 'ỷ': 'y', 'ỹ': 'y',
	'đ': 'd',
}

var slugInvalidChars = regexp.MustCompile(`[^a-z0-9]+`)
var slugTrimDash = regexp.MustCompile(`^-+|-+$`)

// Slugify chuyển một chuỗi tiếng Việt bất kỳ thành slug an toàn cho URL:
// chữ thường, không dấu, chỉ gồm a-z0-9 và dấu gạch ngang.
func Slugify(input string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(input) {
		if repl, ok := vietnameseMap[r]; ok {
			b.WriteRune(repl)
		} else {
			b.WriteRune(r)
		}
	}
	s := slugInvalidChars.ReplaceAllString(b.String(), "-")
	s = slugTrimDash.ReplaceAllString(s, "")
	if s == "" {
		return "muc"
	}
	return s
}
