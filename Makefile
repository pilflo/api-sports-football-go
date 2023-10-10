LINT_COMMAND=golangci-lint run
FILES_LIST=$(shell ls -d */ | grep -v -E "vendor|tools|target|build")
GO_VERSION=1.21.2
TOOLS_DOCKER_IMAGE=go-check:$(GO_VERSION)

.PHONY: check.fmt check.imports check.lint check.test check.get.tools.image check.prepare

check.prepare: check.build

check: check.prepare check.imports check.fmt check.lint check.test

check.build:
	docker build -t $(TOOLS_DOCKER_IMAGE) --build-arg GO_VERSION=$(GO_VERSION) -f ./build/check/Dockerfile .

#help check.fmt: format go code
check.fmt: check.prepare
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(TOOLS_DOCKER_IMAGE) sh -c 'gofumpt -w $(FILES_LIST)'

#help check.imports: fix and format go imports
check.imports: check.prepare
	@docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c 'goimports -w $(FILES_LIST)'

#help check.lint: check if the go code is properly written, rules are in .golangci.yml
check.lint: check.prepare
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c '$(LINT_COMMAND)'

#help check.test: execute go tests, if using test container set TEST_CONTAINER_FLAGS in custom.mk
check.test: check.prepare
	docker run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_DOCKER_IMAGE) sh -c 'go test -mod=vendor ./...'

#help check.get.tools.image: returns the name of the docker image used for the ci tools
check.get.tools.image:
	@echo -n $(TOOLS_DOCKER_IMAGE)
