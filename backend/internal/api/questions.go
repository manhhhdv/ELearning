package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

type questionOptionInput struct {
	Content   string `json:"content"`
	IsCorrect bool   `json:"isCorrect"`
}

type saveQuestionRequest struct {
	// Để trống khi tạo mới thì hệ thống tự sinh; khi sửa thì giữ nguyên mã cũ.
	Code        string                `json:"code"`
	Type        string                `json:"type"`
	Prompt      string                `json:"prompt"`
	Points      float64               `json:"points"`
	Explanation string                `json:"explanation"`
	Options     []questionOptionInput `json:"options"`
}

// toParams kiểm tra tính hợp lệ rồi chuyển sang tham số của tầng store.
func (req *saveQuestionRequest) toParams(assignmentID uuid.UUID) (store.SaveQuestionParams, error) {
	prompt := trimmed(req.Prompt)
	if prompt == "" {
		return store.SaveQuestionParams{}, errValidation("Vui lòng nhập nội dung câu hỏi")
	}
	if !validQuestionType(req.Type) {
		return store.SaveQuestionParams{}, errValidation("Loại câu hỏi không hợp lệ")
	}
	points := req.Points
	if points <= 0 {
		points = 1
	}

	params := store.SaveQuestionParams{
		AssignmentID: assignmentID,
		Code:         trimmed(req.Code),
		Type:         req.Type,
		Prompt:       prompt,
		Points:       points,
		Explanation:  req.Explanation,
	}

	if req.Type == models.QuestionEssay {
		return params, nil
	}

	correct := 0
	for _, o := range req.Options {
		content := trimmed(o.Content)
		if content == "" {
			continue
		}
		if o.IsCorrect {
			correct++
		}
		params.Options = append(params.Options, store.QuestionOptionInput{Content: content, IsCorrect: o.IsCorrect})
	}

	if len(params.Options) < 2 {
		return store.SaveQuestionParams{}, errValidation("Câu trắc nghiệm cần ít nhất 2 phương án")
	}
	if correct == 0 {
		return store.SaveQuestionParams{}, errValidation("Vui lòng đánh dấu ít nhất một phương án đúng")
	}
	if req.Type == models.QuestionSingleChoice && correct > 1 {
		return store.SaveQuestionParams{}, errValidation("Câu một đáp án chỉ được đánh dấu đúng một phương án")
	}
	return params, nil
}

func (s *Server) handleCreateQuestion(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, _, ok := s.loadNode(w, r, nodeID, true)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment {
		writeError(w, http.StatusBadRequest, "Chỉ có thể thêm câu hỏi vào một bài tập")
		return
	}

	var req saveQuestionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	params, err := req.toParams(nodeID)
	if err != nil {
		msg, _ := invalidMessage(err)
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	question, err := s.store.CreateQuestion(r.Context(), params)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}
	writeJSON(w, http.StatusCreated, question)
}

func (s *Server) handleUpdateQuestion(w http.ResponseWriter, r *http.Request) {
	questionID, ok := urlUUID(w, r, "questionID")
	if !ok {
		return
	}
	assignmentID, ok := s.requireQuestionAccess(w, r, questionID)
	if !ok {
		return
	}

	var req saveQuestionRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	params, err := req.toParams(assignmentID)
	if err != nil {
		msg, _ := invalidMessage(err)
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	question, err := s.store.UpdateQuestion(r.Context(), questionID, params)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy câu hỏi")
		return
	}
	writeJSON(w, http.StatusOK, question)
}

func (s *Server) handleDeleteQuestion(w http.ResponseWriter, r *http.Request) {
	questionID, ok := urlUUID(w, r, "questionID")
	if !ok {
		return
	}
	if _, ok := s.requireQuestionAccess(w, r, questionID); !ok {
		return
	}
	if err := s.store.DeleteQuestion(r.Context(), questionID); err != nil {
		writeStoreError(w, err, "Không tìm thấy câu hỏi")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type reorderQuestionsRequest struct {
	QuestionIDs []uuid.UUID `json:"questionIds"`
}

func (s *Server) handleReorderQuestions(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	if _, _, ok := s.loadNode(w, r, nodeID, true); !ok {
		return
	}

	var req reorderQuestionsRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if err := s.store.ReorderQuestions(r.Context(), nodeID, req.QuestionIDs); err != nil {
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Đã cập nhật thứ tự câu hỏi"})
}

// requireQuestionAccess xác định bài tập chứa câu hỏi rồi kiểm tra quyền chỉnh sửa.
func (s *Server) requireQuestionAccess(w http.ResponseWriter, r *http.Request, questionID uuid.UUID) (uuid.UUID, bool) {
	assignmentID, err := s.store.QuestionAssignmentID(r.Context(), questionID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy câu hỏi")
		return uuid.Nil, false
	}
	programID, err := s.store.NodeProgramID(r.Context(), assignmentID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return uuid.Nil, false
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return uuid.Nil, false
	}
	return assignmentID, true
}

type importQuestionsRequest struct {
	Questions []saveQuestionRequest `json:"questions"`
}

type importQuestionsResponse struct {
	Imported  int                `json:"imported"`
	Questions []*models.Question `json:"questions"`
}

// handleImportQuestions nhập nhiều câu hỏi một lần. Toàn bộ nằm trong một transaction:
// sai một câu thì không câu nào được ghi, tránh nhập dở dang.
func (s *Server) handleImportQuestions(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, _, ok := s.loadNode(w, r, nodeID, true)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment {
		writeError(w, http.StatusBadRequest, "Chỉ có thể nhập câu hỏi vào một bài tập")
		return
	}

	var req importQuestionsRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if len(req.Questions) == 0 {
		writeError(w, http.StatusBadRequest, "Không có câu hỏi nào để nhập")
		return
	}
	if len(req.Questions) > 500 {
		writeError(w, http.StatusBadRequest, "Mỗi lần chỉ nhập tối đa 500 câu hỏi")
		return
	}

	items := make([]store.SaveQuestionParams, 0, len(req.Questions))
	for i := range req.Questions {
		params, err := req.Questions[i].toParams(nodeID)
		if err != nil {
			msg, _ := invalidMessage(err)
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Câu thứ %d: %s", i+1, msg))
			return
		}
		params.Code = trimmed(req.Questions[i].Code)
		items = append(items, params)
	}

	count, err := s.store.ImportQuestions(r.Context(), nodeID, items)
	if err != nil {
		if msg, ok := invalidMessage(err); ok {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}

	questions, err := s.store.ListQuestions(r.Context(), nodeID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusCreated, importQuestionsResponse{Imported: count, Questions: questions})
}

// handleAssignmentResults trả về thống kê kết quả theo từng câu hỏi cho người quản lý.
// handleAssignmentResults trả về thống kê theo câu hỏi. Mở cho cả người quản lý
// lẫn Giám sát (chỉ xem, không chấm), nên không dùng loadNode(manage=true).
func (s *Server) handleAssignmentResults(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, err := s.store.GetNode(r.Context(), nodeID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return
	}
	if _, ok := s.requireProgramAudit(w, r, node.ProgramID); !ok {
		return
	}
	if node.Kind != models.KindAssignment {
		writeError(w, http.StatusBadRequest, "Nội dung này không phải bài tập")
		return
	}

	results, err := s.store.AssignmentResults(r.Context(), nodeID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, results)
}

func validQuestionType(t string) bool {
	switch t {
	case models.QuestionSingleChoice, models.QuestionMultiChoice, models.QuestionEssay:
		return true
	}
	return false
}
