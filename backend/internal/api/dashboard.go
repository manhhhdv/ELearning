package api

import "net/http"

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
