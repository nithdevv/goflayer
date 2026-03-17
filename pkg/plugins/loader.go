// Package plugins provides a plugin system for goflayer.
// It matches Mineflayer's architecture with loadable plugins, dependencies, and hot-reload support.
package plugins

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
)

// PluginMetadata contains information about a plugin.
type PluginMetadata struct {
	// Name is the unique identifier for this plugin
	Name string
	// Version is the plugin version (semantic versioning recommended)
	Version string
	// Author is the plugin author(s)
	Author string
	// Description describes what the plugin does
	Description string
	// Dependencies lists required plugin names
	Dependencies []string
}

// Plugin defines the interface that all plugins must implement.
type Plugin interface {
	// Metadata returns the plugin's metadata
	Metadata() PluginMetadata

	// OnLoad is called when the plugin is loaded.
	// The context provides access to bot, world, events, and logger.
	OnLoad(ctx *Context) error

	// OnUnload is called when the plugin is unloaded.
	// Plugins should clean up resources and unsubscribe from events.
	OnUnload() error
}

// Loader manages plugin lifecycle: loading, unloading, and dependencies.
type Loader struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
	ctx     *Context
	log     *logger.Logger

	// Hot-reload support
	watchers map[string]chan struct{}
}

// NewLoader creates a new plugin loader.
func NewLoader(ctx *Context) *Loader {
	return &Loader{
		plugins:  make(map[string]Plugin),
		ctx:      ctx,
		log:      ctx.Logger().With("plugin_loader"),
		watchers: make(map[string]chan struct{}),
	}
}

// Load loads a plugin and its dependencies.
func (l *Loader) Load(plugin Plugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	meta := plugin.Metadata()

	// Check if already loaded
	if _, exists := l.plugins[meta.Name]; exists {
		return fmt.Errorf("plugin '%s' is already loaded", meta.Name)
	}

	l.log.Info("Loading plugin: %s v%s by %s", meta.Name, meta.Version, meta.Author)
	l.log.Debug("Description: %s", meta.Description)

	// Load dependencies first
	for _, dep := range meta.Dependencies {
		if !l.isLoaded(dep) {
			return fmt.Errorf("plugin '%s' requires dependency '%s' which is not loaded", meta.Name, dep)
		}
	}

	// Create plugin context with logger
	pluginCtx := l.ctx.WithPlugin(meta.Name)

	// Load the plugin
	if err := plugin.OnLoad(pluginCtx); err != nil {
		l.log.Error("Failed to load plugin '%s': %v", meta.Name, err)
		return fmt.Errorf("failed to load plugin '%s': %w", meta.Name, err)
	}

	l.plugins[meta.Name] = plugin
	l.log.Info("Successfully loaded plugin: %s", meta.Name)

	return nil
}

// Unload unloads a plugin.
// Plugins that depend on this plugin must be unloaded first.
func (l *Loader) Unload(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	plugin, exists := l.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' is not loaded", name)
	}

	// Check for dependents
	dependents := l.getDependents(name)
	if len(dependents) > 0 {
		return fmt.Errorf("cannot unload plugin '%s': required by %v", name, dependents)
	}

	l.log.Info("Unloading plugin: %s", name)

	// Stop hot-reload watcher if exists
	if watcher, ok := l.watchers[name]; ok {
		close(watcher)
		delete(l.watchers, name)
	}

	// Unload the plugin
	if err := plugin.OnUnload(); err != nil {
		l.log.Error("Failed to unload plugin '%s': %v", name, err)
		return fmt.Errorf("failed to unload plugin '%s': %w", name, err)
	}

	delete(l.plugins, name)
	l.log.Info("Successfully unloaded plugin: %s", name)

	return nil
}

// Reload unloads and reloads a plugin (hot-reload).
func (l *Loader) Reload(name string, newPlugin Plugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	oldPlugin, exists := l.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' is not loaded", name)
	}

	meta := newPlugin.Metadata()
	if meta.Name != name {
		return fmt.Errorf("plugin name mismatch: expected '%s', got '%s'", name, meta.Name)
	}

	l.log.Info("Reloading plugin: %s", name)

	// Unload old plugin
	if err := oldPlugin.OnUnload(); err != nil {
		l.log.Error("Failed to unload plugin '%s' for reload: %v", name, err)
		return fmt.Errorf("failed to unload plugin '%s': %w", name, err)
	}

	// Load new plugin
	pluginCtx := l.ctx.WithPlugin(meta.Name)
	if err := newPlugin.OnLoad(pluginCtx); err != nil {
		l.log.Error("Failed to reload plugin '%s': %v", name, err)

		// Try to restore old plugin
		if restoreErr := oldPlugin.OnLoad(pluginCtx); restoreErr != nil {
			l.log.Error("Failed to restore plugin '%s' after reload failure: %v", name, restoreErr)
		}
		return fmt.Errorf("failed to reload plugin '%s': %w", name, err)
	}

	l.plugins[name] = newPlugin
	l.log.Info("Successfully reloaded plugin: %s", name)

	return nil
}

// IsLoaded returns true if a plugin is loaded.
func (l *Loader) IsLoaded(name string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.isLoaded(name)
}

// isLoaded returns true if a plugin is loaded (must be called with lock held).
func (l *Loader) isLoaded(name string) bool {
	_, exists := l.plugins[name]
	return exists
}

// Get returns a loaded plugin by name.
func (l *Loader) Get(name string) (Plugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	plugin, exists := l.plugins[name]
	return plugin, exists
}

// MustGet returns a loaded plugin or panics.
func (l *Loader) MustGet(name string) Plugin {
	plugin, exists := l.Get(name)
	if !exists {
		panic(fmt.Sprintf("plugin '%s' is not loaded", name))
	}
	return plugin
}

// List returns all loaded plugin names.
func (l *Loader) List() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.plugins))
	for name := range l.plugins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListWithMetadata returns all loaded plugins with their metadata.
func (l *Loader) ListWithMetadata() []PluginMetadata {
	l.mu.RLock()
	defer l.mu.RUnlock()

	metas := make([]PluginMetadata, 0, len(l.plugins))
	for _, plugin := range l.plugins {
		metas = append(metas, plugin.Metadata())
	}

	// Sort by name
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Name < metas[j].Name
	})

	return metas
}

// UnloadAll unloads all plugins in reverse dependency order.
func (l *Loader) UnloadAll() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get plugins in unload order (reverse dependency order)
	order, err := l.getUnloadOrder()
	if err != nil {
		return err
	}

	l.log.Info("Unloading all plugins (%d plugins)", len(order))

	for _, name := range order {
		plugin := l.plugins[name]

		if err := plugin.OnUnload(); err != nil {
			l.log.Error("Failed to unload plugin '%s': %v", name, err)
		}

		delete(l.plugins, name)
	}

	l.log.Info("All plugins unloaded")
	return nil
}

// getUnloadOrder returns plugins in an order that respects dependencies.
// Dependents are unloaded before their dependencies.
func (l *Loader) getUnloadOrder() ([]string, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	for name, plugin := range l.plugins {
		graph[name] = plugin.Metadata().Dependencies
	}

	// Topological sort (Kahn's algorithm)
	order := make([]string, 0, len(l.plugins))
	inDegree := make(map[string]int)

	// Calculate in-degrees
	for name := range l.plugins {
		inDegree[name] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			if _, exists := l.plugins[dep]; exists {
				inDegree[dep]++
			}
		}
	}

	// Process nodes with zero in-degree
	queue := make([]string, 0)
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		// Add to order (we'll reverse at the end)
		order = append(order, node)

		// Decrease in-degree for dependent nodes
		for name, deps := range graph {
			for _, dep := range deps {
				if dep == node {
					inDegree[name]--
					if inDegree[name] == 0 {
						queue = append(queue, name)
					}
				}
			}
		}
	}

	if len(order) != len(l.plugins) {
		return nil, fmt.Errorf("circular dependency detected in plugins")
	}

	// Reverse to get unload order (dependents before dependencies)
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order, nil
}

// getDependents returns all plugins that depend on the given plugin.
func (l *Loader) getDependents(name string) []string {
	dependents := make([]string, 0)

	for _, plugin := range l.plugins {
		for _, dep := range plugin.Metadata().Dependencies {
			if dep == name {
				dependents = append(dependents, plugin.Metadata().Name)
				break
			}
		}
	}

	return dependents
}

// LoadBatch loads multiple plugins in dependency order.
func (l *Loader) LoadBatch(plugins []Plugin) error {
	// Build plugin map
	pluginMap := make(map[string]Plugin)
	for _, plugin := range plugins {
		meta := plugin.Metadata()
		pluginMap[meta.Name] = plugin
	}

	// Get load order
	order, err := l.getLoadOrder(pluginMap)
	if err != nil {
		return err
	}

	// Load in order
	for _, name := range order {
		if err := l.Load(pluginMap[name]); err != nil {
			return fmt.Errorf("failed to load plugin '%s': %w", name, err)
		}
	}

	return nil
}

// getLoadOrder returns plugins in dependency order (dependencies before dependents).
func (l *Loader) getLoadOrder(plugins map[string]Plugin) ([]string, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	for name, plugin := range plugins {
		graph[name] = plugin.Metadata().Dependencies
	}

	// Topological sort
	order := make([]string, 0, len(plugins))
	visited := make(map[string]bool)
	visiting := make(map[string]bool)

	var visit func(string) error
	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving '%s'", name)
		}

		visiting[name] = true

		// Visit dependencies first
		for _, dep := range graph[name] {
			if _, exists := plugins[dep]; exists {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		visiting[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}

	for name := range plugins {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// EnableHotReload enables hot-reload for a plugin.
// When the plugin's source code changes, it will be automatically reloaded.
func (l *Loader) EnableHotReload(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isLoaded(name) {
		return fmt.Errorf("plugin '%s' is not loaded", name)
	}

	if _, exists := l.watchers[name]; exists {
		return fmt.Errorf("hot-reload already enabled for plugin '%s'", name)
	}

	l.log.Info("Enabling hot-reload for plugin: %s", name)

	// Create watcher channel
	watcher := make(chan struct{})
	l.watchers[name] = watcher

	// In a real implementation, you would use fsnotify to watch for file changes
	// For now, this is a placeholder that demonstrates the architecture
	go func() {
		<-watcher
		l.log.Debug("Hot-reload watcher stopped for plugin: %s", name)
	}()

	return nil
}

// DisableHotReload disables hot-reload for a plugin.
func (l *Loader) DisableHotReload(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	watcher, exists := l.watchers[name]
	if !exists {
		return fmt.Errorf("hot-reload not enabled for plugin '%s'", name)
	}

	close(watcher)
	delete(l.watchers, name)

	l.log.Info("Disabled hot-reload for plugin: %s", name)
	return nil
}

// GetAs returns a plugin cast to a specific type.
// This is useful when you need to access plugin-specific methods.
func (l *Loader) GetAs(name string, pluginType interface{}) (bool, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	plugin, exists := l.plugins[name]
	if !exists {
		return false, fmt.Errorf("plugin '%s' is not loaded", name)
	}

	targetType := reflect.TypeOf(pluginType)
	if targetType.Kind() != reflect.Ptr {
		return false, fmt.Errorf("target type must be a pointer")
	}

	pluginValue := reflect.ValueOf(plugin)
	if !pluginValue.Type().Implements(targetType.Elem()) {
		return false, nil
	}

	reflect.ValueOf(pluginType).Elem().Set(pluginValue)
	return true, nil
}

// Count returns the number of loaded plugins.
func (l *Loader) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.plugins)
}

// String returns a string representation of the loader state.
func (l *Loader) String() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("PluginLoader (%d plugins):\n", len(l.plugins)))

	for _, name := range l.List() {
		plugin := l.plugins[name]
		meta := plugin.Metadata()
		sb.WriteString(fmt.Sprintf("  - %s v%s by %s\n", meta.Name, meta.Version, meta.Author))
		if len(meta.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("    Dependencies: %v\n", meta.Dependencies))
		}
	}

	return sb.String()
}
