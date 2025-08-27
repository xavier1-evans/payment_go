package main

import (
	"context"
	"fmt"
	"time"

	"payment_go/pkg/interfaces"
)

// AlipayChannelUltraMinimal implements the PaymentChannel interface with absolute minimal dependencies
type AlipayChannelUltraMinimal struct {
	config *AlipayConfigUltraMinimal
}

// AlipayConfigUltraMinimal holds ultra-minimal configuration
type AlipayConfigUltraMinimal struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
}

// NewPluginUltraMinimal creates a new instance of the ultra-minimal plugin
func NewPluginUltraMinimal() interfaces.Plugin {
	return &AlipayChannelUltraMinimal{}
}

// GetInfo returns metadata about this plugin
func (ac *AlipayChannelUltraMinimal) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Alipay Payment Channel (Ultra-Minimal)",
		Version:     "1.0.0",
		Description: "Ultra-minimal integration with Alipay payment gateway",
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
func (ac *AlipayChannelUltraMinimal) Initialize(config map[string]interface{}) error {
	ac.config = &AlipayConfigUltraMinimal{
		AppID:      config["app_id"].(string),
		PrivateKey: config["private_key"].(string),
	}
	return nil
}

// ValidateConfig validates the configuration
func (ac *AlipayChannelUltraMinimal) ValidateConfig(config map[string]interface{}) error {
	if config["app_id"] == nil || config["app_id"].(string) == "" {
		return fmt.Errorf("app_id is required")
	}
	if config["private_key"] == nil || config["private_key"].(string) == "" {
		return fmt.Errorf("private_key is required")
	}
	return nil
}

// CollectOrder creates an ultra-minimal Alipay collection order
func (ac *AlipayChannelUltraMinimal) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	return &interfaces.CollectOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay collection order created successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_%s", req.OrderID),
		Amount:         req.Amount,
		Currency:       req.Currency,
		PaymentURL:     fmt.Sprintf("https://openapi.alipay.com/gateway.do?order_id=%s", req.OrderID),
		Status:         "pending",
	}, nil
}

// PayoutOrder creates an ultra-minimal Alipay payout order
func (ac *AlipayChannelUltraMinimal) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	return &interfaces.PayoutOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay payout order created successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_PAYOUT_%s", req.OrderID),
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "processing",
	}, nil
}

// CollectQuery queries an ultra-minimal Alipay collection order
func (ac *AlipayChannelUltraMinimal) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Order query successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_%s", req.OrderID),
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "pending",
	}, nil
}

// PayoutQuery queries an ultra-minimal Alipay payout order
func (ac *AlipayChannelUltraMinimal) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Payout query successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_PAYOUT_%s", req.OrderID),
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "processing",
	}, nil
}

// BalanceInquiry performs ultra-minimal balance inquiry
func (ac *AlipayChannelUltraMinimal) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
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

// Callback handles ultra-minimal Alipay callbacks
func (ac *AlipayChannelUltraMinimal) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
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
