.PHONY: build test lint run clean release install uninstall

build:
	@echo "Building binaries..."
	@go build -o bin/vpsm ./cmd/vpsm
	@go build -o bin/vpsmd ./cmd/vpsmd
	@go build -o bin/vpsm-api ./cmd/vpsm-api
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

install: build
	@echo "Installing binaries..."
	@mkdir -p /usr/local/bin
	@rm -f /usr/local/bin/vpsm
	@rm -f /usr/local/bin/vpcm
	@rm -f /usr/local/bin/vpsmd
	@rm -f /usr/local/bin/vpsm-api
	@cp bin/vpsm /usr/local/bin/vpsm
	@cp bin/vpsm /usr/local/bin/vpcm
	@cp bin/vpsmd /usr/local/bin/vpsmd
	@cp bin/vpsm-api /usr/local/bin/vpsm-api
	@echo "Configuring shell wrappers..."
	@if [ -f $$HOME/.zshrc ] && ! grep -q "VPSM ssh wrapper override" $$HOME/.zshrc; then \
		cat scripts/shell_wrapper.sh >> $$HOME/.zshrc; \
		echo "Configured ~/.zshrc"; \
	fi
	@if [ -f $$HOME/.bashrc ] && ! grep -q "VPSM ssh wrapper override" $$HOME/.bashrc; then \
		cat scripts/shell_wrapper.sh >> $$HOME/.bashrc; \
		echo "Configured ~/.bashrc"; \
	fi
	@echo "Installed vpsm, vpcm, vpsmd, and vpsm-api to /usr/local/bin"

uninstall:
	@echo "Uninstalling binaries..."
	@rm -f /usr/local/bin/vpsm
	@rm -f /usr/local/bin/vpcm
	@rm -f /usr/local/bin/vpsmd
	@rm -f /usr/local/bin/vpsm-api
	@echo "Removing shell wrappers..."
	@if [ -f $$HOME/.zshrc ] && grep -q "VPSM ssh wrapper override" $$HOME/.zshrc; then \
		awk ' \
			/# VPSM ssh wrapper override/ { skip = 1; next } \
			skip && /^}/ { skip = 0; next } \
			skip { next } \
			{ print } \
		' $$HOME/.zshrc > $$HOME/.zshrc.tmp && mv $$HOME/.zshrc.tmp $$HOME/.zshrc; \
		echo "Removed wrapper from ~/.zshrc"; \
	fi
	@if [ -f $$HOME/.bashrc ] && grep -q "VPSM ssh wrapper override" $$HOME/.bashrc; then \
		awk ' \
			/# VPSM ssh wrapper override/ { skip = 1; next } \
			skip && /^}/ { skip = 0; next } \
			skip { next } \
			{ print } \
		' $$HOME/.bashrc > $$HOME/.bashrc.tmp && mv $$HOME/.bashrc.tmp $$HOME/.bashrc; \
		echo "Removed wrapper from ~/.bashrc"; \
	fi
	@echo "Cleaned build artifacts..."
	@rm -rf bin/
	@echo "VPSM uninstalled successfully."

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf bin/

release: build
	@echo "Release prepared successfully."
