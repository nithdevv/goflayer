// Package conn управляет TCP соединением с reconect и backoff.
package conn

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/events"
)

// Stats holds connection statistics.
type Stats struct {
	BytesSent     uint64
	BytesReceived uint64
	PacketsSent   uint64
	PacketsRecv   uint64
	Reconnects    uint32
}

// Conn manages a TCP connection with auto-reconnect.
type Conn struct {
	mu         sync.RWMutex
	conn       net.Conn
	connected  atomic.Bool
	closed     atomic.Bool

	// Configuration
	host       string
	port       int
	timeout    time.Duration

	// Events
	events     *events.Bus

	// Stats
	stats      Stats

	// Logger
	log        *logger.Logger
}

// New creates a new connection manager.
func New(host string, port int, timeout time.Duration, ev *events.Bus) *Conn {
	log := logger.Default().With("conn")
	return &Conn{
		host:    host,
		port:    port,
		timeout: timeout,
		events:  ev,
		log:     log,
	}
}

// Connect establishes a connection.
func (c *Conn) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected.Load() {
		return fmt.Errorf("already connected")
	}

	c.log.Info("Connecting to %s:%d", c.host, c.port)

	dialer := net.Dialer{
		Timeout: c.timeout,
	}

	address := fmt.Sprintf("%s:%d", c.host, c.port)
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		c.log.Error("Failed to connect: %v", err)
		return fmt.Errorf("dial failed: %w", err)
	}

	c.conn = conn
	c.connected.Store(true)
	c.log.Info("Connected to %s", address)

	c.events.Emit("connected")

	return nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	if !c.closed.CompareAndSwap(false, true) {
		return nil // already closed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.log.Error("Error closing connection: %v", err)
		}
		c.conn = nil
	}

	c.connected.Store(false)
	c.log.Info("Connection closed")

	c.events.Emit("disconnected")

	return nil
}

// IsConnected returns true if connected.
func (c *Conn) IsConnected() bool {
	return c.connected.Load()
}

// Read reads data from the connection.
func (c *Conn) Read(p []byte) (n int, err error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	// Set read deadline
	if c.timeout > 0 {
		err = conn.SetReadDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return 0, err
		}
	}

	n, err = conn.Read(p)
	if n > 0 {
		atomic.AddUint64(&c.stats.BytesReceived, uint64(n))
	}

	if err != nil {
		c.log.Debug("Read error: %v", err)
		c.handleDisconnect()
		return n, err
	}

	return n, nil
}

// Write writes data to the connection.
func (c *Conn) Write(p []byte) (n int, err error) {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return 0, fmt.Errorf("not connected")
	}

	// Set write deadline
	if c.timeout > 0 {
		err = conn.SetWriteDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return 0, err
		}
	}

	n, err = conn.Write(p)
	if n > 0 {
		atomic.AddUint64(&c.stats.BytesSent, uint64(n))
	}

	if err != nil {
		c.log.Error("Write error: %v", err)
		c.handleDisconnect()
		return n, err
	}

	return n, nil
}

// handleDisconnect handles a disconnection.
func (c *Conn) handleDisconnect() {
	if c.connected.CompareAndSwap(true, false) {
		c.log.Warn("Connection lost")
		c.events.Emit("connection_lost")
	}
}

// SetDeadline sets both read and write deadlines.
func (c *Conn) SetDeadline(t time.Time) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	return c.conn.SetDeadline(t)
}

// RemoteAddr returns the remote address.
func (c *Conn) RemoteAddr() net.Addr {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return nil
	}

	return c.conn.RemoteAddr()
}

// Stats returns connection statistics.
func (c *Conn) Stats() Stats {
	return Stats{
		BytesSent:     atomic.LoadUint64(&c.stats.BytesSent),
		BytesReceived: atomic.LoadUint64(&c.stats.BytesReceived),
		PacketsSent:   atomic.LoadUint64(&c.stats.PacketsSent),
		PacketsRecv:   atomic.LoadUint64(&c.stats.PacketsRecv),
		Reconnects:    atomic.LoadUint32(&c.stats.Reconnects),
	}
}

// IncrementReconnects increments the reconnect counter.
func (c *Conn) IncrementReconnects() {
	atomic.AddUint32(&c.stats.Reconnects, 1)
}
