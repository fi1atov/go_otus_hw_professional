### Calendar app automation recieps
# ¯¯¯¯¯¯¯¯

.DEFAULT_GOAL := help

SHELL=/bin/bash
SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_CONTEXT_PATH:=$(shell pwd)

DOCKER_MAKE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
DOCKER_MAKE_DIR:=$(dir $(DOCKER_MAKE_PATH))

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_IMG="calendar:develop"
DOCKER_BUILDKIT?=1
COMPOSE_DOCKER_CLI_BUILD?=1
DOCKER_REGISTRY?=
export

define __COMPOSE_CMD
source $(DOCKER_MAKE_DIR).env && \
		docker compose -f $(DOCKER_MAKE_DIR)make/deployments/docker-compose.app.yml \
		-f $(DOCKER_MAKE_DIR)make/deployments/docker-compose.yml -p calendar
endef

help: ## Display this help
	@IFS=$$'\n'; for line in `grep -h -E '^[a-zA-Z_#-]+:?.*?## .*$$' $(MAKEFILE_LIST)`; do if [ "$${line:0:2}" = "##" ]; then \
	echo $$line | awk 'BEGIN {FS = "## "}; {printf "\n\033[33m%s\033[0m\n", $$2}'; else \
	echo $$line | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'; fi; \
	done; unset IFS;

generate: ## Generate proto files
	go $@ ./...

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

build: build-scheduler build-sender build-calendar ## Build application

run-calendar: build-calendar ## Run calendar app
	source make/.env && $(BIN_CALENDAR) --config ./configs/config.toml

run-scheduler: build-scheduler ## Run scheduler app
	source make/.env && $(BIN_SCHEDULER) --config ./configs/config_scheduler.toml

run-sender: build-sender ## Run sender app
	source make/.env && $(BIN_SENDER) --config ./configs/config_sender.toml

version: build  ## Project version
	$(BIN_CALENDAR) version

test: ## Execute tests
	go test -race ./internal/... ./api/...

wait-for-db:
	sleep 5

integration: postgres-up wait-for-db migrate ## Execute integration tests
	go test --tags=$@ -race ./internal/... ./api/...


coverage: postgres-up wait-for-db migrate ## test coverage
	go test --tags=integration -coverprofile=coverage.out ./internal/...
	go tool cover -html coverage.out

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.56.2

lint: install-lint-deps ## Run liniter
	golangci-lint run --verbose --config=$$(pwd)/../.golangci.yml \
		--timeout 5m0s \
		--skip-dirs='/opt/hostedtoolcache/go|/go/pkg/mod' \
		--path-prefix=hw12_13_14_15_calendar ./...

goose:
	@(goose --version) || go install github.com/pressly/goose/v3/cmd/goose@latest

define __GOOSE_CMD
source make/.env && \
	goose --dir migrations postgres "postgres://$${DB_USERNAME}:$${DB_PASSWORD}@localhost:5432/$${DB_DATABASE}?sslmode=disable"
endef

migrate: ## Update db items
	$(__GOOSE_CMD) up

downgrade: ## Downgrade db items
	$(__GOOSE_CMD) down

migration_status: ## Current migrations status
	$(__GOOSE_CMD) status

.PHONY: build run version test lint help coverage goose

bi:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=BIN_FILE="/opt/calendar/calendar-app" \
		--build-arg=MAIN_SRC="cmd/calendar/*" \
		--build-arg=CONFIG_SRC="configs/config.toml" \
		-t ${DOCKER_REGISTRY}calendar_app \
		-f ${DOCKER_MAKE_DIR}/build/app/Dockerfile ${DOCKER_CONTEXT_PATH}

.PHONY: run-img
build-img:  ## Create docker image
	docker compose -f $(DOCKER_MAKE_DIR)make/deployments/docker-compose.app.yml build

.PHONY: build-img
run-img: build-img  ## Run  app container
	docker run $(DOCKER_IMG)

.PHONY: service-up
service-up:
	@$(__COMPOSE_CMD) up -d $(_SERVICE)

.PHONY: service-stop
service-stop:
	@$(__COMPOSE_CMD) stop $(_SERVICE)

.PHONY: service-restart
service-restart:
	@$(__COMPOSE_CMD) restart $(_SERVICE)

.PHONY: service-attach
service-attach:
	@$(__COMPOSE_CMD) exec $(_SERVICE) bash

.PHONY: service-logs
service-logs:
	@$(__COMPOSE_CMD) logs --tail 100 $(_SERVICE)

.PHONY: services-status
services-status: ## See current services status
	@$(__COMPOSE_CMD) ps

.PHONY: services-down
services-down: ## Remove all services
	@$(__COMPOSE_CMD) down -v

.PHONY: services-up
services-up: ## Up all services
	@$(__COMPOSE_CMD) up -d

.PHONY: postgres-up
postgres-up: ## Up dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-up

.PHONY: postgres-stop
postgres-stop: ## Stop dev postgress
	$(MAKE) _SERVICE=calendar_postgres service-stop

.PHONY: postgres-logs
postgres-logs: ## See dev postgress logs
	$(MAKE) _SERVICE=calendar_postgres service-logs

.PHONY: postgres-attach
postgres-attach: ## Attach to psql container
	$(MAKE) _SERVICE=calendar_postgres service-attach

.PHONY: mq-up
mq-up: ## Up dev postgress
	$(MAKE) _SERVICE=calendar_mq service-up

.PHONY: postgres-stop
mq-stop: ## Stop dev postgress
	$(MAKE) _SERVICE=calendar_mq service-stop

.PHONY: postgres-logs
mq-logs: ## See dev postgress logs
	$(MAKE) _SERVICE=calendar_mq service-logs

.PHONY: postgres-attach
mq-attach: ## Attach to psql container
	$(MAKE) _SERVICE=calendar_mq service-attach


.PHONY: up
up:	build-img # Up all dev environment
	DURATION=$(or $(DURATION),120s) $(__COMPOSE_CMD) up -d

.PHONY: down # Down all dev environment
down: services-down

__integration-t:
	sleep 3 && ginkgo -p -v ./pkg/tests -- --scheduler-duration="$(DURATION)"

integration-tests: DURATION=5s

integration-tests: down up __integration-t ## Execute smoke-tests
	$(MAKE) down
