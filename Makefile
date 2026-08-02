.PHONY: build test lint run clean release install uninstall build-desktop dev-desktop build-all docker-build docker-run

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
	@INSTALL_DIR="/usr/local/bin"; \
	if [ ! -w "$$INSTALL_DIR" ] && [ -w $$(dirname "$$INSTALL_DIR") ]; then \
		SUDO=""; \
	elif [ ! -w "$$INSTALL_DIR" ]; then \
		SUDO="sudo"; \
	else \
		SUDO=""; \
	fi; \
	$$SUDO mkdir -p "$$INSTALL_DIR" && \
	$$SUDO cp bin/vpsm "$$INSTALL_DIR/vpsm" && \
	$$SUDO cp bin/vpsm "$$INSTALL_DIR/vpcm" && \
	$$SUDO cp bin/vpsmd "$$INSTALL_DIR/vpsmd" && \
	$$SUDO cp bin/vpsm-api "$$INSTALL_DIR/vpsm-api"
	@echo "Configuring shell wrappers..."
	@if [ -f $$HOME/.zshrc ] && ! grep -q "VPSM ssh wrapper override" $$HOME/.zshrc; then \
		cat scripts/shell_wrapper.sh >> $$HOME/.zshrc; \
		echo "Configured ~/.zshrc"; \
	fi
	@if [ -f $$HOME/.bashrc ] && ! grep -q "VPSM ssh wrapper override" $$HOME/.bashrc; then \
		cat scripts/shell_wrapper.sh >> $$HOME/.bashrc; \
		echo "Configured ~/.bashrc"; \
	fi
	@echo "Configuring default config..."
	@CONFIG_DIR="$$HOME/.config/vpsm"; \
	CONFIG_FILE="$$CONFIG_DIR/config.yaml"; \
	if [ ! -f "$$CONFIG_FILE" ]; then \
		mkdir -p "$$CONFIG_DIR" && \
		printf "database:\n  driver: sqlite\n  path: $$HOME/.local/share/vpsm/vpsm.db\napi:\n  enabled: true\n  host: 127.0.0.1\n  port: 8080\n  mode: local\n  global_url: http://127.0.0.1:8080\nssh:\n  timeout: 10s\nlogging:\n  level: info\n  format: pretty\n" > "$$CONFIG_FILE"; \
		echo "Auto-initialized default config in $$CONFIG_FILE"; \
	fi
	@echo "Starting REST API server daemon in background..."
	@pkill -f "vpsm-api" 2>/dev/null || true
	@nohup /usr/local/bin/vpsm-api > "$$HOME/.config/vpsm/vpsm-api.log" 2>&1 &
	@echo "REST API server is running on http://127.0.0.1:8080 (logs at ~/.config/vpsm/vpsm-api.log)"
	@echo "Installed vpsm, vpcm, vpsmd, and vpsm-api to /usr/local/bin"

uninstall:
	@bash scripts/uninstall.sh -y
	@rm -rf bin/

clean:
	@echo "Cleaning up build artifacts..."
	@rm -rf bin/
	@rm -rf app/vpsm-desktop/build/bin/

release: build
	@echo "Release prepared successfully."

build-desktop:
	@echo "Building Wails desktop application..."
	@cd app/vpsm-desktop && ~/go/bin/wails build

dev-desktop:
	@echo "Starting Wails desktop application in dev mode..."
	@cd app/vpsm-desktop && ~/go/bin/wails dev

build-all: build build-desktop
	@echo "All targets (CLI & Desktop) built successfully."

docker-build:
	@echo "Building Docker image..."
	@docker build -t vpsm-api .

docker-run:
	@echo "Running Docker compose..."
	@docker compose up

