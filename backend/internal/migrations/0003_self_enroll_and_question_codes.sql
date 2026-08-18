-- =============================================================
-- 0003: Cho phép học viên tự ghi danh, và gán mã cố định cho câu hỏi
-- =============================================================

-- Bật/tắt theo từng chương trình. Mặc định tắt: chỉ admin/giảng viên ghi danh.
ALTER TABLE programs
    ADD COLUMN allow_self_enroll boolean NOT NULL DEFAULT false;

-- Mã ngắn do người soạn nhìn thấy, dùng để đối chiếu câu hỏi giữa các lần sửa,
-- khi nhập hàng loạt và trên bảng thống kê kết quả.
ALTER TABLE questions
    ADD COLUMN code text NOT NULL DEFAULT '';

-- Gán mã cho các câu hỏi đã có: C01, C02… theo đúng thứ tự hiện tại.
UPDATE questions q
SET code = 'C' || lpad(seq::text, 2, '0')
FROM (
    SELECT id, row_number() OVER (PARTITION BY assignment_id ORDER BY position, created_at) AS seq
    FROM questions
) AS numbered
WHERE q.id = numbered.id;

-- Mã không được trùng trong cùng một bài tập; chuỗi rỗng thì bỏ qua ràng buộc.
CREATE UNIQUE INDEX questions_code_key
    ON questions (assignment_id, lower(code))
    WHERE code <> '';
