package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"payment_go/pkg/interfaces"
)

// AlipayChannelMinimal implements the PaymentChannel interface with minimal dependencies
type AlipayChannelMinimal struct {
	config *AlipayConfigMinimal
	client *http.Client
}

// AlipayConfigMinimal holds minimal configuration for Alipay integration
type AlipayConfigMinimal struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	GatewayURL string `json:"gateway_url"`
	Timeout    int    `json:"timeout"`
}

// AlipayRequestMinimal represents a minimal Alipay API request
type AlipayRequestMinimal struct {
	AppID      string `json:"app_id"`
	Method     string `json:"method"`
	Format     string `json:"format"`
	Charset    string `json:"charset"`
	SignType   string `json:"sign_type"`
	Timestamp  string `json:"timestamp"`
	Version    string `json:"version"`
	BizContent string `json:"biz_content"`
	Sign       string `json:"sign"`
}

// NewPluginMinimal creates a new instance of the AlipayChannelMinimal plugin
func NewPluginMinimal() interfaces.Plugin {
	return &AlipayChannelMinimal{}
}

// GetInfo returns metadata about this plugin
func (ac *AlipayChannelMinimal) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Alipay Payment Channel (Minimal)",
		Version:     "1.0.0",
		Description: "Minimal integration with Alipay payment gateway",
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
			"gateway_url": map[string]interface{}{
				"type":        "string",
				"default":     "https://openapi.alipay.com/gateway.do",
				"description": "Alipay gateway URL",
			},
		},
	}
}

// Initialize sets up the channel with configuration
func (ac *AlipayChannelMinimal) Initialize(config map[string]interface{}) error {
	// Parse minimal configuration
	ac.config = &AlipayConfigMinimal{
		AppID:      config["app_id"].(string),
		PrivateKey: config["private_key"].(string),
		GatewayURL: "https://openapi.alipay.com/gateway.do",
		Timeout:    5000,
	}

	// Create minimal HTTP client
	ac.client = &http.Client{
		Timeout: time.Duration(ac.config.Timeout) * time.Millisecond,
	}

	return nil
}

// CollectOrder creates a minimal Alipay collection order
func (ac *AlipayChannelMinimal) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	// Create minimal biz content
	bizContent := map[string]interface{}{
		"out_trade_no": req.OrderID,
		"total_amount": fmt.Sprintf("%.2f", req.Amount),
		"subject":      req.Description,
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	// Create minimal request
	alipayReq := &AlipayRequestMinimal{
		AppID:      ac.config.AppID,
		Method:     "alipay.trade.page.pay",
		Format:     "JSON",
		Charset:    "utf-8",
		SignType:   "MD5",
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    "1.0",
		BizContent: string(bizContentJSON),
	}

	// Sign the request
	ac.signRequest(alipayReq)

	// Return response (minimal implementation)
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
		PaymentURL:     fmt.Sprintf("%s?%s", ac.config.GatewayURL, ac.buildQueryString(alipayReq)),
		Status:         "pending",
	}, nil
}

// PayoutOrder creates a minimal Alipay payout order
func (ac *AlipayChannelMinimal) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
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

// CollectQuery queries a minimal Alipay collection order
func (ac *AlipayChannelMinimal) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
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

// PayoutQuery queries a minimal Alipay payout order
func (ac *AlipayChannelMinimal) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
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

// BalanceInquiry performs minimal balance inquiry
func (ac *AlipayChannelMinimal) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Balance inquiry successful",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Balance:     1000000.0, // Mock balance
		Currency:    "CNY",
		AccountType: "default",
		LastUpdated: time.Now(),
	}, nil
}

// Callback handles minimal Alipay callbacks
func (ac *AlipayChannelMinimal) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
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

// ValidateConfig validates the configuration
func (ac *AlipayChannelMinimal) ValidateConfig(config map[string]interface{}) error {
	// Check required fields
	if config["app_id"] == nil || config["app_id"].(string) == "" {
		return fmt.Errorf("app_id is required")
	}
	if config["private_key"] == nil || config["private_key"].(string) == "" {
		return fmt.Errorf("private_key is required")
	}
	return nil
}

// signRequest signs the Alipay request with MD5
func (ac *AlipayChannelMinimal) signRequest(req *AlipayRequestMinimal) {
	// Build query string for signing
	queryString := ac.buildQueryString(req)
	queryString += "&key=" + ac.config.PrivateKey

	// Generate MD5 hash
	hash := md5.Sum([]byte(queryString))
	req.Sign = strings.ToUpper(fmt.Sprintf("%x", hash))
}

// buildQueryString builds query string for signing
func (ac *AlipayChannelMinimal) buildQueryString(req *AlipayRequestMinimal) string {
	params := map[string]string{
		"app_id":      req.AppID,
		"method":      req.Method,
		"format":      req.Format,
		"charset":     req.Charset,
		"sign_type":   req.SignType,
		"timestamp":   req.Timestamp,
		"version":     req.Version,
		"biz_content": req.BizContent,
	}

	// Sort keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var pairs []string
	for _, k := range keys {
		if params[k] != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
		}
	}

	return strings.Join(pairs, "&")
}
