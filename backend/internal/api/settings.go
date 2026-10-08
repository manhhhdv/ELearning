package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/manhnv/elearning/backend/internal/ai"
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
		writeError(w, http.StatusBadRequest, "Vai trò tự tạo tài khoản chỉ nhận giá trị học sinh hoặc giáo viên")
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

// ---------------------------------------------------------------------------
// Cấu hình nhà cung cấp AI — chỉ admin xem/sửa được (khoá API là thông tin nhạy cảm).
//
// Cùng nguyên tắc với cấu hình Google: khoá nào có trong DB thì ưu tiên DB, còn
// lại rơi về giá trị .env. Nhờ vậy đổi khoá hay đổi model không cần khởi động
// lại máy chủ.
// ---------------------------------------------------------------------------

// aiConfig là cấu hình AI có hiệu lực sau khi gộp DB với .env.
type aiConfig struct {
	// Channels là các kênh gọi AI theo đúng thứ tự ưu tiên khi xoay vòng.
	Channels []ai.Credential
	// Nguồn cấu hình đang dùng, chỉ để hiển thị cho admin biết.
	Source string // "database" | "env" | "none"
}

// enabledChannels lọc ra các kênh thật sự dùng được.
func (c aiConfig) enabledChannels() []ai.Credential {
	out := make([]ai.Credential, 0, len(c.Channels))
	for _, ch := range c.Channels {
		if ch.APIKey != "" && !ch.Disabled {
			out = append(out, ch)
		}
	}
	return out
}

// resolveAIConfig đọc cấu hình theo thứ tự: danh sách kênh trong DB → cấu hình
// một kênh kiểu cũ trong DB → biến môi trường.
func (s *Server) resolveAIConfig(ctx context.Context) (aiConfig, error) {
	saved, err := s.store.GetSettings(ctx, []string{
		store.SettingAIChannels, store.SettingAIProvider, store.SettingAIAPIKey, store.SettingAIModel,
	})
	if err != nil {
		return aiConfig{}, err
	}

	cfg := aiConfig{Source: "database"}
	switch raw, ok := saved[store.SettingAIChannels]; {
	case ok && strings.TrimSpace(raw) != "":
		var channels []ai.Credential
		if err := json.Unmarshal([]byte(raw), &channels); err != nil {
			// Dữ liệu hỏng thì coi như chưa cấu hình còn hơn làm chết mọi request.
			slog.Error("danh sách kênh AI trong DB không đọc được", "lỗi", err)
			channels = nil
		}
		for _, ch := range channels {
			cfg.Channels = append(cfg.Channels, ch.Normalize())
		}
	case func() bool { _, has := saved[store.SettingAIAPIKey]; return has }():
		// Cấu hình một kênh kiểu cũ: giữ nguyên để bản cài cũ không mất khoá.
		cfg.Channels = []ai.Credential{ai.Credential{
			Provider: saved[store.SettingAIProvider],
			APIKey:   saved[store.SettingAIAPIKey],
			Model:    saved[store.SettingAIModel],
		}.Normalize()}
	default:
		cfg.Source = "env"
		if s.cfg.GeminiAPIKey != "" {
			cfg.Channels = []ai.Credential{ai.Credential{
				Provider: ai.ProviderGemini,
				APIKey:   s.cfg.GeminiAPIKey,
				Model:    s.cfg.GeminiModel,
			}.Normalize()}
		}
	}

	if len(cfg.enabledChannels()) == 0 {
		cfg.Source = "none"
	}
	return cfg, nil
}

// aiClient dựng lớp Service AI theo cấu hình đang có hiệu lực.
// Client rỗng kênh vẫn trả về (không nil) để hàm gọi chỉ cần kiểm tra Enabled().
func (s *Server) aiClient(ctx context.Context) *ai.Client {
	cfg, err := s.resolveAIConfig(ctx)
	if err != nil {
		// Không đọc được cấu hình thì coi như AI chưa bật, đồng thời ghi log để
		// quản trị viên biết — tốt hơn là làm hỏng cả request.
		slog.Error("không đọc được cấu hình AI", "lỗi", err)
		return ai.NewRotatingClient(nil)
	}
	return ai.NewRotatingClient(cfg.Channels)
}

// aiChannelView là một kênh gửi xuống giao diện — không bao giờ kèm khoá thật.
type aiChannelView struct {
	Label     string `json:"label"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	BaseURL   string `json:"baseUrl"`
	Disabled  bool   `json:"disabled"`
	HasApiKey bool   `json:"hasApiKey"`
	// Bốn ký tự cuối của khoá, đủ để admin nhận ra mình đang dùng khoá nào.
	ApiKeyHint string `json:"apiKeyHint"`
}

type aiSettingsResponse struct {
	Enabled bool   `json:"enabled"`
	Source  string `json:"source"`
	// Danh sách kênh theo thứ tự xoay vòng.
	Channels []aiChannelView `json:"channels"`
	// Các trường của kênh đầu tiên, giữ lại cho phần giao diện cũ.
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	HasApiKey  bool   `json:"hasApiKey"`
	ApiKeyHint string `json:"apiKeyHint"`
	// Danh sách nhà cung cấp và model gợi ý cho ô chọn trên giao diện.
	Providers []ai.Provider `json:"providers"`
}

func (s *Server) aiSettingsBody(ctx context.Context) (aiSettingsResponse, error) {
	cfg, err := s.resolveAIConfig(ctx)
	if err != nil {
		return aiSettingsResponse{}, err
	}
	body := aiSettingsResponse{
		Enabled:   len(cfg.enabledChannels()) > 0,
		Source:    cfg.Source,
		Channels:  make([]aiChannelView, 0, len(cfg.Channels)),
		Providers: ai.Providers,
	}
	for _, ch := range cfg.Channels {
		body.Channels = append(body.Channels, aiChannelView{
			Label:      ch.Label,
			Provider:   ch.Provider,
			Model:      ch.Model,
			BaseURL:    ch.BaseURL,
			Disabled:   ch.Disabled,
			HasApiKey:  ch.APIKey != "",
			ApiKeyHint: maskAPIKey(ch.APIKey),
		})
	}
	if active := cfg.enabledChannels(); len(active) > 0 {
		body.Provider = active[0].Provider
		body.Model = active[0].Model
		body.HasApiKey = true
		body.ApiKeyHint = maskAPIKey(active[0].APIKey)
	} else {
		body.Provider = ai.DefaultProvider
		body.Model = ai.DefaultModelFor(ai.DefaultProvider)
	}
	return body, nil
}

func (s *Server) handleGetAISettings(w http.ResponseWriter, r *http.Request) {
	body, err := s.aiSettingsBody(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, body)
}

// saveAIChannelRequest là một kênh do admin gửi lên.
type saveAIChannelRequest struct {
	Label    string `json:"label"`
	Provider string `json:"provider"`
	// Để trống nghĩa là giữ nguyên khoá đã lưu ở kênh cùng vị trí.
	ApiKey   string `json:"apiKey"`
	Model    string `json:"model"`
	BaseURL  string `json:"baseUrl"`
	Disabled bool   `json:"disabled"`
}

type saveAISettingsRequest struct {
	// false: xoá cấu hình đã lưu, quay lại dùng giá trị trong .env (nếu có).
	Enabled bool `json:"enabled"`
	// Danh sách kênh theo thứ tự ưu tiên. Bỏ trống thì dùng ba trường kiểu cũ.
	Channels []saveAIChannelRequest `json:"channels"`
	// Cấu hình một kênh kiểu cũ, vẫn nhận để không phá các bản frontend cũ.
	Provider string `json:"provider"`
	ApiKey   string `json:"apiKey"`
	Model    string `json:"model"`
}

// channels quy về một danh sách kênh, gộp cả dạng cũ lẫn dạng mới.
func (r saveAISettingsRequest) channels() []saveAIChannelRequest {
	if len(r.Channels) > 0 {
		return r.Channels
	}
	if r.Provider == "" && r.ApiKey == "" && r.Model == "" {
		return nil
	}
	return []saveAIChannelRequest{{Provider: r.Provider, ApiKey: r.ApiKey, Model: r.Model}}
}

// buildChannel kiểm tra một kênh và điền khoá cũ nếu admin để trống ô khoá.
// previous là kênh đang lưu ở cùng vị trí (nil nếu là kênh mới thêm).
func buildChannel(in saveAIChannelRequest, previous *ai.Credential) (ai.Credential, error) {
	provider, err := ai.NormalizeProvider(in.Provider)
	if err != nil {
		return ai.Credential{}, err
	}
	prov, _ := ai.FindProvider(provider)

	apiKey := trimmed(in.ApiKey)
	if apiKey == "" && previous != nil {
		apiKey = previous.APIKey
	}
	if apiKey == "" {
		return ai.Credential{}, fmt.Errorf("kênh %q chưa có khoá API", channelName(in, provider))
	}

	baseURL := trimmed(in.BaseURL)
	if baseURL == "" && previous != nil && previous.Provider == provider {
		baseURL = previous.BaseURL
	}
	if prov.RequiresBaseURL && baseURL == "" {
		return ai.Credential{}, fmt.Errorf("kênh %q cần địa chỉ endpoint", channelName(in, provider))
	}

	model := trimmed(in.Model)
	if model == "" {
		model = ai.DefaultModelFor(provider)
	}
	if model == "" {
		return ai.Credential{}, fmt.Errorf("kênh %q chưa chọn model", channelName(in, provider))
	}

	return ai.Credential{
		Label:    trimmed(in.Label),
		Provider: provider,
		APIKey:   apiKey,
		Model:    model,
		BaseURL:  baseURL,
		Disabled: in.Disabled,
	}.Normalize(), nil
}

func channelName(in saveAIChannelRequest, provider string) string {
	if name := trimmed(in.Label); name != "" {
		return name
	}
	return provider
}

func (s *Server) handleSaveAISettings(w http.ResponseWriter, r *http.Request) {
	var req saveAISettingsRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	ctx := r.Context()

	incoming := req.channels()
	if !req.Enabled || len(incoming) == 0 {
		for _, key := range []string{
			store.SettingAIChannels, store.SettingAIProvider,
			store.SettingAIAPIKey, store.SettingAIModel,
		} {
			if err := s.store.DeleteSetting(ctx, key); err != nil {
				writeStoreError(w, err, "")
				return
			}
		}
		s.handleGetAISettings(w, r)
		return
	}

	current, err := s.resolveAIConfig(ctx)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	channels := make([]ai.Credential, 0, len(incoming))
	for i, in := range incoming {
		var previous *ai.Credential
		if i < len(current.Channels) {
			previous = &current.Channels[i]
		}
		cred, err := buildChannel(in, previous)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		channels = append(channels, cred)
	}

	payload, err := json.Marshal(channels)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Không lưu được danh sách kênh AI")
		return
	}
	if err := s.store.SetSetting(ctx, store.SettingAIChannels, string(payload), claims.UserID); err != nil {
		writeStoreError(w, err, "")
		return
	}
	// Cấu hình một kênh kiểu cũ không còn dùng tới; xoá để tránh hai nguồn sự thật.
	for _, key := range []string{store.SettingAIProvider, store.SettingAIAPIKey, store.SettingAIModel} {
		if err := s.store.DeleteSetting(ctx, key); err != nil {
			writeStoreError(w, err, "")
			return
		}
	}
	s.handleGetAISettings(w, r)
}

// handleTestAISettings gọi thử một kênh để admin biết khoá và model có dùng
// được không, trước khi giáo viên gặp lỗi giữa lúc soạn đề.
//
// Body gửi lên là kênh cần thử; ô khoá để trống thì lấy khoá đang lưu ở kênh
// cùng vị trí, nhờ vậy admin kiểm tra được trước khi bấm lưu.
func (s *Server) handleTestAISettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		saveAIChannelRequest
		// Index là vị trí của kênh trong danh sách đang lưu, để lấy lại khoá cũ.
		Index *int `json:"index"`
		// Các trường kiểu cũ vẫn được chấp nhận.
		Channels []saveAIChannelRequest `json:"channels"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	ctx := r.Context()

	in := req.saveAIChannelRequest
	if len(req.Channels) > 0 {
		in = req.Channels[0]
	}

	current, err := s.resolveAIConfig(ctx)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	// Không chỉ rõ vị trí thì mặc định so với kênh đầu tiên.
	idx := 0
	if req.Index != nil {
		idx = *req.Index
	}
	var previous *ai.Credential
	if idx >= 0 && idx < len(current.Channels) {
		previous = &current.Channels[idx]
	}

	cred, err := buildChannel(in, previous)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Kênh đang tắt vẫn phải thử được — admin thường tắt vì nghi nó hỏng.
	cred.Disabled = false

	if err := ai.NewRotatingClient([]ai.Credential{cred}).Ping(ctx); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":      false,
			"message": err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"message": "Kết nối thành công tới model " + cred.Model,
	})
}

// maskAPIKey chỉ giữ lại 4 ký tự cuối để admin nhận ra mình đang dùng khoá nào
// mà không lộ khoá ra ngoài.
func maskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 4 {
		return "••••"
	}
	return "••••" + key[len(key)-4:]
}

// ---------------------------------------------------------------------------
// Cấu hình tự đăng ký tài khoản
// ---------------------------------------------------------------------------

// signupConfig là cấu hình tự đăng ký đang có hiệu lực.
type signupConfig struct {
	Enabled bool
	// Rỗng = mọi email đều đăng ký được.
	AllowedDomains []string
}

func (s *Server) resolveSignupConfig(ctx context.Context) (signupConfig, error) {
	saved, err := s.store.GetSettings(ctx, []string{
		store.SettingSignupEnabled, store.SettingSignupAllowedDomains,
	})
	if err != nil {
		return signupConfig{}, err
	}
	// Mặc định tắt: mở đăng ký là quyết định của admin, không phải mặc định của hệ thống.
	return signupConfig{
		Enabled:        saved[store.SettingSignupEnabled] == "true",
		AllowedDomains: splitDomains(saved[store.SettingSignupAllowedDomains]),
	}, nil
}

type signupSettingsResponse struct {
	Enabled        bool   `json:"enabled"`
	AllowedDomains string `json:"allowedDomains"`
}

func (s *Server) handleGetSignupSettings(w http.ResponseWriter, r *http.Request) {
	cfg, err := s.resolveSignupConfig(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, signupSettingsResponse{
		Enabled:        cfg.Enabled,
		AllowedDomains: strings.Join(cfg.AllowedDomains, ", "),
	})
}

type saveSignupSettingsRequest struct {
	Enabled        bool   `json:"enabled"`
	AllowedDomains string `json:"allowedDomains"`
}

func (s *Server) handleSaveSignupSettings(w http.ResponseWriter, r *http.Request) {
	var req saveSignupSettingsRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	ctx := r.Context()

	enabled := "false"
	if req.Enabled {
		enabled = "true"
	}
	for key, value := range map[string]string{
		store.SettingSignupEnabled:        enabled,
		store.SettingSignupAllowedDomains: strings.Join(splitDomains(req.AllowedDomains), ", "),
	} {
		if err := s.store.SetSetting(ctx, key, value, claims.UserID); err != nil {
			writeStoreError(w, err, "")
			return
		}
	}
	s.handleGetSignupSettings(w, r)
}
