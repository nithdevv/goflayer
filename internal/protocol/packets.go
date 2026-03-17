// Package protocol содержит определения пакетов Minecraft.
package protocol

import (
	"bytes"
)

// Direction represents packet direction.
type Direction int

const (
	Clientbound Direction = iota
	Serverbound
)

// State represents protocol state.
type State int

const (
	Handshaking State = iota
	Status
	Login
	Play
)

func (s State) String() string {
	switch s {
	case Handshaking:
		return "handshaking"
	case Status:
		return "status"
	case Login:
		return "login"
	case Play:
		return "play"
	default:
		return "unknown"
	}
}

// Packet IDs for Minecraft 1.20.1
const (
	// Handshake -> Server
	HandshakePacketID = 0x00

	// Login -> Server
	LoginStartPacketID = 0x00

	// Login -> Client
	LoginSuccessPacketID   = 0x02
	LoginDisconnectPacketID = 0x00

	// Play -> Client
	KeepAliveClientboundPacketID = 0x21

	// Play -> Server
	KeepAliveServerboundPacketID = 0x10
)

// Packet represents a Minecraft packet.
type Packet struct {
	ID        int32
	State     State
	Direction Direction
	Data      []byte
}

// NewPacket creates a new packet.
func NewPacket(id int32, state State, dir Direction) *Packet {
	return &Packet{
		ID:        id,
		State:     state,
		Direction: dir,
		Data:      nil,
	}
}

// ==================== HANDSHAKE ====================

// NewHandshakePacket creates a handshake packet.
func NewHandshakePacket(protocolVersion int, host string, port uint16, nextState int) *Packet {
	w := NewWriter()
	w.WriteVarInt(int32(protocolVersion))
	w.WriteString(host)
	w.WriteUint16(port)
	w.WriteVarInt(int32(nextState))

	pkt := NewPacket(HandshakePacketID, Handshaking, Serverbound)
	pkt.Data = w.Bytes()
	return pkt
}

// ==================== LOGIN ====================

// NewLoginStartPacket creates a login start packet.
func NewLoginStartPacket(username string) *Packet {
	w := NewWriter()
	w.WriteString(username)

	pkt := NewPacket(LoginStartPacketID, Login, Serverbound)
	pkt.Data = w.Bytes()
	return pkt
}

// ParseLoginSuccess parses a login success packet.
func ParseLoginSuccess(data []byte) (uuid, username string, err error) {
	r := NewReader(bytes.NewReader(data))
	uuid, err = r.ReadString()
	if err != nil {
		return "", "", err
	}
	username, err = r.ReadString()
	return uuid, username, err
}

// ==================== PLAY ====================

// ParseKeepAlive parses a keepalive packet.
func ParseKeepAlive(data []byte) (int64, error) {
	r := NewReader(bytes.NewReader(data))
	return r.ReadVarLong()
}

// NewKeepAlivePacket creates a keepalive packet.
func NewKeepAlivePacket(id int64) *Packet {
	w := NewWriter()
	w.WriteVarLong(id)

	pkt := NewPacket(KeepAliveServerboundPacketID, Play, Serverbound)
	pkt.Data = w.Bytes()
	return pkt
}
