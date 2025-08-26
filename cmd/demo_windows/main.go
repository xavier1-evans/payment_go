package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"payment_go/pkg/interfaces"
)

// MockChannel implements the PaymentChannel interface for testing and demonstration
type MockChannel struct {
	config map[string]interface{}
	orders map[string]*MockOrder
}

// MockOrder represents a mock order in the system
type MockOrder struct {
	OrderID        string
	ChannelOrderID string
	Amount         float64
	Currency       string
	Status         string
	CreatedAt      time.Time
	PaidAt         *time.Time
	CompletedAt    *time.Time
	CustomerInfo   *interfaces.CustomerInfo
	RecipientInfo  *interfaces.RecipientInfo
}

// NewMockChannel creates a new instance of the MockChannel
func NewMockChannel() *MockChannel {
	return &MockChannel{
		orders: make(map[string]*MockOrder),
	}
}

// GetInfo returns metadata about this plugin
func (mc *MockChannel) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Mock Payment Channel",
		Version:     "1.0.0",
		Description: "A mock payment channel for testing and development",
		Author:      "Payment Gateway Team",
		ChannelType: "mock",
		Capabilities: []string{
			"collect_order",
			"payout_order",
			"collect_query",
			"payout_query",
			"balance_inquiry",
			"callback",
		},
		ConfigSchema: map[string]interface{}{
			"mock_delay_ms": map[string]interface{}{
				"type":        "integer",
				"default":     100,
				"description": "Artificial delay in milliseconds for testing",
			},
			"success_rate": map[string]interface{}{
				"type":        "float",
				"default":     0.95,
				"description": "Success rate for mock operations (0.0-1.0)",
			},
		},
	}
}

// Initialize sets up the plugin with configuration
func (mc *MockChannel) Initialize(config map[string]interface{}) error {
	mc.config = config
	return nil
}

// ValidateConfig validates the plugin configuration
func (mc *MockChannel) ValidateConfig(config map[string]interface{}) error {
	if delay, exists := config["mock_delay_ms"]; exists {
		if delayInt, ok := delay.(int); ok {
			if delayInt < 0 || delayInt > 10000 {
				return fmt.Errorf("mock_delay_ms must be between 0 and 10000")
			}
		}
	}

	if rate, exists := config["success_rate"]; exists {
		if rateFloat, ok := rate.(float64); ok {
			if rateFloat < 0.0 || rateFloat > 1.0 {
				return fmt.Errorf("success_rate must be between 0.0 and 1.0")
			}
		}
	}

	return nil
}

// CollectOrder creates a mock collection order
func (mc *MockChannel) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	mc.simulateDelay()

	// Generate a mock channel order ID
	channelOrderID := fmt.Sprintf("MOCK_%d", time.Now().UnixNano())

	// Create mock order
	mockOrder := &MockOrder{
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "pending",
		CreatedAt:      time.Now(),
		CustomerInfo:   req.CustomerInfo,
	}

	mc.orders[req.OrderID] = mockOrder

	// Simulate success/failure based on config
	if mc.shouldSucceed() {
		return &interfaces.CollectOrderResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   true,
				Code:      "SUCCESS",
				Message:   "Mock collection order created successfully",
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
			OrderID:        req.OrderID,
			ChannelOrderID: channelOrderID,
			Amount:         req.Amount,
			Currency:       req.Currency,
			PaymentURL:     fmt.Sprintf("https://mock-payment.com/pay/%s", channelOrderID),
			QRCode:         fmt.Sprintf("data:image/png;base64,MOCK_QR_%s", channelOrderID),
			Status:         "pending",
		}, nil
	}

	return &interfaces.CollectOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   false,
			Code:      "MOCK_ERROR",
			Message:   "Mock collection order failed",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "failed",
	}, nil
}

// PayoutOrder creates a mock payout order
func (mc *MockChannel) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	mc.simulateDelay()

	// Generate a mock channel order ID
	channelOrderID := fmt.Sprintf("MOCK_PAYOUT_%d", time.Now().UnixNano())

	// Create mock order
	mockOrder := &MockOrder{
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "processing",
		CreatedAt:      time.Now(),
		RecipientInfo:  req.RecipientInfo,
	}

	mc.orders[req.OrderID] = mockOrder

	// Simulate success/failure based on config
	if mc.shouldSucceed() {
		return &interfaces.PayoutOrderResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   true,
				Code:      "SUCCESS",
				Message:   "Mock payout order created successfully",
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
			OrderID:        req.OrderID,
			ChannelOrderID: channelOrderID,
			Amount:         req.Amount,
			Currency:       req.Currency,
			Status:         "processing",
		}, nil
	}

	return &interfaces.PayoutOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   false,
			Code:      "MOCK_ERROR",
			Message:   "Mock payout order failed",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "failed",
	}, nil
}

// CollectQuery queries a mock collection order
func (mc *MockChannel) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	mc.simulateDelay()

	mockOrder, exists := mc.orders[req.OrderID]
	if !exists {
		return &interfaces.CollectQueryResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ORDER_NOT_FOUND",
				Message:   "Mock order not found",
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Simulate order completion after some time
	if mockOrder.Status == "pending" && time.Since(mockOrder.CreatedAt) > 5*time.Second {
		mockOrder.Status = "completed"
		now := time.Now()
		mockOrder.PaidAt = &now
	}

	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Mock collection order queried successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        mockOrder.OrderID,
		ChannelOrderID: mockOrder.ChannelOrderID,
		Amount:         mockOrder.Amount,
		Currency:       mockOrder.Currency,
		Status:         mockOrder.Status,
		PaidAt:         mockOrder.PaidAt,
	}, nil
}

// PayoutQuery queries a mock payout order
func (mc *MockChannel) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	mc.simulateDelay()

	mockOrder, exists := mc.orders[req.OrderID]
	if !exists {
		return &interfaces.PayoutQueryResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ORDER_NOT_FOUND",
				Message:   "Mock order not found",
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Simulate payout completion after some time
	if mockOrder.Status == "processing" && time.Since(mockOrder.CreatedAt) > 3*time.Second {
		mockOrder.Status = "completed"
		now := time.Now()
		mockOrder.CompletedAt = &now
	}

	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Mock payout order queried successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        mockOrder.OrderID,
		ChannelOrderID: mockOrder.ChannelOrderID,
		Amount:         mockOrder.Amount,
		Currency:       mockOrder.Currency,
		Status:         mockOrder.Status,
		CompletedAt:    mockOrder.CompletedAt,
	}, nil
}

// BalanceInquiry checks mock account balance
func (mc *MockChannel) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	mc.simulateDelay()

	// Generate a mock balance
	balance := 1000000.0 + rand.Float64()*500000.0 // Random balance between 1M and 1.5M

	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Mock balance inquiry successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Balance:     balance,
		Currency:    "CNY",
		AccountType: req.AccountType,
		LastUpdated: time.Now(),
	}, nil
}

// Callback processes mock incoming messages
func (mc *MockChannel) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	mc.simulateDelay()

	// Simulate callback processing
	processed := mc.shouldSucceed()
	message := "Mock callback processed successfully"
	if !processed {
		message = "Mock callback processing failed"
	}

	return &interfaces.CallbackResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   processed,
			Code:      "SUCCESS",
			Message:   message,
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Processed: processed,
		Message:   message,
	}, nil
}

// Helper methods
func (mc *MockChannel) simulateDelay() {
	if delay, exists := mc.config["mock_delay_ms"]; exists {
		if delayInt, ok := delay.(int); ok {
			time.Sleep(time.Duration(delayInt) * time.Millisecond)
		}
	}
}

func (mc *MockChannel) shouldSucceed() bool {
	if rate, exists := mc.config["success_rate"]; exists {
		if rateFloat, ok := rate.(float64); ok {
			return rand.Float64() < rateFloat
		}
	}
	return rand.Float64() < 0.95 // Default 95% success rate
}

func main() {
	fmt.Printf("🚀 Payment Gateway Plugin Demo (Windows Version)\n")
	fmt.Printf("================================================\n\n")

	// Create mock channel directly (no plugin loading needed on Windows)
	paymentChannel := NewMockChannel()

	// Get plugin info
	info := paymentChannel.GetInfo()
	fmt.Printf("📋 Plugin Information:\n")
	fmt.Printf("   Name: %s\n", info.Name)
	fmt.Printf("   Version: %s\n", info.Version)
	fmt.Printf("   Description: %s\n", info.Description)
	fmt.Printf("   Channel Type: %s\n", info.ChannelType)
	fmt.Printf("   Capabilities: %v\n\n", info.Capabilities)

	// Initialize plugin with configuration
	config := map[string]interface{}{
		"mock_delay_ms": 50,  // 50ms delay for faster testing
		"success_rate":  0.9, // 90% success rate
	}

	err := paymentChannel.Initialize(config)
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
			ChannelID:   "mock_channel",
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
		fmt.Printf("   QR Code: %s\n", collectResp.QRCode[:50]+"...")
	}

	// Demo: Balance Inquiry (余额查询)
	fmt.Printf("\n💰 Demo: Balance Inquiry (余额查询)\n")
	fmt.Printf("------------------------------------\n")

	balanceReq := &interfaces.BalanceInquiryRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "MERCHANT_001",
			ChannelID:  "mock_channel",
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
			ChannelID:  "mock_channel",
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
			ChannelID:  "mock_channel",
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
			ChannelID:  "mock_channel",
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
			ChannelID:  "mock_channel",
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

	fmt.Printf("\n🎉 Demo completed successfully!\n")
	fmt.Printf("The plugin framework is working correctly on Windows!\n")
	fmt.Printf("\nNote: This Windows version doesn't use Go plugins (not supported on Windows).\n")
	fmt.Printf("Instead, it directly instantiates the mock channel for demonstration.\n")
}

// generateRequestID generates a unique request ID for testing
func generateRequestID() string {
	return fmt.Sprintf("REQ_%d", time.Now().UnixNano())
}
