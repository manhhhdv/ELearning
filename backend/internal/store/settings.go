package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Các khoá cấu hình lưu trong bảng app_settings, đọc/ghi qua giao diện admin
// thay vì phải sửa file .env và khởi động lại máy chủ.
const (
	SettingGoogleClientID          = "google_client_id"
	SettingGoogleClientSecret      = "google_client_secret"
	SettingGoogleAllowedDomains    = "google_allowed_domains"
	SettingGoogleAutoProvisionRole = "google_auto_provision_role"
)

// GetSetting đọc một giá trị cấu hình; ok=false nếu chưa được đặt trong DB
// (khi đó tầng gọi nên dùng giá trị mặc định từ biến môi trường).
func (s *Store) GetSetting(ctx context.Context, key string) (value string, ok bool, err error) {
	rawErr := s.pool.QueryRow(ctx, `SELECT value FROM app_settings WHERE key = $1`, key).Scan(&value)
	if rawErr == nil {
		return value, true, nil
	}
	if te := translate(rawErr, "đọc cấu hình"); errors.Is(te, ErrNotFound) {
		return "", false, nil
	} else {
		return "", false, te
	}
}

// GetSettings đọc nhiều khoá cùng lúc, trả về map chỉ chứa các khoá đã có giá trị.
func (s *Store) GetSettings(ctx context.Context, keys []string) (map[string]string, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, value FROM app_settings WHERE key = ANY($1)`, keys)
	if err != nil {
		return nil, translate(err, "đọc cấu hình")
	}
	defer rows.Close()

	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, translate(err, "đọc cấu hình")
		}
		out[k] = v
	}
	return out, rows.Err()
}

// SetSetting ghi đè hoặc tạo mới một khoá cấu hình. Giá trị rỗng vẫn được lưu
// (khác với "chưa đặt") — dùng SetSettings với xoá tường minh nếu cần bỏ hẳn.
func (s *Store) SetSetting(ctx context.Context, key, value string, updatedBy uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO app_settings (key, value, updated_by, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = now()`,
		key, value, updatedBy)
	return translate(err, "lưu cấu hình")
}

// DeleteSetting xoá một khoá khỏi DB để quay về dùng giá trị mặc định từ .env.
func (s *Store) DeleteSetting(ctx context.Context, key string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM app_settings WHERE key = $1`, key)
	return translate(err, "xoá cấu hình")
}
