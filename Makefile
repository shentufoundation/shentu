
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::') # grab everything after the space in "github.com/cometbft/cometbft v0.34.7"
DOCKER := $(shell which docker)
BUILDDIR ?= $(CURDIR)/build

GO_SYSTEM_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1-2)
REQUIRE_GO_VERSION = 1.22

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=shentu \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=shentud \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(TM_VERSION)


ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

# The below include contains the tools target.
include devtools/Makefile

###############################################################################
###                              Build                                      ###
###############################################################################

check_version:
ifneq ($(shell [ "$(GO_SYSTEM_VERSION)" \< "$(REQUIRE_GO_VERSION)" ] && echo true),)
	@echo "ERROR: Minimum Go version $(REQUIRE_GO_VERSION) is required for $(VERSION) of Gaia."
	exit 1
endif

all: install release lint test-unit

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): check_version go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

vulncheck: $(BUILDDIR)/
	GOBIN=$(BUILDDIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	$(BUILDDIR)/govulncheck ./...

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@echo "--> Ensure dependencies have not been modified unless suppressed by SKIP_MOD_VERIFY"
ifndef SKIP_MOD_VERIFY
	go mod verify
endif
	go mod tidy
	@echo "--> Download go modules to local cache"
	go mod download

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go install github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/shentud -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf $(BUILDDIR)/ artifacts/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                Release                                  ###
###############################################################################
GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)
GORELEASER_IMAGE := ghcr.io/goreleaser/goreleaser-cross:v$(GO_VERSION)
COSMWASM_VERSION := $(shell go list -m github.com/CosmWasm/wasmvm/v2 | sed 's/.* //')

# create tag and run goreleaser without publishing
# errors are possible while running goreleaser - the process can run for >30 min
# if the build is failing due to timeouts use goreleaser-build-local instead
create-release-dry-run:
ifneq ($(strip $(TAG)),)
	@echo "--> Dry running release for tag: $(TAG)"
	@echo "--> Create tag: $(TAG) dry run"
	git tag -s $(TAG) -m $(TAG)
	git push origin $(TAG) --dry-run
	@echo "--> Delete local tag: $(TAG)"
	@git tag -d $(TAG)
	@echo "--> Running goreleaser"
	@go install github.com/goreleaser/goreleaser@latest
	@docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-e TM_VERSION=$(TM_VERSION) \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v `pwd`:/go/src/shentud \
		-w /go/src/shentud \
		$(GORELEASER_IMAGE) \
		release \
		--snapshot \
		--skip=publish \
		--verbose \
		--clean
	@rm -rf dist/
	@echo "--> Done create-release-dry-run for tag: $(TAG)"
else
	@echo "--> No tag specified, skipping tag release"
endif

# Build static binaries for linux/amd64 using docker buildx
# Pulled from neutron-org/neutron: https://github.com/neutron-org/neutron/blob/v4.2.2/Makefile#L107
build-static-linux-amd64: go.sum $(BUILDDIR)/
	$(DOCKER) buildx create --name shentubuilder || true
	$(DOCKER) buildx use shentubuilder
	$(DOCKER) buildx build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GIT_VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(COMMIT) \
		--build-arg BUILD_TAGS=$(build_tags_comma_sep),muslc \
		--platform linux/amd64 \
		-t shentud-static-amd64 \
		-f Dockerfile . \
		--load
	$(DOCKER) rm -f shentubinary || true
	$(DOCKER) create -ti --name shentubinary shentud-static-amd64
	$(DOCKER) cp shentubinary:/usr/local/bin/ $(BUILDDIR)/shentud-linux-amd64
	$(DOCKER) rm -f shentubinary

# uses goreleaser to create static binaries for darwin on local machine
goreleaser-build-local:
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-e TM_VERSION=$(TM_VERSION) \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v `pwd`:/go/src/shentud \
		-w /go/src/shentud \
		$(GORELEASER_IMAGE) \
		release \
		--snapshot \
		--skip=publish \
		--release-notes ./RELEASE_NOTES.md \
		--timeout 90m \
		--verbose

# uses goreleaser to create static binaries for linux an darwin
# requires access to GITHUB_TOKEN which has to be available in the CI environment
ifdef GITHUB_TOKEN
ci-release:
	docker run \
		--rm \
		-e CGO_ENABLED=1 \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e TM_VERSION=$(TM_VERSION) \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v `pwd`:/go/src/shentud \
		-w /go/src/shentud \
		$(GORELEASER_IMAGE) \
		release \
		--release-notes ./RELEASE_NOTES.md \
		--timeout=90m \
		--clean
else
ci-release:
	@echo "Error: GITHUB_TOKEN is not defined. Please define it before running 'make release'."
endif

# create tag and publish it
create-release:
ifneq ($(strip $(TAG)),)
	@echo "--> Running release for tag: $(TAG)"
	@echo "--> Create release tag: $(TAG)"
	git tag -s $(TAG) -m $(TAG)
	git push origin $(TAG)
	@echo "--> Done creating release tag: $(TAG)"
else
	@echo "--> No tag specified, skipping create-release"
endif

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

include sims.mk

PACKAGES_UNIT=$(shell go list ./... | grep -v -e '/tests/e2e')
PACKAGES_E2E=$(shell go list ./... | grep '/e2e')
TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-unit-cover test-race test-e2e

test-unit: ARGS=-timeout=5m -tags='norace'
test-unit: TEST_PACKAGES=$(PACKAGES_UNIT)
test-unit-cover: ARGS=-timeout=5m -tags='norace' -coverprofile=coverage.txt -covermode=atomic
test-unit-cover: TEST_PACKAGES=$(PACKAGES_UNIT)
test-race: ARGS=-timeout=5m -race
test-race: TEST_PACKAGES=$(PACKAGES_UNIT)
test-e2e: ARGS=-timeout=35m -v
test-e2e: TEST_PACKAGES=$(PACKAGES_E2E)
$(TEST_TARGETS): run-tests

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	@echo "--> Running tests"
	@go test -mod=readonly -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@echo "--> Running tests"
	@go test -mod=readonly $(ARGS) $(TEST_PACKAGES)
endif

.PHONY: run-tests $(TEST_TARGETS)

docker-build-debug:
	@docker build -t shentuchain/shentud-e2e --build-arg IMG_TAG=debug -f Dockerfile .

# TODO: Push this to the Cosmos Dockerhub so we don't have to keep building it
# in CI.
docker-build-hermes:
	@cd tests/e2e/docker; docker build -t cosmos/hermes-e2e:1.0.0 -f hermes.Dockerfile .

docker-build-all: docker-build-debug docker-build-hermes

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
golangci_version=v1.60.1

tidy:
	@gofmt -s -w .
	@go mod tidy

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

format:
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go" -not -path "./crypto/keys/secp256k1/*" | xargs gofumpt -w -l
	$(golangci_lint_cmd) run --fix
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

start-localnet-ci: build
	rm -rf ~/.shentud-liveness
	./build/shentud init liveness -o --chain-id liveness --home ~/.shentud-liveness
	./build/shentud config set client keyring-backend test --home ~/.shentud-liveness
	./build/shentud keys add val --keyring-backend test --home ~/.shentud-liveness
	./build/shentud add-genesis-account val 10000000000000000000stake --home ~/.shentud-liveness --keyring-backend test
	./build/shentud genesis gentx val 1000000000stake --home ~/.shentud-liveness --chain-id liveness --keyring-backend test
	./build/shentud genesis collect-gentxs --home ~/.shentud-liveness
	sed -i 's/minimum-gas-prices = ".*"/minimum-gas-prices = "0stake"/' ~/.shentud-liveness/config/app.toml
	./build/shentud start --home ~/.shentud-liveness --x-crisis-skip-assert-invariants

.PHONY: start-localnet-ci

###############################################################################
###                                Docker                                   ###
###############################################################################

test-docker:
	@docker build -f contrib/Dockerfile.test -t ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) .
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:latest

test-docker-push: test-docker
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD)
	@docker push ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker push ${TEST_DOCKER_REPO}:latest

.PHONY: all build-linux install format lint draw-deps clean build \
	docker-build-debug docker-build-hermes docker-build-all



#build-linux: go.sum
#	CGO_ENABLED=1 LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build
#

