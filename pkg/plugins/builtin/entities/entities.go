// Package entities implements entity tracking.
package entities

import (
	"sync"

	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
)

// Plugin tracks entities in the world.
type Plugin struct {
	*plugins.BasePlugin

	// Entity storage
	entities  map[int64]*Entity
	entitiesMu sync.RWMutex

	// Player entity ID
	playerID  int64
}

// Entity represents an entity in the world.
type Entity struct {
	ID       int64
	Type     int32
	UUID     string
	Position Position
	Motion   Motion
	Rotation Rotation
	OnGround bool
}

// Position represents a 3D position.
type Position struct {
	X, Y, Z float64
}

// Motion represents velocity.
type Motion struct {
	X, Y, Z float64
}

// Rotation represents rotation angles.
type Rotation struct {
	Yaw, Pitch float32
}

// NewPlugin creates a new entities plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("entities", "1.0.0")

	return &Plugin{
		BasePlugin: base,
		entities:   make(map[int64]*Entity),
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Subscribe to entity packets
	p.On("packet", p.handlePacket)

	return nil
}

// handlePacket handles incoming packets.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet := data[0].(*protocol.Packet)

	switch packet.ID {
	case 0x01: // Spawn Entity packet
		p.handleSpawnEntity(packet)
	case 0x02: // Spawn Experience Orb
		// Ignore for now
	case 0x03: // Spawn Player packet
		p.handleSpawnPlayer(packet)
	case 0x04: // Entity Animation
		// Ignore for now
	case 0x26: // Entity Position packet
		p.handleEntityPosition(packet)
	case 0x27: // Entity Position and Rotation packet
		p.handleEntityPositionRotation(packet)
	case 0x28: // Entity Rotation packet
		p.handleEntityRotation(packet)
	case 0x29: // Entity Movement packet
		p.handleEntityMovement(packet)
	case 0x2B: // Remove Entities packet
		p.handleRemoveEntities(packet)
	}
}

// handleSpawnEntity handles the Spawn Entity packet.
func (p *Plugin) handleSpawnEntity(packet *protocol.Packet) {
	// TODO: Parse entity data from packet
	// For now, emit event
	p.Emit("entity_spawned")
}

// handleSpawnPlayer handles the Spawn Player packet.
func (p *Plugin) handleSpawnPlayer(packet *protocol.Packet) {
	// TODO: Parse player data from packet
	p.Emit("player_spawned")
}

// handleEntityPosition handles the Entity Position packet.
func (p *Plugin) handleEntityPosition(packet *protocol.Packet) {
	// TODO: Parse position update
	p.Emit("entity_moved")
}

// handleEntityPositionRotation handles the Entity Position and Rotation packet.
func (p *Plugin) handleEntityPositionRotation(packet *protocol.Packet) {
	// TODO: Parse position and rotation update
	p.Emit("entity_moved")
	p.Emit("entity_rotated")
}

// handleEntityRotation handles the Entity Rotation packet.
func (p *Plugin) handleEntityRotation(packet *protocol.Packet) {
	// TODO: Parse rotation update
	p.Emit("entity_rotated")
}

// handleEntityMovement handles the Entity Movement packet (relative move).
func (p *Plugin) handleEntityMovement(packet *protocol.Packet) {
	// TODO: Parse movement delta
	p.Emit("entity_moved")
}

// handleRemoveEntities handles the Remove Entities packet.
func (p *Plugin) handleRemoveEntities(packet *protocol.Packet) {
	// TODO: Parse entity IDs and remove them
	p.Emit("entity_removed")
}

// Get returns an entity by ID.
func (p *Plugin) Get(id int64) (*Entity, bool) {
	p.entitiesMu.RLock()
	defer p.entitiesMu.RUnlock()

	entity, exists := p.entities[id]
	return entity, exists
}

// All returns all tracked entities.
func (p *Plugin) All() []*Entity {
	p.entitiesMu.RLock()
	defer p.entitiesMu.RUnlock()

	entities := make([]*Entity, 0, len(p.entities))
	for _, entity := range p.entities {
		entities = append(entities, entity)
	}
	return entities
}

// Count returns the number of tracked entities.
func (p *Plugin) Count() int {
	p.entitiesMu.RLock()
	defer p.entitiesMu.RUnlock()
	return len(p.entities)
}

// SetPlayerID sets the player entity ID.
func (p *Plugin) SetPlayerID(id int64) {
	p.entitiesMu.Lock()
	defer p.entitiesMu.Unlock()
	p.playerID = id
}

// PlayerID returns the player entity ID.
func (p *Plugin) PlayerID() int64 {
	p.entitiesMu.RLock()
	defer p.entitiesMu.RUnlock()
	return p.playerID
}

// Add adds an entity to the tracker.
func (p *Plugin) Add(entity *Entity) {
	p.entitiesMu.Lock()
	defer p.entitiesMu.Unlock()

	p.entities[entity.ID] = entity
	p.Emit("entity_added", entity)
}

// Remove removes an entity from the tracker.
func (p *Plugin) Remove(id int64) {
	p.entitiesMu.Lock()
	defer p.entitiesMu.Unlock()

	if entity, exists := p.entities[id]; exists {
		delete(p.entities, id)
		p.Emit("entity_removed", entity)
	}
}

// Clear removes all entities.
func (p *Plugin) Clear() {
	p.entitiesMu.Lock()
	defer p.entitiesMu.Unlock()

	p.entities = make(map[int64]*Entity)
}
