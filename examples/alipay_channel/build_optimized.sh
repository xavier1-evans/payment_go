#!/bin/bash

# Optimized Build Script for Minimal .so File Sizes
# This script builds the Alipay channel plugin with aggressive size optimization

set -e

echo "🚀 Building Alipay Channel Plugin (Optimized for Size)"
echo "======================================================"

# Create output directory
mkdir -p output

# Build for Linux (x86_64) - Optimized for size
echo "📦 Building for Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static" \
    -gcflags="-l=4" \
    -trimpath \
    -o output/alipay_channel_linux.so \
    .

# Build for Linux (ARM64) - Optimized for size
echo "📦 Building for Linux (ARM64)..."
GOOS=linux GOARCH=arm64 go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static" \
    -gcflags="-l=4" \
    -trimpath \
    -o output/alipay_channel_linux_arm64.so \
    .

# Build for macOS (x86_64) - Optimized for size
echo "📦 Building for macOS (x86_64)..."
GOOS=darwin GOARCH=amd64 go build \
    -buildmode=plugin \
    -ldflags="-s -w" \
    -gcflags="-l=4" \
    -trimpath \
    -o output/alipay_channel_darwin.so \
    .

# Build for macOS (ARM64) - Optimized for size
echo "📦 Building for macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 go build \
    -buildmode=plugin \
    -ldflags="-s -w" \
    -gcflags="-l=4" \
    -trimpath \
    -o output/alipay_channel_darwin_arm64.so \
    .

echo ""
echo "📊 Build Results (File Sizes):"
echo "================================"

# Show file sizes
for file in output/*; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "   $(basename "$file"): $size"
    fi
done

echo ""
echo "✅ Build completed successfully!"
echo "🎯 All plugins optimized for minimal size"
echo ""
echo "💡 Size optimization techniques used:"
echo "   • Strip debug symbols (-s -w)"
echo "   • Static linking where possible"
echo "   • Aggressive inlining (-l=4)"
echo "   • Path trimming (-trimpath)"
echo "   • Minimal dependencies"
