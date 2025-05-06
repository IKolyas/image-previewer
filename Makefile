.PHONY: build run test lint clean migrate-up migrate-down docker-build docker-run docker-stop

# Переменные
BINARY_NAME=previewer
DOCKER_IMAGE=previewer
CONFIG_PATH=configs/previewer.json

# Только сборка
build:
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) ./cmd/previewer
	@echo "Build complete"

# Сборка и запуск
run: build
	@echo "Starting application..."
	./bin/$(BINARY_NAME) --config=$(CONFIG_PATH)

# Тесты
test:
	@echo "Running tests..."
	go test -v ./...

# linters
lint:
	@echo "Running linters..."
	golangci-lint run

# Очистка артефактов сборки
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out
	@echo "Clean complete"

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f build/Dockerfile . --no-cache

docker-run:
	@echo "Starting Docker container..."
	docker-compose -f deployments/docker-compose.yaml up -d --build

docker-stop:
	@echo "Stopping Docker container..."
	docker-compose -f deployments/docker-compose.yaml down

# Установка зависимостей
setup:
	@echo "Setting up development environment..."
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Done. Don't forget to install Docker if you haven't already"

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out