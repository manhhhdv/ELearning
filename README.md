# Tập Huấn — Hệ thống đào tạo trực tuyến

Nền tảng đào tạo nội bộ: quản trị viên dựng chương trình đào tạo dưới dạng **cây nội dung**, nhúng
video và slide trực tiếp từ **Google Drive**, kèm **bài tập trắc nghiệm và tự luận**. Học viên và
giảng viên dùng chung một hệ thống, đăng nhập bằng **Google** hoặc **tài khoản do admin cấp**.

- **Backend:** Go 1.25 · chi · pgx · JWT · bcrypt · OAuth2
- **Frontend:** React 19 · Vite · TypeScript
- **Database:** PostgreSQL 16

---

## Chạy thử

```bash
make dev
```

Một lệnh dựng cả ba: khởi động Postgres trong Docker và chờ nó sẵn sàng, cài phụ thuộc nếu thiếu,
rồi chạy song song API và giao diện. Log của hai bên được gắn nhãn `[api]` / `[web]`.
Bấm **Ctrl+C** để dừng cả hai (Postgres vẫn chạy tiếp, dữ liệu giữ nguyên).

Nếu muốn mỗi thứ một cửa sổ terminal riêng: `make db`, `make backend`, `make frontend`.

Xem toàn bộ lệnh có sẵn:

```bash
make help
```

Mở http://localhost:3006 và đăng nhập bằng tài khoản quản trị khởi tạo:

| Email | Mật khẩu |
|---|---|
| `admin@elearning.local` | `Admin@12345` |

> Tài khoản này chỉ được tạo tự động khi database chưa có admin nào. **Đổi mật khẩu ngay sau lần
> đăng nhập đầu tiên**, hoặc đặt `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` trong `backend/.env`
> trước khi chạy lần đầu.

### Cổng sử dụng

| Thành phần | Cổng | Ghi chú |
|---|---|---|
| Giao diện React | 3006 | Vite proxy `/api` sang backend |
| API Go | 8082 | |
| PostgreSQL | 5433 | Trong Docker, không đụng Postgres cài sẵn ở 5432 |

Vite chạy với `strictPort` nên khi cổng bận sẽ báo lỗi thay vì âm thầm nhảy sang cổng khác
(đổi cổng ngầm sẽ làm hỏng CORS và redirect OAuth). `make dev` kiểm tra hai cổng trước khi chạy và báo rõ tiến trình nào đang chiếm, thay vì để server
chết lặng lẽ giữa đống log.

Muốn đổi cổng: sửa `backend/.env` (`PORT`, `FRONTEND_URL`, `ALLOW_ORIGINS`),
`frontend/vite.config.ts` (`server.port`, `server.proxy`) và `Makefile` (`BACKEND_PORT`,
`FRONTEND_PORT`).

Muốn cấu hình Google Login hoặc đổi mật khẩu admin khởi tạo thì tạo file `.env` trước:

```bash
cp backend/.env.example backend/.env
```

Không có `.env` thì backend vẫn chạy bằng giá trị mặc định.

---

## Mô hình phân quyền

Bốn vai trò dùng chung một bảng `users`, một lần đăng nhập:

| Vai trò | Quyền |
|---|---|
| **Quản trị viên** | Toàn quyền: cấp tài khoản, mọi chương trình, chấm bài, cấu hình hệ thống |
| **Giảng viên** | Sửa nội dung và chấm bài ở chương trình mình tạo hoặc được ghi danh làm giảng viên |
| **Giám sát** | Xem toàn bộ chương trình (kể cả bản nháp), cây nội dung, đáp án đúng, danh sách ghi danh, bài nộp/điểm và bảng điều khiển — **không sửa/xoá/ghi danh/chấm bài được ở đâu cả** |
| **Học viên** | Chỉ học và làm bài ở chương trình đã xuất bản mà mình được ghi danh |

Quản trị viên, giảng viên và Giám sát đều vào được khu vực quản lý trên thanh điều hướng; quyền xem
so với quyền sửa được tách riêng ở cả tầng API (`CanManage` khác `CanAudit`) lẫn giao diện (các trang
quản lý tự ẩn nút sửa/xoá/tạo mới và khoá form khi đăng nhập bằng vai trò Giám sát).

Hệ thống chặn hạ quyền, khoá hoặc xoá quản trị viên cuối cùng, và không cho tự xoá tài khoản đang
đăng nhập.

---

## Cấu trúc cây nội dung

Mỗi chương trình là một cây gồm ba loại nút:

| Loại | Vai trò |
|---|---|
| **Thư mục** | Chương / phần — chỉ thư mục mới chứa được nút con |
| **Bài học** | Nội dung nhúng iframe: video, slide, tài liệu, PDF, liên kết ngoài |
| **Bài tập** | Tập câu hỏi trắc nghiệm một đáp án / nhiều đáp án / tự luận |

Kéo-thả để sắp xếp: thả **vào giữa** một thư mục để đưa vào trong, thả **sát mép trên/dưới** để chèn
trước/sau, thả ra **vùng trống** để đưa lên cấp gốc. Hệ thống chặn việc kéo một mục vào chính nhánh
con của nó.

Nút chưa xuất bản (bỏ tick *Hiển thị với học viên*) hiện chữ nghiêng kèm nhãn "ẩn", học viên không
thấy nút đó lẫn toàn bộ nhánh con.

### URL thân thiện (slug)

Khoá học và bài học dùng **slug** thay UUID trên đường dẫn: `/hoc/attp-2026/bai-1-gioi-thieu` thay vì
`/hoc/6edb17bf-.../0a9880b6-...`. Slug khoá học lấy từ mã (`ATTP-2026` → `attp-2026`), slug bài học
sinh từ tiêu đề lúc tạo (bỏ dấu tiếng Việt, chuẩn hoá về chữ thường và gạch ngang), tự thêm hậu tố
`-2`, `-3`… nếu trùng trong cùng chương trình. Slug **chốt lúc tạo** — sửa tiêu đề sau đó không đổi
slug, để link đã chia sẻ không bị hỏng. Đổi mã chương trình thì slug đổi theo (URL luôn khớp mã hiện
tại thay vì giữ mã cũ đã bỏ).

### Xem trước như học viên

Nút **Xem trước** ở trang soạn nội dung (mở tab mới) hiển thị đúng giao diện học viên — video/slide
nhúng, đề bài tập kèm đáp án đúng — nhưng dùng dữ liệu của người quản lý nên thấy được cả nội dung
chưa xuất bản, và **không ghi lại bất kỳ hành động nào**: không có nút đánh dấu hoàn thành, không nộp
bài thật, không tính tiến độ.

---

## Tự ghi danh và khám phá khoá học

Bật *Cho học viên tự ghi danh* khi tạo hoặc sửa chương trình (tab *Cài đặt*) để khoá xuất hiện ở mục
**Khám phá** trên thanh điều hướng — học viên tự bấm *Đăng ký học* mà không cần đợi admin thêm vào,
và tự rời được nếu muốn. Chương trình do admin/giảng viên chủ động ghi danh thủ công (tắt tuỳ chọn
này) thì học viên không tự rời được — đúng logic đối xứng với cách được thêm vào.

### Khoá học mặc định

Bật *Khoá học mặc định* (chỉ admin thấy và đặt được ở tab *Cài đặt*) để khoá **tự động hiện trong
"Khoá học của tôi" của mọi người dùng** ngay khi xuất bản — không cần ghi danh, không cần tự đăng ký.
Dùng cho nội dung bắt buộc như định hướng nhân viên mới hay quy tắc ứng xử. Thẻ khoá học có huy hiệu
**📌 Bắt buộc** để người dùng hiểu vì sao nó xuất hiện mà họ không hề đăng ký.

Kỹ thuật: không tạo dòng `enrollments` nào — quyền xem được cấp trực tiếp ở tầng kiểm tra quyền
(`programAccess`) cho bất kỳ ai chưa ghi danh nhưng chương trình đang là mặc định và đã xuất bản.
Tắt cờ đi thì học viên hết thấy ngay, không để lại enrollment rác cần dọn.

---

## Nhập câu hỏi hàng loạt

Ở trình soạn câu hỏi, nút **Nhập hàng loạt** nhận hai định dạng, xem trước danh sách phân tích được
trước khi nhập thật:

- **Soạn dạng văn bản**: mỗi câu là một khối cách nhau bằng dòng trống, `[C10] Nội dung [2đ]` cho mã
  và điểm, dòng `*` là phương án đúng, `-` là phương án sai, `>` là giải thích. Khối không có phương án
  nào thì thành câu tự luận.
- **Dán từ Excel / Google Sheets**: dán thẳng vùng đã copy (phân tách bằng Tab), thứ tự cột *Mã · Nội
  dung · Điểm · PA1-4 · Đáp án đúng (ghi `1,3` hoặc `A,C`) · Giải thích*.

Nhập hàng loạt là **một transaction** — sai một câu thì không câu nào được ghi, tránh nhập dở dang.
Câu hỏi có **mã cố định** (`C01`, `C02`…, hoặc tự đặt), không đổi khi kéo sắp xếp lại thứ tự — dùng để
đối chiếu qua các lần sửa và hiển thị trên bảng thống kê kết quả theo từng câu (tab *Kết quả học
viên* trong bài tập): tỉ lệ trả lời đúng, điểm trung bình, số lượt bỏ trống cho mỗi câu.

---

## Bảng điều khiển quản trị

Trang `/quan-tri` (mục *Bảng điều khiển* trong menu Quản lý) tổng hợp số liệu toàn hệ thống: số
chương trình theo trạng thái, số người dùng theo vai trò, tổng lượt ghi danh, tổng bài nộp và số đang
chờ chấm, chương trình nhiều học viên nhất, người dùng mới tạo gần đây. Admin và Giám sát đều xem
được.

---

## Đăng nhập Google — cấu hình qua giao diện

Vào **Quản lý → Đăng nhập Google** (chỉ admin) để bật/sửa Client ID, Client Secret, giới hạn domain
email và vai trò tự tạo tài khoản **mà không cần sửa file `.env` hay khởi động lại máy chủ**. Cấu hình
lưu trong bảng `app_settings`, ưu tiên hơn `.env` khi đã bật; tắt đi thì xoá cấu hình đã lưu và quay
lại dùng `.env` (nếu máy chủ có khai báo sẵn). Client Secret không bao giờ trả về nguyên văn sau khi
lưu — để trống ô này khi sửa các trường khác nghĩa là giữ nguyên secret đang có hiệu lực.

---

## Nhúng nội dung từ Google Drive

Dán **link chia sẻ** hoặc **ID file** vào ô *Link Google Drive*, hệ thống tự chuyển thành URL nhúng:

| Dán vào | Kết quả |
|---|---|
| `https://drive.google.com/file/d/ID/view?usp=sharing` | `https://drive.google.com/file/d/ID/preview` |
| `https://docs.google.com/presentation/d/ID/edit#slide=id.p` | `https://docs.google.com/presentation/d/ID/embed?…` |
| `https://drive.google.com/open?id=ID` | `https://drive.google.com/file/d/ID/preview` |
| `ID` (dán thẳng) | `https://drive.google.com/file/d/ID/preview` |
| Link ngoài Drive (YouTube, Vimeo…) | Giữ nguyên |

**Quan trọng:** đặt quyền chia sẻ file trên Drive thành *"Bất kỳ ai có đường liên kết"* (hoặc chia sẻ
cho toàn bộ domain tổ chức), nếu không học viên sẽ thấy khung nhúng trống.

---

## Nhận diện và giao diện

Logo **Tập Huấn** là chiếc mũ tốt nghiệp trên nền xanh bo góc, dùng chung cho thanh điều hướng,
trang đăng nhập, trình học và favicon. Nguồn hình nằm ở hai chỗ và phải sửa song song:

| Tệp | Dùng cho |
|---|---|
| `frontend/public/favicon.svg` | Favicon vector, nguồn để xuất các bản PNG |
| `frontend/src/components/Logo.tsx` | Logo hiển thị trong ứng dụng |

Các bản raster (`favicon.ico`, `favicon-16/32.png`, `apple-touch-icon.png`, `logo-512.png`) được
sinh từ file SVG, không sửa tay.

Toàn hệ thống dùng **một tông sáng duy nhất** (không đổi theo thiết lập sáng/tối của hệ điều hành),
điều hướng bằng **thanh ngang** ở đầu trang: logo, các mục chính, ô tìm kiếm khoá học, chuông báo số
bài đang chờ chấm, nút trợ giúp và menu tài khoản. Màn hình hẹp thu ô tìm kiếm về một nút bấm và gom
các mục điều hướng vào menu tài khoản.

Khu vực học viên đi theo bố cục quen thuộc của các nền tảng MOOC (Coursera):

- **Trang chính** (`/hoc`) gộp cả thống kê lẫn danh sách khoá: lời chào, ba thẻ *Đã ghi danh /
  Đã hoàn thành / Số bài đã nộp*, rồi lưới thẻ khoá học. Ô tìm kiếm trên thanh điều hướng lọc
  ngay tại trang này qua tham số `?q=`.
- **Thẻ khoá học**: ảnh bìa nếu chương trình có `coverUrl`, không thì dùng dải màu suy ra từ mã khoá
  (luôn ổn định qua các lần tải), kèm thanh tiến độ và nhãn *Đã học x/y bài / Đã hoàn thành*.
- **Trình học một khoá** chiếm trọn màn hình, không dùng sidebar chung:
  - Thanh trên: nút quay lại, tên khoá, mã khoá, ảnh đại diện.
  - Cột trái: mục lục theo chương, gập/mở được, mỗi mục có vòng tròn đánh dấu hoàn thành và dòng
    mô tả *Video · 15 phút* / *Bài tập · 3 câu hỏi*. Tiến độ tổng hiển thị theo phần trăm.
  - Cột phải: tiêu đề bài, khung nhúng, ghi chú, và thanh dính đáy có **Bài trước / Đánh dấu hoàn
    thành / Bài tiếp theo**. Bấm *Bài tiếp theo* ở một bài học sẽ đánh dấu hoàn thành rồi chuyển
    luôn sang bài kế.
  - Mã bài nằm trên URL (`/hoc/:maKhoa/:maBai`) nên tải lại trang vẫn giữ đúng vị trí đang học;
    mở khoá học mà không chỉ đích danh bài thì nhảy thẳng tới bài đầu tiên chưa hoàn thành.
- **Làm bài tập**: màn hình mở đầu hiển thị điểm cao nhất, số lượt đã dùng, điểm đạt và hạn nộp;
  khi vào làm thì mỗi câu là một thẻ riêng, thanh đáy đếm *Đã trả lời x/y câu* và nút nộp bài.
- **Xem lại bài nộp**: vòng tròn điểm số theo phần trăm (xanh khi đã chấm, hổ phách khi chờ chấm),
  tách rõ điểm trắc nghiệm và điểm tự luận, từng câu hiện lựa chọn của học viên — đáp án đúng và
  giải thích chỉ hiện sau khi bài được chấm xong.

Giao diện tự đổi sáng/tối theo thiết lập hệ điều hành và dùng được trên điện thoại (mục lục khoá học
thu gọn sau nút *Nội dung*). Font Source Sans 3 được đóng gói sẵn trong ứng dụng, không gọi ra ngoài
Internet.

Khu vực quản trị giữ nguyên giao diện gọn dạng bảng, tối ưu cho thao tác nhập liệu.

---

## Bài tập và chấm điểm

- **Trắc nghiệm** được chấm tự động ngay khi nộp, theo nguyên tắc **đúng trọn vẹn**: tập đáp án chọn
  phải trùng khít tập đáp án đúng mới được điểm.
- **Tự luận** được lưu lại chờ giảng viên vào điểm ở màn hình *Chấm bài*.
- Bài chỉ có trắc nghiệm được đánh dấu **đã chấm** ngay lập tức.
- Học viên **không** thấy đáp án đúng và lời giải cho tới khi bài được chấm xong.
- Giới hạn số lượt làm bài (`0` = không giới hạn) được kiểm tra ở phía máy chủ.

### Giới hạn thời gian làm bài

Đặt *Thời gian làm bài* khác `0` thì bài tập trở thành bài có tính giờ:

- Học viên bấm **Bắt đầu làm bài**, máy chủ mở một *phiên làm bài* và ghi mốc hết giờ. Giao diện
  hiển thị đồng hồ đếm ngược, chuyển vàng khi còn dưới 5 phút và đỏ khi còn dưới 1 phút.
- Đồng hồ **không** reset khi tải lại trang hay bấm lại nút bắt đầu — mốc bắt đầu chỉ ghi một lần.
  Tải lại trang giữa chừng sẽ quay về đúng màn hình làm bài.
- Hết giờ, giao diện tự nộp phần đã làm. Nếu học viên đóng trình duyệt, máy chủ vẫn từ chối mọi bài
  nộp muộn (có khoảng trễ 45 giây cho độ trễ mạng) — đồng hồ **được cưỡng chế ở máy chủ**, không
  phải chỉ trang trí ở trình duyệt.
- Phiên quá hạn mà không nộp thì không tính vào số lượt làm bài: học viên bắt đầu lại được.

### Trộn câu hỏi

Bật *Trộn thứ tự câu hỏi và phương án* trong cấu hình bài tập:

- Mỗi lượt làm bài có một thứ tự riêng, sinh từ ID phiên nên **giữ nguyên khi tải lại trang**.
- Trộn cả thứ tự câu hỏi lẫn thứ tự phương án trả lời; câu tự luận không bị ảnh hưởng.
- Chấm điểm dựa trên ID câu hỏi và ID phương án nên thứ tự hiển thị không ảnh hưởng kết quả.

---

## Cấu trúc mã nguồn

```
backend/
  cmd/server/          Điểm khởi động, seed admin, graceful shutdown
  internal/
    api/               Router chi và toàn bộ handler HTTP (kể cả settings.go, dashboard.go)
    auth/              JWT, bcrypt, OAuth2 Google, middleware phân quyền
    config/            Nạp cấu hình mặc định từ .env (override được qua app_settings trong DB)
    database/          Kết nối pgx pool và bộ chạy migration
    migrations/        File .sql nhúng vào binary
    models/            Kiểu dữ liệu dùng chung
    store/             Toàn bộ truy vấn Postgres (dashboard.go, settings.go, attempts.go…)
    util/              Chuyển link Google Drive sang URL nhúng, sinh slug tiếng Việt
frontend/
  src/
    api/               Client gọi REST và kiểu TypeScript
    components/        Cây kéo-thả, trình soạn nội dung/câu hỏi, trình làm bài, nhập hàng loạt
    pages/             Các màn hình, gồm AdminDashboardPage, GoogleSettingsPage, PreviewPage
    styles/            CSS (tự đổi màu theo sáng/tối của hệ điều hành)
```

Migration chạy tự động lúc khởi động: mỗi file `.sql` trong `backend/internal/migrations/` được áp
dụng đúng một lần trong một transaction, ghi nhận ở bảng `schema_migrations`.

---

## Lược đồ dữ liệu

| Bảng | Nội dung |
|---|---|
| `users` | Tài khoản dùng chung bốn vai trò, hash bcrypt và/hoặc `google_sub` |
| `programs` | Chương trình đào tạo: mã, `slug`, trạng thái, `allow_self_enroll`, `is_default_course` |
| `nodes` | Cây nội dung: `parent_id` + `position` + `slug` (duy nhất trong phạm vi chương trình) |
| `lessons` | Chi tiết bài học 1-1 với nút, giữ `drive_file_id` và `embed_url` |
| `assignments` | Cấu hình bài tập 1-1 với nút |
| `questions`, `question_options` | Câu hỏi (có `code` cố định) và phương án trả lời |
| `enrollments` | Ghi danh học viên/giảng viên vào chương trình |
| `lesson_progress` | Đánh dấu hoàn thành từng bài |
| `submissions`, `submission_answers` | Bài nộp, điểm tự động và điểm chấm tay |
| `attempt_sessions` | Phiên làm bài đang mở — mốc bắt đầu/hết giờ cho bài có tính thời gian |
| `app_settings` | Cấu hình sửa được qua giao diện admin (hiện dùng cho đăng nhập Google) |

---

## Trước khi đưa lên production

> Hướng dẫn triển khai đầy đủ bằng Docker (backend) + Cloudflare Worker (frontend): xem
> [DEPLOY.md](DEPLOY.md).

- [ ] Đặt `JWT_SECRET` bằng chuỗi ngẫu nhiên (`openssl rand -base64 48`) — backend từ chối khởi động
      với `APP_ENV=production` nếu để trống.
- [ ] Đổi mật khẩu tài khoản `SEED_ADMIN_EMAIL`.
- [ ] Đổi mật khẩu Postgres trong `docker-compose.yml` và `DATABASE_URL`.
- [ ] Đặt `APP_ENV=production`, `ALLOW_ORIGINS` và `FRONTEND_URL` theo domain thật.
- [ ] Chạy sau TLS; cập nhật `GOOGLE_REDIRECT_URL` sang `https://` và khai báo lại trên Google Cloud.
      Sau khi lên production, đổi Client ID/Secret hoặc bật/tắt Google login làm được ngay qua
      **Quản lý → Đăng nhập Google** — không cần sửa `.env` hay khởi động lại.

---

## Kiểm thử

```bash
cd backend && go test ./...
```

Phạm vi đang có: chuyển link Google Drive sang URL nhúng, sinh slug tiếng Việt, tính ổn định của
thuật toán trộn đề, chống giả mạo `state` trong luồng OAuth Google (sửa nonce, kéo dài hạn, đổi chữ
ký, ký bằng khoá khác, state hết hạn), và bộ lọc domain email được phép đăng nhập.
