// Package conn manages network connections for Minecraft servers.
package conn

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Conn represents a Minecraft server connection.
// It handles TCP connection, packet framing, and encryption/compression.
type Conn struct {
	// Network connection
	conn net.Conn

	// Connection info
	host string
	port int

	// Encryption
	encryptor io.Writer
	decryptor io.Reader

	// Compression
	compressionThreshold int

	// Synchronization
	mu     sync.RWMutex
	closed bool

	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

// Config holds connection configuration.
type Config struct {
	Host                 string
	Port                 int
	Timeout              time.Duration
	CompressionThreshold int
}

// NewConn creates a new connection.
func NewConn(config Config) *Conn {
	ctx, cancel := context.WithCancel(context.Background())

	return &Conn{
		host:                 config.Host,
		port:                 config.Port,
		compressionThreshold: config.CompressionThreshold,
		ctx:                  ctx,
		cancel:               cancel,
	}
}

// Connect establishes a TCP connection to the server.
func (c *Conn) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return fmt.Errorf("already connected")
	}

	address := fmt.Sprintf("%s:%d", c.host, c.port)

	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	c.conn = conn
	return nil
}

// Close closes the connection.
func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.cancel()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}

	return nil
}

// IsClosed returns true if the connection is closed.
func (c *Conn) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// Read reads data from the connection.
func (c *Conn) Read(p []byte) (n int, err error) {
	c.mu.RLock()
	conn := c.conn
	reader := c.decryptor
	c.mu.RUnlock()

	if conn == nil {
		return 0, io.ErrClosedPipe
	}

	if reader != nil {
		return reader.Read(p)
	}

	// Set read deadline
	err = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		return 0, err
	}

	n, err = conn.Read(p)
	if err != nil {
		return n, err
	}

	return n, nil
}

// Write writes data to the connection.
func (c *Conn) Write(p []byte) (n int, err error) {
	c.mu.RLock()
	conn := c.conn
	writer := c.encryptor
	c.mu.RUnlock()

	if conn == nil {
		return 0, io.ErrClosedPipe
	}

	// Apply compression if needed
	data := p
	if c.compressionThreshold > 0 && len(data) > c.compressionThreshold {
		// TODO: Implement compression
	}

	// Apply encryption if needed
	if writer != nil {
		// TODO: Write through encryptor
		return conn.Write(data)
	}

	return conn.Write(data)
}

// SetReadDeadline sets the read deadline.
func (c *Conn) SetReadDeadline(t time.Time) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return io.ErrClosedPipe
	}

	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return io.ErrClosedPipe
	}

	return c.conn.SetWriteDeadline(t)
}

// SetDeadline sets both read and write deadlines.
func (c *Conn) SetDeadline(t time.Time) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return io.ErrClosedPipe
	}

	return c.conn.SetDeadline(t)
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return nil
	}

	return c.conn.RemoteAddr()
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return nil
	}

	return c.conn.LocalAddr()
}

// Context returns the connection's context.
// It is cancelled when the connection is closed.
func (c *Conn) Context() context.Context {
	return c.ctx
}

// EnableCompression enables packet compression.
func (c *Conn) EnableCompression(threshold int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.compressionThreshold = threshold
}

// SetEncryption sets the encryption cipher streams.
func (c *Conn) SetEncryption(encryptor io.Writer, decryptor io.Reader) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.encryptor = encryptor
	c.decryptor = decryptor
}
