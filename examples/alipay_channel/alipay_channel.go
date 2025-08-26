package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"payment_go/pkg/interfaces"
)

// AlipayChannel implements the PaymentChannel interface for Alipay integration
type AlipayChannel struct {
	config *AlipayConfig
	client *http.Client
}

// AlipayConfig holds the configuration for Alipay integration
type AlipayConfig struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	GatewayURL string `json:"gateway_url"`
	NotifyURL  string `json:"notify_url"`
	ReturnURL  string `json:"return_url"`
	Charset    string `json:"charset"`
	SignType   string `json:"sign_type"`
	Version    string `json:"version"`
	Timeout    int    `json:"timeout"`
}

// AlipayRequest represents a generic Alipay API request
type AlipayRequest struct {
	AppID      string            `json:"app_id"`
	Method     string            `json:"method"`
	Format     string            `json:"format"`
	Charset    string            `json:"charset"`
	SignType   string            `json:"sign_type"`
	Timestamp  string            `json:"timestamp"`
	Version    string            `json:"version"`
	NotifyURL  string            `json:"notify_url,omitempty"`
	ReturnURL  string            `json:"return_url,omitempty"`
	BizContent string            `json:"biz_content"`
	Sign       string            `json:"sign"`
	Extra      map[string]string `json:"-"`
}

// NewPlugin creates a new instance of the AlipayChannel plugin
func NewPlugin() interfaces.Plugin {
	return &AlipayChannel{}
}

// GetInfo returns metadata about this plugin
func (ac *AlipayChannel) GetInfo() *interfaces.PluginInfo {
	return &interfaces.PluginInfo{
		Name:        "Alipay Payment Channel",
		Version:     "1.0.0",
		Description: "Integration with Alipay payment gateway",
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
			"public_key": map[string]interface{}{
				"type":        "string",
				"required":    true,
				"description": "Alipay public key for verification",
			},
			"gateway_url": map[string]interface{}{
				"type":        "string",
				"default":     "https://openapi.alipay.com/gateway.do",
				"description": "Alipay gateway URL",
			},
			"timeout": map[string]interface{}{
				"type":        "integer",
				"default":     5000,
				"description": "Request timeout in milliseconds",
			},
		},
	}
}

// Initialize sets up the plugin with configuration
func (ac *AlipayChannel) Initialize(config map[string]interface{}) error {
	// Parse configuration
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	var alipayConfig AlipayConfig
	if err := json.Unmarshal(configJSON, &alipayConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults
	if alipayConfig.GatewayURL == "" {
		alipayConfig.GatewayURL = "https://openapi.alipay.com/gateway.do"
	}
	if alipayConfig.Charset == "" {
		alipayConfig.Charset = "utf-8"
	}
	if alipayConfig.SignType == "" {
		alipayConfig.SignType = "RSA2"
	}
	if alipayConfig.Version == "" {
		alipayConfig.Version = "1.0"
	}
	if alipayConfig.Timeout == 0 {
		alipayConfig.Timeout = 5000
	}

	ac.config = &alipayConfig
	ac.client = &http.Client{
		Timeout: time.Duration(alipayConfig.Timeout) * time.Millisecond,
	}

	return nil
}

// ValidateConfig validates the plugin configuration
func (ac *AlipayChannel) ValidateConfig(config map[string]interface{}) error {
	required := []string{"app_id", "private_key", "public_key"}
	for _, field := range required {
		if value, exists := config[field]; !exists || value == "" {
			return fmt.Errorf("required field '%s' is missing or empty", field)
		}
	}
	return nil
}

// CollectOrder creates an Alipay collection order
func (ac *AlipayChannel) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	// Create Alipay trade create request
	bizContent := map[string]interface{}{
		"out_trade_no": req.OrderID,
		"total_amount": fmt.Sprintf("%.2f", req.Amount),
		"subject":      req.Description,
		"buyer_id":     req.CustomerInfo.IDNumber, // Alipay user ID
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	alipayReq := &AlipayRequest{
		AppID:      ac.config.AppID,
		Method:     "alipay.trade.create",
		Format:     "JSON",
		Charset:    ac.config.Charset,
		SignType:   ac.config.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    ac.config.Version,
		NotifyURL:  req.NotifyURL,
		ReturnURL:  req.ReturnURL,
		BizContent: string(bizContentJSON),
	}

	// Sign the request
	ac.signRequest(alipayReq)

	// Send request to Alipay
	resp, err := ac.sendRequest(ctx, alipayReq)
	if err != nil {
		return &interfaces.CollectOrderResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ALIPAY_ERROR",
				Message:   fmt.Sprintf("Alipay request failed: %v", err),
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Parse response and create collection order response
	// This is a simplified implementation
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

// PayoutOrder creates an Alipay payout order
func (ac *AlipayChannel) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	// Create Alipay fund transfer request
	bizContent := map[string]interface{}{
		"out_biz_no":    req.OrderID,
		"payee_type":    "ALIPAY_LOGONID",
		"payee_account": req.RecipientInfo.BankAccount, // Alipay account
		"amount":        fmt.Sprintf("%.2f", req.Amount),
		"remark":        req.Description,
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	alipayReq := &AlipayRequest{
		AppID:      ac.config.AppID,
		Method:     "alipay.fund.trans.toaccount.transfer",
		Format:     "JSON",
		Charset:    ac.config.Charset,
		SignType:   ac.config.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    ac.config.Version,
		NotifyURL:  req.NotifyURL,
		BizContent: string(bizContentJSON),
	}

	// Sign the request
	ac.signRequest(alipayReq)

	// Send request to Alipay
	resp, err := ac.sendRequest(ctx, alipayReq)
	if err != nil {
		return &interfaces.PayoutOrderResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ALIPAY_ERROR",
				Message:   fmt.Sprintf("Alipay payout request failed: %v", err),
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Parse response and create payout order response
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

// CollectQuery queries an Alipay collection order
func (ac *AlipayChannel) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	// Create Alipay trade query request
	bizContent := map[string]interface{}{
		"out_trade_no": req.OrderID,
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	alipayReq := &AlipayRequest{
		AppID:      ac.config.AppID,
		Method:     "alipay.trade.query",
		Format:     "JSON",
		Charset:    ac.config.Charset,
		SignType:   ac.config.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    ac.config.Version,
		BizContent: string(bizContentJSON),
	}

	// Sign the request
	ac.signRequest(alipayReq)

	// Send request to Alipay
	resp, err := ac.sendRequest(ctx, alipayReq)
	if err != nil {
		return &interfaces.CollectQueryResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ALIPAY_ERROR",
				Message:   fmt.Sprintf("Alipay query request failed: %v", err),
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Parse response and create query response
	// This is a simplified implementation
	status := "pending"
	var paidAt *time.Time

	// In a real implementation, parse the actual Alipay response
	if resp != nil {
		// Parse resp to determine actual status and paid time
		status = "completed" // Simplified
		now := time.Now()
		paidAt = &now
	}

	return &interfaces.CollectQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay collection order queried successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_%s", req.OrderID),
		Amount:         0, // Would be parsed from response
		Currency:       "CNY",
		Status:         status,
		PaidAt:         paidAt,
	}, nil
}

// PayoutQuery queries an Alipay payout order
func (ac *AlipayChannel) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	// Create Alipay fund transfer query request
	bizContent := map[string]interface{}{
		"out_biz_no": req.OrderID,
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	alipayReq := &AlipayRequest{
		AppID:      ac.config.AppID,
		Method:     "alipay.fund.trans.order.query",
		Format:     "JSON",
		Charset:    ac.config.Charset,
		SignType:   ac.config.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    ac.config.Version,
		BizContent: string(bizContentJSON),
	}

	// Sign the request
	ac.signRequest(alipayReq)

	// Send request to Alipay
	resp, err := ac.sendRequest(ctx, alipayReq)
	if err != nil {
		return &interfaces.PayoutQueryResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "ALIPAY_ERROR",
				Message:   fmt.Sprintf("Alipay payout query request failed: %v", err),
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
		}, nil
	}

	// Parse response and create query response
	status := "processing"
	var completedAt *time.Time

	// In a real implementation, parse the actual Alipay response
	if resp != nil {
		// Parse resp to determine actual status and completion time
		status = "completed" // Simplified
		now := time.Now()
		completedAt = &now
	}

	return &interfaces.PayoutQueryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   "Alipay payout order queried successfully",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		OrderID:        req.OrderID,
		ChannelOrderID: fmt.Sprintf("ALIPAY_PAYOUT_%s", req.OrderID),
		Amount:         0, // Would be parsed from response
		Currency:       "CNY",
		Status:         status,
		CompletedAt:    completedAt,
	}, nil
}

// BalanceInquiry checks Alipay account balance
func (ac *AlipayChannel) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	// Note: Alipay doesn't provide a direct balance inquiry API
	// This is a placeholder implementation
	return &interfaces.BalanceInquiryResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   false,
			Code:      "NOT_SUPPORTED",
			Message:   "Balance inquiry not supported by Alipay",
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Balance:     0,
		Currency:    "CNY",
		AccountType: req.AccountType,
		LastUpdated: time.Now(),
	}, nil
}

// Callback processes Alipay notifications
func (ac *AlipayChannel) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	// Verify Alipay callback signature
	if !ac.verifyCallback(req) {
		return &interfaces.CallbackResponse{
			BaseResponse: interfaces.BaseResponse{
				Success:   false,
				Code:      "SIGNATURE_VERIFICATION_FAILED",
				Message:   "Alipay callback signature verification failed",
				RequestID: req.RequestID,
				Timestamp: time.Now(),
			},
			Processed: false,
			Message:   "Signature verification failed",
		}, nil
	}

	// Process the callback data
	// In a real implementation, parse the callback data and update order status
	processed := true
	message := "Alipay callback processed successfully"

	return &interfaces.CallbackResponse{
		BaseResponse: interfaces.BaseResponse{
			Success:   true,
			Code:      "SUCCESS",
			Message:   message,
			RequestID: req.RequestID,
			Timestamp: time.Now(),
		},
		Processed: processed,
		Message:   message,
	}, nil
}

// Helper methods for Alipay integration

func (ac *AlipayChannel) signRequest(req *AlipayRequest) {
	// In a real implementation, this would use RSA signing
	// This is a simplified MD5 signature for demonstration
	params := ac.buildQueryString(req)
	req.Sign = fmt.Sprintf("%x", md5.Sum([]byte(params+ac.config.PrivateKey)))
}

func (ac *AlipayChannel) verifyCallback(req *interfaces.CallbackRequest) bool {
	// In a real implementation, this would verify RSA signatures
	// This is a simplified verification for demonstration
	return req.Signature != ""
}

func (ac *AlipayChannel) buildQueryString(req *AlipayRequest) string {
	params := make(map[string]string)

	// Add all fields to params map
	params["app_id"] = req.AppID
	params["method"] = req.Method
	params["format"] = req.Format
	params["charset"] = req.Charset
	params["sign_type"] = req.SignType
	params["timestamp"] = req.Timestamp
	params["version"] = req.Version

	if req.NotifyURL != "" {
		params["notify_url"] = req.NotifyURL
	}
	if req.ReturnURL != "" {
		params["return_url"] = req.ReturnURL
	}
	params["biz_content"] = req.BizContent

	// Sort keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
	}

	return strings.Join(pairs, "&")
}

func (ac *AlipayChannel) sendRequest(ctx context.Context, req *AlipayRequest) (interface{}, error) {
	// In a real implementation, this would send HTTP requests to Alipay
	// This is a placeholder that simulates a successful response
	return map[string]interface{}{"success": true}, nil
}
