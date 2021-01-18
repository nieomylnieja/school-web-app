export CGO_ENABLED=0

APP_NAME = backend
ENVFILE=../school-web-app.env.local
BUILD_DIR=.build

.PHONY: build

all: test build run

verify:
	go vet ./...

test: verify
	go test ./...

build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) .

build-docker:
	env GOOS=linux go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)_linux

run:
	env $$(cat $(ENVFILE)) $(BUILD_DIR)/$(APP_NAME)

help: Makefile
	@echo " Choose a command run in \033[32m"$(APP_NAME)"\033[0m:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

