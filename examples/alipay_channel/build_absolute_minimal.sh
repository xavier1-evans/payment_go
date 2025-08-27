#!/bin/bash

echo "=== Building Absolute Minimal Alipay Channel Plugin ==="
echo "Target: alipay_channel_absolute_minimal_linux.so"
echo "Goal: Under 3MB file size"

cd "$(dirname "$0")"

# Check if source file exists
if [ ! -f "alipay_channel_absolute_minimal.go" ]; then
    echo "❌ Source file not found: alipay_channel_absolute_minimal.go"
    exit 1
fi

echo "Source file found. Building with maximum size optimization..."

# Build with ultra-aggressive size optimization
go build -buildmode=plugin \
  -ldflags="-s -w" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -a \
  -installsuffix absolute_minimal \
  -o alipay_channel_absolute_minimal_linux.so \
  alipay_channel_absolute_minimal.go

# Check build result
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    
    # Get file size
    if [ -f "alipay_channel_absolute_minimal_linux.so" ]; then
        SIZE_BYTES=$(stat -c%s alipay_channel_absolute_minimal_linux.so 2>/dev/null || stat -f%z alipay_channel_absolute_minimal_linux.so 2>/dev/null || echo "0")
        SIZE_MB=$(echo "scale=2; $SIZE_BYTES / 1048576" | bc 2>/dev/null || echo "0")
        
        echo "📁 File: alipay_channel_absolute_minimal_linux.so"
        echo "📏 Size: ${SIZE_BYTES} bytes (${SIZE_MB} MB)"
        
        # Check if under 3MB
        if [ "$SIZE_BYTES" -lt 3145728 ]; then
            echo "🎯 SUCCESS: File size is under 3MB!"
        else
            echo "⚠️  File size is over 3MB target"
        fi
        
        # Show file info
        echo ""
        echo "📋 File details:"
        ls -lh alipay_channel_absolute_minimal_linux.so
        file alipay_channel_absolute_minimal_linux.so
        
    else
        echo "❌ Build file not found"
    fi
else
    echo "❌ Build failed!"
    echo "Trying to build without optimization to see errors..."
    go build -buildmode=plugin alipay_channel_absolute_minimal.go
fi
