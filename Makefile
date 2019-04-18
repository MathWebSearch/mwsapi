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

all: test build

build: $(OUT_DIR)/mwsapid $(OUT_DIR)/temaquery $(OUT_DIR)/elasticsync

$(OUT_DIR)/mwsapid: deps
	cd cmd/mwsapid && CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) -o ../../$(OUT_DIR)/mwsapid

$(OUT_DIR)/temaquery: deps
	cd cmd/temaquery && CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) -o ../../$(OUT_DIR)/temaquery

$(OUT_DIR)/elasticsync: deps
	cd cmd/elasticsync && CGO_ENABLED=0 $(GOBUILD) $(GOFLAGS) -o ../../$(OUT_DIR)/elasticsync

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