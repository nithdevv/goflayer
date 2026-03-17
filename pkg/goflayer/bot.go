package goflayer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-flayer/goflayer/pkg/math"
	"github.com/go-flayer/goflayer/pkg/protocol"
)

// Bot represents a Minecraft bot.
//
// Bot is the main interface for interacting with Minecraft servers.
// It provides methods for controlling the bot, querying game state,
// and communicating with other players.
type Bot interface {
	// Connection methods
	Connect(ctx context.Context) error
	Disconnect() error

	// Event methods
	On(event string, handler EventHandler) Subscription
	Emit(event string, data ...interface{})

	// Plugin methods
	LoadPlugin(plugin Plugin) error
	HasPlugin(plugin Plugin) bool
	UnloadPlugin(plugin Plugin) error

	// State getters
	Entity() *Entity
	Players() map[string]*Player
	Entities() map[int32]*Entity
	World() World

	// Actions
	Chat(message string) error
	Attack(target Entity) error
	MoveTo(pos *math.Vec3) error
	LookAt(pos *math.Vec3) error

	// Getters
	Username() string
	Version() string
	Registry() Registry
	ProtocolClient() *protocol.Client
	EventBus() *EventBus

	// Lifecycle
	Start() error
	Stop() error
}

// bot implements the Bot interface.
type bot struct {
	// Core components
	client   *protocol.Client
	eventBus *EventBus
	registry Registry

	// State
	entity   *Entity
	players  map[string]*Player
	entities map[int32]*Entity
	world    World
	options  Options

	// Plugins
	plugins map[string]Plugin
	// FIXED: Use public interface for plugin loader
	pluginLoader PluginLoader

	// Channels
	packetChan chan *protocol.Packet
	eventChan  chan Event

	// Concurrency
	mu     sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	// Status
	connected bool
	stopped   bool
}

// CreateBot creates a new Minecraft bot with the given options.
func CreateBot(options Options) (Bot, error) {
	// Validate options
	if err := options.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options: %w", err)
	}

	// Set default values
	if options.Auth == "" {
		options.Auth = "offline"
	}

	ctx, cancel := context.WithCancel(context.Background())

	b := &bot{
		eventBus:   NewEventBus(),
		players:    make(map[string]*Player),
		entities:   make(map[int32]*Entity),
		plugins:    make(map[string]Plugin),
		options:    options,
		packetChan: make(chan *protocol.Packet, 256),
		eventChan:  make(chan Event, 100),
		ctx:        ctx,
		cancel:     cancel,
		connected:  false,
		stopped:    false,
	}

	// Create protocol client
	clientConfig := protocol.ClientConfig{
		Version:              options.Version,
		HideErrors:           options.HideErrors,
		CompressionThreshold: -1,
	}
	b.client = protocol.NewClient(clientConfig)

	// FIXED: Create plugin loader directly, not from external package
	b.pluginLoader = NewPluginLoader(b, options)

	// Load plugins
	if options.LoadInternalPlugins {
		if err := b.loadInternalPlugins(); err != nil {
			cancel()
			return nil, fmt.Errorf("failed to load internal plugins: %w", err)
		}
	}

	// Load external plugins
	if len(options.Plugins) > 0 {
		if err := b.loadExternalPlugins(); err != nil {
			cancel()
			return nil, fmt.Errorf("failed to load external plugins: %w", err)
		}
	}

	return b, nil
}

// Connect connects the bot to the Minecraft server.
func (b *bot) Connect(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.connected {
		return ErrBotAlreadyConnected
	}

	// Connect to server
	if err := b.client.Connect(ctx, b.options.Host, b.options.Port); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	b.connected = true

	// Start event loop
	b.wg.Add(1)
	go b.eventLoop()

	// Start packet processing
	b.wg.Add(1)
	go b.packetLoop()

	b.eventBus.Emit("connect")

	return nil
}

// Disconnect disconnects the bot from the server.
func (b *bot) Disconnect() error {
	b.mu.Lock()
	if !b.connected {
		b.mu.Unlock()
		return ErrBotNotConnected
	}

	b.connected = false
	b.mu.Unlock()

	// Cancel context first to stop goroutines
	b.cancel()

	// Disconnect client
	b.client.Disconnect("bot disconnect")

	// Wait for goroutines to finish
	b.wg.Wait()

	b.eventBus.Emit("disconnect")

	return nil
}

// On subscribes to an event.
func (b *bot) On(event string, handler EventHandler) Subscription {
	return b.eventBus.On(event, handler)
}

// Emit emits an event to all registered handlers.
func (b *bot) Emit(event string, data ...interface{}) {
	b.eventBus.Emit(event, data...)
}

// LoadPlugin loads a plugin into the bot.
func (b *bot) LoadPlugin(plugin Plugin) error {
	return b.pluginLoader.LoadPlugin(plugin)
}

// HasPlugin checks if a plugin is loaded.
func (b *bot) HasPlugin(plugin Plugin) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, exists := b.plugins[plugin.Name()]
	return exists
}

// UnloadPlugin unloads a plugin from the bot.
func (b *bot) UnloadPlugin(plugin Plugin) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// FIXED: Check if plugin exists before trying to unload
	if _, exists := b.plugins[plugin.Name()]; !exists {
		return fmt.Errorf("plugin %s is not loaded", plugin.Name())
	}

	// Cleanup plugin
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("plugin cleanup error: %w", err)
	}

	delete(b.plugins, plugin.Name())
	return nil
}

// Entity returns the bot's entity.
func (b *bot) Entity() *Entity {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.entity
}

// Players returns all players currently loaded.
func (b *bot) Players() map[string]*Player {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Return a copy to prevent race conditions
	players := make(map[string]*Player, len(b.players))
	for k, v := range b.players {
		players[k] = v
	}
	return players
}

// Entities returns all entities currently loaded.
func (b *bot) Entities() map[int32]*Entity {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Return a copy to prevent race conditions
	entities := make(map[int32]*Entity, len(b.entities))
	for k, v := range b.entities {
		entities[k] = v
	}
	return entities
}

// World returns the world interface.
func (b *bot) World() World {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.world
}

// Chat sends a chat message.
func (b *bot) Chat(message string) error {
	b.mu.RLock()
	connected := b.connected
	b.mu.RUnlock()

	if !connected {
		return ErrBotNotConnected
	}

	return b.client.Write("chat", map[string]interface{}{
		"message": message,
	})
}

// Attack attacks an entity.
func (b *bot) Attack(target Entity) error {
	b.mu.RLock()
	connected := b.connected
	b.mu.RUnlock()

	if !connected {
		return ErrBotNotConnected
	}

	return b.client.Write("use_entity", map[string]interface{}{
		"target": target.ID(),
	})
}

// MoveTo moves the bot to a position.
func (b *bot) MoveTo(pos *math.Vec3) error {
	// This will be implemented by the physics plugin
	return fmt.Errorf("not implemented")
}

// LookAt makes the bot look at a position.
func (b *bot) LookAt(pos *math.Vec3) error {
	// This will be implemented by the physics plugin
	return fmt.Errorf("not implemented")
}

// Username returns the bot's username.
func (b *bot) Username() string {
	return b.options.Username
}

// Version returns the Minecraft version being used.
func (b *bot) Version() string {
	return b.client.Version()
}

// Registry returns the version registry.
func (b *bot) Registry() Registry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.registry
}

// ProtocolClient returns the protocol client.
func (b *bot) ProtocolClient() *protocol.Client {
	return b.client
}

// EventBus returns the event bus.
func (b *bot) EventBus() *EventBus {
	return b.eventBus
}

// Start starts the bot's main loop.
func (b *bot) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.stopped {
		return ErrBotStopped
	}

	b.eventBus.Emit("start")
	return nil
}

// Stop stops the bot.
func (b *bot) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.stopped {
		return ErrBotStopped
	}

	b.stopped = true
	b.cancel()

	b.eventBus.Emit("stop")
	return nil
}

// eventLoop processes events from the event channel.
func (b *bot) eventLoop() {
	defer b.wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			return
		case event := <-b.eventChan:
			b.eventBus.Emit(event.Name, event.Data...)
		}
	}
}

// packetLoop processes packets from the protocol client.
func (b *bot) packetLoop() {
	defer b.wg.Done()

	// FIXED: Subscribe outside the loop to avoid duplicate subscriptions
	sub := b.client.On("packet", func(packet *protocol.Packet) {
		// Check if still connected before handling
		b.mu.RLock()
		connected := b.connected
		b.mu.RUnlock()

		if connected {
			b.handlePacket(packet)
		}
	})
	defer sub.Unsubscribe()

	// FIXED: Simply wait for context cancel instead of busy-wait
	<-b.ctx.Done()
}

// handlePacket handles an incoming packet.
func (b *bot) handlePacket(packet *protocol.Packet) {
	// Emit the packet event
	b.eventBus.Emit(packet.Name, packet.Data)

	// Handle specific packets
	switch packet.Name {
	case "login":
		b.handleLogin(packet)
	case "chat":
		b.handleChat(packet)
	}
}

// handleLogin handles the login packet.
func (b *bot) handleLogin(packet *protocol.Packet) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Create bot entity
	b.entity = &Entity{
		id:       0,
		username: b.options.Username,
		position: math.NewVec3(0, 0, 0),
	}

	b.eventBus.Emit("login")
	b.eventBus.Emit("spawn")
}

// handleChat handles the chat packet.
func (b *bot) handleChat(packet *protocol.Packet) {
	message, _ := packet.GetString("message")
	b.eventBus.Emit("chat", message)
}

// loadInternalPlugins loads all internal plugins.
func (b *bot) loadInternalPlugins() error {
	// TODO: Load internal plugins
	return nil
}

// loadExternalPlugins loads external plugins from options.
func (b *bot) loadExternalPlugins() error {
	for name, plugin := range b.options.Plugins {
		// FIXED: Type assertion to Plugin interface
		if p, ok := plugin.(Plugin); ok {
			if err := b.LoadPlugin(p); err != nil {
				return fmt.Errorf("failed to load plugin %s: %w", name, err)
			}
		} else {
			return fmt.Errorf("plugin %s does not implement Plugin interface", name)
		}
	}
	return nil
}

// Entity represents a game entity (player, mob, object).
type Entity struct {
	id       int32
	username string
	position *math.Vec3
	rotation *math.Vec3
	health   float64
}

// ID returns the entity ID.
func (e *Entity) ID() int32 {
	return e.id
}

// Username returns the entity's username (for players).
func (e *Entity) Username() string {
	return e.username
}

// Position returns the entity's position.
func (e *Entity) Position() *math.Vec3 {
	return e.position
}

// Rotation returns the entity's rotation (yaw, pitch).
func (e *Entity) Rotation() *math.Vec3 {
	return e.rotation
}

// Health returns the entity's health.
func (e *Entity) Health() float64 {
	return e.health
}

// Player represents a player entity.
type Player struct {
	Entity
	uuid string
}

// UUID returns the player's UUID.
func (p *Player) UUID() string {
	return p.uuid
}

// World represents the Minecraft world.
type World interface {
	// World methods will be implemented by the world plugin
}

// Registry represents a version-specific data registry.
type Registry interface {
	// Registry methods will be implemented by the registry package
}

// PluginLoader is the interface for loading and managing plugins.
// FIXED: Define this as a public interface in the goflayer package
type PluginLoader interface {
	LoadPlugin(plugin Plugin) error
	UnloadPlugin(plugin Plugin) error
	HasPlugin(name string) bool
	GetPlugin(name string) (Plugin, bool)
	CleanupAll() error
}

// NewPluginLoader creates a new plugin loader.
// FIXED: Factory function to create loader
func NewPluginLoader(bot Bot, options Options) PluginLoader {
	return &defaultPluginLoader{
		bot:     bot,
		options: options,
		plugins: make(map[string]Plugin),
		loaded:  make(map[string]bool),
	}
}

// defaultPluginLoader is the default implementation of PluginLoader.
type defaultPluginLoader struct {
	bot     Bot
	options Options
	plugins map[string]Plugin
	loaded  map[string]bool
	mu      sync.RWMutex
}

// LoadPlugin loads a single plugin into the bot.
func (l *defaultPluginLoader) LoadPlugin(plugin Plugin) error {
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

	// Also store in bot's plugin map
	if b, ok := l.bot.(*bot); ok {
		b.plugins[name] = plugin
	}

	return nil
}

// UnloadPlugin unloads a plugin from the bot.
func (l *defaultPluginLoader) UnloadPlugin(plugin Plugin) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	name := plugin.Name()

	// Check if loaded
	if !l.loaded[name] {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	// Cleanup the plugin
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("plugin %s cleanup error: %w", name, err)
	}

	// Remove from maps
	delete(l.plugins, name)
	delete(l.loaded, name)

	// Also remove from bot's plugin map
	if b, ok := l.bot.(*bot); ok {
		delete(b.plugins, name)
	}

	return nil
}

// HasPlugin checks if a plugin is loaded.
func (l *defaultPluginLoader) HasPlugin(name string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.loaded[name]
}

// GetPlugin returns a loaded plugin by name.
func (l *defaultPluginLoader) GetPlugin(name string) (Plugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	plugin, ok := l.plugins[name]
	return plugin, ok
}

// CleanupAll cleans up all loaded plugins.
func (l *defaultPluginLoader) CleanupAll() error {
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

	// Also clear bot's plugin map
	if b, ok := l.bot.(*bot); ok {
		b.plugins = make(map[string]Plugin)
	}

	// Return errors if any
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}
