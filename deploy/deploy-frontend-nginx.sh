#!/usr/bin/env bash
# Build frontend rồi rsync sang webroot trên server, dùng khi frontend được phục vụ bởi Nginx
# chạy trực tiếp trên server (domain hoctap.space) thay vì Cloudflare Worker.
# Xem deploy/nginx-hoctap.space.conf và DEPLOY.md mục 4C.
#
# Dùng:
#   deploy/deploy-frontend-nginx.sh user@server [--webroot /var/www/hoctap.space] [--skip-build] [--ssh-opts "-p 2222"]
#
# Yêu cầu: thư mục webroot trên server đã tồn tại và user SSH có quyền ghi vào đó, ví dụ:
#   ssh user@server "sudo mkdir -p /var/www/hoctap.space && sudo chown \$(whoami) /var/www/hoctap.space"

set -euo pipefail

usage() {
  echo "Dùng: $0 user@server [--webroot DIR] [--skip-build] [--ssh-opts \"opts\"]" >&2
  exit 1
}

[ $# -ge 1 ] || usage
TARGET="$1"; shift

WEBROOT="/var/www/hoctap.space"
SKIP_BUILD=0
SSH_OPTS=""

while [ $# -gt 0 ]; do
  case "$1" in
    --webroot) WEBROOT="$2"; shift 2 ;;
    --skip-build) SKIP_BUILD=1; shift ;;
    --ssh-opts) SSH_OPTS="$2"; shift 2 ;;
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

echo "==> rsync dist/ -> ${TARGET}:${WEBROOT}/dist"
# shellcheck disable=SC2086
rsync -az --delete -e "ssh ${SSH_OPTS}" dist/ "${TARGET}:${WEBROOT}/dist/"

echo "==> Xong. Kiểm tra: https://hoctap.space"
