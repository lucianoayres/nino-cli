# Define variables for the binary name and build directory
BINARY_NAME := gollama
BUILD_DIR := src

# Default target
.PHONY: all
all: test build

# Build the project
.PHONY: build
build:
	go build -C $(BUILD_DIR) -o ../$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -C $(BUILD_DIR) ./...

# Clean up the binary
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
