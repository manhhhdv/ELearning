-- =============================================================
-- 0005: Khoá học mặc định — tự động hiện trong "Khoá học của tôi"
--       của mọi người dùng mà không cần ghi danh thủ công.
-- =============================================================

ALTER TABLE programs
    ADD COLUMN is_default_course boolean NOT NULL DEFAULT false;

-- Chỉ có ý nghĩa khi đã xuất bản, nên chỉ cần lọc nhanh những chương trình mặc định.
CREATE INDEX programs_default_course_idx ON programs (is_default_course) WHERE is_default_course;
