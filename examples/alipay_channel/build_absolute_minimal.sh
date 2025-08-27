#!/bin/bash
# Build absolute minimal Alipay plugin with multiple aggressive size reduction approaches

echo "🔧 Building absolute minimal Alipay plugin with multiple approaches..."

# Set environment variables
export GOOS=linux
export GOARCH=amd64

echo "📦 Trying multiple build approaches..."

# Approach 1: Standard minimal build
echo "=== Approach 1: Standard Minimal ==="
export CGO_ENABLED=1
go build -buildmode=plugin \
  -ldflags="-s -w" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -o alipay_channel_approach1.so \
  alipay_channel_absolute_minimal.go

if [ -f "alipay_channel_approach1.so" ]; then
    SIZE1=$(stat -c%s alipay_channel_approach1.so 2>/dev/null || stat -f%z alipay_channel_approach1.so 2>/dev/null || echo "unknown")
    echo "Size: $SIZE1 bytes ($(echo "scale=2; $SIZE1/1024/1024" | bc -l 2>/dev/null || echo "unknown") MB)"
else
    echo "Build failed"
fi

# Approach 2: No CGO
echo "=== Approach 2: No CGO ==="
export CGO_ENABLED=0
go build -buildmode=plugin \
  -ldflags="-s -w" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -o alipay_channel_approach2.so \
  alipay_channel_absolute_minimal.go

if [ -f "alipay_channel_approach2.so" ]; then
    SIZE2=$(stat -c%s alipay_channel_approach2.so 2>/dev/null || stat -f%z alipay_channel_approach2.so 2>/dev/null || echo "unknown")
    echo "Size: $SIZE2 bytes ($(echo "scale=2; $SIZE2/1024/1024" | bc -l 2>/dev/null || echo "unknown") MB)"
else
    echo "Build failed"
fi

# Approach 3: Ultra-aggressive stripping
echo "=== Approach 3: Ultra-Aggressive Stripping ==="
export CGO_ENABLED=1
go build -buildmode=plugin \
  -ldflags="-s -w -H linux -E" \
  -gcflags="-l=4 -B -N -shared" \
  -trimpath \
  -a \
  -installsuffix stripped \
  -o alipay_channel_approach3.so \
  alipay_channel_absolute_minimal.go

if [ -f "alipay_channel_approach3.so" ]; then
    SIZE3=$(stat -c%s alipay_channel_approach3.so 2>/dev/null || stat -f%z alipay_channel_approach3.so 2>/dev/null || echo "unknown")
    echo "Size: $SIZE3 bytes ($(echo "scale=2; $SIZE3/1024/1024" | bc -l 2>/dev/null || echo "unknown") MB)"
else
    echo "Build failed"
fi

# Approach 4: Try building as static library instead of plugin
echo "=== Approach 4: Static Library (Not Plugin) ==="
export CGO_ENABLED=0
go build \
  -ldflags="-s -w" \
  -gcflags="-l=4 -B -N" \
  -trimpath \
  -o alipay_channel_static \
  alipay_channel_absolute_minimal.go

if [ -f "alipay_channel_static" ]; then
    SIZE4=$(stat -c%s alipay_channel_static 2>/dev/null || stat -f%z alipay_channel_static 2>/dev/null || echo "unknown")
    echo "Size: $SIZE4 bytes ($(echo "scale=2; $SIZE4/1024/1024" | bc -l 2>/dev/null || echo "unknown") MB)"
    echo "Note: This is a static binary, not a plugin"
else
    echo "Build failed"
fi

echo ""
echo "=== Summary ==="
echo "Approach 1 (Standard): $SIZE1 bytes"
echo "Approach 2 (No CGO): $SIZE2 bytes"
echo "Approach 3 (Ultra-Stripped): $SIZE3 bytes"
echo "Approach 4 (Static): $SIZE4 bytes"

# Find the smallest successful build
echo ""
echo "=== Smallest Successful Build ==="
if [ "$SIZE1" != "unknown" ] && [ "$SIZE1" -lt 3145728 ]; then
    echo "✅ Approach 1: Under 3MB at $SIZE1 bytes"
fi
if [ "$SIZE2" != "unknown" ] && [ "$SIZE2" -lt 3145728 ]; then
    echo "✅ Approach 2: Under 3MB at $SIZE2 bytes"
fi
if [ "$SIZE3" != "unknown" ] && [ "$SIZE3" -lt 3145728 ]; then
    echo "✅ Approach 3: Under 3MB at $SIZE3 bytes"
fi
if [ "$SIZE4" != "unknown" ] && [ "$SIZE4" -lt 3145728 ]; then
    echo "✅ Approach 4: Under 3MB at $SIZE4 bytes"
fi
