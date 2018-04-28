GOCMD=go
GOPKG=github.com/mkenney/git-status/pkg
GOBIN=bin
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test pkg/...
BINARY_NAME=git-status

all: clean build
build: build-linux-64 build-darwin-64
clean:
	$(GOCLEAN)
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-darwin-amd64
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Target architectures
build-linux-64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-amd64  -v $(GOPKG)
build-darwin-64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 -v $(GOPKG)
