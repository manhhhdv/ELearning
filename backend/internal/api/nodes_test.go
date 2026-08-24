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

func TestBuildLessonMaterials(t *testing.T) {
	lesson, err := buildLesson(&lessonInput{
		ContentType: "materials",
		Source:      "https://drive.google.com/file/d/ABC1234567/view",
		Attachments: []attachmentInput{
			{Name: " Giáo trình ", URL: " https://example.com/giao-trinh.pdf "},
			{Name: "", URL: ""},
			{Name: "Slide buổi 1", URL: "1a2b3c4d5e6f7g8h9i"},
			{Name: "Không tên", URL: "https://example.com/phu-luc.docx"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Bài tài liệu không nhúng file ngoài, kể cả khi form còn giữ link cũ.
	if lesson.EmbedURL != "" || lesson.DriveFileID != "" {
		t.Fatalf("expected no embed for materials lesson: %+v", lesson)
	}
	if len(lesson.Attachments) != 3 {
		t.Fatalf("expected 3 attachments (dòng trống bị bỏ), got %+v", lesson.Attachments)
	}
	if lesson.Attachments[0].Name != "Giáo trình" || lesson.Attachments[0].URL != "https://example.com/giao-trinh.pdf" {
		t.Fatalf("unexpected first attachment: %+v", lesson.Attachments[0])
	}
	// Dán mỗi ID file Drive vẫn ra link tải dùng được.
	if lesson.Attachments[1].URL != "https://drive.google.com/file/d/1a2b3c4d5e6f7g8h9i/view" {
		t.Fatalf("unexpected drive link: %+v", lesson.Attachments[1])
	}
}

func TestBuildLessonMaterialsRejectsBadLink(t *testing.T) {
	_, err := buildLesson(&lessonInput{
		ContentType: "materials",
		Attachments: []attachmentInput{{Name: "Giáo trình", URL: "chưa có link"}},
	})
	if err == nil {
		t.Fatal("expected error for unusable attachment link")
	}
}

func TestBuildAttachmentsFallsBackToLinkAsName(t *testing.T) {
	list, err := buildAttachments([]attachmentInput{{URL: "https://example.com/tai-lieu.pdf"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 || list[0].Name != "https://example.com/tai-lieu.pdf" {
		t.Fatalf("unexpected attachments: %+v", list)
	}
}
