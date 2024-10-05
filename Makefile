# Define variables for the binary name, build directory, and test directory
BINARY_NAME := nino
BUILD_DIR := ./src/cmd/nino
TEST_DIR := src
DEFAULT_MODEL := "llama3.2"
WARNING_MESSAGE := "Warning: Ollama not detected. Please run 'make install-deps' to install it."

# Default target
.PHONY: all
all: test build

# Reusable target for installing Ollama
.PHONY: install-ollama
install-ollama:
	@echo "Installing or updating Ollama..."
	@curl -fsSL https://ollama.com/install.sh | sh

# Install dependencies (Ollama - https://github.com/ollama/ollama)
.PHONY: install-deps
install-deps:
	@echo "Installing Ollama if not detected..."
	@which ollama > /dev/null || $(MAKE) install-ollama

# Build the project
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build -C $(BUILD_DIR) -o ../../../$(BINARY_NAME)
	@echo "Binary generated successfully: $(BINARY_NAME)"
	@which ollama > /dev/null || echo $(WARNING_MESSAGE)

# Update nino version on the user local binary directory
.PHONY: 
update-local-bin:
	sudo cp ./$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

# Run tests with verbose mode
.PHONY: test
test:
	@echo "Running tests..."
	go test -C $(TEST_DIR) ./...
	@echo "Tests complete"

# Run tests with verbose mode
.PHONY: test-verbose
test-verbose:
	@echo "Running tests..."
	go test -C $(TEST_DIR) ./... -v
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

# Start the Ollama server with default model
.PHONY: start-ollamma
start-ollama:
	@echo "Starting Ollama server and model..."
	ollama serve & ollama run $(DEFAULT_MODEL)
	@echo "Ollama server is running with $(DEFAULT_MODEL)."
	@echo "Open a new terminal window to run 'nino'."