@echo off
REM Simple Build Script for Windows Users
REM This script builds plugins without CGO dependencies

echo 🚀 Building Alipay Channel Plugins (Simple Method)
echo ==================================================

REM Create output directory
if not exist output mkdir output

echo.
echo 📦 Building for Linux (x86_64) - No CGO...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux.so .

echo.
echo 📦 Building for Linux (ARM64) - No CGO...
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_arm64.so .

echo.
echo 📦 Building for macOS (x86_64) - No CGO...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin.so .

echo.
echo 📦 Building for macOS (ARM64) - No CGO...
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=0
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin_arm64.so .

echo.
echo 📊 Build Results:
echo =================

REM Show file sizes
for %%f in (output\*.so) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Simple build completed!
echo 🎯 Plugins ready for Linux and macOS deployment
echo.
echo 💡 Note: These .so files cannot run on Windows
echo    They are designed for Linux/macOS servers
echo    Built without CGO for maximum compatibility
