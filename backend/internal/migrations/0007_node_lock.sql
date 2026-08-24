-- =============================================================
-- 0007: Cho phép quản trị viên khoá/mở từng nội dung
-- =============================================================

ALTER TABLE nodes
    ADD COLUMN is_locked boolean NOT NULL DEFAULT false;
