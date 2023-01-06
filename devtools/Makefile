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

##################################################################################
###                                 Protobuf                                   ###
##################################################################################
containerProtoVer=v0.3
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)

proto-all: proto-format proto-lint proto-check-breaking proto-gen

proto-gen:
	@./devtools/protocgen.sh

proto-gen-docker:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./devtools/protocgen.sh

proto-format:
	find ./ -not -path "./third_party/*" -name *.proto -exec clang-format -i {} \;

proto-lint:
	@buf check lint --error-format=json

proto-check-breaking:
	@buf check breaking --against '.git#branch=master'

proto-lint-docker:
	@$(DOCKER_BUF) check lint --error-format=json

proto-check-breaking-docker:
	@$(DOCKER_BUF) check breaking --against-input $(HTTPS_GIT)#branch=master

GOGO_PROTO_URL   = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL   = https://raw.githubusercontent.com/cosmos/cosmos-sdk/master/proto/cosmos

GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES    = third_party/proto/cosmos

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)/base/query/v1beta1/
	@curl -sSL $(COSMOS_PROTO_URL)/base/query/v1beta1/pagination.proto > $(COSMOS_PROTO_TYPES)/base/query/v1beta1/pagination.proto