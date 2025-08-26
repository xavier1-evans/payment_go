package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"payment_go/pkg/interfaces"
	"payment_go/pkg/plugin"
)

func main() {
	// Initialize the plugin loader
	loader := plugin.NewPluginLoader()

	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/demo/main.go <plugin_path>")
		fmt.Println("Example: go run cmd/demo/main.go examples/mock_channel/output/mock_channel.so")
		os.Exit(1)
	}

	pluginPath := os.Args[1]
	channelID := "mock_channel"

	fmt.Printf("🚀 Payment Gateway Plugin Demo\n")
	fmt.Printf("================================\n\n")

	// Load the plugin
	fmt.Printf("📦 Loading plugin from: %s\n", pluginPath)
	err := loader.LoadPlugin(pluginPath, channelID)
	if err != nil {
		log.Fatalf("❌ Failed to load plugin: %v", err)
	}
	fmt.Printf("✅ Plugin loaded successfully!\n\n")

	// Get plugin info
	info, err := loader.GetPluginInfo(channelID)
	if err != nil {
		log.Fatalf("❌ Failed to get plugin info: %v", err)
	}

	fmt.Printf("📋 Plugin Information:\n")
	fmt.Printf("   Name: %s\n", info.Name)
	fmt.Printf("   Version: %s\n", info.Version)
	fmt.Printf("   Description: %s\n", info.Description)
	fmt.Printf("   Channel Type: %s\n", info.ChannelType)
	fmt.Printf("   Capabilities: %v\n\n", info.Capabilities)

	// Get the plugin instance
	paymentChannel, err := loader.GetPlugin(channelID)
	if err != nil {
		log.Fatalf("❌ Failed to get plugin instance: %v", err)
	}

	// Initialize plugin with configuration
	config := map[string]interface{}{
		"mock_delay_ms": 50,  // 50ms delay for faster testing
		"success_rate":  0.9, // 90% success rate
	}

	err = paymentChannel.Initialize(config)
	if err != nil {
		log.Fatalf("❌ Failed to initialize plugin: %v", err)
	}

	// Validate configuration
	err = paymentChannel.ValidateConfig(config)
	if err != nil {
		log.Fatalf("❌ Configuration validation failed: %v", err)
	}

	fmt.Printf("⚙️  Plugin initialized with configuration:\n")
	configJSON, _ := json.MarshalIndent(config, "   ", "  ")
	fmt.Printf("%s\n\n", string(configJSON))

	// Demo: Collection Order (代收下单)
	fmt.Printf("💳 Demo: Collection Order (代收下单)\n")
	fmt.Printf("------------------------------------\n")
	
	collectReq := &interfaces.CollectOrderRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID:  "MERCHANT_001",
			ChannelID:   channelID,
			RequestID:   generateRequestID(),
			Timestamp:   time.Now(),
			ExtraParams: map[string]string{"test": "true"},
		},
		OrderID:     "ORDER_001",
		Amount:      100.50,
		Currency:    "CNY",
		Description: "Test payment for demo",
		ReturnURL:   "https://example.com/return",
		NotifyURL:   "https://example.com/notify",
		CustomerInfo: &interfaces.CustomerInfo{
			Name:     "John Doe",
			Email:    "john@example.com",
			Phone:    "+86-138-0013-8000",
			IDNumber: "110101199001011234",
		},
	}

	collectResp, err := paymentChannel.CollectOrder(context.Background(), collectReq)
	if err != nil {
		log.Printf("❌ Collection order failed: %v", err)
	} else {
		fmt.Printf("✅ Collection order created:\n")
		fmt.Printf("   Order ID: %s\n", collectResp.OrderID)
		fmt.Printf("   Channel Order ID: %s\n", collectResp.ChannelOrderID)
		fmt.Printf("   Amount: %.2f %s\n", collectResp.Amount, collectResp.Currency)
		fmt.Printf("   Status: %s\n", collectResp.Status)
		fmt.Printf("   Payment URL: %s\n", collectResp.PaymentURL)
		fmt.Printf("   QR Code: %s\n", collectResp.QRCode[:50] + "...")
	}

	// Demo: Balance Inquiry (余额查询)
	fmt.Printf("\n💰 Demo: Balance Inquiry (余额查询)\n")
	fmt.Printf("------------------------------------\n")
	
	balanceReq := &interfaces.BalanceInquiryRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  channelID,
			RequestID:  generateRequestID(),
			Timestamp:  time.Now(),
		},
		AccountType: "settlement",
	}

	balanceResp, err := paymentChannel.BalanceInquiry(context.Background(), balanceReq)
	if err != nil {
		log.Printf("❌ Balance inquiry failed: %v", err)
	} else {
		fmt.Printf("✅ Balance inquiry successful:\n")
		fmt.Printf("   Balance: %.2f %s\n", balanceResp.Balance, balanceResp.Currency)
		fmt.Printf("   Account Type: %s\n", balanceResp.AccountType)
		fmt.Printf("   Last Updated: %s\n", balanceResp.LastUpdated.Format(time.RFC3339))
	}

	// Demo: Payout Order (代付下单)
	fmt.Printf("\n💸 Demo: Payout Order (代付下单)\n")
	fmt.Printf("------------------------------------\n")
	
	payoutReq := &interfaces.PayoutOrderRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  channelID,
			RequestID:  generateRequestID(),
			Timestamp:  time.Now(),
		},
		OrderID:     "PAYOUT_001",
		Amount:      50.00,
		Currency:    "CNY",
		Description: "Test payout for demo",
		NotifyURL:   "https://example.com/payout-notify",
		RecipientInfo: &interfaces.RecipientInfo{
			Name:        "Jane Smith",
			BankAccount: "6222021234567890123",
			BankCode:    "ICBC",
			BankName:    "Industrial and Commercial Bank of China",
			Phone:       "+86-139-0013-9000",
			IDNumber:    "110101199002021234",
		},
	}

	payoutResp, err := paymentChannel.PayoutOrder(context.Background(), payoutReq)
	if err != nil {
		log.Printf("❌ Payout order failed: %v", err)
	} else {
		fmt.Printf("✅ Payout order created:\n")
		fmt.Printf("   Order ID: %s\n", payoutResp.OrderID)
		fmt.Printf("   Channel Order ID: %s\n", payoutResp.ChannelOrderID)
		fmt.Printf("   Amount: %.2f %s\n", payoutResp.Amount, payoutResp.Currency)
		fmt.Printf("   Status: %s\n", payoutResp.Status)
	}

	// Demo: Query Orders
	fmt.Printf("\n🔍 Demo: Query Orders\n")
	fmt.Printf("----------------------\n")

	// Wait a bit for orders to potentially complete
	time.Sleep(2 * time.Second)

	// Query collection order
	collectQueryReq := &interfaces.CollectQueryRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  channelID,
			RequestID:  generateRequestID(),
			Timestamp:  time.Now(),
		},
		OrderID: "ORDER_001",
	}

	collectQueryResp, err := paymentChannel.CollectQuery(context.Background(), collectQueryReq)
	if err != nil {
		log.Printf("❌ Collection query failed: %v", err)
	} else {
		fmt.Printf("✅ Collection order query:\n")
		fmt.Printf("   Order ID: %s\n", collectQueryResp.OrderID)
		fmt.Printf("   Status: %s\n", collectQueryResp.Status)
		if collectQueryResp.PaidAt != nil {
			fmt.Printf("   Paid At: %s\n", collectQueryResp.PaidAt.Format(time.RFC3339))
		}
	}

	// Query payout order
	payoutQueryReq := &interfaces.PayoutQueryRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  channelID,
			RequestID:  generateRequestID(),
			Timestamp:  time.Now(),
		},
		OrderID: "PAYOUT_001",
	}

	payoutQueryResp, err := paymentChannel.PayoutQuery(context.Background(), payoutQueryReq)
	if err != nil {
		log.Printf("❌ Payout query failed: %v", err)
	} else {
		fmt.Printf("✅ Payout order query:\n")
		fmt.Printf("   Order ID: %s\n", payoutQueryResp.OrderID)
		fmt.Printf("   Status: %s\n", payoutQueryResp.Status)
		if payoutQueryResp.CompletedAt != nil {
			fmt.Printf("   Completed At: %s\n", payoutQueryResp.CompletedAt.Format(time.RFC3339))
		}
	}

	// Demo: Callback Processing
	fmt.Printf("\n📞 Demo: Callback Processing (消息回调)\n")
	fmt.Printf("----------------------------------------\n")
	
	callbackReq := &interfaces.CallbackRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  channelID,
			RequestID:  generateRequestID(),
			Timestamp:  time.Now(),
		},
		CallbackType: "payment_notification",
		CallbackData: map[string]interface{}{
			"order_id": "ORDER_001",
			"status":   "paid",
			"amount":   100.50,
		},
		Signature: "mock_signature_12345",
	}

	callbackResp, err := paymentChannel.Callback(context.Background(), callbackReq)
	if err != nil {
		log.Printf("❌ Callback processing failed: %v", err)
	} else {
		fmt.Printf("✅ Callback processed:\n")
		fmt.Printf("   Processed: %t\n", callbackResp.Processed)
		fmt.Printf("   Message: %s\n", callbackResp.Message)
	}

	// Plugin health check
	fmt.Printf("\n🏥 Plugin Health Check\n")
	fmt.Printf("------------------------\n")
	health := loader.HealthCheck()
	for channelID, healthy := range health {
		status := "✅ Healthy"
		if !healthy {
			status = "❌ Unhealthy"
		}
		fmt.Printf("   %s: %s\n", channelID, status)
	}

	// List loaded plugins
	fmt.Printf("\n📋 Loaded Plugins\n")
	fmt.Printf("------------------\n")
	plugins := loader.ListPlugins()
	for channelID, loadedPlugin := range plugins {
		fmt.Printf("   %s:\n", channelID)
		fmt.Printf("     Path: %s\n", loadedPlugin.Path)
		fmt.Printf("     Loaded At: %s\n", loadedPlugin.LoadedAt.Format(time.RFC3339))
		fmt.Printf("     Last Used: %s\n", loadedPlugin.LastUsed.Format(time.RFC3339))
		fmt.Printf("     Usage Count: %d\n", loadedPlugin.UsageCount)
	}

	fmt.Printf("\n🎉 Demo completed successfully!\n")
	fmt.Printf("The plugin framework is working correctly.\n")
}

// generateRequestID generates a unique request ID for testing
func generateRequestID() string {
	return fmt.Sprintf("REQ_%d", time.Now().UnixNano())
}
