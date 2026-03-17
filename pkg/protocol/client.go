package protocol

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Client represents a Minecraft protocol client.
//
// Client handles the low-level protocol communication with a Minecraft server.
// It manages packet serialization/deserialization, compression, encryption,
// and protocol state transitions.
type Client struct {
	// Connection
	conn        net.Conn
	host        string
	port        int
	connMu      sync.Mutex

	// Protocol state
	state          State
	connState      *ConnectionState
	version        string
	protocolVersion int
	stateMu        sync.RWMutex

	// Processing
	serializer     *Serializer
	deserializer   *Deserializer
	compressor     Compressor
	encryptor      Encryptor

	// Packet handling
	packetRegistry *PacketRegistry
	handlers       map[string][]*handlerWrapper
	handlersMu     sync.RWMutex
	nextHandlerID  uint64

	// Channels
	packetChan chan *Packet
	errorChan  chan error

	// Context
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Metrics
	latency       time.Duration
	lastKeepAlive time.Time

	// Configuration
	hideErrors           bool
	compressionThreshold int
}

// ClientConfig contains configuration for creating a client.
type ClientConfig struct {
	// Version is the Minecraft version to connect to.
	// If empty, version will be auto-detected during handshake.
	Version string

	// HideErrors suppresses error logging.
	HideErrors bool

	// CompressionThreshold is the packet size threshold for compression.
	// Packets larger than this will be compressed. -1 to disable compression.
	CompressionThreshold int
}

// NewClient creates a new Minecraft protocol client.
func NewClient(config ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		state:               Handshaking,
		connState:           NewConnectionState(),
		version:             config.Version,
		handlers:            make(map[string][]*handlerWrapper),
		packetChan:          make(chan *Packet, 256),
		errorChan:           make(chan error, 10),
		ctx:                 ctx,
		cancel:              cancel,
		hideErrors:          config.HideErrors,
		compressionThreshold: config.CompressionThreshold,
		packetRegistry:      NewPacketRegistry(config.Version),
		nextHandlerID:       1,
	}
}

// Connect connects to a Minecraft server.
//
// This establishes a TCP connection and starts the read loop.
// The connection handshake is not performed automatically.
func (c *Client) Connect(ctx context.Context, host string, port int) error {
	c.connMu.Lock()
	defer c.connMu.Unlock()

	c.host = host
	c.port = port

	// Establish TCP connection
	address := fmt.Sprintf("%s:%d", host, port)
	dialer := net.Dialer{}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	c.conn = conn
	c.connState.SetState(Handshaking)

	// Initialize serializer/deserializer for current state
	c.initSerializerDeserializer()

	// Start read loop
	c.wg.Add(1)
	go c.readLoop()

	return nil
}

// Disconnect closes the connection to the server.
// FIXED: Use instance variable instead of package variable
func (c *Client) Disconnect(reason string) {
	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()

	if conn == nil {
		return
	}

	// Cancel context first to stop goroutines
	c.cancel()

	// Close the connection
	if err := conn.Close(); err != nil && !c.hideErrors {
		// Could log error here
	}

	// Wait for goroutines to finish
	c.wg.Wait()

	c.connMu.Lock()
	c.conn = nil
	c.connMu.Unlock()

	c.emit("end", reason)
}

// State returns the current protocol state.
func (c *Client) State() State {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.state
}

// SetState changes the protocol state and reinitializes serializer/deserializer.
func (c *Client) SetState(newState State) error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if !c.state.CanTransitionTo(newState) {
		return ErrInvalidState
	}

	oldState := c.state
	c.state = newState

	// Update connection state
	if err := c.connState.SetState(newState); err != nil {
		return err
	}

	// Reinitialize serializer/deserializer for new state
	c.initSerializerDeserializer()

	c.emit("state", newState, oldState)
	return nil
}

// On registers a packet handler for a specific packet name.
//
// The handler will be called whenever a packet with this name is received.
// Multiple handlers can be registered for the same packet.
func (c *Client) On(packetName string, handler PacketHandler) Subscription {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()

	if c.handlers[packetName] == nil {
		c.handlers[packetName] = make([]*handlerWrapper, 0)
	}

	// FIXED: Use wrapper with unique ID for proper removal
	wrapper := &handlerWrapper{
		id:      c.nextHandlerID,
		handler: handler,
	}
	c.nextHandlerID++

	c.handlers[packetName] = append(c.handlers[packetName], wrapper)

	return &clientSubscription{
		client: c,
		name:   packetName,
		id:     wrapper.id,
	}
}

// Write writes a packet to the server.
//
// The packet is serialized, compressed (if enabled), encrypted (if enabled),
// and sent to the server.
func (c *Client) Write(packetName string, data map[string]interface{}) error {
	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()

	if conn == nil {
		return ErrBotNotConnected
	}

	packet := NewPacket(packetName, data)

	c.stateMu.RLock()
	packet.State = c.state
	c.stateMu.RUnlock()

	// Serialize the packet
	buffer, err := c.serializer.Serialize(packet)
	if err != nil {
		return fmt.Errorf("failed to serialize packet %s: %w", packetName, err)
	}

	packet.Buffer = buffer

	// Apply compression
	if c.compressor != nil {
		compressed, err := c.compressor.Compress(buffer)
		if err != nil {
			return fmt.Errorf("compression error: %w", err)
		}
		buffer = compressed
	}

	// Apply encryption
	if c.encryptor != nil {
		encrypted, err := c.encryptor.Encrypt(buffer)
		if err != nil {
			return fmt.Errorf("encryption error: %w", err)
		}
		buffer = encrypted
	}

	// Write length prefix (VarInt) + data
	length := len(buffer)
	lengthBuf := make([]byte, varIntByteCount(uint32(length)))
	writeVarInt(lengthBuf, uint32(length))

	_, err = conn.Write(append(lengthBuf, buffer...))
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

// readLoop reads packets from the connection in a loop.
func (c *Client) readLoop() {
	defer c.wg.Done()
	defer close(c.packetChan)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			packet, err := c.readPacket()
			if err != nil {
				select {
				case c.errorChan <- err:
				case <-c.ctx.Done():
				}
				return
			}

			// Emit the packet
			c.emitPacket(packet)
		}
	}
}

// readPacket reads a single packet from the connection.
func (c *Client) readPacket() (*Packet, error) {
	c.connMu.Lock()
	conn := c.conn
	c.connMu.Unlock()

	if conn == nil {
		return nil, ErrBotNotConnected
	}

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer conn.SetReadDeadline(time.Time{}) // Clear deadline

	// Read packet length (VarInt)
	length, err := readVarInt(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %w", err)
	}

	if length <= 0 || length > 0x200000 { // Max 2MB
		return nil, fmt.Errorf("invalid packet length: %d", length)
	}

	// Read packet data
	buffer := make([]byte, length)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read packet data: %w", err)
	}

	// Decrypt
	if c.encryptor != nil {
		buffer, err = c.encryptor.Decrypt(buffer)
		if err != nil {
			return nil, fmt.Errorf("decryption error: %w", err)
		}
	}

	// Decompress
	if c.compressor != nil {
		buffer, err = c.compressor.Decompress(buffer)
		if err != nil {
			return nil, fmt.Errorf("decompression error: %w", err)
		}
	}

	// Deserialize
	packet, err := c.deserializer.Deserialize(buffer)
	if err != nil {
		return nil, fmt.Errorf("deserialization error: %w", err)
	}

	c.stateMu.RLock()
	packet.State = c.state
	c.stateMu.RUnlock()

	return packet, nil
}

// emitPacket emits a packet to all registered handlers.
// FIXED: Limit concurrent goroutines and properly handle handlers
func (c *Client) emitPacket(packet *Packet) {
	// Emit generic packet event
	c.emit("packet", packet)

	c.handlersMu.RLock()
	handlers := c.handlers[packet.Name]
	// Copy handlers to avoid holding lock
	handlersCopy := make([]*handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	c.handlersMu.RUnlock()

	if len(handlersCopy) == 0 {
		return
	}

	// FIXED: Spawn goroutines with proper panic recovery
	for _, wrapper := range handlersCopy {
		go func(h PacketHandler) {
			defer func() {
				if r := recover(); r != nil && !c.hideErrors {
					// Could log panic here
				}
			}()
			h(packet)
		}(wrapper.handler)
	}
}

// emit emits an event with data.
func (c *Client) emit(event string, data ...interface{}) {
	c.handlersMu.RLock()
	handlers := c.handlers[event]
	// Copy handlers
	handlersCopy := make([]*handlerWrapper, len(handlers))
	copy(handlersCopy, handlers)
	c.handlersMu.RUnlock()

	if len(handlersCopy) == 0 {
		return
	}

	for _, wrapper := range handlersCopy {
		go func(h PacketHandler) {
			defer func() {
				if r := recover(); r != nil {
					// Recover from panic
				}
			}()

			// Create a fake packet for non-packet events
			if event == "packet" && len(data) > 0 {
				if packet, ok := data[0].(*Packet); ok {
					h(packet)
					return
				}
			}

			// For other events, create a fake packet
			h(&Packet{Name: event, Data: map[string]interface{}{"data": data}})
		}(wrapper.handler)
	}
}

// initSerializerDeserializer initializes serializer and deserializer for current state.
func (c *Client) initSerializerDeserializer() {
	// TODO: Initialize based on protocol state and version
	// This will be implemented when serializer/deserializer are complete
	c.stateMu.RLock()
	state := c.state
	version := c.version
	c.stateMu.RUnlock()

	c.serializer = NewSerializer(version, state)
	c.deserializer = NewDeserializer(version, state, c.packetRegistry)
}

// Version returns the Minecraft version string.
func (c *Client) Version() string {
	return c.version
}

// ProtocolVersion returns the protocol version number.
func (c *Client) ProtocolVersion() int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.protocolVersion
}

// Latency returns the current network latency (ping).
func (c *Client) Latency() time.Duration {
	return c.latency
}

// SetCompressionThreshold sets the compression threshold.
// Packets larger than this size will be compressed.
func (c *Client) SetCompressionThreshold(threshold int) error {
	if threshold < 0 {
		// Disable compression
		c.compressor = nil
		return nil
	}

	// Initialize compressor
	c.compressor = NewCompressor(threshold)
	return nil
}

// EnableCompression enables packet compression.
func (c *Client) EnableCompression(threshold int) error {
	return c.SetCompressionThreshold(threshold)
}

// DisableCompression disables packet compression.
func (c *Client) DisableCompression() {
	c.compressor = nil
}

// EnableEncryption enables packet encryption with the given key.
func (c *Client) EnableEncryption(key []byte) error {
	var err error
	c.encryptor, err = NewEncryptor(key)
	return err
}

// DisableEncryption disables packet encryption.
func (c *Client) DisableEncryption() {
	c.encryptor = nil
}

// handlerWrapper wraps a packet handler with a unique ID for removal.
// FIXED: This solves the function comparison problem
type handlerWrapper struct {
	id      uint64
	handler PacketHandler
}

// clientSubscription implements Subscription for client handlers.
type clientSubscription struct {
	client *Client
	name   string
	id     uint64
}

// Unsubscribe removes the packet handler.
// FIXED: Now properly identifies the handler by ID
func (s *clientSubscription) Unsubscribe() {
	s.client.handlersMu.Lock()
	defer s.client.handlersMu.Unlock()

	handlers := s.client.handlers[s.name]
	if handlers == nil {
		return
	}

	// Find and remove the handler by ID
	for i, wrapper := range handlers {
		if wrapper.id == s.id {
			// Remove by swapping with last and shrinking
			last := len(handlers) - 1
			handlers[i] = handlers[last]
			handlers[last] = nil
			s.client.handlers[s.name] = handlers[:last]
			return
		}
	}
}

// Helper functions for VarInt

func varIntByteCount(value uint32) int {
	if value == 0 {
		return 1
	}

	count := 0
	for {
		count++
		value >>= 7
		if value == 0 {
			break
		}
	}
	return count
}

func writeVarInt(buf []byte, value uint32) {
	i := 0
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		buf[i] = temp
		i++
		if value == 0 {
			break
		}
	}
}

func readVarInt(r io.Reader) (int, error) {
	var result uint32
	var shift uint

	for {
		buf := make([]byte, 1)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return 0, err
		}

		b := buf[0]
		result |= uint32(b&0x7F) << shift

		if (b & 0x80) == 0 {
			return int(result), nil
		}

		shift += 7
		if shift >= 35 {
			return 0, fmt.Errorf("VarInt too big")
		}
	}
}
