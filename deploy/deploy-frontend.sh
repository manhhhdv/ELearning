#!/usr/bin/env bash
# Build và deploy frontend lên Cloudflare bằng wrangler CLI, dạng Worker phục vụ static assets
# (không phải Pages) — cấu hình ở frontend/wrangler.jsonc, logic proxy /api/* ở
# frontend/worker/index.js. Xem DEPLOY.md mục 4.
#
# Dùng:
#   deploy/deploy-frontend.sh [--env ENV] [--name WORKER_NAME] [--backend-origin URL] [--skip-build]
#
# Đăng nhập Cloudflare — chọn một trong hai cách trước khi dùng lần đầu:
#   a) `npx wrangler login` một lần (mở trình duyệt xác thực OAuth, lưu session cục bộ ở máy này).
#   b) Đặt biến môi trường CLOUDFLARE_API_TOKEN (tạo ở Cloudflare dashboard → My Profile → API
#      Tokens → dùng template "Edit Cloudflare Workers") và CLOUDFLARE_ACCOUNT_ID nếu tài khoản có
#      nhiều account — hợp để chạy tự động/CI không cần trình duyệt.
#
# --env             Deploy theo môi trường đặt tên trong frontend/wrangler.jsonc (vd "staging").
#                   Bỏ qua thì deploy cấu hình gốc (top-level) — dùng cho production. Cloudflare tự
#                   đặt tên Worker "<name>-<env>" trừ khi env đó tự khai "name" riêng.
# --name            Ghi đè tên Worker (mặc định: theo "name"/env đã chọn trong wrangler.jsonc).
# --backend-origin  Ghi đè biến BACKEND_ORIGIN lúc deploy thay vì sửa wrangler.jsonc — tiện cho một
#                   lần deploy thử, không cần thêm environment mới.
# --skip-build      Bỏ qua bước `npm install`/`npm run build`, dùng thẳng frontend/dist đã build sẵn.
#
# Ví dụ:
#   deploy/deploy-frontend.sh                       # production (cấu hình gốc)
#   deploy/deploy-frontend.sh --env staging         # môi trường "staging" trong wrangler.jsonc
#   deploy/deploy-frontend.sh --backend-origin https://api-preview.yourdomain.com --name tap-huan-preview

set -euo pipefail

usage() {
  echo "Dùng: $0 [--env ENV] [--name WORKER_NAME] [--backend-origin URL] [--skip-build]" >&2
  exit 1
}

ENV_NAME=""
WORKER_NAME=""
BACKEND_ORIGIN=""
SKIP_BUILD=0

while [ $# -gt 0 ]; do
  case "$1" in
    --env) ENV_NAME="$2"; shift 2 ;;
    --name) WORKER_NAME="$2"; shift 2 ;;
    --backend-origin) BACKEND_ORIGIN="$2"; shift 2 ;;
    --skip-build) SKIP_BUILD=1; shift ;;
    -h|--help) usage ;;
    *) echo "Không hiểu tuỳ chọn: $1" >&2; usage ;;
  esac
done

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="${REPO_ROOT}/frontend"
cd "${FRONTEND_DIR}"

if [ "${SKIP_BUILD}" -eq 0 ]; then
  echo "==> Cài phụ thuộc (nếu thiếu)"
  [ -d node_modules ] || npm install
  echo "==> Build frontend (tsc -b && vite build)"
  npm run build
fi

if [ ! -f dist/index.html ]; then
  echo "Không thấy frontend/dist/index.html — build lỗi hoặc chưa build. Bỏ --skip-build để build lại." >&2
  exit 1
fi

DEPLOY_ARGS=(deploy)
[ -n "${ENV_NAME}" ] && DEPLOY_ARGS+=(--env "${ENV_NAME}")
[ -n "${WORKER_NAME}" ] && DEPLOY_ARGS+=(--name "${WORKER_NAME}")
[ -n "${BACKEND_ORIGIN}" ] && DEPLOY_ARGS+=(--var "BACKEND_ORIGIN:${BACKEND_ORIGIN}")

echo "==> Deploy Worker${ENV_NAME:+ (env: ${ENV_NAME})} (assets: frontend/dist, script: frontend/worker/index.js)"
npx wrangler "${DEPLOY_ARGS[@]}"
