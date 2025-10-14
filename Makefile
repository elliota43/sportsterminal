.PHONY: build run clean install test

# Build the application
build:
	go build -o sportsterminal .

# Run the application
run: build
	./sportsterminal

# Clean build artifacts
clean:
	rm -f sportsterminal
	go clean

# Install dependencies
install:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o dist/sportsterminal-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/sportsterminal-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o dist/sportsterminal-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o dist/sportsterminal-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/sportsterminal-windows-amd64.exe .

# Development mode with hot reload (requires air)
dev:
	air

.DEFAULT_GOAL := build

