# Включение других Makefile
include Makefile.default.mk
include Makefile.generators.mk

# Версии инструментов
LINT_VERSION ?= 1.59.0
MOCKERY_VERSION ?= 2.43.1

# Запуск приложения
run: ## Run app
	${DOCKER_COMPOSE} run --rm --service-ports app /bin/sh -c "go run cmd/photoloader/main.go server"

# Очистка сгенерированных файлов
clean: ## Clean up generated files
	rm -rf gen/
	rm -rf bin/
	rm -rf pkg/epgx/gen/

# Очистка снапшотов тестов
clean-snapshots: ## Remove test snapshots
	${DOCKER_COMPOSE} run --rm --no-deps app /bin/sh -c "rm -rf /src/tests/snapshots"

# Обновление снапшотов тестов
update-snapshots: clean-snapshots test ## Update test snapshots

# Сборка приложения в CI
build-ci: bin/ ## Build the application in CI
	CGO_ENABLED=0 go build -ldflags="-w -s $(LDFLAGS)" -o bin/$(APP_CMD) cmd/$(APP_CMD)/main.go
	bin/$(APP_CMD) ver

.PHONY: run clean clean-snapshots update-snapshots up build-ci 
