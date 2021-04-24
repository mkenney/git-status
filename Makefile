GOCMD=go
GOPKG=./pkg/
GOBIN=bin
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test pkg/...
BINARY_NAME=git-status

all: clean build
build: build-darwin-amd64 build-linux-amd64 build-linux-arm7
clean:
	$(GOCLEAN)
	rm -f $(GOBIN)/$(BINARY_NAME)-darwin-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-arm7
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Target architectures
build-darwin-amd64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 -v $(GOPKG)
build-linux-amd64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-amd64  -v $(GOPKG)
build-linux-arm7:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux  GOARCH=arm   GOARM=7 $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-arm7   -v $(GOPKG)
