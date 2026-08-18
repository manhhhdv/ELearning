# syntax=docker/dockerfile:1

# ---- Build: biên dịch binary tĩnh, không phụ thuộc glibc ----
# FROM --platform=$BUILDPLATFORM: stage này luôn chạy bằng kiến trúc máy build (không qua QEMU),
# rồi dùng chính khả năng cross-compile của Go (GOOS/GOARCH) để ra binary cho kiến trúc đích —
# nhanh và ổn định hơn hẳn build toàn bộ image dưới QEMU (từng gặp Go compiler crash "illegal
# instruction" khi ép cả stage build chạy qua QEMU trên máy Apple Silicon).
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Runtime: ảnh nhỏ gọn, chạy bằng user thường ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build --chown=app:app /out/server ./server

USER app
EXPOSE 8082
ENTRYPOINT ["./server"]
