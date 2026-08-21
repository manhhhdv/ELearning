package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// maxStructureImportItems chặn file quá lớn tạo ra hàng nghìn nút trong một lần gọi.
const maxStructureImportItems = 500

// structureItemInput là một dòng trong file cấu trúc: Level cho biết nút này nằm
// ở cấp nào (1 = ngay dưới nút gốc được chọn), nhờ đó tái dựng lại cây từ bảng phẳng.
type structureItemInput struct {
	Level           int    `json:"level"`
	Kind            string `json:"kind"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ContentType     string `json:"contentType"`
	Source          string `json:"source"`
	DurationMinutes int    `json:"durationMinutes"`
}

type importStructureRequest struct {
	ParentID *uuid.UUID           `json:"parentId"`
	Items    []structureItemInput `json:"items"`
}

type importStructureResponse struct {
	Imported int `json:"imported"`
}

// handleImportStructureFile đọc file .xlsx/.csv thành văn bản dạng bảng (ô cách
// nhau bằng Tab) để giao diện dùng chung một bộ phân tích với luồng dán tay.
func (s *Server) handleImportStructureFile(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}

	rows, errMsg := readUploadedTable(w, r)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	writeJSON(w, http.StatusOK, importQuestionsFileResponse{Text: rowsToTabText(rows)})
}

// handleImportStructure tạo hàng loạt nút theo đúng thứ tự gửi lên. Mỗi mục được
// gắn vào thư mục gần nhất ở cấp trên nó, nên chỉ cần một danh sách phẳng kèm Level.
func (s *Server) handleImportStructure(w http.ResponseWriter, r *http.Request) {
	programID, ok := urlUUID(w, r, "programID")
	if !ok {
		return
	}
	if _, ok := s.requireProgramAccess(w, r, programID, true); !ok {
		return
	}

	var req importStructureRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Không có mục nào để nhập")
		return
	}
	if len(req.Items) > maxStructureImportItems {
		writeError(w, http.StatusBadRequest, "Mỗi lần chỉ nhập được tối đa 500 mục")
		return
	}
	if req.ParentID != nil {
		parent, err := s.store.GetNode(r.Context(), *req.ParentID)
		if err != nil {
			writeStoreError(w, err, "Không tìm thấy thư mục cha")
			return
		}
		if parent.ProgramID != programID || parent.Kind != models.KindFolder {
			writeError(w, http.StatusBadRequest, "Chỉ có thể nhập vào một thư mục của chương trình này")
			return
		}
	}

	// parents[i] là nút cha cho các mục ở cấp i+1; parents[0] là gốc được chọn.
	parents := []*uuid.UUID{req.ParentID}

	for i, item := range req.Items {
		params, errMsg := buildStructureNode(programID, item)
		if errMsg != "" {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("Dòng %d: %s", i+1, errMsg))
			return
		}

		level := item.Level
		if level < 1 {
			level = 1
		}
		if level > len(parents) {
			// Nhảy cấp (ví dụ từ cấp 1 xuống cấp 3) thì coi như nằm ngay dưới cấp sâu nhất đang có.
			level = len(parents)
		}
		params.ParentID = parents[level-1]

		node, err := s.store.CreateNode(r.Context(), params)
		if err != nil {
			if msg, ok := invalidMessage(err); ok {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("Dòng %d: %s", i+1, msg))
				return
			}
			writeStoreError(w, err, "Không tìm thấy thư mục cha")
			return
		}

		// Chỉ thư mục mới nhận con; cắt bớt các cấp sâu hơn để cấp sau bám đúng nút vừa tạo.
		parents = parents[:level]
		if node.Kind == models.KindFolder {
			id := node.ID
			parents = append(parents, &id)
		}
	}

	writeJSON(w, http.StatusCreated, importStructureResponse{Imported: len(req.Items)})
}

func buildStructureNode(programID uuid.UUID, item structureItemInput) (store.SaveNodeParams, string) {
	title := trimmed(item.Title)
	if title == "" {
		return store.SaveNodeParams{}, "thiếu tiêu đề"
	}
	if !validNodeKind(item.Kind) {
		return store.SaveNodeParams{}, "loại nội dung không hợp lệ"
	}

	params := store.SaveNodeParams{
		ProgramID:   programID,
		Kind:        item.Kind,
		Title:       title,
		Description: trimmed(item.Description),
		IsPublished: true,
	}

	switch item.Kind {
	case models.KindLesson:
		lesson, err := buildLesson(&lessonInput{
			ContentType:     item.ContentType,
			Source:          item.Source,
			DurationMinutes: item.DurationMinutes,
		})
		if err != nil {
			return store.SaveNodeParams{}, err.Error()
		}
		params.Lesson = lesson
	case models.KindAssignment:
		params.Assignment = buildAssignment(nil)
	}
	return params, ""
}
