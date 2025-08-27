#!/bin/bash
# Build ultra-stripped Alipay plugin with maximum size reduction

echo "🔧 Building ultra-stripped Alipay plugin for maximum size reduction..."

# Set environment variables for ultra-minimal build
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

echo "📦 Using ultra-stripped build flags..."

# Build with maximum stripping and no symbolic compilation
go build -buildmode=plugin \
  -ldflags="-s -w -H linux -E" \
  -gcflags="-l=4 -B -N -shared" \
  -trimpath \
  -a \
  -installsuffix stripped \
  -o alipay_channel_ultra_stripped.so \
  alipay_channel_ultra_stripped.go

# Check file size
if [ -f "alipay_channel_ultra_stripped.so" ]; then
    echo "✅ Build successful!"
    echo "📁 File: alipay_channel_ultra_stripped.so"
    echo "📏 Size: $(du -h alipay_channel_ultra_stripped.so | cut -f1)"
    
    # Show size in bytes for comparison
    SIZE_BYTES=$(stat -c%s alipay_channel_ultra_stripped.so 2>/dev/null || stat -f%z alipay_channel_ultra_stripped.so 2>/dev/null || echo "unknown")
    echo "📏 Size (bytes): $SIZE_BYTES"
    
    if [ "$SIZE_BYTES" != "unknown" ] && [ "$SIZE_BYTES" -lt 3145728 ]; then
        echo "🎯 SUCCESS: File size is under 3MB!"
        echo "💡 Ultra-stripped optimizations used:"
        echo "   • No symbolic compilation (-H linux -E)"
        echo "   • Maximum stripping (-s -w)"
        echo "   • No bounds checking (-B)"
        echo "   • No nil checks (-N)"
        echo "   • Shared mode (-shared)"
        echo "   • Force rebuild (-a)"
        echo "   • Custom install suffix"
    else
        echo "⚠️  File size is still over 3MB."
        echo "💡 Trying with CGO_ENABLED=0 for even smaller size"
        
        # Try without CGO
        export CGO_ENABLED=0
        go build -buildmode=plugin \
          -ldflags="-s -w -H linux -E" \
          -gcflags="-l=4 -B -N -shared" \
          -trimpath \
          -a \
          -installsuffix stripped_nocgo \
          -o alipay_channel_ultra_stripped_nocgo.so \
          alipay_channel_ultra_stripped.go
        
        if [ -f "alipay_channel_ultra_stripped_nocgo.so" ]; then
            SIZE_NOCGO=$(stat -c%s alipay_channel_ultra_stripped_nocgo.so 2>/dev/null || stat -f%z alipay_channel_ultra_stripped_nocgo.so 2>/dev/null || echo "unknown")
            echo "📏 No-CGO Size: $SIZE_NOCGO bytes"
        fi
    fi
else
    echo "❌ Build failed!"
fi
