package plugin

import (
	"context"
	"testing"
	"time"

	"payment_go/pkg/interfaces"
)

// MockPlugin implements the interfaces.Plugin for testing
type MockPlugin struct {
	info *interfaces.PluginInfo
}

func (mp *MockPlugin) GetInfo() *interfaces.PluginInfo {
	return mp.info
}

func (mp *MockPlugin) Initialize(config map[string]interface{}) error {
	return nil
}

func (mp *MockPlugin) ValidateConfig(config map[string]interface{}) error {
	return nil
}

func (mp *MockPlugin) CollectOrder(ctx context.Context, req *interfaces.CollectOrderRequest) (*interfaces.CollectOrderResponse, error) {
	return &interfaces.CollectOrderResponse{}, nil
}

func (mp *MockPlugin) PayoutOrder(ctx context.Context, req *interfaces.PayoutOrderRequest) (*interfaces.PayoutOrderResponse, error) {
	return &interfaces.PayoutOrderResponse{}, nil
}

func (mp *MockPlugin) CollectQuery(ctx context.Context, req *interfaces.CollectQueryRequest) (*interfaces.CollectQueryResponse, error) {
	return &interfaces.CollectQueryResponse{}, nil
}

func (mp *MockPlugin) PayoutQuery(ctx context.Context, req *interfaces.PayoutQueryRequest) (*interfaces.PayoutQueryResponse, error) {
	return &interfaces.PayoutQueryResponse{}, nil
}

func (mp *MockPlugin) BalanceInquiry(ctx context.Context, req *interfaces.BalanceInquiryRequest) (*interfaces.BalanceInquiryResponse, error) {
	return &interfaces.BalanceInquiryResponse{}, nil
}

func (mp *MockPlugin) Callback(ctx context.Context, req *interfaces.CallbackRequest) (*interfaces.CallbackResponse, error) {
	return &interfaces.CallbackResponse{}, nil
}

func TestNewPluginLoader(t *testing.T) {
	loader := NewPluginLoader()
	if loader == nil {
		t.Fatal("NewPluginLoader returned nil")
	}

	if loader.plugins == nil {
		t.Fatal("PluginLoader plugins map is nil")
	}
}

func TestValidatePluginInfo(t *testing.T) {
	loader := NewPluginLoader()

	// Test valid plugin info
	validInfo := &interfaces.PluginInfo{
		Name:         "Test Plugin",
		Version:      "1.0.0",
		ChannelType:  "test",
		Capabilities: []string{"collect_order"},
	}

	err := loader.validatePluginInfo(validInfo)
	if err != nil {
		t.Errorf("Valid plugin info should not cause error: %v", err)
	}

	// Test invalid plugin info
	testCases := []struct {
		name        string
		info        *interfaces.PluginInfo
		shouldError bool
	}{
		{
			name:        "nil info",
			info:        nil,
			shouldError: true,
		},
		{
			name: "empty name",
			info: &interfaces.PluginInfo{
				Name:         "",
				Version:      "1.0.0",
				ChannelType:  "test",
				Capabilities: []string{"collect_order"},
			},
			shouldError: true,
		},
		{
			name: "empty version",
			info: &interfaces.PluginInfo{
				Name:         "Test Plugin",
				Version:      "",
				ChannelType:  "test",
				Capabilities: []string{"collect_order"},
			},
			shouldError: true,
		},
		{
			name: "empty channel type",
			info: &interfaces.PluginInfo{
				Name:         "Test Plugin",
				Version:      "1.0.0",
				ChannelType:  "",
				Capabilities: []string{"collect_order"},
			},
			shouldError: true,
		},
		{
			name: "empty capabilities",
			info: &interfaces.PluginInfo{
				Name:         "Test Plugin",
				Version:      "1.0.0",
				ChannelType:  "test",
				Capabilities: []string{},
			},
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := loader.validatePluginInfo(tc.info)
			if tc.shouldError && err == nil {
				t.Errorf("Expected error for %s, but got none", tc.name)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for %s, but got: %v", tc.name, err)
			}
		})
	}
}

func TestPluginLoaderOperations(t *testing.T) {
	loader := NewPluginLoader()

	// Test initial state
	plugins := loader.ListPlugins()
	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins initially, got %d", len(plugins))
	}

	// Test health check on empty loader
	health := loader.HealthCheck()
	if len(health) != 0 {
		t.Errorf("Expected 0 health checks initially, got %d", len(health))
	}

	// Test getting non-existent plugin
	_, err := loader.GetPlugin("non_existent")
	if err == nil {
		t.Error("Expected error when getting non-existent plugin")
	}

	// Test getting info for non-existent plugin
	_, err = loader.GetPluginInfo("non_existent")
	if err == nil {
		t.Error("Expected error when getting info for non-existent plugin")
	}

	// Test unloading non-existent plugin
	err = loader.UnloadPlugin("non_existent")
	if err == nil {
		t.Error("Expected error when unloading non-existent plugin")
	}
}

func TestLoadedPluginFields(t *testing.T) {
	plugin := &LoadedPlugin{
		Path:       "/test/path",
		Instance:   nil,
		Info:       nil,
		LoadedAt:   time.Now(),
		LastUsed:   time.Now(),
		UsageCount: 42,
	}

	if plugin.Path != "/test/path" {
		t.Errorf("Expected path /test/path, got %s", plugin.Path)
	}

	if plugin.UsageCount != 42 {
		t.Errorf("Expected usage count 42, got %d", plugin.UsageCount)
	}
}
