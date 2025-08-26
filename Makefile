# Payment Gateway Plugin Framework Makefile
# Provides common commands for building, testing, and managing the framework

.PHONY: help build clean test demo performance mock-plugin all

# Default target
help:
	@echo "Payment Gateway Plugin Framework"
	@echo "================================"
	@echo ""
	@echo "Available commands:"
	@echo "  help         - Show this help message"
	@echo "  build        - Build all components"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run all tests"
	@echo "  demo         - Build and run demo application"
	@echo "  performance  - Build and run performance tests"
	@echo "  mock-plugin  - Build mock payment channel plugin"
	@echo "  all          - Build everything and run tests"
	@echo ""

# Build all components
build: mock-plugin
	@echo "✅ All components built successfully"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf examples/mock_channel/output/
	rm -rf examples/mock_channel/build/
	@echo "✅ Clean completed"

# Run all tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...
	@echo "✅ Tests completed"

# Build and run demo application
demo: mock-plugin
	@echo "🚀 Running demo application..."
	@if [ ! -f "examples/mock_channel/output/mock_channel.so" ]; then \
		echo "❌ Mock plugin not found. Run 'make mock-plugin' first."; \
		exit 1; \
	fi
	go run cmd/demo/main.go examples/mock_channel/output/mock_channel.so

# Build and run performance tests
performance: mock-plugin
	@echo "📊 Running performance tests..."
	@if [ ! -f "examples/mock_channel/output/mock_channel.so" ]; then \
		echo "❌ Mock plugin not found. Run 'make mock-plugin' first."; \
		exit 1; \
	fi
	go run cmd/performance/main.go examples/mock_channel/output/mock_channel.so

# Build mock payment channel plugin
mock-plugin:
	@echo "🔌 Building mock payment channel plugin..."
	@cd examples/mock_channel && ./build.sh
	@echo "✅ Mock plugin built successfully"

# Build everything and run tests
all: clean build test
	@echo "🎉 All tasks completed successfully!"

# Platform-specific builds
build-linux:
	@echo "🐧 Building for Linux..."
	cd examples/mock_channel && GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o output/mock_channel_linux_amd64.so .

build-windows:
	@echo "🪟 Building for Windows..."
	cd examples/mock_channel && GOOS=windows GOARCH=amd64 go build -buildmode=plugin -o output/mock_channel_windows_amd64.dll .

build-macos:
	@echo "🍎 Building for macOS..."
	cd examples/mock_channel && GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -o output/mock_channel_darwin_amd64.so .

# Cross-platform build
build-all-platforms: build-linux build-windows build-macos
	@echo "✅ Cross-platform builds completed"

# Development helpers
dev-setup:
	@echo "🔧 Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✅ Development setup completed"

# Check code quality
lint:
	@echo "🔍 Running linters..."
	gofmt -d .
	golint ./...
	govet ./...
	@echo "✅ Linting completed"

# Format code
format:
	@echo "🎨 Formatting code..."
	gofmt -w .
	@echo "✅ Code formatting completed"

# Show project structure
tree:
	@echo "📁 Project structure:"
	@tree -I 'vendor|.git|output|build' || echo "Install 'tree' command to view project structure"

# Show Go module information
mod-info:
	@echo "📦 Go module information:"
	go mod graph
	go list -m all

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	go install golang.org/x/lint/golint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ Development tools installed"

# Show help for specific target
%:
	@echo "Target '$*' not found. Run 'make help' for available commands."

