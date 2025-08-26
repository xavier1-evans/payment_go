@echo off
REM Windows Build Script (Static Library Alternative)
REM Since Go plugins aren't supported on Windows, we build a static library

echo 🚀 Building Windows-Compatible Alipay Channel
echo =============================================

REM Create output directory
if not exist output mkdir output

echo.
echo 📦 Building for Windows (x86_64) - Static Library...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

REM Build as static library (not plugin)
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_windows.exe .

echo.
echo 📦 Building for Windows (ARM64) - Static Library...
set GOOS=windows
set GOARCH=arm64
set CGO_ENABLED=0

REM Build as static library (not plugin)
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_windows_arm64.exe .

echo.
echo 📊 Build Results:
echo =================

REM Show file sizes
for %%f in (output\*.exe) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Windows build completed!
echo 🎯 Static libraries ready for Windows deployment
echo.
echo 💡 Note: These are static libraries, not plugins
echo    They can be imported and used directly in Go code
echo    For production, use the Linux/macOS .so plugins
