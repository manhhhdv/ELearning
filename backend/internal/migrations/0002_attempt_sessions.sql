-- =============================================================
-- 0002: Phiên làm bài — mốc bắt đầu để cưỡng chế thời gian làm bài
--       và giữ thứ tự trộn câu hỏi ổn định khi học viên tải lại trang.
-- =============================================================

CREATE TABLE attempt_sessions (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    assignment_id uuid        NOT NULL REFERENCES assignments (node_id) ON DELETE CASCADE,
    user_id       uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    started_at    timestamptz NOT NULL DEFAULT now(),
    -- NULL khi bài tập không giới hạn thời gian.
    expires_at    timestamptz,
    submitted_at  timestamptz
);

-- Mỗi học viên chỉ có một phiên đang mở cho mỗi bài tập.
CREATE UNIQUE INDEX attempt_sessions_open_idx
    ON attempt_sessions (assignment_id, user_id)
    WHERE submitted_at IS NULL;

CREATE INDEX attempt_sessions_user_idx ON attempt_sessions (user_id);

-- Gắn bài nộp với phiên đã tạo ra nó, phục vụ đối chiếu về sau.
ALTER TABLE submissions
    ADD COLUMN session_id uuid REFERENCES attempt_sessions (id) ON DELETE SET NULL;
