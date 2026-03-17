// Package bot implements the core Minecraft bot.
package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/pkg/event"
	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
	"github.com/nithdevv/goflayer/pkg/protocol/states"
)

// botImpl implements the Bot interface.
type botImpl struct {
	// Configuration
	options *Options

	// Protocol client
	client *protocol.Client

	// Event bus
	events *event.Bus

	// Plugins
	plugins      map[string]plugins.Plugin
	pluginsMu    sync.RWMutex
	pluginLoader *plugins.Loader

	// State
	connected bool
	mu        sync.RWMutex

	// Context
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Options holds bot configuration.
type Options struct {
	// Host is the server address.
	Host string

	// Port is the server port.
	Port int

	// Username is the bot's username.
	Username string

	// Password is the bot's password (optional, for online auth).
	Password string

	// Version is the Minecraft version to connect to.
	Version string

	// Auth is the authentication type ("offline" or "microsoft").
	Auth string

	// Plugins are plugins to load on startup.
	Plugins map[string]plugins.Plugin
}

// New creates a new bot with the given options.
func New(options *Options) (Bot, error) {
	if options == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	// Set defaults
	if options.Version == "" {
		options.Version = "1.20.1"
	}
	if options.Auth == "" {
		options.Auth = "offline"
	}

	ctx, cancel := context.WithCancel(context.Background())

	bot := &botImpl{
		options:     options,
		events:      event.NewBus(),
		plugins:     make(map[string]plugins.Plugin),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Create protocol client
	bot.client = protocol.NewClient(protocol.Config{
		Host:            options.Host,
		Port:            options.Port,
		Version:         options.Version,
		ProtocolVersion: 763, // 1.20.1
	})

	// Create plugin loader
	bot.pluginLoader = plugins.NewLoader(bot)

	return bot, nil
}

// Connect connects to the server.
func (b *botImpl) Connect(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.connected {
		return fmt.Errorf("already connected")
	}

	// Connect protocol client
	if err := b.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Perform handshake
	if err := b.handshake(); err != nil {
		b.client.Disconnect()
		return fmt.Errorf("handshake failed: %w", err)
	}

	// Start packet processing
	b.wg.Add(1)
	go b.processPackets()

	// Load plugins
	if err := b.loadPlugins(); err != nil {
		b.client.Disconnect()
		return fmt.Errorf("failed to load plugins: %w", err)
	}

	b.connected = true

	// Emit connected event
	b.events.Emit("connected")

	return nil
}

// Disconnect disconnects from the server.
func (b *botImpl) Disconnect() error {
	b.mu.Lock()
	if !b.connected {
		b.mu.Unlock()
		return fmt.Errorf("not connected")
	}
	b.connected = false
	b.mu.Unlock()

	// Unload plugins
	b.unloadPlugins()

	// Cancel context
	b.cancel()

	// Wait for goroutines
	b.wg.Wait()

	// Disconnect client
	if err := b.client.Disconnect(); err != nil {
		return err
	}

	// Emit disconnected event
	b.events.Emit("disconnected")

	return nil
}

// handshake performs the login handshake.
func (b *botImpl) handshake() error {
	// Send handshake packet
	// TODO: Implement handshake

	// Transition to login state
	if err := b.client.SetState(states.Login); err != nil {
		return err
	}

	// Send login start packet
	// TODO: Implement login

	// Wait for login success
	// TODO: Implement login handling

	// Transition to play state
	if err := b.client.SetState(states.Play); err != nil {
		return err
	}

	return nil
}

// loadPlugins loads all configured plugins.
func (b *botImpl) loadPlugins() error {
	// Load configured plugins
	if b.options.Plugins != nil {
		for name, plugin := range b.options.Plugins {
			if err := b.pluginLoader.Load(plugin); err != nil {
				return fmt.Errorf("failed to load plugin %s: %w", name, err)
			}
			b.plugins[name] = plugin
		}
	}
	return nil
}

// unloadPlugins unloads all plugins.
func (b *botImpl) unloadPlugins() {
	b.pluginsMu.Lock()
	defer b.pluginsMu.Unlock()

	for name, plugin := range b.plugins {
		if err := plugin.Unload(); err != nil {
			// Log error but continue
			_ = err
		}
		delete(b.plugins, name)
	}
}

// processPackets processes incoming packets.
func (b *botImpl) processPackets() {
	defer b.wg.Done()

	for {
		select {
		case <-b.ctx.Done():
			return
		case packet := <-b.client.Incoming():
			if packet == nil {
				return
			}

			// Emit packet event
			b.events.Emit("packet", packet)

			// Handle packet
			b.handlePacket(packet)
		}
	}
}

// handlePacket handles an incoming packet.
func (b *botImpl) handlePacket(packet *protocol.Packet) {
	// Emit specific packet event
	packetEvent := fmt.Sprintf("packet:%d", packet.ID)
	b.events.Emit(packetEvent, packet)
}

// On subscribes to an event.
// Returns a subscription that can be used to unsubscribe.
func (b *botImpl) On(event string, handler func(...interface{})) event.Subscription {
	return b.events.Subscribe(event, handler)
}

// Emit emits an event.
func (b *botImpl) Emit(event string, data ...interface{}) {
	b.events.Emit(event, data...)
}

// Chat sends a chat message.
func (b *botImpl) Chat(message string) error {
	b.mu.RLock()
	connected := b.connected
	b.mu.RUnlock()

	if !connected {
		return fmt.Errorf("not connected")
	}

	// TODO: Implement chat packet
	return nil
}

// Events returns the event bus.
func (b *botImpl) Events() *event.Bus {
	return b.events
}

// Client returns the protocol client.
func (b *botImpl) Client() *protocol.Client {
	return b.client
}

// IsConnected returns true if connected to the server.
func (b *botImpl) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.connected
}

// Username returns the bot's username.
func (b *botImpl) Username() string {
	return b.options.Username
}

// Options returns the bot options.
func (b *botImpl) Options() *Options {
	return b.options
}
