# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Flags for local go
GOFLAGS=-a

# Binary paths
OUT_DIR=out
BINARY_NAME=mwsapi
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_MACOS=$(BINARY_NAME)_mac
BINARY_WINDOWS=$(BINARY_NAME)_windows.exe

all: test build

build: build-local build-linux build-macos
build-local: deps
	CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) -o $(OUT_DIR)/$(BINARY_NAME)
build-linux: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_UNIX)
build-macos: deps
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUT_DIR)/$(BINARY_MACOS)

test: testdeps
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(OUT_DIR)
run: build-local
	./$(BINARY_NAME)
deps:
	$(GOGET) -v ./...
testdeps:
	$(GOGET) -v -t ./...