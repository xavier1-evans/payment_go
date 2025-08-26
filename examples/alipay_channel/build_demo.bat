@echo off
REM Demo Build Script - Shows Size Optimization
REM This builds regular Go binaries to demonstrate size reduction techniques

echo 🎯 Demo Build - Size Optimization Showcase
echo ==========================================

REM Create output directory
if not exist output mkdir output

REM Check if Go is available
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go is not installed
    echo Please install Go from: https://golang.org/dl/
    pause
    exit /b 1
)

echo ✅ Go is available: 
go version

echo.
echo 📦 Building with NO optimization...
go build -o output\alipay_channel_no_opt.exe .

echo.
echo 📦 Building with STANDARD optimization...
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_std_opt.exe .

echo.
echo 📦 Building with ULTRA optimization...
go build -ldflags="-s -w" -gcflags="-l=4 -B -N" -trimpath -o output\alipay_channel_ultra_opt.exe .

echo.
echo 📦 Building for Linux (x86_64) - Cross-compile...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_bin .

echo.
echo 📦 Building for macOS (x86_64) - Cross-compile...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_macos_bin .

echo.
echo 📊 Size Comparison Results:
echo ===========================

REM Show file sizes
for %%f in (output\*) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo 📈 Size Reduction Analysis:
echo ===========================

REM Calculate size reductions
for %%f in (output\*no_opt.exe) do set "no_opt_size=%%~zf"
for %%f in (output\*std_opt.exe) do set "std_opt_size=%%~zf"
for %%f in (output\*ultra_opt.exe) do set "ultra_opt_size=%%~zf"

if defined no_opt_size if defined std_opt_size (
    set /a "std_reduction=(%no_opt_size%-%std_opt_size%)*100/%no_opt_size%"
    echo    Standard optimization: %std_reduction%%% smaller
)

if defined no_opt_size if defined ultra_opt_size (
    set /a "ultra_reduction=(%no_opt_size%-%ultra_opt_size%)*100/%no_opt_size%"
    echo    Ultra optimization: %ultra_reduction%%% smaller
)

echo.
echo ✅ Demo build completed!
echo.
echo 💡 What this shows:
echo    • Size optimization techniques work on Windows
echo    • Cross-platform builds are possible
echo    • The same optimization flags apply to .so files
echo.
echo 🚀 For actual .so plugins:
echo    1. Install Docker Desktop (recommended)
echo    2. Use: .\build_with_docker.bat
echo    3. Or build on Linux/macOS machine
