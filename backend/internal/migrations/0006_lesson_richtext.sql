-- =============================================================
-- 0006: Bài học tự soạn — nội dung viết thẳng trong hệ thống
--       (markdown + công thức LaTeX + media nhúng) thay vì
--       nhúng file Google Drive.
-- =============================================================

ALTER TABLE lessons
    DROP CONSTRAINT IF EXISTS lessons_content_type_check;

ALTER TABLE lessons
    ADD CONSTRAINT lessons_content_type_check
    CHECK (content_type IN ('video', 'slide', 'document', 'pdf', 'link', 'richtext'));
