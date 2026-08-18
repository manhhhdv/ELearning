package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
)

const userInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"

// GoogleProfile là thông tin tài khoản lấy về từ Google sau khi đăng nhập.
type GoogleProfile struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	HostedDomain  string `json:"hd"`
}

// GoogleAuthenticator bọc luồng OAuth2 Authorization Code của Google.
type GoogleAuthenticator struct {
	cfg    *oauth2.Config
	secret []byte
}

func NewGoogleAuthenticator(clientID, clientSecret, redirectURL, stateSecret string) *GoogleAuthenticator {
	return &GoogleAuthenticator{
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     googleoauth.Endpoint,
		},
		secret: []byte(stateSecret),
	}
}

// NewState sinh chuỗi state có chữ ký và hạn dùng để chống CSRF trên luồng OAuth.
func (g *GoogleAuthenticator) NewState() (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("sinh nonce: %w", err)
	}
	payload := fmt.Sprintf("%s.%d", hex.EncodeToString(nonce), time.Now().Add(10*time.Minute).Unix())
	return payload + "." + g.sign(payload), nil
}

// VerifyState kiểm tra chữ ký và hạn dùng của state nhận lại từ Google.
func (g *GoogleAuthenticator) VerifyState(state string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return fmt.Errorf("state không hợp lệ")
	}
	payload := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(g.sign(payload)), []byte(parts[2])) {
		return fmt.Errorf("chữ ký state không khớp")
	}
	var exp int64
	if _, err := fmt.Sscanf(parts[1], "%d", &exp); err != nil {
		return fmt.Errorf("state không hợp lệ")
	}
	if time.Now().Unix() > exp {
		return fmt.Errorf("state đã hết hạn, vui lòng đăng nhập lại")
	}
	return nil
}

func (g *GoogleAuthenticator) sign(payload string) string {
	mac := hmac.New(sha256.New, g.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// AuthCodeURL trả về URL để chuyển hướng người dùng sang trang đăng nhập Google.
func (g *GoogleAuthenticator) AuthCodeURL(state string) string {
	return g.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))
}

// Exchange đổi authorization code lấy access token rồi đọc hồ sơ người dùng.
func (g *GoogleAuthenticator) Exchange(ctx context.Context, code string) (*GoogleProfile, error) {
	token, err := g.cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("đổi mã xác thực Google thất bại: %w", err)
	}

	client := g.cfg.Client(ctx, token)
	client.Timeout = 10 * time.Second
	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("gọi userinfo Google thất bại: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo Google trả về mã %d", resp.StatusCode)
	}

	var profile GoogleProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("đọc userinfo Google: %w", err)
	}
	if profile.Sub == "" || profile.Email == "" {
		return nil, fmt.Errorf("Google không trả về email của tài khoản")
	}
	if !profile.EmailVerified {
		return nil, fmt.Errorf("email Google chưa được xác minh")
	}
	profile.Email = strings.ToLower(profile.Email)
	return &profile, nil
}
