package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/manhnv/elearning/backend/internal/ai"
	"github.com/manhnv/elearning/backend/internal/auth"
	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/store"
)

// Only server-loaded, accessible lessons may enter the model context.
func (s *Server) loadChatContext(w http.ResponseWriter, r *http.Request, programID, nodeID *uuid.UUID, query string) (string, []ai.CourseSource, bool) {
	if programID == nil {
		if nodeID != nil {
			writeError(w, http.StatusBadRequest, "Bài học phải thuộc một lớp học")
			return "", nil, false
		}
		return "", []ai.CourseSource{}, true
	}
	acc, ok := s.requireProgramAccess(w, r, *programID, false)
	if !ok {
		return "", nil, false
	}
	claims, _ := auth.FromContext(r.Context())
	p, err := s.store.GetProgram(r.Context(), *programID, claims.UserID)
	if err != nil {
		writeStoreError(w, err, "Không tìm thấy lớp học")
		return "", nil, false
	}
	if !acc.CanAudit && p.Status != "published" {
		writeError(w, http.StatusForbidden, "Lớp học chưa được xuất bản")
		return "", nil, false
	}
	nodes, err := s.store.ListNodes(r.Context(), *programID, false)
	if err != nil {
		writeStoreError(w, err, "")
		return "", nil, false
	}
	visible := chatVisibleNodes(store.BuildTree(nodes), acc.CanAudit)
	if nodeID != nil {
		found := false
		for _, n := range visible {
			if n.ID == *nodeID {
				found = true
				break
			}
		}
		if !found {
			writeError(w, http.StatusForbidden, "Bài học không thuộc lớp này hoặc không còn được phép truy cập")
			return "", nil, false
		}
	}
	content, sources := buildChatContext(p, visible, nodeID, query)
	return content, sources, true
}

func chatVisibleNodes(tree []*models.Node, audit bool) []*models.Node {
	out := []*models.Node{}
	var visit func([]*models.Node)
	visit = func(nodes []*models.Node) {
		for _, n := range nodes {
			if !audit && (!n.IsPublished || n.IsLocked) {
				continue
			}
			out = append(out, n)
			visit(n.Children)
		}
	}
	visit(tree)
	return out
}

const chatContextBudget = 24000

func chatExcerpt(value string, limit int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) > limit {
		return string(runes[:limit]) + "… [đã rút gọn]"
	}
	return string(runes)
}

// Rank accessible lessons by question terms, always prioritizing the open lesson.
// Assignment answers, grading rubrics and remote file contents are never included.
func buildChatContext(p *models.Program, nodes []*models.Node, nodeID *uuid.UUID, query string) (string, []ai.CourseSource) {
	type excerpt struct {
		ai.CourseSource
		Description string `json:"description"`
		Content     string `json:"content"`
	}
	type candidate struct {
		node  *models.Node
		score int
	}
	terms := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })
	ranked := []candidate{}
	current := ""
	for _, n := range nodes {
		if nodeID != nil && n.ID == *nodeID {
			current = n.Title
		}
		if n.Kind != models.KindLesson || n.Lesson == nil {
			continue
		}
		score := 0
		title := strings.ToLower(n.Title)
		body := strings.ToLower(n.Description + " " + n.Lesson.Body)
		for _, term := range terms {
			if len([]rune(term)) < 3 {
				continue
			}
			if strings.Contains(title, term) {
				score += 4
			}
			if strings.Contains(body, term) {
				score++
			}
		}
		if nodeID != nil && n.ID == *nodeID {
			score += 10000
		}
		ranked = append(ranked, candidate{n, score})
	}
	sort.SliceStable(ranked, func(i, j int) bool { return ranked[i].score > ranked[j].score })
	excerpts := []excerpt{}
	sources := []ai.CourseSource{}
	remaining := chatContextBudget
	for _, c := range ranked {
		if len(excerpts) == 6 || remaining <= 0 {
			break
		}
		n := c.node
		source := ai.CourseSource{NodeID: n.ID.String(), Title: chatExcerpt(n.Title, 200), URL: "/hoc/" + p.Slug + "/" + n.Slug, ContentAvailable: strings.TrimSpace(n.Lesson.Body) != ""}
		limit := 6000
		if remaining < limit {
			limit = remaining
		}
		content := chatExcerpt(n.Lesson.Body, limit)
		remaining -= len([]rune(content))
		excerpts = append(excerpts, excerpt{source, chatExcerpt(n.Description, 500), content})
		sources = append(sources, source)
	}
	data, _ := json.Marshal(struct {
		Course        string    `json:"course"`
		Description   string    `json:"description"`
		CurrentLesson string    `json:"currentLesson"`
		Lessons       []excerpt `json:"lessons"`
	}{chatExcerpt(p.Title, 200), chatExcerpt(p.Description, 1000), chatExcerpt(current, 200), excerpts})
	return string(data), sources
}

func sameChatScope(a, b *uuid.UUID) bool {
	return a == nil && b == nil || a != nil && b != nil && *a == *b
}
