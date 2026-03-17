package plugins

import (
	"github.com/go-flayer/goflayer/pkg/goflayer"
	"github.com/go-flayer/goflayer/pkg/math"
)

// Plugin is the base interface that all plugins must implement.
//
// Plugins extend the bot's functionality by injecting methods, events,
// and behaviors into the bot instance.
//
// Example plugin:
//
//	type MyPlugin struct {
//	    bot goflayer.Bot
//	}
//
//	func (p *MyPlugin) Name() string {
//	    return "myPlugin"
//	}
//
//	func (p *MyPlugin) Version() string {
//	    return "1.0.0"
//	}
//
//	func (p *MyPlugin) Inject(bot goflayer.Bot, options goflayer.Options) error {
//	    p.bot = bot
//	    // Register handlers, add methods, etc.
//	    return nil
//	}
//
//	func (p *MyPlugin) Cleanup() error {
//	    // Cleanup resources
//	    return nil
//	}
type Plugin interface {
	// Name returns the unique name of the plugin.
	// This must be unique across all plugins.
	Name() string

	// Version returns the version string of the plugin.
	Version() string

	// Inject is called when the plugin is loaded into a bot.
	// This is where the plugin should:
	// - Register event handlers
	// - Add methods to the bot
	// - Initialize plugin state
	// - Start background goroutines if needed
	Inject(bot goflayer.Bot, options goflayer.Options) error

	// Cleanup is called when the bot is shutting down or the plugin is unloaded.
	// This is where the plugin should:
	// - Unregister event handlers
	// - Stop background goroutines
	// - Release resources
	// - Save any necessary state
	Cleanup() error
}

// EntityPlugin handles entity-related functionality.
//
// Plugins that work with entities (players, mobs, items) should implement this interface.
type EntityPlugin interface {
	Plugin

	// OnEntitySpawn is called when a new entity spawns.
	OnEntitySpawn(entity *Entity) error

	// OnEntityDespawn is called when an entity is removed from the world.
	OnEntityDespawn(entity *Entity) error

	// OnEntityMove is called when an entity moves.
	OnEntityMove(entity *Entity) error

	// OnEntityRotate is called when an entity rotates.
	OnEntityRotate(entity *Entity) error
}

// PhysicsPlugin handles physics simulation.
//
// Plugins that simulate physics (movement, collision, gravity) should implement this interface.
type PhysicsPlugin interface {
	Plugin

	// Update updates the physics simulation.
	Update(delta float64) error

	// SetControlState sets a control state (forward, back, jump, etc.).
	SetControlState(control Control, state bool) error

	// GetControlState gets the current state of a control.
	GetControlState(control Control) bool

	// Position returns the current position.
	Position() *math.Vec3

	// Velocity returns the current velocity.
	Velocity() *math.Vec3

	// OnGround returns true if on the ground.
	OnGround() bool

	// InWater returns true if in water.
	InWater() bool
}

// Control represents a movement control state.
type Control int

const (
	// Forward - move forward control
	Forward Control = iota

	// Back - move backward control
	Back

	// Left - strafe left control
	Left

	// Right - strafe right control
	Right

	// Jump - jump control
	Jump

	// Sprint - sprint control
	Sprint

	// Sneak - sneak control
	Sneak
)

// String returns the string representation of the control.
func (c Control) String() string {
	switch c {
	case Forward:
		return "Forward"
	case Back:
		return "Back"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case Jump:
		return "Jump"
	case Sprint:
		return "Sprint"
	case Sneak:
		return "Sneak"
	default:
		return "Unknown"
	}
}

// ChatPlugin handles chat messages.
//
// Plugins that work with chat should implement this interface.
type ChatPlugin interface {
	Plugin

	// OnMessage is called when a chat message is received.
	OnMessage(username, message string, position ChatPosition) error

	// SendMessage sends a chat message.
	SendMessage(message string) error

	// AddPattern adds a chat message pattern handler.
	AddPattern(pattern string, handler func(username, message string)) error

	// RemovePattern removes a chat message pattern handler.
	RemovePattern(pattern string) error
}

// ChatPosition represents where a chat message came from.
type ChatPosition int

const (
	// ChatPositionChatBox - regular chat message
	ChatPositionChatBox ChatPosition = iota

	// ChatPositionSystemInfo - system message
	ChatPositionSystemInfo

	// ChatPositionGameInfo - game info (e.g., "You died")
	ChatPositionGameInfo

	// ChatPositionAboveActionBar - above action bar
	ChatPositionAboveActionBar
)

// String returns the string representation of the chat position.
func (p ChatPosition) String() string {
	switch p {
	case ChatPositionChatBox:
		return "ChatBox"
	case ChatPositionSystemInfo:
		return "SystemInfo"
	case ChatPositionGameInfo:
		return "GameInfo"
	case ChatPositionAboveActionBar:
		return "AboveActionBar"
	default:
		return "Unknown"
	}
}

// BlockPlugin handles block-related functionality.
//
// Plugins that work with blocks (digging, placing) should implement this interface.
type BlockPlugin interface {
	Plugin

	// OnBlockUpdate is called when a block changes.
	OnBlockUpdate(position *math.Vec3, oldBlock, newBlock *Block) error

	// DigBlock starts digging a block at the position.
	DigBlock(position *math.Vec3) error

	// PlaceBlock places a block at the position.
	PlaceBlock(position *math.Vec3, block *Block) error

	// GetBlock returns the block at the position.
	GetBlock(position *math.Vec3) (*Block, error)
}

// Block represents a block in the world.
type Block struct {
	// Type is the block type ID
	Type int

	// Position is the block position
	Position *math.Vec3

	// Metadata contains block-specific metadata
	Metadata map[string]interface{}
}

// InventoryPlugin handles inventory management.
//
// Plugins that work with inventories should implement this interface.
type InventoryPlugin interface {
	Plugin

	// OnWindowOpen is called when a window is opened.
	OnWindowOpen(window Window) error

	// OnWindowClose is called when a window is closed.
	OnWindowClose(windowID int) error

	// OnSlotUpdate is called when a slot changes.
	OnSlotUpdate(windowID, slot int, item *Item) error

	// GetItems returns all items in a window.
	GetItems(windowID int) ([]*Item, error)

	// ClickSlot clicks a slot in a window.
	ClickSlot(windowID, slot, mouseButton, mode int) error

	// CloseWindow closes a window.
	CloseWindow(windowID int) error
}

// Window represents an inventory window.
type Window interface {
	// ID returns the window ID
	ID() int

	// Type returns the window type
	Type() string

	// Title returns the window title
	Title() string

	// Slots returns all slots in the window
	Slots() []*Item

	// Slot returns the item in a specific slot
	Slot(slot int) (*Item, error)
}

// Item represents an item stack.
type Item struct {
	// Type is the item type ID
	Type int

	// Count is the number of items in the stack
	Count int

	// DisplayName is the item's display name
	DisplayName string

	// NBT contains the item's NBT data
	NBT map[string]interface{}
}

// Entity represents a generic entity.
type Entity struct {
	// ID is the entity ID
	ID int32

	// Type is the entity type
	Type string

	// Position is the entity position
	Position *math.Vec3

	// Velocity is the entity velocity
	Velocity *math.Vec3

	// Rotation is the entity rotation (yaw, pitch)
	Rotation *math.Vec3

	// Metadata contains entity-specific metadata
	Metadata map[string]interface{}
}

// Player represents a player entity.
type Player struct {
	Entity

	// Username is the player's username
	Username string

	// UUID is the player's UUID
	UUID string

	// GameMode is the player's game mode
	GameMode int
}

// GameMode values
const (
	GameModeSurvival = 0
	GameModeCreative = 1
	GameModeAdventure = 2
	GameModeSpectator = 3
)

// WorldPlugin handles world-related functionality.
//
// Plugins that work with the world (chunks, biomes) should implement this interface.
type WorldPlugin interface {
	Plugin

	// OnChunkLoad is called when a chunk is loaded.
	OnChunkLoad(x, z int) error

	// OnChunkUnload is called when a chunk is unloaded.
	OnChunkUnload(x, z int) error

	// GetBlockAt returns the block at a world position.
	GetBlockAt(x, y, z int) (*Block, error)

	// GetBiomeAt returns the biome at a position.
	GetBiomeAt(x, z int) (string, error)

	// IsChunkLoaded checks if a chunk is loaded.
	IsChunkLoaded(x, z int) bool
}

// Entity reference for avoiding circular import
type Entity = Entity
type Player = Player
