-- =============================================================
-- 0004: Vai trò Giám sát, slug cho URL, cấu hình hệ thống lưu trong DB
-- =============================================================

-- Tìm mọi ràng buộc CHECK nhắc tới cột role rồi xoá, bất kể tên hay cú pháp nội bộ.
-- Postgres viết lại "role IN (...)" thành "role = ANY (ARRAY[...])" nên không dò
-- theo chuỗi "IN" được — dò theo tên cột trong pg_constraint.conkey thay vào đó.
DO $$
DECLARE
    c record;
BEGIN
    FOR c IN
        SELECT con.conname
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_attribute att
             ON att.attrelid = rel.oid AND att.attnum = ANY (con.conkey)
        WHERE rel.relname = 'users' AND con.contype = 'c' AND att.attname = 'role'
    LOOP
        EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', c.conname);
    END LOOP;
END $$;

ALTER TABLE users
    ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'trainer', 'supervisor', 'student'));

-- -------------------------------------------------------------
-- Slug cho URL: chương trình dùng chính mã code; bài học/bài tập
-- sinh riêng từ tiêu đề. Việc sinh slug (bỏ dấu, chuẩn hoá) nằm ở
-- tầng ứng dụng — cột ở đây chỉ lưu kết quả và đảm bảo tính duy nhất.
-- -------------------------------------------------------------
ALTER TABLE programs ADD COLUMN slug text NOT NULL DEFAULT '';
UPDATE programs SET slug = lower(regexp_replace(code, '[^a-zA-Z0-9]+', '-', 'g'));
CREATE UNIQUE INDEX programs_slug_key ON programs (slug);

ALTER TABLE nodes ADD COLUMN slug text NOT NULL DEFAULT '';
-- Backfill đơn giản cho dữ liệu demo hiện có: chuẩn hoá ASCII, thêm số thứ tự
-- khi trùng trong cùng chương trình. Dữ liệu thật về sau luôn được sinh bởi Go.
WITH base AS (
    SELECT id, program_id,
           lower(regexp_replace(regexp_replace(title, '[^a-zA-Z0-9]+', '-', 'g'), '^-+|-+$', '', 'g')) AS s
    FROM nodes
),
numbered AS (
    SELECT id, program_id,
           CASE WHEN s = '' THEN 'muc' ELSE s END AS s,
           row_number() OVER (PARTITION BY program_id,
               CASE WHEN s = '' THEN 'muc' ELSE s END ORDER BY id) AS rn
    FROM base
)
UPDATE nodes n
SET slug = CASE WHEN numbered.rn = 1 THEN numbered.s ELSE numbered.s || '-' || numbered.rn END
FROM numbered
WHERE n.id = numbered.id;

CREATE UNIQUE INDEX nodes_slug_key ON nodes (program_id, slug);

-- -------------------------------------------------------------
-- Cấu hình hệ thống lưu trong DB (ví dụ: Google OAuth), để admin
-- chỉnh qua giao diện thay vì phải sửa file .env và khởi động lại.
-- -------------------------------------------------------------
CREATE TABLE app_settings (
    key        text PRIMARY KEY,
    value      text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    updated_by uuid REFERENCES users (id) ON DELETE SET NULL
);
