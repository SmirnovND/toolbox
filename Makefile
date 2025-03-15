.ONESHELL:
.PHONY: help test test-verbose test-coverage lint build clean

# По умолчанию выполняется команда help
.DEFAULT_GOAL := help

# Переменные
GO = go
GOFLAGS = -v
PACKAGES = ./...
COVERAGE_FILE = coverage.out

# Запуск тестов
test:
	$(GO) test $(PACKAGES)

# Запуск тестов с подробным выводом
test-verbose:
	$(GO) test $(GOFLAGS) $(PACKAGES)

# Запуск тестов с покрытием кода
test-coverage:
	$(GO) test -coverprofile=$(COVERAGE_FILE) $(PACKAGES)
	$(GO) tool cover -html=$(COVERAGE_FILE)

# Сборка проекта
build:
	$(GO) build -o bin/toolbox $(GOFLAGS) ./...

# Очистка временных файлов
clean:
	rm -rf bin
	rm -f $(COVERAGE_FILE)

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make test            - Запуск тестов"
	@echo "  make test-verbose    - Запуск тестов с подробным выводом"
	@echo "  make test-coverage   - Запуск тестов с покрытием кода"
	@echo "  make build           - Сборка проекта"
	@echo "  make clean           - Очистка временных файлов"
	@echo "  make help            - Показать эту справку"
	@echo ""
	@echo "Примечание: Для команды 'lint' требуется установленный golangci-lint."
	@echo "Установить его можно командой: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"