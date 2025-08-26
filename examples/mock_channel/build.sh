#!/bin/bash

# Build script for the Mock Payment Channel Plugin
# This script compiles the plugin into a .so file that can be loaded by the payment gateway

set -e

# Configuration
PLUGIN_NAME="mock_channel"
PLUGIN_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$PLUGIN_DIR/build"
OUTPUT_DIR="$PLUGIN_DIR/output"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Mock Payment Channel Plugin...${NC}"

# Create build directories
mkdir -p "$BUILD_DIR"
mkdir -p "$OUTPUT_DIR"

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf "$BUILD_DIR"/*
rm -rf "$OUTPUT_DIR"/*

# Build the plugin
echo -e "${YELLOW}Compiling plugin...${NC}"
cd "$PLUGIN_DIR"

# Build for Linux (most common deployment target)
echo -e "${YELLOW}Building for Linux...${NC}"
GOOS=linux GOARCH=amd64 go build \
    -buildmode=plugin \
    -o "$OUTPUT_DIR/${PLUGIN_NAME}_linux_amd64.so" \
    .

# Build for Windows (for development/testing)
echo -e "${YELLOW}Building for Windows...${NC}"
GOOS=windows GOARCH=amd64 go build \
    -buildmode=plugin \
    -o "$OUTPUT_DIR/${PLUGIN_NAME}_windows_amd64.dll" \
    .

# Build for macOS (for development/testing)
echo -e "${YELLOW}Building for macOS...${NC}"
GOOS=darwin GOARCH=amd64 go build \
    -buildmode=plugin \
    -o "$OUTPUT_DIR/${PLUGIN_NAME}_darwin_amd64.so" \
    .

# Create a symlink for the current platform
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    ln -sf "${PLUGIN_NAME}_linux_amd64.so" "$OUTPUT_DIR/${PLUGIN_NAME}.so"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    ln -sf "${PLUGIN_NAME}_darwin_amd64.so" "$OUTPUT_DIR/${PLUGIN_NAME}.so"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]]; then
    ln -sf "${PLUGIN_NAME}_windows_amd64.dll" "$OUTPUT_DIR/${PLUGIN_NAME}.dll"
fi

echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${YELLOW}Output files:${NC}"
ls -la "$OUTPUT_DIR"

echo -e "${GREEN}Plugin is ready to be loaded by the payment gateway!${NC}"
echo -e "${YELLOW}Note: Use the .so file for Linux, .dll for Windows, and .so for macOS${NC}"
