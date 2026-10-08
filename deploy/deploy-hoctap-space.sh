#!/usr/bin/env bash
# Deploy/update toàn bộ (backend + frontend) cho domain hoctap.space, chạy trên server
# thviet@210.245.53.88:2222, thư mục /workspace/thviet/2026/hoctap-space — TÁCH BIỆT hoàn toàn
# với deployment daotao.space đang chạy cùng server (khác container/port/image/database).
#
# Kiến trúc (xem thêm DEPLOY.md mục 4C):
#   Nginx (host, TLS certbot) --/api/*--> 127.0.0.1:8092 (container elearning_api_hoctap)
#                              --còn lại--> /workspace/thviet/2026/hoctap-space/frontend/dist
#
# Dùng:
#   deploy/deploy-hoctap-space.sh                # deploy cả backend + frontend
#   deploy/deploy-hoctap-space.sh --backend-only # chỉ build + ship backend
#   deploy/deploy-hoctap-space.sh --frontend-only# chỉ build + đẩy frontend
#
# Yêu cầu trên máy chạy script:
#   - Docker (build image backend), Node/npm (build frontend), rsync, ssh, scp.
#   - Đã add key SSH (id_rsa) vào ssh-agent, hoặc sẵn sàng nhập passphrase khi được hỏi:
#       eval "$(ssh-agent -s)" && ssh-add ~/.ssh/id_rsa
#
# Lần đầu deploy (tạo DB/container/nginx/TLS mới) thì làm thủ công theo DEPLOY.md mục 4C —
# script này chỉ dùng để CẬP NHẬT phiên bản mới cho deployment đã có sẵn.

set -euo pipefail

SERVER="thviet@210.245.53.88"
SSH_PORT="2222"
REMOTE_DIR="/workspace/thviet/2026/hoctap-space"
IMAGE_NAME="elearning-backend-hoctap"
BACKEND_PORT_HOST="8092"   # cổng backend bind trên server (Nginx proxy sang cổng này)
PLATFORM="linux/amd64"

DO_BACKEND=1
DO_FRONTEND=1

while [ $# -gt 0 ]; do
  case "$1" in
    --backend-only) DO_FRONTEND=0; shift ;;
    --frontend-only) DO_BACKEND=0; shift ;;
    -h|--help)
      echo "Dùng: $0 [--backend-only] [--frontend-only]" >&2
      exit 1
      ;;
    *) echo "Không hiểu tuỳ chọn: $1" >&2; exit 1 ;;
  esac
done

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SSH_OPTS=(-p "${SSH_PORT}")
SCP_OPTS=(-P "${SSH_PORT}")
TAG="$(date +%Y%m%d-%H%M%S)"

if [ "${DO_BACKEND}" -eq 1 ]; then
  echo "==> [backend] Build ${IMAGE_NAME}:${TAG} cho ${PLATFORM}"
  docker buildx build --platform "${PLATFORM}" \
    -f "${REPO_ROOT}/backend/Dockerfile" \
    -t "${IMAGE_NAME}:${TAG}" \
    --load \
    "${REPO_ROOT}/backend"
  docker tag "${IMAGE_NAME}:${TAG}" "${IMAGE_NAME}:latest"

  TMPDIR_LOCAL="$(mktemp -d)"
  trap 'rm -rf "${TMPDIR_LOCAL}"' EXIT
  ARCHIVE="${TMPDIR_LOCAL}/${IMAGE_NAME}-${TAG}.tar.gz"

  echo "==> [backend] Đóng gói image (kèm tag latest)"
  docker save "${IMAGE_NAME}:${TAG}" "${IMAGE_NAME}:latest" | gzip > "${ARCHIVE}"
  du -h "${ARCHIVE}"

  echo "==> [backend] scp lên ${SERVER}:${REMOTE_DIR}/"
  scp "${SCP_OPTS[@]}" "${ARCHIVE}" "${SERVER}:${REMOTE_DIR}/$(basename "${ARCHIVE}")"

  echo "==> [backend] docker load + restart container trên server"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "
    set -e
    cd '${REMOTE_DIR}'
    gunzip -c '$(basename "${ARCHIVE}")' | docker load
    rm -f '$(basename "${ARCHIVE}")'
    docker compose -f docker-compose.prod.yml up -d backend
  "

  echo "==> [backend] Kiểm tra health (cổng ${BACKEND_PORT_HOST})"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "sleep 3 && curl -sS http://127.0.0.1:${BACKEND_PORT_HOST}/api/health && echo"
fi

if [ "${DO_FRONTEND}" -eq 1 ]; then
  echo "==> [frontend] Cài phụ thuộc (nếu thiếu) + build"
  cd "${REPO_ROOT}/frontend"
  [ -d node_modules ] || npm install
  npm run build

  if [ ! -f dist/index.html ]; then
    echo "Không thấy frontend/dist/index.html — build lỗi." >&2
    exit 1
  fi

  echo "==> [frontend] rsync dist/ -> ${SERVER}:${REMOTE_DIR}/frontend/dist"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "mkdir -p '${REMOTE_DIR}/frontend/dist'"
  rsync -az --delete -e "ssh -p ${SSH_PORT}" dist/ "${SERVER}:${REMOTE_DIR}/frontend/dist/"
fi

echo
echo "==> Xong. Kiểm tra: https://hoctap.space"
