package api

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildStructureNodeRejectsBadInput(t *testing.T) {
	programID := uuid.New()

	if _, msg := buildStructureNode(programID, structureItemInput{Kind: "lesson", Title: "  "}); msg == "" {
		t.Fatal("expected error for empty title")
	}
	if _, msg := buildStructureNode(programID, structureItemInput{Kind: "chapter", Title: "Chương 1"}); msg == "" {
		t.Fatal("expected error for unknown kind")
	}
	if _, msg := buildStructureNode(programID, structureItemInput{
		Kind: "lesson", Title: "Bài 1", ContentType: "video", Source: "không-phải-link",
	}); msg == "" {
		t.Fatal("expected error for unrecognized drive source")
	}
}

func TestBuildStructureNodeLesson(t *testing.T) {
	programID := uuid.New()
	params, msg := buildStructureNode(programID, structureItemInput{
		Kind: "lesson", Title: " Bài 1 ", ContentType: "slide", DurationMinutes: 30,
		Source: "https://drive.google.com/file/d/ABC123/view",
	})
	if msg != "" {
		t.Fatalf("unexpected error: %s", msg)
	}
	if params.Title != "Bài 1" || params.ProgramID != programID || !params.IsPublished {
		t.Fatalf("unexpected params: %+v", params)
	}
	if params.Lesson == nil || params.Lesson.ContentType != "slide" || params.Lesson.DurationMinutes != 30 {
		t.Fatalf("unexpected lesson: %+v", params.Lesson)
	}
	if params.Lesson.EmbedURL == "" {
		t.Fatal("expected embed URL to be built from drive link")
	}
}

func TestBuildStructureNodeFolderHasNoLesson(t *testing.T) {
	params, msg := buildStructureNode(uuid.New(), structureItemInput{Kind: "folder", Title: "Chương 1"})
	if msg != "" {
		t.Fatalf("unexpected error: %s", msg)
	}
	if params.Lesson != nil || params.Assignment != nil {
		t.Fatal("folder should not carry lesson/assignment payload")
	}
}
