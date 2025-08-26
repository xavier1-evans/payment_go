#!/bin/bash

# Dependency Analysis Script for Plugin Size Optimization
# This script helps identify what's contributing to file sizes

set -e

echo "🔍 Analyzing Dependencies for Size Optimization"
echo "================================================"

# Create output directory
mkdir -p output

# Build with dependency analysis
echo "📦 Building with dependency analysis..."
GOOS=linux GOARCH=amd64 go build \
    -buildmode=plugin \
    -ldflags="-s -w" \
    -gcflags="-l=4" \
    -trimpath \
    -o output/alipay_channel_analysis.so \
    .

echo ""
echo "📊 File Size Analysis:"
echo "======================"

# Get file size
size=$(du -h output/alipay_channel_analysis.so | cut -f1)
echo "   Plugin size: $size"

# Analyze binary sections
echo ""
echo "🔍 Binary Section Analysis:"
echo "==========================="

if command -v objdump &> /dev/null; then
    echo "   Analyzing with objdump..."
    objdump -h output/alipay_channel_analysis.so | grep -E "^\s*[0-9]+" | awk '{print "   " $2 ": " $3 " bytes"}'
else
    echo "   objdump not available, using alternative analysis..."
fi

# Analyze Go dependencies
echo ""
echo "📦 Go Module Dependencies:"
echo "=========================="
go list -m all | grep -v "payment_go" | head -10

# Show import graph
echo ""
echo "🔄 Import Graph (Top Level):"
echo "============================="
go list -f '{{.ImportPath}} -> {{join .Imports ", "}}' . | head -10

echo ""
echo "💡 Size Optimization Recommendations:"
echo "===================================="
echo "   • Review imported packages for alternatives"
echo "   • Consider using interfaces instead of concrete types"
echo "   • Minimize use of reflection and dynamic features"
echo "   • Use build tags to exclude unused code"
echo "   • Consider splitting large plugins into smaller ones"
