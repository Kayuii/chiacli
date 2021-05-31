GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOBUILD:=$(GOCMD) build
GOINSTALL:=$(GOCMD) install
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOVET:=$(GOCMD) vet
GOPATH:=$(shell $(GOCMD) env GOPATH)
u := $(if $(update),-u)

BINARY_NAME:=chiacli
PACKAGES:=$(shell $(GOLIST) github.com/kayuii/chiacli github.com/kayuii/chiacli/cmd/chiacli github.com/kayuii/chiacli/plot )
GOFILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE := on

all: test build

mini: test build-mini

.PHONY: build
build: deps
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/chiacli

.PHONY: build-mini
build-mini: deps
	$(GOBUILD) -ldflags "-s -w" -o $(BINARY_NAME)-mini ./cmd/chiacli

.PHONY: build-static
build-static: deps
	CGO_ENABLED=0 $(GOBUILD) -ldflags '-extldflags=-static -w -s ' -o $(BINARY_NAME)-static ./cmd/chiacli

.PHONY: install
install: deps
	$(GOINSTALL) ./cmd/chiacli

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: deps
deps:
	$(GOGET) github.com/urfave/cli
	$(GOGET) github.com/urfave/cli/v2
	$(GOGET) github.com/mackerelio/go-osstat/memory
	$(GOGET) golang.org/x/sys/unix
	$(GOGET) github.com/google/logger
	$(GOGET) github.com/go-cmd/cmd
	$(GOGET) github.com/kilic/bls12-381
	$(GOGET) golang.org/x/crypto/hkdf
	$(GOGET) github.com/stretchr/testify/require

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off $(GOGET) -v -u \
		golang.org/x/lint/golint

.PHONY: lint
lint: devel-deps
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: vet
vet: deps devel-deps
	$(GOVET) $(PACKAGES)

.PHONY: fmt
fmt:
	$(GOFMT) -s -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;
