###
# Find OS and Go environment
# GO contains the Go binary
# FS contains the OS file separator
###
ifeq ($(OS),Windows_NT)
  GO := $(shell where go.exe 2> NUL)
  FS := "\\"
else
  GO := $(shell command -v go 2> /dev/null)
  FS := "/"
endif

ifeq ($(GO),)
  $(error could not find go. Is it in PATH? $(GO))
endif

GOPATH ?= $(shell $(GO) env GOPATH)
GITHUBDIR := $(GOPATH)$(FS)src$(FS)github.com

###############################################################################
###                                Functions                                ###
###############################################################################

go_get = $(if $(findstring Windows_NT,$(OS)),\
IF NOT EXIST $(GITHUBDIR)$(FS)$(1)$(FS) ( mkdir $(GITHUBDIR)$(FS)$(1) ) else (cd .) &\
IF NOT EXIST $(GITHUBDIR)$(FS)$(1)$(FS)$(2)$(FS) ( cd $(GITHUBDIR)$(FS)$(1) && git clone https://github.com/$(1)/$(2) ) else (cd .) &\
,\
mkdir -p $(GITHUBDIR)$(FS)$(1) &&\
(test ! -d $(GITHUBDIR)$(FS)$(1)$(FS)$(2) && cd $(GITHUBDIR)$(FS)$(1) && git clone https://github.com/$(1)/$(2)) || true &&\
)\
cd $(GITHUBDIR)$(FS)$(1)$(FS)$(2) && git fetch origin && git checkout -q $(3)

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell cd $(shell dirname $(mkfile_path)); pwd)

###############################################################################
###                                 Tools                                   ###
###############################################################################

TOOLS_DESTDIR  ?= $(GOPATH)/bin

STATIK        = $(TOOLS_DESTDIR)/statik
RUNSIM        = $(TOOLS_DESTDIR)/runsim

all: tools

tools: statik runsim

golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(mkfile_dir)/install-golangci-lint.sh
	@echo "Installing golangci-lint..."
	@bash $(mkfile_dir)/install-golangci-lint.sh $(TOOLS_DESTDIR)

# Install the statik binary for which the version is specified in go.mod
statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@(go install github.com/rakyll/statik)

# Install the runsim binary for which the version is specified in go.mod
runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(go install github.com/cosmos/tools/cmd/runsim)

tools-clean:
	rm -f $(STATIK) $(RUNSIM)
	rm -f tools-stamp

.PHONY: all tools tools-clean golangci-lint

###############################################################################
###                                Protobuf                                 ###
###############################################################################

containerProtoVer=v0.7
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(containerProtoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./devtools/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \; ; fi

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

.PHONY: proto-all proto-gen proto-format proto-lint
