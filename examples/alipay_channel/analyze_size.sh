#!/bin/bash
# Analyze .so file size and identify optimization opportunities

echo "🔍 Analyzing .so file size and dependencies..."

# Build with size analysis
echo "📦 Building with size analysis..."

go build -buildmode=plugin \
  -ldflags="-s -w" \
  -gcflags="-l=4" \
  -trimpath \
  -o alipay_channel_analyze.so \
  alipay_channel.go

if [ -f "alipay_channel_analyze.so" ]; then
    echo "✅ Build successful!"
    
    # Show file size
    SIZE=$(du -h alipay_channel_analyze.so | cut -f1)
    SIZE_BYTES=$(stat -c%s alipay_channel_analyze.so 2>/dev/null || stat -f%z alipay_channel_analyze.so 2>/dev/null || echo "unknown")
    echo "📏 File size: $SIZE ($SIZE_BYTES bytes)"
    
    # Analyze dependencies that contribute to size
    echo ""
    echo "🔍 Dependencies contributing to size:"
    echo "====================================="
    
    # Show heavy dependencies
    go list -deps . | grep -E "(crypto|net|encoding|compress|vendor)" | head -20
    
    echo ""
    echo "💡 Size optimization strategies:"
    echo "================================"
    echo "1. Remove unused crypto packages"
    echo "2. Use minimal HTTP client"
    echo "3. Strip debug symbols"
    echo "4. Disable bounds checking"
    echo "5. Use static linking"
    
    # Try to identify specific large packages
    echo ""
    echo "📊 Package size analysis:"
    echo "========================="
    
    # Check if we can use go tool nm to analyze symbols
    if command -v nm >/dev/null 2>&1; then
        echo "Symbols analysis:"
        nm -D alipay_channel_analyze.so 2>/dev/null | head -10 || echo "Symbol analysis not available"
    fi
    
else
    echo "❌ Build failed!"
fi
