// Package physics implements movement and collision physics.
package physics

import (
	stdmath "math"
	"sync"

	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
)

// Plugin implements physics simulation.
type Plugin struct {
	*plugins.BasePlugin

	// Physics state
	state     State
	stateMu   sync.RWMutex
}

// State represents the physics state.
type State struct {
	// Position
	Position *math.Vec3

	// Velocity
	Velocity *math.Vec3

	// Rotation (yaw, pitch)
	Yaw   float32
	Pitch float32

	// On ground
	OnGround bool
}

// NewPlugin creates a new physics plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("physics", "1.0.0")

	return &Plugin{
		BasePlugin: base,
		state: State{
			Position: math.NewVec3(0, 0, 0),
			Velocity: math.NewVec3(0, 0, 0),
			OnGround:  true,
		},
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Subscribe to position packets
	p.On("packet", p.handlePacket)
	p.On("position_updated", p.onPositionUpdated)

	return nil
}

// handlePacket handles incoming packets.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet := data[0].(*protocol.Packet)

	switch packet.ID {
	case 0x2E: // Set Player Position
		p.handleSetPosition(packet)
	case 0x2F: // Set Position and Rotation
		p.handleSetPositionRotation(packet)
	case 0x30: // Set Rotation
		p.handleSetRotation(packet)
	}
}

// handleSetPosition handles the Set Player Position packet.
func (p *Plugin) handleSetPosition(packet *protocol.Packet) {
	// TODO: Parse position from packet data
	// For now, just emit event
	p.Emit("physics_position_updated")
}

// handleSetPositionRotation handles the Set Position and Rotation packet.
func (p *Plugin) handleSetPositionRotation(packet *protocol.Packet) {
	// TODO: Parse position and rotation from packet data
	p.Emit("physics_position_updated")
	p.Emit("rotation_updated")
}

// handleSetRotation handles the Set Rotation packet.
func (p *Plugin) handleSetRotation(packet *protocol.Packet) {
	// TODO: Parse rotation from packet data
	p.Emit("rotation_updated")
}

// onPositionUpdated handles position updates.
func (p *Plugin) onPositionUpdated(data ...interface{}) {
	// TODO: Update physics state from event data
}

// setPosition sets the player position.
func (p *Plugin) setPosition(x, y, z float64) {
	p.stateMu.Lock()
	p.state.Position.Set(x, y, z)
	p.stateMu.Unlock()

	p.Emit("physics_position_updated", x, y, z)
}

// Position returns the player position.
func (p *Plugin) Position() *math.Vec3 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Position.Clone()
}

// Velocity returns the player velocity.
func (p *Plugin) Velocity() *math.Vec3 {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Velocity.Clone()
}

// Rotation returns the player rotation (yaw, pitch).
func (p *Plugin) Rotation() (yaw, pitch float32) {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.Yaw, p.state.Pitch
}

// IsOnGround returns true if the player is on the ground.
func (p *Plugin) IsOnGround() bool {
	p.stateMu.RLock()
	defer p.stateMu.RUnlock()
	return p.state.OnGround
}

// LookAt makes the player look at a position.
func (p *Plugin) LookAt(target *math.Vec3) {
	p.stateMu.Lock()
	defer p.stateMu.Unlock()

	// Calculate direction
	dx := target.X - p.state.Position.X
	dy := target.Y - p.state.Position.Y
	dz := target.Z - p.state.Position.Z

	// Calculate yaw and pitch
	p.state.Yaw = float32(stdmath.Atan2(dz, dx) * 180.0 / stdmath.Pi) - 90.0
	p.state.Pitch = float32(stdmath.Atan2(-dy, stdmath.Sqrt(dx*dx+dz*dz)) * 180.0 / stdmath.Pi)

	// TODO: Send rotation packet to server
	p.Emit("look_at", target)
}

// SetPosition sets the player position and sends it to the server.
func (p *Plugin) SetPosition(pos *math.Vec3) {
	p.stateMu.Lock()
	p.state.Position.Set(pos.X, pos.Y, pos.Z)
	p.stateMu.Unlock()

	// TODO: Send position packet to server
	p.Emit("position_set", pos)
}

// SetRotation sets the player rotation and sends it to the server.
func (p *Plugin) SetRotation(yaw, pitch float32) {
	p.stateMu.Lock()
	p.state.Yaw = yaw
	p.state.Pitch = pitch
	p.stateMu.Unlock()

	// TODO: Send rotation packet to server
	p.Emit("rotation_set", yaw, pitch)
}
