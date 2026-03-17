// Package bot предоставляет главный Bot для Minecraft.
package bot

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/internal/conn"
	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/internal/protocol"
	"github.com/nithdevv/goflayer/internal/session"
	"github.com/nithdevv/goflayer/internal/types"
	"github.com/nithdevv/goflayer/internal/worker"
	"github.com/nithdevv/goflayer/pkg/events"
)

// Bot represents a Minecraft bot.
type Bot struct {
	mu     sync.RWMutex
	config types.BotConfig

	// Components
	conn    *conn.Conn
	session *session.Manager
	worker  *worker.Pool
	events  *events.Bus

	// Lifecycle
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	running  bool

	// Logger
	log *logger.Logger
}

// New creates a new bot with the given configuration.
func New(config types.BotConfig) (*Bot, error) {
	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Initialize logger
	logger.Init(os.Stdout, logger.INFO)

	log := logger.Default().With("bot")
	log.Info("Initializing bot...")
	log.Debug("Server: %s:%d", config.Server.Host, config.Server.Port)
	log.Debug("Username: %s", config.Player.Username)

	ctx, cancel := context.WithCancel(context.Background())

	// Create event bus
	ev := events.NewBus()

	// Create connection manager
	c := conn.New(
		config.Server.Host,
		config.Server.Port,
		config.ReadTimeout,
		ev,
	)

	// Create session manager
	sess := session.New(c, config.Player.Username, config.Server.Protocol, ev)

	// Create worker pool (packet processor wraps session)
	wp := worker.New(config.WorkerCount, &packetProcessor{session: sess})

	bot := &Bot{
		config:  config,
		conn:    c,
		session: sess,
		worker:  wp,
		events:  ev,
		ctx:     ctx,
		cancel:  cancel,
		log:     log,
	}

	log.Info("Bot initialized successfully")
	return bot, nil
}

// validateConfig validates the bot configuration.
func validateConfig(cfg *types.BotConfig) error {
	if cfg.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.Player.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}

// Connect connects the bot to the server.
func (b *Bot) Connect(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return fmt.Errorf("already running")
	}

	b.log.Info("Connecting to server...")

	// Start worker pool
	if err := b.worker.Start(b.ctx); err != nil {
		return fmt.Errorf("failed to start worker pool: %w", err)
	}

	// Connect to server
	if err := b.conn.Connect(ctx); err != nil {
		b.worker.Stop()
		return fmt.Errorf("connection failed: %w", err)
	}

	// Start packet reader
	b.wg.Add(1)
	go b.packetReader()

	// Start session (handshake + login)
	if err := b.session.Start(
		b.config.Server.Host,
		uint16(b.config.Server.Port),
	); err != nil {
		b.conn.Close()
		b.worker.Stop()
		return fmt.Errorf("session failed: %w", err)
	}

	b.running = true
	b.log.Info("Bot connected and running")
	b.events.Emit("connected")

	return nil
}

// Disconnect disconnects the bot from the server.
func (b *Bot) Disconnect() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return nil
	}

	b.log.Info("Disconnecting...")

	// Signal shutdown
	b.cancel()

	// Close session
	b.session.Close()

	// Stop worker pool
	b.worker.Stop()

	// Close connection
	b.conn.Close()

	// Wait for goroutines
	b.wg.Wait()

	b.running = false
	b.log.Info("Disconnected")
	b.events.Emit("disconnected")

	return nil
}

// packetReader reads packets from the connection.
func (b *Bot) packetReader() {
	defer b.wg.Done()
	b.log.Debug("Packet reader started")

	defer func() {
		if r := recover(); r != nil {
			b.log.Error("Packet reader panic: %v", r)
		}
	}()

	for {
		select {
		case <-b.ctx.Done():
			b.log.Debug("Packet reader stopped by context")
			return

		default:
			pkt, err := b.readPacket()
			if err != nil {
				if !b.isContextCancelled() {
					b.log.Error("Failed to read packet: %v", err)
					b.handleConnectionError(err)
				}
				return
			}

			// Submit to worker pool
			if err := b.worker.Submit(pkt); err != nil {
				b.log.Error("Failed to submit packet: %v", err)
			}
		}
	}
}

// readPacket reads a single packet from the connection.
func (b *Bot) readPacket() (*protocol.Packet, error) {
	// Read packet length (VarInt)
	r := protocol.NewReader(b.conn)
	length, err := r.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}

	if length < 0 {
		return nil, fmt.Errorf("invalid length: %d", length)
	}

	if length > 0x200000 { // 2MB limit
		return nil, fmt.Errorf("packet too large: %d", length)
	}

	// Read packet data
	data := make([]byte, length)
	_, err = b.conn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("read data: %w", err)
	}

	// Parse packet ID
	packetReader := protocol.NewReader(bytes.NewReader(data))
	packetID, err := packetReader.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("read packet ID: %w", err)
	}

	pkt := &protocol.Packet{
		ID:    packetID,
		Data:  data,
		State: b.session.GetState(),
	}

	b.log.Debug("Read packet 0x%02X (%d bytes)", pkt.ID, length)

	return pkt, nil
}

// handleConnectionError handles a connection error.
func (b *Bot) handleConnectionError(err error) {
	b.events.Emit("error", err)

	if b.config.EnableReconnect {
		b.log.Warn("Connection lost, attempting to reconnect...")
		b.wg.Add(1)
		go b.reconnect()
	}
}

// reconnect attempts to reconnect to the server.
func (b *Bot) reconnect() {
	defer b.wg.Done()

	delay := b.config.ReconnectDelay
	attempts := 0

	for attempts < b.config.MaxReconnects {
		if b.isContextCancelled() {
			return
		}

		attempts++
		b.log.Info("Reconnect attempt %d/%d", attempts, b.config.MaxReconnects)

		time.Sleep(delay)

		ctx, cancel := context.WithTimeout(b.ctx, b.config.ConnectTimeout)
		err := b.conn.Connect(ctx)
		cancel()

		if err == nil {
			b.log.Info("Reconnected successfully")
			b.conn.IncrementReconnects()

			// Restart session
			if err := b.session.Start(
				b.config.Server.Host,
				uint16(b.config.Server.Port),
			); err != nil {
				b.log.Error("Failed to restart session: %v", err)
				b.conn.Close()
			} else {
				return
			}
		}

		delay = time.Duration(float64(delay) * b.config.ReconnectBackoff)
	}

	b.log.Error("Failed to reconnect after %d attempts", attempts)
	b.events.Emit("reconnect_failed")
}

// isContextCancelled checks if context is cancelled.
func (b *Bot) isContextCancelled() bool {
	select {
	case <-b.ctx.Done():
		return true
	default:
		return false
	}
}

// IsConnected returns true if the bot is connected.
func (b *Bot) IsConnected() bool {
	return b.conn.IsConnected()
}

// On subscribes to an event (convenience method).
func (b *Bot) On(event string, handler func(...interface{})) *events.Subscription {
	return b.events.Subscribe(event, handler)
}

// Emit emits an event (convenience method).
func (b *Bot) Emit(event string, data ...interface{}) {
	b.events.Emit(event, data...)
}

// Events returns the event bus.
func (b *Bot) Events() *events.Bus {
	return b.events
}

// Config returns the bot configuration.
func (b *Bot) Config() *types.BotConfig {
	return &b.config
}

// packetProcessor processes packets by forwarding to session.
type packetProcessor struct {
	session *session.Manager
}

func (p *packetProcessor) Process(pkt *protocol.Packet) error {
	p.session.HandlePacket(pkt)
	return nil
}
