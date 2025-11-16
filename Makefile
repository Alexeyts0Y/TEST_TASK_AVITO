APP_NAME=avito-app
BIN_DIR=bin
GO_VERSION=1.24
VERSION=$(shell git describe --tags 2>/dev/null || echo "v0.0.0")
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

all: deps build

build:
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) ./cmd/server

run:
	go run $(LDFLAGS) ./cmd/server

clean:
	rm -rf $(BIN_DIR)

deps:
	go mod download
	go mod verify

lint:
	golangci-lint run ./...

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-clean:
	docker-compose down -v

docker-run: docker-build docker-up

install-tools:
	@echo "Установка golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Установка oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

help:
	@echo "Доступные цели:"
	@echo "  all           - запустить deps и build"
	@echo "  build         - собрать приложение"
	@echo "  run           - запустить приложение"
	@echo "  clean         - удалить собранные файлы"
	@echo "  deps          - загрузить зависимости"
	@echo "  lint          - запустить линтер"
	@echo "  docker-build  - собрать Docker-образ"
	@echo "  docker-up     - запустить контейнеры"
	@echo "  docker-down   - остановить контейнеры"
	@echo "  docker-clean  - остановить контейнеры и удалить volumes"
	@echo "  docker-run    - собрать и запустить контейнеры"
	@echo "  install-tools - установить инструменты разработки"

.PHONY: all build run clean deps lint docker-build docker-up docker-down docker-clean docker-run install-tools help