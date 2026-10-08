#!/usr/bin/env bash
# Deploy/update toàn bộ (backend + frontend) cho domain troly.hoctap.space, chạy trên server
# thviet@210.245.53.88:2222, thư mục /workspace/thviet/2026/troly-hoctap-space — TÁCH BIỆT hoàn toàn
# với các deployment khác cùng server (khác project compose, container, image, cổng, database).
#
# Kiến trúc (xem DEPLOY.md mục 4C):
#   Nginx (host, TLS certbot) --/api/*--> 127.0.0.1:8093 (container elearning_api_troly)
#                              --còn lại--> /workspace/thviet/2026/troly-hoctap-space/frontend/dist
#
# Dùng:
#   ADMIN_EMAIL=ban@example.com deploy/deploy-troly-hoctap-space.sh --init   # lần đầu
#   deploy/deploy-troly-hoctap-space.sh                  # cập nhật cả backend + frontend
#   deploy/deploy-troly-hoctap-space.sh --backend-only   # chỉ build + ship backend
#   deploy/deploy-troly-hoctap-space.sh --frontend-only  # chỉ build + đẩy frontend
#
# --init (chạy một lần): tạo thư mục trên server, chép compose + cấu hình Nginx, sinh `.env` và
# `backend.env` với mật khẩu Postgres / JWT_SECRET / mật khẩu admin ngẫu nhiên (KHÔNG in ra màn
# hình — đọc lại bằng `ssh ... cat <REMOTE_DIR>/backend.env`), rồi cài site Nginx + xin TLS bằng
# certbot (cần sudo trên server, sẽ hỏi mật khẩu). File env đã có thì giữ nguyên, không ghi đè.
#
# Yêu cầu trên máy chạy script: Docker (buildx), Node/npm, rsync, ssh, scp; key SSH đã nạp vào agent.

set -euo pipefail

SERVER="thviet@210.245.53.88"
SSH_PORT="2222"
DOMAIN="troly.hoctap.space"
REMOTE_DIR="/workspace/thviet/2026/troly-hoctap-space"
IMAGE_NAME="elearning-backend-troly"
BACKEND_PORT_HOST="8093"   # phải khớp deploy/troly/env.example và cấu hình Nginx
PLATFORM="linux/amd64"

DO_INIT=0
DO_BACKEND=1
DO_FRONTEND=1

while [ $# -gt 0 ]; do
  case "$1" in
    --init) DO_INIT=1; shift ;;
    --backend-only) DO_FRONTEND=0; shift ;;
    --frontend-only) DO_BACKEND=0; shift ;;
    -h|--help)
      echo "Dùng: $0 [--init] [--backend-only] [--frontend-only]" >&2
      exit 1
      ;;
    *) echo "Không hiểu tuỳ chọn: $1" >&2; exit 1 ;;
  esac
done

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SSH_OPTS=(-p "${SSH_PORT}")
SCP_OPTS=(-P "${SSH_PORT}")
TAG="$(date +%Y%m%d-%H%M%S)"

if [ "${DO_INIT}" -eq 1 ]; then
  : "${ADMIN_EMAIL:?Đặt ADMIN_EMAIL=<email quản trị> khi chạy --init}"

  echo "==> [init] Kiểm tra cổng ${BACKEND_PORT_HOST} còn trống trên server"
  if ssh "${SSH_OPTS[@]}" "${SERVER}" "ss -ltn | awk '{print \$4}' | grep -q ':${BACKEND_PORT_HOST}\$'"; then
    if ! ssh "${SSH_OPTS[@]}" "${SERVER}" "docker ps --format '{{.Names}}' | grep -qx elearning_api_troly"; then
      echo "Cổng ${BACKEND_PORT_HOST} đã bị dịch vụ khác chiếm — đổi BACKEND_PORT_HOST ở script, env.example và Nginx." >&2
      exit 1
    fi
  fi

  echo "==> [init] Tạo ${REMOTE_DIR} và chép compose + cấu hình Nginx"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "mkdir -p '${REMOTE_DIR}/frontend/dist'"
  scp "${SCP_OPTS[@]}" \
    "${REPO_ROOT}/deploy/troly/docker-compose.yml" \
    "${REPO_ROOT}/deploy/troly/nginx-troly.hoctap.space.conf" \
    "${SERVER}:${REMOTE_DIR}/"

  echo "==> [init] Sinh .env và backend.env (nếu chưa có)"
  scp "${SCP_OPTS[@]}" "${REPO_ROOT}/deploy/troly/backend.env.example" "${SERVER}:${REMOTE_DIR}/backend.env.example"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "
    set -e
    cd '${REMOTE_DIR}'
    rand() { head -c 48 /dev/urandom | base64 | tr -d '/+=\n' | cut -c1-\$1; }
    if [ ! -f .env ]; then
      umask 077
      printf 'POSTGRES_USER=elearning\nPOSTGRES_PASSWORD=%s\nPOSTGRES_DB=elearning\nBACKEND_PORT_HOST=${BACKEND_PORT_HOST}\n' \"\$(rand 32)\" > .env
      echo '   đã tạo .env'
    fi
    if [ ! -f backend.env ]; then
      umask 077
      sed -e \"s|^JWT_SECRET=.*|JWT_SECRET=\$(rand 64)|\" \
          -e \"s|^SEED_ADMIN_EMAIL=.*|SEED_ADMIN_EMAIL=${ADMIN_EMAIL}|\" \
          -e \"s|^SEED_ADMIN_PASSWORD=.*|SEED_ADMIN_PASSWORD=\$(rand 16)|\" \
          backend.env.example > backend.env
      echo '   đã tạo backend.env (mật khẩu admin nằm trong file này)'
    fi
    rm -f backend.env.example
  "

  echo "==> [init] Cài site Nginx + xin TLS (sudo trên server)"
  ssh -t "${SSH_OPTS[@]}" "${SERVER}" "
    set -e
    sudo cp '${REMOTE_DIR}/nginx-troly.hoctap.space.conf' /etc/nginx/sites-available/${DOMAIN}
    sudo ln -sf /etc/nginx/sites-available/${DOMAIN} /etc/nginx/sites-enabled/${DOMAIN}
    sudo nginx -t && sudo systemctl reload nginx
    sudo certbot --nginx -d ${DOMAIN} --non-interactive --agree-tos --redirect -m '${ADMIN_EMAIL}'
  "
fi

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
  scp "${SCP_OPTS[@]}" "${REPO_ROOT}/deploy/troly/docker-compose.yml" "${SERVER}:${REMOTE_DIR}/docker-compose.yml"

  echo "==> [backend] docker load + khởi động lại trên server"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "
    set -e
    cd '${REMOTE_DIR}'
    gunzip -c '$(basename "${ARCHIVE}")' | docker load
    rm -f '$(basename "${ARCHIVE}")'
    docker compose up -d
  "

  echo "==> [backend] Kiểm tra health (cổng ${BACKEND_PORT_HOST})"
  ssh "${SSH_OPTS[@]}" "${SERVER}" "
    for i in \$(seq 1 30); do
      curl -fsS http://127.0.0.1:${BACKEND_PORT_HOST}/api/health && echo && exit 0
      sleep 2
    done
    echo 'Backend chưa sẵn sàng — xem: docker logs elearning_api_troly' >&2
    exit 1
  "
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
echo "==> Xong. Kiểm tra: https://${DOMAIN}"
