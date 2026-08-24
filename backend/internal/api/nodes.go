package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
	"github.com/manhnv/elearning/backend/internal/util"
)

// lessonInput là phần dữ liệu bài học gửi lên từ form soạn thảo.
// Source nhận link chia sẻ Google Drive hoặc ID file, hệ thống tự dựng URL nhúng.
type lessonInput struct {
	ContentType     string `json:"contentType"`
	Source          string `json:"source"`
	DurationMinutes int    `json:"durationMinutes"`
	Body            string `json:"body"`
}

type assignmentInput struct {
	Instructions     string     `json:"instructions"`
	TimeLimitMinutes int        `json:"timeLimitMinutes"`
	MaxAttempts      int        `json:"maxAttempts"`
	PassScore        float64    `json:"passScore"`
	ShuffleQuestions bool       `json:"shuffleQuestions"`
	DueAt            *time.Time `json:"dueAt"`
}

type createNodeRequest struct {
	ParentID    *uuid.UUID       `json:"parentId"`
	Kind        string           `json:"kind"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	IsPublished *bool            `json:"isPublished"`
	Lesson      *lessonInput     `json:"lesson"`
	Assignment  *assignmentInput `json:"assignment"`
}

func (s *Server) handleCreateNode(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}

	var req createNodeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	req.Title = trimmed(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "Vui lòng nhập tiêu đề")
		return
	}
	if !validNodeKind(req.Kind) {
		writeError(w, http.StatusBadRequest, "Loại nội dung không hợp lệ")
		return
	}

	published := true
	if req.IsPublished != nil {
		published = *req.IsPublished
	}

	params := store.SaveNodeParams{
		ProgramID:   programID,
		ParentID:    req.ParentID,
		Kind:        req.Kind,
		Title:       req.Title,
		Description: req.Description,
		IsPublished: published,
	}
	if req.Kind == models.KindLesson {
		lesson, err := buildLesson(req.Lesson)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		params.Lesson = lesson
	}
	if req.Kind == models.KindAssignment {
		params.Assignment = buildAssignment(req.Assignment)
	}

	node, err := s.store.CreateNode(r.Context(), params)
	if err != nil {
		if msg, ok := invalidMessage(err); ok {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeStoreError(w, err, "Không tìm thấy nút cha")
		return
	}
	writeJSON(w, http.StatusCreated, node)
}

func (s *Server) handleGetNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, acc, ok := s.loadNode(w, r, nodeID, false)
	if !ok {
		return
	}
	// Học viên không được thấy đáp án đúng khi xem nội dung bài tập.
	// Giám sát xem như người quản lý: thấy cả nội dung chưa xuất bản lẫn đáp án đúng.
	if !acc.CanManage && !acc.CanAudit {
		if !node.IsPublished {
			writeError(w, http.StatusForbidden, "Nội dung này chưa được xuất bản")
			return
		}
		stripAnswers(node)
	}
	writeJSON(w, http.StatusOK, node)
}

type updateNodeRequest struct {
	Title       *string          `json:"title"`
	Description *string          `json:"description"`
	IsPublished *bool            `json:"isPublished"`
	Lesson      *lessonInput     `json:"lesson"`
	Assignment  *assignmentInput `json:"assignment"`
}

func (s *Server) handleUpdateNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	if _, _, ok := s.loadNode(w, r, nodeID, true); !ok {
		return
	}

	var req updateNodeRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	params := store.UpdateNodeParams{
		Title:       req.Title,
		Description: req.Description,
		IsPublished: req.IsPublished,
	}
	if req.Lesson != nil {
		lesson, err := buildLesson(req.Lesson)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		params.Lesson = lesson
	}
	if req.Assignment != nil {
		params.Assignment = buildAssignment(req.Assignment)
	}

	node, err := s.store.UpdateNode(r.Context(), nodeID, params)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return
	}
	writeJSON(w, http.StatusOK, node)
}

func (s *Server) handleDeleteNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	if _, _, ok := s.loadNode(w, r, nodeID, true); !ok {
		return
	}
	if err := s.store.DeleteNode(r.Context(), nodeID); err != nil {
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type moveNodeRequest struct {
	ParentID *uuid.UUID `json:"parentId"`
	Position int        `json:"position"`
}

func (s *Server) handleMoveNode(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	if _, _, ok := s.loadNode(w, r, nodeID, true); !ok {
		return
	}

	var req moveNodeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := s.store.MoveNode(r.Context(), nodeID, req.ParentID, req.Position); err != nil {
		if msg, ok := invalidMessage(err); ok {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Đã cập nhật vị trí"})
}

// loadNode đọc một nút và kiểm tra quyền truy cập chương trình chứa nó.
func (s *Server) loadNode(w http.ResponseWriter, r *http.Request, nodeID uuid.UUID, manage bool) (*models.Node, access, bool) {
	node, err := s.store.GetNode(r.Context(), nodeID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return nil, access{}, false
	}
	acc, ok := s.requireProgramAccess(w, r, node.ProgramID, manage)
	if !ok {
		return nil, access{}, false
	}
	return node, acc, true
}

func buildLesson(in *lessonInput) (*models.Lesson, error) {
	if in == nil {
		return &models.Lesson{ContentType: "video"}, nil
	}
	if in.ContentType == "" {
		in.ContentType = "video"
	}
	if !validContentType(in.ContentType) {
		return nil, errValidation("Loại nội dung bài học không hợp lệ")
	}

	// Bài tự soạn không nhúng file bên ngoài: toàn bộ nội dung nằm trong Body.
	if in.ContentType == "richtext" {
		return &models.Lesson{
			ContentType:     in.ContentType,
			DurationMinutes: in.DurationMinutes,
			Body:            in.Body,
		}, nil
	}

	driveID, embedURL := util.BuildEmbedURL(in.ContentType, in.Source)
	if trimmed(in.Source) != "" && embedURL == "" {
		return nil, errValidation("Không nhận diện được link Google Drive. Hãy dán link chia sẻ hoặc ID file.")
	}
	return &models.Lesson{
		ContentType:     in.ContentType,
		DriveFileID:     driveID,
		EmbedURL:        embedURL,
		DurationMinutes: in.DurationMinutes,
		Body:            in.Body,
	}, nil
}

func buildAssignment(in *assignmentInput) *models.Assignment {
	if in == nil {
		return &models.Assignment{}
	}
	return &models.Assignment{
		Instructions:     in.Instructions,
		TimeLimitMinutes: in.TimeLimitMinutes,
		MaxAttempts:      in.MaxAttempts,
		PassScore:        in.PassScore,
		ShuffleQuestions: in.ShuffleQuestions,
		DueAt:            in.DueAt,
	}
}

// stripAnswers loại bỏ cờ đáp án đúng và lời giải trước khi trả nội dung cho học viên.
func stripAnswers(n *models.Node) {
	if n.Assignment == nil {
		return
	}
	for _, q := range n.Assignment.Questions {
		q.Explanation = ""
		for _, o := range q.Options {
			o.IsCorrect = false
		}
	}
}

func validNodeKind(kind string) bool {
	switch kind {
	case models.KindFolder, models.KindLesson, models.KindAssignment:
		return true
	}
	return false
}

func validContentType(ct string) bool {
	switch ct {
	case "video", "slide", "document", "pdf", "link", "richtext":
		return true
	}
	return false
}
