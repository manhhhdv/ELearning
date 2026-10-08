-- =============================================================
-- 0010: Tự đăng ký tài khoản và yêu cầu đặt lại mật khẩu
--   - Người dùng tự đăng ký ở trang đăng nhập (vai trò Học sinh)
--   - Quên mật khẩu: gửi yêu cầu để admin cấp lại, không cần SMTP
-- =============================================================

-- -------------------------------------------------------------
-- Yêu cầu đặt lại mật khẩu.
--
-- Hệ thống chưa gửi được email nên không dùng token đặt lại; thay vào đó
-- người dùng gửi yêu cầu, admin thấy ở màn hình Quản lý người dùng và cấp
-- mật khẩu mới.
--
-- user_id có thể NULL: người dùng gõ nhầm email không tồn tại thì vẫn ghi
-- nhận để admin biết mà hỗ trợ, đồng thời trang đăng nhập không lộ email nào
-- có trong hệ thống.
-- -------------------------------------------------------------
CREATE TABLE IF NOT EXISTS password_reset_requests (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email       text        NOT NULL,
    user_id     uuid        REFERENCES users (id) ON DELETE CASCADE,
    note        text        NOT NULL DEFAULT '',
    status      text        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'done', 'rejected')),
    handled_by  uuid        REFERENCES users (id) ON DELETE SET NULL,
    handled_at  timestamptz,
    created_at  timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS password_reset_requests_status_idx
    ON password_reset_requests (status, created_at DESC);

-- Mỗi email chỉ có một yêu cầu đang chờ: bấm nhiều lần không tạo ra hàng loạt
-- dòng rác cho admin phải dọn.
CREATE UNIQUE INDEX IF NOT EXISTS password_reset_requests_pending_email_idx
    ON password_reset_requests (lower(email)) WHERE status = 'pending';
