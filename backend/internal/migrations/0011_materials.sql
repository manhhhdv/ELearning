-- =============================================================
-- 0011: Kho tài liệu dùng chung (mục 3.2.3 và 4.4 của báo cáo)
--
-- Khác với tài liệu đính kèm trong từng bài học, đây là kho dùng chung
-- toàn hệ thống: quản trị viên và giảng viên đăng tài liệu, mọi người
-- đăng nhập đều xem được.
--
-- Hệ thống không lưu file nhị phân mà lưu đường dẫn (Google Drive hoặc
-- link ngoài), thống nhất với cách bài học đang nhúng nội dung.
-- =============================================================

CREATE TABLE IF NOT EXISTS materials (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title       text        NOT NULL,
    description text        NOT NULL DEFAULT '',
    -- Phân loại để lọc nhanh: môn học, khối lớp hoặc chủ đề tuỳ đơn vị đặt.
    category    text        NOT NULL DEFAULT '',
    -- Link tải hoặc xem tài liệu.
    url         text        NOT NULL,
    -- ID file Drive tách sẵn từ url (rỗng nếu là link ngoài), dùng để nhúng xem trước.
    drive_file_id text      NOT NULL DEFAULT '',
    -- Loại nội dung để chọn biểu tượng và cách mở: pdf, slide, document, video, link.
    kind        text        NOT NULL DEFAULT 'link'
                CHECK (kind IN ('pdf', 'slide', 'document', 'video', 'link')),
    -- Tắt cờ này thì tài liệu chỉ người quản lý thấy, dùng khi đang soạn dở.
    is_published boolean    NOT NULL DEFAULT true,
    created_by  uuid        REFERENCES users (id) ON DELETE SET NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS materials_created_idx ON materials (created_at DESC);
CREATE INDEX IF NOT EXISTS materials_category_idx ON materials (category);

DROP TRIGGER IF EXISTS materials_set_updated_at ON materials;
CREATE TRIGGER materials_set_updated_at BEFORE UPDATE ON materials
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
