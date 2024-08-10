USER ?= defaultUser
PASS ?= defaultPass

login:
	@echo "login with: ${USER}"
	@go run . login --username $(USER) --password $(PASS)

logout:
	@go run . logout --username $(USER)

build:
	@go build -o bin/todo-cli-app

run: build
	@./bin/todo-cli-app

test:
	@go test -v ./...