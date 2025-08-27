#!/bin/bash
# Ultra-minimal build script to reduce .so file size to under 3MB

echo "🔧 Building ultra-minimal Alipay plugin..."

# Set environment variables for minimal build
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

# Ultra-aggressive optimization flags
echo "📦 Using ultra-minimal build flags..."

go build -buildmode=plugin \
  -ldflags="-s -w -extldflags=-static" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -o alipay_channel_ultra_minimal.so \
  alipay_channel.go

# Check file size
if [ -f "alipay_channel_ultra_minimal.so" ]; then
    echo "✅ Build successful!"
    echo "📁 File: alipay_channel_ultra_minimal.so"
    echo "📏 Size: $(du -h alipay_channel_ultra_minimal.so | cut -f1)"
    
    # Show size in bytes for comparison
    SIZE_BYTES=$(stat -c%s alipay_channel_ultra_minimal.so 2>/dev/null || stat -f%z alipay_channel_ultra_minimal.so 2>/dev/null || echo "unknown")
    echo "📏 Size (bytes): $SIZE_BYTES"
    
    if [ "$SIZE_BYTES" != "unknown" ] && [ "$SIZE_BYTES" -lt 3145728 ]; then
        echo "🎯 SUCCESS: File size is under 3MB!"
    else
        echo "⚠️  File size is still over 3MB. Trying more aggressive optimization..."
        
        # Try even more aggressive optimization
        go build -buildmode=plugin \
          -ldflags="-s -w -extldflags=-static -H linux -E" \
          -gcflags="-l=4 -B -N -shared" \
          -trimpath \
          -o alipay_channel_ultra_minimal_v2.so \
          alipay_channel.go
        
        if [ -f "alipay_channel_ultra_minimal_v2.so" ]; then
            SIZE_V2=$(stat -c%s alipay_channel_ultra_minimal_v2.so 2>/dev/null || stat -f%z alipay_channel_ultra_minimal_v2.so 2>/dev/null || echo "unknown")
            echo "📏 V2 Size: $SIZE_V2 bytes"
        fi
    fi
else
    echo "❌ Build failed!"
fi
