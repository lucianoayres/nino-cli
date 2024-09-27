# Define variables for the binary name, build directory, and test directory
BINARY_NAME := gollama
BUILD_DIR := ./src/cmd/gollama
TEST_DIR := src

# Default target
.PHONY: all
all: test build

# Build the project
.PHONY: build
build:
	go build -C $(BUILD_DIR) -o ../../../$(BINARY_NAME)

# Update gollama version on the user local binary directory
.PHONY: 
update-local-bin:
	sudo cp ./gollama /usr/local/bin/gollama

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -C $(TEST_DIR) ./...
	@echo "Tests complete"

# Run tests with coverage
.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	go test -C $(TEST_DIR) -coverprofile=coverage.out ./...
	@echo "Generating HTML coverage report..."
	go tool -C $(TEST_DIR) cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean up the binary and coverage files
.PHONY: clean
clean:
	@echo "Cleaning up binaries and coverage reports..."
	rm -f $(BINARY_NAME)
	(cd ./$(TEST_DIR) && rm -f coverage.out coverage.html)
	@echo "Cleanup complete."

# Clean only coverage files
.PHONY: clean-coverage
clean-coverage:
	@echo "Cleaning up coverage reports..."
	(cd ./$(TEST_DIR) && rm -f coverage.out coverage.html)
	@echo "Coverage cleanup complete."

# Set the local Git config to use the custom hooks directory
.PHONY: setup-git-hooks
setup-git-hooks:
	@echo "Setting up git hooks..."
	git config core.hooksPath git-hooks
	chmod +x git-hooks/*
	@echo "Git hooks have been configured successfully."