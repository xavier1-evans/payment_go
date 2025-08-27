package main

import (
	"time"

	"payment_go/pkg/interfaces"
)

// AlipayChannelUltraStripped implements the PaymentChannel interface with absolute minimal dependencies
type AlipayChannelUltraStripped struct {
	config *AlipayConfigUltraStripped
}

// AlipayConfigUltraStripped holds ultra-stripped configuration
type AlipayConfigUltraStripped struct {
	AppID      string
	PrivateKey string
}

// NewPluginUltraStripped creates a new instance of the ultra-stripped plugin
func NewPluginUltraStripped() interfaces.Plugin {
	return &AlipayChannelUltraStripped{}
}

// GetInfo returns metadata about this plugin
func (ac *AlipayChannelUltraStripped) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Alipay Payment Channel (Ultra-Stripped)",
		Version:     "1.0.0",
		Description: "Ultra-stripped integration with Alipay payment gateway",
		Author:      "Payment Gateway Team",
		ChannelType: "alipay",
		Capabilities: []string{
			"collect_order",
			"payout_order",
			"collect_query",
			"payout_query",
			"balance_inquiry",
			"callback",
		},
		ConfigSchema: map[string]interface{}{
			"app_id": map[string]interface{}{
				"type":        "string",
				"required":    true,
				"description": "Alipay application ID",
			},
			"private_key": map[string]interface{}{
				"type":        "string",
				"required":    true,
				"description": "Alipay private key for signing",
			},
		},
	}
}

// Initialize sets up the channel with configuration
func (ac *AlipayChannelUltraStripped) Initialize(config map[string]interface{}) error {
	ac.config = &AlipayConfigUltraStripped{
		AppID:      config["app_id"].(string),
		PrivateKey: config["private_key"].(string),
	}
	return nil
}

// ValidateConfig validates the configuration
func (ac *AlipayChannelUltraStripped) ValidateConfig(config map[string]interface{}) error {
	if config["app_id"] == nil || config["app_id"].(string) == "" {
		return &validationError{field: "app_id"}
	}
	if config["private_key"] == nil || config["private_key"].(string) == "" {
		return &validationError{field: "private_key"}
	}
	return nil
}

// CollectOrder creates an ultra-stripped Alipay collection order
func (ac *AlipayChannelUltraStripped) CollectOrder(ctx interface{}, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	channelOrderID := "ALIPAY_" + req.OrderID
	paymentURL := "https://openapi.alipay.com/gateway.do?order_id=" + req.OrderID
	
	return &interfaces.CollectOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay collection order created successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		PaymentURL:     paymentURL,
		Status:         "pending",
	}, nil
}

// PayoutOrder creates an ultra-stripped Alipay payout order
func (ac *AlipayChannelUltraStripped) PayoutOrder(ctx interface{}, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	channelOrderID := "ALIPAY_PAYOUT_" + req.OrderID
	
	return &interfaces.PayoutOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay payout order created successfully",
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

// CollectQuery queries an ultra-stripped Alipay collection order
func (ac *AlipayChannelUltraStripped) CollectQuery(ctx interface{}, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	channelOrderID := "ALIPAY_" + req.OrderID
	
	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Order query successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "pending",
	}, nil
}

// PayoutQuery queries an ultra-stripped Alipay payout order
func (ac *AlipayChannelUltraStripped) PayoutQuery(ctx interface{}, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	channelOrderID := "ALIPAY_PAYOUT_" + req.OrderID
	
	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Payout query successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: channelOrderID,
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "processing",
	}, nil
}

// BalanceInquiry performs ultra-stripped balance inquiry
func (ac *AlipayChannelUltraStripped) BalanceInquiry(ctx interface{}, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Balance inquiry successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Balance:     1000000.0,
		Currency:    "CNY",
		AccountType: "default",
		LastUpdated: time.Now(),
	}, nil
}

// Callback handles ultra-stripped Alipay callbacks
func (ac *AlipayChannelUltraStripped) Callback(ctx interface{}, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	return &interfaces.CallbackResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Callback processed successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Processed: true,
	}, nil
}

// validationError is a simple error type without fmt dependency
type validationError struct {
	field string
}

func (e *validationError) Error() string {
	return e.field + " is required"
}
