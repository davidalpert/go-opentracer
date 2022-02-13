PROJECT_NAME=$(shell basename "$(PWD)")

# version of go tooling to use === must match version in go.mod ===
GO=go1.16

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# CGO?
CGO=0

# project name
PROJECT_NAME=ddtracer

# project version
VERSION=$(shell cat VERSION)

# ---------------------- targets -------------------------------------

.PHONY: default
default: help

.PHONY: clean
clean:
	rm -rf ./bin
	rm -rf ./internal/version/detail.go

.PHONY: tidy
tidy: ## runs 'go mod tidy' with the current versioned go command
	$(GO) mod tidy

./internal/version/detail.go: VERSION
	$(MAKE) gen

.PHONY: gen
gen: ## invoke go generate
	@CGO_ENABLED=$(CGO) $(GO) generate ./...

run: ./internal/version/detail.go ## run direct from source
	#@echo "Running $(PROJECT_NAME) $(VERSION)"
	@CGO_ENABLED=$(CGO) $(GO) run ./main.go $(filter-out $@,$(MAKECMDGOALS))

build: ./internal/version/detail.go ## build
	@echo "Building $(PROJECT_NAME) $(VERSION)"
	mkdir -p bin
	@CGO_ENABLED=$(CGO) $(GO) build -o bin/$(PROJECT_NAME) ./main.go

ship:
	@echo "Building $(PROJECT_NAME) $(VERSION) for linux/amd64"
	mkdir -p bin
	@CGO_ENABLED=$(CGO) GOOS=linux GOARCH=amd64 $(GO) build -o bin/linux-amd64/$(PROJECT_NAME) ./main.go
	@echo "Copying $(PROJECT_NAME) $(VERSION) into the container"
	#docker cp "$(HOME)/dev/personal/go-opentracer/bin/linux-amd64/$(PROJECT_NAME)" "$(shell docker ps | grep otel | awk '{print $$1}'):/$(PROJECT_NAME)"
	docker cp "$(HOME)/dev/personal/go-opentracer/bin/linux-amd64/$(PROJECT_NAME)" "$(shell docker ps | grep default-centos-7 | awk '{print $$1}'):/$(PROJECT_NAME)"

rebuild: clean build ## rebuild

test: ./internal/version/detail.go ## run all tests
	CGO_ENABLED=$(CGO) $(GO) test ./...

test-verbose: ./internal/version/detail.go ## run all tests (with verbose flag)
	CGO_ENABLED=$(CGO) $(GO) test -v -timeout 1s ./...

.PHONY: help
help: Makefile
	@echo
	@echo " $(PROJECT_NAME) $(VERSION) - available targets:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo
