#!/usr/bin/env bash
# Build và deploy frontend lên Cloudflare Pages bằng wrangler CLI (không cần kết nối Git, dùng khi
# muốn deploy thủ công hoặc từ CI). Xem DEPLOY.md mục 4.
#
# Dùng:
#   deploy/deploy-frontend.sh PROJECT_NAME [--branch BRANCH] [--skip-build]
#
# Đăng nhập Cloudflare — chọn một trong hai cách trước khi chạy lần đầu:
#   a) `npx wrangler login` một lần (mở trình duyệt xác thực OAuth, lưu session cục bộ ở máy này).
#   b) Đặt biến môi trường CLOUDFLARE_API_TOKEN (tạo ở Cloudflare dashboard → My Profile → API
#      Tokens → dùng template "Edit Cloudflare Workers", thêm quyền Pages) và CLOUDFLARE_ACCOUNT_ID
#      nếu tài khoản có nhiều account — hợp để chạy tự động/CI không cần trình duyệt.
#
# --branch      Nhánh deploy tới (mặc định: nhánh hiện tại theo git, hoặc "main" nếu không rõ).
#               Deploy lên nhánh khác "production branch" của project sẽ ra một Preview URL riêng,
#               không cập nhật domain production/domain tuỳ chỉnh đã gắn.
# --skip-build  Bỏ qua bước `npm install`/`npm run build`, dùng thẳng frontend/dist đã build sẵn.
#
# Ví dụ:
#   deploy/deploy-frontend.sh tap-huan                     # deploy bản production
#   deploy/deploy-frontend.sh tap-huan --branch preview-x  # deploy bản xem trước, có URL riêng

set -euo pipefail

usage() {
  echo "Dùng: $0 PROJECT_NAME [--branch BRANCH] [--skip-build]" >&2
  exit 1
}

[ $# -ge 1 ] || usage
PROJECT_NAME="$1"; shift

BRANCH=""
SKIP_BUILD=0

while [ $# -gt 0 ]; do
  case "$1" in
    --branch) BRANCH="$2"; shift 2 ;;
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

echo "==> Đảm bảo project Pages \"${PROJECT_NAME}\" tồn tại (bỏ qua lỗi nếu đã có sẵn)"
npx wrangler pages project create "${PROJECT_NAME}" \
  --production-branch "${BRANCH:-main}" || true

echo "==> Deploy frontend/dist lên Cloudflare Pages (project: ${PROJECT_NAME})"
DEPLOY_ARGS=(pages deploy dist --project-name "${PROJECT_NAME}")
[ -n "${BRANCH}" ] && DEPLOY_ARGS+=(--branch "${BRANCH}")
npx wrangler "${DEPLOY_ARGS[@]}"
