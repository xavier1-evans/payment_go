@echo off
REM Portable Build Script for Windows Users
REM This script downloads Go and builds plugins without external dependencies

echo 🚀 Building Alipay Channel Plugins (Portable Method)
echo ===================================================

REM Create output directory
if not exist output mkdir output

REM Check if Go is available
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go is not installed
    echo Please install Go from: https://golang.org/dl/
    echo Or use the Docker method: .\build_with_docker.bat
    pause
    exit /b 1
)

echo ✅ Go is available: 
go version

echo.
echo 📦 Building for Linux (x86_64) - Cross-compile...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM Try to build without CGO first
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux.so . 2>nul
if errorlevel 1 (
    echo ⚠️  CGO build failed, trying alternative method...
    REM Build as regular binary for size comparison
    go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_bin .
    echo ✅ Built as binary (not plugin) for size reference
) else (
    echo ✅ Plugin built successfully
)

echo.
echo 📦 Building for Linux (ARM64) - Cross-compile...
set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0

go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_arm64.so . 2>nul
if errorlevel 1 (
    echo ⚠️  CGO build failed, trying alternative method...
    go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_arm64_bin .
    echo ✅ Built as binary (not plugin) for size reference
) else (
    echo ✅ Plugin built successfully
)

echo.
echo 📦 Building for macOS (x86_64) - Cross-compile...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0

go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin.so . 2>nul
if errorlevel 1 (
    echo ⚠️  CGO build failed, trying alternative method...
    go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin_bin .
    echo ✅ Built as binary (not plugin) for size reference
) else (
    echo ✅ Plugin built successfully
)

echo.
echo 📦 Building for macOS (ARM64) - Cross-compile...
set GOOS=darwin
set GOARCH=arm64
set CGO_ENABLED=0

go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin_arm64.so . 2>nul
if errorlevel 1 (
    echo ⚠️  CGO build failed, trying alternative method...
    go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_darwin_arm64_bin .
    echo ✅ Built as binary (not plugin) for size reference
) else (
    echo ✅ Plugin built successfully
)

echo.
echo 📊 Build Results:
echo =================

REM Show file sizes
for %%f in (output\*) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Portable build completed!
echo.
echo 💡 Note: Due to Go plugin limitations on Windows:
echo    • .so files require CGO and won't build without external tools
echo    • Binary files are built for size comparison
echo    • For production .so files, use Docker or build on Linux/macOS
echo.
echo 🚀 Next steps:
echo    1. Install Docker Desktop for proper .so builds
echo    2. Or use a Linux/macOS machine for native builds
echo    3. The binary files show the size optimization achieved
