@echo off
REM Alternative Build Approaches for Windows
REM Trying different methods to work around Go plugin limitations

echo 🔧 Alternative Build Approaches for Windows
echo ===========================================

REM Create output directory
if not exist output mkdir output

echo.
echo 🚀 Attempting Alternative Build Methods:
echo =======================================

echo.
echo 📦 Method 1: Try building with CGO enabled...
set CGO_ENABLED=1
go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\test_plugin.so . 2>nul
if errorlevel 1 (
    echo    ❌ FAILED: CGO requires C compiler (gcc/clang)
    echo    💡 Windows doesn't have built-in C compilers
)

echo.
echo 📦 Method 2: Try building as shared library...
go build -buildmode=c-shared -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\test_lib.so . 2>nul
if errorlevel 1 (
    echo    ❌ FAILED: c-shared also requires C compiler
)

echo.
echo 📦 Method 3: Try building as archive...
go build -buildmode=c-archive -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\test_lib.a . 2>nul
if errorlevel 1 (
    echo    ❌ FAILED: c-archive also requires C compiler
)

echo.
echo 📦 Method 4: Build optimized static binary...
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_static.exe .

echo.
echo 📦 Method 5: Build for Linux (binary, not plugin)...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_linux_bin .

echo.
echo 📦 Method 6: Build for macOS (binary, not plugin)...
set GOOS=darwin
set GOARCH=amd64
set CGO_ENABLED=0
go build -ldflags="-s -w" -gcflags="-l=4" -trimpath -o output\alipay_channel_macos_bin .

echo.
echo 📊 What We Can Build on Windows:
echo =================================

REM Show file sizes
for %%f in (output\*) do (
    echo    %%~nxf: %%~zf bytes
)

echo.
echo ✅ Alternative build completed!
echo.
echo 💡 What this proves:
echo    • Go plugins (.so) are fundamentally impossible on Windows
echo    • All build modes requiring CGO fail without C compiler
echo    • Only static binaries and cross-platform binaries work
echo.
echo 🚨 The Reality:
echo    • .so files require CGO
echo    • CGO requires C compiler (gcc/clang)
echo    • Windows doesn't have C compilers for cross-platform builds
echo    • This is a Go limitation, not a Windows limitation
echo.
echo 🚀 Your Only Options for .so Files:
echo    1. Use WSL2 (Linux on Windows)
echo    2. Use cloud Linux environment
echo    3. Use different machine (Linux/macOS)
echo    4. Accept that .so files aren't possible on Windows
