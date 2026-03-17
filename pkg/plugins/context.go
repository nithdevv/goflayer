// Package plugins provides plugin context for accessing bot components.
package plugins

import (
	"context"
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/bot"
	"github.com/nithdevv/goflayer/pkg/events"
)

// World represents the game world (placeholder for future implementation).
// This will include chunks, blocks, entities, etc.
type World struct {
	// TODO: Implement world state
	// - Chunks map[ChunkPos]*Chunk
	// - Entities map[int32]*Entity
	// - Blocks map[BlockPos]*Block
	// - Players map[uuid]*Player
	mu sync.RWMutex
}

// Context provides plugins with access to bot components.
// Each plugin gets its own context instance with a namespaced logger.
type Context struct {
	bot    *bot.Bot
	world  *World
	events *events.Bus
	log    *logger.Logger
	plugin string // plugin name for namespaced logging
}

// NewContext creates a new plugin context.
func NewContext(b *bot.Bot, ev *events.Bus) *Context {
	return &Context{
		bot:    b,
		world:  &World{},
		events: ev,
		log:    logger.Default(),
	}
}

// Bot returns the bot instance.
func (c *Context) Bot() *bot.Bot {
	return c.bot
}

// World returns the world instance.
func (c *Context) World() *World {
	return c.world
}

// Events returns the event bus.
func (c *Context) Events() *events.Bus {
	return c.events
}

// Logger returns the logger.
func (c *Context) Logger() *logger.Logger {
	return c.log
}

// WithPlugin returns a new context with a plugin-specific logger.
func (c *Context) WithPlugin(name string) *Context {
	return &Context{
		bot:    c.bot,
		world:  c.world,
		events: c.events,
		log:    c.log.With(name),
		plugin: name,
	}
}

// Plugin returns the plugin name for this context.
func (c *Context) Plugin() string {
	return c.plugin
}

// On subscribes to an event (convenience method).
func (c *Context) On(event string, handler func(...interface{})) *events.Subscription {
	return c.events.Subscribe(event, handler)
}

// Emit emits an event (convenience method).
func (c *Context) Emit(event string, data ...interface{}) {
	c.events.Emit(event, data...)
}

// EmitAsync emits an event asynchronously (convenience method).
func (c *Context) EmitAsync(event string, data ...interface{}) {
	c.events.EmitAsync(event, data...)
}

// Once waits for an event once (convenience method).
func (c *Context) Once(ctx context.Context, event string, predicate func(...interface{}) bool) ([]interface{}, error) {
	return c.events.Once(ctx, event, predicate)
}

// Info logs an info message (convenience method).
func (c *Context) Info(format string, args ...interface{}) {
	c.log.Info(format, args...)
}

// Debug logs a debug message (convenience method).
func (c *Context) Debug(format string, args ...interface{}) {
	c.log.Debug(format, args...)
}

// Warn logs a warning message (convenience method).
func (c *Context) Warn(format string, args ...interface{}) {
	c.log.Warn(format, args...)
}

// Error logs an error message (convenience method).
func (c *Context) Error(format string, args ...interface{}) {
	c.log.Error(format, args...)
}

// IsConnected returns true if the bot is connected (convenience method).
func (c *Context) IsConnected() bool {
	return c.bot != nil && c.bot.IsConnected()
}

// String returns a string representation of the context.
func (c *Context) String() string {
	if c.plugin != "" {
		return fmt.Sprintf("PluginContext{plugin=%s}", c.plugin)
	}
	return "PluginContext{}"
}

// Context extensions for World
// These methods will be expanded when World is fully implemented

// GetBlockAt returns the block at a position (placeholder).
func (w *World) GetBlockAt(x, y, z int) interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// TODO: Implement block lookup
	return nil
}

// SetBlockAt sets a block at a position (placeholder).
func (w *World) SetBlockAt(x, y, z int, block interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// TODO: Implement block setting
}

// GetEntity returns an entity by ID (placeholder).
func (w *World) GetEntity(id int32) interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// TODO: Implement entity lookup
	return nil
}

// AddEntity adds an entity to the world (placeholder).
func (w *World) AddEntity(id int32, entity interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// TODO: Implement entity addition
}

// RemoveEntity removes an entity from the world (placeholder).
func (w *World) RemoveEntity(id int32) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// TODO: Implement entity removal
}

// GetChunk returns a chunk at a position (placeholder).
func (w *World) GetChunk(x, z int) interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()
	// TODO: Implement chunk lookup
	return nil
}

// SetChunk sets a chunk at a position (placeholder).
func (w *World) SetChunk(x, z int, chunk interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// TODO: Implement chunk setting
}
