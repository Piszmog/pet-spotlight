# Parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test

# Environments
ARCH=amd64
ENV_LINUX=linux
ENV_MAC=darwin
ENV_WINDOWS=windows

# Binary names
BINARY_NAME=pet-spotlight
BINARY_LINUX=$(BINARY_NAME)-linux
BINARY_MAC=$(BINARY_NAME)-mac
BINARY_WINDOWS=$(BINARY_NAME).exe

# Commands
all: clean vet test build
clean:
    $(GO_CLEAN)
vet:
    $(GO_CMD) vet
test:
    $(GO_TEST) -v ./...
build:
    set GOARCH=$(ARCH)
    set GOOS=$(ENV_LINUX)
    $(GO_BUILD) -o $(BINARY_LINUX)
    set GOOS=$(ENV_MAC)
    $(GO_BUILD) -o $(BINARY_MAC)
    set GOOS=$(ENV_WINDOWS)
    $(GO_BUILD) -o $(BINARY_WINDOWS)