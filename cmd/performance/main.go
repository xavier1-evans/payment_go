package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"payment_go/pkg/interfaces"
	"payment_go/pkg/plugin"
)

// PerformanceTestResult holds the results of a performance test
type PerformanceTestResult struct {
	TotalRequests    int64
	SuccessfulRequests int64
	FailedRequests   int64
	TotalDuration    time.Duration
	AverageLatency   time.Duration
	MinLatency       time.Duration
	MaxLatency       time.Duration
	RequestsPerSecond float64
	Concurrency      int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/performance/main.go <plugin_path>")
		fmt.Println("Example: go run cmd/performance/main.go examples/mock_channel/output/mock_channel.so")
		os.Exit(1)
	}

	pluginPath := os.Args[1]
	channelID := "mock_channel"

	fmt.Printf("🚀 Payment Gateway Performance Test\n")
	fmt.Printf("====================================\n\n")

	// Load the plugin
	loader := plugin.NewPluginLoader()
	err := loader.LoadPlugin(pluginPath, channelID)
	if err != nil {
		log.Fatalf("❌ Failed to load plugin: %v", err)
	}

	paymentChannel, err := loader.GetPlugin(channelID)
	if err != nil {
		log.Fatalf("❌ Failed to get plugin instance: %v", err)
	}

	// Initialize with minimal delay for performance testing
	config := map[string]interface{}{
		"mock_delay_ms": 1,   // Minimal delay for performance testing
		"success_rate":  1.0, // 100% success rate for consistent testing
	}

	err = paymentChannel.Initialize(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize plugin: %v", err)
	}

	// Performance test configurations
	testConfigs := []struct {
		concurrency int
		totalRequests int
		description string
	}{
		{1, 100, "Single-threaded (100 requests)"},
		{10, 1000, "Low concurrency (10 workers, 1000 requests)"},
		{50, 5000, "Medium concurrency (50 workers, 5000 requests)"},
		{100, 10000, "High concurrency (100 workers, 10000 requests)"},
		{200, 20000, "Very high concurrency (200 workers, 20000 requests)"},
	}

	for _, testConfig := range testConfigs {
		fmt.Printf("🧪 Running Test: %s\n", testConfig.description)
		fmt.Printf("   Concurrency: %d workers\n", testConfig.concurrency)
		fmt.Printf("   Total Requests: %d\n", testConfig.totalRequests)
		fmt.Printf("   Target: Collection Order (代收下单) - the busiest operation\n\n")

		result := runPerformanceTest(paymentChannel, testConfig.concurrency, testConfig.totalRequests)
		printPerformanceResults(result)

		fmt.Printf("\n" + strings.Repeat("-", 60) + "\n\n")
	}
}

// runPerformanceTest executes a performance test with the given parameters
func runPerformanceTest(paymentChannel interfaces.Plugin, concurrency, totalRequests int) *PerformanceTestResult {
	var (
		successCount int64
		failedCount int64
		totalLatency int64
		minLatency   int64 = 1<<63 - 1
		maxLatency   int64
		startTime    = time.Now()
		latencyMutex sync.Mutex
	)

	// Create a worker pool
	requestChan := make(chan int, totalRequests)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for requestID := range requestChan {
				start := time.Now()
				
				// Create test request
				req := &interfaces.CollectOrderRequest{
					BaseRequest: interfaces.BaseRequest{
						MerchantID:  "PERF_TEST",
						ChannelID:   "mock_channel",
						RequestID:   fmt.Sprintf("PERF_%d_%d", workerID, requestID),
						Timestamp:   time.Now(),
						ExtraParams: map[string]string{"performance_test": "true"},
					},
					OrderID:     fmt.Sprintf("PERF_ORDER_%d_%d", workerID, requestID),
					Amount:      100.00,
					Currency:    "CNY",
					Description: "Performance test payment",
					ReturnURL:   "https://example.com/return",
					NotifyURL:   "https://example.com/notify",
					CustomerInfo: &interfaces.CustomerInfo{
						Name:     "Performance Tester",
						Email:    "perf@example.com",
						Phone:    "+86-138-0000-0000",
						IDNumber: "110101199001011234",
					},
				}

				// Execute the request
				_, err := paymentChannel.CollectOrder(context.Background(), req)
				
				latency := time.Since(start).Nanoseconds()
				
				// Update counters atomically
				if err != nil {
					atomic.AddInt64(&failedCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}

				// Update latency statistics (need mutex for min/max)
				latencyMutex.Lock()
				atomic.AddInt64(&totalLatency, latency)
				if latency < minLatency {
					minLatency = latency
				}
				if latency > maxLatency {
					maxLatency = latency
				}
				latencyMutex.Unlock()
			}
		}(i)
	}

	// Send requests to workers
	for i := 0; i < totalRequests; i++ {
		requestChan <- i
	}
	close(requestChan)

	// Wait for all workers to complete
	wg.Wait()
	totalDuration := time.Since(startTime)

	// Calculate results
	var avgLatency int64
	if successCount > 0 {
		avgLatency = totalLatency / successCount
	}

	requestsPerSecond := float64(totalRequests) / totalDuration.Seconds()

	return &PerformanceTestResult{
		TotalRequests:     int64(totalRequests),
		SuccessfulRequests: successCount,
		FailedRequests:    failedCount,
		TotalDuration:     totalDuration,
		AverageLatency:    time.Duration(avgLatency),
		MinLatency:        time.Duration(minLatency),
		MaxLatency:        time.Duration(maxLatency),
		RequestsPerSecond: requestsPerSecond,
		Concurrency:       concurrency,
	}
}

// printPerformanceResults displays the performance test results
func printPerformanceResults(result *PerformanceTestResult) {
	fmt.Printf("📊 Performance Test Results:\n")
	fmt.Printf("   Total Requests: %d\n", result.TotalRequests)
	fmt.Printf("   Successful: %d\n", result.SuccessfulRequests)
	fmt.Printf("   Failed: %d\n", result.FailedRequests)
	fmt.Printf("   Success Rate: %.2f%%\n", float64(result.SuccessfulRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("   Total Duration: %s\n", result.TotalDuration)
	fmt.Printf("   Average Latency: %s\n", result.AverageLatency)
	fmt.Printf("   Min Latency: %s\n", result.MinLatency)
	fmt.Printf("   Max Latency: %s\n", result.MaxLatency)
	fmt.Printf("   Requests/Second: %.2f\n", result.RequestsPerSecond)
	fmt.Printf("   Concurrency: %d workers\n", result.Concurrency)

	// Performance analysis
	fmt.Printf("\n💡 Performance Analysis:\n")
	if result.RequestsPerSecond > 1000 {
		fmt.Printf("   ✅ Excellent performance: >1000 req/s\n")
	} else if result.RequestsPerSecond > 500 {
		fmt.Printf("   ✅ Good performance: >500 req/s\n")
	} else if result.RequestsPerSecond > 100 {
		fmt.Printf("   ⚠️  Acceptable performance: >100 req/s\n")
	} else {
		fmt.Printf("   ❌ Poor performance: <100 req/s\n")
	}

	if result.AverageLatency < 10*time.Millisecond {
		fmt.Printf("   ✅ Excellent latency: <10ms\n")
	} else if result.AverageLatency < 50*time.Millisecond {
		fmt.Printf("   ✅ Good latency: <50ms\n")
	} else if result.AverageLatency < 100*time.Millisecond {
		fmt.Printf("   ⚠️  Acceptable latency: <100ms\n")
	} else {
		fmt.Printf("   ❌ High latency: >100ms\n")
	}

	// Scalability analysis
	if result.Concurrency > 1 {
		efficiency := float64(result.RequestsPerSecond) / float64(result.Concurrency)
		fmt.Printf("   📈 Efficiency per worker: %.2f req/s\n", efficiency)
		
		if efficiency > 50 {
			fmt.Printf("   ✅ Excellent scalability\n")
		} else if efficiency > 20 {
			fmt.Printf("   ✅ Good scalability\n")
		} else {
			fmt.Printf("   ⚠️  Limited scalability\n")
		}
	}
}

// strings.Repeat is a simple implementation for the performance test
func strings.Repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
