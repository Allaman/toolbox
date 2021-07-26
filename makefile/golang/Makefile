.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

GO_VERSION:=1.16
PROJECT_NAME:=FooBar
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

prepare: ## Prepare your go setup
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go mod tidy

format: ## Format the project
	go fmt $(shell go list ./...)
	go vet $(shell go list ./...)

lint: ## Lint the project
	staticcheck $(shell go list ./...)

build: ## Build the project
	go clean
	go build

all: prepare format lint release  ## Format, Lint, and release the project

run: ## Run the project
	go run main.go $(ARGS)

release: ## Build releases for different plattforms
	GOOS=linux GOARCH=amd64 go build -o release/$(PROJECTNAME)-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o release/$(PROJECTNAME)-windows-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o release/$(PROJECTNAME)-darwin-amd64 main.go

clean: ## Clean project dir
	go clean
	sudo rm -fr release

docker-build: ## Build using a docker container
	docker run --rm -v "$(PWD)":/usr/src/foobar -w /usr/src/foobar golang:${GO_VERSION} make all

build-image: ## Build docker image
	docker build . -t ${PROJECT_NAME}
