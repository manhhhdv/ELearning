package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
	"github.com/manhnv/elearning/backend/internal/util"
)

// lessonInput là phần dữ liệu bài học gửi lên từ form soạn thảo.
// Source nhận link chia sẻ Google Drive hoặc ID file, hệ thống tự dựng URL nhúng.
type lessonInput struct {
	ContentType     string            `json:"contentType"`
	Source          string            `json:"source"`
	DurationMinutes int               `json:"durationMinutes"`
	Body            string            `json:"body"`
	Attachments     []attachmentInput `json:"attachments"`
}

// attachmentInput là một tài liệu tải về của bài học dạng "materials".
type attachmentInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Chặn payload rác; một bài học thực tế không cần nhiều hơn số này.
const maxLessonAttachments = 50

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
	IsLocked    *bool            `json:"isLocked"`
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
		IsLocked:    req.IsLocked != nil && *req.IsLocked,
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
	IsLocked    *bool            `json:"isLocked"`
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
		IsLocked:    req.IsLocked,
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
	if !manage && !acc.CanAudit && node.IsLocked {
		writeError(w, http.StatusForbidden, "Nội dung này đang bị khoá")
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

	attachments, err := buildAttachments(in.Attachments)
	if err != nil {
		return nil, err
	}

	// Bài tự soạn và bài tài liệu không nhúng file bên ngoài.
	if in.ContentType == "richtext" || in.ContentType == "materials" {
		return &models.Lesson{
			ContentType:     in.ContentType,
			DurationMinutes: in.DurationMinutes,
			Body:            in.Body,
			Attachments:     attachments,
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
		Attachments:     attachments,
	}, nil
}

// buildAttachments chuẩn hoá danh sách tài liệu tải về: bỏ dòng trống,
// kiểm tra link và lấy tên hiển thị mặc định từ chính link khi để trống.
func buildAttachments(in []attachmentInput) ([]models.LessonAttachment, error) {
	list := make([]models.LessonAttachment, 0, len(in))
	for _, a := range in {
		name, raw := trimmed(a.Name), trimmed(a.URL)
		if name == "" && raw == "" {
			continue
		}
		link := normalizeAttachmentURL(raw)
		if link == "" {
			label := name
			if label == "" {
				label = "không tên"
			}
			return nil, errValidation("Tài liệu “" + label + "” chưa có link tải hợp lệ (http/https hoặc link Google Drive)")
		}
		if name == "" {
			name = link
		}
		list = append(list, models.LessonAttachment{Name: name, URL: link})
	}
	if len(list) > maxLessonAttachments {
		return nil, errValidation(fmt.Sprintf("Mỗi bài học chỉ đính kèm tối đa %d tài liệu", maxLessonAttachments))
	}
	return list, nil
}

// normalizeAttachmentURL nhận link http(s) bất kỳ; nếu người soạn chỉ dán ID file
// Google Drive thì dựng sẵn link Drive tương ứng. Trả về "" khi không dùng được.
func normalizeAttachmentURL(raw string) string {
	if raw == "" {
		return ""
	}
	if u, err := url.Parse(raw); err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" {
		return raw
	}
	if id := util.ExtractDriveID(raw); id != "" {
		return "https://drive.google.com/file/d/" + id + "/view"
	}
	return ""
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
	case "video", "slide", "document", "pdf", "link", "richtext", "materials":
		return true
	}
	return false
}
