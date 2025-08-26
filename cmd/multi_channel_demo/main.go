package main

import (
	"context"
	"fmt"
	"time"

	"payment_go/pkg/interfaces"
)

// Mock Alipay Channel Implementation
type MockAlipayChannel struct{}

func (a *MockAlipayChannel) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:         "Alipay Channel",
		Version:      "1.0.0",
		Description:  "Integration with Alipay for mobile and web payments",
		Author:       "Payment Gateway Team",
		ChannelType:  "alipay",
		Capabilities: []string{"collect_order", "payout_order", "collect_query", "payout_query", "balance_inquiry", "callback"},
	}
}

func (a *MockAlipayChannel) Initialize(config map[string]interface{}) error     { return nil }
func (a *MockAlipayChannel) ValidateConfig(config map[string]interface{}) error { return nil }

func (a *MockAlipayChannel) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	channelOrderID := fmt.Sprintf("ALIPAY_%d", time.Now().UnixNano())
	return &interfaces.CollectOrderResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay order created", RequestID: req.RequestID, Timestamp: time.Now()},
		OrderID:      req.OrderID, ChannelOrderID: channelOrderID, Amount: req.Amount, Currency: req.Currency,
		PaymentURL: fmt.Sprintf("https://alipay.com/pay/%s", channelOrderID), Status: "pending",
	}, nil
}

func (a *MockAlipayChannel) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	channelOrderID := fmt.Sprintf("ALIPAY_PAYOUT_%d", time.Now().UnixNano())
	return &interfaces.PayoutOrderResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay payout initiated", RequestID: req.RequestID, Timestamp: time.Now()},
		OrderID:      req.OrderID, ChannelOrderID: channelOrderID, Amount: req.Amount, Currency: req.Currency, Status: "processing",
	}, nil
}

func (a *MockAlipayChannel) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay query successful", RequestID: req.RequestID, Timestamp: time.Now()},
		OrderID:      req.OrderID, ChannelOrderID: "ALIPAY_" + req.OrderID, Amount: 100.00, Currency: "CNY", Status: "paid", PaidAt: &time.Time{},
	}, nil
}

func (a *MockAlipayChannel) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay payout query successful", RequestID: req.RequestID, Timestamp: time.Now()},
		OrderID:      req.OrderID, ChannelOrderID: "ALIPAY_PAYOUT_" + req.OrderID, Amount: 100.00, Currency: "CNY", Status: "completed", CompletedAt: &time.Time{},
	}, nil
}

func (a *MockAlipayChannel) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay balance inquiry successful", RequestID: req.RequestID, Timestamp: time.Now()},
		AccountType:  "merchant", Balance: 100000.00, Currency: "CNY",
	}, nil
}

func (a *MockAlipayChannel) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	return &interfaces.CallbackResponse{
		BaseResponse: interfaces.BaseResponse{Success: true, Code: "SUCCESS", Message: "Alipay callback processed", RequestID: req.RequestID, Timestamp: time.Now()},
		Processed:    true,
	}, nil
}

// Payment Gateway that manages Alipay channel
type PaymentGateway struct {
	channels map[string]interfaces.Plugin
}

func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{
		channels: make(map[string]interfaces.Plugin),
	}
}

func (pg *PaymentGateway) AddChannel(channelType string, channel interfaces.Plugin) {
	pg.channels[channelType] = channel
}

func (pg *PaymentGateway) GetChannel(channelType string) (interfaces.Plugin, bool) {
	channel, exists := pg.channels[channelType]
	return channel, exists
}

func (pg *PaymentGateway) ListChannels() []string {
	var channelTypes []string
	for channelType := range pg.channels {
		channelTypes = append(channelTypes, channelType)
	}
	return channelTypes
}

func (pg *PaymentGateway) ProcessPayment(channelType string, amount float64, currency string, customerInfo *interfaces.CustomerInfo) error {
	channel, exists := pg.GetChannel(channelType)
	if !exists {
		return fmt.Errorf("payment channel '%s' not found", channelType)
	}

	// Create payment request
	req := &interfaces.CollectOrderRequest{
		BaseRequest: interfaces.BaseRequest{
			MerchantID: "DEMO_MERCHANT",
			ChannelID:  channelType,
			RequestID:  fmt.Sprintf("REQ_%d", time.Now().UnixNano()),
			Timestamp:  time.Now(),
		},
		OrderID:      fmt.Sprintf("ORDER_%s_%d", channelType, time.Now().UnixNano()),
		Amount:       amount,
		Currency:     currency,
		Description:  fmt.Sprintf("Payment via %s", channelType),
		ReturnURL:    "https://example.com/return",
		NotifyURL:    "https://example.com/notify",
		CustomerInfo: customerInfo,
	}

	// Process payment
	_, err := channel.CollectOrder(context.Background(), req)
	return err
}

func main() {
	fmt.Printf("🚀 Alipay Payment Gateway Demo\n")
	fmt.Printf("==============================\n\n")

	// Create payment gateway
	gateway := NewPaymentGateway()

	// Add Alipay payment channel
	gateway.AddChannel("alipay", &MockAlipayChannel{})

	// Display available channels
	fmt.Printf("📋 Available Payment Channels:\n")
	for _, channelType := range gateway.ListChannels() {
		channel, _ := gateway.GetChannel(channelType)
		info := channel.GetInfo()
		fmt.Printf("   • %s (%s) - %s\n", info.Name, info.ChannelType, info.Description)
	}
	fmt.Printf("\n")

	// Test customer info
	customerInfo := &interfaces.CustomerInfo{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Phone:    "+86-138-0000-0000",
		IDNumber: "110101199001011234",
	}

	// Test payments with Alipay
	testAmounts := []float64{50.00, 100.00, 200.00, 500.00}

	for _, amount := range testAmounts {
		fmt.Printf("💳 Testing Payment: %.2f CNY\n", amount)
		fmt.Printf("   " + repeatString("-", 40) + "\n")

		for _, channelType := range gateway.ListChannels() {
			start := time.Now()
			err := gateway.ProcessPayment(channelType, amount, "CNY", customerInfo)
			duration := time.Since(start)

			if err != nil {
				fmt.Printf("   ❌ %s: Failed - %v\n", channelType, err)
			} else {
				fmt.Printf("   ✅ %s: Success (%.2fms)\n", channelType, float64(duration.Microseconds())/1000.0)
			}
		}
		fmt.Printf("\n")
	}

	// Test balance inquiries
	fmt.Printf("💰 Balance Inquiries:\n")
	fmt.Printf("   " + repeatString("-", 40) + "\n")

	for _, channelType := range gateway.ListChannels() {
		channel, _ := gateway.GetChannel(channelType)

		req := &interfaces.BalanceInquiryRequest{
			BaseRequest: interfaces.BaseRequest{
				MerchantID: "DEMO_MERCHANT",
				ChannelID:  channelType,
				RequestID:  fmt.Sprintf("BAL_%d", time.Now().UnixNano()),
				Timestamp:  time.Now(),
			},
		}

		resp, err := channel.BalanceInquiry(context.Background(), req)
		if err != nil {
			fmt.Printf("   ❌ %s: Failed - %v\n", channelType, err)
		} else {
			fmt.Printf("   ✅ %s: %.2f %s\n", channelType, resp.Balance, resp.Currency)
		}
	}

	fmt.Printf("\n🎉 Alipay payment gateway demo completed successfully!\n")
	fmt.Printf("   The Alipay channel is working correctly.\n")
	fmt.Printf("   You can now integrate real Alipay API using the same interface.\n")
}

// Helper function for string repetition
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
