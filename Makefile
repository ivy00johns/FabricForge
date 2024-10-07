# Binary name
BINARY_NAME=fabric-cli

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean

# Main Go package
MAIN_PACKAGE=.

# Prettier command
PRETTIER=npx prettier --write .

.PHONY: all build run clean format merge

all: build

build: format
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

run:
	./$(BINARY_NAME)

dev:
	$(GORUN) $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Rebuild and run
rebuild: clean build run

# Run the script that merges JSON files in utils/merge_metadata.go
merge:
	$(GORUN) ./utils/merge_metadata/merge_metadata.go

# Update metadata .json
update_json:
	$(GORUN) ./utils/update_json/update_json.go

# Run tests
test:
	$(GOCMD) test ./...

# Format code (Go)
fmt:
	$(GOCMD) fmt ./...

# Format code (Prettier)
format:
	$(PRETTIER)

# Run linter
lint:
	golangci-lint run

# Install dependencies
deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

# Help command
help:
	@echo "Available commands:"
	@echo "  make build       - Format code and build the binary"
	@echo "  make run         - Build and run the binary"
	@echo "  make dev         - Run the application without building a binary"
	@echo "  make clean       - Remove built binary"
	@echo "  make rebuild     - Clean, build, and run"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format Go code"
	@echo "  make format      - Format code using Prettier"
	@echo "  make lint        - Run linter"
	@echo "  make deps        - Ensure dependencies are up to date"
	@echo "  make build-all   - Build for multiple platforms"
	@echo "  make merge       - Run the JSON merge script in ./utils/merge_patterns.go"
