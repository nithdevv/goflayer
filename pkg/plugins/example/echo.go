// Package example provides example plugins for goflayer.
package example

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/plugins/core"
)

// EchoPlugin is an example plugin that echoes chat messages.
// It demonstrates how to create a plugin that:
// - Depends on other plugins (chat)
// - Subscribes to events
// - Uses the plugin context
// - Provides custom API methods
type EchoPlugin struct {
	mu       sync.RWMutex
	ctx      *plugins.Context
	log      *logger.Logger
	enabled  bool
	prefix   string
	chat     *core.ChatPlugin
	messages []string
}

// Metadata returns the plugin metadata.
func (p *EchoPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "echo",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Example plugin that echoes chat messages",
		Dependencies: []string{"chat"},
	}
}

// OnLoad initializes the echo plugin.
func (p *EchoPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.enabled = true
	p.prefix = "[Echo] "
	p.messages = make([]string, 0, 100)

	p.log.Info("Echo plugin loaded")

	// Get dependency plugins
	loader := ctx.Bot().Events()
	if loader == nil {
		return fmt.Errorf("plugin loader not available")
	}

	// The chat plugin will be available through event handlers
	// No need to reference it directly here

	// Register event handlers
	ctx.On("chat_parsed", p.handleChat)

	// Register command handlers
	ctx.On("echo_enable", p.handleEnable)
	ctx.On("echo_disable", p.handleDisable)
	ctx.On("echo_clear", p.handleClear)
	ctx.On("echo_stats", p.handleStats)

	p.log.Info("Echo plugin listening for chat messages")
	return nil
}

// OnUnload cleans up the echo plugin.
func (p *EchoPlugin) OnUnload() error {
	p.log.Info("Echo plugin unloaded")
	return nil
}

// Enable enables echoing.
func (p *EchoPlugin) Enable() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.enabled = true
	p.log.Info("Echo enabled")
	p.ctx.Emit("echo_state_changed", true)
}

// Disable disables echoing.
func (p *EchoPlugin) Disable() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.enabled = false
	p.log.Info("Echo disabled")
	p.ctx.Emit("echo_state_changed", false)
}

// IsEnabled returns true if echoing is enabled.
func (p *EchoPlugin) IsEnabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

// SetPrefix sets the echo prefix.
func (p *EchoPlugin) SetPrefix(prefix string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.prefix = prefix
	p.log.Debug("Echo prefix set to: %s", prefix)
}

// GetPrefix returns the echo prefix.
func (p *EchoPlugin) GetPrefix() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.prefix
}

// GetMessages returns all echoed messages.
func (p *EchoPlugin) GetMessages() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	messages := make([]string, len(p.messages))
	copy(messages, p.messages)
	return messages
}

// GetMessageCount returns the number of echoed messages.
func (p *EchoPlugin) GetMessageCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.messages)
}

// Clear clears all echoed messages.
func (p *EchoPlugin) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.messages = make([]string, 0, 100)
	p.log.Info("Echo messages cleared")
	p.ctx.Emit("echo_cleared")
}

// Echo sends a message through chat.
func (p *EchoPlugin) Echo(message string) error {
	if !p.IsEnabled() {
		return fmt.Errorf("echo is disabled")
	}

	// Use the chat plugin to send the message
	return p.chat.Send(p.prefix + message)
}

// EchoTo echoes a message to a specific player.
func (p *EchoPlugin) EchoTo(player, message string) error {
	if !p.IsEnabled() {
		return fmt.Errorf("echo is disabled")
	}

	// Use the chat plugin to whisper
	return p.chat.Whisper(player, p.prefix+message)
}

// Event handlers

func (p *EchoPlugin) handleChat(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	msg, ok := args[0].(*core.ChatMessage)
	if !ok {
		return
	}

	// Only echo player chat messages
	if msg.Type != "chat" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.enabled {
		return
	}

	// Store message
	p.messages = append(p.messages, msg.Content)

	// Limit message history
	if len(p.messages) > 100 {
		p.messages = p.messages[len(p.messages)-100:]
	}

	// Echo the message
	response := fmt.Sprintf("%s%s: %s", p.prefix, msg.Username, msg.Content)
	p.log.Debug("Echoing: %s", response)

	// Send response through chat
	p.ctx.Emit("echo_message", response)
}

func (p *EchoPlugin) handleEnable(args ...interface{}) {
	p.Enable()
	p.chat.Send("Echo enabled")
}

func (p *EchoPlugin) handleDisable(args ...interface{}) {
	p.Disable()
	p.chat.Send("Echo disabled")
}

func (p *EchoPlugin) handleClear(args ...interface{}) {
	p.Clear()
	p.chat.Send("Echo messages cleared")
}

func (p *EchoPlugin) handleStats(args ...interface{}) {
	count := p.GetMessageCount()
	msg := fmt.Sprintf("Echo stats: %d messages recorded", count)
	p.chat.Send(msg)
}

// FilterMessages returns messages containing a substring.
func (p *EchoPlugin) FilterMessages(substring string) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	filtered := make([]string, 0)
	for _, msg := range p.messages {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(substring)) {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}

// String returns a string representation of the echo plugin.
func (p *EchoPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	state := "disabled"
	if p.enabled {
		state = "enabled"
	}

	return fmt.Sprintf("Echo{%s, messages=%d, prefix='%s'}",
		state, len(p.messages), p.prefix)
}
