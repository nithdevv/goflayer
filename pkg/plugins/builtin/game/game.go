// Package game implements the game state plugin.
// It tracks basic game state like player position, health, level, etc.
package game

import (
	"sync"

	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
)

// Plugin tracks game state.
type Plugin struct {
	*plugins.BasePlugin

	// Game state
	state     State
	stateMu   sync.RWMutex
}

// State represents the current game state.
type State struct {
	// Player state
	Level          int32
	GameMode       int32
	Dimension      int32
	Difficulty     int32
	MaxPlayers     int32

	// Player properties
	Health         float32
	Food           int32
	Saturation     float32
	Experience     float32
	LevelProgress  int32

	// Position
	X, Y, Z        float64
	Yaw, Pitch     float32
	OnGround       bool
}

// NewPlugin creates a new game plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("game", "1.0.0")

	return &Plugin{
		BasePlugin: base,
		state: State{
			GameMode:    0, // Survival
			Dimension:   0, // Overworld
			Difficulty:  1, // Normal
			Health:      20.0,
			Food:        20,
			Saturation:  5.0,
		},
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Subscribe to relevant packets
	p.On("packet", p.handlePacket)

	return nil
}

// handlePacket handles incoming packets.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet := data[0].(*protocol.Packet)

	switch packet.ID {
	case 0x2E: // Set Player Position packet (Play, 1.20.1)
		p.handleSetPosition(packet)
	case 0x2F: // Set Position and Rotation packet
		p.handleSetPositionRotation(packet)
	case 0x30: // Set Rotation packet
		p.handleSetRotation(packet)
	case 0x08: // Set Health packet
		p.handleSetHealth(packet)
	}
}

// handleSetPosition handles the Set Player Position packet.
func (p *Plugin) handleSetPosition(packet *protocol.Packet) {
	// TODO: Parse position from packet data
	// For now, we'll emit an event
	p.Emit("position_updated")
}

// handleSetPositionRotation handles the Set Position and Rotation packet.
func (p *Plugin) handleSetPositionRotation(packet *protocol.Packet) {
	// TODO: Parse position and rotation from packet data
	p.Emit("position_updated")
	p.Emit("rotation_updated")
}

// handleSetRotation handles the Set Rotation packet.
func (p *Plugin) handleSetRotation(packet *protocol.Packet) {
	// TODO: Parse rotation from packet data
	p.Emit("rotation_updated")
}

// handleSetHealth handles the Set Health packet.
func (p *Plugin) handleSetHealth(packet *protocol.Packet) {
	// TODO: Parse health from packet data
	p.Emit("health_updated")
}

// State returns the current game state.
func (p *Plugin) State() *State {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()

	// Return a copy
	state := p.state
	return &state
}

// Position returns the player's position.
func (p *Plugin) Position() (x, y, z float64) {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.X, p.state.Y, p.state.Z
}

// Rotation returns the player's rotation.
func (p *Plugin) Rotation() (yaw, pitch float32) {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Yaw, p.state.Pitch
}

// Health returns the player's health.
func (p *Plugin) Health() float32 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Health
}

// Food returns the player's food level.
func (p *Plugin) Food() int32 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Food
}

// GameMode returns the player's game mode.
func (p *Plugin) GameMode() int32 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.GameMode
}

// Dimension returns the player's dimension.
func (p *Plugin) Dimension() int32 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Dimension
}
