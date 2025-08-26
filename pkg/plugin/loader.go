package plugin

import (
	"fmt"
	"plugin"
	"sync"
	"time"

	"payment_go/pkg/interfaces"
)

// PluginLoader manages the loading and lifecycle of payment channel plugins
type PluginLoader struct {
	plugins map[string]*LoadedPlugin
	mutex   sync.RWMutex
}

// LoadedPlugin represents a loaded plugin with its metadata and instance
type LoadedPlugin struct {
	Path       string
	Plugin     *plugin.Plugin
	Instance   interfaces.Plugin
	Info       *interfaces.PluginInfo
	LoadedAt   time.Time
	LastUsed   time.Time
	UsageCount int64
}

// NewPluginLoader creates a new plugin loader instance
func NewPluginLoader() *PluginLoader {
	return &PluginLoader{
		plugins: make(map[string]*LoadedPlugin),
	}
}

// LoadPlugin loads a payment channel plugin from a .so file
func (pl *PluginLoader) LoadPlugin(pluginPath, channelID string) error {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	// Check if plugin is already loaded
	if _, exists := pl.plugins[channelID]; exists {
		return fmt.Errorf("plugin for channel %s is already loaded", channelID)
	}

	// Open the .so file
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", pluginPath, err)
	}

	// Look up the required symbols
	newPluginFunc, err := p.Lookup("NewPlugin")
	if err != nil {
		return fmt.Errorf("plugin %s missing NewPlugin function: %w", pluginPath, err)
	}

	// Type assert the function
	newPlugin, ok := newPluginFunc.(func() interfaces.Plugin)
	if !ok {
		return fmt.Errorf("plugin %s NewPlugin function has wrong signature", pluginPath)
	}

	// Create plugin instance
	instance := newPlugin()

	// Get plugin info
	info := instance.GetInfo()

	// Validate plugin info
	if err := pl.validatePluginInfo(info); err != nil {
		return fmt.Errorf("plugin %s validation failed: %w", pluginPath, err)
	}

	// Store the loaded plugin
	pl.plugins[channelID] = &LoadedPlugin{
		Path:     pluginPath,
		Plugin:   p,
		Instance: instance,
		Info:     info,
		LoadedAt: time.Now(),
	}

	return nil
}

// GetPlugin retrieves a loaded plugin by channel ID
func (pl *PluginLoader) GetPlugin(channelID string) (interfaces.Plugin, error) {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	loadedPlugin, exists := pl.plugins[channelID]
	if !exists {
		return nil, fmt.Errorf("plugin for channel %s not found", channelID)
	}

	// Update usage statistics
	loadedPlugin.LastUsed = time.Now()
	loadedPlugin.UsageCount++

	return loadedPlugin.Instance, nil
}

// UnloadPlugin unloads a plugin and removes it from memory
func (pl *PluginLoader) UnloadPlugin(channelID string) error {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	_, exists := pl.plugins[channelID]
	if !exists {
		return fmt.Errorf("plugin for channel %s not found", channelID)
	}

	// Note: Go plugins cannot be fully unloaded from memory
	// We can only remove the reference
	delete(pl.plugins, channelID)

	return nil
}

// ListPlugins returns information about all loaded plugins
func (pl *PluginLoader) ListPlugins() map[string]*LoadedPlugin {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	result := make(map[string]*LoadedPlugin)
	for k, v := range pl.plugins {
		result[k] = v
	}
	return result
}

// GetPluginInfo returns metadata for a specific plugin
func (pl *PluginLoader) GetPluginInfo(channelID string) (*interfaces.PluginInfo, error) {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	loadedPlugin, exists := pl.plugins[channelID]
	if !exists {
		return nil, fmt.Errorf("plugin for channel %s not found", channelID)
	}

	return loadedPlugin.Instance.GetInfo(), nil
}

// validatePluginInfo validates that a plugin has the required metadata
func (pl *PluginLoader) validatePluginInfo(info *interfaces.PluginInfo) error {
	if info == nil {
		return fmt.Errorf("plugin info is nil")
	}
	if info.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if info.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if info.ChannelType == "" {
		return fmt.Errorf("plugin channel type is required")
	}
	if len(info.Capabilities) == 0 {
		return fmt.Errorf("plugin must declare at least one capability")
	}

	return nil
}

// ReloadPlugin reloads a plugin from disk (useful for development/testing)
func (pl *PluginLoader) ReloadPlugin(channelID string) error {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()

	loadedPlugin, exists := pl.plugins[channelID]
	if !exists {
		return fmt.Errorf("plugin for channel %s not found", channelID)
	}

	// Unload current plugin
	delete(pl.plugins, channelID)

	// Reload from disk
	return pl.LoadPlugin(loadedPlugin.Path, channelID)
}

// HealthCheck performs a basic health check on all loaded plugins
func (pl *PluginLoader) HealthCheck() map[string]bool {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()

	health := make(map[string]bool)
	for channelID, loadedPlugin := range pl.plugins {
		// Try to get plugin info as a basic health check
		info := loadedPlugin.Instance.GetInfo()
		health[channelID] = info != nil
	}
	return health
}
