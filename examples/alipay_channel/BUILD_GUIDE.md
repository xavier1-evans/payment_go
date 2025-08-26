# 🏗️ Build Guide for Alipay Channel Plugins

## 🚨 **Important: Go Plugin Limitations**

**Go plugins ALWAYS require CGO, even for cross-platform builds.** This is a fundamental limitation of Go's plugin system.

## 🎯 **Build Options (Ranked by Recommendation)**

### **Option 1: Build on Linux/macOS (Best)**
```bash
# Navigate to the directory
cd examples/alipay_channel

# Make scripts executable
chmod +x build_optimized.sh
chmod +x build_ultra_minimal.sh

# Build with standard optimization
./build_optimized.sh

# Or build with ultra-minimal optimization
./build_ultra_minimal.sh
```

**Expected Results:**
- ✅ Full CGO support
- ✅ Proper .so files
- ✅ All optimization flags work
- ✅ File sizes: 500KB - 1.5MB

### **Option 2: Docker on Windows (Good)**
```bash
# Run the Docker build script
.\build_with_docker.bat
```

**Requirements:**
- Docker Desktop installed and running
- Internet connection for pulling Go image

**Expected Results:**
- ✅ Full CGO support
- ✅ Proper .so files
- ✅ All optimization flags work
- ✅ File sizes: 500KB - 1.5MB

### **Option 3: Manual Commands (Advanced)**
```bash
# For Linux (x86_64)
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static" \
    -gcflags="-l=4" \
    -trimpath \
    -o alipay_channel_linux.so \
    .

# For macOS (x86_64)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build \
    -buildmode=plugin \
    -ldflags="-s -w" \
    -gcflags="-l=4" \
    -trimpath \
    -o alipay_channel_darwin.so \
    .
```

## 🔧 **Build Flag Explanations**

### **Essential Flags for Minimal Size**
- **`-buildmode=plugin`**: Creates a Go plugin (.so file)
- **`-ldflags="-s -w"`**: Strips debug symbols (30-40% size reduction)
- **`-gcflags="-l=4"`**: Aggressive inlining
- **`-trimpath`**: Removes file paths (security + size)

### **Advanced Flags for Ultra-Minimal Size**
- **`-extldflags=-static`**: Static linking (Linux only)
- **`-gcflags="-l=4 -B -N"`**: Disable bounds/nil checks
- **`CGO_ENABLED=1`**: Required for plugins

## 📱 **Platform-Specific Builds**

### **Linux Production Build**
```bash
#!/bin/bash
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

go build \
    -buildmode=plugin \
    -ldflags="-s -w -extldflags=-static" \
    -gcflags="-l=4" \
    -trimpath \
    -tags prod \
    -o alipay_channel_linux.so \
    .
```

### **macOS Production Build**
```bash
#!/bin/bash
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=amd64

go build \
    -buildmode=plugin \
    -ldflags="-s -w" \
    -gcflags="-l=4" \
    -trimpath \
    -tags prod \
    -o alipay_channel_darwin.so \
    .
```

## 🐳 **Docker Build Commands**

### **Single Platform Build**
```bash
# Linux (x86_64)
docker run --rm -v "$(pwd):/workspace" -w /workspace golang:1.21-alpine sh -c "
    apk add --no-cache gcc musl-dev
    GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w -extldflags=-static' -gcflags='-l=4' -trimpath -o alipay_channel_linux.so .
"
```

### **Multi-Platform Build**
```bash
# Build for all platforms
for os in linux darwin; do
    for arch in amd64 arm64; do
        echo "Building for $os/$arch..."
        docker run --rm -v "$(pwd):/workspace" -w /workspace golang:1.21-alpine sh -c "
            apk add --no-cache gcc musl-dev
            GOOS=$os GOARCH=$arch CGO_ENABLED=1 go build -buildmode=plugin -ldflags='-s -w' -gcflags='-l=4' -trimpath -o alipay_channel_${os}_${arch}.so .
        "
    done
done
```

## 📊 **Expected File Sizes**

### **Size Comparison**
```
Original build:           ~2.5 MB
Standard optimization:    ~1.2 MB  (50% reduction)
Ultra-minimal:           ~800 KB   (70% reduction)
```

### **Size by Platform**
```
Linux (x86_64):          ~800 KB - 1.2 MB
Linux (ARM64):           ~800 KB - 1.2 MB
macOS (x86_64):          ~800 KB - 1.2 MB
macOS (ARM64):           ~800 KB - 1.2 MB
```

## 🚀 **Quick Start Commands**

### **For Windows Users (Docker)**
```bash
# Install Docker Desktop first, then:
cd examples/alipay_channel
.\build_with_docker.bat
```

### **For Linux/macOS Users**
```bash
cd examples/alipay_channel
chmod +x build_optimized.sh
./build_optimized.sh
```

### **For Manual Build**
```bash
cd examples/alipay_channel
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -ldflags="-s -w" -gcflags="-l=4" -trimpath -o alipay_channel.so .
```

## 🔍 **Troubleshooting**

### **Common Errors**
1. **"CGO not enabled"**: Set `CGO_ENABLED=1`
2. **"C compiler not found"**: Install gcc/clang or use Docker
3. **"Plugin not supported"**: Go plugins don't work on Windows
4. **Large file sizes**: Use `-ldflags="-s -w"` and `-gcflags="-l=4"`

### **Solutions**
- **Windows**: Use Docker or build on Linux/macOS
- **Missing compiler**: Install build tools or use Docker
- **Large files**: Apply all optimization flags

## 📝 **Next Steps**

1. **Choose your build method** (Docker recommended for Windows)
2. **Run the build script** for your platform
3. **Verify the .so files** are created and sized correctly
4. **Test the plugins** on a Linux/macOS system
5. **Deploy to production** servers

## 🎯 **Production Recommendations**

- **Use standard optimization** for most production deployments
- **Use ultra-minimal** only if you need absolute smallest size
- **Test thoroughly** before deploying to production
- **Monitor performance** after deployment
- **Keep build scripts** for consistent rebuilds
