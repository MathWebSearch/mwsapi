# Go itself
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Flags for local go
GOFLAGS=-a

# The Executables
EXECUTABLES=mwsapid mwsquery elasticquery elasticsync temaquery

all: test build
.PHONY : all build test integrationtest integrationpull clean run deps testdeps

build: $(EXECUTABLES)

$(EXECUTABLES): %: deps
	CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) ./cmd/$@

test: testdeps
	CGO_ENABLED=0 $(GOTEST) -short -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(EXECUTABLES)
run: build-local
	./$(BINARY_NAME)
deps:
	$(GOGET) -v ./...
testdeps:
	$(GOGET) -v -t ./...

# Integration Tests
integrationdeps:
	cd test && docker-compose pull
integrationtest: testdeps
	CGO_ENABLED=0 $(GOTEST) -v ./cmd/...