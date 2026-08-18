package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword băm mật khẩu bằng bcrypt với cost mặc định.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("băm mật khẩu: %w", err)
	}
	return string(b), nil
}

// CheckPassword so khớp mật khẩu người dùng nhập với hash đã lưu.
func CheckPassword(hash, plain string) bool {
	if hash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

const passwordAlphabet = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789@#$%"

// GeneratePassword sinh mật khẩu ngẫu nhiên để admin cấp cho người dùng mới.
func GeneratePassword(length int) (string, error) {
	if length < 8 {
		length = 12
	}
	var sb strings.Builder
	limit := big.NewInt(int64(len(passwordAlphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", fmt.Errorf("sinh mật khẩu ngẫu nhiên: %w", err)
		}
		sb.WriteByte(passwordAlphabet[n.Int64()])
	}
	return sb.String(), nil
}

// ValidatePassword kiểm tra độ mạnh tối thiểu của mật khẩu.
func ValidatePassword(p string) error {
	if len(p) < 8 {
		return fmt.Errorf("mật khẩu phải có ít nhất 8 ký tự")
	}
	var hasLetter, hasDigit bool
	for _, r := range p {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			hasLetter = true
		}
	}
	if !hasLetter || !hasDigit {
		return fmt.Errorf("mật khẩu phải gồm cả chữ và số")
	}
	return nil
}
