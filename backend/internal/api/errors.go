package api

import (
	"errors"

	"github.com/manhnv/elearning/backend/internal/store"
)

// errValidation tạo lỗi dữ liệu đầu vào ngay tại tầng HTTP.
func errValidation(msg string) error { return &store.InvalidError{Msg: msg} }

// invalidMessage trả về thông báo dành cho người dùng nếu err là vi phạm quy tắc nghiệp vụ.
func invalidMessage(err error) (string, bool) {
	var ie *store.InvalidError
	if errors.As(err, &ie) {
		return ie.Msg, true
	}
	return "", false
}
