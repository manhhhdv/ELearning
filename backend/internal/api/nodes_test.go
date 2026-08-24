package api

import "testing"

func TestBuildLessonRichText(t *testing.T) {
	body := "# Bài đọc\n\nCông thức $E = mc^2$."
	lesson, err := buildLesson(&lessonInput{
		ContentType:     "richtext",
		Source:          "https://drive.google.com/file/d/ABC1234567/view",
		DurationMinutes: 12,
		Body:            body,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lesson.Body != body || lesson.DurationMinutes != 12 {
		t.Fatalf("unexpected lesson: %+v", lesson)
	}
	// Bài tự soạn không nhúng file ngoài, kể cả khi form còn giữ link cũ.
	if lesson.EmbedURL != "" || lesson.DriveFileID != "" {
		t.Fatalf("expected no embed for richtext lesson: %+v", lesson)
	}
}

func TestBuildLessonRichTextAllowsEmptySource(t *testing.T) {
	if _, err := buildLesson(&lessonInput{ContentType: "richtext", Body: "Nội dung"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Loại nhúng vẫn bắt buộc link đọc được như trước.
	if _, err := buildLesson(&lessonInput{ContentType: "video", Source: "không-phải-link"}); err == nil {
		t.Fatal("expected error for unrecognized drive source")
	}
}
