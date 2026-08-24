-- =============================================================
-- 0008: Bài học dạng "Tài liệu tải về" — đính kèm nhiều tài liệu,
--       mỗi tài liệu có tên hiển thị và link tải riêng.
-- =============================================================

ALTER TABLE lessons
    ADD COLUMN IF NOT EXISTS attachments jsonb NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE lessons
    DROP CONSTRAINT IF EXISTS lessons_content_type_check;

ALTER TABLE lessons
    ADD CONSTRAINT lessons_content_type_check
    CHECK (content_type IN ('video', 'slide', 'document', 'pdf', 'link', 'richtext', 'materials'));
