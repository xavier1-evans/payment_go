#!/bin/bash
# Build script for online Go environments
# You can run this on: https://go.dev/play/, https://replit.com/, or any cloud IDE

echo "Building Alipay Channel Plugin..."

# Set environment variables
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

# Build with size optimization
go build -buildmode=plugin \
  -ldflags="-s -w" \
  -gcflags="-l=4" \
  -trimpath \
  -o alipay_channel_linux.so \
  alipay_channel.go

# Check file size
if [ -f "alipay_channel_linux.so" ]; then
    echo "✅ Build successful!"
    echo "📁 File: alipay_channel_linux.so"
    echo "📏 Size: $(du -h alipay_channel_linux.so | cut -f1)"
    echo "🔍 Dependencies:"
    go list -deps . | head -20
else
    echo "❌ Build failed!"
    echo "💡 Try running this on a Linux environment or cloud IDE"
fi
