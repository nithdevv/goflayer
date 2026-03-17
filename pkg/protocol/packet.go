package protocol

import (
	"bytes"
	"time"
)

// Packet represents a Minecraft protocol packet.
//
// Packets are the basic unit of communication in the Minecraft protocol.
// Each packet has a name, data, and optionally the raw buffer.
//
// Example:
//
//	packet := Packet{
//	    Name: "chat",
//	    Data: map[string]interface{}{
//	        "message": "Hello, World!",
//	    },
//	}
type Packet struct {
	// Name is the packet identifier (e.g., "chat", "block_change").
	Name string

	// Data contains the parsed packet data as key-value pairs.
	// The structure depends on the packet type.
	// Example: {"message": "hello", "position": 0}
	Data map[string]interface{}

	// Buffer contains the raw packet bytes (after length prefix).
	// This is useful for debugging or when you need the raw data.
	Buffer []byte

	// Length is the total packet length including the length prefix.
	Length int

	// ID is the packet ID for the current protocol state and direction.
	// Packet IDs vary by protocol version and state (Handshaking, Login, Play).
	ID int

	// Timestamp is when the packet was received or created.
	Timestamp time.Time

	// State is the protocol state when this packet was sent/received.
	State State
}

// NewPacket creates a new packet with the given name and data.
func NewPacket(name string, data map[string]interface{}) *Packet {
	return &Packet{
		Name:      name,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// NewPacketWithBuffer creates a new packet with the given name, data, and raw buffer.
func NewPacketWithBuffer(name string, data map[string]interface{}, buffer []byte) *Packet {
	return &Packet{
		Name:      name,
		Data:      data,
		Buffer:    buffer,
		Length:    len(buffer),
		Timestamp: time.Now(),
	}
}

// GetInt safely gets an int value from packet data.
// Returns the value and true if the key exists and is an int, false otherwise.
func (p *Packet) GetInt(key string) (int, bool) {
	val, ok := p.Data[key]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

// MustGetInt gets an int value from packet data or panics if not found.
func (p *Packet) MustGetInt(key string) int {
	val, ok := p.GetInt(key)
	if !ok {
		panic("packet: key not found or not an int: " + key)
	}
	return val
}

// GetString safely gets a string value from packet data.
func (p *Packet) GetString(key string) (string, bool) {
	val, ok := p.Data[key]
	if !ok {
		return "", false
	}

	switch v := val.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	default:
		return "", false
	}
}

// MustGetString gets a string value from packet data or panics.
func (p *Packet) MustGetString(key string) string {
	val, ok := p.GetString(key)
	if !ok {
		panic("packet: key not found or not a string: " + key)
	}
	return val
}

// GetBool safely gets a bool value from packet data.
func (p *Packet) GetBool(key string) (bool, bool) {
	val, ok := p.Data[key]
	if !ok {
		return false, false
	}

	switch v := val.(type) {
	case bool:
		return v, true
	default:
		return false, false
	}
}

// MustGetBool gets a bool value from packet data or panics.
func (p *Packet) MustGetBool(key string) bool {
	val, ok := p.GetBool(key)
	if !ok {
		panic("packet: key not found or not a bool: " + key)
	}
	return val
}

// GetFloat64 safely gets a float64 value from packet data.
func (p *Packet) GetFloat64(key string) (float64, bool) {
	val, ok := p.Data[key]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

// MustGetFloat64 gets a float64 value from packet data or panics.
func (p *Packet) MustGetFloat64(key string) float64 {
	val, ok := p.GetFloat64(key)
	if !ok {
		panic("packet: key not found or not a float64: " + key)
	}
	return val
}

// GetBytes safely gets a byte slice value from packet data.
func (p *Packet) GetBytes(key string) ([]byte, bool) {
	val, ok := p.Data[key]
	if !ok {
		return nil, false
	}

	switch v := val.(type) {
	case []byte:
		return v, true
	case string:
		return []byte(v), true
	default:
		return nil, false
	}
}

// MustGetBytes gets a byte slice value from packet data or panics.
func (p *Packet) MustGetBytes(key string) []byte {
	val, ok := p.GetBytes(key)
	if !ok {
		panic("packet: key not found or not bytes: " + key)
	}
	return val
}

// Has checks if a key exists in packet data.
func (p *Packet) Has(key string) bool {
	_, ok := p.Data[key]
	return ok
}

// Set sets a value in packet data.
func (p *Packet) Set(key string, value interface{}) {
	if p.Data == nil {
		p.Data = make(map[string]interface{})
	}
	p.Data[key] = value
}

// Delete removes a key from packet data.
func (p *Packet) Delete(key string) {
	delete(p.Data, key)
}

// Clone creates a deep copy of the packet.
func (p *Packet) Clone() *Packet {
	clone := &Packet{
		Name:      p.Name,
		Length:    p.Length,
		ID:        p.ID,
		Timestamp: p.Timestamp,
		State:     p.State,
	}

	// Clone buffer
	if p.Buffer != nil {
		clone.Buffer = make([]byte, len(p.Buffer))
		copy(clone.Buffer, p.Buffer)
	}

	// Clone data
	if p.Data != nil {
		clone.Data = make(map[string]interface{})
		for k, v := range p.Data {
			clone.Data[k] = v
		}
	}

	return clone
}

// String returns a string representation of the packet.
func (p *Packet) String() string {
	return p.Name
}

// PacketDirection indicates the direction of packet flow.
type PacketDirection int

const (
	// ServerBound packets are sent from the client to the server.
	ServerBound PacketDirection = iota

	// ClientBound packets are sent from the server to the client.
	ClientBound
)

// String returns the string representation of the packet direction.
func (d PacketDirection) String() string {
	switch d {
	case ServerBound:
		return "ServerBound"
	case ClientBound:
		return "ClientBound"
	default:
		return "Unknown"
	}
}

// PacketHandler is a function that handles incoming packets.
type PacketHandler func(packet *Packet)

// PacketMiddleware wraps a packet handler with pre/post processing.
type PacketMiddleware func(next PacketHandler) PacketHandler

// PacketRegistry maps packet IDs to packet names for a given version and state.
type PacketRegistry struct {
	version       string
	serverBound   map[int]string
	clientBound   map[int]string
	nameToID      map[string]int
	middleware    []PacketMiddleware
}

// NewPacketRegistry creates a new packet registry for a given version.
func NewPacketRegistry(version string) *PacketRegistry {
	return &PacketRegistry{
		version:     version,
		serverBound: make(map[int]string),
		clientBound: make(map[int]string),
		nameToID:    make(map[string]int),
		middleware:  make([]PacketMiddleware, 0),
	}
}

// RegisterPacket registers a packet with its ID and name.
func (pr *PacketRegistry) RegisterPacket(direction PacketDirection, id int, name string) {
	switch direction {
	case ServerBound:
		pr.serverBound[id] = name
	case ClientBound:
		pr.clientBound[id] = name
	}
	pr.nameToID[name] = id
}

// GetPacketName returns the packet name for a given ID and direction.
func (pr *PacketRegistry) GetPacketName(direction PacketDirection, id int) (string, bool) {
	var nameMap map[int]string
	switch direction {
	case ServerBound:
		nameMap = pr.serverBound
	case ClientBound:
		nameMap = pr.clientBound
	}

	name, ok := nameMap[id]
	return name, ok
}

// GetPacketID returns the packet ID for a given name.
func (pr *PacketRegistry) GetPacketID(name string) (int, bool) {
	id, ok := pr.nameToID[name]
	return id, ok
}

// AddMiddleware adds a middleware to the packet processing chain.
func (pr *PacketRegistry) AddMiddleware(middleware PacketMiddleware) {
	pr.middleware = append(pr.middleware, middleware)
}

// ApplyMiddleware applies all middleware to a handler.
func (pr *PacketRegistry) ApplyMiddleware(handler PacketHandler) PacketHandler {
	// Apply middleware in reverse order so first added is outermost
	for i := len(pr.middleware) - 1; i >= 0; i-- {
		handler = pr.middleware[i](handler)
	}
	return handler
}

// PacketBuffer is a buffer for building and parsing packets.
type PacketBuffer struct {
	buffer *bytes.Buffer
	offset int
}

// NewPacketBuffer creates a new packet buffer.
func NewPacketBuffer() *PacketBuffer {
	return &PacketBuffer{
		buffer: bytes.NewBuffer(nil),
	}
}

// NewPacketBufferFromBytes creates a packet buffer from existing bytes.
func NewPacketBufferFromBytes(data []byte) *PacketBuffer {
	return &PacketBuffer{
		buffer: bytes.NewBuffer(data),
	}
}

// Bytes returns the buffer contents.
func (pb *PacketBuffer) Bytes() []byte {
	return pb.buffer.Bytes()
}

// Len returns the current length of the buffer.
func (pb *PacketBuffer) Len() int {
	return pb.buffer.Len()
}

// Reset clears the buffer.
func (pb *PacketBuffer) Reset() {
	pb.buffer.Reset()
	pb.offset = 0
}

// Write writes bytes to the buffer.
func (pb *PacketBuffer) Write(data []byte) (int, error) {
	return pb.buffer.Write(data)
}

// WriteByte writes a single byte to the buffer.
func (pb *PacketBuffer) WriteByte(b byte) error {
	return pb.buffer.WriteByte(b)
}

// Read reads bytes from the buffer.
func (pb *PacketBuffer) Read(data []byte) (int, error) {
	return pb.buffer.Read(data)
}

// ReadByte reads a single byte from the buffer.
func (pb *PacketBuffer) ReadByte() (byte, error) {
	return pb.buffer.ReadByte()
}

// Remaining returns the number of bytes remaining in the buffer.
func (pb *PacketBuffer) Remaining() int {
	return pb.buffer.Len()
}
