# Payment Gateway Plugin Framework

A high-performance, plugin-based framework for integrating multiple payment channels into a unified payment gateway system.

## 🎯 Overview

This framework provides a standardized interface for payment channel plugins, allowing the main payment gateway to communicate with different upstream payment providers through a unified API. Each payment channel is implemented as a lightweight `.so` plugin that can be dynamically loaded at runtime.

**Current Focus**: Alipay integration with a flexible plugin architecture that can easily accommodate additional payment channels in the future.

### 🔑 Key Features

- **Unified Interface**: Standardized API for all payment operations
- **Dynamic Plugin Loading**: Load/unload payment channels without restarting the gateway
- **High Performance**: Optimized for high-concurrency collection orders (代收下单)
- **Minimal Dependencies**: Lightweight plugins with minimal external dependencies
- **Cross-Platform**: Support for Linux, Windows, and macOS
- **Comprehensive Testing**: Demo applications and performance testing tools

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Payment Gateway Core                     │
│  (Business Logic, Risk Control, Logging, etc.)            │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                 Plugin Framework Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   Plugin Loader │  │  Plugin Manager │  │   Interface │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Payment Channel Plugins                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Alipay    │  │   Mock      │  │   Custom    │  ...   │
│  │   Plugin    │  │   Plugin    │  │   Plugin    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## 📋 Standard Interface

The framework defines 6 core payment operations:

1. **CollectOrder** (代收下单) - Create payment collection orders
2. **PayoutOrder** (代付下单) - Create payout orders  
3. **CollectQuery** (代收查单) - Query collection order status
4. **PayoutQuery** (代付查单) - Query payout order status
5. **BalanceInquiry** (余额查询) - Check account balance
6. **Callback** (消息回调) - Process incoming notifications

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Git

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd payment_go

# Install dependencies
go mod tidy
```

### Build the Mock Plugin

```bash
# Navigate to the mock channel directory
cd examples/mock_channel

# Make the build script executable (Linux/macOS)
chmod +x build.sh

# Run the build script
./build.sh

# Or build manually for your platform
go build -buildmode=plugin -o output/mock_channel.so .
```

### Run the Demo

```bash
# From the project root
go run cmd/demo/main.go examples/mock_channel/output/mock_channel.so
```

### Run Performance Tests

```bash
# Test collection order performance (the busiest operation)
go run cmd/performance/main.go examples/mock_channel/output/mock_channel.so
```

### Run Alipay Demo

```bash
# Test the Alipay payment channel
go run cmd/multi_channel_demo/main.go
```

## 🔌 Creating Custom Plugins

### Plugin Structure

Each plugin must implement the `interfaces.Plugin` interface:

```go
package main

import "payment_go/pkg/interfaces"

type MyPaymentChannel struct {
    // Your implementation
}

// Required: Export this exact function name
func NewPlugin() interfaces.Plugin {
    return &MyPaymentChannel{}
}

// Required: Implement all interface methods
func (m *MyPaymentChannel) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
    // Your implementation
}

// ... implement other methods
func (m *MyPaymentChannel) GetInfo() *interfaces.PluginInfo { ... }
func (m *MyPaymentChannel) Initialize(config map[string]interface{}) error { ... }
func (m *MyPaymentChannel) ValidateConfig(config map[string]interface{}) error { ... }
```

### Plugin Requirements

1. **Export `NewPlugin()` function**: Must be named exactly `NewPlugin`
2. **Implement all methods**: All 6 payment operations + metadata methods
3. **Minimal dependencies**: Keep external packages to a minimum
4. **Error handling**: Return meaningful errors for debugging
5. **Configuration**: Support runtime configuration via `Initialize()`

### Building Your Plugin

```bash
# For Linux deployment
GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o my_plugin.so .

# For Windows development
GOOS=windows GOARCH=amd64 go build -buildmode=plugin -o my_plugin.dll .

# For macOS development  
GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -o my_plugin.so .
```

## 🧪 Testing and Development

### Demo Application

The demo app (`cmd/demo/main.go`) shows:
- Plugin loading and initialization
- All 6 payment operations
- Configuration validation
- Plugin health checks
- Usage statistics

### Performance Testing

The performance test (`cmd/performance/main.go`) focuses on:
- **Collection Order performance** (代收下单) - the busiest operation
- Concurrency testing (1 to 200 workers)
- Latency measurements
- Throughput analysis
- Scalability assessment

### Mock Plugin Features

The included mock plugin provides:
- Configurable artificial delays
- Adjustable success rates
- Realistic order simulation
- Order state progression over time
- Comprehensive error scenarios

## 📊 Performance Considerations

### Collection Orders (代收下单)

Since collection orders are the busiest operations, the framework is optimized for:

- **High Concurrency**: Support for hundreds of concurrent workers
- **Low Latency**: Minimal overhead per request
- **Efficient Memory Usage**: Lightweight request/response structures
- **Connection Pooling**: Reuse connections where possible

### Plugin Optimization Tips

1. **Minimize allocations**: Reuse objects when possible
2. **Async operations**: Use goroutines for I/O operations
3. **Connection pooling**: Maintain persistent connections to upstream providers
4. **Caching**: Cache frequently accessed data
5. **Batch operations**: Group multiple requests when possible

## 🔧 Configuration

### Plugin Configuration Schema

Each plugin should define its configuration schema in `GetInfo()`:

```go
func (p *MyPlugin) GetInfo() *interfaces.PluginInfo {
    return &interfaces.PluginInfo{
        // ... other fields
        ConfigSchema: map[string]interface{}{
            "api_key": map[string]interface{}{
                "type":        "string",
                "required":    true,
                "description": "API key for authentication",
            },
            "timeout_ms": map[string]interface{}{
                "type":        "integer",
                "default":     5000,
                "description": "Request timeout in milliseconds",
            },
        },
    }
}
```

### Runtime Configuration

```go
config := map[string]interface{}{
    "api_key":    "your_api_key_here",
    "timeout_ms": 3000,
    "base_url":   "https://api.payment-provider.com",
}

err := plugin.Initialize(config)
if err != nil {
    log.Fatalf("Plugin initialization failed: %v", err)
}
```

## 🚨 Error Handling

### Standard Error Codes

The framework uses consistent error codes across all operations:

- `SUCCESS`: Operation completed successfully
- `INVALID_REQUEST`: Request parameters are invalid
- `AUTHENTICATION_FAILED`: Authentication/authorization failed
- `INSUFFICIENT_BALANCE`: Insufficient funds for operation
- `ORDER_NOT_FOUND`: Requested order doesn't exist
- `SYSTEM_ERROR`: Internal system error
- `TIMEOUT`: Operation timed out
- `RATE_LIMITED`: Too many requests

### Error Response Structure

```go
type BaseResponse struct {
    Success   bool   `json:"success"`
    Code      string `json:"code"`
    Message   string `json:"message"`
    RequestID string `json:"request_id"`
    Timestamp time.Time `json:"timestamp"`
    ExtraData map[string]string `json:"extra_data,omitempty"`
}
```

## 🔒 Security Considerations

### Plugin Security

1. **Signature Verification**: Always verify callback signatures
2. **Input Validation**: Validate all incoming data
3. **Rate Limiting**: Implement request rate limiting
4. **Audit Logging**: Log all payment operations
5. **Secure Configuration**: Use environment variables for sensitive data

### Best Practices

- Never hardcode API keys or secrets
- Use HTTPS for all external communications
- Implement proper authentication and authorization
- Validate all callback data before processing
- Use timeouts for all external API calls

## 📁 Project Structure

```
payment_go/
├── pkg/
│   ├── interfaces/          # Core payment interfaces
│   │   └── payment_channel.go
│   └── plugin/             # Plugin loading and management
│       └── loader.go
├── examples/
│   └── mock_channel/       # Sample plugin implementation
│       ├── mock_channel.go
│       ├── build.sh
│       └── output/         # Compiled plugins
├── cmd/
│   ├── demo/               # Demo application
│   │   └── main.go
│   └── performance/        # Performance testing
│       └── main.go
├── go.mod
└── README.md
```

## 🚀 Deployment

### Production Considerations

1. **Plugin Versioning**: Implement plugin version management
2. **Hot Reloading**: Support plugin updates without downtime
3. **Monitoring**: Monitor plugin health and performance
4. **Backup**: Keep backup copies of all plugins
5. **Rollback**: Ability to rollback to previous plugin versions

### Deployment Checklist

- [ ] Compile plugins for target platform
- [ ] Verify plugin signatures and integrity
- [ ] Test plugins in staging environment
- [ ] Configure monitoring and alerting
- [ ] Document plugin configuration
- [ ] Plan rollback procedures

## 🤝 Contributing

### Development Guidelines

1. **Follow Go conventions**: Use `gofmt`, `golint`, `govet`
2. **Write tests**: Include unit tests for new features
3. **Documentation**: Update README and add code comments
4. **Error handling**: Provide meaningful error messages
5. **Performance**: Consider performance impact of changes

### Testing Your Changes

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Check code coverage
go test -cover ./...
```

## 📚 Additional Resources

### Go Plugin Documentation

- [Go Plugin Package](https://golang.org/pkg/plugin/)
- [Building Go Plugins](https://golang.org/cmd/go/#hdr-Build_modes)
- [Plugin Best Practices](https://golang.org/doc/faq#plugin)

### Payment Industry Standards

- [ISO 20022](https://www.iso20022.org/) - Financial messaging standards
- [PCI DSS](https://www.pcisecuritystandards.org/) - Payment security standards
- [EMV](https://www.emvco.com/) - Chip card standards

## 📄 License

[Add your license information here]

## 🆘 Support

For questions, issues, or contributions:

1. Check the existing issues
2. Create a new issue with detailed information
3. Provide reproduction steps and error messages
4. Include your Go version and platform information

---

**Note**: This framework is designed to be the plugin layer underneath a payment gateway. It does not include business logic, risk control, merchant onboarding, or other gateway core functionality.
