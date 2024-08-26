# Define variables
BINARY_NAME=flowlog-processor
GO_FILES=$(shell find . -name '*.go')
GO_TEST_FILES=$(shell find . -name '*_test.go')

# Build the project
build:
	go build -o $(BINARY_NAME) $(GO_FILES)

# Run tests
test:
	go test -v $(GO_TEST_FILES)
