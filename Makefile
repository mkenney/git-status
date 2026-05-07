GOCMD=go
GOPKG=./pkg/
GOBIN=bin
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=git-status

all: clean build
build: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-linux-armv7
test:
	GO111MODULE=off $(GOCMD) test -v $(GOPKG)
clean:
	$(GOCLEAN)
	rm -f $(GOBIN)/$(BINARY_NAME)-darwin-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-darwin-arm64
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-amd64
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-arm64
	rm -f $(GOBIN)/$(BINARY_NAME)-linux-armv7
run:
	GO111MODULE=off $(GOBUILD) -o $(BINARY_NAME) -v $(GOPKG)
	./$(BINARY_NAME)

# Target architectures
build-darwin-amd64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-darwin-amd64 -v $(GOPKG)
build-darwin-arm64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-darwin-arm64 -v $(GOPKG)
build-linux-amd64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-amd64  -v $(GOPKG)
build-linux-arm64:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 GOARM=  $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-arm64  -v $(GOPKG)
build-linux-armv7:
	GO111MODULE=off CGO_ENABLED=0 GOOS=linux  GOARCH=arm   GOARM=7 $(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)-linux-armv7  -v $(GOPKG)
