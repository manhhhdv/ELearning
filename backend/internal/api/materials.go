package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
	"github.com/manhnv/elearning/backend/internal/util"
)

// canManageMaterials cho biết người dùng có quyền thêm/sửa/xoá tài liệu dùng
// chung không. Quản trị viên và giáo viên được phép; học viên và vai trò
// Giám sát chỉ xem.
func canManageMaterials(role string) bool {
	return role == models.RoleAdmin || role == models.RoleTrainer
}

type materialRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	URL         string `json:"url"`
	Kind        string `json:"kind"`
	IsPublished bool   `json:"isPublished"`
}

func validMaterialKind(kind string) bool {
	switch kind {
	case models.MaterialPDF, models.MaterialSlide, models.MaterialDocument,
		models.MaterialVideo, models.MaterialLink:
		return true
	}
	return false
}

// toParams kiểm tra dữ liệu và tách sẵn ID file Drive từ đường dẫn.
func (req *materialRequest) toParams(userID uuid.UUID) (store.SaveMaterialParams, error) {
	title := trimmed(req.Title)
	if title == "" {
		return store.SaveMaterialParams{}, errValidation("Vui lòng nhập tên tài liệu")
	}
	rawURL := trimmed(req.URL)
	if rawURL == "" {
		return store.SaveMaterialParams{}, errValidation("Vui lòng nhập đường dẫn tài liệu")
	}
	kind := req.Kind
	if kind == "" {
		kind = models.MaterialLink
	}
	if !validMaterialKind(kind) {
		return store.SaveMaterialParams{}, errValidation("Loại tài liệu không hợp lệ")
	}

	// Link Drive được tách ID để giao diện nhúng xem trước; link ngoài giữ nguyên.
	driveID, embedURL := util.BuildEmbedURL(kind, rawURL)
	if driveID == "" && embedURL == "" {
		return store.SaveMaterialParams{}, errValidation(
			"Đường dẫn không hợp lệ. Hãy dán link Google Drive hoặc link http(s).")
	}

	return store.SaveMaterialParams{
		Title:       title,
		Description: trimmed(req.Description),
		Category:    trimmed(req.Category),
		URL:         rawURL,
		DriveFileID: driveID,
		Kind:        kind,
		IsPublished: req.IsPublished,
		CreatedBy:   userID,
	}, nil
}

// handleListMaterials trả về kho tài liệu. Học viên chỉ thấy tài liệu đã xuất
// bản; người quản lý thấy cả bản nháp để sửa tiếp.
func (s *Server) handleListMaterials(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.FromContext(r.Context())
	items, err := s.store.ListMaterials(r.Context(), store.ListMaterialsFilter{
		Search:             r.URL.Query().Get("search"),
		Category:           r.URL.Query().Get("category"),
		IncludeUnpublished: canManageMaterials(claims.Role),
	})
	if err != nil {
		writeStoreError(w, err, "")
		return
	}

	categories, err := s.store.MaterialCategories(r.Context())
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":      items,
		"categories": categories,
		"canManage":  canManageMaterials(claims.Role),
	})
}

func (s *Server) handleCreateMaterial(w http.ResponseWriter, r *http.Request) {
	var req materialRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	params, err := req.toParams(claims.UserID)
	if err != nil {
		msg, _ := invalidMessage(err)
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	m, err := s.store.CreateMaterial(r.Context(), params)
	if err != nil {
		writeStoreError(w, err, "")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (s *Server) handleUpdateMaterial(w http.ResponseWriter, r *http.Request) {
	materialID, ok := urlUUID(w, r, "materialID")
	if !ok {
		return
	}
	var req materialRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	claims, _ := auth.FromContext(r.Context())
	params, err := req.toParams(claims.UserID)
	if err != nil {
		msg, _ := invalidMessage(err)
		writeError(w, http.StatusBadRequest, msg)
		return
	}

	m, err := s.store.UpdateMaterial(r.Context(), materialID, params)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy tài liệu")
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleDeleteMaterial(w http.ResponseWriter, r *http.Request) {
	materialID, ok := urlUUID(w, r, "materialID")
	if !ok {
		return
	}
	if err := s.store.DeleteMaterial(r.Context(), materialID); err != nil {
		writeStoreError(w, err, "Không tìm thấy tài liệu")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
