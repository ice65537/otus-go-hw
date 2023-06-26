.PHONY: build run version test lint generate down

BIN := "./bin/calendar"
#DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./api

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	docker compose -f ./deployments/docker-compose-test.yml up -d 
	sleep 5s
	$(BIN) --config ./configs/config.toml
	
down: 
	docker compose -f ./deployments/* down -v

#build-img:
#	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

#run-img: build-img
#	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	docker compose -f ./deployments/docker-compose-test.yml up -d
	sleep 5s
	go test -race ./internal/storage/...
	go test ./internal/server/...
	docker compose -f ./deployments/docker-compose-test.yml down -v

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...


