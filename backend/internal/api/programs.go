package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// access mô tả quyền của người dùng hiện tại đối với một chương trình.
type access struct {
	CanManage bool // sửa cây nội dung, ghi danh, chấm bài
	CanView   bool // xem nội dung đã xuất bản (học viên)
	// Xem được nội dung chưa xuất bản, đáp án đúng, danh sách ghi danh và bài nộp —
	// giống người quản lý ở khả năng xem, nhưng không sửa/xoá/chấm bài được.
	// Đúng cho cả người quản lý thật (CanManage kéo theo CanAudit) lẫn vai trò Giám sát.
	CanAudit bool
}

// isDefaultCourseVisible cho biết một khoá học mặc định có nên hiện với người chưa
// ghi danh hay không — chỉ khi đã xuất bản, tránh lộ nội dung nháp cho người ngoài.
func isDefaultCourseVisible(p *models.Program) bool {
	return p.IsDefaultCourse && p.Status == "published"
}

// programAccess quyết định quyền dựa trên vai trò hệ thống và vai trò ghi danh.
func (s *Server) programAccess(r *http.Request, programID uuid.UUID) (access, error) {
	claims, _ := auth.FromContext(r.Context())

	switch claims.Role {
	case models.RoleAdmin:
		return access{CanManage: true, CanView: true, CanAudit: true}, nil
	case models.RoleSupervisor:
		// Giám sát xem được mọi chương trình, kể cả bản nháp, ghi danh và bài nộp,
		// nhưng không sửa/xoá/chấm bài — CanManage luôn false.
		return access{CanView: true, CanAudit: true}, nil
	case models.RoleTrainer:
		program, err := s.store.GetProgram(r.Context(), programID, uuid.Nil)
		if err != nil {
			return access{}, err
		}
		if program.CreatedBy != nil && *program.CreatedBy == claims.UserID {
			return access{CanManage: true, CanView: true, CanAudit: true}, nil
		}
		role, err := s.store.EnrollmentRole(r.Context(), programID, claims.UserID)
		if err != nil {
			return access{}, err
		}
		manage := role == "trainer"
		canView := role != "" || isDefaultCourseVisible(program)
		return access{CanManage: manage, CanView: canView, CanAudit: manage}, nil
	default:
		role, err := s.store.EnrollmentRole(r.Context(), programID, claims.UserID)
		if err != nil {
			return access{}, err
		}
		if role != "" {
			return access{CanView: true}, nil
		}
		// Chưa ghi danh nhưng có thể vẫn xem được nếu đây là khoá học mặc định —
		// tự động hiện với mọi người, không cần dòng enrollment.
		program, err := s.store.GetProgram(r.Context(), programID, uuid.Nil)
		if err != nil {
			return access{}, err
		}
		return access{CanView: isDefaultCourseVisible(program)}, nil
	}
}

// requireProgramAccess trả về false và ghi lỗi nếu người dùng không đủ quyền.
func (s *Server) requireProgramAccess(w http.ResponseWriter, r *http.Request, programID uuid.UUID, manage bool) (access, bool) {
	acc, err := s.programAccess(r, programID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return access{}, false
	}
	if manage && !acc.CanManage {
		writeError(w, http.StatusForbidden, "Bạn không có quyền chỉnh sửa chương trình này")
		return access{}, false
	}
	if !manage && !acc.CanView {
		writeError(w, http.StatusForbidden, "Bạn chưa được ghi danh vào chương trình này")
		return access{}, false
	}
	return acc, true
}

// requireProgramAudit trả về false nếu người dùng không có quyền xem thông tin quản lý
// (ghi danh, bài nộp, kết quả) — dành cho người quản lý thật lẫn vai trò Giám sát.
func (s *Server) requireProgramAudit(w http.ResponseWriter, r *http.Request, programID uuid.UUID) (access, bool) {
	acc, err := s.programAccess(r, programID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return access{}, false
	}
	if !acc.CanAudit {
		writeError(w, http.StatusForbidden, "Bạn không có quyền xem thông tin này")
		return access{}, false
	}
	return acc, true
}

func (s *Server) handleListPrograms(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	q := r.URL.Query()

	filter := store.ListProgramsFilter{Search: q.Get("search"), Status: q.Get("status")}
	switch claims.Role {
	case models.RoleStudent:
		// Học viên chỉ thấy chương trình đã xuất bản mà mình được ghi danh.
		filter.EnrolledUserID = claims.UserID
		filter.ViewerID = claims.UserID
		filter.Status = "published"
	case models.RoleTrainer:
		// Giảng viên chỉ thấy chương trình mình tạo hoặc được ghi danh vào.
		filter.VisibleToUserID = claims.UserID
	}

	programs, err := s.store.ListPrograms(r.Context(), filter)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, programs)
}

type saveProgramRequest struct {
	Code            string `json:"code"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	CoverURL        string `json:"coverUrl"`
	Status          string `json:"status"`
	AllowSelfEnroll bool   `json:"allowSelfEnroll"`
	// Ảnh hưởng tới mọi người dùng trong hệ thống nên chỉ admin đặt được — xem checkAdminOnlyDefaultFlag.
	IsDefaultCourse bool `json:"isDefaultCourse"`
}

func (s *Server) handleCreateProgram(w http.ResponseWriter, r *http.Request) {
	var req saveProgramRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Code, req.Title = trimmed(req.Code), trimmed(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập tên chương trình")
		return
	}
	if req.Code == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập mã chương trình")
		return
	}
	if req.Status == "" {
		req.Status = "draft"
	}
	if !validProgramStatus(req.Status) {
		writeError(w, http.StatusBadRequest, "Trạng thái chương trình không hợp lệ")
		return
	}

	claims, _ := auth.FromContext(r.Context())
	if req.IsDefaultCourse && claims.Role != models.RoleAdmin {
		writeError(w, http.StatusForbidden, "Chỉ quản trị viên mới đặt được khoá học mặc định")
		return
	}
	program, err := s.store.CreateProgram(r.Context(), store.CreateProgramParams{
		Code:            req.Code,
		Title:           req.Title,
		Description:     req.Description,
		CoverURL:        req.CoverURL,
		Status:          req.Status,
		AllowSelfEnroll: req.AllowSelfEnroll,
		IsDefaultCourse: req.IsDefaultCourse,
		CreatedBy:       claims.UserID,
	})
	if err != nil {
		if errors.Is(err, store.ErrConflict) {
			writeError(w, http.StatusConflict, "Mã chương trình đã tồn tại")
			return
		}
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusCreated, program)
}

func (s *Server) handleGetProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, false); !ok {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	program, err := s.store.GetProgram(r.Context(), programID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	writeJSON(w, http.StatusOK, program)
}

// handleGetProgramBySlug tra chương trình theo slug trên URL (VD: /hoc/attp-2026)
// thay vì UUID nội bộ. Frontend dùng ID trả về để gọi tiếp các API còn lại.
func (s *Server) handleGetProgramBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	claims, _ := auth.FromContext(r.Context())
	program, err := s.store.GetProgramBySlug(r.Context(), slug, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	if _, ok := s.requireProgramAccess(w, r, program.ID, false); !ok {
		return
	}
	writeJSON(w, http.StatusOK, program)
}

type updateProgramRequest struct {
	Code            *string `json:"code"`
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	CoverURL        *string `json:"coverUrl"`
	Status          *string `json:"status"`
	AllowSelfEnroll *bool   `json:"allowSelfEnroll"`
	// Ảnh hưởng tới mọi người dùng trong hệ thống nên chỉ admin đổi được.
	IsDefaultCourse *bool `json:"isDefaultCourse"`
}

func (s *Server) handleUpdateProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}
	var req updateProgramRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Status != nil && !validProgramStatus(*req.Status) {
		writeError(w, http.StatusBadRequest, "Trạng thái chương trình không hợp lệ")
		return
	}
	claims, _ := auth.FromContext(r.Context())
	if req.IsDefaultCourse != nil && claims.Role != models.RoleAdmin {
		writeError(w, http.StatusForbidden, "Chỉ quản trị viên mới đổi được khoá học mặc định")
		return
	}

	program, err := s.store.UpdateProgram(r.Context(), programID, store.UpdateProgramParams{
		Code:            req.Code,
		Title:           req.Title,
		Description:     req.Description,
		CoverURL:        req.CoverURL,
		Status:          req.Status,
		AllowSelfEnroll: req.AllowSelfEnroll,
		IsDefaultCourse: req.IsDefaultCourse,
	})
	if err != nil {
		if errors.Is(err, store.ErrConflict) {
			writeError(w, http.StatusConflict, "Mã chương trình đã tồn tại")
			return
		}
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	writeJSON(w, http.StatusOK, program)
}

func (s *Server) handleDeleteProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}
	if err := s.store.DeleteProgram(r.Context(), programID); err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handleGetTree trả về cây nội dung kèm trạng thái hoàn thành của người xem.
// Người quản lý thấy cả nút chưa xuất bản, học viên chỉ thấy nút đã xuất bản.
func (s *Server) handleGetTree(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	acc, ok := s.requireProgramAccess(w, r, programID, false)
	if !ok {
		return
	}
	// Giám sát xem cây như người quản lý, thấy cả nội dung chưa xuất bản.
	seesAll := acc.CanManage || acc.CanAudit

	nodes, err := s.store.ListNodes(r.Context(), programID, !seesAll)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	// Tiến độ hoàn thành cá nhân áp dụng cho mọi vai trò xem được cây, kể cả
	// người quản lý — họ vẫn có thể tự đánh dấu hoàn thành khi xem trước bài học.
	claims, _ := auth.FromContext(r.Context())
	done, err := s.store.CompletedNodeIDs(r.Context(), claims.UserID, programID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	for _, n := range nodes {
		n.Completed = done[n.ID]
	}

	writeJSON(w, http.StatusOK, store.BuildTree(nodes))
}

func validProgramStatus(status string) bool {
	switch status {
	case "draft", "published", "archived":
		return true
	}
	return false
}
