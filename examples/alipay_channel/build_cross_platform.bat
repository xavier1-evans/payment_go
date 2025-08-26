@echo off
REM Cross-Platform Build Script for Windows Users
REM This script builds Linux and macOS plugins from Windows

echo 🚀 Building Cross-Platform Plugins from Windows
echo ================================================

REM Create output directory
if not exist output mkdir output

echo.
echo 📦 Building for Linux (x86_64)...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=1
go build -buildmode=plugin -ldflags="-s -w -extldflags=-static" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux.so .

echo.
echo 📦 Building for Linux (ARM64)...
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=1
go build -buildmode=plugin -ldflags="-s -w -extldflags=-static" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_arm64.so .

echo.
echo 📦 Building for macOS (x86_64)...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=1
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin.so .

echo.
echo 📦 Building for macOS (ARM64)...
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=1
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin_arm64.so .

echo.
echo 📊 Build Results:
echo =================

REM Show file sizes
for %%f in (output\*.so) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Cross-platform build completed!
echo 🎯 Plugins ready for Linux and macOS deployment
echo.
echo 💡 Note: These .so files cannot run on Windows
echo    They are designed for Linux/macOS servers
