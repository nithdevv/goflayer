package core

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// ChatMessage represents a parsed chat message.
type ChatMessage struct {
	Content    string
	Username   string
	Type       string // "chat", "system", "whisper", etc.
	Timestamp  int64
	MessageID  string
	Hash       string
	Signature string
}

// ChatPlugin handles chat send/receive with message parsing.
type ChatPlugin struct {
	mu         sync.RWMutex
	ctx        *plugins.Context
	log        *logger.Logger
	messageID  int
	patterns   map[string]*regexp.Regexp
}

// Metadata returns the plugin metadata.
func (p *ChatPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "chat",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Chat send/receive with message parsing",
		Dependencies: []string{},
	}
}

// OnLoad initializes the chat plugin.
func (p *ChatPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.messageID = 0
	p.patterns = make(map[string]*regexp.Regexp)

	// Compile common regex patterns
	p.patterns["username"] = regexp.MustCompile(`<[^>]+>`)
	p.patterns["whisper"] = regexp.MustCompile(`[^ ]+ whispers`)
	p.patterns["emote"] = regexp.MustCompile(`[^ ]+ [^ ]+ to you`)

	p.log.Info("Chat plugin loaded")

	// Register event handlers
	p.ctx.On("chat", p.handleChat)
	p.ctx.On("system_chat", p.handleSystemChat)

	return nil
}

// OnUnload cleans up the chat plugin.
func (p *ChatPlugin) OnUnload() error {
	p.log.Info("Chat plugin unloaded")
	return nil
}

// Send sends a chat message.
func (p *ChatPlugin) Send(message string) error {
	p.mu.Lock()
	p.messageID++
	id := p.messageID
	p.mu.Unlock()

	p.log.Debug("Sending chat message: %s", message)
	p.ctx.Emit("chat_send", message)

	// TODO: Send actual chat packet to server
	// This will be implemented when protocol packets are fully integrated

	p.ctx.Emit("chat_sent", id, message)
	return nil
}

// Say is an alias for Send.
func (p *ChatPlugin) Say(message string) error {
	return p.Send(message)
}

// Message sends a message (alias for Say).
func (p *ChatPlugin) Message(message string) error {
	return p.Send(message)
}

// Whisper sends a private message to a player.
func (p *ChatPlugin) Whisper(target, message string) error {
	msg := fmt.Sprintf("/tell %s %s", target, message)
	return p.Send(msg)
}

// Reply replies to the last message.
func (p *ChatPlugin) Reply(message string) error {
	// TODO: Track last message sender
	return p.Send(message)
}

// SetPattern compiles and stores a custom regex pattern.
func (p *ChatPlugin) SetPattern(name, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	p.mu.Lock()
	p.patterns[name] = re
	p.mu.Unlock()

	return nil
}

// Match checks if a message matches a pattern.
func (p *ChatPlugin) Match(patternName, message string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	re, exists := p.patterns[patternName]
	if !exists {
		return false
	}

	return re.MatchString(message)
}

// FindAll returns all matches for a pattern in a message.
func (p *ChatPlugin) FindAll(patternName, message string) [][]string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	re, exists := p.patterns[patternName]
	if !exists {
		return nil
	}

	return re.FindAllStringSubmatch(message, -1)
}

// ParseMessage parses a raw chat message into structured data.
func (p *ChatPlugin) ParseMessage(raw string) *ChatMessage {
	msg := &ChatMessage{
		Content:   raw,
		Timestamp: 0, // TODO: Get actual timestamp
	}

	// Parse message type and username
	p.parseMessageType(msg)

	return msg
}

// parseMessageType determines the message type and extracts username.
func (p *ChatPlugin) parseMessageType(msg *ChatMessage) {
	content := msg.Content

	// Check for whisper
	if strings.HasPrefix(content, "[") {
		if idx := strings.Index(content, " whispers"); idx > 0 {
			msg.Type = "whisper"
			// Extract username between [ and whispers
			start := strings.LastIndex(content[:idx], " ") + 1
			if start > 0 && start < idx {
				msg.Username = content[start:idx]
			}
			return
		}
	}

	// Check for regular chat with username prefix
	if strings.HasPrefix(content, "<") {
		if idx := strings.Index(content, "> "); idx > 0 {
			msg.Type = "chat"
			msg.Username = content[1:idx]
			msg.Content = content[idx+2:]
			return
		}
	}

	// System message
	msg.Type = "system"
}

// WaitFor waits for a message matching a predicate.
func (p *ChatPlugin) WaitFor(predicate func(*ChatMessage) bool) *ChatMessage {
	resultChan := make(chan *ChatMessage, 1)

	sub := p.ctx.On("chat_parsed", func(args ...interface{}) {
		if len(args) > 0 {
			if msg, ok := args[0].(*ChatMessage); ok {
				if predicate == nil || predicate(msg) {
					select {
					case resultChan <- msg:
					default:
					}
				}
			}
		}
	})
	defer sub.Unsubscribe()

	// TODO: Add timeout support
	msg := <-resultChan
	return msg
}

// WaitForString waits for a message containing specific text.
func (p *ChatPlugin) WaitForString(text string) *ChatMessage {
	return p.WaitFor(func(msg *ChatMessage) bool {
		return strings.Contains(msg.Content, text)
	})
}

// WaitForUsername waits for a message from a specific user.
func (p *ChatPlugin) WaitForUsername(username string) *ChatMessage {
	return p.WaitFor(func(msg *ChatMessage) bool {
		return msg.Username == username
	})
}

// AddPatternHandler adds a handler for messages matching a pattern.
func (p *ChatPlugin) AddPatternHandler(pattern string, handler func(*ChatMessage)) {
	sub := p.ctx.On("chat_parsed", func(args ...interface{}) {
		if len(args) > 0 {
			if msg, ok := args[0].(*ChatMessage); ok {
				if p.Match(pattern, msg.Content) {
					handler(msg)
				}
			}
		}
	})

	p.log.Debug("Added pattern handler for: %s", pattern)
	// Note: In production, you'd want to track this subscription for cleanup
	_ = sub
}

// Event handlers

func (p *ChatPlugin) handleChat(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	raw, ok := args[0].(string)
	if !ok {
		return
	}

	msg := p.ParseMessage(raw)
	p.log.Info("Chat received: [%s] %s: %s", msg.Type, msg.Username, msg.Content)

	p.ctx.Emit("chat_parsed", msg)
}

func (p *ChatPlugin) handleSystemChat(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	raw, ok := args[0].(string)
	if !ok {
		return
	}

	msg := &ChatMessage{
		Content:   raw,
		Type:      "system",
		Timestamp: 0,
	}

	p.log.Info("System chat: %s", msg.Content)
	p.ctx.Emit("chat_parsed", msg)
}

// String returns a string representation of the chat plugin.
func (p *ChatPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Chat{patterns=%d, messageID=%d}", len(p.patterns), p.messageID)
}
