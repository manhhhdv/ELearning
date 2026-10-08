package api

import (
	"net/http"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
)

// handleDashboard trả về số liệu tổng quan toàn hệ thống — dành cho admin và
// vai trò Giám sát (chỉ xem, không có nút hành động nào đi kèm số liệu này).
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.DashboardStats(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// handleTrainerAnalytics trả về số liệu học tập trong phạm vi các lớp học mà
// giáo viên đang đăng nhập phụ trách — bản thu hẹp của handleDashboard, không
// lộ số liệu của lớp học khác.
func (s *Server) handleTrainerAnalytics(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	stats, err := s.store.TrainerAnalytics(r.Context(), claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// ---------------------------------------------------------------------------
// Trang tổng quan — "không gian làm việc AI" (mục 4.2)
// ---------------------------------------------------------------------------

// handleWorkspace trả về dữ liệu cho trang tổng quan của người đang đăng nhập.
// Nội dung khác nhau theo vai trò: học viên thấy lịch học và tiến độ, người dạy
// thấy thêm số bài chờ chấm và cảnh báo học tập.
func (s *Server) handleWorkspace(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	ctx := r.Context()

	// Ai cũng có thể là người học, kể cả giáo viên được ghi danh vào lớp khác.
	ws, err := s.store.StudentWorkspace(ctx, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	if claims.Role == models.RoleTrainer || claims.Role == models.RoleAdmin {
		// Admin nhìn toàn hệ thống; giáo viên chỉ nhìn lớp mình phụ trách.
		trainer, err := s.store.TrainerWorkspace(ctx, claims.UserID, claims.Role == models.RoleAdmin)
		if err != nil {
			writeStoreError(w, err, "")
			return
		}
		ws.PendingGrading = trainer.PendingGrading
		ws.Alerts = trainer.Alerts
	}

	writeJSON(w, http.StatusOK, ws)
}
