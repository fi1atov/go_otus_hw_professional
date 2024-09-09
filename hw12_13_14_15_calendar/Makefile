BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -v -count=1 -race -timeout=30s ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.56.2

lint: install-lint-deps
	golangci-lint run ./...

install-goose:
ifeq ($(wildcard $(GOOSE_BIN)),)
	go get -u github.com/pressly/goose/cmd/goose
	go build -i -o bin/goose github.com/pressly/goose/cmd/goose
	go mod tidy
GOOSE_BIN:=bin/goose
endif

migrate-up:
	goose -dir migrations postgres "host=localhost port=6543 user=postgres password=postgres dbname=postgres sslmode=disable search_path=calendar" up

generate:
	protoc --proto_path=api/ --go_out=internal/server/grpc --go-grpc_out=internal/server/grpc api/EventService.proto

.PHONY: build run build-img run-img version test lint generate
