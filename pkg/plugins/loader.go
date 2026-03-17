package plugins

import (
	"fmt"
	"sync"

	"github.com/go-flayer/goflayer/pkg/goflayer"
)

// pluginLoader manages plugin loading and unloading.
type pluginLoader struct {
	bot     goflayer.Bot
	options goflayer.Options
	plugins map[string]Plugin
	loaded  map[string]bool
	mu      sync.RWMutex
}

// newPluginLoader creates a new plugin loader.
func newPluginLoader(bot goflayer.Bot, options goflayer.Options) *pluginLoader {
	return &pluginLoader{
		bot:     bot,
		options: options,
		plugins: make(map[string]Plugin),
		loaded:  make(map[string]bool),
	}
}

// LoadPlugin loads a single plugin into the bot.
//
// The plugin's Inject method will be called with the bot instance.
// If the plugin fails to inject, it will not be loaded.
func (l *pluginLoader) LoadPlugin(plugin Plugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	name := plugin.Name()

	// Check if already loaded
	if l.loaded[name] {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	// Inject the plugin
	if err := plugin.Inject(l.bot, l.options); err != nil {
		return fmt.Errorf("failed to inject plugin %s: %w", name, err)
	}

	// Store and mark as loaded
	l.plugins[name] = plugin
	l.loaded[name] = true

	return nil
}

// LoadPlugins loads multiple plugins at once.
// If any plugin fails to load, returns an error without loading any further plugins.
func (l *pluginLoader) LoadPlugins(plugins []Plugin) error {
	for _, plugin := range plugins {
		if err := l.LoadPlugin(plugin); err != nil {
			return err
		}
	}
	return nil
}

// UnloadPlugin unloads a plugin from the bot.
//
// The plugin's Cleanup method will be called.
func (l *pluginLoader) UnloadPlugin(plugin Plugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	name := plugin.Name()

	// Check if loaded
	if !l.loaded[name] {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	// Get the plugin instance
	p, ok := l.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Cleanup the plugin
	if err := p.Cleanup(); err != nil {
		return fmt.Errorf("plugin %s cleanup error: %w", name, err)
	}

	// Remove from maps
	delete(l.plugins, name)
	delete(l.loaded, name)

	return nil
}

// HasPlugin checks if a plugin is loaded.
func (l *pluginLoader) HasPlugin(name string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.loaded[name]
}

// GetPlugin returns a loaded plugin by name.
func (l *pluginLoader) GetPlugin(name string) (Plugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	plugin, ok := l.plugins[name]
	return plugin, ok
}

// Plugins returns all loaded plugins.
func (l *pluginLoader) Plugins() map[string]Plugin {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Return a copy to prevent race conditions
	plugins := make(map[string]Plugin, len(l.plugins))
	for k, v := range l.plugins {
		plugins[k] = v
	}
	return plugins
}

// CleanupAll cleans up all loaded plugins.
//
// This is called when the bot is shutting down.
func (l *pluginLoader) CleanupAll() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var errs []error

	// Cleanup all plugins
	for name, plugin := range l.plugins {
		if err := plugin.Cleanup(); err != nil {
			errs = append(errs, fmt.Errorf("plugin %s cleanup error: %w", name, err))
		}
	}

	// Clear maps
	l.plugins = make(map[string]Plugin)
	l.loaded = make(map[string]bool)

	// Return errors if any
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}

// Count returns the number of loaded plugins.
func (l *pluginLoader) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.plugins)
}

// Names returns the names of all loaded plugins.
func (l *pluginLoader) Names() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0, len(l.loaded))
	for name := range l.loaded {
		names = append(names, name)
	}
	return names
}

// ReloadPlugin reloads a plugin by unloading and loading it again.
func (l *pluginLoader) ReloadPlugin(plugin Plugin) error {
	if err := l.UnloadPlugin(plugin); err != nil {
		// Ignore errors if plugin wasn't loaded
		_ = err
	}

	return l.LoadPlugin(plugin)
}

// PluginDependency represents a dependency between plugins.
type PluginDependency struct {
	Plugin     string
	DependsOn  string
	Optional   bool
}

// LoadWithDependencies loads plugins respecting their dependencies.
//
// Plugins are loaded in dependency order - dependencies are loaded
// before the plugins that depend on them.
func (l *pluginLoader) LoadWithDependencies(plugins []Plugin, dependencies []PluginDependency) error {
	// Build dependency graph
	depGraph := make(map[string][]string)
	for _, dep := range dependencies {
		depGraph[dep.Plugin] = append(depGraph[dep.Plugin], dep.DependsOn)
	}

	// Sort plugins by dependency order
	sorted, err := l.topologicalSort(plugins, depGraph)
	if err != nil {
		return err
	}

	// Load in order
	for _, plugin := range sorted {
		if err := l.LoadPlugin(plugin); err != nil {
			return err
		}
	}

	return nil
}

// topologicalSort performs topological sort on plugins based on dependencies.
func (l *pluginLoader) topologicalSort(plugins []Plugin, dependencies map[string][]string) ([]Plugin, error) {
	// Create plugin map
	pluginMap := make(map[string]Plugin)
	for _, plugin := range plugins {
		pluginMap[plugin.Name()] = plugin
	}

	// Kahn's algorithm for topological sorting
	sorted := make([]Plugin, 0)
	inDegree := make(map[string]int)

	// Calculate in-degrees
	for name := range pluginMap {
		inDegree[name] = 0
	}
	for _, deps := range dependencies {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	// Start with nodes that have no dependencies
	queue := make([]string, 0)
	for name := range pluginMap {
		if inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	// Process queue
	for len(queue) > 0 {
		// Get a node with no dependencies
		name := queue[0]
		queue = queue[1:]

		// Add to sorted list
		sorted = append(sorted, pluginMap[name])

		// Remove edges from this node
		for _, dep := range dependencies[name] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
	}

	// Check for cycles
	if len(sorted) != len(plugins) {
		return nil, fmt.Errorf("circular dependency detected in plugins")
	}

	return sorted, nil
}
