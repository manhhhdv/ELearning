package export

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/manhnv/elearning/backend/internal/docx"
	"github.com/manhnv/elearning/backend/internal/models"
)

func TestProgramArchive(t *testing.T) {
	tree := []*models.Node{
		{Kind: models.KindLesson, Title: "Mở đầu", Lesson: &models.Lesson{ContentType: "richtext", Body: "Chào mừng"}},
		{Kind: models.KindFolder, Title: "Chương 1: Động học", Children: []*models.Node{
			{Kind: models.KindLesson, Title: "Tốc độ và vận tốc", Lesson: &models.Lesson{ContentType: "richtext", Body: "$v = s/t$"}},
			{Kind: models.KindAssignment, Title: "Bài tập chương 1", Assignment: &models.Assignment{Questions: sampleQuestions()}},
			// Bài tập chưa có câu hỏi thì bỏ qua.
			{Kind: models.KindAssignment, Title: "Bài tập trống", Assignment: &models.Assignment{}},
		}},
		// Thư mục rỗng không sinh mục nào trong tệp nén.
		{Kind: models.KindFolder, Title: "Chương 2"},
	}

	raw, err := ProgramArchive("Vật lý lớp 10", "", tree)
	if err != nil {
		t.Fatalf("ProgramArchive: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}

	var names []string
	contents := map[string]string{}
	for _, f := range zr.File {
		names = append(names, f.Name)
		rc, _ := f.Open()
		data, _ := io.ReadAll(rc)
		rc.Close()
		text, err := docx.ExtractText(data)
		if err != nil {
			t.Fatalf("%s is not a valid docx: %v", f.Name, err)
		}
		contents[f.Name] = text
	}

	want := []string{
		"vat-ly-lop-10/01-mo-dau.docx",
		"vat-ly-lop-10/02-chuong-1-dong-hoc/01-toc-do-va-van-toc.docx",
		"vat-ly-lop-10/02-chuong-1-dong-hoc/02-bai-tap-chuong-1-de.docx",
		"vat-ly-lop-10/02-chuong-1-dong-hoc/02-bai-tap-chuong-1-de-va-dap-an.docx",
	}
	if !slices.Equal(names, want) {
		t.Fatalf("entries = %v\nwant %v", names, want)
	}
	if strings.Contains(contents[want[2]], "ĐÁP ÁN") || !strings.Contains(contents[want[3]], "ĐÁP ÁN VÀ HƯỚNG DẪN CHẤM") {
		t.Errorf("only the answer-key file should contain the answer key")
	}
	if !strings.Contains(contents[want[0]], "Lớp học: Vật lý lớp 10") {
		t.Errorf("lesson should carry the program title: %q", contents[want[0]])
	}
}

func TestProgramArchiveEmpty(t *testing.T) {
	tree := []*models.Node{{Kind: models.KindFolder, Title: "Trống"}}
	if _, err := ProgramArchive("Lớp", "", tree); !errors.Is(err, ErrEmptyArchive) {
		t.Errorf("want ErrEmptyArchive, got %v", err)
	}
}
