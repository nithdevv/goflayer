package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// Entity represents any entity in the game world.
type Entity struct {
	ID       int32
	Type     string
	Position *math.Vec3
	Rotation *math.Vec3 // yaw, pitch, roll
	Velocity *math.Vec3
	OnGround bool
	Metadata map[string]interface{}
}

// Player represents a player entity.
type Player struct {
	Entity
	Username string
	UUID     string
	Gamemode int
	Ping     int
}

// Mob represents a mob entity.
type Mob struct {
	Entity
	Name      string
	Health    float32
	MaxHealth float32
}

// Object represents an item/object entity (dropped items, arrows, etc.).
type Object struct {
	Entity
	ItemType string
	Count    int
}

// EntitiesPlugin tracks all entities in the world.
type EntitiesPlugin struct {
	mu       sync.RWMutex
	ctx      *plugins.Context
	log      *logger.Logger
	entities map[int32]*Entity
	players  map[string]*Player // keyed by UUID
	mobs     map[int32]*Mob
	objects  map[int32]*Object
	nextID   int32
}

// Metadata returns the plugin metadata.
func (p *EntitiesPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "entities",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Entity tracking and queries",
		Dependencies: []string{},
	}
}

// OnLoad initializes the entities plugin.
func (p *EntitiesPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.entities = make(map[int32]*Entity)
	p.players = make(map[string]*Player)
	p.mobs = make(map[int32]*Mob)
	p.objects = make(map[int32]*Object)
	p.nextID = 0

	p.log.Info("Entities plugin loaded")

	// Register event handlers
	p.ctx.On("spawn_entity", p.handleSpawnEntity)
	p.ctx.On("spawn_player", p.handleSpawnPlayer)
	p.ctx.On("spawn_mob", p.handleSpawnMob)
	p.ctx.On("spawn_object", p.handleSpawnObject)
	p.ctx.On("entity_destroy", p.handleEntityDestroy)
	p.ctx.On("entity_position", p.handleEntityPosition)
	p.ctx.On("entity_metadata", p.handleEntityMetadata)

	return nil
}

// OnUnload cleans up the entities plugin.
func (p *EntitiesPlugin) OnUnload() error {
	p.log.Info("Entities plugin unloaded")
	return nil
}

// GetAll returns all entities.
func (p *EntitiesPlugin) GetAll() []*Entity {
	p.mu.RLock()
	defer p.mu.RUnlock()

	entities := make([]*Entity, 0, len(p.entities))
	for _, entity := range p.entities {
		entities = append(entities, entity)
	}
	return entities
}

// GetByID returns an entity by ID.
func (p *EntitiesPlugin) GetByID(id int32) (*Entity, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	entity, exists := p.entities[id]
	return entity, exists
}

// GetByType returns all entities of a specific type.
func (p *EntitiesPlugin) GetByType(entityType string) []*Entity {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]*Entity, 0)
	for _, entity := range p.entities {
		if entity.Type == entityType {
			result = append(result, entity)
		}
	}
	return result
}

// GetNearby returns entities within a radius of a position.
func (p *EntitiesPlugin) GetNearby(pos *math.Vec3, radius float64) []*Entity {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]*Entity, 0)
	radiusSq := radius * radius

	for _, entity := range p.entities {
		if entity.Position.DistanceSquaredTo(pos) <= radiusSq {
			result = append(result, entity)
		}
	}
	return result
}

// GetPlayers returns all players.
func (p *EntitiesPlugin) GetPlayers() []*Player {
	p.mu.RLock()
	defer p.mu.RUnlock()

	players := make([]*Player, 0, len(p.players))
	for _, player := range p.players {
		players = append(players, player)
	}
	return players
}

// GetPlayerByName returns a player by username.
func (p *EntitiesPlugin) GetPlayerByName(username string) (*Player, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, player := range p.players {
		if player.Username == username {
			return player, true
		}
	}
	return nil, false
}

// GetPlayerByUUID returns a player by UUID.
func (p *EntitiesPlugin) GetPlayerByUUID(uuid string) (*Player, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	player, exists := p.players[uuid]
	return player, exists
}

// GetMobs returns all mobs.
func (p *EntitiesPlugin) GetMobs() []*Mob {
	p.mu.RLock()
	defer p.mu.RUnlock()

	mobs := make([]*Mob, 0, len(p.mobs))
	for _, mob := range p.mobs {
		mobs = append(mobs, mob)
	}
	return mobs
}

// GetMobByName returns all mobs of a specific name.
func (p *EntitiesPlugin) GetMobByName(name string) []*Mob {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]*Mob, 0)
	for _, mob := range p.mobs {
		if mob.Name == name {
			result = append(result, mob)
		}
	}
	return result
}

// GetNearestMob returns the nearest mob of a specific type.
func (p *EntitiesPlugin) GetNearestMob(pos *math.Vec3, mobName string) *Mob {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var nearest *Mob
	nearestDist := float64(-1)

	for _, mob := range p.mobs {
		if mobName != "" && mob.Name != mobName {
			continue
		}

		dist := mob.Position.DistanceTo(pos)
		if nearestDist < 0 || dist < nearestDist {
			nearest = mob
			nearestDist = dist
		}
	}

	return nearest
}

// Count returns the total number of entities.
func (p *EntitiesPlugin) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entities)
}

// CountPlayers returns the number of players.
func (p *EntitiesPlugin) CountPlayers() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.players)
}

// CountMobs returns the number of mobs.
func (p *EntitiesPlugin) CountMobs() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.mobs)
}

// Event handlers

func (p *EntitiesPlugin) handleSpawnEntity(args ...interface{}) {
	// TODO: Parse entity spawn packet
	p.log.Debug("Entity spawned")
}

func (p *EntitiesPlugin) handleSpawnPlayer(args ...interface{}) {
	// TODO: Parse player spawn packet
	p.log.Debug("Player spawned")
}

func (p *EntitiesPlugin) handleSpawnMob(args ...interface{}) {
	// TODO: Parse mob spawn packet
	p.log.Debug("Mob spawned")
}

func (p *EntitiesPlugin) handleSpawnObject(args ...interface{}) {
	// TODO: Parse object spawn packet
	p.log.Debug("Object spawned")
}

func (p *EntitiesPlugin) handleEntityDestroy(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	if entityIDs, ok := args[0].([]int32); ok {
		p.mu.Lock()
		for _, id := range entityIDs {
			delete(p.entities, id)
			delete(p.mobs, id)
			delete(p.objects, id)
		}
		p.mu.Unlock()

		p.log.Debug("Destroyed %d entities", len(entityIDs))
		p.ctx.Emit("entities_removed", entityIDs)
	}
}

func (p *EntitiesPlugin) handleEntityPosition(args ...interface{}) {
	// TODO: Update entity position
	p.log.Debug("Entity position updated")
}

func (p *EntitiesPlugin) handleEntityMetadata(args ...interface{}) {
	// TODO: Update entity metadata
	p.log.Debug("Entity metadata updated")
}

// String returns a string representation of the entities plugin.
func (p *EntitiesPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Entities{total=%d, players=%d, mobs=%d, objects=%d}",
		len(p.entities), len(p.players), len(p.mobs), len(p.objects))
}
