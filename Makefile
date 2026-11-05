
MODULE  := github.com/CodeIter/webdav-server-golang
BINARY  := webdav-server
BIN_DIR := bin
CMD_DIR := ./cmd

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(MODULE)/version.Version=$(VERSION) \
	-X $(MODULE)/version.Commit=$(COMMIT) \
	-X $(MODULE)/version.Date=$(DATE)

# Extra arguments for `make run`, for example:
#   make run ARGS="-s ./data -l 8080"
ARGS ?=

.PHONY: all fmt vet test build run clean

all: fmt vet test build

fmt:
	gofmt -s -w .
	go fmt ./...

vet:
	go vet ./...

test:
	go test -race ./...

build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD_DIR)

run:
	go run -ldflags "$(LDFLAGS)" $(CMD_DIR) $(ARGS)

clean:
	rm -rf $(BIN_DIR)

