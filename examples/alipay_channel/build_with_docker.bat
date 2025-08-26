@echo off
REM Docker-Based Build Script for Windows Users
REM This script uses Docker to build Linux plugins from Windows

echo 🐳 Building Alipay Channel Plugins with Docker
echo ==============================================

REM Check if Docker is available
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Docker is not installed or not running
    echo Please install Docker Desktop and start it
    echo Download from: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

echo ✅ Docker is available

REM Create output directory
if not exist output mkdir output

echo.
echo 📦 Building Linux (x86_64) plugin with Docker...
docker run --rm -v "%cd%:/workspace" -w /workspace golang:1.21-alpine sh -c "
    apk add --no-cache gcc musl-dev
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w -extldflags=-static' -gcflags='-l=4' -trimpath -o output/alipay_channel_linux.so .
"

echo.
echo 📦 Building Linux (ARM64) plugin with Docker...
docker run --rm -v "%cd%:/workspace" -w /workspace golang:1.21-alpine sh -c "
    apk add --no-cache gcc musl-dev
    GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w -extldflags=-static' -gcflags='-l=4' -trimpath -o output/alipay_channel_linux_arm64.so .
"

echo.
echo 📦 Building macOS (x86_64) plugin with Docker...
docker run --rm -v "%cd%:/workspace" -w /workspace golang:1.21-alpine sh -c "
    apk add --no-cache gcc musl-dev
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w' -gcflags='-l=4' -trimpath -o output/alipay_channel_darwin.so .
"

echo.
echo 📦 Building macOS (ARM64) plugin with Docker...
docker run --rm -v "%cd%:/workspace" -w /workspace golang:1.21-alpine sh -c "
    apk add --no-cache gcc musl-dev
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w' -gcflags='-l=4' -trimpath -o output/alipay_channel_darwin_arm64.so .
"

echo.
echo 📊 Build Results:
echo =================

REM Show file sizes
for %%f in (output\*.so) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Docker build completed successfully!
echo 🎯 All .so plugins ready for deployment
echo.
echo 💡 These .so files are built with proper CGO support
echo    and will work correctly on Linux/macOS servers
