package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// googleConfig là cấu hình Google OAuth có hiệu lực sau khi gộp giá trị lưu
// trong DB (do admin chỉnh qua giao diện) với giá trị mặc định từ .env.
type googleConfig struct {
	ClientID          string
	ClientSecret      string
	AllowedDomains    []string
	AutoProvisionRole string
	// Nguồn cấu hình đang dùng, chỉ để hiển thị cho admin biết.
	Source string // "database" | "env" | "none"
}

// resolveGoogleConfig gộp cấu hình: khoá nào có trong DB thì ưu tiên DB (kể cả
// khi được lưu là chuỗi rỗng — admin có thể chủ động tắt qua giao diện mà
// không cần đụng vào .env), khoá nào không có thì rơi về giá trị .env.
func (s *Server) resolveGoogleConfig(ctx context.Context) (googleConfig, error) {
	saved, err := s.store.GetSettings(ctx, []string{
		store.SettingGoogleClientID, store.SettingGoogleClientSecret,
		store.SettingGoogleAllowedDomains, store.SettingGoogleAutoProvisionRole,
	})
	if err != nil {
		return googleConfig{}, err
	}

	gc := googleConfig{
		ClientID:          s.cfg.GoogleClientID,
		ClientSecret:      s.cfg.GoogleClientSecret,
		AllowedDomains:    s.cfg.GoogleAllowedDomains,
		AutoProvisionRole: s.cfg.GoogleAutoProvisionRole,
		Source:            "env",
	}
	_, hasID := saved[store.SettingGoogleClientID]
	_, hasSecret := saved[store.SettingGoogleClientSecret]
	if hasID || hasSecret {
		gc.Source = "database"
		gc.ClientID = saved[store.SettingGoogleClientID]
		gc.ClientSecret = saved[store.SettingGoogleClientSecret]
		gc.AllowedDomains = splitDomains(saved[store.SettingGoogleAllowedDomains])
		gc.AutoProvisionRole = saved[store.SettingGoogleAutoProvisionRole]
	}
	if gc.ClientID == "" || gc.ClientSecret == "" {
		gc.Source = "none"
	}
	return gc, nil
}

// googleAuthenticator dựng authenticator từ cấu hình hiện có; enabled=false nếu
// chưa đủ Client ID và Secret ở cả DB lẫn .env.
func (s *Server) googleAuthenticator(ctx context.Context) (*auth.GoogleAuthenticator, googleConfig, bool, error) {
	gc, err := s.resolveGoogleConfig(ctx)
	if err != nil {
		return nil, googleConfig{}, false, err
	}
	if gc.ClientID == "" || gc.ClientSecret == "" {
		return nil, gc, false, nil
	}
	// redirectURL là địa chỉ callback của chính máy chủ này — cố định theo triển khai,
	// không phải thứ admin nên tự sửa qua giao diện nên vẫn lấy từ cấu hình khởi động.
	g := auth.NewGoogleAuthenticator(gc.ClientID, gc.ClientSecret, s.cfg.GoogleRedirectURL, s.cfg.JWTSecret)
	return g, gc, true, nil
}

func splitDomains(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// Cấu hình đăng nhập Google — chỉ admin xem/sửa được, vì đây là thông tin nhạy cảm.
// ---------------------------------------------------------------------------

type googleSettingsResponse struct {
	Enabled bool `json:"enabled"`
	// Client ID không phải bí mật (Google coi đây là định danh công khai của ứng dụng).
	ClientID string `json:"clientId"`
	// Không bao giờ trả secret thật ra ngoài — chỉ báo đã có hay chưa.
	HasSecret         bool   `json:"hasSecret"`
	Source            string `json:"source"`
	RedirectURL       string `json:"redirectUrl"`
	AllowedDomains    string `json:"allowedDomains"`
	AutoProvisionRole string `json:"autoProvisionRole"`
}

func (s *Server) handleGetGoogleSettings(w http.ResponseWriter, r *http.Request) {
	gc, err := s.resolveGoogleConfig(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, googleSettingsResponse{
		Enabled:           gc.ClientID != "" && gc.ClientSecret != "",
		ClientID:          gc.ClientID,
		HasSecret:         gc.ClientSecret != "",
		Source:            gc.Source,
		RedirectURL:       s.cfg.GoogleRedirectURL,
		AllowedDomains:    strings.Join(gc.AllowedDomains, ", "),
		AutoProvisionRole: gc.AutoProvisionRole,
	})
}

type saveGoogleSettingsRequest struct {
	// false: xoá cấu hình đã lưu trong hệ thống, quay lại dùng giá trị trong .env (nếu có).
	Enabled bool `json:"enabled"`
	// Bắt buộc khi Enabled=true.
	ClientID string `json:"clientId"`
	// Để trống khi Enabled=true nghĩa là giữ nguyên secret đã lưu trước đó.
	ClientSecret      string `json:"clientSecret"`
	AllowedDomains    string `json:"allowedDomains"`
	AutoProvisionRole string `json:"autoProvisionRole"`
}

func (s *Server) handleSaveGoogleSettings(w http.ResponseWriter, r *http.Request) {
	var req saveGoogleSettingsRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	ctx := r.Context()

	if !req.Enabled {
		for _, key := range []string{
			store.SettingGoogleClientID, store.SettingGoogleClientSecret,
			store.SettingGoogleAllowedDomains, store.SettingGoogleAutoProvisionRole,
		} {
			if err := s.store.DeleteSetting(ctx, key); err != nil {
				writeStoreError(w, err, "")
				return
			}
		}
		s.handleGetGoogleSettings(w, r)
		return
	}

	clientID := trimmed(req.ClientID)
	if clientID == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập Client ID")
		return
	}
	if req.AutoProvisionRole != "" && req.AutoProvisionRole != models.RoleStudent && req.AutoProvisionRole != models.RoleTrainer {
		writeError(w, http.StatusBadRequest, "Vai trò tự tạo tài khoản chỉ nhận giá trị học viên hoặc giảng viên")
		return
	}

	secret := trimmed(req.ClientSecret)
	if secret == "" {
		// Giữ nguyên secret đang có hiệu lực (dù đang lấy từ DB hay từ .env) nếu admin
		// chỉ sửa các trường khác mà để trống ô này — khớp với placeholder trên giao diện.
		current, err := s.resolveGoogleConfig(ctx)
		if err != nil {
			writeStoreError(w, err, "")
			return
		}
		if current.ClientSecret == "" {
			writeError(w, http.StatusBadRequest, "Vui lòng nhập Client Secret")
			return
		}
		secret = current.ClientSecret
	}

	// Chuẩn hoá danh sách domain: cắt khoảng trắng từng phần tử rồi nối lại,
	// để lần đọc sau không phải xử lý lại chuỗi thô người dùng gõ.
	domains := strings.Join(splitDomains(req.AllowedDomains), ", ")

	for key, value := range map[string]string{
		store.SettingGoogleClientID:          clientID,
		store.SettingGoogleClientSecret:      secret,
		store.SettingGoogleAllowedDomains:    domains,
		store.SettingGoogleAutoProvisionRole: req.AutoProvisionRole,
	} {
		if err := s.store.SetSetting(ctx, key, value, claims.UserID); err != nil {
			writeStoreError(w, err, "")
			return
		}
	}
	s.handleGetGoogleSettings(w, r)
}
