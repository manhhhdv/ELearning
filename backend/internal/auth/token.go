package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims là payload của access token.
type Claims struct {
	UserID uuid.UUID `json:"uid"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager phát hành và kiểm tra JWT.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

func (m *TokenManager) TTL() time.Duration { return m.ttl }

func (m *TokenManager) Issue(userID uuid.UUID, email, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.ttl)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("ký token: %w", err)
	}
	return signed, expiresAt, nil
}

func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("thuật toán ký không được hỗ trợ: %v", t.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("token không hợp lệ: %w", err)
	}
	return claims, nil
}
