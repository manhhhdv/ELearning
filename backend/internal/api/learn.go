package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// ---------------------------------------------------------------------------
// Ghi danh
// ---------------------------------------------------------------------------

func (s *Server) handleListEnrollments(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAudit(w, r, programID); !ok {
		return
	}
	items, err := s.store.ListProgramEnrollments(r.Context(), programID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type enrollRequest struct {
	UserIDs []uuid.UUID `json:"userIds"`
	Role    string      `json:"role"`
}

func (s *Server) handleEnroll(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}

	var req enrollRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Role == "" {
		req.Role = models.RoleStudent
	}
	if req.Role != models.RoleStudent && req.Role != models.RoleTrainer {
		writeError(w, http.StatusBadRequest, "Vai trò ghi danh chỉ có thể là học viên hoặc giảng viên")
		return
	}
	if len(req.UserIDs) == 0 {
		writeError(w, http.StatusBadRequest, "Vui lòng chọn ít nhất một người dùng")
		return
	}

	for _, userID := range req.UserIDs {
		if err := s.store.Enroll(r.Context(), programID, userID, req.Role); err != nil {
			writeStoreError(w, err, "Không tìm thấy người dùng")
			return
		}
	}
	items, err := s.store.ListProgramEnrollments(r.Context(), programID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleUnenroll(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	userID, ok := urlUUID(w, r, "userID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}
	if err := s.store.Unenroll(r.Context(), programID, userID); err != nil {
		writeStoreError(w, err, "Người dùng chưa được ghi danh")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Khu vực học viên
// ---------------------------------------------------------------------------

func (s *Server) handleMyPrograms(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	programs, err := s.store.ListPrograms(r.Context(), store.ListProgramsFilter{
		EnrolledUserID: claims.UserID,
		ViewerID:       claims.UserID,
		Status:         "published",
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, programs)
}

// handleCatalog liệt kê các khoá đang mở cho học viên tự ghi danh.
func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	programs, err := s.store.ListPrograms(r.Context(), store.ListProgramsFilter{
		Status:             "published",
		OnlySelfEnrollable: true,
		ViewerID:           claims.UserID,
		Search:             r.URL.Query().Get("search"),
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, programs)
}

// handleSelfEnroll cho học viên tự ghi danh vào khoá đã bật tuỳ chọn này.
func (s *Server) handleSelfEnroll(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	program, err := s.store.GetProgram(r.Context(), programID, uuid.Nil)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	if !program.AllowSelfEnroll {
		writeError(w, http.StatusForbidden, "Khoá học này không cho tự ghi danh, vui lòng liên hệ quản trị viên")
		return
	}
	if program.Status != "published" {
		writeError(w, http.StatusForbidden, "Khoá học này chưa được mở")
		return
	}

	claims, _ := auth.FromContext(r.Context())
	if err := s.store.Enroll(r.Context(), programID, claims.UserID, models.RoleStudent); err != nil {
		writeStoreError(w, err, "Không ghi danh được")
		return
	}

	enrolled, err := s.store.GetProgram(r.Context(), programID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, enrolled)
}

// handleSelfUnenroll cho học viên tự rời khoá mình đã tự ghi danh.
func (s *Server) handleSelfUnenroll(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	program, err := s.store.GetProgram(r.Context(), programID, uuid.Nil)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}
	// Khoá do admin xếp lớp thì học viên không được tự rút.
	if !program.AllowSelfEnroll {
		writeError(w, http.StatusForbidden, "Khoá học này do quản trị viên xếp, bạn không thể tự rời")
		return
	}

	claims, _ := auth.FromContext(r.Context())
	if err := s.store.Unenroll(r.Context(), programID, claims.UserID); err != nil {
		writeStoreError(w, err, "Bạn chưa ghi danh khoá này")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type markCompleteRequest struct {
	Completed bool `json:"completed"`
}

func (s *Server) handleMarkComplete(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	if _, _, ok := s.loadNode(w, r, nodeID, false); !ok {
		return
	}

	var req markCompleteRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	if err := s.store.MarkLessonComplete(r.Context(), claims.UserID, nodeID, req.Completed); err != nil {
		writeStoreError(w, err, "Không tìm thấy nội dung")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"completed": req.Completed})
}

type attemptResponse struct {
	Node *models.Node `json:"node"`
	// Số lượt đã nộp và số lượt tối đa, để giao diện biết còn được làm bài nữa không.
	AttemptsUsed int                  `json:"attemptsUsed"`
	MaxAttempts  int                  `json:"maxAttempts"`
	Submissions  []*models.Submission `json:"submissions"`
	// Lượt đang làm dở, nil khi học viên chưa bấm bắt đầu.
	Session *store.AttemptSession `json:"session"`
}

// handleGetAttempt trả về đề bài cho học viên (đã ẩn đáp án) kèm lịch sử làm bài.
func (s *Server) handleGetAttempt(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, acc, ok := s.loadNode(w, r, nodeID, false)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment || node.Assignment == nil {
		writeError(w, http.StatusBadRequest, "Nội dung này không phải bài tập")
		return
	}
	if !acc.CanManage && !node.IsPublished {
		writeError(w, http.StatusForbidden, "Bài tập này chưa được xuất bản")
		return
	}
	stripAnswers(node)

	claims, _ := auth.FromContext(r.Context())
	used, err := s.store.CountAttempts(r.Context(), nodeID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	history, err := s.store.ListSubmissions(r.Context(), store.ListSubmissionsFilter{
		AssignmentID: nodeID,
		UserID:       claims.UserID,
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	session, err := s.store.OpenAttempt(r.Context(), nodeID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	// Phiên đã quá giờ coi như không còn: giao diện sẽ mời học viên bắt đầu lượt mới.
	if session != nil && session.Expired(time.Now()) {
		session = nil
	}
	if session != nil && node.Assignment.ShuffleQuestions {
		store.ShuffleForSession(node.Assignment.Questions, session.ID)
	}

	writeJSON(w, http.StatusOK, attemptResponse{
		Node:         node,
		AttemptsUsed: used,
		MaxAttempts:  node.Assignment.MaxAttempts,
		Submissions:  history,
		Session:      session,
	})
}

// handleStartAttempt mở lượt làm bài và trả về đề đã sắp xếp theo đúng lượt đó.
func (s *Server) handleStartAttempt(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, acc, ok := s.loadNode(w, r, nodeID, false)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment || node.Assignment == nil {
		writeError(w, http.StatusBadRequest, "Nội dung này không phải bài tập")
		return
	}
	if !acc.CanManage && !node.IsPublished {
		writeError(w, http.StatusForbidden, "Bài tập này chưa được xuất bản")
		return
	}
	if len(node.Assignment.Questions) == 0 {
		writeError(w, http.StatusBadRequest, "Bài tập này chưa có câu hỏi nào")
		return
	}

	claims, _ := auth.FromContext(r.Context())
	session, err := s.store.StartAttempt(r.Context(), nodeID, claims.UserID)
	if err != nil {
		if msg, ok := invalidMessage(err); ok {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}

	stripAnswers(node)
	if node.Assignment.ShuffleQuestions {
		store.ShuffleForSession(node.Assignment.Questions, session.ID)
	}

	used, err := s.store.CountAttempts(r.Context(), nodeID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	history, err := s.store.ListSubmissions(r.Context(), store.ListSubmissionsFilter{
		AssignmentID: nodeID,
		UserID:       claims.UserID,
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	writeJSON(w, http.StatusOK, attemptResponse{
		Node:         node,
		AttemptsUsed: used,
		MaxAttempts:  node.Assignment.MaxAttempts,
		Submissions:  history,
		Session:      session,
	})
}

type submitRequest struct {
	Answers []struct {
		QuestionID        uuid.UUID   `json:"questionId"`
		SelectedOptionIDs []uuid.UUID `json:"selectedOptionIds"`
		EssayText         string      `json:"essayText"`
	} `json:"answers"`
}

func (s *Server) handleSubmitAssignment(w http.ResponseWriter, r *http.Request) {
	nodeID, ok := urlUUID(w, r, "nodeID")
	if !ok {
		return
	}
	node, _, ok := s.loadNode(w, r, nodeID, false)
	if !ok {
		return
	}
	if node.Kind != models.KindAssignment {
		writeError(w, http.StatusBadRequest, "Nội dung này không phải bài tập")
		return
	}

	var req submitRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	answers := make([]store.AnswerInput, 0, len(req.Answers))
	for _, a := range req.Answers {
		answers = append(answers, store.AnswerInput{
			QuestionID:        a.QuestionID,
			SelectedOptionIDs: a.SelectedOptionIDs,
			EssayText:         a.EssayText,
		})
	}

	claims, _ := auth.FromContext(r.Context())
	submission, err := s.store.SubmitAssignment(r.Context(), nodeID, claims.UserID, answers)
	if err != nil {
		if msg, ok := invalidMessage(err); ok {
			writeError(w, http.StatusBadRequest, msg)
			return
		}
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}
	// Bài chưa chấm xong thì chưa cho học viên thấy điểm từng câu.
	writeJSON(w, http.StatusCreated, submission)
}

func (s *Server) handleMySubmissions(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	items, err := s.store.ListSubmissions(r.Context(), store.ListSubmissionsFilter{UserID: claims.UserID})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// ---------------------------------------------------------------------------
// Bài nộp và chấm bài
// ---------------------------------------------------------------------------

func (s *Server) handleListProgramSubmissions(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAudit(w, r, programID); !ok {
		return
	}
	items, err := s.store.ListSubmissions(r.Context(), store.ListSubmissionsFilter{
		ProgramID:   programID,
		OnlyPending: r.URL.Query().Get("pending") == "true",
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type submissionDetailResponse struct {
	Submission *models.Submission `json:"submission"`
	Questions  []*models.Question `json:"questions"`
	CanGrade   bool               `json:"canGrade"`
}

func (s *Server) handleGetSubmission(w http.ResponseWriter, r *http.Request) {
	submissionID, ok := urlUUID(w, r, "submissionID")
	if !ok {
		return
	}
	submission, err := s.store.GetSubmission(r.Context(), submissionID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài nộp")
		return
	}

	programID, err := s.store.NodeProgramID(r.Context(), submission.AssignmentID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}
	acc, err := s.programAccess(r, programID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy chương trình")
		return
	}

	claims, _ := auth.FromContext(r.Context())
	isOwner := submission.UserID == claims.UserID
	canReview := acc.CanManage || acc.CanAudit
	if !canReview && !isOwner {
		writeError(w, http.StatusForbidden, "Bạn không có quyền xem bài nộp này")
		return
	}

	questions, err := s.store.ListQuestions(r.Context(), submission.AssignmentID)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	// Học viên chỉ thấy đáp án đúng sau khi bài đã được chấm xong; người xem có
	// quyền quản lý hoặc giám sát thì thấy ngay, không phải chờ chấm.
	if !canReview && submission.Status != "graded" {
		for _, q := range questions {
			q.Explanation = ""
			for _, o := range q.Options {
				o.IsCorrect = false
			}
		}
	}

	writeJSON(w, http.StatusOK, submissionDetailResponse{
		Submission: submission,
		Questions:  questions,
		CanGrade:   acc.CanManage,
	})
}

type gradeRequest struct {
	Feedback string `json:"feedback"`
	Answers  []struct {
		AnswerID uuid.UUID `json:"answerId"`
		Score    float64   `json:"score"`
		Comment  string    `json:"comment"`
	} `json:"answers"`
}

func (s *Server) handleGradeSubmission(w http.ResponseWriter, r *http.Request) {
	submissionID, ok := urlUUID(w, r, "submissionID")
	if !ok {
		return
	}
	submission, err := s.store.GetSubmission(r.Context(), submissionID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài nộp")
		return
	}
	programID, err := s.store.NodeProgramID(r.Context(), submission.AssignmentID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài tập")
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}

	var req gradeRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	answers := make([]store.GradeAnswerInput, 0, len(req.Answers))
	for _, a := range req.Answers {
		if a.Score < 0 {
			writeError(w, http.StatusBadRequest, "Điểm không được là số âm")
			return
		}
		answers = append(answers, store.GradeAnswerInput{AnswerID: a.AnswerID, Score: a.Score, Comment: a.Comment})
	}

	claims, _ := auth.FromContext(r.Context())
	graded, err := s.store.GradeSubmission(r.Context(), submissionID, claims.UserID, req.Feedback, answers)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy bài nộp")
		return
	}
	writeJSON(w, http.StatusOK, graded)
}
