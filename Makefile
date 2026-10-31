# Makefile — TaxiLux

.PHONY: up down dev test lint logs

compose := docker compose
BACKEND_SRC := $(shell pwd)/backend
FRONTEND_SRC := $(shell pwd)/frontend

up:
	$(compose) up --build

down:
	$(compose) down

dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

test:
	docker run --rm -v "$(BACKEND_SRC):/src" -w /src golang:1.22-alpine sh -c "apk add --no-cache git >/dev/null 2>&1; go mod tidy && gofmt -l . && go vet ./... && go test ./..."
	cd "$(FRONTEND_SRC)" && npm ci && npm run build

lint:
	docker run --rm -v "$(BACKEND_SRC):/src" -w /src golang:1.22-alpine sh -c "apk add --no-cache git >/dev/null 2>&1; gofmt -l . && go vet ./..."
	cd "$(FRONTEND_SRC)" && npm ci && npm run check

logs:
	$(compose) logs -f --tail=100
