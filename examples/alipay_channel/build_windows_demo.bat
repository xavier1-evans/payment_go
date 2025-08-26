@echo off
REM Windows Demo Build Script
REM Shows what's possible on Windows vs what requires Docker/Linux

echo 🎯 Windows Build Capabilities Demo
echo ==================================

REM Create output directory
if not exist output mkdir output

echo.
echo ✅ What WORKS on Windows:
echo =========================

echo 📦 Building Windows executable (no optimization)...
go build -o output\alipay_channel_windows.exe .

echo 📦 Building Windows executable (with optimization)...
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_windows_opt.exe .

echo.
echo ✅ What WORKS on Windows (Cross-compile):
echo =========================================

echo 📦 Building Linux binary (not plugin)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_binary .

echo 📦 Building macOS binary (not plugin)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_macos_binary .

echo.
echo ❌ What DOESN'T WORK on Windows:
echo =================================

echo 📦 Attempting to build Linux plugin...
set GOOS=linux
set GOARCH=amd64
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux.so . 2>nul
if errorlevel 1 (
    echo    ❌ FAILED: Go plugins not supported on Windows
    echo    💡 This is why you need Docker or Linux/macOS
)

echo.
echo 📊 Build Results (What We Can Show):
echo ====================================

REM Show file sizes
for %%f in (output\*) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo 📈 Size Optimization Results:
echo =============================

REM Calculate size reductions
for %%f in (output\*windows.exe) do set "no_opt_size=%%~zf"
for %%f in (output\*windows_opt.exe) do set "opt_size=%%~zf"

if defined no_opt_size if defined opt_size (
    set /a "reduction=(%no_opt_size%-%opt_size%)*100/%no_opt_size%"
    echo    Windows optimization: %reduction%%% smaller
)

echo.
echo ✅ Windows demo completed!
echo.
echo 💡 What this proves:
echo    • Size optimization techniques work on Windows
echo    • Cross-platform binaries can be built
echo    • The same optimization flags apply to .so files
echo.
echo 🚨 What you CANNOT do on Windows:
echo    • Build Go plugins (.so files)
echo    • Use -buildmode=plugin
echo    • Create production-ready .so files
echo.
echo 🚀 To get real .so files, you need:
echo    1. Docker Desktop (recommended for Windows)
echo    2. WSL2 (Windows Subsystem for Linux)
echo    3. A Linux/macOS machine
echo    4. Cloud Linux environment
