// Package protocol implements the Minecraft protocol.
package protocol

import (
	"github.com/nithdevv/goflayer/pkg/protocol/states"
)

// Packet represents a Minecraft protocol packet.
type Packet struct {
	// ID is the packet ID (varies by protocol state and version)
	ID int32

	// State is the protocol state this packet belongs to
	State states.State

	// Direction is the packet direction
	Direction states.PacketDirection

	// Data contains the raw packet data
	Data []byte

	// Fields contains parsed packet fields (when decoded)
	Fields map[string]interface{}
}

// NewPacket creates a new packet.
func NewPacket(id int32, state states.State, direction states.PacketDirection) *Packet {
	return &Packet{
		ID:        id,
		State:     state,
		Direction: direction,
		Fields:    make(map[string]interface{}),
	}
}

// PacketHandler handles incoming packets.
type PacketHandler func(*Packet)

// PacketDefinition defines how to encode/decode a packet.
type PacketDefinition struct {
	// ID is the packet ID
	ID int32

	// Name is the human-readable packet name
	Name string

	// State is the protocol state
	State states.State

	// Direction is the packet direction
	Direction states.PacketDirection

	// Fields defines the packet structure
	Fields []FieldDefinition
}

// FieldDefinition defines a field in a packet.
type FieldDefinition struct {
	// Name is the field name
	Name string

	// Type is the field type
	Type FieldType

	// Optional annotations
	Annotations map[string]interface{}
}

// FieldType represents the type of a field.
type FieldType int

const (
	// Field types
	TypeBool FieldType = iota
	TypeByte
	TypeUByte
	TypeShort
	TypeUShort
	TypeInt
	TypeLong
	TypeFloat
	TypeDouble
	TypeString
	TypeChat
	TypeIdentifier
	TypeVarInt
	TypeVarLong
	TypeEntityMetadata
	TypeSlot
	TypeBoolean
	TypeUUID
	TypePosition
	TypeNBT
	TypeNBTPath
	TypeParticle
	TypeVillagerData
	TypeItemID
	TypeBlockID
	TypeParticleID
	TypePotionID
	TypeRecipe
	TypeSoundID
	TypeRotation
	TypeBlockFace
)

// String returns the string representation of the field type.
func (t FieldType) String() string {
	switch t {
	case TypeBool:
		return "bool"
	case TypeByte:
		return "byte"
	case TypeUByte:
		return "ubyte"
	case TypeShort:
		return "short"
	case TypeUShort:
		return "ushort"
	case TypeInt:
		return "int"
	case TypeLong:
		return "long"
	case TypeFloat:
		return "float"
	case TypeDouble:
		return "double"
	case TypeString:
		return "string"
	case TypeChat:
		return "chat"
	case TypeIdentifier:
		return "identifier"
	case TypeVarInt:
		return "varint"
	case TypeVarLong:
		return "varlong"
	case TypeEntityMetadata:
		return "entity_metadata"
	case TypeSlot:
		return "slot"
	case TypeBoolean:
		return "boolean"
	case TypeUUID:
		return "uuid"
	case TypePosition:
		return "position"
	case TypeNBT:
		return "nbt"
	case TypeNBTPath:
		return "nbtpath"
	case TypeParticle:
		return "particle"
	case TypeVillagerData:
		return "villager_data"
	case TypeItemID:
		return "item_id"
	case TypeBlockID:
		return "block_id"
	case TypeParticleID:
		return "particle_id"
	case TypePotionID:
		return "potion_id"
	case TypeRecipe:
		return "recipe"
	case TypeSoundID:
		return "sound_id"
	case TypeRotation:
		return "rotation"
	case TypeBlockFace:
		return "block_face"
	default:
		return "unknown"
	}
}
