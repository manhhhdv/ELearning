package api

import (
	"net/http"
	"strconv"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

type listUsersResponse struct {
	Items []*models.User `json:"items"`
	Total int            `json:"total"`
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	role := q.Get("role")
	if role != "" && !validRole(role) {
		writeError(w, http.StatusBadRequest, "Vai trò không hợp lệ")
		return
	}

	users, total, err := s.store.ListUsers(r.Context(), store.ListUsersFilter{
		Search: q.Get("search"),
		Role:   role,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, listUsersResponse{Items: users, Total: total})
}

type createUserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	// Bỏ trống để hệ thống sinh mật khẩu ngẫu nhiên và trả về cho admin.
	Password string `json:"password"`
}

type createUserResponse struct {
	User *models.User `json:"user"`
	// Mật khẩu thô chỉ xuất hiện đúng một lần ngay sau khi tạo hoặc đặt lại.
	Password string `json:"password"`
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Email = trimmed(req.Email)
	req.FullName = trimmed(req.FullName)

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập email")
		return
	}
	if req.Role == "" {
		req.Role = models.RoleStudent
	}
	if !validRole(req.Role) {
		writeError(w, http.StatusBadRequest, "Vai trò không hợp lệ")
		return
	}

	password := trimmed(req.Password)
	generated := password == ""
	if generated {
		var err error
		if password, err = auth.GeneratePassword(12); err != nil {
			writeStoreError(w, err, "")
			return
		}
	} else if err := auth.ValidatePassword(password); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	user, err := s.store.CreateUser(r.Context(), store.CreateUserParams{
		Email:              req.Email,
		FullName:           req.FullName,
		Role:               req.Role,
		PasswordHash:       hash,
		MustChangePassword: true,
	})
	if err != nil {
		writeStoreError(w, err, "Không tạo được tài khoản")
		return
	}
	writeJSON(w, http.StatusCreated, createUserResponse{User: user, Password: password})
}

type updateUserRequest struct {
	FullName *string `json:"fullName"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"isActive"`
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := urlUUID(w, r, "userID")
	if !ok {
		return
	}
	var req updateUserRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Role != nil && !validRole(*req.Role) {
		writeError(w, http.StatusBadRequest, "Vai trò không hợp lệ")
		return
	}

	// Giữ lại ít nhất một admin đang hoạt động.
	demoting := (req.Role != nil && *req.Role != models.RoleAdmin) || (req.IsActive != nil && !*req.IsActive)
	if demoting {
		target, err := s.store.GetUserByID(r.Context(), userID)
		if err != nil {
			writeStoreError(w, err, "Không tìm thấy tài khoản")
			return
		}
		if target.Role == models.RoleAdmin && target.IsActive {
			others, err := s.store.CountActiveAdmins(r.Context(), userID)
			if err != nil {
				writeStoreError(w, err, "")
				return
			}
			if others == 0 {
				writeError(w, http.StatusBadRequest, "Không thể hạ quyền hoặc khoá quản trị viên cuối cùng")
				return
			}
		}
	}

	user, err := s.store.UpdateUser(r.Context(), userID, store.UpdateUserParams{
		FullName: req.FullName,
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

type resetPasswordRequest struct {
	Password string `json:"password"`
}

func (s *Server) handleResetPassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := urlUUID(w, r, "userID")
	if !ok {
		return
	}
	var req resetPasswordRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	password := trimmed(req.Password)
	if password == "" {
		var err error
		if password, err = auth.GeneratePassword(12); err != nil {
			writeStoreError(w, err, "")
			return
		}
	} else if err := auth.ValidatePassword(password); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	if err := s.store.SetPassword(r.Context(), userID, hash, true); err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}

	user, err := s.store.GetUserByID(r.Context(), userID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	writeJSON(w, http.StatusOK, createUserResponse{User: user, Password: password})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := urlUUID(w, r, "userID")
	if !ok {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	if claims.UserID == userID {
		writeError(w, http.StatusBadRequest, "Không thể tự xoá tài khoản đang đăng nhập")
		return
	}

	target, err := s.store.GetUserByID(r.Context(), userID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	if target.Role == models.RoleAdmin {
		others, err := s.store.CountActiveAdmins(r.Context(), userID)
		if err != nil {
			writeStoreError(w, err, "")
			return
		}
		if others == 0 {
			writeError(w, http.StatusBadRequest, "Không thể xoá quản trị viên cuối cùng")
			return
		}
	}

	if err := s.store.DeleteUser(r.Context(), userID); err != nil {
		writeStoreError(w, err, "Không tìm thấy tài khoản")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validRole(role string) bool {
	switch role {
	case models.RoleAdmin, models.RoleTrainer, models.RoleSupervisor, models.RoleStudent:
		return true
	}
	return false
}
