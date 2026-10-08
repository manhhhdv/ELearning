# Đào Tạo — Hệ thống đào tạo trực tuyến

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
| **Bài học** | Nội dung nhúng iframe (video, slide, tài liệu, PDF, liên kết ngoài) hoặc **bài đọc tự soạn** viết thẳng trong hệ thống |
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

**Giải thích là bắt buộc** với mọi câu hỏi (kể cả tự luận, dùng làm gợi ý chấm bài) — thiếu giải thích
bị chặn ngay khi lưu một câu đơn lẻ lẫn khi nhập hàng loạt, cả ở form lẫn API.

Nhập hàng loạt là **một transaction** — sai một câu thì không câu nào được ghi, tránh nhập dở dang.
Câu hỏi có **mã cố định** (`C01`, `C02`…, hoặc tự đặt), không đổi khi kéo sắp xếp lại thứ tự — dùng để
đối chiếu qua các lần sửa và hiển thị trên bảng thống kê kết quả theo từng câu (tab *Kết quả học
viên* trong bài tập): tỉ lệ trả lời đúng, điểm trung bình, số lượt bỏ trống cho mỗi câu.

---

## Nhập tài khoản hàng loạt

Nút **Nhập từ file** ở trang *Quản lý → Người dùng* (chỉ admin) cấp nhiều tài khoản một lần: tải lên
file `.xlsx`/`.csv` hoặc dán thẳng vùng bảng đã copy từ Excel / Google Sheets. Thứ tự cột *Email · Họ
và tên · Vai trò · Mật khẩu*; dòng tiêu đề được bỏ qua tự động, tối đa 500 dòng mỗi lần.

- **Vai trò** nhận cả tiếng Việt lẫn tiếng Anh (`Học viên`, `Giảng viên`, `Giám sát`, `Quản trị viên`,
  hay `student`, `trainer`…); để trống thì mặc định là Học viên.
- **Mật khẩu** để trống thì hệ thống sinh ngẫu nhiên cho riêng dòng đó; tự đặt thì phải đủ 8 ký tự
  gồm cả chữ và số. Mọi tài khoản mới đều bị buộc đổi mật khẩu ở lần đăng nhập đầu tiên.
- Email sai định dạng, vai trò lạ hay email lặp trong file bị bắt lỗi **trước khi nhập**, kèm số dòng.

Khác với nhập câu hỏi, mỗi dòng ở đây được xử lý **độc lập** chứ không gói trong một transaction —
danh sách nhân sự thường có sẵn vài người đã có tài khoản, huỷ cả lô chỉ vì một dòng trùng sẽ rất bất
tiện. Email đã tồn tại được **bỏ qua** (không ghi đè dữ liệu cũ), và bảng kết quả cuối cùng liệt kê
từng dòng: đã tạo / bỏ qua / lỗi kèm lý do. Mật khẩu chỉ hiển thị **đúng một lần** ở bảng đó — nút
*Tải file mật khẩu (.csv)* xuất danh sách tài khoản vừa tạo để gửi cho người dùng.

---

## Bảng điều khiển quản trị

Trang `/quan-tri` (mục *Bảng điều khiển* trong menu Quản lý) tổng hợp số liệu toàn hệ thống: số
chương trình theo trạng thái, số người dùng theo vai trò, tổng lượt ghi danh, tổng bài nộp và số đang
chờ chấm, chương trình nhiều học viên nhất, người dùng mới tạo gần đây. Admin và Giám sát đều xem
được.

---

## Kho tài liệu dùng chung

Mục **Tài liệu** trên thanh điều hướng là kho dùng chung cho toàn hệ thống, **không thuộc khoá học
nào** — khác với tài liệu đính kèm trong từng bài học.

- Ai đăng nhập cũng tra cứu được; **giảng viên và quản trị viên** thêm, sửa, xoá.
- Hệ thống **lưu đường dẫn** (Google Drive hoặc link ngoài) chứ không lưu bản sao tệp, thống nhất
  với cách bài học đang nhúng nội dung. Nhớ đặt quyền chia sẻ cho người học xem được.
- Mỗi tài liệu có **phân loại** tự do (môn học, khối lớp…) dùng để lọc, và cờ **hiện cho người học**
  để soạn dở mà chưa công bố.
- Lọc theo từ khoá và phân loại ngay trên trang.

---

## Trang tổng quan — "không gian làm việc AI"

Trang chính sau khi đăng nhập (`/tong-quan`) dùng chung một bố cục cho mọi vai trò:

- **Ô nhập yêu cầu cho AI** kèm gợi ý nhanh — hỏi trợ lý ngay tại trang chủ, không cần mở khung chat.
  Gợi ý khác nhau theo vai trò (học viên hỏi bài, giảng viên hỏi ý tưởng dạy học).
- **Lối tắt** tới các chức năng chính của vai trò đó.
- **Bài cần làm**: bài tập được giao còn hạn trong 7 ngày tới hoặc **đã quá hạn mà chưa nộp**. Bài
  đã nộp tự biến mất khỏi danh sách.
- **Tác vụ AI gần đây** — 6 dòng mới nhất, xem đầy đủ ở mục riêng.
- **Tiến độ học tập**: số bài đã hoàn thành trên tổng số bài của từng khoá, kèm điểm trung bình
  (mục 3.5).

**Bố cục**: toàn bộ nội dung nằm trong container rộng tối đa 1360px canh giữa (class `.page-body` —
trang nào quên bọc thì nội dung sẽ tràn sát mép màn hình). Trên màn rộng chia hai cột: cột chính giữ
*Bài cần làm* và *Tiến độ học tập*, cột phụ giữ *Tác vụ AI gần đây*. Cột phụ rỗng — chẳng hạn khi
chưa bật chức năng AI — thì cột chính tự chiếm trọn bề ngang.

Người dạy và quản trị viên thấy thêm:

- **Số bài đang chờ chấm**, hiện ngay trên lối tắt *Chấm bài*.
- **Cảnh báo học tập**: học viên có điểm trung bình dưới 50% hoặc có bài quá hạn chưa nộp. Đây là
  **luật cố định**, không phải phân tích bằng AI — việc dùng AI để cảnh báo sớm là hướng phát triển
  tiếp theo. Giảng viên chỉ thấy khoá mình phụ trách, admin thấy toàn hệ thống.

---

## Đăng nhập, đăng ký và quên mật khẩu

Trang đăng nhập có: **chọn vai trò**, đăng nhập bằng email/mật khẩu, **ghi nhớ đăng nhập**,
**quên mật khẩu**, đăng nhập Google và **đăng ký tài khoản mới**.

**Chọn vai trò** chỉ là lựa chọn ở giao diện. Vai trò thật nằm ở tài khoản và do admin cấp — chọn
nhầm tab vẫn đăng nhập được, hệ thống mở đúng khu vực tương ứng và báo cho người dùng biết.

**Ghi nhớ đăng nhập** quyết định nơi lưu token: bật thì `localStorage` (còn sau khi đóng trình
duyệt), tắt thì `sessionStorage` (mất khi đóng tab).

**Đăng ký tài khoản** mặc định **tắt**, bật ở **Quản lý → Đăng ký tài khoản**. Tài khoản tự đăng ký
luôn có vai trò **Học viên** và dùng được ngay; quyền Giảng viên hoặc Quản trị viên vẫn phải do
admin nâng cấp. Nên đặt **giới hạn domain email** của đơn vị — để trống nghĩa là bất kỳ ai biết địa
chỉ trang web cũng tạo được tài khoản.

**Quên mật khẩu** không gửi email (hệ thống chưa cấu hình SMTP). Người dùng gửi yêu cầu kèm lời
nhắn; admin thấy danh sách ngay đầu mục **Quản lý → Người dùng**, bấm *Tìm tài khoản* rồi *Đặt lại
mật khẩu* — yêu cầu **tự đóng** sau khi cấp mật khẩu mới. Phản hồi cho người gửi luôn giống nhau dù
email có tồn tại hay không, để trang đăng nhập không thành công cụ dò tài khoản.

---

## Cấu hình sửa được qua giao diện admin

Các cấu hình nhạy cảm không cần sửa file `.env`: **đăng nhập Google**, **nhà cung cấp AI** và
**cho phép tự đăng ký**. Tất cả theo cùng một nguyên tắc — giá trị lưu trong bảng `app_settings` được ưu tiên, thiếu thì rơi
về `.env`, thiếu nốt thì chức năng tự tắt. Tắt ở giao diện sẽ xoá bản ghi trong DB để quay lại dùng
`.env`. Bí mật (Client Secret, khoá API) **không bao giờ được trả về trình duyệt**.

---

## Đăng nhập Google — cấu hình qua giao diện

Vào **Quản lý → Đăng nhập Google** (chỉ admin) để bật/sửa Client ID, Client Secret, giới hạn domain
email và vai trò tự tạo tài khoản **mà không cần sửa file `.env` hay khởi động lại máy chủ**. Cấu hình
lưu trong bảng `app_settings`, ưu tiên hơn `.env` khi đã bật; tắt đi thì xoá cấu hình đã lưu và quay
lại dùng `.env` (nếu máy chủ có khai báo sẵn). Client Secret không bao giờ trả về nguyên văn sau khi
lưu — để trống ô này khi sửa các trường khác nghĩa là giữ nguyên secret đang có hiệu lực.

---

## Bài đọc tự soạn (markdown, LaTeX, media)

Chọn *Loại nội dung* → **Bài đọc tự soạn** khi muốn viết nội dung ngay trong hệ thống thay vì nhúng
file Drive. Ô soạn thảo có thanh công cụ, phím tắt (`Ctrl/Cmd + B` đậm, `+ I` nghiêng, `+ K` liên kết)
và khung **xem trước trực tiếp** cạnh bên; nội dung lưu xuống là markdown thuần nên dễ sao chép,
so sánh phiên bản và nhập từ file.

| Gõ | Kết quả |
|---|---|
| `**đậm**`, `*nghiêng*`, `~~gạch ngang~~`, `` `mã` `` | Định dạng chữ |
| `# Tiêu đề`, `- gạch đầu dòng`, `1. đánh số`, `- [ ] việc cần làm` | Cấu trúc, danh sách, ô tích |
| `> trích dẫn`, ```` ```khối mã``` ````, bảng kiểu markdown | Trích dẫn, khối mã, bảng (tự cuộn ngang) |
| `$a^2+b^2=c^2$`, `$$…$$` (hoặc `\(…\)`, `\[…\]`) | Công thức LaTeX dựng bằng KaTeX |
| `![mô tả](link)` | Ảnh, video, audio, hoặc khung nhúng YouTube / Vimeo / Google Drive |
| `> [!NOTE]` mở đầu trích dẫn | Hộp nhấn mạnh (NOTE, TIP, IMPORTANT, WARNING, CAUTION) |

Loại media nhận theo đuôi file và tên miền: `.mp4/.webm/.mov` thành trình phát video, `.mp3/.wav/.m4a`
thành trình phát audio, link YouTube/Vimeo/Drive thành iframe 16:9, còn lại là ảnh. Chú thích trong
ngoặc kép — `![mô tả](link "chú thích")` — hiện thành dòng caption dưới media.

Ghi chú của các loại bài học khác (ô *Ghi chú cho học viên*) dùng chung trình soạn thảo này, nên
markdown và công thức cũng hiển thị được trong hộp *Ghi chú bài học* của trang học.

Ký hiệu `$` đứng một mình (VD: `giá $150`) không bị hiểu nhầm thành công thức, và `$` bên trong khối
mã giữ nguyên. HTML kết xuất luôn đi qua bộ lọc DOMPurify: chỉ iframe từ YouTube, Vimeo và Google
Drive được giữ lại, link ngoài tự thêm `target="_blank" rel="noopener noreferrer"`.

Nhập cấu trúc từ Excel: ghi `richtext` (hoặc `tự soạn`, `bài đọc`) ở cột *loại nội dung*, để trống cột
link, rồi soạn nội dung sau khi nhập.

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

## Bài học "Tài liệu tải về"

Chọn *Loại nội dung* → **Tài liệu tải về** khi bài học chỉ là bộ tài liệu cho học viên tải về thay vì
một file nhúng. Mỗi dòng trong ô *Tài liệu tải về* gồm **tên hiển thị** và **link tải**; dùng nút ↑ ↓
để sắp thứ tự, nút thùng rác để xoá, tối đa 50 tài liệu cho một bài.

- Link nhận URL `http(s)` bất kỳ, hoặc link chia sẻ / ID file Google Drive (dán ID thì hệ thống dựng
  thành `https://drive.google.com/file/d/ID/view`).
- Bỏ trống tên thì học viên thấy chính đường link; dòng để trống cả hai ô sẽ bị bỏ qua khi lưu.
- File Drive vẫn phải đặt quyền chia sẻ *"Bất kỳ ai có đường liên kết"*, nếu không học viên bấm vào
  sẽ bị chặn.
- Ô *Ghi chú cho học viên* vẫn dùng được để hướng dẫn thêm, hiện dưới danh sách tài liệu.

Trang học hiển thị danh sách thành các thẻ bấm được, mở ở tab mới. Nhập cấu trúc từ Excel: ghi
`materials` (hoặc `tài liệu tải về`, `đính kèm`) ở cột *loại nội dung*, rồi thêm tài liệu sau khi nhập.

---

## Nhận diện và giao diện

Logo **Đào Tạo** là chiếc mũ tốt nghiệp trên nền xanh bo góc, dùng chung cho thanh điều hướng,
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

Bố cục đi theo các nền tảng MOOC quen thuộc (Coursera):

- **Trang chính** (`/hoc`) gộp cả thống kê lẫn danh sách khoá: lời chào, ba thẻ *Đã ghi danh /
  Đã hoàn thành / Số bài đã nộp*, rồi lưới thẻ khoá học. Ô tìm kiếm trên thanh điều hướng lọc
  ngay tại trang này qua tham số `?q=`.
- **Thẻ khoá học**: ảnh bìa nếu chương trình có `coverUrl`, không thì dùng dải màu suy ra từ mã khoá
  (luôn ổn định qua các lần tải), tên khoá, mô tả ngắn hai dòng, thanh tiến độ và nhãn *Đã học x/y
  bài / Đã hoàn thành*.
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
- **Chấm bài tự luận**: mỗi câu hiện *Hướng dẫn chấm* gồm **đáp án gợi ý** và **tiêu chí chấm** do
  chính người soạn đề viết, ngay cạnh bài làm — không phải mở đề ở tab khác để đối chiếu. Ô cho điểm
  có ba nút chấm nhanh (*Tối đa / Một nửa / 0*), ô nhận xét nhiều dòng. Thanh dính đáy màn hình luôn
  hiện tổng điểm hiện tại và **còn mấy câu chưa cho điểm**, kèm nút lưu — bài dài không phải cuộn
  xuống tận cuối. Câu chưa chấm hiện *Chưa chấm · tối đa N điểm* thay vì *0/N điểm*.
- Sau khi chấm xong, học viên xem được **đáp án gợi ý** để tự đối chiếu, còn **tiêu chí chấm là công
  cụ nội bộ của giáo viên nên luôn bị gỡ khỏi dữ liệu gửi xuống trình duyệt**.

Giao diện dùng được trên điện thoại (mục lục khoá học thu gọn sau nút *Nội dung*). Font Source
Sans 3 được đóng gói sẵn trong ứng dụng, không gọi ra ngoài Internet.

### Hệ thiết kế dùng chung

Trước đây khu học viên và khu quản trị chạy hai bộ style khác nhau, nên cùng một nút lại có màu và
cỡ chữ khác nhau tuỳ trang. Nay **cả ứng dụng dùng chung một hệ thiết kế**:

- Biến màu đặt ở `:root` trong `learner.css`; `app.css` ánh xạ các biến `--primary`, `--border`,
  `--text`… sang đúng bộ màu đó. Sửa màu ở một chỗ là cả hệ thống đổi theo.
- Màu nhấn `#0056d2`, chữ `#1f1f1f`, viền `#e3e8ee`. Màu chỉ dùng cho liên kết, nút và nhãn trạng
  thái — **tiêu đề luôn là chữ tối**, không tô màu.
- Chữ nền 15px, dòng 1.6; thang tiêu đề 28 / 21 / 17 / 15px.
- `.panel` là thẻ nội dung chuẩn: viền mảnh, bo 12px, đệm 22–24px. Thẻ chỉ chứa đúng một bảng thì
  bảng chạm mép thẻ; thẻ có thêm tiêu đề vẫn giữ đệm.
- **Mọi điều khiển cao đúng 40px** (32px với biến thể `-sm`), đặt ở biến
  `--control-h`. Ô nhập, ô chọn và nút đứng cùng một hàng thì thẳng mép trên dưới.

### Điều khiển biểu mẫu

| Thành phần | Cách dùng |
|---|---|
| Nút | `.btn` · `.btn-primary` (nền xanh) · `.btn-outline` (viền xanh) · `.btn-ghost` · `.btn-danger` · `.btn-sm` · `.btn-block` · `.btn-icon` |
| Ô nhập | Không cần class — `input`, `select`, `textarea` được định dạng sẵn |
| Ô chọn tệp | Component `FilePicker` |
| Trạng thái lỗi | Đặt `aria-invalid="true"`, kèm `<div className="field-error">` |
| Dấu bắt buộc | `<span className="req">*</span>` trong `<label>` |

Vài điểm đáng lưu ý:

- **Một hệ nút duy nhất.** `.cbtn*` của khu học viên trước đây là hệ riêng với kích thước và màu
  khác; nay chỉ còn là tên gọi khác trỏ về cùng định nghĩa với `.btn`, nên sửa một chỗ là đổi khắp
  nơi. `.cbtn` không kèm biến thể vẫn giữ dáng viền xanh như cũ.
- **Ô nhập chọn theo phép loại trừ** (`input:not([type=checkbox])…`) thay vì liệt kê từng `type`.
  Trước đây `type="search"` không nằm trong danh sách nên ô tìm kiếm hiện theo mặc định của trình
  duyệt, lệch hẳn với các ô khác.
- **Ô chọn** bị gỡ mũi tên mặc định của hệ điều hành và vẽ lại bằng SVG nhúng, nên trông giống nhau
  trên mọi trình duyệt.
- **Ô chọn tệp** dùng `FilePicker`: trình duyệt tự vẽ nút `input[type=file]` với chữ tiếng Anh
  ("Choose File", "No file chosen") mà CSS không đổi được, nên input thật bị ẩn và nhãn `<label>`
  đóng vai nút bấm — ra tiếng Việt và đúng kiểu nút chung.
- **Vòng focus** chỉ hiện khi đi bằng bàn phím (`:focus-visible`), bấm chuột không bị nhấp nháy.

Chỉ có **một nút tròn nổi** ở góc phải — trợ lý AI. Nút trợ giúp nằm trên thanh điều hướng, không
lặp lại thành nút nổi thứ hai.

---

## Bài tập và chấm điểm

Hệ thống hỗ trợ **bốn dạng câu hỏi**, mỗi dạng có cách chấm riêng:

| Dạng | Cấu trúc | Cách chấm |
|---|---|---|
| **Trắc nghiệm nhiều lựa chọn** | Câu dẫn và 4 phương án, một phương án đúng | Tự động: so khớp phương án đã chọn |
| **Đúng / Sai** | Nhiều phát biểu, mỗi phát biểu chọn Đúng hoặc Sai | Tự động theo từng phát biểu, **có điểm thành phần** |
| **Điền khuyết** | Câu có chỗ trống `___`, kèm danh sách đáp án được chấp nhận | Tự động sau khi chuẩn hoá câu trả lời |
| **Tự luận** | Câu hỏi mở, kèm đáp án gợi ý và tiêu chí chấm | AI gợi ý điểm, **giáo viên xác nhận** |

- Ba dạng có đáp án xác định được chấm **tự động ngay khi nộp**. Trắc nghiệm chấm theo nguyên tắc
  **đúng trọn vẹn**: tập đáp án chọn phải trùng khít tập đáp án đúng.
- **Đúng/Sai** chấm từng phát biểu: đúng 3/4 phát biểu được 3/4 số điểm của câu. Phát biểu học viên
  không đánh dấu được hiểu là chọn *Sai*.
- **Điền khuyết** chuẩn hoá câu trả lời trước khi so sánh: bỏ qua khác biệt hoa/thường, khoảng trắng
  thừa, dấu câu ở hai đầu, và chuẩn hoá Unicode (tiếng Việt gõ tổ hợp khớp với đáp án gõ dựng sẵn).
- **Tự luận** được lưu lại chờ giảng viên vào điểm ở màn hình *Chấm bài*.
- Mỗi câu có thể gắn **mức độ nhận thức**: Nhận biết / Thông hiểu / Vận dụng.
- Bài không có câu tự luận được đánh dấu **đã chấm** ngay lập tức.
- Học viên **không** thấy đáp án đúng và lời giải cho tới khi bài được chấm xong. Với câu điền
  khuyết, máy chủ **gỡ bỏ hẳn** danh sách đáp án khỏi dữ liệu gửi xuống trình duyệt, vì chính nội
  dung phương án là đáp án; câu tự luận cũng bị gỡ đáp án gợi ý và tiêu chí chấm.
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

## Trí tuệ nhân tạo

Mọi lời gọi AI đều đi **qua backend**. Khoá API chỉ nằm trên máy chủ, không bao giờ gửi xuống trình
duyệt. Lớp `internal/ai` (Service AI) chịu trách nhiệm dựng prompt, gọi Gemini, kiểm tra và chuẩn
hoá kết quả.

### Bật chức năng AI — cấu hình qua giao diện

Vào **Quản lý → Cấu hình AI** (chỉ admin). Tại đây khai báo một hoặc nhiều **kênh**, mỗi kênh gồm
nhà cung cấp, khoá API và model. Nút **Kiểm tra kết nối** gọi thử một yêu cầu rất ngắn bằng khoá vừa
gõ, để biết khoá và model có dùng được không **trước khi lưu** — thay vì để giáo viên gặp lỗi giữa
lúc đang soạn đề.

Nhà cung cấp được hỗ trợ: **Google Gemini**, **OpenAI**, **Anthropic Claude**, **Groq**,
**OpenRouter** và mục **Khác (tương thích OpenAI)** cho proxy nội bộ hoặc dịch vụ tự dựng — mục này
cần nhập thêm địa chỉ API.

**Xoay vòng khi lỗi.** Hệ thống gọi lần lượt từ kênh trên cùng xuống dưới theo thứ tự trong danh
sách. Kênh nào báo hết hạn mức (HTTP 429, hoặc thông báo lỗi có chữ *quota* / *overloaded*) sẽ bị
tạm ngưng 5 phút; lỗi tạm thời khác (HTTP 5xx, lỗi mạng, phản hồi rỗng) bị tạm ngưng 1 phút — và
lượt gọi tự chuyển sang kênh kế tiếp. Ngược lại, lỗi do cấu hình sai hoặc do nội dung (khoá sai,
model không tồn tại, bị chặn nội dung) trả về ngay, vì đổi kênh cũng không cứu được. Con trỏ xoay
vòng dịch một bước sau mỗi lần gọi nên tải được rải đều thay vì luôn đập vào kênh đầu.

Nên khai báo ít nhất hai kênh của **hai nhà cung cấp khác nhau**: giờ cao điểm — cả lớp cùng dùng
trợ lý AI — là lúc hạn mức của một khoá dễ chạm trần nhất. Kênh nào nghi hỏng có thể bật
**Tạm ngưng** để giữ lại cấu hình mà không gọi tới.

Giống cấu hình đăng nhập Google: khoá lưu trong bảng `app_settings`, **có hiệu lực ngay**, không cần
sửa `.env` hay khởi động lại máy chủ. Client gọi AI được dựng lại theo từng request từ cấu hình đang
có hiệu lực.

Khoá API **không bao giờ được gửi xuống trình duyệt**. Màn hình cấu hình chỉ hiển thị 4 ký tự cuối
(`••••1234`) để admin nhận ra mình đang dùng khoá nào.

Cách còn lại là đặt sẵn trong `backend/.env` — tiện cho môi trường phát triển hoặc khi triển khai
tự động. Đường này chỉ tạo được **một kênh Gemini**; muốn xoay vòng nhiều nhà cung cấp thì phải khai
báo qua giao diện:

```bash
GEMINI_API_KEY=...            # lấy tại https://aistudio.google.com/apikey
GEMINI_MODEL=gemini-2.5-flash # không bắt buộc
```

Thứ tự ưu tiên: **cấu hình trong hệ thống → `.env` → tắt**. Chưa có khoá ở đâu cả thì hệ thống vẫn
chạy bình thường, các nút AI **tự ẩn** khỏi giao diện thay vì báo lỗi khi người dùng bấm vào
(frontend hỏi `GET /api/ai/status` để biết).

### Các chức năng

| Chức năng | Vị trí trên giao diện | Ai dùng được |
|---|---|---|
| **Soạn giáo án** | Quản lý → *Soạn giáo án bằng AI* | Giảng viên, admin |
| **Tạo câu hỏi và đề kiểm tra** | Trình soạn bài tập → *Tạo bằng AI* | Giảng viên, admin |
| **Gợi ý điểm tự luận** | Màn hình chấm bài → *Nhờ AI gợi ý điểm* | Người có quyền chấm |
| **Trợ lý AI** | Nút trò chuyện ở góc phải màn hình | Mọi vai trò |
| **Nhật ký tác vụ AI** | Menu tài khoản → *Tác vụ AI gần đây* | Mọi vai trò |
| **Cấu hình nhà cung cấp AI** | Quản lý → *Cấu hình AI* | Chỉ admin |

**Tạo câu hỏi** nhận nội dung từ hai nguồn: chủ đề/nội dung bài học do giáo viên gõ, hoặc **tài liệu
PDF** tải lên (tối đa 12 MB). Tệp PDF chỉ đi qua máy chủ để gửi sang Gemini, **không được lưu lại**.
Giáo viên chọn số câu cho từng dạng và các mức độ nhận thức cần có.

**Trợ lý AI** dùng chung một cơ chế nhưng đổi hướng dẫn hệ thống theo vai trò người đăng nhập: với
học viên, trợ lý ưu tiên gợi ý và giải thích từng bước thay vì đưa ngay đáp án, và từ chối nội dung
ngoài phạm vi học tập; với giảng viên, trợ lý hỗ trợ tìm ý tưởng bài học, thiết kế hoạt động và gợi
ý câu hỏi.

### Hai lớp kiểm soát nội dung AI

Nội dung do AI tạo ra không mặc nhiên chính xác, nên hệ thống áp dụng hai lớp kiểm soát:

1. **Kiểm tra tự động (máy chủ).** Kết quả phải đúng cấu trúc JSON và **đủ số câu của từng dạng**;
   câu trắc nghiệm phải có đúng 4 phương án và một đáp án hợp lệ; câu đúng/sai cần ít nhất 2 phát
   biểu, mỗi phát biểu có giá trị rõ ràng; câu điền khuyết phải có dấu `___` và ít nhất một đáp án;
   câu tự luận phải có đáp án gợi ý và tiêu chí chấm; không câu nào được trùng nội dung. Không đạt
   thì máy chủ **tự yêu cầu AI sinh lại** (tối đa 3 lần), kèm mô tả lỗi của lần trước.
2. **Kiểm tra của con người.** Giáo viên duyệt và chỉnh sửa câu hỏi, giáo án trước khi giao cho học
   viên; điểm tự luận do AI gợi ý **chỉ được ghi nhận sau khi giáo viên xác nhận** — cột `score`
   không thay đổi khi bấm *Nhờ AI gợi ý điểm*, gợi ý được lưu riêng ở `ai_score` / `ai_comment`.

Hệ thống kiểm tra được **hình thức** của câu hỏi nhưng không tự khẳng định được tính đúng đắn về
kiến thức — vai trò duyệt của giáo viên là bắt buộc.

### Thiết kế prompt

Mỗi chức năng có một mẫu prompt riêng gồm năm phần: **Vai trò · Bối cảnh · Nhiệm vụ · Ràng buộc ·
Định dạng đầu ra**. Hệ thống tự điền thông tin từ lựa chọn của giáo viên, nên người dùng không cần
biết cách viết prompt. Xem `backend/internal/ai/prompt.go`.

---

## Thống kê điểm bài tập của học viên

Tab **Thống kê điểm** trong trang chương trình là bảng điểm của cả lớp: mỗi dòng một học viên, mỗi
cột một bài tập (đúng thứ tự trên cây nội dung), ô giao nhau là **điểm của lượt cao nhất** — đúng con
số học viên nhìn thấy ở màn hình bài tập. Bấm vào ô để mở thẳng bài nộp đó (chấm bài hoặc xem lại).

- Đầu bảng có năm số tổng: số học viên, số bài tập, tỉ lệ đã nộp, điểm tổng trung bình và số ô còn
  **chờ chấm**. Dòng cuối bảng là **trung bình lớp** từng bài kèm số lượt đã nộp và số học viên đạt.
- Điểm hiện màu xanh khi đạt ngưỡng *Điểm đạt* của bài (không đặt ngưỡng thì lấy mốc 50% thang điểm),
  đỏ khi chưa đạt, và vàng kèm chữ *chờ chấm* khi bài còn câu tự luận chưa vào điểm.
- Ô để trống (`—`) nghĩa là học viên chưa nộp bài đó lần nào.
- Có ô tìm theo tên/email, sắp xếp theo tên · điểm tổng · số bài đã nộp, và nút **Xuất CSV** (kèm BOM
  nên mở bằng Excel không vỡ dấu tiếng Việt).
- Học viên đã bị gỡ ghi danh nhưng còn bài nộp cũ vẫn hiện, kèm nhãn *đã gỡ ghi danh*, để điểm cũ
  không biến mất khỏi báo cáo. Giảng viên của chương trình không tính vào bảng.

Vai trò **Giám sát** xem được tab này như admin và giảng viên (chỉ xem, không chấm bài). Thống kê chi
tiết **theo từng câu hỏi** của một bài tập vẫn nằm ở tab *Kết quả học viên* trong chính bài tập đó.

---

## Cấu trúc mã nguồn

```
backend/
  cmd/server/          Điểm khởi động, seed admin, graceful shutdown
  internal/
    ai/                Lớp Service AI: dựng prompt, gọi Gemini, kiểm tra và chuẩn hoá kết quả
    api/               Router chi và toàn bộ handler HTTP (kể cả settings.go, dashboard.go, ai.go)
    auth/              JWT, bcrypt, OAuth2 Google, middleware phân quyền
    config/            Nạp cấu hình mặc định từ .env (override được qua app_settings trong DB)
    database/          Kết nối pgx pool và bộ chạy migration
    migrations/        File .sql nhúng vào binary
    models/            Kiểu dữ liệu dùng chung
    store/             Toàn bộ truy vấn Postgres (dashboard.go, settings.go, attempts.go, ai.go…)
                       grading.go: chấm điểm bốn dạng câu hỏi
                       workspace.go: lịch học, tiến độ, cảnh báo học tập
    util/              Chuyển link Google Drive sang URL nhúng, sinh slug tiếng Việt
frontend/
  src/
    api/               Client gọi REST và kiểu TypeScript
    components/        Cây kéo-thả, trình soạn nội dung/câu hỏi, trình làm bài, nhập hàng loạt,
                       GenerateQuestionsModal (tạo đề bằng AI), AIAssistant (trợ lý AI)
    pages/             Các màn hình, gồm AdminDashboardPage, GoogleSettingsPage, PreviewPage,
                       LessonPlanPage (soạn giáo án bằng AI), AITasksPage,
                       WorkspacePage (trang tổng quan), AISettingsPage, SignupSettingsPage
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
| `questions`, `question_options` | Câu hỏi (có `code` cố định, `level`, `sample_answer`, `rubric`) và phương án trả lời |
| `enrollments` | Ghi danh học viên/giảng viên vào chương trình |
| `lesson_progress` | Đánh dấu hoàn thành từng bài |
| `submissions`, `submission_answers` | Bài nộp, điểm tự động và điểm chấm tay |
| `attempt_sessions` | Phiên làm bài đang mở — mốc bắt đầu/hết giờ cho bài có tính thời gian |
| `app_settings` | Cấu hình sửa được qua giao diện admin (đăng nhập Google, nhà cung cấp AI, tự đăng ký) |
| `password_reset_requests` | Yêu cầu quên mật khẩu chờ admin xử lý |
| `materials` | Kho tài liệu dùng chung toàn hệ thống (lưu đường dẫn, không lưu tệp) |
| `lesson_plans` | Giáo án do AI soạn, giáo viên chỉnh sửa và lưu lại |
| `ai_tasks` | Nhật ký các lần gọi Gemini — hiện ở mục *Tác vụ AI gần đây* |
| `ai_conversations`, `ai_messages` | Hội thoại với trợ lý AI, dùng làm ngữ cảnh cho lượt hỏi sau |

Ý nghĩa của `question_options` thay đổi theo `questions.type`:

| `type` | Mỗi dòng `question_options` là | `is_correct` |
|---|---|---|
| `single_choice` | một phương án trả lời | phương án đúng |
| `true_false` | một phát biểu | giá trị Đúng/Sai của phát biểu |
| `fill_blank` | một đáp án được chấp nhận | luôn `true` |
| `essay` | *(không có dòng nào)* | — |

---

## Trước khi đưa lên production

> Hướng dẫn triển khai đầy đủ bằng Docker (backend) + Cloudflare Worker (frontend): xem
> [DEPLOY.md](DEPLOY.md).

- [ ] Đặt `JWT_SECRET` bằng chuỗi ngẫu nhiên (`openssl rand -base64 48`) — backend từ chối khởi động
      với `APP_ENV=production` nếu để trống.
- [ ] Đổi mật khẩu tài khoản `SEED_ADMIN_EMAIL`.
- [ ] Nếu bật **tự đăng ký tài khoản**, nhớ đặt giới hạn domain email — để trống nghĩa là bất kỳ ai
      biết địa chỉ trang web cũng tạo được tài khoản học viên.
- [ ] Đổi mật khẩu Postgres trong `docker-compose.yml` và `DATABASE_URL`.
- [ ] Đặt `APP_ENV=production`, `ALLOW_ORIGINS` và `FRONTEND_URL` theo domain thật.
- [ ] Đặt khoá AI qua **Quản lý → Cấu hình AI** (khuyến nghị) hoặc bằng biến môi trường
      `GEMINI_API_KEY` trên máy chủ — không commit khoá vào repo. Kiểm tra hạn mức của khoá: mỗi lần
      tạo đề có thể tốn tới 3 lượt gọi API nếu kết quả đầu chưa đạt.
- [ ] Chạy sau TLS; cập nhật `GOOGLE_REDIRECT_URL` sang `https://` và khai báo lại trên Google Cloud.
      Sau khi lên production, đổi Client ID/Secret hoặc bật/tắt Google login làm được ngay qua
      **Quản lý → Đăng nhập Google** — không cần sửa `.env` hay khởi động lại.

---

## Kiểm thử

```bash
cd backend && go test ./...
```

Phạm vi đang có: chấm điểm cả bốn dạng câu hỏi (điểm thành phần của câu đúng/sai, chuẩn hoá câu trả
lời của câu điền khuyết kể cả tiếng Việt gõ tổ hợp), kiểm tra tự động kết quả AI trả về trước khi
nhận, chuyển link Google Drive sang URL nhúng, sinh slug tiếng Việt, tính ổn định của thuật toán
trộn đề, chống giả mạo `state` trong luồng OAuth Google (sửa nonce, kéo dài hạn, đổi chữ ký, ký bằng
khoá khác, state hết hạn), và bộ lọc domain email được phép đăng nhập (dùng chung cho cả đăng nhập
Google lẫn tự đăng ký).

Một số test cần database thật (chạy migration, lưu đủ bốn dạng câu hỏi rồi nộp bài để đối chiếu
điểm, luồng quên mật khẩu, các truy vấn của trang tổng quan). Chúng tự bỏ qua khi thiếu biến môi
trường, nên bật riêng khi cần:

```bash
cd backend && TEST_DATABASE_URL="$DATABASE_URL" go test ./internal/store/ -v
```

## Bài giảng AI và trợ lý theo khóa học

- Trong **Quản lý chương trình → Tạo bài giảng AI**, nhập chủ đề, trình độ, mục tiêu,
  thời lượng và nội dung tham khảo. Giáo viên xem trước/chỉnh Markdown rồi **lưu bản nháp**
  vào chương đã chọn. Mở bài trong trình soạn thảo để xuất bản sau khi duyệt.
- Trong bài tập, **Tạo câu hỏi bằng AI → Bài giảng trong khóa học** lấy văn bản bài học đã
  lưu làm nguồn. Có thể ưu tiên một bài, hoặc lấy các bài liên quan theo chủ đề; các nguồn
  nhập nội dung và PDF vẫn dùng được. Câu hỏi luôn cần giáo viên duyệt trước khi lưu.
- Khi đang học, trợ lý nhận ID khóa học và bài đang mở. Máy chủ kiểm tra quyền, trạng thái
  xuất bản và cả nhánh cha bị ẩn/khóa trước khi chọn tối đa 6 nguồn văn bản liên quan.
  Ưu tiên bài đang mở; giới hạn khoảng 24.000 ký tự nội dung để tránh ngữ cảnh quá dài.
  Trợ lý có danh sách nguồn và hướng dẫn dẫn liên kết bài học trong câu trả lời.
- Mỗi bài có lịch sử hội thoại riêng. Hội thoại chung vẫn dùng ở trang tổng quan.
- Link video/PDF/Drive không tự động được đọc hoặc phiên âm: trợ lý sử dụng mô tả và phần
  nội dung/ghi chú văn bản đã lưu. Không cung cấp đáp án, lời giải hoặc tiêu chí chấm bài tập
  làm nguồn cho trợ lý học tập. Việc bám nguồn của mô hình vẫn cần kiểm chứng thực tế.
- Có 8 bìa minh họa SVG cục bộ; chọn trong **Cài đặt → Ảnh bìa** hoặc để tự động theo mã khóa.
  Link ảnh riêng vẫn được ưu tiên và có bìa dự phòng khi ảnh không tải được.

Migration `0012_chat_course_context.sql` bổ sung phạm vi hội thoại, tự chạy khi khởi động
backend mới. Bật nhà cung cấp AI trong trang cấu hình để dùng các nút tạo nội dung.
