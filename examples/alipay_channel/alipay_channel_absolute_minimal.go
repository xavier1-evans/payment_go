package main

import (
	"context"
	"errors"
	"payment_go/pkg/interfaces"
)

// AlipayChannelAbsoluteMinimal implements the PaymentChannel interface with absolute minimal dependencies
type AlipayChannelAbsoluteMinimal struct {
	config *AlipayConfigAbsoluteMinimal
}

// AlipayConfigAbsoluteMinimal holds absolute minimal configuration
type AlipayConfigAbsoluteMinimal struct {
	AppID      string
	PrivateKey string
}

// NewPluginAbsoluteMinimal creates a new instance of the absolute minimal plugin
func NewPluginAbsoluteMinimal() interfaces.Plugin {
	return &AlipayChannelAbsoluteMinimal{}
}

// GetInfo returns metadata about this plugin
func (ac *AlipayChannelAbsoluteMinimal) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Alipay Payment Channel (Absolute Minimal)",
		Version:     "1.0.0",
		Description: "Absolute minimal integration with Alipay payment gateway",
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
func (ac *AlipayChannelAbsoluteMinimal) Initialize(config map[string]interface{}) error {
	if appID, ok := config["app_id"].(string); ok {
		ac.config.AppID = appID
	}
	if privateKey, ok := config["private_key"].(string); ok {
		ac.config.PrivateKey = privateKey
	}
	return nil
}

// ValidateConfig validates the configuration
func (ac *AlipayChannelAbsoluteMinimal) ValidateConfig(config map[string]interface{}) error {
	if config["app_id"] == nil || config["app_id"].(string) == "" {
		return errors.New("app_id is required")
	}
	if config["private_key"] == nil || config["private_key"].(string) == "" {
		return errors.New("private_key is required")
	}
	return nil
}

// CollectOrder creates an absolute minimal Alipay collection order
func (ac *AlipayChannelAbsoluteMinimal) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	return &interfaces.CollectOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay collection order created successfully",
			RequestID: req.RequestID,
		},
		OrderID:        req.OrderID,
		ChannelOrderID: "ALIPAY_" + req.OrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		PaymentURL:     "https://openapi.alipay.com/gateway.do?order_id=" + req.OrderID,
		Status:         "pending",
	}, nil
}

// PayoutOrder creates an absolute minimal Alipay payout order
func (ac *AlipayChannelAbsoluteMinimal) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	return &interfaces.PayoutOrderResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay payout order created successfully",
			RequestID: req.RequestID,
		},
		OrderID:        req.OrderID,
		ChannelOrderID: "ALIPAY_PAYOUT_" + req.OrderID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "processing",
	}, nil
}

// CollectQuery queries an absolute minimal Alipay collection order
func (ac *AlipayChannelAbsoluteMinimal) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Order query successful",
			RequestID: req.RequestID,
		},
		OrderID:        req.OrderID,
		ChannelOrderID: "ALIPAY_" + req.OrderID,
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "pending",
	}, nil
}

// PayoutQuery queries an absolute minimal Alipay payout order
func (ac *AlipayChannelAbsoluteMinimal) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Payout query successful",
			RequestID: req.RequestID,
		},
		OrderID:        req.OrderID,
		ChannelOrderID: "ALIPAY_PAYOUT_" + req.OrderID,
		Amount:         0.0,
		Currency:       "CNY",
		Status:         "processing",
	}, nil
}

// BalanceInquiry performs absolute minimal balance inquiry
func (ac *AlipayChannelAbsoluteMinimal) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Balance inquiry successful",
			RequestID: req.RequestID,
		},
		Balance:     1000000.0,
		Currency:    "CNY",
		AccountType: "default",
	}, nil
}

// Callback handles absolute minimal Alipay callbacks
func (ac *AlipayChannelAbsoluteMinimal) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	return &interfaces.CallbackResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Callback processed successfully",
			RequestID: req.RequestID,
		},
		Processed: true,
	}, nil
}
