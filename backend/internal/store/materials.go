package store

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

// SaveMaterialParams là dữ liệu một tài liệu dùng chung được lưu từ giao diện.
type SaveMaterialParams struct {
	Title       string
	Description string
	Category    string
	URL         string
	DriveFileID string
	Kind        string
	IsPublished bool
	CreatedBy   uuid.UUID
}

const materialColumns = `m.id, m.title, m.description, m.category, m.url,
	m.drive_file_id, m.kind, m.is_published, m.created_by, m.created_at, m.updated_at,
	COALESCE(u.full_name, '')`

func scanMaterial(row rowScanner) (*models.Material, error) {
	var m models.Material
	err := row.Scan(&m.ID, &m.Title, &m.Description, &m.Category, &m.URL,
		&m.DriveFileID, &m.Kind, &m.IsPublished, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		&m.CreatedByName)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// ListMaterialsFilter lọc kho tài liệu theo từ khoá và phân loại.
type ListMaterialsFilter struct {
	Search   string
	Category string
	// false: chỉ trả tài liệu đã xuất bản (dành cho học viên).
	IncludeUnpublished bool
}

func (s *Store) ListMaterials(ctx context.Context, f ListMaterialsFilter) ([]*models.Material, error) {
	search := strings.TrimSpace(f.Search)
	rows, err := s.pool.Query(ctx, `
		SELECT `+materialColumns+`
		FROM materials m
		LEFT JOIN users u ON u.id = m.created_by
		WHERE ($1 OR m.is_published = true)
		  AND ($2 = '' OR m.title ILIKE '%' || $2 || '%' OR m.description ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR m.category = $3)
		ORDER BY m.created_at DESC
		LIMIT 500`, f.IncludeUnpublished, search, strings.TrimSpace(f.Category))
	if err != nil {
		return nil, translate(err, "liệt kê tài liệu")
	}
	defer rows.Close()

	out := []*models.Material{}
	for rows.Next() {
		m, err := scanMaterial(rows)
		if err != nil {
			return nil, translate(err, "đọc tài liệu")
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// MaterialCategories trả về danh sách phân loại đang được dùng, để giao diện
// dựng ô lọc mà không cần một bảng danh mục riêng.
func (s *Store) MaterialCategories(ctx context.Context) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT DISTINCT category FROM materials
		WHERE category <> '' ORDER BY category`)
	if err != nil {
		return nil, translate(err, "đọc phân loại tài liệu")
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, translate(err, "đọc phân loại tài liệu")
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) GetMaterial(ctx context.Context, id uuid.UUID) (*models.Material, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+materialColumns+`
		FROM materials m LEFT JOIN users u ON u.id = m.created_by
		WHERE m.id = $1`, id)
	m, err := scanMaterial(row)
	if err != nil {
		return nil, translate(err, "đọc tài liệu")
	}
	return m, nil
}

func (s *Store) CreateMaterial(ctx context.Context, p SaveMaterialParams) (*models.Material, error) {
	var id uuid.UUID
	err := s.pool.QueryRow(ctx, `
		INSERT INTO materials (title, description, category, url, drive_file_id, kind, is_published, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		p.Title, p.Description, p.Category, p.URL, p.DriveFileID, p.Kind, p.IsPublished, p.CreatedBy).Scan(&id)
	if err != nil {
		return nil, translate(err, "tạo tài liệu")
	}
	return s.GetMaterial(ctx, id)
}

func (s *Store) UpdateMaterial(ctx context.Context, id uuid.UUID, p SaveMaterialParams) (*models.Material, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE materials
		SET title = $2, description = $3, category = $4, url = $5,
		    drive_file_id = $6, kind = $7, is_published = $8
		WHERE id = $1`,
		id, p.Title, p.Description, p.Category, p.URL, p.DriveFileID, p.Kind, p.IsPublished)
	if err != nil {
		return nil, translate(err, "cập nhật tài liệu")
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return s.GetMaterial(ctx, id)
}

func (s *Store) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM materials WHERE id = $1`, id)
	if err != nil {
		return translate(err, "xoá tài liệu")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
