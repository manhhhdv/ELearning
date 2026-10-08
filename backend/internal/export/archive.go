package export

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/util"
)

// ArchiveContentType là kiểu MIME của tệp .zip.
const ArchiveContentType = "application/zip"

// ErrEmptyArchive báo lớp học không có bài học hay bài tập có câu hỏi nào để xuất.
var ErrEmptyArchive = errors.New("lớp học chưa có bài học hay bài tập có câu hỏi nào để xuất")

// ArchiveName là tên tệp .zip của một lớp học, VD: "hoa-hoc-lop-10.zip".
func ArchiveName(programTitle string) string {
	return slugPart(programTitle) + ".zip"
}

// ProgramArchive đóng gói toàn bộ nội dung một lớp học thành tệp .zip: mỗi
// bài học một tệp Word, mỗi bài tập hai tệp (đề, và đề kèm đáp án). Thư mục
// trên cây trở thành thư mục trong tệp nén; tên tệp đánh số theo thứ tự trên
// cây và viết không dấu để giải nén bằng công cụ có sẵn của Windows không bị
// lỗi phông — tiêu đề có dấu vẫn nằm nguyên trong từng tệp.
//
// Bài tập trong tree phải có sẵn câu hỏi (Assignment.Questions).
func ProgramArchive(programTitle, author string, tree []*models.Node) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	now := time.Now()
	files := 0

	add := func(name string, data []byte) error {
		w, err := zw.CreateHeader(&zip.FileHeader{Name: name, Method: zip.Deflate, Modified: now})
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		files++
		return err
	}

	var walk func(dir string, nodes []*models.Node) error
	walk = func(dir string, nodes []*models.Node) error {
		for i, n := range nodes {
			name := dir + fmt.Sprintf("%02d-%s", i+1, slugPart(n.Title))
			switch n.Kind {
			case models.KindLesson:
				data, err := Lesson(programTitle, author, n)
				if err != nil {
					return fmt.Errorf("bài học %q: %w", n.Title, err)
				}
				if err := add(name+".docx", data); err != nil {
					return err
				}

			case models.KindAssignment:
				a := n.Assignment
				if a == nil || len(a.Questions) == 0 {
					break
				}
				info := QuestionSetInfo{
					ProgramTitle:     programTitle,
					Title:            n.Title,
					Instructions:     a.Instructions,
					TimeLimitMinutes: a.TimeLimitMinutes,
					Author:           author,
				}
				for _, withAnswers := range []bool{false, true} {
					data, err := Questions(info, a.Questions, withAnswers)
					if err != nil {
						return fmt.Errorf("bài tập %q: %w", n.Title, err)
					}
					suffix := "-de.docx"
					if withAnswers {
						suffix = "-de-va-dap-an.docx"
					}
					if err := add(name+suffix, data); err != nil {
						return err
					}
				}
			}
			if len(n.Children) > 0 {
				if err := walk(name+"/", n.Children); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Mọi tệp nằm trong một thư mục mang tên lớp học, giải nén ra không bị vương vãi.
	if err := walk(slugPart(programTitle)+"/", tree); err != nil {
		return nil, err
	}
	if files == 0 {
		return nil, ErrEmptyArchive
	}
	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("đóng gói tệp nén: %w", err)
	}
	return buf.Bytes(), nil
}

// slugPart là tiêu đề viết không dấu, cắt gọn để đường dẫn trong tệp nén không quá dài.
func slugPart(title string) string {
	slug := util.Slugify(title)
	if len(slug) > 60 {
		slug = strings.TrimRight(slug[:60], "-")
	}
	return slug
}
