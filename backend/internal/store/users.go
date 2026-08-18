package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/manhnv/elearning/backend/internal/models"
)

const userColumns = `
	u.id, u.email, u.full_name, u.avatar_url, u.role, u.is_active, u.must_change_password,
	(u.password_hash IS NOT NULL) AS has_password, (u.google_sub IS NOT NULL) AS has_google,
	u.last_login_at, u.created_at`

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(row userScanner) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &u.AvatarURL, &u.Role, &u.IsActive,
		&u.MustChangePassword, &u.HasPassword, &u.HasGoogle, &u.LastLoginAt, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUserParams là dữ liệu admin nhập khi cấp tài khoản mới.
type CreateUserParams struct {
	Email        string
	FullName     string
	Role         string
	PasswordHash string
	GoogleSub    string
	AvatarURL    string
	// Buộc đổi mật khẩu ở lần đăng nhập đầu tiên.
	MustChangePassword bool
}

func (s *Store) CreateUser(ctx context.Context, p CreateUserParams) (*models.User, error) {
	var passwordHash, googleSub *string
	if p.PasswordHash != "" {
		passwordHash = &p.PasswordHash
	}
	if p.GoogleSub != "" {
		googleSub = &p.GoogleSub
	}

	row := s.pool.QueryRow(ctx, `
		INSERT INTO users (email, full_name, role, password_hash, google_sub, avatar_url, must_change_password)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+strings.ReplaceAll(userColumns, "u.", "")+`
	`, strings.ToLower(strings.TrimSpace(p.Email)), p.FullName, p.Role, passwordHash, googleSub, p.AvatarURL, p.MustChangePassword)

	u, err := scanUser(row)
	if err != nil {
		return nil, translate(err, "tạo người dùng")
	}
	return u, nil
}

func (s *Store) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users u WHERE u.id = $1`, id)
	u, err := scanUser(row)
	if err != nil {
		return nil, translate(err, "đọc người dùng")
	}
	return u, nil
}

// Credentials là bản ghi người dùng kèm hash mật khẩu, chỉ dùng cho luồng đăng nhập.
type Credentials struct {
	models.User
	PasswordHash string
}

func (s *Store) GetCredentialsByEmail(ctx context.Context, email string) (*Credentials, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+userColumns+`, COALESCE(u.password_hash, '')
		FROM users u WHERE lower(u.email) = lower($1)`, strings.TrimSpace(email))

	var c Credentials
	err := row.Scan(&c.ID, &c.Email, &c.FullName, &c.AvatarURL, &c.Role, &c.IsActive,
		&c.MustChangePassword, &c.HasPassword, &c.HasGoogle, &c.LastLoginAt, &c.CreatedAt, &c.PasswordHash)
	if err != nil {
		return nil, translate(err, "đọc thông tin đăng nhập")
	}
	return &c, nil
}

func (s *Store) GetCredentialsByID(ctx context.Context, id uuid.UUID) (*Credentials, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+userColumns+`, COALESCE(u.password_hash, '')
		FROM users u WHERE u.id = $1`, id)

	var c Credentials
	err := row.Scan(&c.ID, &c.Email, &c.FullName, &c.AvatarURL, &c.Role, &c.IsActive,
		&c.MustChangePassword, &c.HasPassword, &c.HasGoogle, &c.LastLoginAt, &c.CreatedAt, &c.PasswordHash)
	if err != nil {
		return nil, translate(err, "đọc thông tin đăng nhập")
	}
	return &c, nil
}

func (s *Store) GetUserByGoogleSub(ctx context.Context, sub string) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users u WHERE u.google_sub = $1`, sub)
	u, err := scanUser(row)
	if err != nil {
		return nil, translate(err, "đọc người dùng theo Google")
	}
	return u, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+userColumns+` FROM users u WHERE lower(u.email) = lower($1)`, strings.TrimSpace(email))
	u, err := scanUser(row)
	if err != nil {
		return nil, translate(err, "đọc người dùng theo email")
	}
	return u, nil
}

type ListUsersFilter struct {
	Search string
	Role   string
	Limit  int
	Offset int
}

func (s *Store) ListUsers(ctx context.Context, f ListUsersFilter) ([]*models.User, int, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	search := "%" + strings.ToLower(strings.TrimSpace(f.Search)) + "%"

	var total int
	err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM users u
		WHERE ($1 = '' OR u.role = $1)
		  AND ($2 = '%%' OR lower(u.email) LIKE $2 OR lower(u.full_name) LIKE $2)
	`, f.Role, search).Scan(&total)
	if err != nil {
		return nil, 0, translate(err, "đếm người dùng")
	}

	rows, err := s.pool.Query(ctx, `
		SELECT `+userColumns+` FROM users u
		WHERE ($1 = '' OR u.role = $1)
		  AND ($2 = '%%' OR lower(u.email) LIKE $2 OR lower(u.full_name) LIKE $2)
		ORDER BY u.created_at DESC
		LIMIT $3 OFFSET $4
	`, f.Role, search, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, translate(err, "liệt kê người dùng")
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, translate(err, "đọc người dùng")
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

type UpdateUserParams struct {
	FullName *string
	Role     *string
	IsActive *bool
}

func (s *Store) UpdateUser(ctx context.Context, id uuid.UUID, p UpdateUserParams) (*models.User, error) {
	row := s.pool.QueryRow(ctx, `
		UPDATE users u SET
			full_name = COALESCE($2, u.full_name),
			role      = COALESCE($3, u.role),
			is_active = COALESCE($4, u.is_active)
		WHERE u.id = $1
		RETURNING `+userColumns+`
	`, id, p.FullName, p.Role, p.IsActive)

	u, err := scanUser(row)
	if err != nil {
		return nil, translate(err, "cập nhật người dùng")
	}
	return u, nil
}

// SetPassword đặt lại mật khẩu; mustChange=true buộc người dùng đổi ở lần đăng nhập kế tiếp.
func (s *Store) SetPassword(ctx context.Context, id uuid.UUID, hash string, mustChange bool) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE users SET password_hash = $2, must_change_password = $3 WHERE id = $1`, id, hash, mustChange)
	if err != nil {
		return translate(err, "đặt mật khẩu")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// LinkGoogle gắn tài khoản Google vào một người dùng đã có sẵn trong hệ thống.
func (s *Store) LinkGoogle(ctx context.Context, id uuid.UUID, sub, avatarURL, fullName string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE users SET
			google_sub = $2,
			avatar_url = CASE WHEN avatar_url = '' THEN $3 ELSE avatar_url END,
			full_name  = CASE WHEN full_name = '' THEN $4 ELSE full_name END
		WHERE id = $1`, id, sub, avatarURL, fullName)
	return translate(err, "liên kết tài khoản Google")
}

func (s *Store) TouchLastLogin(ctx context.Context, id uuid.UUID) {
	// Không chặn đăng nhập nếu cập nhật mốc thời gian thất bại.
	_, _ = s.pool.Exec(ctx, `UPDATE users SET last_login_at = now() WHERE id = $1`, id)
}

func (s *Store) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return translate(err, "xoá người dùng")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CountActiveAdmins dùng để chặn thao tác vô hiệu hoá / hạ quyền admin cuối cùng.
func (s *Store) CountActiveAdmins(ctx context.Context, excluding uuid.UUID) (int, error) {
	var n int
	err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM users WHERE role = 'admin' AND is_active AND id <> $1`, excluding).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("đếm admin: %w", err)
	}
	return n, nil
}

// EnsureSeedAdmin tạo tài khoản quản trị đầu tiên nếu hệ thống chưa có admin nào.
func (s *Store) EnsureSeedAdmin(ctx context.Context, email, fullName, passwordHash string) (bool, error) {
	var exists bool
	if err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE role = 'admin')`).Scan(&exists); err != nil {
		return false, fmt.Errorf("kiểm tra admin: %w", err)
	}
	if exists {
		return false, nil
	}
	_, err := s.CreateUser(ctx, CreateUserParams{
		Email:        email,
		FullName:     fullName,
		Role:         models.RoleAdmin,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
