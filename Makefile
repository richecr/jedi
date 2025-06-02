BINARY_NAME := jedi-scan

.PHONY: all build run test lint fmt vet clean help

all: build

build:
	@echo ">> Building..."
	go build

run: build
	@echo ">> Running..."
	./$(BINARY_NAME)

test:
	@echo ">> Running tests..."
	go test ./...

fmt:
	@echo ">> Formatting..."
	go fmt ./...

clean:
	@echo ">> Cleaning..."
	rm -rf bin/
