# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Flags for local go
GOFLAGS=-a

# The Executables
EXECUTABLES=mwsapid mwsquery temaquery elasticsync temasearchquery

all: test build
.PHONY : all build test clean run deps testdeps

build: $(EXECUTABLES)

$(EXECUTABLES): %: deps
	CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) ./cmd/$@

test: testdeps
	CGO_ENABLED=0 $(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(OUT_DIR)
run: build-local
	./$(BINARY_NAME)
deps:
	$(GOGET) -v ./...
testdeps:
	$(GOGET) -v -t ./...