#!/bin/bash
# Build minimal Alipay plugin to achieve under 3MB file size

echo "🔧 Building minimal Alipay plugin for size optimization..."

# Set environment variables for minimal build
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

echo "📦 Using minimal build flags..."

# Build minimal version with aggressive optimization
go build -buildmode=plugin \
  -ldflags="-s -w -extldflags=-static" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -o alipay_channel_minimal.so \
  alipay_channel_minimal.go

# Check file size
if [ -f "alipay_channel_minimal.so" ]; then
    echo "✅ Build successful!"
    echo "📁 File: alipay_channel_minimal.so"
    echo "📏 Size: $(du -h alipay_channel_minimal.so | cut -f1)"
    
    # Show size in bytes for comparison
    SIZE_BYTES=$(stat -c%s alipay_channel_minimal.so 2>/dev/null || stat -f%z alipay_channel_minimal.so 2>/dev/null || echo "unknown")
    echo "📏 Size (bytes): $SIZE_BYTES"
    
    if [ "$SIZE_BYTES" != "unknown" ] && [ "$SIZE_BYTES" -lt 3145728 ]; then
        echo "🎯 SUCCESS: File size is under 3MB!"
        echo "💡 This minimal version removed:"
        echo "   • Heavy crypto packages (TLS, RSA, ECDSA)"
        echo "   • HTTP/2 support"
        echo "   • Compression libraries"
        echo "   • Unused encoding packages"
    else
        echo "⚠️  File size is still over 3MB."
        echo "💡 Try building with CGO_ENABLED=0 for even smaller size"
        
        # Try without CGO
        export CGO_ENABLED=0
        go build -buildmode=plugin \
          -ldflags="-s -w" \
          -gcflags="-l=4 -B -N" \
          -trimpath \
          -o alipay_channel_minimal_nocgo.so \
          alipay_channel_minimal.go
        
        if [ -f "alipay_channel_minimal_nocgo.so" ]; then
            SIZE_NOCGO=$(stat -c%s alipay_channel_minimal_nocgo.so 2>/dev/null || stat -f%z alipay_channel_minimal_nocgo.so 2>/dev/null || echo "unknown")
            echo "📏 No-CGO Size: $SIZE_NOCGO bytes"
        fi
    fi
else
    echo "❌ Build failed!"
fi
