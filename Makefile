# Makefile

# Define the binary name
BINARY_NAME=flowlog-processor

# Build the project
build:
	go build -o $(BINARY_NAME) ./cmd/flowlog-processor