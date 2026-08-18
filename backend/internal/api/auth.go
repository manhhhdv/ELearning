package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

type authConfigResponse struct {
	GoogleEnabled bool `json:"googleEnabled"`
}

func (s *Server) handleAuthConfig(w http.ResponseWriter, r *http.Request) {
	_, _, enabled, err := s.googleAuthenticator(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, authConfigResponse{GoogleEnabled: enabled})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Email = trimmed(req.Email)
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập email và mật khẩu")
		return
	}

	creds, err := s.store.GetCredentialsByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Không tiết lộ email nào tồn tại trong hệ thống.
			writeError(w, http.StatusUnauthorized, "Email hoặc mật khẩu không đúng")
			return
		}
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	if !auth.CheckPassword(creds.PasswordHash, req.Password) {
		writeError(w, http.StatusUnauthorized, "Email hoặc mật khẩu không đúng")
		return
	}
	if !creds.IsActive {
		writeError(w, http.StatusForbidden, "Tài khoản đã bị khoá, vui lòng liên hệ quản trị viên")
		return
	}

	s.issueToken(w, r, &creds.User)
}

func (s *Server) issueToken(w http.ResponseWriter, r *http.Request, u *models.User) {
	token, _, err := s.tokens.Issue(u.ID, u.Email, u.Role)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	s.store.TouchLastLogin(r.Context(), u.ID)
	writeJSON(w, http.StatusOK, loginResponse{Token: token, User: u})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	u, err := s.store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	if !u.IsActive {
		writeError(w, http.StatusForbidden, "Tài khoản đã bị khoá")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req changePasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())

	creds, err := s.store.GetCredentialsByID(r.Context(), claims.UserID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	// Tài khoản đăng nhập bằng Google chưa từng có mật khẩu thì được đặt mới mà không cần mật khẩu cũ.
	if creds.PasswordHash != "" && !auth.CheckPassword(creds.PasswordHash, req.CurrentPassword) {
		writeError(w, http.StatusBadRequest, "Mật khẩu hiện tại không đúng")
		return
	}
	if err := auth.ValidatePassword(req.NewPassword); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	if err := s.store.SetPassword(r.Context(), claims.UserID, hash, false); err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Đã đổi mật khẩu"})
}

// ---------------------------------------------------------------------------
// Đăng nhập bằng Google
// ---------------------------------------------------------------------------

func (s *Server) handleGoogleStart(w http.ResponseWriter, r *http.Request) {
	g, _, enabled, err := s.googleAuthenticator(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	if !enabled {
		writeError(w, http.StatusNotImplemented, "Chưa cấu hình đăng nhập Google trên máy chủ")
		return
	}
	state, err := g.NewState()
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	http.Redirect(w, r, g.AuthCodeURL(state), http.StatusFound)
}

// handleGoogleCallback nhận mã từ Google rồi chuyển hướng về frontend kèm token hoặc thông báo lỗi.
func (s *Server) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	g, gc, enabled, err := s.googleAuthenticator(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	if !enabled {
		writeError(w, http.StatusNotImplemented, "Chưa cấu hình đăng nhập Google trên máy chủ")
		return
	}
	q := r.URL.Query()
	if errMsg := q.Get("error"); errMsg != "" {
		s.redirectToFrontend(w, r, "", "Bạn đã huỷ đăng nhập bằng Google")
		return
	}
	if err := g.VerifyState(q.Get("state")); err != nil {
		s.redirectToFrontend(w, r, "", err.Error())
		return
	}

	profile, err := g.Exchange(r.Context(), q.Get("code"))
	if err != nil {
		s.redirectToFrontend(w, r, "", "Đăng nhập Google thất bại, vui lòng thử lại")
		return
	}
	if !domainAllowed(profile.Email, gc.AllowedDomains) {
		s.redirectToFrontend(w, r, "", "Email này không thuộc tổ chức được phép truy cập")
		return
	}

	user, err := s.resolveGoogleUser(r, profile, gc.AutoProvisionRole)
	if err != nil {
		s.redirectToFrontend(w, r, "", err.Error())
		return
	}
	if !user.IsActive {
		s.redirectToFrontend(w, r, "", "Tài khoản đã bị khoá, vui lòng liên hệ quản trị viên")
		return
	}

	token, _, err := s.tokens.Issue(user.ID, user.Email, user.Role)
	if err != nil {
		s.redirectToFrontend(w, r, "", "Không tạo được phiên đăng nhập")
		return
	}
	s.store.TouchLastLogin(r.Context(), user.ID)
	s.redirectToFrontend(w, r, token, "")
}

// resolveGoogleUser tìm tài khoản khớp với hồ sơ Google, gắn liên kết nếu tài khoản đã tồn tại
// theo email, hoặc tạo mới khi hệ thống cho phép tự đăng ký.
func (s *Server) resolveGoogleUser(r *http.Request, profile *auth.GoogleProfile, autoProvisionRole string) (*models.User, error) {
	ctx := r.Context()

	if u, err := s.store.GetUserByGoogleSub(ctx, profile.Sub); err == nil {
		return u, nil
	} else if !errors.Is(err, store.ErrNotFound) {
		return nil, errors.New("Không truy vấn được tài khoản")
	}

	if u, err := s.store.GetUserByEmail(ctx, profile.Email); err == nil {
		// Tài khoản do admin cấp sẵn: gắn Google vào để lần sau đăng nhập thẳng.
		if err := s.store.LinkGoogle(ctx, u.ID, profile.Sub, profile.Picture, profile.Name); err != nil {
			return nil, errors.New("Không liên kết được tài khoản Google")
		}
		u.HasGoogle = true
		return u, nil
	} else if !errors.Is(err, store.ErrNotFound) {
		return nil, errors.New("Không truy vấn được tài khoản")
	}

	if autoProvisionRole == "" {
		return nil, errors.New("Tài khoản chưa được cấp quyền truy cập, vui lòng liên hệ quản trị viên")
	}
	u, err := s.store.CreateUser(ctx, store.CreateUserParams{
		Email:     profile.Email,
		FullName:  profile.Name,
		Role:      autoProvisionRole,
		GoogleSub: profile.Sub,
		AvatarURL: profile.Picture,
	})
	if err != nil {
		return nil, errors.New("Không tạo được tài khoản mới")
	}
	return u, nil
}

// domainAllowed kiểm tra email có thuộc danh sách domain được phép hay không.
// Danh sách rỗng nghĩa là không giới hạn — mọi email đều qua được.
func domainAllowed(email string, allowedDomains []string) bool {
	if len(allowedDomains) == 0 {
		return true
	}
	_, domain, ok := strings.Cut(email, "@")
	if !ok {
		return false
	}
	for _, allowed := range allowedDomains {
		if strings.EqualFold(strings.TrimPrefix(allowed, "@"), domain) {
			return true
		}
	}
	return false
}

func (s *Server) redirectToFrontend(w http.ResponseWriter, r *http.Request, token, errMsg string) {
	target, err := url.Parse(s.cfg.FrontendURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "FRONTEND_URL cấu hình không hợp lệ")
		return
	}
	target.Path = strings.TrimSuffix(target.Path, "/") + "/dang-nhap"

	q := url.Values{}
	if token != "" {
		q.Set("token", token)
	}
	if errMsg != "" {
		q.Set("error", errMsg)
	}
	target.RawQuery = q.Encode()
	http.Redirect(w, r, target.String(), http.StatusFound)
}
