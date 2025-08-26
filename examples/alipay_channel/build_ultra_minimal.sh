#!/bin/bash

# Ultra-Minimal Build Script for Smallest Possible .so File Sizes
# This script uses the most aggressive optimization techniques

set -e

echo "🚀 Building Alipay Channel Plugin (Ultra-Minimal Size)"
echo "======================================================"

# Create output directory
mkdir -p output

# Ultra-optimized build for Linux (x86_64)
echo "📦 Building for Linux (x86_64) - Ultra-optimized..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static -H windowsgui" \
    -gcflags="-l=4 -B -N" \
    -trimpath \
    -a \
    -installsuffix cgo \
    -o output/alipay_channel_linux_minimal.so \
    .

# Ultra-optimized build for Linux (ARM64)
echo "📦 Building for Linux (ARM64) - Ultra-optimized..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static -H windowsgui" \
    -gcflags="-l=4 -B -N" \
    -trimpath \
    -a \
    -installsuffix cgo \
    -o output/alipay_channel_linux_arm64_minimal.so \
    .

# Ultra-optimized build for macOS (x86_64)
echo "📦 Building for macOS (x86_64) - Ultra-optimized..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
    -buildmode=plugin \
    -ldflags="-s -w -H windowsgui" \
    -gcflags="-l=4 -B -N" \
    -trimpath \
    -a \
    -installsuffix cgo \
    -o output/alipay_channel_darwin_minimal.so \
    .

# Ultra-optimized build for macOS (ARM64)
echo "📦 Building for macOS (ARM64) - Ultra-optimized..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
    -buildmode=plugin \
    -ldflags="-s -w -H windowsgui" \
    -gcflags="-l=4 -B -N" \
    -trimpath \
    -a \
    -installsuffix cgo \
    -o output/alipay_channel_darwin_arm64_minimal.so \
    .

echo ""
echo "📊 Build Results (Ultra-Minimal File Sizes):"
echo "============================================="

# Show file sizes
for file in output/*minimal*; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "   $(basename "$file"): $size"
    fi
done

echo ""
echo "✅ Ultra-minimal build completed successfully!"
echo "🎯 All plugins optimized for absolute minimal size"
echo ""
echo "💡 Ultra-minimal optimization techniques used:"
echo "   • CGO disabled (CGO_ENABLED=0)"
echo "   • Strip ALL symbols (-s -w)"
echo "   • Static linking (-extldflags=-static)"
echo "   • Aggressive inlining (-l=4)"
echo "   • Disable bounds checking (-B)"
echo "   • Disable nil checks (-N)"
echo "   • Path trimming (-trimpath)"
echo "   • Force rebuild (-a)"
echo "   • Minimal install suffix"
