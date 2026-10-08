-- =============================================================
-- 0009: Lớp chức năng AI (EduFlow)
--   - Bổ sung 2 dạng câu hỏi còn thiếu: đúng/sai, điền khuyết
--   - Mức độ nhận thức, đáp án gợi ý và tiêu chí chấm cho tự luận
--   - Điểm/nhận xét do AI gợi ý cho câu tự luận (giáo viên xác nhận)
--   - Giáo án do AI sinh, nhật ký tác vụ AI, hội thoại trợ lý AI
-- =============================================================

-- -------------------------------------------------------------
-- Bốn dạng câu hỏi theo Bảng 3.1
--   single_choice : trắc nghiệm nhiều lựa chọn (1 đáp án đúng)
--   true_false    : mỗi option là một phát biểu, is_correct = giá trị Đúng/Sai
--   fill_blank    : mỗi option là một đáp án được chấp nhận
--   essay         : tự luận, chấm theo tiêu chí
-- multi_choice giữ lại để không làm hỏng dữ liệu đã có.
-- -------------------------------------------------------------
ALTER TABLE questions DROP CONSTRAINT IF EXISTS questions_type_check;
ALTER TABLE questions
    ADD CONSTRAINT questions_type_check
    CHECK (type IN ('single_choice', 'multi_choice', 'true_false', 'fill_blank', 'essay'));

-- Mức độ nhận thức: Nhận biết / Thông hiểu / Vận dụng (rỗng = chưa phân loại).
ALTER TABLE questions ADD COLUMN IF NOT EXISTS level text NOT NULL DEFAULT '';
-- Đáp án gợi ý và tiêu chí chấm, dùng cho câu tự luận.
ALTER TABLE questions ADD COLUMN IF NOT EXISTS sample_answer text NOT NULL DEFAULT '';
ALTER TABLE questions ADD COLUMN IF NOT EXISTS rubric text NOT NULL DEFAULT '';

-- -------------------------------------------------------------
-- Điểm tự luận do AI gợi ý: chỉ là gợi ý, điểm chính thức vẫn
-- nằm ở cột score và do giáo viên xác nhận.
-- -------------------------------------------------------------
ALTER TABLE submission_answers ADD COLUMN IF NOT EXISTS ai_score numeric(6, 2);
ALTER TABLE submission_answers ADD COLUMN IF NOT EXISTS ai_comment text NOT NULL DEFAULT '';

-- -------------------------------------------------------------
-- Giáo án do AI sinh, giáo viên chỉnh sửa và lưu lại
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS lesson_plans (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id         uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    program_id       uuid        REFERENCES programs (id) ON DELETE SET NULL,
    subject          text        NOT NULL DEFAULT '',
    grade            text        NOT NULL DEFAULT '',
    topic            text        NOT NULL DEFAULT '',
    objectives       text        NOT NULL DEFAULT '',
    duration_minutes integer     NOT NULL DEFAULT 45,
    title            text        NOT NULL DEFAULT '',
    content          jsonb       NOT NULL DEFAULT '{}'::jsonb,
    created_at       timestamptz NOT NULL DEFAULT now(),
    updated_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS lesson_plans_owner_idx ON lesson_plans (owner_id, created_at DESC);
DROP TRIGGER IF EXISTS lesson_plans_set_updated_at ON lesson_plans;
CREATE TRIGGER lesson_plans_set_updated_at BEFORE UPDATE ON lesson_plans
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- -------------------------------------------------------------
-- Nhật ký tác vụ AI: hiện ở "tác vụ AI gần đây" và dùng để theo dõi
-- mức độ sử dụng Gemini API.
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ai_tasks (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kind       text        NOT NULL,
    title      text        NOT NULL DEFAULT '',
    status     text        NOT NULL DEFAULT 'ok' CHECK (status IN ('ok', 'error')),
    detail     text        NOT NULL DEFAULT '',
    attempts   integer     NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ai_tasks_user_idx ON ai_tasks (user_id, created_at DESC);

-- -------------------------------------------------------------
-- Hội thoại với trợ lý AI. Lịch sử được gửi kèm mỗi lượt hỏi
-- để AI hiểu ngữ cảnh (mục 3.4.4).
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ai_conversations (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      text        NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ai_conversations_user_idx ON ai_conversations (user_id, updated_at DESC);
DROP TRIGGER IF EXISTS ai_conversations_set_updated_at ON ai_conversations;
CREATE TRIGGER ai_conversations_set_updated_at BEFORE UPDATE ON ai_conversations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS ai_messages (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id uuid        NOT NULL REFERENCES ai_conversations (id) ON DELETE CASCADE,
    role            text        NOT NULL CHECK (role IN ('user', 'model')),
    content         text        NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ai_messages_conversation_idx ON ai_messages (conversation_id, created_at);
