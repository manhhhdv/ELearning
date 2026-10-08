package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/manhnv/elearning/backend/internal/ai"
	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

func (s *Server) handleGenerateLesson(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProgramID       uuid.UUID `json:"programId"`
		Topic           string    `json:"topic"`
		Grade           string    `json:"grade"`
		Objectives      string    `json:"objectives"`
		Content         string    `json:"content"`
		DurationMinutes int       `json:"durationMinutes"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.ProgramID == uuid.Nil || trimmed(req.Topic) == "" || len([]rune(req.Topic)) > 300 || len([]rune(req.Content)) > maxDocxContentRunes || len([]rune(req.Objectives)) > 4000 || len([]rune(req.Grade)) > 100 || req.DurationMinutes < 1 || req.DurationMinutes > 240 {
		writeError(w, http.StatusBadRequest, "Nhập lớp học, chủ đề (tối đa 300 ký tự), tài liệu (tối đa 20.000 ký tự) và thời lượng từ 1–240 phút")
		return
	}
	s.generateLessonDraft(w, r, req.ProgramID, trimmed(req.Topic), trimmed(req.Grade), trimmed(req.Objectives), trimmed(req.Content), req.DurationMinutes)
}

// handleGenerateLessonFromFile giống handleGenerateLesson, nhưng nội dung
// tham khảo được bóc tách từ một tệp Word (.docx) tải lên thay vì dán tay —
// tiện cho giáo viên đã có sẵn tài liệu soạn trước đó.
func (s *Server) handleGenerateLessonFromFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxDocxBytes); err != nil {
		writeError(w, http.StatusBadRequest, "Không đọc được tệp tải lên: "+err.Error())
		return
	}
	programID, err := uuid.Parse(r.FormValue("programId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Vui lòng chọn lớp học")
		return
	}
	topic := trimmed(r.FormValue("topic"))
	grade := trimmed(r.FormValue("grade"))
	objectives := trimmed(r.FormValue("objectives"))
	duration, _ := strconv.Atoi(r.FormValue("durationMinutes"))
	if topic == "" || len([]rune(topic)) > 300 || len([]rune(objectives)) > 4000 || len([]rune(grade)) > 100 || duration < 1 || duration > 240 {
		writeError(w, http.StatusBadRequest, "Nhập lớp học, chủ đề (tối đa 300 ký tự) và thời lượng từ 1–240 phút")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Vui lòng chọn một tệp Word (.docx)")
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".docx") {
		writeError(w, http.StatusBadRequest, "Chỉ hỗ trợ tài liệu định dạng Word (.docx)")
		return
	}
	content, ok := readDocxContent(w, file, maxDocxContentRunes)
	if !ok {
		return
	}

	s.generateLessonDraft(w, r, programID, topic, grade, objectives, content, duration)
}

// generateLessonDraft là phần dùng chung của hai lối vào tạo bài giảng bằng
// AI (dán nội dung tham khảo, hoặc tải lên tệp Word): kiểm tra quyền, gọi
// Service AI rồi ghi nhật ký tác vụ.
func (s *Server) generateLessonDraft(
	w http.ResponseWriter, r *http.Request,
	programID uuid.UUID, topic, grade, objectives, content string, durationMinutes int,
) {
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}
	p, err := s.store.GetProgram(r.Context(), programID, uuid.Nil)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy lớp học")
		return
	}
	client, ok := s.requireAI(w, r)
	if !ok {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	draft, err := client.GenerateLesson(r.Context(), ai.LessonParams{
		Course: p.Title, Topic: topic, Grade: grade, Objectives: objectives,
		Content: content, DurationMinutes: durationMinutes,
	})
	task := store.LogAITaskParams{UserID: claims.UserID, Kind: models.AIKindLesson, Title: topic}
	if err != nil {
		task.Status = "error"
		task.Detail = err.Error()
		s.logAI(r.Context(), task)
		writeAIError(w, err)
		return
	}
	s.logAI(r.Context(), task)
	writeJSON(w, http.StatusOK, draft)
}
