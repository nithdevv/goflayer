// Package chat implements chat message handling.
package chat

import (
	"strings"

	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
)

// Plugin handles chat messages.
type Plugin struct {
	*plugins.BasePlugin
}

// Message represents a chat message.
type Message struct {
	// The message text
	Text string

	// The sender's username (for player messages)
	Username string

	// The message type (chat, system, etc.)
	Type string

	// The raw JSON data
	JSON string
}

// NewPlugin creates a new chat plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("chat", "1.0.0")
	return &Plugin{
		BasePlugin: base,
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Subscribe to chat packet
	p.On("packet", p.handlePacket)

	return nil
}

// handlePacket handles incoming packets.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet := data[0].(*protocol.Packet)

	// Player Chat packet (1.19+)
	if packet.ID == 0x33 {
		p.handlePlayerChat(packet)
	}

	// System Chat packet (1.19+)
	if packet.ID == 0x24 {
		p.handleSystemChat(packet)
	}
}

// handlePlayerChat handles the Player Chat packet.
func (p *Plugin) handlePlayerChat(packet *protocol.Packet) {
	// TODO: Parse chat message from packet
	// For now, emit a basic event

	message := &Message{
		Type: "chat",
		Text: "chat message", // Placeholder
	}

	p.Emit("chat", message)
	p.Emit("message", message)
}

// handleSystemChat handles the System Chat packet.
func (p *Plugin) handleSystemChat(packet *protocol.Packet) {
	// TODO: Parse system message from packet

	message := &Message{
		Type: "system",
		Text: "system message", // Placeholder
	}

	p.Emit("chat", message)
	p.Emit("message", message)
}

// Send sends a chat message to the server.
func (p *Plugin) Send(message string) error {
	if message == "" {
		return nil
	}

	// TODO: Create and send chat packet
	// For now, just emit an event
	p.Emit("chat_send", message)
	return nil
}

// Say is an alias for Send.
func (p *Plugin) Say(message string) error {
	return p.Send(message)
}

// Whisper sends a private message to a player.
func (p *Plugin) Whisper(username, message string) error {
	return p.Send("/tell " + username + " " + message)
}

// OnPattern registers a handler for chat messages matching a pattern.
func (p *Plugin) OnPattern(pattern string, handler func(*Message)) {
	p.Bot().On("chat", func(data ...interface{}) {
		message := data[0].(*Message)

		if pattern == "" || strings.Contains(message.Text, pattern) {
			handler(message)
		}
	})
}

// OnUsername registers a handler for messages from a specific player.
func (p *Plugin) OnUsername(username string, handler func(*Message)) {
	p.Bot().On("chat", func(data ...interface{}) {
		message := data[0].(*Message)

		if message.Username == username {
			handler(message)
		}
	})
}

// OnRegex registers a handler for messages matching a regex pattern.
func (p *Plugin) OnRegex(pattern string, handler func(*Message, []string)) {
	// TODO: Implement regex matching
}

// AddPattern adds a chat pattern handler (for compatibility).
func (p *Plugin) AddPattern(pattern string, handler func(*Message)) {
	p.OnPattern(pattern, handler)
}
