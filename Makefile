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
.PHONY : all build test integrationtest integrationpull clean deps testdeps

build: $(EXECUTABLES)

$(EXECUTABLES): %: deps
	CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) ./cmd/$@

test: testdeps
	CGO_ENABLED=0 $(GOTEST) -short -v ./...
clean: 
	$(GOCLEAN)
	rm -f integrationtest/testdata/lockfile
	rm -f $(EXECUTABLES)
deps:
	$(GOGET) -v ./...
testdeps:
	$(GOGET) -v -t ./...

# Integration Tests
integrationdeps:
	cd integrationtest/testdata && docker-compose -f docker-compose-elasticquery.yml pull
	cd integrationtest/testdata && docker-compose -f docker-compose-elasticsync.yml pull
	cd integrationtest/testdata && docker-compose -f docker-compose-mwsquery.yml pull
	cd integrationtest/testdata && docker-compose -f docker-compose-temaquery.yml pull
integrationtest: testdeps
	CGO_ENABLED=0 $(GOTEST) -v -p 1 ./cmd/...