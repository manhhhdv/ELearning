package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateLessonDraft(t *testing.T) {
	for _, tc := range []struct {
		name, output string
		valid        bool
	}{
		{"valid", `{"title":"Định luật Ôm","body":"## Mục tiêu\nHiểu mối quan hệ giữa cường độ dòng điện, hiệu điện thế và điện trở.\n## Kiến thức\nCông thức I = U/R.\n## Ví dụ\nVới U = 12 V và R = 6 ôm thì I = 2 A.\n## Tự kiểm tra\nKhi điện trở tăng thì cường độ dòng điện thay đổi thế nào?"}`, true},
		{"empty", `{"title":"","body":""}`, false},
		{"invalid json", "not json", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var payload geminiRequest
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Error(err)
				}
				json.NewEncoder(w).Encode(map[string]any{"candidates": []any{map[string]any{"content": map[string]any{"parts": []any{map[string]string{"text": tc.output}}}}}})
			}))
			defer server.Close()
			client := NewRotatingClient([]Credential{{Provider: "gemini", APIKey: "test-key", Model: "test-model", BaseURL: server.URL}})
			draft, err := client.GenerateLesson(context.Background(), LessonParams{Course: "Vật lý", Topic: "Định luật Ôm", DurationMinutes: 45})
			if tc.valid && (err != nil || draft == nil || !strings.Contains(draft.Body, "I = U/R")) {
				t.Fatalf("draft=%+v err=%v", draft, err)
			}
			if !tc.valid && err == nil {
				t.Fatal("invalid draft accepted")
			}
		})
	}
}

func TestCourseAssistantPromptRequiresGrounding(t *testing.T) {
	prompt := assistantSystemPrompt(ChatParams{Persona: PersonaStudent, CourseContext: `{"course":"test"}`})
	for _, wanted := range []string{"PHẠM VI LỚP HỌC", "contentAvailable=false", "liên kết Markdown", "không phải toàn bộ lớp học"} {
		if !strings.Contains(prompt, wanted) {
			t.Errorf("missing %s", wanted)
		}
	}
	if strings.Contains(assistantSystemPrompt(ChatParams{}), "PHẠM VI LỚP HỌC") {
		t.Fatal("general chat should remain available")
	}
}
