package api

import (
	"github.com/google/uuid"
	"github.com/manhnv/elearning/backend/internal/models"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestChatContextFiltersHiddenAndLockedBranches(t *testing.T) {
	visible := &models.Node{ID: uuid.New(), IsPublished: true, Title: "visible", Kind: models.KindLesson, Lesson: &models.Lesson{Body: "Kiến thức được phép đọc"}}
	hidden := &models.Node{ID: uuid.New(), Title: "hidden", Children: []*models.Node{{ID: uuid.New(), IsPublished: true, Title: "secret child"}}}
	locked := &models.Node{ID: uuid.New(), IsPublished: true, IsLocked: true, Title: "locked", Children: []*models.Node{{ID: uuid.New(), IsPublished: true, Title: "locked child"}}}
	nodes := []*models.Node{visible, hidden, locked}
	got := chatVisibleNodes(nodes, false)
	if len(got) != 1 || got[0].ID != visible.ID {
		t.Fatalf("restricted nodes leaked: %+v", got)
	}
	if len(chatVisibleNodes(nodes, true)) != 5 {
		t.Fatal("teacher audit context should include drafts")
	}
}

func TestChatContextPrioritizesCurrentAndExcludesAnswers(t *testing.T) {
	p := &models.Program{Title: "Vật lý", Slug: "vat-ly"}
	nodes := []*models.Node{}
	for i := 0; i < 10; i++ {
		nodes = append(nodes, &models.Node{ID: uuid.New(), Kind: models.KindLesson, Title: "Điện trở", Slug: "dien-tro", Lesson: &models.Lesson{Body: strings.Repeat("điện trở ", 2000)}})
	}
	current := nodes[9].ID
	nodes = append(nodes, &models.Node{ID: uuid.New(), Kind: models.KindAssignment, Assignment: &models.Assignment{Instructions: "secret instructions", Questions: []*models.Question{{SampleAnswer: "SECRET_ANSWER", Rubric: "SECRET_RUBRIC"}}}})
	content, sources := buildChatContext(p, nodes, &current, "điện trở")
	if len(sources) == 0 || sources[0].NodeID != current.String() {
		t.Fatal("open lesson must be first")
	}
	if len(sources) > 6 || len([]rune(content)) > 30000 || !utf8.ValidString(content) {
		t.Fatal("context is not bounded UTF-8")
	}
	if strings.Contains(content, "SECRET") || strings.Contains(content, "secret instructions") {
		t.Fatal("assessment content leaked")
	}
	if sources[0].URL != "/hoc/vat-ly/dien-tro" {
		t.Fatal("invalid source URL")
	}
}

func TestChatContextMarksUnreadExternalMedia(t *testing.T) {
	id := uuid.New()
	content, sources := buildChatContext(&models.Program{Slug: "course"}, []*models.Node{{ID: id, Kind: models.KindLesson, Lesson: &models.Lesson{ContentType: "video", EmbedURL: "https://example.com/private"}}}, &id, "")
	if len(sources) != 1 || sources[0].ContentAvailable {
		t.Fatal("video must not be presented as read")
	}
	if strings.Contains(content, "https://example.com/private") {
		t.Fatal("external URL should not be fetched or sent as knowledge")
	}
}

func TestChatScope(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if !sameChatScope(nil, nil) || !sameChatScope(&a, &a) || sameChatScope(&a, nil) || sameChatScope(nil, &a) || sameChatScope(&a, &b) {
		t.Fatal("scope matching failed")
	}
}
