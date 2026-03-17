package protocol

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/nithdevv/goflayer/pkg/event"
	"github.com/nithdevv/goflayer/pkg/net/conn"
	"github.com/nithdevv/goflayer/pkg/protocol/codec"
	"github.com/nithdevv/goflayer/pkg/protocol/states"
)

var (
	ErrAlreadyConnected = errors.New("already connected")
	ErrNotConnected     = errors.New("not connected")
	ErrInvalidState     = errors.New("invalid state transition")
)

// Client represents a Minecraft protocol client.
// It handles the low-level protocol communication with a Minecraft server.
type Client struct {
	// Connection
	conn *conn.Conn
	host string
	port int

	// Protocol state
	state        states.State
	version      string
	protocolVer  int
	stateMu      sync.RWMutex

	// Compression
	compressionThreshold int
	compressionMu        sync.RWMutex

	// Encryption
	encryptionEnabled bool

	// Event bus
	events *event.Bus

	// Packet handlers
	handlersMu sync.RWMutex
	handlers   map[states.State]map[int32][]*handlerWrapper
	nextID    uint64

	// Context
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Incoming packets
	incoming chan *Packet

	// Read loop control
	readLoopDone chan struct{}
}

// Config holds client configuration.
type Config struct {
	Host                 string
	Port                 int
	Version              string
	ProtocolVersion      int
	CompressionThreshold int
}

// NewClient creates a new protocol client.
func NewClient(config Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		host:                 config.Host,
		port:                 config.Port,
		version:              config.Version,
		protocolVer:          config.ProtocolVersion,
		state:                states.Handshaking,
		compressionThreshold: config.CompressionThreshold,
		events:               event.NewBus(),
		handlers:             make(map[states.State]map[int32][]*handlerWrapper),
		ctx:                  ctx,
		cancel:               cancel,
		incoming:             make(chan *Packet, 256),
		readLoopDone:         make(chan struct{}),
		nextID:               1,
	}

	// Initialize handler maps for all states
	c.handlers[states.Handshaking] = make(map[int32][]*handlerWrapper)
	c.handlers[states.Status] = make(map[int32][]*handlerWrapper)
	c.handlers[states.Login] = make(map[int32][]*handlerWrapper)
	c.handlers[states.Play] = make(map[int32][]*handlerWrapper)

	return c
}

// Connect establishes a connection to the server.
func (c *Client) Connect(ctx context.Context) error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if c.conn != nil {
		return ErrAlreadyConnected
	}

	// Create connection
	c.conn = conn.NewConn(conn.Config{
		Host:                 c.host,
		Port:                 c.port,
		CompressionThreshold: c.compressionThreshold,
	})

	// Connect to server
	if err := c.conn.Connect(ctx); err != nil {
		c.conn = nil
		return err
	}

	// Enable compression if threshold is set
	if c.compressionThreshold > 0 {
		c.conn.EnableCompression(c.compressionThreshold)
	}

	// Start read loop
	c.wg.Add(1)
	go c.readLoop()

	return nil
}

// Disconnect closes the connection to the server.
func (c *Client) Disconnect() error {
	c.stateMu.Lock()

	if c.conn == nil {
		c.stateMu.Unlock()
		return ErrNotConnected
	}

	// Cancel context to stop goroutines
	c.cancel()

	c.stateMu.Unlock()

	// Wait for read loop to finish
	<-c.readLoopDone
	c.wg.Wait()

	// Close connection
	c.stateMu.Lock()
	err := c.conn.Close()
	c.conn = nil
	c.stateMu.Unlock()

	return err
}

// State returns the current protocol state.
func (c *Client) State() states.State {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

// SetState changes the protocol state.
func (c *Client) SetState(state states.State) error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if !c.state.CanTransitionTo(state) {
		return ErrInvalidState
	}

	oldState := c.state
	c.state = state

	// Emit state change event
	c.events.Emit("state", c.state, oldState)

	return nil
}

// readLoop reads packets from the connection in a loop.
func (c *Client) readLoop() {
	defer c.wg.Done()
	defer close(c.readLoopDone)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			packet, err := c.readPacket()
			if err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) {
					// Emit error event
					c.events.Emit("error", err)
				}
				return
			}

			// Emit packet event
			c.events.Emit("packet", packet)

			// Send to incoming channel
			select {
			case c.incoming <- packet:
			case <-c.ctx.Done():
				return
			}
		}
	}
}

// readPacket reads a single packet from the connection.
func (c *Client) readPacket() (*Packet, error) {
	// Read packet length (VarInt)
	reader := codec.NewReader(c.conn)
	length, err := reader.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %w", err)
	}

	if length < 0 {
		return nil, fmt.Errorf("invalid packet length: %d", length)
	}

	// Sanity check: limit packet size to 2MB
	if length > 0x200000 {
		return nil, fmt.Errorf("packet too large: %d bytes", length)
	}

	// Read packet data
	data := make([]byte, length)
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet data: %w", err)
	}

	// Handle compression
	if c.compressionThreshold > 0 {
		data, err = c.decompressPacket(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress packet: %w", err)
		}
	}

	// Read packet ID
	buf := bytes.NewBuffer(data)
	packetReader := codec.NewReader(buf)
	packetID, err := packetReader.ReadVarInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %w", err)
	}

	// Get current state
	c.stateMu.RLock()
	state := c.state
	c.stateMu.RUnlock()

	// Create packet
	packet := &Packet{
		ID:        packetID,
		State:     state,
		Direction: states.Clientbound,
		Data:      data,
		Fields:    make(map[string]interface{}),
	}

	return packet, nil
}

// decompressPacket decompresses a packet if it's compressed.
func (c *Client) decompressPacket(data []byte) ([]byte, error) {
	buf := bytes.NewReader(data)
	reader := codec.NewReader(buf)

	// Read data length (VarInt)
	// 0 means uncompressed, >0 means compressed
	dataLength, err := reader.ReadVarInt()
	if err != nil {
		return nil, err
	}

	if dataLength == 0 {
		// Uncompressed, return remaining data
		remaining := make([]byte, buf.Len())
		_, err = buf.Read(remaining)
		return remaining, err
	}

	// Compressed, decompress
	compressedData := make([]byte, buf.Len())
	_, err = buf.Read(compressedData)
	if err != nil {
		return nil, err
	}

	// Decompress using zlib
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	decompressed := make([]byte, dataLength)
	_, err = io.ReadFull(r, decompressed)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

// Write writes a packet to the server.
func (c *Client) Write(packet *Packet) error {
	c.stateMu.RLock()
	conn := c.conn
	c.stateMu.RUnlock()

	if conn == nil {
		return ErrNotConnected
	}

	// Serialize packet
	data, err := c.serializePacket(packet)
	if err != nil {
		return fmt.Errorf("failed to serialize packet: %w", err)
	}

	// Apply compression if enabled
	if c.compressionThreshold > 0 {
		data, err = c.compressPacket(data)
		if err != nil {
			return fmt.Errorf("failed to compress packet: %w", err)
		}
	}

	// Write length prefix + data
	lengthBuf := make([]byte, 0)
	length := codec.VarIntSize(int32(len(data)))
	lengthBuf = make([]byte, length)
	codec.WriteVarInt(lengthBuf[:0], int32(len(data)))

	fullData := append(lengthBuf, data...)

	_, err = conn.Write(fullData)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	return nil
}

// serializePacket serializes a packet to bytes.
func (c *Client) serializePacket(packet *Packet) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := codec.NewWriter(buf)

	// Write packet ID
	if err := writer.WriteVarInt(packet.ID); err != nil {
		return nil, err
	}

	// Write packet data
	if _, err := buf.Write(packet.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// compressPacket compresses a packet if it's large enough.
func (c *Client) compressPacket(data []byte) ([]byte, error) {
	if len(data) < c.compressionThreshold {
		// Don't compress small packets
		// Write data length as 0 (uncompressed)
		buf := &bytes.Buffer{}
		writer := codec.NewWriter(buf)
		writer.WriteVarInt(0)
		writer.WriteBytes(data)
		return buf.Bytes(), nil
	}

	// Compress
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	// Write data length + compressed data
	buf := &bytes.Buffer{}
	writer := codec.NewWriter(buf)
	writer.WriteVarInt(int32(len(data)))
	writer.WriteBytes(compressed.Bytes())

	return buf.Bytes(), nil
}

// On registers a packet handler.
func (c *Client) On(packetID int32, handler PacketHandler) {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()

	c.stateMu.RLock()
	state := c.state
	c.stateMu.RUnlock()

	wrapper := &handlerWrapper{
		id:      c.nextID,
		handler: handler,
	}
	c.nextID++

	c.handlers[state][packetID] = append(c.handlers[state][packetID], wrapper)
}

// handlerWrapper wraps a handler with a unique ID.
type handlerWrapper struct {
	id      uint64
	handler PacketHandler
}

// Events returns the event bus.
func (c *Client) Events() *event.Bus {
	return c.events
}

// Incoming returns the incoming packet channel.
func (c *Client) Incoming() <-chan *Packet {
	return c.incoming
}
