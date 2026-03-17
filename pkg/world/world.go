// Package world implements Minecraft world tracking and management.
package world

import (
	"fmt"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/pkg/events"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/registry"
)

// World represents a Minecraft world with chunks, entities, and players.
type World struct {
	mu       sync.RWMutex
	name     string
	chunks   map[ChunkPos]*Chunk
	entities map[int32]Entity
	players  map[uuid]*Player
	events   *events.Bus

	// World properties
	dimension      Dimension
	difficulty     Difficulty
	worldType      WorldType
	seed           int64
	spawnPosition  *math.BlockPos
	time           int64
	age            int64
	raining        bool
	thundering     bool

	// Tracking
	lastUpdate     time.Time
	totalChunks    int
	totalEntities  int
	totalPlayers   int
}

// Dimension represents a world dimension.
type Dimension int32

const (
	DimensionNether   Dimension = -1
	DimensionOverworld Dimension = 0
	DimensionEnd      Dimension = 1
)

func (d Dimension) String() string {
	switch d {
	case DimensionNether:
		return "nether"
	case DimensionOverworld:
		return "overworld"
	case DimensionEnd:
		return "end"
	default:
		return fmt.Sprintf("unknown(%d)", d)
	}
}

// Difficulty represents world difficulty.
type Difficulty byte

const (
	DifficultyPeaceful Difficulty = 0
	DifficultyEasy     Difficulty = 1
	DifficultyNormal   Difficulty = 2
	DifficultyHard     Difficulty = 3
)

func (d Difficulty) String() string {
	switch d {
	case DifficultyPeaceful:
		return "peaceful"
	case DifficultyEasy:
		return "easy"
	case DifficultyNormal:
		return "normal"
	case DifficultyHard:
		return "hard"
	default:
		return fmt.Sprintf("unknown(%d)", d)
	}
}

// WorldType represents the world type.
type WorldType string

const (
	WorldTypeDefault   WorldType = "default"
	WorldTypeFlat      WorldType = "flat"
	WorldTypeLargeBiomes WorldType = "largeBiomes"
	WorldTypeAmplified WorldType = "amplified"
	WorldTypeDebug     WorldType = "debug_all_block_states"
)

// NewWorld creates a new world instance.
func NewWorld(name string) *World {
	return &World{
		name:     name,
		chunks:   make(map[ChunkPos]*Chunk),
		entities: make(map[int32]Entity),
		players:  make(map[uuid]*Player),
		events:   events.NewBus(),
		dimension: DimensionOverworld,
		difficulty: DifficultyNormal,
		worldType: WorldTypeDefault,
		lastUpdate: time.Now(),
	}
}

// Name returns the world name.
func (w *World) Name() string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.name
}

// SetName sets the world name.
func (w *World) SetName(name string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.name = name
}

// Events returns the event bus.
func (w *World) Events() *events.Bus {
	return w.events
}

// Dimension returns the current dimension.
func (w *World) Dimension() Dimension {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.dimension
}

// SetDimension sets the dimension.
func (w *World) SetDimension(dim Dimension) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.dimension = dim
}

// Difficulty returns the world difficulty.
func (w *World) Difficulty() Difficulty {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.difficulty
}

// SetDifficulty sets the difficulty.
func (w *World) SetDifficulty(diff Difficulty) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.difficulty = diff
}

// WorldType returns the world type.
func (w *World) WorldType() WorldType {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.worldType
}

// SetWorldType sets the world type.
func (w *World) SetWorldType(worldType WorldType) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.worldType = worldType
}

// Seed returns the world seed.
func (w *World) Seed() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.seed
}

// SetSeed sets the world seed.
func (w *World) SetSeed(seed int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.seed = seed
}

// SpawnPosition returns the spawn position.
func (w *World) SpawnPosition() *math.BlockPos {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.spawnPosition == nil {
		return nil
	}
	return &math.BlockPos{
		X: w.spawnPosition.X,
		Y: w.spawnPosition.Y,
		Z: w.spawnPosition.Z,
	}
}

// SetSpawnPosition sets the spawn position.
func (w *World) SetSpawnPosition(pos *math.BlockPos) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.spawnPosition = pos
}

// Time returns the world time.
func (w *World) Time() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.time
}

// SetTime sets the world time.
func (w *World) SetTime(time int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.time = time
}

// Age returns the world age.
func (w *World) Age() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.age
}

// SetAge sets the world age.
func (w *World) SetAge(age int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.age = age
}

// IsRaining returns whether it's raining.
func (w *World) IsRaining() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.raining
}

// SetRaining sets whether it's raining.
func (w *World) SetRaining(raining bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.raining = raining
}

// IsThundering returns whether it's thundering.
func (w *World) IsThundering() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.thundering
}

// SetThundering sets whether it's thundering.
func (w *World) SetThundering(thundering bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.thundering = thundering
}

// GetChunk retrieves a chunk at the given position.
func (w *World) GetChunk(x, z int32) (*Chunk, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	pos := ChunkPos{X: x, Z: z}
	chunk, ok := w.chunks[pos]
	return chunk, ok
}

// SetChunk sets a chunk at the given position.
func (w *World) SetChunk(x, z int32, chunk *Chunk) {
	w.mu.Lock()
	defer w.mu.Unlock()

	pos := ChunkPos{X: x, Z: z}

	// Check if chunk already exists
	_, exists := w.chunks[pos]
	w.chunks[pos] = chunk

	if !exists {
		w.totalChunks++
		w.lastUpdate = time.Now()

		// Emit chunk_load event
		w.events.Emit("chunk_load", x, z, chunk)
	}
}

// RemoveChunk removes a chunk at the given position.
func (w *World) RemoveChunk(x, z int32) {
	w.mu.Lock()
	defer w.mu.Unlock()

	pos := ChunkPos{X: x, Z: z}
	if chunk, ok := w.chunks[pos]; ok {
		delete(w.chunks, pos)
		w.totalChunks--
		w.lastUpdate = time.Now()

		// Emit chunk_unload event
		w.events.Emit("chunk_unload", x, z, chunk)
	}
}

// GetBlock retrieves a block at the given world position.
func (w *World) GetBlock(x, y, z int32) (registry.BlockID, error) {
	chunkX := x >> 4
	chunkZ := z >> 4

	chunk, ok := w.GetChunk(chunkX, chunkZ)
	if !ok {
		return registry.BlockAir, fmt.Errorf("chunk not loaded at %d, %d", chunkX, chunkZ)
	}

	return chunk.GetBlock(x&15, y, z&15)
}

// SetBlock sets a block at the given world position.
func (w *World) SetBlock(x, y, z int32, blockID registry.BlockID) error {
	chunkX := x >> 4
	chunkZ := z >> 4

	chunk, ok := w.GetChunk(chunkX, chunkZ)
	if !ok {
		return fmt.Errorf("chunk not loaded at %d, %d", chunkX, chunkZ)
	}

	oldBlock, _ := chunk.GetBlock(x&15, y, z&15)
	err := chunk.SetBlock(x&15, y, z&15, blockID)
	if err != nil {
		return err
	}

	// Emit block_update event
	pos := &math.BlockPos{X: int(x), Y: int(y), Z: int(z)}
	w.events.Emit("block_update", pos, oldBlock, blockID)

	return nil
}

// GetBlockEntity retrieves a block entity at the given position.
func (w *World) GetBlockEntity(x, y, z int32) (interface{}, bool) {
	chunkX := x >> 4
	chunkZ := z >> 4

	chunk, ok := w.GetChunk(chunkX, chunkZ)
	if !ok {
		return nil, false
	}

	return chunk.GetBlockEntity(x&15, y, z&15)
}

// SetBlockEntity sets a block entity at the given position.
func (w *World) SetBlockEntity(x, y, z int32, entity interface{}) {
	chunkX := x >> 4
	chunkZ := z >> 4

	chunk, ok := w.GetChunk(chunkX, chunkZ)
	if !ok {
		return
	}

	chunk.SetBlockEntity(x&15, y, z&15, entity)
}

// AddEntity adds an entity to the world.
func (w *World) AddEntity(entity Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.entities[entity.ID()] = entity
	w.totalEntities++
	w.lastUpdate = time.Now()

	// Emit entity_spawn event
	w.events.Emit("entity_spawn", entity)
}

// GetEntity retrieves an entity by ID.
func (w *World) GetEntity(id int32) (Entity, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	entity, ok := w.entities[id]
	return entity, ok
}

// RemoveEntity removes an entity from the world.
func (w *World) RemoveEntity(id int32) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if entity, ok := w.entities[id]; ok {
		delete(w.entities, id)
		w.totalEntities--
		w.lastUpdate = time.Now()

		// Emit entity_despawn event
		w.events.Emit("entity_despawn", entity)
	}
}

// GetEntitiesInRange returns all entities within the given range.
func (w *World) GetEntitiesInRange(center *math.Vec3, distance float64) []Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var entities []Entity
	distanceSq := distance * distance

	for _, entity := range w.entities {
		pos := entity.Position()
		if center.DistanceSquaredTo(pos) <= distanceSq {
			entities = append(entities, entity)
		}
	}

	return entities
}

// AddPlayer adds a player to the world.
func (w *World) AddPlayer(player *Player) {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := player.UUID()
	w.players[id] = player
	w.totalPlayers++
	w.lastUpdate = time.Now()

	// Also add as entity
	w.entities[player.entityID] = player
}

// GetPlayer retrieves a player by UUID.
func (w *World) GetPlayer(id uuid) (*Player, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	player, ok := w.players[id]
	return player, ok
}

// GetPlayerByName retrieves a player by name.
func (w *World) GetPlayerByName(name string) (*Player, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, player := range w.players {
		if player.Name() == name {
			return player, true
		}
	}

	return nil, false
}

// RemovePlayer removes a player from the world.
func (w *World) RemovePlayer(id uuid) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if player, ok := w.players[id]; ok {
		delete(w.players, id)
		delete(w.entities, player.entityID)
		w.totalPlayers--
		w.lastUpdate = time.Now()
	}
}

// GetPlayersInRange returns all players within the given range.
func (w *World) GetPlayersInRange(center *math.Vec3, distance float64) []*Player {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var players []*Player
	distanceSq := distance * distance

	for _, player := range w.players {
		pos := player.Position()
		if center.DistanceSquaredTo(pos) <= distanceSq {
			players = append(players, player)
		}
	}

	return players
}

// GetAllPlayers returns all players in the world.
func (w *World) GetAllPlayers() []*Player {
	w.mu.RLock()
	defer w.mu.RUnlock()

	players := make([]*Player, 0, len(w.players))
	for _, player := range w.players {
		players = append(players, player)
	}

	return players
}

// ChunkCount returns the number of loaded chunks.
func (w *World) ChunkCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.totalChunks
}

// EntityCount returns the number of entities in the world.
func (w *World) EntityCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.totalEntities
}

// PlayerCount returns the number of players in the world.
func (w *World) PlayerCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.totalPlayers
}

// LastUpdate returns the last update time.
func (w *World) LastUpdate() time.Time {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.lastUpdate
}

// Clear clears all world data.
func (w *World) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.chunks = make(map[ChunkPos]*Chunk)
	w.entities = make(map[int32]Entity)
	w.players = make(map[uuid]*Player)
	w.totalChunks = 0
	w.totalEntities = 0
	w.totalPlayers = 0
	w.lastUpdate = time.Now()
}

// Close closes the world and cleans up resources.
func (w *World) Close() error {
	w.Clear()
	return w.events.Close()
}

// uuid represents a player UUID (128-bit).
type uuid [16]byte

// Entity represents any entity in the world.
type Entity interface {
	ID() int32
	Type() registry.EntityType
	Position() *math.Vec3
		SetPosition(pos *math.Vec3)
	Rotation() (yaw, pitch float64)
	SetRotation(yaw, pitch float64)
	Velocity() *math.Vec3
	SetVelocity(vel *math.Vec3)
	OnGround() bool
	SetOnGround(onGround bool)
}

// BaseEntity implements common entity functionality.
type BaseEntity struct {
	id        int32
	entityType registry.EntityType
	position  *math.Vec3
	velocity  *math.Vec3
	yaw       float64
	pitch     float64
	onGround  bool
}

// NewBaseEntity creates a new base entity.
func NewBaseEntity(id int32, entityType registry.EntityType, pos *math.Vec3) *BaseEntity {
	return &BaseEntity{
		id:         id,
		entityType: entityType,
		position:   pos,
		velocity:   &math.Vec3{X: 0, Y: 0, Z: 0},
		yaw:        0,
		pitch:      0,
		onGround:   false,
	}
}

// ID returns the entity ID.
func (e *BaseEntity) ID() int32 {
	return e.id
}

// Type returns the entity type.
func (e *BaseEntity) Type() registry.EntityType {
	return e.entityType
}

// Position returns the entity position.
func (e *BaseEntity) Position() *math.Vec3 {
	return e.position
}

// SetPosition sets the entity position.
func (e *BaseEntity) SetPosition(pos *math.Vec3) {
	e.position = pos
}

// Rotation returns the entity rotation.
func (e *BaseEntity) Rotation() (yaw, pitch float64) {
	return e.yaw, e.pitch
}

// SetRotation sets the entity rotation.
func (e *BaseEntity) SetRotation(yaw, pitch float64) {
	e.yaw = yaw
	e.pitch = pitch
}

// Velocity returns the entity velocity.
func (e *BaseEntity) Velocity() *math.Vec3 {
	return e.velocity
}

// SetVelocity sets the entity velocity.
func (e *BaseEntity) SetVelocity(vel *math.Vec3) {
	e.velocity = vel
}

// OnGround returns whether the entity is on the ground.
func (e *BaseEntity) OnGround() bool {
	return e.onGround
}

// SetOnGround sets whether the entity is on the ground.
func (e *BaseEntity) SetOnGround(onGround bool) {
	e.onGround = onGround
}

// Player represents a player entity.
type Player struct {
	*BaseEntity
	entityID   int32
	uuid       uuid
	name       string
	health     float32
	food       int32
	saturation float32
	gameMode   GameMode
}

// GameMode represents a player's game mode.
type GameMode byte

const (
	GameModeSurvival GameMode = 0
	GameModeCreative GameMode = 1
	GameModeAdventure GameMode = 2
	GameModeSpectator GameMode = 3
)

func (g GameMode) String() string {
	switch g {
	case GameModeSurvival:
		return "survival"
	case GameModeCreative:
		return "creative"
	case GameModeAdventure:
		return "adventure"
	case GameModeSpectator:
		return "spectator"
	default:
		return fmt.Sprintf("unknown(%d)", g)
	}
}

// NewPlayer creates a new player.
func NewPlayer(id int32, playerUUID uuid, name string, pos *math.Vec3) *Player {
	return &Player{
		BaseEntity: NewBaseEntity(id, registry.EntityTypePlayer, pos),
		uuid:       playerUUID,
		name:       name,
		health:     20,
		food:       20,
		saturation: 5,
		gameMode:   GameModeSurvival,
	}
}

// UUID returns the player UUID.
func (p *Player) UUID() uuid {
	return p.uuid
}

// Name returns the player name.
func (p *Player) Name() string {
	return p.name
}

// SetName sets the player name.
func (p *Player) SetName(name string) {
	p.name = name
}

// Health returns the player health.
func (p *Player) Health() float32 {
	return p.health
}

// SetHealth sets the player health.
func (p *Player) SetHealth(health float32) {
	p.health = health
}

// Food returns the player food level.
func (p *Player) Food() int32 {
	return p.food
}

// SetFood sets the player food level.
func (p *Player) SetFood(food int32) {
	p.food = food
}

// Saturation returns the player saturation level.
func (p *Player) Saturation() float32 {
	return p.saturation
}

// SetSaturation sets the player saturation level.
func (p *Player) SetSaturation(saturation float32) {
	p.saturation = saturation
}

// GameMode returns the player game mode.
func (p *Player) GameMode() GameMode {
	return p.gameMode
}

// SetGameMode sets the player game mode.
func (p *Player) SetGameMode(mode GameMode) {
	p.gameMode = mode
}
