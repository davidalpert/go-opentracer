# project name
#PROJECT_NAME=$(shell basename "$(PWD)")
PROJECT_NAME=opentracer

# version of go tooling to use === must match version in go.mod ===
GO=go1.16

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# CGO?
CGO=0

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

rebuild: clean build ## rebuild

build: ./internal/version/detail.go ## build
	@echo "Building $(PROJECT_NAME) $(VERSION)"
	mkdir -p bin
	@CGO_ENABLED=$(CGO) $(GO) build -o bin/$(PROJECT_NAME) ./main.go

xbuild: ./internal/version/detail.go ## xbuild
	@echo "Building $(PROJECT_NAME) $(VERSION)"
	mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO) $(GO) build -o bin/linux-amd64/$(PROJECT_NAME) ./main.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO) $(GO) build -o bin/darwin-amd64/$(PROJECT_NAME) ./main.go
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=$(CGO) $(GO) build -o bin/darwin-arm64/$(PROJECT_NAME) ./main.go

ship:
	@echo "Building $(PROJECT_NAME) $(VERSION) for linux/amd64"
	mkdir -p bin
	@CGO_ENABLED=$(CGO) GOOS=linux GOARCH=amd64 $(GO) build -o bin/linux-amd64/$(PROJECT_NAME) ./main.go
	@echo "Copying $(PROJECT_NAME) $(VERSION) into the container"
	docker cp "$(HOME)/dev/personal/go-opentracer/bin/linux-amd64/$(PROJECT_NAME)" "$(shell docker ps | grep otel | awk '{print $$1}'):/$(PROJECT_NAME)"
	#docker cp "$(HOME)/dev/personal/go-opentracer/bin/linux-amd64/$(PROJECT_NAME)" "$(shell docker ps | grep default-centos-7 | awk '{print $$1}'):/$(PROJECT_NAME)"

rebuild: clean build ## rebuild

test: ./internal/version/detail.go ## run all tests
	CGO_ENABLED=$(CGO) $(GO) test ./...

test-verbose: ./internal/version/detail.go ## run all tests (with verbose flag)
	CGO_ENABLED=$(CGO) $(GO) test -v -timeout 1s ./...

changelog: ## Generate/update CHANGELOG.md
	git-chglog --output CHANGELOG.md

eq = $(and $(findstring $(1),$(2)),$(findstring $(2),$(1)))

release:
	$(if $(call eq,0,$(shell git diff-files --quiet; echo $$?)),, \
		$(error There are unstaged changes; clean your working directory before releasing.) \
	)
	$(if $(call eq,0,$(shell git diff-index --quiet --cached HEAD --; echo $$?)),, \
		$(error There are uncomitted changes; clean your working directory before releasing.) \
	)
	$(eval next_version := $(shell sbot predict version --mode ${BUMP_TYPE}))
	# echo "Current Version: ${VERSION}"
	# echo "   Next Version: ${next_version}"
	git-chglog --next-tag v$(next_version) --output CHANGELOG.md
	git add -f CHANGELOG.md
	git commit --message "docs: release notes for v$(next_version)"
	sbot release version --mode ${BUMP_TYPE}
	git show --no-patch --format=short v$(next_version)

SEMVER_TYPES := major minor patch
BUMP_TARGETS := $(addprefix release-,$(SEMVER_TYPES))
.PHONY: $(BUMP_TARGETS)
$(BUMP_TARGETS): ## bump version
	$(eval BUMP_TYPE := $(strip $(word 2,$(subst -, ,$@))))
	$(MAKE) release BUMP_TYPE=$(BUMP_TYPE)

.PHONY: help
help: Makefile
	@echo
	@echo " $(PROJECTNAME) $(VERSION) - available targets:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	printf "\033[36m%-30s\033[0m %s\n" '----------' '------------------'
	@echo $(BUMP_TARGETS) | tr ' ' '\n' | sort | sed -E 's/((.+)\-(.+))/\1: ## \2 \3 version/' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo
