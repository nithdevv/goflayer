// Package play implements Minecraft 1.20.1 Play state packets.
package play

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/protocol"
)

// Packet direction constants
const (
	Clientbound = true
	Serverbound = false
)

// ==================== PACKET IDS ====================

// Clientbound Packet IDs (Server -> Client)
const (
	JoinGamePacketID                      = 0x26
	PluginMessageClientboundPacketID      = 0x00
	ServerDataPacketID                    = 0x27
	SynchronizePlayerPositionPacketID     = 0x2F
	EntitySpawnPacketID                   = 0x01
	EntitySpawnExperienceOrbPacketID      = 0x02
	EntitySpawnMobPacketID                = 0x03
	EntitySpawnPaintingPacketID           = 0x04
	EntitySpawnPlayerPacketID             = 0x05
	EntityAnimationClientboundPacketID    = 0x06
	EntityStatisticsPacketID              = 0x07
	EntityPositionPacketID                = 0x20
	EntityPositionAndRotationPacketID     = 0x21
	EntityVelocityPacketID                = 0x1A
	EntityEquipmentPacketID               = 0x4C
	EntityUpdateAttributesPacketID        = 0x4B
	EntityEffectPacketID                  = 0x4F
	EntityMetadataPacketID                = 0x44
	EntityTeleportPacketID                = 0x56
	EntityStatusPacketID                  = 0x1B
	EntityDamagePacketID                  = 0x1A
	EntityDeathPacketID                   = 0x1C
	SetExperiencePacketID                 = 0x48
	UpdateHealthPacketID                  = 0x49
	SetActionBarTextPacketID              = 0x1F
	SetTitleTextPacketID                  = 0x4D
	SetTimePacketID                       = 0x58
	SetSlotPacketID                       = 0x15
	SetItemsPacketID                      = 0x14
	OpenScreenPacketID                    = 0x2E
	CloseScreenClientboundPacketID        = 0x13
	ContainerSetContentPacketID           = 0x07
	ContainerSetDataPacketID              = 0x08
	ContainerSetSlotPacketID              = 0x09
	ContainerClosePacketID                = 0x0A
	BlockUpdatePacketID                   = 0x09
	ChunkDataPacketID                     = 0x1F
	UnloadChunkPacketID                   = 0x1D
	BlockChangedAckPacketID               = 0x0B
	SectionBlocksUpdatePacketID           = 0x3B
	GameEventPacketID                     = 0x1C
	LevelChunkPacketID                    = 0x20
	LevelChunkWithLightPacketID           = 0x20
	BlockEntityDataPacketID               = 0x07
	BlockEventPacketID                    = 0x06
	SoundPacketID                         = 0x5D
	ParticlePacketID                      = 0x23
	ExplosionPacketID                     = 0x01
	DisconnectClientboundPacketID         = 0x19
	ServerPlayerPacketID                  = 0x2E
	KeepAliveClientboundPacketID          = 0x21
)

// Serverbound Packet IDs (Client -> Server)
const (
	PluginMessageServerboundPacketID      = 0x0B
	ClientInformationPacketID             = 0x00
	ClientCommandPacketID                 = 0x02
	PlayerChatMessagePacketID             = 0x05
	PlayerPositionPacketID                = 0x12
	PlayerPositionAndLookServerboundID    = 0x12
	SetCreativeModeSlotPacketID           = 0x23
	ClickContainerPacketID                = 0x09
	SetHeldItemPacketID                   = 0x21
	SetPlayerPositionAndRotationPacketID  = 0x12
	UpdateSelectedSlotPacketID            = 0x13
	CloseContainerServerboundPacketID     = 0x0A
	KeepAliveServerboundPacketID          = 0x10
)

// ==================== GAME MODES ====================

type GameMode int8

const (
	SurvivalGameMode GameMode = iota
	CreativeGameMode
	AdventureGameMode
	SpectatorGameMode
)

// ==================== DIMENSIONS ====================

type Dimension int8

const (
	NetherDimension Dimension = iota
	OverworldDimension
	EndDimension
)

// ==================== GAMBLE EVENTS ====================

type GameEventType int

const (
	WelcomeGameEvent GameEventType = iota
	StartWaitingChunks
)

// ==================== HANDS ====================

type Hand int

const (
	MainHand Hand = iota
	OffHand
)

// ==================== INVENTORY ACTIONS ====================

type ClickType int

const (
	PickupClick ClickType = iota
	QuickMove
	Swap
	Clone
	Throw
	QuickCraft
	PickupAll
)

// ==================== CONTAINERS ====================

type ContainerType int

const (
	PlayerInventoryContainer ContainerType = iota
	Generic9x1
	Generic9x2
	Generic9x3
	Generic9x4
	Generic9x5
	Generic9x6
	Generic3x3
	Anvil
	Beacon
	BlastFurnace
	BrewingStand
	Crafting
	Furnace
	Grindstone
	Hopper
	Lectern
	Loom
	Merchant
	ShulkerBox
	Smoker
	Cartography
	Stonecutter
)

// ==================== PACKET INTERFACE ====================

// Packet represents a Minecraft Play state packet.
type Packet interface {
	ID() int32
	IsClientbound() bool
}

// ==================== CLIENTBOUND PACKETS ====================

// JoinGame is sent when the player joins the game.
type JoinGame struct {
	EntityID               int32
	IsHardcore             bool
	GameMode               GameMode
	PreviousGameMode       GameMode
	WorldNames             []string
	DimensionCodec         []byte
	Dimension              []byte
	WorldName              string
	HashedSeed             int64
	MaxPlayers             int32
	ViewDistance           int32
	SimulationDistance     int32
	ReducedDebugInfo       bool
	EnableRespawnScreen    bool
	DoLimitedCrafting      bool
	DimensionType          []byte
	DimensionName          string
}

func (p *JoinGame) ID() int32                              { return JoinGamePacketID }
func (p *JoinGame) IsClientbound() bool                     { return true }

func (p *JoinGame) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.IsHardcore, err = r.ReadBool()
	if err != nil {
		return err
	}

	gameMode, err := r.ReadUint8()
	if err != nil {
		return err
	}
	p.GameMode = GameMode(gameMode)

	previousGameMode, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.PreviousGameMode = GameMode(previousGameMode)

	worldCount, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.WorldNames = make([]string, worldCount)
	for i := int32(0); i < worldCount; i++ {
		p.WorldNames[i], err = r.ReadString()
		if err != nil {
			return err
		}
	}

	p.DimensionCodec, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.Dimension, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.WorldName, err = r.ReadString()
	if err != nil {
		return err
	}

	p.HashedSeed, err = r.ReadVarLong()
	if err != nil {
		return err
	}

	p.MaxPlayers, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.ViewDistance, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.SimulationDistance, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.ReducedDebugInfo, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.EnableRespawnScreen, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.DoLimitedCrafting, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.DimensionType, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.DimensionName, err = r.ReadString()
	return err
}

// PluginMessageClientbound is used for plugin messaging.
type PluginMessageClientbound struct {
	Channel string
	Data    []byte
}

func (p *PluginMessageClientbound) ID() int32          { return PluginMessageClientboundPacketID }
func (p *PluginMessageClientbound) IsClientbound() bool { return true }

func (p *PluginMessageClientbound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Channel, err = r.ReadString()
	if err != nil {
		return err
	}

	p.Data, err = r.ReadBytes()
	return err
}

// ServerData contains server information.
type ServerData struct {
	HasMotD             bool
	MotD                []byte
	HasIcon             bool
	Icon                []byte
	EnforcesSecureChat  bool
}

func (p *ServerData) ID() int32               { return ServerDataPacketID }
func (p *ServerData) IsClientbound() bool      { return true }

func (p *ServerData) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.HasMotD, err = r.ReadBool()
	if err != nil {
		return err
	}

	if p.HasMotD {
		p.MotD, err = r.ReadBytes()
		if err != nil {
			return err
		}
	}

	p.HasIcon, err = r.ReadBool()
	if err != nil {
		return err
	}

	if p.HasIcon {
		p.Icon, err = r.ReadBytes()
		if err != nil {
			return err
		}
	}

	p.EnforcesSecureChat, err = r.ReadBool()
	return err
}

// SynchronizePlayerPosition syncs the player's position.
type SynchronizePlayerPosition struct {
	X                    float64
	Y                    float64
	Z                    float64
	Yaw                  float32
	Pitch                float32
	Flags                byte
	TeleportID           int32
	DismountVehicle      bool
}

func (p *SynchronizePlayerPosition) ID() int32          { return SynchronizePlayerPositionPacketID }
func (p *SynchronizePlayerPosition) IsClientbound() bool { return true }

func (p *SynchronizePlayerPosition) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Yaw, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Pitch, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Flags, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.TeleportID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.DismountVehicle, err = r.ReadBool()
	return err
}

// EntitySpawn spawns an entity.
type EntitySpawn struct {
	EntityID int32
	UUID     []byte
	Type     int32
	X        float64
	Y        float64
	Z        float64
	Pitch    byte
	Yaw      byte
	Data     int32
}

func (p *EntitySpawn) ID() int32               { return EntitySpawnPacketID }
func (p *EntitySpawn) IsClientbound() bool      { return true }

func (p *EntitySpawn) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.UUID, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.Type, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Pitch, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.Yaw, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.Data, err = r.ReadInt32()
	return err
}

// EntitySpawnExperienceOrb spawns an experience orb.
type EntitySpawnExperienceOrb struct {
	EntityID int32
	X        float64
	Y        float64
	Z        float64
	Count    int16
}

func (p *EntitySpawnExperienceOrb) ID() int32          { return EntitySpawnExperienceOrbPacketID }
func (p *EntitySpawnExperienceOrb) IsClientbound() bool { return true }

func (p *EntitySpawnExperienceOrb) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Count, err = r.ReadInt16()
	return err
}

// EntitySpawnMob spawns a mob.
type EntitySpawnMob struct {
	EntityID  int32
	UUID      []byte
	Type      int32
	X         float64
	Y         float64
	Z         float64
	Yaw       int8
	Pitch     int8
	HeadPitch int8
	VelocityX int16
	VelocityY int16
	VelocityZ int16
}

func (p *EntitySpawnMob) ID() int32               { return EntitySpawnMobPacketID }
func (p *EntitySpawnMob) IsClientbound() bool      { return true }

func (p *EntitySpawnMob) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.UUID, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.Type, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	yaw, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Yaw = int8(yaw)

	pitch, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Pitch = int8(pitch)

	headPitch, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.HeadPitch = int8(headPitch)

	p.VelocityX, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.VelocityY, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.VelocityZ, err = r.ReadInt16()
	return err
}

// EntitySpawnPainting spawns a painting.
type EntitySpawnPainting struct {
	EntityID   int32
	UUID       []byte
	Variant    string
	Position   struct{ X, Y, Z int32 }
	Direction  byte
}

func (p *EntitySpawnPainting) ID() int32          { return EntitySpawnPaintingPacketID }
func (p *EntitySpawnPainting) IsClientbound() bool { return true }

func (p *EntitySpawnPainting) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.UUID, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.Variant, err = r.ReadString()
	if err != nil {
		return err
	}

	pos, err := r.ReadBlockPos()
	if err != nil {
		return err
	}
	p.Position.X, p.Position.Y, p.Position.Z = pos.X(), pos.Y(), pos.Z()

	p.Direction, err = r.ReadByte()
	return err
}

// EntitySpawnPlayer spawns a player.
type EntitySpawnPlayer struct {
	EntityID int32
	UUID     []byte
	X        float64
	Y        float64
	Z        float64
	Yaw      int8
	Pitch    int8
}

func (p *EntitySpawnPlayer) ID() int32          { return EntitySpawnPlayerPacketID }
func (p *EntitySpawnPlayer) IsClientbound() bool { return true }

func (p *EntitySpawnPlayer) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.UUID, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	yaw, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Yaw = int8(yaw)

	pitch, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Pitch = int8(pitch)

	return nil
}

// EntityAnimationClientbound plays an entity animation.
type EntityAnimationClientbound struct {
	EntityID int32
	Animation byte
}

func (p *EntityAnimationClientbound) ID() int32          { return EntityAnimationClientboundPacketID }
func (p *EntityAnimationClientbound) IsClientbound() bool { return true }

func (p *EntityAnimationClientbound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Animation, err = r.ReadByte()
	return err
}

// EntityStatistics sends statistics.
type EntityStatistics struct {
	Statistics []Statistic
}

type Statistic struct {
	CategoryID int32
	StatisticID int32
	Value      int64
}

func (p *EntityStatistics) ID() int32               { return EntityStatisticsPacketID }
func (p *EntityStatistics) IsClientbound() bool      { return true }

func (p *EntityStatistics) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Statistics = make([]Statistic, count)
	for i := int32(0); i < count; i++ {
		p.Statistics[i].CategoryID, err = r.ReadVarInt()
		if err != nil {
			return err
		}
		p.Statistics[i].StatisticID, err = r.ReadVarInt()
		if err != nil {
			return err
		}
		p.Statistics[i].Value, err = r.ReadVarLong()
		if err != nil {
			return err
		}
	}

	return nil
}

// EntityPosition updates entity position.
type EntityPosition struct {
	EntityID int32
	DeltaX   int16
	DeltaY   int16
	DeltaZ   int16
	OnGround bool
}

func (p *EntityPosition) ID() int32               { return EntityPositionPacketID }
func (p *EntityPosition) IsClientbound() bool      { return true }

func (p *EntityPosition) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.DeltaX, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.DeltaY, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.DeltaZ, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.OnGround, err = r.ReadBool()
	return err
}

// EntityPositionAndRotation updates entity position and rotation.
type EntityPositionAndRotation struct {
	EntityID int32
	DeltaX   int16
	DeltaY   int16
	DeltaZ   int16
	Yaw      byte
	Pitch    byte
	OnGround bool
}

func (p *EntityPositionAndRotation) ID() int32          { return EntityPositionAndRotationPacketID }
func (p *EntityPositionAndRotation) IsClientbound() bool { return true }

func (p *EntityPositionAndRotation) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.DeltaX, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.DeltaY, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.DeltaZ, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.Yaw, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.Pitch, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.OnGround, err = r.ReadBool()
	return err
}

// EntityVelocity updates entity velocity.
type EntityVelocity struct {
	EntityID int32
	VelocityX int16
	VelocityY int16
	VelocityZ int16
}

func (p *EntityVelocity) ID() int32               { return EntityVelocityPacketID }
func (p *EntityVelocity) IsClientbound() bool      { return true }

func (p *EntityVelocity) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.VelocityX, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.VelocityY, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.VelocityZ, err = r.ReadInt16()
	return err
}

// EntityEquipment updates entity equipment.
type EntityEquipment struct {
	EntityID  int32
	Equipment []EquipmentItem
}

type EquipmentItem struct {
	Slot int8
	Item []byte
}

func (p *EntityEquipment) ID() int32               { return EntityEquipmentPacketID }
func (p *EntityEquipment) IsClientbound() bool      { return true }

func (p *EntityEquipment) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Equipment = make([]EquipmentItem, count)
	for i := int32(0); i < count; i++ {
		slot, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.Equipment[i].Slot = int8(slot)

		p.Equipment[i].Item, err = r.ReadBytes()
		if err != nil {
			return err
		}
	}

	return nil
}

// EntityUpdateAttributes updates entity attributes.
type EntityUpdateAttributes struct {
	EntityID  int32
	Attributes []Attribute
}

type Attribute struct {
	Key    string
	Value  float64
	Modifiers []AttributeModifier
}

type AttributeModifier struct {
	UUID      []byte
	Amount    float64
	Operation byte
}

func (p *EntityUpdateAttributes) ID() int32          { return EntityUpdateAttributesPacketID }
func (p *EntityUpdateAttributes) IsClientbound() bool { return true }

func (p *EntityUpdateAttributes) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Attributes = make([]Attribute, count)
	for i := int32(0); i < count; i++ {
		p.Attributes[i].Key, err = r.ReadString()
		if err != nil {
			return err
		}

		p.Attributes[i].Value, err = r.ReadDouble()
		if err != nil {
			return err
		}

		modCount, err := r.ReadVarInt()
		if err != nil {
			return err
		}

		p.Attributes[i].Modifiers = make([]AttributeModifier, modCount)
		for j := int32(0); j < modCount; j++ {
			p.Attributes[i].Modifiers[j].UUID, err = r.ReadBytes()
			if err != nil {
				return err
			}

			p.Attributes[i].Modifiers[j].Amount, err = r.ReadDouble()
			if err != nil {
				return err
			}

			p.Attributes[i].Modifiers[j].Operation, err = r.ReadByte()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// EntityEffect applies an effect to an entity.
type EntityEffect struct {
	EntityID   int32
	EffectID   int8
	Amplifier  int8
	Duration   int32
	Flags      byte
}

func (p *EntityEffect) ID() int32               { return EntityEffectPacketID }
func (p *EntityEffect) IsClientbound() bool      { return true }

func (p *EntityEffect) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	effectID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.EffectID = int8(effectID)

	amplifier, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Amplifier = int8(amplifier)

	p.Duration, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Flags, err = r.ReadByte()
	return err
}

// EntityMetadata updates entity metadata.
type EntityMetadata struct {
	EntityID int32
	Metadata []byte
}

func (p *EntityMetadata) ID() int32               { return EntityMetadataPacketID }
func (p *EntityMetadata) IsClientbound() bool      { return true }

func (p *EntityMetadata) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Metadata, err = r.ReadBytes()
	return err
}

// EntityTeleport teleports an entity.
type EntityTeleport struct {
	EntityID int32
	X        float64
	Y        float64
	Z        float64
	Yaw      byte
	Pitch    byte
	OnGround bool
}

func (p *EntityTeleport) ID() int32               { return EntityTeleportPacketID }
func (p *EntityTeleport) IsClientbound() bool      { return true }

func (p *EntityTeleport) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Yaw, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.Pitch, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.OnGround, err = r.ReadBool()
	return err
}

// EntityStatus sends entity status.
type EntityStatus struct {
	EntityID int32
	Status   int8
}

func (p *EntityStatus) ID() int32               { return EntityStatusPacketID }
func (p *EntityStatus) IsClientbound() bool      { return true }

func (p *EntityStatus) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadInt32()
	if err != nil {
		return err
	}

	status, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Status = int8(status)
	return nil
}

// EntityDamage indicates entity damage.
type EntityDamage struct {
	EntityID    int32
	SourceID    int32
	SourceCause int32
	SourceDirec int32
	Amount      float32
}

func (p *EntityDamage) ID() int32               { return EntityDamagePacketID }
func (p *EntityDamage) IsClientbound() bool      { return true }

func (p *EntityDamage) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.SourceID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.SourceCause, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.SourceDirec, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Amount, err = r.ReadFloat32()
	return err
}

// EntityDeath indicates entity death.
type EntityDeath struct {
	EntityID int32
	Message  []byte
}

func (p *EntityDeath) ID() int32               { return EntityDeathPacketID }
func (p *EntityDeath) IsClientbound() bool      { return true }

func (p *EntityDeath) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.EntityID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Message, err = r.ReadBytes()
	return err
}

// SetExperience sets player experience.
type SetExperience struct {
	ExperienceBar float32
	Level         int32
	TotalPoints   int32
}

func (p *SetExperience) ID() int32               { return SetExperiencePacketID }
func (p *SetExperience) IsClientbound() bool      { return true }

func (p *SetExperience) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.ExperienceBar, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Level, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.TotalPoints, err = r.ReadVarInt()
	return err
}

// UpdateHealth updates player health.
type UpdateHealth struct {
	Health         float32
	Food           int32
	FoodSaturation float32
}

func (p *UpdateHealth) ID() int32               { return UpdateHealthPacketID }
func (p *UpdateHealth) IsClientbound() bool      { return true }

func (p *UpdateHealth) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Health, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Food, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.FoodSaturation, err = r.ReadFloat32()
	return err
}

// SetActionBarText sets the action bar text.
type SetActionBarText struct {
	Text []byte
}

func (p *SetActionBarText) ID() int32               { return SetActionBarTextPacketID }
func (p *SetActionBarText) IsClientbound() bool      { return true }

func (p *SetActionBarText) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Text, err = r.ReadBytes()
	return err
}

// SetTitleText sets the title text.
type SetTitleText struct {
	Text []byte
}

func (p *SetTitleText) ID() int32               { return SetTitleTextPacketID }
func (p *SetTitleText) IsClientbound() bool      { return true }

func (p *SetTitleText) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Text, err = r.ReadBytes()
	return err
}

// SetTime sets the world time.
type SetTime struct {
	AgeOfWorld  int64
	TimeOfDay   int64
}

func (p *SetTime) ID() int32               { return SetTimePacketID }
func (p *SetTime) IsClientbound() bool      { return true }

func (p *SetTime) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.AgeOfWorld, err = r.ReadVarLong()
	if err != nil {
		return err
	}

	p.TimeOfDay, err = r.ReadVarLong()
	return err
}

// SetSlot sets a slot in a container.
type SetSlot struct {
	WindowID int8
	Slot     int16
	ItemData []byte
}

func (p *SetSlot) ID() int32               { return SetSlotPacketID }
func (p *SetSlot) IsClientbound() bool      { return true }

func (p *SetSlot) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	p.Slot, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.ItemData, err = r.ReadBytes()
	return err
}

// SetItems sets multiple slots in a container.
type SetItems struct {
	WindowID int8
	Slots    []int16
	Items    [][]byte
}

func (p *SetItems) ID() int32               { return SetItemsPacketID }
func (p *SetItems) IsClientbound() bool      { return true }

func (p *SetItems) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Slots = make([]int16, count)
	p.Items = make([][]byte, count)

	for i := int32(0); i < count; i++ {
		p.Slots[i], err = r.ReadInt16()
		if err != nil {
			return err
		}

		p.Items[i], err = r.ReadBytes()
		if err != nil {
			return err
		}
	}

	return nil
}

// OpenScreen opens a container screen.
type OpenScreen struct {
	WindowID    int8
	WindowType  []byte
	WindowTitle []byte
}

func (p *OpenScreen) ID() int32               { return OpenScreenPacketID }
func (p *OpenScreen) IsClientbound() bool      { return true }

func (p *OpenScreen) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	p.WindowType, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.WindowTitle, err = r.ReadBytes()
	return err
}

// CloseScreenClientbound closes a container screen.
type CloseScreenClientbound struct {
	WindowID int8
}

func (p *CloseScreenClientbound) ID() int32               { return CloseScreenClientboundPacketID }
func (p *CloseScreenClientbound) IsClientbound() bool      { return true }

func (p *CloseScreenClientbound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	return nil
}

// ContainerSetContent sets the contents of a container.
type ContainerSetContent struct {
	WindowID    int8
	StateID     int32
	Items       [][]byte
	CarriedItem []byte
}

func (p *ContainerSetContent) ID() int32          { return ContainerSetContentPacketID }
func (p *ContainerSetContent) IsClientbound() bool { return true }

func (p *ContainerSetContent) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	p.StateID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Items = make([][]byte, count)
	for i := int32(0); i < count; i++ {
		p.Items[i], err = r.ReadBytes()
		if err != nil {
			return err
		}
	}

	p.CarriedItem, err = r.ReadBytes()
	return err
}

// ContainerSetData sets data for a container.
type ContainerSetData struct {
	WindowID int8
	Key      int16
	Value    int16
}

func (p *ContainerSetData) ID() int32          { return ContainerSetDataPacketID }
func (p *ContainerSetData) IsClientbound() bool { return true }

func (p *ContainerSetData) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	p.Key, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.Value, err = r.ReadInt16()
	return err
}

// ContainerSetSlot sets a slot in a container.
type ContainerSetSlot struct {
	WindowID int8
	StateID  int32
	Slot     int16
	ItemData []byte
}

func (p *ContainerSetSlot) ID() int32          { return ContainerSetSlotPacketID }
func (p *ContainerSetSlot) IsClientbound() bool { return true }

func (p *ContainerSetSlot) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	p.StateID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Slot, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.ItemData, err = r.ReadBytes()
	return err
}

// ContainerClose indicates a container was closed.
type ContainerClose struct {
	WindowID int8
}

func (p *ContainerClose) ID() int32               { return ContainerClosePacketID }
func (p *ContainerClose) IsClientbound() bool      { return true }

func (p *ContainerClose) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))

	windowID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.WindowID = int8(windowID)

	return nil
}

// BlockUpdate updates a single block.
type BlockUpdate struct {
	Position struct{ X, Y, Z int32 }
	BlockID  int32
}

func (p *BlockUpdate) ID() int32               { return BlockUpdatePacketID }
func (p *BlockUpdate) IsClientbound() bool      { return true }

func (p *BlockUpdate) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	pos, err := r.ReadBlockPos()
	if err != nil {
		return err
	}
	p.Position.X, p.Position.Y, p.Position.Z = pos.X(), pos.Y(), pos.Z()

	p.BlockID, err = r.ReadVarInt()
	return err
}

// ChunkData sends chunk data.
type ChunkData struct {
	X          int32
	Z          int32
	OldData    []byte
}

func (p *ChunkData) ID() int32               { return ChunkDataPacketID }
func (p *ChunkData) IsClientbound() bool      { return true }

func (p *ChunkData) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.OldData, err = r.ReadBytes()
	return err
}

// UnloadChunk unloads a chunk.
type UnloadChunk struct {
	X int32
	Z int32
}

func (p *UnloadChunk) ID() int32               { return UnloadChunkPacketID }
func (p *UnloadChunk) IsClientbound() bool      { return true }

func (p *UnloadChunk) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadInt32()
	return err
}

// BlockChangedAck acknowledges a block change.
type BlockChangedAck struct {
	Sequence int32
}

func (p *BlockChangedAck) ID() int32               { return BlockChangedAckPacketID }
func (p *BlockChangedAck) IsClientbound() bool      { return true }

func (p *BlockChangedAck) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Sequence, err = r.ReadVarInt()
	return err
}

// SectionBlocksUpdate updates multiple blocks in a section.
type SectionBlocksUpdate struct {
	ChunkCoordX int32
	ChunkCoordZ int32
	SectionY    int8
	Blocks      []BlockEntry
}

type BlockEntry struct {
	Offset uint16
	State  int32
}

func (p *SectionBlocksUpdate) ID() int32          { return SectionBlocksUpdatePacketID }
func (p *SectionBlocksUpdate) IsClientbound() bool { return true }

func (p *SectionBlocksUpdate) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.ChunkCoordX, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ChunkCoordZ, err = r.ReadInt32()
	if err != nil {
		return err
	}

	sectionY, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.SectionY = int8(sectionY)

	count, err := r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Blocks = make([]BlockEntry, count)
	for i := int32(0); i < count; i++ {
		offset, err := r.ReadUint16()
		if err != nil {
			return err
		}
		p.Blocks[i].Offset = offset

		p.Blocks[i].State, err = r.ReadVarInt()
		if err != nil {
			return err
		}
	}

	return nil
}

// GameEvent sends a game event.
type GameEvent struct {
	Event   GameEventType
	Value   float32
}

func (p *GameEvent) ID() int32               { return GameEventPacketID }
func (p *GameEvent) IsClientbound() bool      { return true }

func (p *GameEvent) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	eventType, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Event = GameEventType(eventType)

	p.Value, err = r.ReadFloat32()
	return err
}

// LevelChunk sends level chunk data.
type LevelChunk struct {
	X        int32
	Z        int32
	Data     []byte
}

func (p *LevelChunk) ID() int32               { return LevelChunkPacketID }
func (p *LevelChunk) IsClientbound() bool      { return true }

func (p *LevelChunk) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Data, err = r.ReadBytes()
	return err
}

// LevelChunkWithLight sends level chunk data with lighting.
type LevelChunkWithLight struct {
	X         int32
	Z         int32
	ChunkData []byte
	LightData []byte
}

func (p *LevelChunkWithLight) ID() int32          { return LevelChunkWithLightPacketID }
func (p *LevelChunkWithLight) IsClientbound() bool { return true }

func (p *LevelChunkWithLight) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ChunkData, err = r.ReadBytes()
	if err != nil {
		return err
	}

	p.LightData, err = r.ReadBytes()
	return err
}

// BlockEntityData sends block entity data.
type BlockEntityData struct {
	Position struct{ X, Y, Z int32 }
	Type     int32
	Data     []byte
}

func (p *BlockEntityData) ID() int32               { return BlockEntityDataPacketID }
func (p *BlockEntityData) IsClientbound() bool      { return true }

func (p *BlockEntityData) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	pos, err := r.ReadBlockPos()
	if err != nil {
		return err
	}
	p.Position.X, p.Position.Y, p.Position.Z = pos.X(), pos.Y(), pos.Z()

	p.Type, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	p.Data, err = r.ReadBytes()
	return err
}

// BlockEvent sends a block event.
type BlockEvent struct {
	Position struct{ X, Y, Z int32 }
	BlockID  int32
	ActionID int8
	ActionParam int8
}

func (p *BlockEvent) ID() int32               { return BlockEventPacketID }
func (p *BlockEvent) IsClientbound() bool      { return true }

func (p *BlockEvent) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	pos, err := r.ReadBlockPos()
	if err != nil {
		return err
	}
	p.Position.X, p.Position.Y, p.Position.Z = pos.X(), pos.Y(), pos.Z()

	p.BlockID, err = r.ReadInt32()
	if err != nil {
		return err
	}

	actionID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.ActionID = int8(actionID)

	actionParam, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.ActionParam = int8(actionParam)

	return nil
}

// Sound plays a sound.
type Sound struct {
	SoundID     int32
	SoundCategory int8
	X           int32
	Y           int32
	Z           int32
	Volume      float32
	Pitch       float32
	Seed        int64
}

func (p *Sound) ID() int32               { return SoundPacketID }
func (p *Sound) IsClientbound() bool      { return true }

func (p *Sound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.SoundID, err = r.ReadVarInt()
	if err != nil {
		return err
	}

	soundCat, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.SoundCategory = int8(soundCat)

	p.X, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Volume, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Pitch, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Seed, err = r.ReadVarLong()
	return err
}

// Particle displays particles.
type Particle struct {
	ParticleID   int32
	LongDistance bool
	X            float64
	Y            float64
	Z            float64
	OffsetX      float32
	OffsetY      float32
	OffsetZ      float32
	ParticleData float32
	Count        int32
	Data         []byte
}

func (p *Particle) ID() int32               { return ParticlePacketID }
func (p *Particle) IsClientbound() bool      { return true }

func (p *Particle) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.ParticleID, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.LongDistance, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.X, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat64()
	if err != nil {
		return err
	}

	p.OffsetX, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.OffsetY, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.OffsetZ, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.ParticleData, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Count, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Data, err = r.ReadBytes()
	return err
}

// Explosion creates an explosion.
type Explosion struct {
	X        float32
	Y        float32
	Z        float32
	Strength float32
	Records  []ExplosionRecord
	PlayerMotionX float32
	PlayerMotionY float32
	PlayerMotionZ float32
}

type ExplosionRecord struct {
	X int8
	Y int8
	Z int8
}

func (p *Explosion) ID() int32               { return ExplosionPacketID }
func (p *Explosion) IsClientbound() bool      { return true }

func (p *Explosion) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.X, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Y, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Z, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Strength, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	count, err := r.ReadInt32()
	if err != nil {
		return err
	}

	p.Records = make([]ExplosionRecord, count)
	for i := int32(0); i < count; i++ {
		x, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.Records[i].X = int8(x)

		y, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.Records[i].Y = int8(y)

		z, err := r.ReadByte()
		if err != nil {
			return err
		}
		p.Records[i].Z = int8(z)
	}

	p.PlayerMotionX, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.PlayerMotionY, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.PlayerMotionZ, err = r.ReadFloat32()
	return err
}

// DisconnectClientbound disconnects the player.
type DisconnectClientbound struct {
	Reason []byte
}

func (p *DisconnectClientbound) ID() int32               { return DisconnectClientboundPacketID }
func (p *DisconnectClientbound) IsClientbound() bool      { return true }

func (p *DisconnectClientbound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.Reason, err = r.ReadBytes()
	return err
}

// ServerPlayer contains server player information.
type ServerPlayer struct {
	PlayerID int32
}

func (p *ServerPlayer) ID() int32               { return ServerPlayerPacketID }
func (p *ServerPlayer) IsClientbound() bool      { return true }

func (p *ServerPlayer) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.PlayerID, err = r.ReadVarInt()
	return err
}

// KeepAliveClientbound is a keepalive packet from server.
type KeepAliveClientbound struct {
	KeepAliveID int64
}

func (p *KeepAliveClientbound) ID() int32               { return KeepAliveClientboundPacketID }
func (p *KeepAliveClientbound) IsClientbound() bool      { return true }

func (p *KeepAliveClientbound) Parse(data []byte) error {
	r := protocol.NewReader(bytes.NewReader(data))
	var err error

	p.KeepAliveID, err = r.ReadVarLong()
	return err
}

// ==================== SERVERBOUND PACKETS ====================

// PluginMessageServerbound sends a plugin message to the server.
type PluginMessageServerbound struct {
	Channel string
	Data    []byte
}

func (p *PluginMessageServerbound) ID() int32          { return PluginMessageServerboundPacketID }
func (p *PluginMessageServerbound) IsClientbound() bool { return false }

func (p *PluginMessageServerbound) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteString(p.Channel)
	w.WriteBytes(p.Data)
	return w.Bytes()
}

// ClientInformation sends client settings.
type ClientInformation struct {
	Locale             string
	ViewDistance       int8
	ChatMode           int32
	ChatColors         bool
	DisplayedSkinParts byte
	MainHand           int8
	EnableTextFiltering bool
	AllowServerListings bool
}

func (p *ClientInformation) ID() int32          { return ClientInformationPacketID }
func (p *ClientInformation) IsClientbound() bool { return false }

func (p *ClientInformation) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteString(p.Locale)
	w.WriteUint8(uint8(p.ViewDistance))
	w.WriteVarInt(p.ChatMode)
	w.WriteBool(p.ChatColors)
	w.WriteByte(p.DisplayedSkinParts)
	w.WriteUint8(uint8(p.MainHand))
	w.WriteBool(p.EnableTextFiltering)
	w.WriteBool(p.AllowServerListings)
	return w.Bytes()
}

// ClientCommand sends client commands.
type ClientCommand struct {
	ActionID int32
}

func (p *ClientCommand) ID() int32          { return ClientCommandPacketID }
func (p *ClientCommand) IsClientbound() bool { return false }

func (p *ClientCommand) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteVarInt(p.ActionID)
	return w.Bytes()
}

// PlayerChatMessage sends a chat message.
type PlayerChatMessage struct {
	Message string
}

func (p *PlayerChatMessage) ID() int32          { return PlayerChatMessagePacketID }
func (p *PlayerChatMessage) IsClientbound() bool { return false }

func (p *PlayerChatMessage) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteString(p.Message)
	return w.Bytes()
}

// PlayerPosition updates player position.
type PlayerPosition struct {
	X        float64
	Y        float64
	Z        float64
	OnGround bool
}

func (p *PlayerPosition) ID() int32          { return PlayerPositionPacketID }
func (p *PlayerPosition) IsClientbound() bool { return false }

func (p *PlayerPosition) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
	w.WriteFloat64(p.Z)
	w.WriteBool(p.OnGround)
	return w.Bytes()
}

// PlayerPositionAndLook updates player position and look.
type PlayerPositionAndLook struct {
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (p *PlayerPositionAndLook) ID() int32          { return PlayerPositionAndLookServerboundID }
func (p *PlayerPositionAndLook) IsClientbound() bool { return false }

func (p *PlayerPositionAndLook) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
	w.WriteFloat64(p.Z)
	w.WriteFloat32(p.Yaw)
	w.WriteFloat32(p.Pitch)
	w.WriteBool(p.OnGround)
	return w.Bytes()
}

// SetCreativeModeSlot sets an item in creative mode.
type SetCreativeModeSlot struct {
	Slot     int16
	ItemData []byte
}

func (p *SetCreativeModeSlot) ID() int32          { return SetCreativeModeSlotPacketID }
func (p *SetCreativeModeSlot) IsClientbound() bool { return false }

func (p *SetCreativeModeSlot) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteInt16(p.Slot)
	w.WriteBytes(p.ItemData)
	return w.Bytes()
}

// ClickContainer clicks on a container slot.
type ClickContainer struct {
	WindowID    int8
	Slot        int16
	Button      int8
	Mode        int32
	Slots       []SlotChange
	CarriedItem []byte
}

type SlotChange struct {
	Slot int16
	Item []byte
}

func (p *ClickContainer) ID() int32          { return ClickContainerPacketID }
func (p *ClickContainer) IsClientbound() bool { return false }

func (p *ClickContainer) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteByte(byte(p.WindowID))
	w.WriteInt16(p.Slot)
	w.WriteByte(byte(p.Button))
	w.WriteVarInt(p.Mode)
	w.WriteVarInt(int32(len(p.Slots)))

	for _, slot := range p.Slots {
		w.WriteInt16(slot.Slot)
		w.WriteBytes(slot.Item)
	}

	w.WriteBytes(p.CarriedItem)
	return w.Bytes()
}

// SetHeldItem changes the held item slot.
type SetHeldItem struct {
	Slot int8
}

func (p *SetHeldItem) ID() int32          { return SetHeldItemPacketID }
func (p *SetHeldItem) IsClientbound() bool { return false }

func (p *SetHeldItem) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteByte(byte(p.Slot))
	return w.Bytes()
}

// SetPlayerPositionAndRotation sets player position and rotation.
type SetPlayerPositionAndRotation struct {
	X        float64
	Y        float64
	Z        float64
	Yaw      float32
	Pitch    float32
	OnGround bool
}

func (p *SetPlayerPositionAndRotation) ID() int32          { return SetPlayerPositionAndRotationPacketID }
func (p *SetPlayerPositionAndRotation) IsClientbound() bool { return false }

func (p *SetPlayerPositionAndRotation) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteFloat64(p.X)
	w.WriteFloat64(p.Y)
	w.WriteFloat64(p.Z)
	w.WriteFloat32(p.Yaw)
	w.WriteFloat32(p.Pitch)
	w.WriteBool(p.OnGround)
	return w.Bytes()
}

// UpdateSelectedSlot updates the selected hotbar slot.
type UpdateSelectedSlot struct {
	Slot int8
}

func (p *UpdateSelectedSlot) ID() int32          { return UpdateSelectedSlotPacketID }
func (p *UpdateSelectedSlot) IsClientbound() bool { return false }

func (p *UpdateSelectedSlot) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteByte(byte(p.Slot))
	return w.Bytes()
}

// CloseContainerServerbound closes a container.
type CloseContainerServerbound struct {
	WindowID int8
}

func (p *CloseContainerServerbound) ID() int32          { return CloseContainerServerboundPacketID }
func (p *CloseContainerServerbound) IsClientbound() bool { return false }

func (p *CloseContainerServerbound) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteByte(byte(p.WindowID))
	return w.Bytes()
}

// KeepAliveServerbound is a keepalive response to the server.
type KeepAliveServerbound struct {
	KeepAliveID int64
}

func (p *KeepAliveServerbound) ID() int32          { return KeepAliveServerboundPacketID }
func (p *KeepAliveServerbound) IsClientbound() bool { return false }

func (p *KeepAliveServerbound) Serialize() []byte {
	w := protocol.NewWriter()
	w.WriteVarLong(p.KeepAliveID)
	return w.Bytes()
}

// ==================== PACKET REGISTRY ====================

// PacketRegistry maps packet IDs to packet types.
type PacketRegistry struct {
	mu                sync.RWMutex
	clientboundPackets map[int32]func() Packet
	serverboundPackets map[int32]func() Packet
}

// NewPacketRegistry creates a new packet registry.
func NewPacketRegistry() *PacketRegistry {
	registry := &PacketRegistry{
		clientboundPackets: make(map[int32]func() Packet),
		serverboundPackets: make(map[int32]func() Packet),
	}

	// Register clientbound packets
	registry.RegisterClientbound(JoinGamePacketID, func() Packet { return &JoinGame{} })
	registry.RegisterClientbound(PluginMessageClientboundPacketID, func() Packet { return &PluginMessageClientbound{} })
	registry.RegisterClientbound(ServerDataPacketID, func() Packet { return &ServerData{} })
	registry.RegisterClientbound(SynchronizePlayerPositionPacketID, func() Packet { return &SynchronizePlayerPosition{} })
	registry.RegisterClientbound(EntitySpawnPacketID, func() Packet { return &EntitySpawn{} })
	registry.RegisterClientbound(EntitySpawnExperienceOrbPacketID, func() Packet { return &EntitySpawnExperienceOrb{} })
	registry.RegisterClientbound(EntitySpawnMobPacketID, func() Packet { return &EntitySpawnMob{} })
	registry.RegisterClientbound(EntitySpawnPaintingPacketID, func() Packet { return &EntitySpawnPainting{} })
	registry.RegisterClientbound(EntitySpawnPlayerPacketID, func() Packet { return &EntitySpawnPlayer{} })
	registry.RegisterClientbound(EntityAnimationClientboundPacketID, func() Packet { return &EntityAnimationClientbound{} })
	registry.RegisterClientbound(EntityStatisticsPacketID, func() Packet { return &EntityStatistics{} })
	registry.RegisterClientbound(EntityPositionPacketID, func() Packet { return &EntityPosition{} })
	registry.RegisterClientbound(EntityPositionAndRotationPacketID, func() Packet { return &EntityPositionAndRotation{} })
	registry.RegisterClientbound(EntityVelocityPacketID, func() Packet { return &EntityVelocity{} })
	registry.RegisterClientbound(EntityEquipmentPacketID, func() Packet { return &EntityEquipment{} })
	registry.RegisterClientbound(EntityUpdateAttributesPacketID, func() Packet { return &EntityUpdateAttributes{} })
	registry.RegisterClientbound(EntityEffectPacketID, func() Packet { return &EntityEffect{} })
	registry.RegisterClientbound(EntityMetadataPacketID, func() Packet { return &EntityMetadata{} })
	registry.RegisterClientbound(EntityTeleportPacketID, func() Packet { return &EntityTeleport{} })
	registry.RegisterClientbound(EntityStatusPacketID, func() Packet { return &EntityStatus{} })
	registry.RegisterClientbound(EntityDamagePacketID, func() Packet { return &EntityDamage{} })
	registry.RegisterClientbound(EntityDeathPacketID, func() Packet { return &EntityDeath{} })
	registry.RegisterClientbound(SetExperiencePacketID, func() Packet { return &SetExperience{} })
	registry.RegisterClientbound(UpdateHealthPacketID, func() Packet { return &UpdateHealth{} })
	registry.RegisterClientbound(SetActionBarTextPacketID, func() Packet { return &SetActionBarText{} })
	registry.RegisterClientbound(SetTitleTextPacketID, func() Packet { return &SetTitleText{} })
	registry.RegisterClientbound(SetTimePacketID, func() Packet { return &SetTime{} })
	registry.RegisterClientbound(SetSlotPacketID, func() Packet { return &SetSlot{} })
	registry.RegisterClientbound(SetItemsPacketID, func() Packet { return &SetItems{} })
	registry.RegisterClientbound(OpenScreenPacketID, func() Packet { return &OpenScreen{} })
	registry.RegisterClientbound(CloseScreenClientboundPacketID, func() Packet { return &CloseScreenClientbound{} })
	registry.RegisterClientbound(ContainerSetContentPacketID, func() Packet { return &ContainerSetContent{} })
	registry.RegisterClientbound(ContainerSetDataPacketID, func() Packet { return &ContainerSetData{} })
	registry.RegisterClientbound(ContainerSetSlotPacketID, func() Packet { return &ContainerSetSlot{} })
	registry.RegisterClientbound(ContainerClosePacketID, func() Packet { return &ContainerClose{} })
	registry.RegisterClientbound(BlockUpdatePacketID, func() Packet { return &BlockUpdate{} })
	registry.RegisterClientbound(ChunkDataPacketID, func() Packet { return &ChunkData{} })
	registry.RegisterClientbound(UnloadChunkPacketID, func() Packet { return &UnloadChunk{} })
	registry.RegisterClientbound(BlockChangedAckPacketID, func() Packet { return &BlockChangedAck{} })
	registry.RegisterClientbound(SectionBlocksUpdatePacketID, func() Packet { return &SectionBlocksUpdate{} })
	registry.RegisterClientbound(GameEventPacketID, func() Packet { return &GameEvent{} })
	registry.RegisterClientbound(LevelChunkPacketID, func() Packet { return &LevelChunk{} })
	registry.RegisterClientbound(LevelChunkWithLightPacketID, func() Packet { return &LevelChunkWithLight{} })
	registry.RegisterClientbound(BlockEntityDataPacketID, func() Packet { return &BlockEntityData{} })
	registry.RegisterClientbound(BlockEventPacketID, func() Packet { return &BlockEvent{} })
	registry.RegisterClientbound(SoundPacketID, func() Packet { return &Sound{} })
	registry.RegisterClientbound(ParticlePacketID, func() Packet { return &Particle{} })
	registry.RegisterClientbound(ExplosionPacketID, func() Packet { return &Explosion{} })
	registry.RegisterClientbound(DisconnectClientboundPacketID, func() Packet { return &DisconnectClientbound{} })
	registry.RegisterClientbound(ServerPlayerPacketID, func() Packet { return &ServerPlayer{} })
	registry.RegisterClientbound(KeepAliveClientboundPacketID, func() Packet { return &KeepAliveClientbound{} })

	// Register serverbound packets
	registry.RegisterServerbound(PluginMessageServerboundPacketID, func() Packet { return &PluginMessageServerbound{} })
	registry.RegisterServerbound(ClientInformationPacketID, func() Packet { return &ClientInformation{} })
	registry.RegisterServerbound(ClientCommandPacketID, func() Packet { return &ClientCommand{} })
	registry.RegisterServerbound(PlayerChatMessagePacketID, func() Packet { return &PlayerChatMessage{} })
	registry.RegisterServerbound(PlayerPositionPacketID, func() Packet { return &PlayerPosition{} })
	registry.RegisterServerbound(PlayerPositionAndLookServerboundID, func() Packet { return &PlayerPositionAndLook{} })
	registry.RegisterServerbound(SetCreativeModeSlotPacketID, func() Packet { return &SetCreativeModeSlot{} })
	registry.RegisterServerbound(ClickContainerPacketID, func() Packet { return &ClickContainer{} })
	registry.RegisterServerbound(SetHeldItemPacketID, func() Packet { return &SetHeldItem{} })
	registry.RegisterServerbound(SetPlayerPositionAndRotationPacketID, func() Packet { return &SetPlayerPositionAndRotation{} })
	registry.RegisterServerbound(UpdateSelectedSlotPacketID, func() Packet { return &UpdateSelectedSlot{} })
	registry.RegisterServerbound(CloseContainerServerboundPacketID, func() Packet { return &CloseContainerServerbound{} })
	registry.RegisterServerbound(KeepAliveServerboundPacketID, func() Packet { return &KeepAliveServerbound{} })

	return registry
}

// RegisterClientbound registers a clientbound packet.
func (r *PacketRegistry) RegisterClientbound(id int32, factory func() Packet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clientboundPackets[id] = factory
}

// RegisterServerbound registers a serverbound packet.
func (r *PacketRegistry) RegisterServerbound(id int32, factory func() Packet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serverboundPackets[id] = factory
}

// CreateClientboundPacket creates a clientbound packet from ID.
func (r *PacketRegistry) CreateClientboundPacket(id int32) (Packet, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.clientboundPackets[id]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// CreateServerboundPacket creates a serverbound packet from ID.
func (r *PacketRegistry) CreateServerboundPacket(id int32) (Packet, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.serverboundPackets[id]
	if !ok {
		return nil, false
	}
	return factory(), true
}

// ParseClientboundPacket parses a clientbound packet.
func (r *PacketRegistry) ParseClientboundPacket(id int32, data []byte) (Packet, error) {
	packet, ok := r.CreateClientboundPacket(id)
	if !ok {
		return nil, fmt.Errorf("unknown clientbound packet ID: %d", id)
	}

	if parser, ok := packet.(interface{ Parse([]byte) error }); ok {
		if err := parser.Parse(data); err != nil {
			return nil, fmt.Errorf("failed to parse packet %d: %w", id, err)
		}
	}

	return packet, nil
}

// SerializeServerboundPacket serializes a serverbound packet.
func (r *PacketRegistry) SerializeServerboundPacket(packet Packet) ([]byte, error) {
	if serializer, ok := packet.(interface{ Serialize() []byte }); ok {
		return serializer.Serialize(), nil
	}
	return nil, fmt.Errorf("packet %d does not support serialization", packet.ID())
}

// DefaultRegistry is the default packet registry.
var DefaultRegistry = NewPacketRegistry()
