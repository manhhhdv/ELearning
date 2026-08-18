# Makefile của hệ thống đào tạo trực tuyến.
# Dùng bash để có trap dọn tiến trình con khi dừng bằng Ctrl+C.
SHELL := /bin/bash

BACKEND_PORT  := 8082
FRONTEND_PORT := 3006

.PHONY: help dev deps db down clean backend frontend test build fmt reset-db favicon

help:
	@echo "Phát triển"
	@echo "  make dev        Chạy tất cả: Postgres + API + giao diện (Ctrl+C để dừng hết)"
	@echo "  make backend    Chỉ chạy API Go              (cổng $(BACKEND_PORT))"
	@echo "  make frontend   Chỉ chạy giao diện React     (cổng $(FRONTEND_PORT))"
	@echo ""
	@echo "Cơ sở dữ liệu"
	@echo "  make db         Khởi động Postgres và chờ sẵn sàng (cổng 5433)"
	@echo "  make down       Dừng Postgres, giữ nguyên dữ liệu"
	@echo "  make reset-db   Xoá sạch dữ liệu rồi dựng lại từ đầu"
	@echo ""
	@echo "Khác"
	@echo "  make deps       Cài phụ thuộc Go và npm"
	@echo "  make test       Chạy test backend"
	@echo "  make build      Build backend và frontend cho production"
	@echo "  make fmt        Định dạng lại mã Go"
	@echo "  make favicon    Sinh lại favicon PNG/ICO từ favicon.svg"

# --- Môi trường phát triển ---------------------------------------------------

dev: db deps
	@echo ""
	@echo "  API       http://localhost:$(BACKEND_PORT)"
	@echo "  Giao diện http://localhost:$(FRONTEND_PORT)"
	@echo "  Ctrl+C để dừng tất cả"
	@echo ""
# 	@$(MAKE) --no-print-directory check-ports
	@# Dấu - ở đầu: Ctrl+C kết thúc bằng mã 130, không phải lỗi thật nên đừng để make kêu.
	-@stop() { \
		[ -n "$$stopping" ] && return; stopping=1; \
		echo; echo "Đang dừng API và giao diện…"; \
		kill 0; \
	}; \
	trap stop EXIT INT TERM; \
	( cd backend  && go run ./cmd/server 2>&1 | awk '{ print "\033[36m[api]\033[0m " $$0; fflush() }' ) & \
	( cd frontend && npm run dev        2>&1 | awk '{ print "\033[35m[web]\033[0m " $$0; fflush() }' ) & \
	wait

# Báo sớm nếu cổng đã bị chiếm, thay vì để tiến trình chết lặng lẽ giữa đống log.
check-ports:
	@for port in $(BACKEND_PORT) $(FRONTEND_PORT); do \
		if lsof -nP -iTCP:$$port -sTCP:LISTEN >/dev/null 2>&1; then \
			echo "Cổng $$port đang bị tiến trình khác chiếm:"; \
			lsof -nP -iTCP:$$port -sTCP:LISTEN | tail -n +2 | awk '{ print "   " $$1, "(PID " $$2 ")" }'; \
			echo "   Dừng tiến trình đó, hoặc đổi cổng trong backend/.env và frontend/vite.config.ts"; \
			exit 1; \
		fi; \
	done

deps:
	@cd backend && go mod download
	@if [ ! -d frontend/node_modules ]; then \
		echo "Cài phụ thuộc npm lần đầu…"; \
		cd frontend && npm install; \
	fi

backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

# --- Cơ sở dữ liệu -----------------------------------------------------------

db:
	@docker compose up -d --wait

down:
	docker compose down

reset-db:
	docker compose down -v
	@docker compose up -d --wait
	@echo "Đã xoá sạch dữ liệu. Chạy lại backend để tạo bảng và tài khoản quản trị đầu tiên."

clean: reset-db

# --- Kiểm tra và đóng gói ----------------------------------------------------

test:
	cd backend && go test ./...

build:
	cd backend && go build -o bin/server ./cmd/server
	cd frontend && npm run build

fmt:
	cd backend && gofmt -w .

# Sinh lại favicon PNG/ICO từ frontend/public/favicon.svg sau khi sửa logo.
favicon:
	cd frontend && python3 scripts/make-favicon.py
