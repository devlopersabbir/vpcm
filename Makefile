.PHONY: build test lint run clean release

build:
	@echo "Building binaries..."
	@go build -o bin/vpsm cmd/vpsm/main.go
	@go build -o bin/vpsmd cmd/vpsmd/main.go
	@go build -o bin/vpsm-api cmd/vpsm-api/main.go
	@echo "Binaries compiled inside bin/"

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Checking formatting..."
	@go fmt ./...

run: build
	@echo "Running CLI doctor check..."
	@./bin/vpsm doctor

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf bin/

release: build
	@echo "Release prepared successfully."
