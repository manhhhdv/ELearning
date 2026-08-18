package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type ctxKey int

const claimsKey ctxKey = iota

// ErrorWriter cho phép package api quyết định định dạng thân phản hồi lỗi.
type ErrorWriter func(w http.ResponseWriter, status int, message string)

// Middleware kiểm tra JWT trên mỗi request và gắn claims vào context.
type Middleware struct {
	tokens   *TokenManager
	writeErr ErrorWriter
}

func NewMiddleware(tokens *TokenManager, writeErr ErrorWriter) *Middleware {
	return &Middleware{tokens: tokens, writeErr: writeErr}
}

// RequireAuth chặn request không có token hợp lệ.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := bearerToken(r)
		if raw == "" {
			m.writeErr(w, http.StatusUnauthorized, "Thiếu token đăng nhập")
			return
		}
		claims, err := m.tokens.Parse(raw)
		if err != nil {
			m.writeErr(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ hoặc đã hết hạn")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), claimsKey, claims)))
	})
}

// RequireRole chặn request của người dùng không thuộc các vai trò cho phép.
func (m *Middleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := FromContext(r.Context())
			if !ok {
				m.writeErr(w, http.StatusUnauthorized, "Thiếu token đăng nhập")
				return
			}
			if !allowed[claims.Role] {
				m.writeErr(w, http.StatusForbidden, "Bạn không có quyền thực hiện thao tác này")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// FromContext lấy claims đã được RequireAuth gắn vào context.
func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}

// UserID trả về ID người dùng hiện tại, uuid.Nil nếu chưa đăng nhập.
func UserID(ctx context.Context) uuid.UUID {
	if claims, ok := FromContext(ctx); ok {
		return claims.UserID
	}
	return uuid.Nil
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}
	prefix, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(prefix, "Bearer") {
		return ""
	}
	return strings.TrimSpace(token)
}
