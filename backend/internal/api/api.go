// Package api định nghĩa router và các handler HTTP của hệ thống.
package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/config"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

type Server struct {
	cfg    *config.Config
	store  *store.Store
	tokens *auth.TokenManager
	mw     *auth.Middleware
}

func NewServer(cfg *config.Config, st *store.Store) *Server {
	tokens := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTTTL)
	s := &Server{
		cfg:    cfg,
		store:  st,
		tokens: tokens,
		mw:     auth.NewMiddleware(tokens, writeError),
	}
	return s
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(requestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.cfg.AllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", s.handleHealth)

		// --- Công khai: đăng nhập ---
		r.Route("/auth", func(r chi.Router) {
			r.Get("/config", s.handleAuthConfig)
			r.Post("/login", s.handleLogin)
			r.Get("/google/start", s.handleGoogleStart)
			r.Get("/google/callback", s.handleGoogleCallback)

			r.Group(func(r chi.Router) {
				r.Use(s.mw.RequireAuth)
				r.Get("/me", s.handleMe)
				r.Post("/change-password", s.handleChangePassword)
			})
		})

		// --- Cần đăng nhập ---
		r.Group(func(r chi.Router) {
			r.Use(s.mw.RequireAuth)

			// Quản trị người dùng: chỉ admin.
			r.Route("/users", func(r chi.Router) {
				r.Use(s.mw.RequireRole(models.RoleAdmin))
				r.Get("/", s.handleListUsers)
				r.Post("/", s.handleCreateUser)
				r.Patch("/{userID}", s.handleUpdateUser)
				r.Post("/{userID}/password", s.handleResetPassword)
				r.Delete("/{userID}", s.handleDeleteUser)
			})

			// Chương trình đào tạo.
			r.Route("/programs", func(r chi.Router) {
				r.Get("/", s.handleListPrograms)
				r.With(s.mw.RequireRole(models.RoleAdmin, models.RoleTrainer)).Post("/", s.handleCreateProgram)
				// Đặt trước {programID}: chi ưu tiên khớp đoạn tĩnh "slug" hơn tham số động.
				r.Get("/slug/{slug}", s.handleGetProgramBySlug)

				r.Route("/{programID}", func(r chi.Router) {
					r.Get("/", s.handleGetProgram)
					r.Get("/tree", s.handleGetTree)
					r.Patch("/", s.handleUpdateProgram)
					r.Delete("/", s.handleDeleteProgram)
					r.Post("/nodes", s.handleCreateNode)
					r.Get("/enrollments", s.handleListEnrollments)
					r.Post("/enrollments", s.handleEnroll)
					r.Delete("/enrollments/{userID}", s.handleUnenroll)
					r.Get("/submissions", s.handleListProgramSubmissions)
					r.Post("/self-enroll", s.handleSelfEnroll)
					r.Delete("/self-enroll", s.handleSelfUnenroll)
				})
			})

			// Nút trên cây: bài học, bài tập, thư mục.
			r.Route("/nodes/{nodeID}", func(r chi.Router) {
				r.Get("/", s.handleGetNode)
				r.Patch("/", s.handleUpdateNode)
				r.Delete("/", s.handleDeleteNode)
				r.Post("/move", s.handleMoveNode)
				r.Post("/questions", s.handleCreateQuestion)
				r.Post("/questions/reorder", s.handleReorderQuestions)
				r.Post("/questions/import", s.handleImportQuestions)
				r.Post("/questions/import-file", s.handleImportQuestionsFile)
				r.Get("/results", s.handleAssignmentResults)
				r.Post("/complete", s.handleMarkComplete)
				r.Get("/attempt", s.handleGetAttempt)
				r.Post("/attempt/start", s.handleStartAttempt)
				r.Post("/submit", s.handleSubmitAssignment)
			})

			r.Patch("/questions/{questionID}", s.handleUpdateQuestion)
			r.Delete("/questions/{questionID}", s.handleDeleteQuestion)

			// Bài nộp.
			r.Get("/submissions/{submissionID}", s.handleGetSubmission)
			r.Post("/submissions/{submissionID}/grade", s.handleGradeSubmission)

			// Khu vực của học viên.
			r.Get("/catalog", s.handleCatalog)
			r.Get("/my/programs", s.handleMyPrograms)
			r.Get("/my/submissions", s.handleMySubmissions)

			// Khu vực quản trị hệ thống.
			r.Route("/admin", func(r chi.Router) {
				// Dashboard chỉ xem, nên vai trò Giám sát cũng vào được như admin.
				r.With(s.mw.RequireRole(models.RoleAdmin, models.RoleSupervisor)).
					Get("/dashboard", s.handleDashboard)

				// Cấu hình đăng nhập Google chứa thông tin nhạy cảm — chỉ admin.
				r.Route("/settings/google", func(r chi.Router) {
					r.Use(s.mw.RequireRole(models.RoleAdmin))
					r.Get("/", s.handleGetGoogleSettings)
					r.Put("/", s.handleSaveGoogleSettings)
				})
			})
		})
	})

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ---------------------------------------------------------------------------
// Tiện ích chung cho handler
// ---------------------------------------------------------------------------

type errorBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("ghi phản hồi JSON thất bại", "lỗi", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Error: message})
}

// writeStoreError quy đổi lỗi từ tầng store sang mã HTTP tương ứng.
func writeStoreError(w http.ResponseWriter, err error, notFoundMsg string) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, notFoundMsg)
	case errors.Is(err, store.ErrConflict):
		writeError(w, http.StatusConflict, "Dữ liệu đã tồn tại")
	default:
		slog.Error("lỗi máy chủ", "lỗi", err)
		writeError(w, http.StatusInternalServerError, "Đã xảy ra lỗi, vui lòng thử lại")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ: "+err.Error())
		return false
	}
	return true
}

// urlUUID đọc một tham số UUID trên đường dẫn.
func urlUUID(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Định danh trên đường dẫn không hợp lệ")
		return uuid.Nil, false
	}
	return id, true
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		if ww.Status() >= 400 {
			slog.Warn("request lỗi", "method", r.Method, "path", r.URL.Path, "status", ww.Status())
		}
	})
}

func trimmed(s string) string { return strings.TrimSpace(s) }
