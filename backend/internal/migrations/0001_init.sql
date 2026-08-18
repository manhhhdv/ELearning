-- =============================================================
-- 0001_init: Lược đồ khởi tạo hệ thống đào tạo trực tuyến
-- =============================================================

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- -------------------------------------------------------------
-- Người dùng: admin / giảng viên / học viên dùng chung một bảng
-- -------------------------------------------------------------
CREATE TABLE users (
    id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email                text        NOT NULL,
    full_name            text        NOT NULL DEFAULT '',
    password_hash        text,                       -- NULL khi tài khoản chỉ đăng nhập bằng Google
    google_sub           text UNIQUE,                -- subject ID do Google cấp
    avatar_url           text        NOT NULL DEFAULT '',
    role                 text        NOT NULL DEFAULT 'student'
                                     CHECK (role IN ('admin', 'trainer', 'student')),
    is_active            boolean     NOT NULL DEFAULT true,
    must_change_password boolean     NOT NULL DEFAULT false,
    last_login_at        timestamptz,
    created_at           timestamptz NOT NULL DEFAULT now(),
    updated_at           timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_key ON users (lower(email));
CREATE TRIGGER users_set_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- -------------------------------------------------------------
-- Chương trình đào tạo (gốc của cây nội dung)
-- -------------------------------------------------------------
CREATE TABLE programs (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code        text        NOT NULL,
    title       text        NOT NULL,
    description text        NOT NULL DEFAULT '',
    cover_url   text        NOT NULL DEFAULT '',
    status      text        NOT NULL DEFAULT 'draft'
                            CHECK (status IN ('draft', 'published', 'archived')),
    created_by  uuid REFERENCES users (id) ON DELETE SET NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX programs_code_key ON programs (lower(code));
CREATE TRIGGER programs_set_updated_at BEFORE UPDATE ON programs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- -------------------------------------------------------------
-- Nút trên cây nội dung: thư mục / bài học / bài tập
-- Dùng danh sách kề (parent_id) + position để sắp thứ tự anh em
-- -------------------------------------------------------------
CREATE TABLE nodes (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id   uuid        NOT NULL REFERENCES programs (id) ON DELETE CASCADE,
    parent_id    uuid REFERENCES nodes (id) ON DELETE CASCADE,
    kind         text        NOT NULL CHECK (kind IN ('folder', 'lesson', 'assignment')),
    title        text        NOT NULL,
    description  text        NOT NULL DEFAULT '',
    position     integer     NOT NULL DEFAULT 0,
    is_published boolean     NOT NULL DEFAULT true,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX nodes_program_idx ON nodes (program_id);
CREATE INDEX nodes_parent_idx ON nodes (parent_id);
CREATE INDEX nodes_sibling_idx ON nodes (program_id, parent_id, position);
CREATE TRIGGER nodes_set_updated_at BEFORE UPDATE ON nodes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- -------------------------------------------------------------
-- Bài học: nội dung nhúng iframe từ Google Drive
-- -------------------------------------------------------------
CREATE TABLE lessons (
    node_id          uuid PRIMARY KEY REFERENCES nodes (id) ON DELETE CASCADE,
    content_type     text    NOT NULL DEFAULT 'video'
                             CHECK (content_type IN ('video', 'slide', 'document', 'pdf', 'link')),
    drive_file_id    text    NOT NULL DEFAULT '',  -- ID file trên Google Drive
    embed_url        text    NOT NULL DEFAULT '',  -- URL nhúng đã chuẩn hoá
    duration_minutes integer NOT NULL DEFAULT 0,
    body             text    NOT NULL DEFAULT ''   -- ghi chú / mô tả thêm cho bài học
);

-- -------------------------------------------------------------
-- Bài tập: cấu hình chung của một bài (MCQ và/hoặc tự luận)
-- -------------------------------------------------------------
CREATE TABLE assignments (
    node_id            uuid PRIMARY KEY REFERENCES nodes (id) ON DELETE CASCADE,
    instructions       text          NOT NULL DEFAULT '',
    time_limit_minutes integer       NOT NULL DEFAULT 0,   -- 0 = không giới hạn
    max_attempts       integer       NOT NULL DEFAULT 0,   -- 0 = không giới hạn
    pass_score         numeric(6, 2) NOT NULL DEFAULT 0,
    shuffle_questions  boolean       NOT NULL DEFAULT false,
    due_at             timestamptz
);

-- -------------------------------------------------------------
-- Câu hỏi và đáp án
-- -------------------------------------------------------------
CREATE TABLE questions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id uuid          NOT NULL REFERENCES assignments (node_id) ON DELETE CASCADE,
    type          text          NOT NULL CHECK (type IN ('single_choice', 'multi_choice', 'essay')),
    prompt        text          NOT NULL,
    points        numeric(6, 2) NOT NULL DEFAULT 1,
    position      integer       NOT NULL DEFAULT 0,
    explanation   text          NOT NULL DEFAULT '',
    created_at    timestamptz   NOT NULL DEFAULT now(),
    updated_at    timestamptz   NOT NULL DEFAULT now()
);
CREATE INDEX questions_assignment_idx ON questions (assignment_id, position);
CREATE TRIGGER questions_set_updated_at BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE question_options (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id uuid    NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    content     text    NOT NULL,
    is_correct  boolean NOT NULL DEFAULT false,
    position    integer NOT NULL DEFAULT 0
);
CREATE INDEX question_options_question_idx ON question_options (question_id, position);

-- -------------------------------------------------------------
-- Ghi danh học viên / giảng viên vào chương trình
-- -------------------------------------------------------------
CREATE TABLE enrollments (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id  uuid        NOT NULL REFERENCES programs (id) ON DELETE CASCADE,
    user_id     uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role        text        NOT NULL DEFAULT 'student' CHECK (role IN ('student', 'trainer')),
    enrolled_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (program_id, user_id)
);
CREATE INDEX enrollments_user_idx ON enrollments (user_id);

-- -------------------------------------------------------------
-- Tiến độ học từng bài
-- -------------------------------------------------------------
CREATE TABLE lesson_progress (
    user_id      uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    node_id      uuid        NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,
    completed_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, node_id)
);

-- -------------------------------------------------------------
-- Bài nộp và câu trả lời
-- -------------------------------------------------------------
CREATE TABLE submissions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id uuid          NOT NULL REFERENCES assignments (node_id) ON DELETE CASCADE,
    user_id       uuid          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    attempt_no    integer       NOT NULL DEFAULT 1,
    status        text          NOT NULL DEFAULT 'submitted'
                                CHECK (status IN ('submitted', 'graded')),
    auto_score    numeric(6, 2) NOT NULL DEFAULT 0,   -- điểm phần trắc nghiệm chấm tự động
    manual_score  numeric(6, 2),                      -- điểm phần tự luận do giảng viên chấm
    max_score     numeric(6, 2) NOT NULL DEFAULT 0,
    feedback      text          NOT NULL DEFAULT '',
    graded_by     uuid REFERENCES users (id) ON DELETE SET NULL,
    graded_at     timestamptz,
    submitted_at  timestamptz   NOT NULL DEFAULT now(),
    UNIQUE (assignment_id, user_id, attempt_no)
);
CREATE INDEX submissions_user_idx ON submissions (user_id);
CREATE INDEX submissions_assignment_idx ON submissions (assignment_id, submitted_at DESC);

CREATE TABLE submission_answers (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id       uuid          NOT NULL REFERENCES submissions (id) ON DELETE CASCADE,
    question_id         uuid          NOT NULL REFERENCES questions (id) ON DELETE CASCADE,
    selected_option_ids uuid[]        NOT NULL DEFAULT '{}',
    essay_text          text          NOT NULL DEFAULT '',
    is_correct          boolean,
    score               numeric(6, 2) NOT NULL DEFAULT 0,
    comment             text          NOT NULL DEFAULT '',
    UNIQUE (submission_id, question_id)
);
