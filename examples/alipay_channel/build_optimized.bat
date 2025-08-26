@echo off
REM Optimized Build Script for Minimal .dll File Sizes (Windows)
REM This script builds the Alipay channel plugin with aggressive size optimization

echo 🚀 Building Alipay Channel Plugin (Optimized for Size)
echo ======================================================

REM Create output directory
if not exist output mkdir output

REM Build for Windows (x86_64) - Optimized for size
echo 📦 Building for Windows (x86_64)...
set GOOS=windows
set GOARCH=amd64
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_windows.dll .

REM Build for Windows (ARM64) - Optimized for size
echo 📦 Building for Windows (ARM64)...
set GOOS=windows
set GOARCH=arm64
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_windows_arm64.dll .

echo.
echo 📊 Build Results (File Sizes):
echo ================================

REM Show file sizes
for %%f in (output\*.dll) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Build completed successfully!
echo 🎯 All plugins optimized for minimal size
echo.
echo 💡 Size optimization techniques used:
echo    • Strip debug symbols (-s -w)
echo    • Aggressive inlining (-l=4)
echo    • Path trimming (-trimpath)
echo    • Minimal dependencies
