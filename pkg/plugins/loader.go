// Package plugins implements the plugin system.
package plugins

import (
	"fmt"
)

// Plugin represents a bot plugin.
// Plugins can extend bot functionality with new features.
type Plugin interface {
	// Name returns the plugin name.
	Name() string

	// Version returns the plugin version.
	Version() string

	// Load is called when the plugin is loaded.
	// The plugin should register event handlers and initialize resources.
	Load(b Bot) error

	// Unload is called when the plugin is unloaded.
	// The plugin should clean up resources and unregister handlers.
	Unload() error
}

// Loader manages plugin loading and unloading.
type Loader struct {
	bot     Bot
	plugins map[string]Plugin
}

// NewLoader creates a new plugin loader.
func NewLoader(b Bot) *Loader {
	return &Loader{
		bot:     b,
		plugins: make(map[string]Plugin),
	}
}

// Load loads a plugin.
func (l *Loader) Load(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}

	name := plugin.Name()

	if _, exists := l.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Call plugin's Load method
	if err := plugin.Load(l.bot); err != nil {
		return fmt.Errorf("failed to load plugin %s: %w", name, err)
	}

	l.plugins[name] = plugin

	// Emit plugin loaded event
	l.bot.Emit("plugin_loaded", name)

	return nil
}

// Unload unloads a plugin.
func (l *Loader) Unload(name string) error {
	plugin, exists := l.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not loaded", name)
	}

	// Call plugin's Unload method
	if err := plugin.Unload(); err != nil {
		return fmt.Errorf("failed to unload plugin %s: %w", name, err)
	}

	delete(l.plugins, name)

	// Emit plugin unloaded event
	l.bot.Emit("plugin_unloaded", name)

	return nil
}

// Get returns a loaded plugin by name.
func (l *Loader) Get(name string) (Plugin, bool) {
	plugin, exists := l.plugins[name]
	return plugin, exists
}

// Has returns true if a plugin is loaded.
func (l *Loader) Has(name string) bool {
	_, exists := l.plugins[name]
	return exists
}

// List returns all loaded plugin names.
func (l *Loader) List() []string {
	names := make([]string, 0, len(l.plugins))
	for name := range l.plugins {
		names = append(names, name)
	}
	return names
}

// UnloadAll unloads all plugins.
func (l *Loader) UnloadAll() error {
	for name := range l.plugins {
		if err := l.Unload(name); err != nil {
			return err
		}
	}
	return nil
}
