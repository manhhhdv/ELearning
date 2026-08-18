#!/usr/bin/env bash
# Build image backend cho kiến trúc server (mặc định linux/amd64), đóng gói rồi scp lên server —
# dùng khi build ngay trên server không tiện (server yếu, không có mã nguồn, muốn build ở CI...).
# Xem DEPLOY.md mục 2b.
#
# Dùng:
#   deploy/build-and-ship.sh user@server [--dir REMOTE_DIR] [--tag TAG] [--platform PLATFORM] \
#     [--ssh-opts "OPTS"] [--load] [--restart]
#
# Mặc định:
#   REMOTE_DIR = elearning (thư mục chứa docker-compose.prod.yml trên server, tính từ $HOME)
#   TAG        = ngày giờ hiện tại, dạng 20260819-171500
#   PLATFORM   = linux/amd64 — đổi thành linux/arm64 nếu server dùng CPU ARM (Graviton, Ampere...)
#
#   --load     sau khi scp xong, SSH vào server và `docker load` luôn image (mặc định: không).
#   --restart  bao gồm --load, rồi chạy `IMAGE_TAG=<tag> docker compose ... up -d backend` trên
#              server để chuyển sang phiên bản mới ngay (mặc định: không tự restart production).
#
# Yêu cầu một lần trên server trước khi dùng --restart: đã scp sẵn docker-compose.prod.yml,
# .env và backend/.env.production vào REMOTE_DIR (xem DEPLOY.md mục 2b).
#
# Ví dụ:
#   deploy/build-and-ship.sh deploy@1.2.3.4                     # chỉ build + scp
#   deploy/build-and-ship.sh deploy@1.2.3.4 --restart           # build + scp + load + chuyển đổi
#   deploy/build-and-ship.sh deploy@1.2.3.4 --dir /opt/elearning --tag v1.2.0 --load

set -euo pipefail

REMOTE_DIR="elearning"
TAG="$(date +%Y%m%d-%H%M%S)"
PLATFORM="linux/amd64"
SSH_OPTS=""
DO_LOAD=0
DO_RESTART=0

usage() {
  echo "Dùng: $0 user@server [--dir REMOTE_DIR] [--tag TAG] [--platform PLATFORM] [--ssh-opts \"OPTS\"] [--load] [--restart]" >&2
  exit 1
}

[ $# -ge 1 ] || usage
SERVER="$1"; shift

while [ $# -gt 0 ]; do
  case "$1" in
    --dir) REMOTE_DIR="$2"; shift 2 ;;
    --tag) TAG="$2"; shift 2 ;;
    --platform) PLATFORM="$2"; shift 2 ;;
    --ssh-opts) SSH_OPTS="$2"; shift 2 ;;
    --load) DO_LOAD=1; shift ;;
    --restart) DO_LOAD=1; DO_RESTART=1; shift ;;
    -h|--help) usage ;;
    *) echo "Không hiểu tuỳ chọn: $1" >&2; usage ;;
  esac
done

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="elearning-backend:${TAG}"

TMPDIR_LOCAL="$(mktemp -d)"
trap 'rm -rf "${TMPDIR_LOCAL}"' EXIT
ARCHIVE="${TMPDIR_LOCAL}/elearning-backend-${TAG}.tar.gz"

echo "==> Build ${IMAGE} cho ${PLATFORM}"
docker buildx build --platform "${PLATFORM}" \
  -f "${REPO_ROOT}/backend/Dockerfile" \
  -t "${IMAGE}" \
  --load \
  "${REPO_ROOT}/backend"

echo "==> Đóng gói thành $(basename "${ARCHIVE}")"
docker save "${IMAGE}" | gzip > "${ARCHIVE}"
du -h "${ARCHIVE}"

echo "==> scp lên ${SERVER}:${REMOTE_DIR}/images/"
ssh ${SSH_OPTS} "${SERVER}" "mkdir -p ${REMOTE_DIR}/images"
scp ${SSH_OPTS} "${ARCHIVE}" "${SERVER}:${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz"

if [ "${DO_LOAD}" -eq 1 ]; then
  echo "==> docker load trên server"
  ssh ${SSH_OPTS} "${SERVER}" "gunzip -c ${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz | docker load"
fi

if [ "${DO_RESTART}" -eq 1 ]; then
  echo "==> Chuyển backend sang ${IMAGE} trên server"
  ssh ${SSH_OPTS} "${SERVER}" "cd ${REMOTE_DIR} && IMAGE_TAG=${TAG} docker compose -f docker-compose.prod.yml up -d backend"
  echo
  echo "Đã chạy phiên bản mới (tag ${TAG})."
  echo "Muốn tag này thành mặc định cho các lần \`up -d\` sau: sửa IMAGE_TAG=${TAG} trong ${REMOTE_DIR}/.env trên server."
else
  echo
  echo "Image đã ở trên server với tag ${TAG}."
  echo "Để nạp và chuyển sang chạy phiên bản này:"
  echo "  ssh ${SERVER} 'gunzip -c ${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz | docker load'"
  echo "  ssh ${SERVER} 'cd ${REMOTE_DIR} && IMAGE_TAG=${TAG} docker compose -f docker-compose.prod.yml up -d backend'"
fi
