package interfaces

import (
	"context"
	"time"
)

// PaymentChannel defines the standard interface for payment channel plugins
// This interface allows the payment gateway to communicate with different upstream providers
// through a unified API, regardless of the specific payment channel implementation.
type PaymentChannel interface {
	// CollectOrder creates a collection order (代收下单)
	// This is typically the busiest operation and should be highly optimized
	CollectOrder(ctx context.Context, req *CollectOrderRequest) (*CollectOrderResponse, error)
	
	// PayoutOrder creates a payout order (代付下单)
	PayoutOrder(ctx context.Context, req *PayoutOrderRequest) (*PayoutOrderResponse, error)
	
	// CollectQuery queries a collection order status (代收查单)
	CollectQuery(ctx context.Context, req *CollectQueryRequest) (*CollectQueryResponse, error)
	
	// PayoutQuery queries a payout order status (代付查单)
	PayoutQuery(ctx context.Context, req *PayoutQueryRequest) (*PayoutQueryResponse, error)
	
	// BalanceInquiry checks account balance (余额查询)
	BalanceInquiry(ctx context.Context, req *BalanceInquiryRequest) (*BalanceInquiryResponse, error)
	
	// Callback processes incoming messages from upstream providers (消息回调)
	Callback(ctx context.Context, req *CallbackRequest) (*CallbackResponse, error)
}

// Common request/response structures
type BaseRequest struct {
	MerchantID   string            `json:"merchant_id"`
	ChannelID    string            `json:"channel_id"`
	RequestID    string            `json:"request_id"`
	Timestamp    time.Time         `json:"timestamp"`
	ExtraParams  map[string]string `json:"extra_params,omitempty"`
}

type BaseResponse struct {
	Success      bool              `json:"success"`
	Code         string            `json:"code"`
	Message      string            `json:"message"`
	RequestID    string            `json:"request_id"`
	Timestamp    time.Time         `json:"timestamp"`
	ExtraData    map[string]string `json:"extra_data,omitempty"`
}

// Collection Order (代收下单)
type CollectOrderRequest struct {
	BaseRequest
	OrderID      string  `json:"order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
	ReturnURL    string  `json:"return_url"`
	NotifyURL    string  `json:"notify_url"`
	CustomerInfo *CustomerInfo `json:"customer_info,omitempty"`
}

type CollectOrderResponse struct {
	BaseResponse
	OrderID      string  `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	PaymentURL   string  `json:"payment_url,omitempty"`
	QRCode       string  `json:"qr_code,omitempty"`
	Status       string  `json:"status"`
}

// Payout Order (代付下单)
type PayoutOrderRequest struct {
	BaseRequest
	OrderID      string  `json:"order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Description  string  `json:"description"`
	NotifyURL    string  `json:"notify_url"`
	RecipientInfo *RecipientInfo `json:"recipient_info"`
}

type PayoutOrderResponse struct {
	BaseResponse
	OrderID      string  `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
}

// Query Requests
type CollectQueryRequest struct {
	BaseRequest
	OrderID      string `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id,omitempty"`
}

type CollectQueryResponse struct {
	BaseResponse
	OrderID      string  `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
}

type PayoutQueryRequest struct {
	BaseRequest
	OrderID      string `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id,omitempty"`
}

type PayoutQueryResponse struct {
	BaseResponse
	OrderID      string  `json:"order_id"`
	ChannelOrderID string `json:"channel_order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Status       string  `json:"status"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// Balance Inquiry (余额查询)
type BalanceInquiryRequest struct {
	BaseRequest
	AccountType  string `json:"account_type,omitempty"`
}

type BalanceInquiryResponse struct {
	BaseResponse
	Balance      float64 `json:"balance"`
	Currency     string  `json:"currency"`
	AccountType  string  `json:"account_type"`
	LastUpdated  time.Time `json:"last_updated"`
}

// Callback (消息回调)
type CallbackRequest struct {
	BaseRequest
	CallbackType string            `json:"callback_type"`
	CallbackData map[string]interface{} `json:"callback_data"`
	Signature    string            `json:"signature"`
}

type CallbackResponse struct {
	BaseResponse
	Processed    bool   `json:"processed"`
	Message      string `json:"message"`
}

// Supporting structures
type CustomerInfo struct {
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	IDNumber    string `json:"id_number,omitempty"`
}

type RecipientInfo struct {
	Name        string `json:"name"`
	BankAccount string `json:"bank_account"`
	BankCode    string `json:"bank_code"`
	BankName    string `json:"bank_name"`
	Phone       string `json:"phone,omitempty"`
	IDNumber    string `json:"id_number,omitempty"`
}

// Plugin metadata and configuration
type PluginInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	ChannelType string            `json:"channel_type"`
	Capabilities []string         `json:"capabilities"`
	ConfigSchema map[string]interface{} `json:"config_schema"`
}

// Plugin interface for metadata
type Plugin interface {
	PaymentChannel
	GetInfo() *PluginInfo
	Initialize(config map[string]interface{}) error
	ValidateConfig(config map[string]interface{}) error
}
