# Makefile for Alipay Channel Plugin (Size Optimized)
# Provides different build targets for various optimization levels

.PHONY: all clean build-optimized build-ultra-minimal build-all analyze

# Default target
all: build-optimized

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf output/
	@echo "✅ Clean completed"

# Build with standard optimization
build-optimized:
	@echo "📦 Building with standard optimization..."
	@chmod +x build_optimized.sh
	@./build_optimized.sh

# Build with ultra-minimal optimization
build-ultra-minimal:
	@echo "📦 Building with ultra-minimal optimization..."
	@chmod +x build_ultra_minimal.sh
	@./build_ultra_minimal.sh

# Build for all platforms with all optimization levels
build-all: build-optimized build-ultra-minimal
	@echo "🎯 All builds completed!"

# Analyze dependencies and file sizes
analyze:
	@echo "🔍 Analyzing dependencies..."
	@chmod +x analyze_dependencies.sh
	@./analyze_dependencies.sh

# Build for specific platform (Linux x86_64)
build-linux:
	@echo "📦 Building for Linux (x86_64)..."
	@mkdir -p output
	GOOS=linux GOARCH=amd64 go build \
		-buildmode=plugin \
		-ldflags="-s -w -extldflags=-static" \
		-gcflags="-l=4" \
		-trimpath \
		-o output/alipay_channel_linux.so \
		.

# Build for specific platform (Windows x86_64)
build-windows:
	@echo "📦 Building for Windows (x86_64)..."
	@mkdir -p output
	GOOS=windows GOARCH=amd64 go build \
		-buildmode=plugin \
		-ldflags="-s -w" \
		-gcflags="-l=4" \
		-trimpath \
		-o output/alipay_channel_windows.dll \
		.

# Build for specific platform (macOS x86_64)
build-macos:
	@echo "📦 Building for macOS (x86_64)..."
	@mkdir -p output
	GOOS=darwin GOARCH=amd64 go build \
		-buildmode=plugin \
		-ldflags="-s -w" \
		-gcflags="-l=4" \
		-trimpath \
		-o output/alipay_channel_darwin.so \
		.

# Show file sizes
size:
	@echo "📊 Current file sizes:"
	@if [ -d output ]; then \
		for file in output/*; do \
			if [ -f "$$file" ]; then \
				size=$$(du -h "$$file" | cut -f1); \
				echo "   $$(basename $$file): $$size"; \
			fi; \
		done; \
	else \
		echo "   No output directory found. Run 'make build-optimized' first."; \
	fi

# Help
help:
	@echo "🚀 Alipay Channel Plugin Build Targets"
	@echo "======================================"
	@echo "  all              - Build with standard optimization (default)"
	@echo "  build-optimized  - Build with standard optimization"
	@echo "  build-ultra-minimal - Build with ultra-minimal optimization"
	@echo "  build-all        - Build for all platforms and optimization levels"
	@echo "  build-linux      - Build for Linux (x86_64) only"
	@echo "  build-windows    - Build for Windows (x86_64) only"
	@echo "  build-macos      - Build for macOS (x86_64) only"
	@echo "  analyze          - Analyze dependencies and file sizes"
	@echo "  size             - Show current file sizes"
	@echo "  clean            - Clean build artifacts"
	@echo "  help             - Show this help message"
