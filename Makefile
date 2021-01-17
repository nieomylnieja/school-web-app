APP_NAME = school-web-app
ENVFILE=school-web-app.env.local

.PHONY: build

all: build run

build:
	cd backend && make build-docker
	cd -
	env $$(cat $(ENVFILE)) docker-compose build

run:
	env $$(cat $(ENVFILE)) docker-compose up

local-env:
	@grep -ve "^##.*" $(APP_NAME).env  > $(ENVFILE)

help: Makefile
	@echo " Choose a command run in \033[32m"$(APP_NAME)"\033[0m:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
