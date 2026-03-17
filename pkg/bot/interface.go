// Package bot provides the Bot interface and core bot implementation.
package bot

import (
	"context"

	"github.com/nithdevv/goflayer/pkg/event"
	"github.com/nithdevv/goflayer/pkg/protocol"
)

// Bot represents a Minecraft bot.
type Bot interface {
	// Connect connects to the server.
	Connect(ctx context.Context) error

	// Disconnect disconnects from the server.
	Disconnect() error

	// On subscribes to an event.
	On(event string, handler func(...interface{})) event.Subscription

	// Emit emits an event.
	Emit(event string, data ...interface{})

	// Chat sends a chat message.
	Chat(message string) error

	// Events returns the event bus.
	Events() *event.Bus

	// Client returns the protocol client.
	Client() *protocol.Client

	// IsConnected returns true if connected to the server.
	IsConnected() bool

	// Username returns the bot's username.
	Username() string

	// Options returns the bot options.
	Options() *Options
}
