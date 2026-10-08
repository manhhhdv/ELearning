package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/ai"
	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/export"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// Xuất tệp Word (.docx) cho giáo viên: giáo án, bộ câu hỏi và bài giảng.

type exportLessonPlanRequest struct {
	Subject         string        `json:"subject"`
	Grade           string        `json:"grade"`
	DurationMinutes int           `json:"durationMinutes"`
	Content         ai.LessonPlan `json:"content"`
}

// handleExportLessonPlanDraft xuất đúng bản giáo án đang hiện trên màn hình
// soạn — kể cả khi chưa lưu hoặc vừa sửa tay — nên nhận nội dung trong body
// thay vì đọc lại từ cơ sở dữ liệu.
func (s *Server) handleExportLessonPlanDraft(w http.ResponseWriter, r *http.Request) {
	var req exportLessonPlanRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	info := export.LessonPlanInfo{
		Subject:         trimmed(req.Subject),
		Grade:           trimmed(req.Grade),
		DurationMinutes: req.DurationMinutes,
		Author:          s.authorName(r.Context()),
	}
	data, err := export.LessonPlan(info, &req.Content)
	writeDownload(w, export.ContentType, export.FileName("giao an", req.Content.Title), data, err)
}

// handleExportLessonPlan xuất một giáo án đã lưu.
func (s *Server) handleExportLessonPlan(w http.ResponseWriter, r *http.Request) {
	planID, ok := urlUUID(w, r, "planID")
	if !ok {
		return
	}
	plan, err := s.store.GetLessonPlan(r.Context(), planID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy giáo án")
		return
	}
	claims, _ := auth.FromContext(r.Context())
	if plan.OwnerID != claims.UserID {
		writeError(w, http.StatusForbidden, "Bạn không có quyền xem giáo án này")
		return
	}
	var content ai.LessonPlan
	if err := json.Unmarshal(plan.Content, &content); err != nil {
		slog.Error("đọc nội dung giáo án", "lỗi", err, "planID", plan.ID)
		writeError(w, http.StatusInternalServerError, "Nội dung giáo án bị lỗi, không xuất được")
		return
	}
	if content.Title == "" {
		content.Title = plan.Title
	}
	info := export.LessonPlanInfo{
		Subject:         plan.Subject,
		Grade:           plan.Grade,
		DurationMinutes: plan.DurationMinutes,
		Author:          s.authorName(r.Context()),
	}
	data, err := export.LessonPlan(info, &content)
	writeDownload(w, export.ContentType, export.FileName("giao an", content.Title), data, err)
}

// handleExportNode xuất một bài học (bài giảng) hoặc bộ câu hỏi của một bài
// tập. Với bài tập, ?answers=1 in thêm phần đáp án và hướng dẫn chấm.
// Chỉ người quản lý hoặc giám sát lớp học mới xuất được — tệp có thể chứa đáp án.
func (s *Server) handleExportNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, acc, ok := s.loadNode(w, r, nodeID, false)
	if !ok {
		return
	}
	if !acc.CanAudit {
		writeError(w, http.StatusForbidden, "Chỉ giáo viên phụ trách mới xuất được nội dung này")
		return
	}
	programTitle := ""
	if p, err := s.store.GetProgram(r.Context(), node.ProgramID, uuid.Nil); err == nil {
		programTitle = p.Title
	}
	author := s.authorName(r.Context())

	switch node.Kind {
	case models.KindLesson:
		data, err := export.Lesson(programTitle, author, node)
		writeDownload(w, export.ContentType, export.FileName("bai giang", node.Title), data, err)

	case models.KindAssignment:
		a := node.Assignment
		if a == nil || len(a.Questions) == 0 {
			writeError(w, http.StatusBadRequest, "Bài tập chưa có câu hỏi nào để xuất")
			return
		}
		withAnswers, _ := strconv.ParseBool(r.URL.Query().Get("answers"))
		info := export.QuestionSetInfo{
			ProgramTitle:     programTitle,
			Title:            node.Title,
			Instructions:     a.Instructions,
			TimeLimitMinutes: a.TimeLimitMinutes,
			Author:           author,
		}
		data, err := export.Questions(info, a.Questions, withAnswers)
		prefix := "de"
		if withAnswers {
			prefix = "de va dap an"
		}
		writeDownload(w, export.ContentType, export.FileName(prefix, node.Title), data, err)

	default:
		writeError(w, http.StatusBadRequest, "Chỉ xuất được bài học hoặc bài tập")
	}
}

// authorName là họ tên người đang xuất, ghi vào thuộc tính "Tác giả" của tệp.
func (s *Server) authorName(ctx context.Context) string {
	claims, ok := auth.FromContext(ctx)
	if !ok {
		return ""
	}
	u, err := s.store.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return ""
	}
	return u.FullName
}

// handleExportProgram xuất toàn bộ lớp học thành một tệp .zip gồm các tệp
// Word của từng bài học và bài tập. Tệp có đáp án nên chỉ người quản lý hoặc
// giám sát lớp học mới tải được.
func (s *Server) handleExportProgram(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAudit(w, r, programID); !ok {
		return
	}
	ctx := r.Context()
	program, err := s.store.GetProgram(ctx, programID, uuid.Nil)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy lớp học")
		return
	}
	nodes, err := s.store.ListNodes(ctx, programID, false)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	// Cây chỉ có số câu hỏi; nạp đủ câu hỏi (kèm đáp án) cho bài tập nào có câu.
	for _, n := range nodes {
		if n.Assignment == nil || n.Assignment.QuestionCount == 0 {
			continue
		}
		if n.Assignment.Questions, err = s.store.ListQuestions(ctx, n.ID); err != nil {
			writeStoreError(w, err, "")
			return
		}
	}

	data, err := export.ProgramArchive(program.Title, s.authorName(ctx), store.BuildTree(nodes))
	if errors.Is(err, export.ErrEmptyArchive) {
		writeError(w, http.StatusBadRequest, "Lớp học chưa có bài học hay bài tập có câu hỏi nào để xuất")
		return
	}
	writeDownload(w, export.ArchiveContentType, export.ArchiveName(program.Title), data, err)
}

// writeDownload gửi tệp về để trình duyệt tải xuống.
func writeDownload(w http.ResponseWriter, contentType, filename string, data []byte, err error) {
	if err != nil {
		slog.Error("dựng tệp xuất", "lỗi", err, "tệp", filename)
		writeError(w, http.StatusInternalServerError, "Không tạo được tệp, vui lòng thử lại")
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
