# Go params
GOCMD=go
GOPKG=github.com/mkenney/git-status/pkg
GOBIN=bin
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test pkg/...
GODEP=dep ensure
BINARY_NAME=git-status

all: clean build
build: build-linux build-darwin
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-darwin-amd64
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	cd $(GOPKG)
	$(GODEP)
	cd -

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-amd64 -v $(GOPKG)
build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 -v $(GOPKG)
