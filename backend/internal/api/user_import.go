package api

import (
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// maxUserImportItems chặn một lần nhập tạo ra quá nhiều tài khoản.
const maxUserImportItems = 500

// Trạng thái của từng dòng sau khi nhập.
const (
	importUserCreated = "created"
	importUserSkipped = "skipped"
	importUserFailed  = "failed"
)

// userItemInput là một dòng trong bảng danh sách nhân sự.
type userItemInput struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	// Bỏ trống để hệ thống sinh mật khẩu ngẫu nhiên cho riêng dòng này.
	Password string `json:"password"`
}

type importUsersRequest struct {
	Items []userItemInput `json:"items"`
}

type importUserResult struct {
	Row      int    `json:"row"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	// Mật khẩu thô chỉ có ở dòng tạo mới thành công và chỉ trả về đúng một lần.
	Password string `json:"password"`
	Message  string `json:"message"`
}

type importUsersResponse struct {
	Created int                `json:"created"`
	Skipped int                `json:"skipped"`
	Failed  int                `json:"failed"`
	Results []importUserResult `json:"results"`
}

// handleImportUsersFile đọc file .xlsx/.csv thành văn bản dạng bảng (ô cách nhau
// bằng Tab) để giao diện dùng chung một bộ phân tích với luồng dán tay.
func (s *Server) handleImportUsersFile(w http.ResponseWriter, r *http.Request) {
	rows, errMsg := readUploadedTable(w, r)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	writeJSON(w, http.StatusOK, importQuestionsFileResponse{Text: rowsToTabText(rows)})
}

// handleImportUsers cấp hàng loạt tài khoản. Khác với nhập câu hỏi hay cấu trúc,
// mỗi dòng được xử lý độc lập: danh sách nhân sự thường có sẵn vài người đã có
// tài khoản, nên huỷ cả lô chỉ vì một dòng trùng sẽ rất bất tiện. Kết quả từng
// dòng được trả về để admin biết chính xác ai được tạo, ai bị bỏ qua.
func (s *Server) handleImportUsers(w http.ResponseWriter, r *http.Request) {
	var req importUsersRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Không có tài khoản nào để nhập")
		return
	}
	if len(req.Items) > maxUserImportItems {
		writeError(w, http.StatusBadRequest, "Mỗi lần chỉ nhập được tối đa 500 tài khoản")
		return
	}

	resp := importUsersResponse{Results: make([]importUserResult, 0, len(req.Items))}
	// Email đã gặp trong chính lô này: file danh sách hay bị lặp dòng.
	seen := make(map[string]int, len(req.Items))

	for i, raw := range req.Items {
		item, errMsg := cleanUserItem(raw)
		res := importUserResult{Row: i + 1, Email: item.Email, FullName: item.FullName, Role: item.Role}

		switch {
		case errMsg != "":
			res.Status, res.Message = importUserFailed, errMsg
		case seen[item.Email] > 0:
			res.Status = importUserSkipped
			res.Message = "Trùng với dòng " + strconv.Itoa(seen[item.Email]) + " trong danh sách"
		default:
			seen[item.Email] = i + 1
			password, status, message := s.createImportedUser(r, item)
			res.Password, res.Status, res.Message = password, status, message
		}

		switch res.Status {
		case importUserCreated:
			resp.Created++
		case importUserSkipped:
			resp.Skipped++
		default:
			resp.Failed++
		}
		resp.Results = append(resp.Results, res)
	}

	writeJSON(w, http.StatusOK, resp)
}

// createImportedUser tạo một tài khoản và quy đổi lỗi thành trạng thái của dòng.
func (s *Server) createImportedUser(r *http.Request, item userItemInput) (password, status, message string) {
	password = item.Password
	if password == "" {
		var err error
		if password, err = auth.GeneratePassword(12); err != nil {
			slog.Error("sinh mật khẩu khi nhập tài khoản", "lỗi", err)
			return "", importUserFailed, "Không sinh được mật khẩu"
		}
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		slog.Error("băm mật khẩu khi nhập tài khoản", "lỗi", err)
		return "", importUserFailed, "Không tạo được mật khẩu"
	}

	_, err = s.store.CreateUser(r.Context(), store.CreateUserParams{
		Email:              item.Email,
		FullName:           item.FullName,
		Role:               item.Role,
		PasswordHash:       hash,
		MustChangePassword: true,
	})
	switch {
	case err == nil:
		return password, importUserCreated, ""
	case errors.Is(err, store.ErrConflict):
		return "", importUserSkipped, "Email đã có tài khoản trong hệ thống"
	default:
		slog.Error("tạo tài khoản khi nhập hàng loạt", "email", item.Email, "lỗi", err)
		return "", importUserFailed, "Không tạo được tài khoản"
	}
}

// cleanUserItem chuẩn hoá và kiểm tra một dòng; thông báo trả về được viết cho
// admin đọc trực tiếp trong bảng kết quả.
func cleanUserItem(in userItemInput) (userItemInput, string) {
	out := userItemInput{
		Email:    strings.ToLower(trimmed(in.Email)),
		FullName: trimmed(in.FullName),
		Role:     strings.ToLower(trimmed(in.Role)),
		Password: trimmed(in.Password),
	}
	if out.Role == "" {
		out.Role = models.RoleStudent
	}

	if out.Email == "" {
		return out, "Thiếu email"
	}
	if !validEmail(out.Email) {
		return out, "Email không hợp lệ"
	}
	if !validRole(out.Role) {
		return out, "Vai trò không hợp lệ"
	}
	if out.Password != "" {
		if err := auth.ValidatePassword(out.Password); err != nil {
			return out, err.Error()
		}
	}
	return out, ""
}

// validEmail chỉ chấp nhận địa chỉ trần (không kèm tên hiển thị) và có tên miền
// đủ dạng — dữ liệu dán từ bảng tính rất hay lẫn ô ghi chú.
func validEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}
	at := strings.LastIndex(email, "@")
	domain := email[at+1:]
	return !strings.ContainsAny(email, " \t") && strings.Contains(domain, ".") &&
		!strings.HasPrefix(domain, ".") && !strings.HasSuffix(domain, ".")
}
