.PHONY: build test clean install lint fmt help

# Binary name
BINARY_NAME=ahrefs
BINARY_PATH=./$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build variables
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Binary built: $(BINARY_PATH)"

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.out

clean: ## Remove build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

install: ## Install the binary to GOPATH/bin
	$(GOCMD) install $(LDFLAGS) .

fmt: ## Format code
	$(GOFMT) ./...

lint: ## Run linters (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run

tidy: ## Tidy go modules
	$(GOMOD) tidy

run: build ## Build and run with example command
	./$(BINARY_NAME) --help

demo: build ## Run demo commands
	@echo "=== Ahrefs CLI Demo ==="
	@echo "\n1. Show help:"
	./$(BINARY_NAME) --help
	@echo "\n2. List commands:"
	./$(BINARY_NAME) --list-commands | head -20
	@echo "\n3. Site Explorer help:"
	./$(BINARY_NAME) site-explorer --help
	@echo "\n4. Domain rating help:"
	./$(BINARY_NAME) site-explorer domain-rating --help
	@echo "\n5. Dry run example:"
	./$(BINARY_NAME) site-explorer domain-rating --target example.com --api-key test --dry-run

deps: ## Download dependencies
	$(GOGET) -v ./...
	$(GOMOD) tidy

all: clean fmt test build ## Run fmt, test, and build

.DEFAULT_GOAL := help
