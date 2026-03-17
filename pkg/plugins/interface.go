package plugins

import (
	"github.com/nithdevv/goflayer/pkg/event"
)

// Bot is a minimal interface for plugins to interact with the bot.
// This avoids circular dependencies.
type Bot interface {
	On(event string, handler func(...interface{})) event.Subscription
	Emit(event string, data ...interface{})
}

// BasePlugin provides a base implementation for plugins.
// It handles common functionality like event handler management.
type BasePlugin struct {
	bot         Bot
	name        string
	version     string
	subscriptions []event.Subscription
}

// NewBasePlugin creates a new base plugin.
func NewBasePlugin(name, version string) *BasePlugin {
	return &BasePlugin{
		name:        name,
		version:     version,
		subscriptions: make([]event.Subscription, 0),
	}
}

// Name returns the plugin name.
func (p *BasePlugin) Name() string {
	return p.name
}

// Version returns the plugin version.
func (p *BasePlugin) Version() string {
	return p.version
}

// Load is called when the plugin is loaded.
// Subclasses should override this method.
func (p *BasePlugin) Load(b Bot) error {
	p.bot = b
	return nil
}

// Unload is called when the plugin is unloaded.
// It unsubscribes all event handlers.
func (p *BasePlugin) Unload() error {
	for _, sub := range p.subscriptions {
		sub.Unsubscribe()
	}
	p.subscriptions = p.subscriptions[:0]
	return nil
}

// Bot returns the bot instance.
func (p *BasePlugin) Bot() Bot {
	return p.bot
}

// On subscribes to an event and tracks the subscription.
func (p *BasePlugin) On(event string, handler func(...interface{})) event.Subscription {
	sub := p.bot.On(event, handler)
	p.subscriptions = append(p.subscriptions, sub)
	return sub
}

// Emit emits an event.
func (p *BasePlugin) Emit(event string, data ...interface{}) {
	p.bot.Emit(event, data...)
}
