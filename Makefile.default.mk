# Цель по умолчанию
.DEFAULT_GOAL := default
default: tidy env docker-compose create_volumes mod lint test

# Пользовательские идентификаторы
USER_ID := $(shell id -u):$(shell id -g)

# Переменные для Docker Compose
COMPOSE_PREFIX ?= GIT_DOMAIN=${GIT_DOMAIN} GO_VERSION=${GO_VERSION} LINT_VERSION=${LINT_VERSION} MOCKERY_VERSION=${MOCKERY_VERSION} GOPRIVATE=${GO_PRIVATE}
DOCKER_COMPOSE ?= ${COMPOSE_PREFIX} docker-compose

# Команды для тестирования
TEST_CMD ?= go test --race -tags=integration ./...
WITH_COVER ?= -coverprofile=bin/cover.out && go tool cover -html=bin/cover.out -o bin/cover.html

# Название пакета и команды приложения
PACKAGE_NAME ?= $(shell head -n 1 go.mod | cut -d " " -f 2)
APP_CMD ?= $(shell echo ${PACKAGE_NAME} | awk -F '/' '{print $$NF}')
GIT_DOMAIN ?= $(shell echo ${PACKAGE_NAME} | cut -d "/" -f 1)
GO_VERSION ?= $(shell grep "^go " go.mod | awk '{print $$2}')

# Условие для установки переменной GOPRIVATE
ifeq ($(GIT_DOMAIN),github.com)
	GO_PRIVATE := $(shell head -n 1 go.mod | awk '{split($$2,a,"/"); print a[1]"/"a[2]}')
else
	GO_PRIVATE := $(GIT_DOMAIN)
endif

# Флаги компиляции
LDFLAGS := -X $(PACKAGE_NAME)/config.ServiceName=$(PACKAGE_NAME) \
    -X $(PACKAGE_NAME)/config.AppName=$(APP_CMD) \
    -X $(PACKAGE_NAME)/config.GitHash=$$(git rev-parse HEAD) \
    -X $(PACKAGE_NAME)/config.Version=$$(git describe --tags) \
    -X $(PACKAGE_NAME)/config.BuildAt=$$(date --utc +%FT%TZ)

# Цели для работы с Docker и сборкой проекта
build: bin/ ## Build the application
	$(DOCKER_COMPOSE) run --rm --no-deps app /bin/sh -c "go build -v -ldflags=\"-w -s $(LDFLAGS)\" -o bin/$(APP_CMD) cmd/$(APP_CMD)/main.go"

mod: ## Download go modules
	$(DOCKER_COMPOSE) run --rm --no-deps app with_creds /bin/sh -c "go mod download"

lint: ## Run the linter
	$(DOCKER_COMPOSE) run --rm linter /bin/sh -c "golangci-lint run ./... -c .golangci.yml -v"

test: bin/ ## Run tests
	$(DOCKER_COMPOSE) run --rm app with_creds /bin/sh -c "$(TEST_CMD) $(WITH_COVER)"
	$(DOCKER_COMPOSE) down --volumes

down: ## Stop and remove infrastructure
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

create_volumes: ## Create docker cache volumes
	docker volume create go-mod-cache
	docker volume create go-build-cache
	docker volume create go-lint-cache

docker-compose: ## Generate local docker-compose.override.yml file
	test -s docker-compose.override.yml || cp docker-compose.override.sample.yml docker-compose.override.yml

env: ## Generate local .env file
	test -s .env || cp .env.sample .env

tidy: ## Run go mod tidy
	go mod tidy

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_\/-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

bin/: ## Create bin folder
	mkdir -p $@

rebuild: ## Rebuild Docker image
	$(DOCKER_COMPOSE) build --no-cache

.PHONY: build mod lint test down create_volumes docker-compose tidy env help bin/ rebuild

%:
	@true
