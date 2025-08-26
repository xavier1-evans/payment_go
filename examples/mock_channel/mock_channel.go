package main

import (
	"context"
	"fmt"
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

// NewPlugin creates a new instance of the MockChannel plugin
// This function must be exported and named exactly "NewPlugin" for the plugin loader
func NewPlugin() interfaces.Plugin {
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
