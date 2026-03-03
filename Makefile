.PHONY: build test run clean install

# Build the binary
build:
	go build -ldflags "-X main.version=dev -X main.commit=$(shell git rev-parse --short HEAD) -X main.date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o issue2md ./cmd/issue2md

# Run tests
test:
	go test -v ./...

# Run the binary
run:
	go run ./cmd/issue2md

# Clean build artifacts
clean:
	rm -f issue2md

# Install dependencies
deps:
	go mod download
	go mod tidy

# Lint code
lint:
	go vet ./...

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o issue2md-darwin-amd64 ./cmd/issue2md
	GOOS=darwin GOARCH=arm64 go build -o issue2md-darwin-arm64 ./cmd/issue2md
	GOOS=linux GOARCH=amd64 go build -o issue2md-linux-amd64 ./cmd/issue2md
	GOOS=windows GOARCH=amd64 go build -o issue2md-windows-amd64.exe ./cmd/issue2md
