# Triển khai production

Backend chạy bằng Docker trên một server riêng (VPS, máy nội bộ...). Frontend build tĩnh và
deploy lên **Cloudflare Pages**, không chạy trong Docker.

## Kiến trúc

```
Trình duyệt
   │  https://your-app.pages.dev
   ▼
Cloudflare Pages (frontend, static)
   │  fetch('/api/...') — cùng origin với Pages
   │  → _redirects proxy /api/* sang backend (mã 200 = rewrite, không phải redirect)
   ▼
https://api.yourdomain.com                       ← Cloudflare Tunnel HOẶC Caddy (TLS)
   │
   ▼
Container `backend` (Go, cổng 8082, trong Docker)
   │
   ▼
Container `postgres` (chỉ truy cập nội bộ trong mạng Docker)
```

Frontend gọi API bằng đường dẫn tương đối `fetch('/api/...')`
([client.ts](frontend/src/api/client.ts:32)), không có biến `VITE_API_URL` nào để đổi domain
đích. Vì vậy khi frontend và backend nằm ở hai domain khác nhau (Cloudflare Pages vs server riêng),
bắt buộc phải có một lớp proxy `/api/*` ở phía Pages — đó là việc file
[`frontend/public/_redirects`](frontend/public/_redirects) đã làm sẵn. Xem mục 5.

## Chuẩn bị

- Server Linux đã cài **Docker** + **Docker Compose plugin** (`docker compose version`).
- Một domain cho backend, ví dụ `api.yourdomain.com` (bắt buộc nếu dùng Caddy; Cloudflare Tunnel
  cũng dùng domain này nhưng không cần trỏ DNS thủ công).
- Tài khoản Cloudflare (miễn phí) để deploy Pages, và Tunnel nếu chọn cách đó.

Các file cấu hình đã có sẵn trong repo:

| File | Dùng để |
|---|---|
| [`backend/Dockerfile`](backend/Dockerfile) | Build binary Go thành image chạy production |
| [`docker-compose.prod.yml`](docker-compose.prod.yml) | Dựng `postgres` + `backend` |
| [`docker-compose.tunnel.yml`](docker-compose.tunnel.yml) | Overlay thêm Cloudflare Tunnel (cách A) |
| [`docker-compose.proxy.yml`](docker-compose.proxy.yml) | Overlay thêm Caddy tự xin TLS (cách B) |
| [`deploy/build-and-ship.sh`](deploy/build-and-ship.sh) | Build image ở máy khác rồi scp lên server (mục 2b) |
| [`.env.example`](.env.example) | Biến cho `docker compose` (mật khẩu Postgres, domain, token tunnel) |
| [`backend/.env.production.example`](backend/.env.production.example) | Biến riêng của backend (JWT, admin, Google OAuth...) |
| [`deploy/Caddyfile`](deploy/Caddyfile) | Cấu hình Caddy nếu chọn tự xin TLS thay vì Cloudflare Tunnel |
| [`frontend/public/_redirects`](frontend/public/_redirects) | Cloudflare Pages proxy `/api/*` sang backend |

---

## 1. Chuẩn bị biến môi trường

Trên server, ở thư mục gốc repo:

```bash
cp .env.example .env
cp backend/.env.production.example backend/.env.production
```

Sửa `.env` (dùng bởi `docker compose` để dựng Postgres và proxy):

| Biến | Ý nghĩa |
|---|---|
| `POSTGRES_PASSWORD` | Đổi khỏi giá trị mặc định — **bắt buộc** |
| `DOMAIN` | Domain backend, chỉ cần nếu dùng Caddy (`docker-compose.proxy.yml`) |
| `CLOUDFLARE_TUNNEL_TOKEN` | Token tunnel, chỉ cần nếu dùng Cloudflare Tunnel (`docker-compose.tunnel.yml`) |

Sửa `backend/.env.production` (biến riêng của ứng dụng backend, xem chú thích đầy đủ trong file):

| Biến | Ý nghĩa |
|---|---|
| `JWT_SECRET` | **Bắt buộc** — sinh bằng `openssl rand -base64 48`. Backend từ chối khởi động nếu để trống khi `APP_ENV=production`. |
| `FRONTEND_URL`, `ALLOW_ORIGINS` | Domain thật của Cloudflare Pages, ví dụ `https://your-app.pages.dev`. Sai giá trị này thì CORS chặn mọi request từ frontend. |
| `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` | Tài khoản quản trị được tạo tự động lần chạy đầu tiên (chỉ khi DB chưa có admin nào) |
| `GOOGLE_REDIRECT_URL` | `https://api.yourdomain.com/api/auth/google/callback` — xem mục 6 |

`DATABASE_URL` trong file này **không cần sửa**: `docker-compose.prod.yml` tự ghi đè nó để trỏ vào
service `postgres` nội bộ, dùng đúng mật khẩu đã đặt ở `.env`.

---

## 2. Đưa image backend lên server

Chọn **một** trong hai cách nạp image, rồi chạy chung một lệnh khởi động.

### 2a. Build ngay trên server (đơn giản nhất)

Clone repo lên server (đã có `.env` và `backend/.env.production` từ bước 1), rồi:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Lệnh này build image backend từ [`backend/Dockerfile`](backend/Dockerfile) ngay trên server.

### 2b. Build ở máy khác rồi scp lên (không cần mã nguồn trên server)

Dùng khi server yếu, không muốn cài Go/toolchain lên server, hoặc build ở máy CI/máy dev rồi đẩy
image sang. Trên **máy build** (đã cài Docker + `docker buildx`, đã clone repo):

```bash
./deploy/build-and-ship.sh user@server --restart
```

Script [`deploy/build-and-ship.sh`](deploy/build-and-ship.sh):

1. Build image bằng `docker buildx build --platform linux/amd64` (mặc định amd64 — hầu hết VPS;
   đổi bằng `--platform linux/arm64` nếu server dùng CPU ARM). Stage build luôn chạy ở kiến trúc
   máy bạn (native, không qua giả lập QEMU) rồi dùng cross-compile của Go để ra binary đúng kiến
   trúc đích — nhanh và tránh được lỗi crash trình biên dịch Go khi ép chạy dưới QEMU.
2. `docker save | gzip` rồi `scp` file (vài MB) lên `~/elearning/images/` trên server.
3. Với `--restart`: SSH vào server, `docker load` rồi chạy `IMAGE_TAG=<tag> docker compose ...
   up -d backend` để chuyển ngay sang bản mới. Bỏ `--restart` (hoặc dùng `--load`) nếu chỉ muốn
   đẩy image lên trước, tự quyết định lúc nào cắt sang.

Chạy `./deploy/build-and-ship.sh --help` để xem đầy đủ tuỳ chọn (`--dir`, `--tag`, `--platform`,
`--ssh-opts`).

> Lần đầu dùng cách này: server cần có sẵn `docker-compose.prod.yml`, `.env` và
> `backend/.env.production` ở `REMOTE_DIR` (mặc định `~/elearning`) — scp thủ công 3 file này một
> lần trước khi chạy script, vì server không cần (và không có) phần còn lại của mã nguồn.
>
> Mỗi lần build ra một tag mới (`YYYYMMDD-HHMMSS`) chứ không ghi đè `latest`, nên rollback chỉ cần
> đổi `IMAGE_TAG` trong `.env` trên server rồi `docker compose -f docker-compose.prod.yml up -d
> backend` lại — miễn bản cũ vẫn còn nằm trong `docker images` trên server (dọn bớt bằng `docker
> image prune` nếu đầy đĩa).

---

Dù chọn cách nào ở trên, Postgres khởi động và chờ khoẻ (`healthcheck`) trước, backend tự chạy
migration (`backend/internal/migrations/*.sql`) và tạo tài khoản admin đầu tiên lúc khởi động —
không cần bước nào thêm.

Kiểm tra:

```bash
docker compose -f docker-compose.prod.yml ps
curl http://127.0.0.1:8082/api/health
docker compose -f docker-compose.prod.yml logs -f backend
```

Ở bước này backend mới chỉ nghe ở `127.0.0.1:8082` trên server — chưa ra được Internet. Sang
mục 3 để chọn cách expose ra ngoài.

---

## 3. Đưa backend ra Internet

Chọn **một** trong hai cách.

### Cách A — Cloudflare Tunnel (khuyến nghị)

Không cần mở cổng nào trên server/firewall, TLS do Cloudflare edge lo, hợp với việc frontend cũng
đang dùng Cloudflare.

1. Vào **Cloudflare Zero Trust → Networks → Tunnels → Create a tunnel**, chọn kiểu **Docker**,
   copy token hiển thị.
2. Trong tunnel vừa tạo, thêm **Public Hostname**: domain `api.yourdomain.com`, service
   `http://backend:8082` (đúng tên container trong mạng Docker nội bộ, không phải `localhost`).
3. Điền token vào `CLOUDFLARE_TUNNEL_TOKEN` trong `.env`, rồi chạy (backend/Postgres đã chạy sẵn
   từ mục 2 nên không cần `--build` ở đây, kể cả khi dùng cách 2b không có mã nguồn trên server):

```bash
docker compose -f docker-compose.prod.yml -f docker-compose.tunnel.yml up -d
```

`api.yourdomain.com` sẽ hoạt động ngay, không cần trỏ DNS thủ công (Cloudflare tự tạo bản ghi khi
thêm Public Hostname).

### Cách B — Caddy tự xin TLS (Let's Encrypt)

Dùng khi không muốn phụ thuộc Cloudflare Tunnel — cần mở cổng 80/443 và domain trỏ **DNS only**
(tắt proxy cam ☁️) thẳng vào IP server để Caddy xác thực ACME HTTP‑01.

```bash
docker compose -f docker-compose.prod.yml -f docker-compose.proxy.yml up -d
```

Caddy đọc `DOMAIN` từ `.env` ([`deploy/Caddyfile`](deploy/Caddyfile)), tự xin và gia hạn chứng chỉ.
Có thể bật lại proxy cam của Cloudflare sau khi chứng chỉ đã cấp xong nếu muốn (Caddy vẫn tự gia
hạn qua HTTP‑01 miễn cổng 80 còn mở tới server).

---

## 4. Deploy frontend lên Cloudflare Pages

Trong Cloudflare dashboard → **Workers & Pages → Create → Pages → Connect to Git**, chọn repo, rồi
cấu hình build:

| Trường | Giá trị |
|---|---|
| Framework preset | Vite |
| Build command | `npm run build` |
| Build output directory | `dist` |
| Root directory | `frontend` |

Không cần biến môi trường build nào — frontend không đọc `VITE_API_URL` hay biến tương tự, mọi
cấu hình API nằm ở `_redirects` (mục 5).

---

## 5. Nối frontend ↔ backend (`_redirects`)

Sửa domain trong [`frontend/public/_redirects`](frontend/public/_redirects) thành domain backend
thật rồi commit:

```
/api/*  https://api.yourdomain.com/api/:splat  200
```

Mã **200** báo Cloudflare Pages đây là **proxy** (rewrite ở edge), không phải redirect — trình
duyệt vẫn thấy origin của Pages nên `fetch('/api/...')` trong
[`client.ts`](frontend/src/api/client.ts:32) chạy đúng mà không cần sửa code, và không dính CORS
vì trình duyệt không thấy có cross-origin request nào. Vite copy nguyên file trong `public/` vào
gốc thư mục build, Cloudflare Pages tự nhận `_redirects` ở đó.

> Lưu ý: Cloudflare Pages giới hạn dung lượng và thời gian cho request đi qua proxy này (phù hợp
> cho API JSON thông thường, không phù hợp để proxy tải file lớn). Nếu sau này cần bỏ qua giới hạn
> này, có thể đổi kiến trúc sang gọi thẳng domain backend (cần sửa `client.ts` để dùng URL tuyệt đối
> và bật CORS qua `ALLOW_ORIGINS` — backend đã hỗ trợ sẵn CORS, có thể làm việc này sau nếu cần).

---

## 6. Cập nhật Google OAuth cho production

Nếu dùng đăng nhập Google, vào
[Google Cloud Console → Credentials](https://console.cloud.google.com/apis/credentials), mở OAuth
client đang dùng, thêm vào **Authorized redirect URIs**:

```
https://api.yourdomain.com/api/auth/google/callback
```

Giá trị này phải khớp chính xác với `GOOGLE_REDIRECT_URL` trong `backend/.env.production`.

---

## 7. Vận hành

**Xem log:**

```bash
docker compose -f docker-compose.prod.yml logs -f backend
```

**Cập nhật phiên bản mới** (migration mới tự chạy khi backend khởi động lại):

Nếu build ngay trên server (cách 2a):

```bash
git pull
docker compose -f docker-compose.prod.yml up -d --build backend
```

Nếu build ở máy khác rồi scp (cách 2b), chạy lại từ máy build:

```bash
./deploy/build-and-ship.sh user@server --restart
```

**Backup Postgres:**

```bash
docker compose -f docker-compose.prod.yml exec -T postgres \
  pg_dump -U elearning elearning | gzip > backup-$(date +%F).sql.gz
```

**Restore:**

```bash
gunzip -c backup-2026-08-19.sql.gz | \
  docker compose -f docker-compose.prod.yml exec -T postgres psql -U elearning elearning
```

**Vào thẳng Postgres** (không expose cổng ra host nên phải qua `exec`):

```bash
docker compose -f docker-compose.prod.yml exec postgres psql -U elearning
```

---

## 8. Checklist trước khi bật production

Xem checklist chung trong [README.md](README.md#trước-khi-đưa-lên-production), cộng thêm phần
riêng cho Docker:

- [ ] `POSTGRES_PASSWORD` trong `.env` đã đổi khỏi giá trị mặc định.
- [ ] `JWT_SECRET` trong `backend/.env.production` đã đặt (sinh bằng `openssl rand -base64 48`).
- [ ] `ALLOW_ORIGINS` / `FRONTEND_URL` trỏ đúng domain Cloudflare Pages thật (kể cả domain tuỳ
      chỉnh nếu có gắn thêm, không chỉ `*.pages.dev`).
- [ ] Cổng Postgres **không** được map ra host trong `docker-compose.prod.yml` (mặc định đã đúng).
- [ ] `GOOGLE_REDIRECT_URL` dùng `https://` và đã khai báo lại trên Google Cloud Console.
- [ ] Đã chọn và chạy một trong hai overlay (`docker-compose.tunnel.yml` hoặc `docker-compose.proxy.yml`) để backend có TLS.
- [ ] `frontend/public/_redirects` đã trỏ đúng domain backend thật, không còn `api.yourdomain.com`.

---

## Xử lý sự cố

| Triệu chứng | Nguyên nhân thường gặp |
|---|---|
| Backend thoát ngay lúc khởi động, log báo thiếu `JWT_SECRET` | Chưa điền `JWT_SECRET` trong `backend/.env.production` |
| Frontend gọi API bị lỗi CORS trong console trình duyệt | `_redirects` sai domain (proxy không chạy) hoặc `ALLOW_ORIGINS` không khớp domain Pages |
| Đăng nhập Google báo `redirect_uri_mismatch` | `GOOGLE_REDIRECT_URL` không khớp Authorized redirect URI trên Google Cloud Console |
| `docker compose up` báo thiếu biến (`variable is not set` / `required variable ... is missing`) | Chưa `cp .env.example .env`, hoặc chạy overlay `proxy`/`tunnel` mà thiếu `DOMAIN`/`CLOUDFLARE_TUNNEL_TOKEN` tương ứng |
| Caddy không xin được chứng chỉ | Cổng 80 chưa mở ra Internet, hoặc DNS domain đang bật proxy cam Cloudflare (cần chuyển tạm về DNS only) |
| Container backend thoát ngay với `exec format error` (dùng `build-and-ship.sh`) | Build sai kiến trúc CPU — kiểm tra CPU server bằng `ssh user@server uname -m` (`x86_64` → `--platform linux/amd64`, `aarch64`/`arm64` → `--platform linux/arm64`) rồi build lại |
| `build-and-ship.sh --restart` báo lỗi không thấy `docker-compose.prod.yml` trên server | Chưa scp 3 file `docker-compose.prod.yml`, `.env`, `backend/.env.production` vào `REMOTE_DIR` trên server lần đầu (xem mục 2b) |
