#!/usr/bin/env bash
# Build image backend cho kiến trúc server (mặc định linux/amd64), đóng gói rồi scp lên server —
# dùng khi build ngay trên server không tiện (server yếu, không có mã nguồn, muốn build ở CI...).
# Xem DEPLOY.md mục 2b.
#
# Dùng:
#   deploy/build-and-ship.sh user@server [-p PORT] [--dir REMOTE_DIR] [--tag TAG] \
#     [--platform PLATFORM] [--ssh-opts "OPTS"] [--load] [--restart]
#
# Mặc định:
#   REMOTE_DIR = elearning (thư mục chứa docker-compose.prod.yml trên server, tính từ $HOME)
#   TAG        = ngày giờ hiện tại, dạng 20260819-171500
#   PLATFORM   = linux/amd64 — đổi thành linux/arm64 nếu server dùng CPU ARM (Graviton, Ampere...)
#   PORT       = 22 (dùng -p 2222 nếu server nghe trên cổng 2222)
#
#   --load     sau khi scp xong, SSH vào server và `docker load` luôn image (mặc định: không).
#   --restart  bao gồm --load, rồi chạy `docker compose ... up -d backend` trên server để chuyển
#              sang phiên bản mới ngay (mặc định: không tự restart production).
#
# Mỗi lần chạy, image được gắn cả tag theo thời gian (rollback) lẫn tag `latest` (dùng mặc định) —
# nhờ vậy mọi lệnh `docker compose ... up -d` sau này trên server, kể cả khi thêm Tunnel/Caddy,
# đều tự dùng đúng bản mới nhất mà không cần nhớ truyền IMAGE_TAG.
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
  echo "Dùng: $0 user@server [-p PORT] [--dir REMOTE_DIR] [--tag TAG] [--platform PLATFORM] [--ssh-opts \"OPTS\"] [--load] [--restart]" >&2
  exit 1
}

[ $# -ge 1 ] || usage
SERVER="$1"; shift
SSH_PORT=""
SCP_OPTS=""

while [ $# -gt 0 ]; do
  case "$1" in
    -p|--port)
      SSH_PORT="$2"
      shift 2
      ;;
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

if [ -n "${SSH_PORT}" ]; then
  SSH_OPTS="${SSH_OPTS:+${SSH_OPTS} }-p ${SSH_PORT}"
  SCP_OPTS="${SCP_OPTS:+${SCP_OPTS} }-P ${SSH_PORT}"
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE="elearning-backend:${TAG}"
IMAGE_LATEST="elearning-backend:latest"

TMPDIR_LOCAL="$(mktemp -d)"
trap 'rm -rf "${TMPDIR_LOCAL}"' EXIT
ARCHIVE="${TMPDIR_LOCAL}/elearning-backend-${TAG}.tar.gz"

echo "==> Build ${IMAGE} cho ${PLATFORM}"
docker buildx build --platform "${PLATFORM}" \
  -f "${REPO_ROOT}/backend/Dockerfile" \
  -t "${IMAGE}" \
  --load \
  "${REPO_ROOT}/backend"

# Luôn gắn thêm tag `latest` trỏ vào đúng bản vừa build. docker-compose.prod.yml mặc định dùng
# `elearning-backend:${IMAGE_TAG:-latest}` — nếu không làm bước này, một lệnh `docker compose up -d`
# bất kỳ chạy sau đó mà quên truyền IMAGE_TAG (ví dụ lúc bật thêm Caddy/Tunnel) sẽ tìm `:latest`,
# không thấy, rồi cố `pull`/`build` và fail vì server không có mã nguồn hay quyền truy cập registry.
docker tag "${IMAGE}" "${IMAGE_LATEST}"

echo "==> Đóng gói thành $(basename "${ARCHIVE}") (kèm tag latest)"
docker save "${IMAGE}" "${IMAGE_LATEST}" | gzip > "${ARCHIVE}"
du -h "${ARCHIVE}"

echo "==> scp lên ${SERVER}:${REMOTE_DIR}/images/"
ssh ${SSH_OPTS} "${SERVER}" "mkdir -p ${REMOTE_DIR}/images"
scp ${SCP_OPTS} "${ARCHIVE}" "${SERVER}:${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz"

if [ "${DO_LOAD}" -eq 1 ]; then
  echo "==> docker load trên server"
  ssh ${SSH_OPTS} "${SERVER}" "gunzip -c ${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz | docker load"
fi

if [ "${DO_RESTART}" -eq 1 ]; then
  echo "==> Chuyển backend sang ${IMAGE} trên server"
  ssh ${SSH_OPTS} "${SERVER}" "cd ${REMOTE_DIR} && docker compose -f docker-compose.prod.yml up -d backend"
  echo
  echo "Đã chạy phiên bản mới (tag ${TAG}, cũng là \`latest\` hiện tại)."
  echo "Mọi lệnh \`docker compose ... up -d\` sau này trên server (kể cả khi thêm Tunnel/Caddy) sẽ"
  echo "tự dùng đúng bản này — không cần truyền IMAGE_TAG. Muốn rollback về bản này sau khi đã ship"
  echo "bản mới hơn: IMAGE_TAG=${TAG} docker compose -f docker-compose.prod.yml up -d backend"
else
  echo
  echo "Image đã ở trên server với tag ${TAG} (và latest)."
  echo "Để nạp và chuyển sang chạy phiên bản này:"
  echo "  ssh ${SERVER} 'gunzip -c ${REMOTE_DIR}/images/elearning-backend-${TAG}.tar.gz | docker load'"
  echo "  ssh ${SERVER} 'cd ${REMOTE_DIR} && docker compose -f docker-compose.prod.yml up -d backend'"
fi
