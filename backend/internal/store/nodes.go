package store

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/manhnv/elearning/backend/internal/models"
	"github.com/manhnv/elearning/backend/internal/util"
)

const nodeColumns = `
	n.id, n.program_id, n.parent_id, n.kind, n.slug, n.title, n.description, n.position,
	n.is_published, n.is_locked, n.created_at, n.updated_at`

type nodeScanner interface {
	Scan(dest ...any) error
}

func scanNode(row nodeScanner) (*models.Node, error) {
	var n models.Node
	err := row.Scan(&n.ID, &n.ProgramID, &n.ParentID, &n.Kind, &n.Slug, &n.Title, &n.Description,
		&n.Position, &n.IsPublished, &n.IsLocked, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	n.Children = []*models.Node{}
	return &n, nil
}

// uniqueNodeSlug thêm hậu tố -2, -3… nếu slug cơ sở đã tồn tại trong cùng chương trình.
// Phạm vi duy nhất là theo chương trình (không toàn hệ thống) vì URL luôn đi kèm mã chương trình.
func uniqueNodeSlug(ctx context.Context, tx pgx.Tx, programID uuid.UUID, base string) (string, error) {
	candidate := base
	for i := 2; ; i++ {
		var exists bool
		if err := tx.QueryRow(ctx, `
			SELECT EXISTS (SELECT 1 FROM nodes WHERE program_id = $1 AND slug = $2)`,
			programID, candidate).Scan(&exists); err != nil {
			return "", translate(err, "kiểm tra slug")
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// ListNodes trả về toàn bộ nút của một chương trình kèm chi tiết bài học / bài tập.
// publishedOnly dùng cho phía học viên: chỉ lấy nút đã xuất bản.
func (s *Store) ListNodes(ctx context.Context, programID uuid.UUID, publishedOnly bool) ([]*models.Node, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+nodeColumns+`
		FROM nodes n
		WHERE n.program_id = $1 AND ($2 = false OR n.is_published)
		ORDER BY n.position, n.created_at`, programID, publishedOnly)
	if err != nil {
		return nil, translate(err, "liệt kê nút")
	}
	defer rows.Close()

	nodes := []*models.Node{}
	byID := map[uuid.UUID]*models.Node{}
	for rows.Next() {
		n, err := scanNode(rows)
		if err != nil {
			return nil, translate(err, "đọc nút")
		}
		nodes = append(nodes, n)
		byID[n.ID] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}

	if err := s.attachLessons(ctx, programID, byID); err != nil {
		return nil, err
	}
	if err := s.attachAssignments(ctx, programID, byID); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (s *Store) attachLessons(ctx context.Context, programID uuid.UUID, byID map[uuid.UUID]*models.Node) error {
	rows, err := s.pool.Query(ctx, `
		SELECT l.node_id, l.content_type, l.drive_file_id, l.embed_url, l.duration_minutes, l.body, l.attachments
		FROM lessons l JOIN nodes n ON n.id = l.node_id
		WHERE n.program_id = $1`, programID)
	if err != nil {
		return translate(err, "đọc bài học")
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var l models.Lesson
		if err := rows.Scan(&id, &l.ContentType, &l.DriveFileID, &l.EmbedURL, &l.DurationMinutes, &l.Body, &l.Attachments); err != nil {
			return translate(err, "đọc bài học")
		}
		l.Attachments = lessonAttachments(l.Attachments)
		if n, ok := byID[id]; ok {
			n.Lesson = &l
		}
	}
	return rows.Err()
}

// lessonAttachments đảm bảo cột jsonb luôn là mảng, không phải null,
// để JSON trả về cho giao diện luôn có dạng danh sách.
func lessonAttachments(list []models.LessonAttachment) []models.LessonAttachment {
	if list == nil {
		return []models.LessonAttachment{}
	}
	return list
}

func (s *Store) attachAssignments(ctx context.Context, programID uuid.UUID, byID map[uuid.UUID]*models.Node) error {
	rows, err := s.pool.Query(ctx, `
		SELECT a.node_id, a.instructions, a.time_limit_minutes, a.max_attempts, a.pass_score,
		       a.shuffle_questions, a.due_at,
		       COALESCE((SELECT count(*) FROM questions q WHERE q.assignment_id = a.node_id), 0)
		FROM assignments a JOIN nodes n ON n.id = a.node_id
		WHERE n.program_id = $1`, programID)
	if err != nil {
		return translate(err, "đọc bài tập")
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var a models.Assignment
		if err := rows.Scan(&id, &a.Instructions, &a.TimeLimitMinutes, &a.MaxAttempts, &a.PassScore,
			&a.ShuffleQuestions, &a.DueAt, &a.QuestionCount); err != nil {
			return translate(err, "đọc bài tập")
		}
		if n, ok := byID[id]; ok {
			n.Assignment = &a
		}
	}
	return rows.Err()
}

// BuildTree ghép danh sách phẳng thành cây, sắp xếp anh em theo position.
func BuildTree(nodes []*models.Node) []*models.Node {
	byID := make(map[uuid.UUID]*models.Node, len(nodes))
	for _, n := range nodes {
		n.Children = []*models.Node{}
		byID[n.ID] = n
	}

	roots := []*models.Node{}
	for _, n := range nodes {
		if n.ParentID == nil {
			roots = append(roots, n)
			continue
		}
		parent, ok := byID[*n.ParentID]
		if !ok {
			// Cha bị ẩn (chế độ publishedOnly) thì nút con cũng không hiển thị.
			continue
		}
		parent.Children = append(parent.Children, n)
	}

	var sortRec func(list []*models.Node)
	sortRec = func(list []*models.Node) {
		sort.SliceStable(list, func(i, j int) bool { return list[i].Position < list[j].Position })
		for _, n := range list {
			sortRec(n.Children)
		}
	}
	sortRec(roots)
	return roots
}

func (s *Store) GetNode(ctx context.Context, id uuid.UUID) (*models.Node, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+nodeColumns+` FROM nodes n WHERE n.id = $1`, id)
	n, err := scanNode(row)
	if err != nil {
		return nil, translate(err, "đọc nút")
	}

	switch n.Kind {
	case models.KindLesson:
		var l models.Lesson
		err = s.pool.QueryRow(ctx, `
			SELECT content_type, drive_file_id, embed_url, duration_minutes, body, attachments
			FROM lessons WHERE node_id = $1`, id).
			Scan(&l.ContentType, &l.DriveFileID, &l.EmbedURL, &l.DurationMinutes, &l.Body, &l.Attachments)
		if err != nil && err != pgx.ErrNoRows {
			return nil, translate(err, "đọc bài học")
		}
		if err == nil {
			l.Attachments = lessonAttachments(l.Attachments)
			n.Lesson = &l
		}
	case models.KindAssignment:
		a, err := s.getAssignment(ctx, id)
		if err != nil && err != ErrNotFound {
			return nil, err
		}
		n.Assignment = a
	}
	return n, nil
}

type SaveNodeParams struct {
	ProgramID   uuid.UUID
	ParentID    *uuid.UUID
	Kind        string
	Title       string
	Description string
	IsPublished bool
	IsLocked    bool

	Lesson     *models.Lesson
	Assignment *models.Assignment
}

func (s *Store) CreateNode(ctx context.Context, p SaveNodeParams) (*models.Node, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if p.ParentID != nil {
		var parentKind string
		var parentProgram uuid.UUID
		err := tx.QueryRow(ctx, `SELECT kind, program_id FROM nodes WHERE id = $1`, *p.ParentID).
			Scan(&parentKind, &parentProgram)
		if err != nil {
			return nil, translate(err, "đọc nút cha")
		}
		if parentKind != models.KindFolder {
			return nil, Invalidf("Chỉ có thể thêm nội dung vào bên trong một thư mục")
		}
		if parentProgram != p.ProgramID {
			return nil, Invalidf("Mục cha không thuộc chương trình này")
		}
	}

	var position int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(max(position) + 1, 0) FROM nodes
		WHERE program_id = $1 AND parent_id IS NOT DISTINCT FROM $2`, p.ProgramID, p.ParentID).Scan(&position)
	if err != nil {
		return nil, translate(err, "tính thứ tự nút")
	}

	slug, err := uniqueNodeSlug(ctx, tx, p.ProgramID, util.Slugify(p.Title))
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO nodes (program_id, parent_id, kind, slug, title, description, position, is_published, is_locked)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		p.ProgramID, p.ParentID, p.Kind, slug, p.Title, p.Description, position, p.IsPublished, p.IsLocked).Scan(&id)
	if err != nil {
		return nil, translate(err, "tạo nút")
	}

	switch p.Kind {
	case models.KindLesson:
		l := p.Lesson
		if l == nil {
			l = &models.Lesson{ContentType: "video"}
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO lessons (node_id, content_type, drive_file_id, embed_url, duration_minutes, body, attachments)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id, l.ContentType, l.DriveFileID, l.EmbedURL, l.DurationMinutes, l.Body, lessonAttachments(l.Attachments))
	case models.KindAssignment:
		a := p.Assignment
		if a == nil {
			a = &models.Assignment{}
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO assignments (node_id, instructions, time_limit_minutes, max_attempts, pass_score, shuffle_questions, due_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id, a.Instructions, a.TimeLimitMinutes, a.MaxAttempts, a.PassScore, a.ShuffleQuestions, a.DueAt)
	}
	if err != nil {
		return nil, translate(err, "tạo chi tiết nút")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetNode(ctx, id)
}

type UpdateNodeParams struct {
	Title       *string
	Description *string
	IsPublished *bool
	IsLocked    *bool

	Lesson     *models.Lesson
	Assignment *models.Assignment
}

func (s *Store) UpdateNode(ctx context.Context, id uuid.UUID, p UpdateNodeParams) (*models.Node, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var kind string
	if err := tx.QueryRow(ctx, `SELECT kind FROM nodes WHERE id = $1`, id).Scan(&kind); err != nil {
		return nil, translate(err, "đọc nút")
	}

	_, err = tx.Exec(ctx, `
		UPDATE nodes SET
			title        = COALESCE($2, title),
			description  = COALESCE($3, description),
			is_published = COALESCE($4, is_published),
			is_locked    = COALESCE($5, is_locked)
		WHERE id = $1`, id, p.Title, p.Description, p.IsPublished, p.IsLocked)
	if err != nil {
		return nil, translate(err, "cập nhật nút")
	}

	if p.Lesson != nil && kind == models.KindLesson {
		l := p.Lesson
		_, err = tx.Exec(ctx, `
			INSERT INTO lessons (node_id, content_type, drive_file_id, embed_url, duration_minutes, body, attachments)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (node_id) DO UPDATE SET
				content_type = EXCLUDED.content_type,
				drive_file_id = EXCLUDED.drive_file_id,
				embed_url = EXCLUDED.embed_url,
				duration_minutes = EXCLUDED.duration_minutes,
				body = EXCLUDED.body,
				attachments = EXCLUDED.attachments`,
			id, l.ContentType, l.DriveFileID, l.EmbedURL, l.DurationMinutes, l.Body, lessonAttachments(l.Attachments))
		if err != nil {
			return nil, translate(err, "cập nhật bài học")
		}
	}

	if p.Assignment != nil && kind == models.KindAssignment {
		a := p.Assignment
		_, err = tx.Exec(ctx, `
			INSERT INTO assignments (node_id, instructions, time_limit_minutes, max_attempts, pass_score, shuffle_questions, due_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (node_id) DO UPDATE SET
				instructions = EXCLUDED.instructions,
				time_limit_minutes = EXCLUDED.time_limit_minutes,
				max_attempts = EXCLUDED.max_attempts,
				pass_score = EXCLUDED.pass_score,
				shuffle_questions = EXCLUDED.shuffle_questions,
				due_at = EXCLUDED.due_at`,
			id, a.Instructions, a.TimeLimitMinutes, a.MaxAttempts, a.PassScore, a.ShuffleQuestions, a.DueAt)
		if err != nil {
			return nil, translate(err, "cập nhật bài tập")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return s.GetNode(ctx, id)
}

// MoveNode chuyển một nút sang nút cha khác và/hoặc đổi vị trí trong danh sách anh em.
func (s *Store) MoveNode(ctx context.Context, id uuid.UUID, newParent *uuid.UUID, position int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var programID uuid.UUID
	var oldParent *uuid.UUID
	err = tx.QueryRow(ctx, `SELECT program_id, parent_id FROM nodes WHERE id = $1`, id).Scan(&programID, &oldParent)
	if err != nil {
		return translate(err, "đọc nút")
	}

	if newParent != nil {
		// Không cho phép thả một nút vào chính nó hoặc vào nhánh con của nó.
		var isDescendant bool
		err = tx.QueryRow(ctx, `
			WITH RECURSIVE subtree AS (
				SELECT id FROM nodes WHERE id = $1
				UNION ALL
				SELECT n.id FROM nodes n JOIN subtree st ON n.parent_id = st.id
			)
			SELECT EXISTS (SELECT 1 FROM subtree WHERE id = $2)`, id, *newParent).Scan(&isDescendant)
		if err != nil {
			return translate(err, "kiểm tra vòng lặp cây")
		}
		if isDescendant {
			return Invalidf("Không thể chuyển một mục vào bên trong chính nó")
		}

		var parentKind string
		var parentProgram uuid.UUID
		err = tx.QueryRow(ctx, `SELECT kind, program_id FROM nodes WHERE id = $1`, *newParent).Scan(&parentKind, &parentProgram)
		if err != nil {
			return translate(err, "đọc nút cha")
		}
		if parentKind != models.KindFolder {
			return Invalidf("Chỉ có thể chuyển nội dung vào bên trong một thư mục")
		}
		if parentProgram != programID {
			return Invalidf("Không thể chuyển nội dung sang chương trình khác")
		}
	}

	if position < 0 {
		position = 0
	}

	// Nhân đôi thứ tự hiện tại để chèn nút vào đúng khe mong muốn, sau đó đánh số lại liên tục.
	if _, err = tx.Exec(ctx, `
		UPDATE nodes SET position = position * 2
		WHERE program_id = $1 AND parent_id IS NOT DISTINCT FROM $2`, programID, newParent); err != nil {
		return translate(err, "sắp xếp lại thứ tự")
	}
	if _, err = tx.Exec(ctx, `
		UPDATE nodes SET parent_id = $2, position = $3 WHERE id = $1`, id, newParent, position*2-1); err != nil {
		return translate(err, "chuyển nút")
	}

	if err := renumber(ctx, tx, programID, newParent); err != nil {
		return err
	}
	if !sameParent(oldParent, newParent) {
		if err := renumber(ctx, tx, programID, oldParent); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// renumber đánh lại position thành 0,1,2,... cho một nhóm anh em.
func renumber(ctx context.Context, tx pgx.Tx, programID uuid.UUID, parent *uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		WITH ordered AS (
			SELECT id, (row_number() OVER (ORDER BY position, created_at) - 1)::int AS rn
			FROM nodes
			WHERE program_id = $1 AND parent_id IS NOT DISTINCT FROM $2
		)
		UPDATE nodes n SET position = ordered.rn
		FROM ordered WHERE n.id = ordered.id AND n.position <> ordered.rn`, programID, parent)
	return translate(err, "đánh số lại thứ tự")
}

func sameParent(a, b *uuid.UUID) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func (s *Store) DeleteNode(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM nodes WHERE id = $1`, id)
	if err != nil {
		return translate(err, "xoá nút")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// NodeProgramID trả về chương trình chứa nút, dùng để kiểm tra quyền truy cập.
func (s *Store) NodeProgramID(ctx context.Context, nodeID uuid.UUID) (uuid.UUID, error) {
	var programID uuid.UUID
	err := s.pool.QueryRow(ctx, `SELECT program_id FROM nodes WHERE id = $1`, nodeID).Scan(&programID)
	if err != nil {
		return uuid.Nil, translate(err, "đọc chương trình của nút")
	}
	return programID, nil
}
