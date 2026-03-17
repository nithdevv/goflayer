package core

import (
	"fmt"
	stdmath "math"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// PhysicsConfig contains physics configuration.
type PhysicsConfig struct {
	Gravity        float64
	TerminalVelocity float64
	WalkSpeed       float64
	JumpForce       float64
	StepHeight      float64
	SprintSpeed     float64
}

// PhysicsPlugin handles basic physics simulation.
type PhysicsPlugin struct {
	mu        sync.RWMutex
	ctx       *plugins.Context
	log       *logger.Logger
	config    PhysicsConfig
	position  *math.Vec3
	velocity  *math.Vec3
	onGround  bool
	isSprinting bool
	lastUpdate time.Time
}

// Metadata returns the plugin metadata.
func (p *PhysicsPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "physics",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Basic physics (gravity, collision detection)",
		Dependencies: []string{},
	}
}

// OnLoad initializes the physics plugin.
func (p *PhysicsPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.config = PhysicsConfig{
		Gravity:          0.08,
		TerminalVelocity: -3.92,
		WalkSpeed:        0.1,
		JumpForce:        0.42,
		StepHeight:       0.6,
		SprintSpeed:      0.13,
	}
	p.position = math.NewVec3(0, 0, 0)
	p.velocity = math.NewVec3(0, 0, 0)
	p.onGround = false
	p.lastUpdate = time.Now()

	p.log.Info("Physics plugin loaded")

	// Register event handlers
	p.ctx.On("player_position", p.handlePlayerPosition)
	p.ctx.On("entity_velocity", p.handleEntityVelocity)

	return nil
}

// OnUnload cleans up the physics plugin.
func (p *PhysicsPlugin) OnUnload() error {
	p.log.Info("Physics plugin unloaded")
	return nil
}

// GetPosition returns the current position.
func (p *PhysicsPlugin) GetPosition() *math.Vec3 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.position.Clone()
}

// GetVelocity returns the current velocity.
func (p *PhysicsPlugin) GetVelocity() *math.Vec3 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.velocity.Clone()
}

// IsOnGround returns true if on the ground.
func (p *PhysicsPlugin) IsOnGround() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.onGround
}

// SetPosition sets the player position.
func (p *PhysicsPlugin) SetPosition(pos *math.Vec3) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.position = pos.Clone()
	p.ctx.Emit("physics_position_update", p.position)
}

// SetVelocity sets the player velocity.
func (p *PhysicsPlugin) SetVelocity(vel *math.Vec3) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.velocity = vel.Clone()
	p.ctx.Emit("physics_velocity_update", p.velocity)
}

// ApplyGravity applies gravity to velocity.
func (p *PhysicsPlugin) ApplyGravity(delta float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.onGround {
		p.velocity.Y -= p.config.Gravity * delta
		if p.velocity.Y < p.config.TerminalVelocity {
			p.velocity.Y = p.config.TerminalVelocity
		}
	}
}

// ApplyFriction applies friction to horizontal velocity.
func (p *PhysicsPlugin) ApplyFriction(friction float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.velocity.X *= friction
	p.velocity.Z *= friction
}

// SimulateStep simulates one physics step.
func (p *PhysicsPlugin) SimulateStep(delta float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Apply gravity
	if !p.onGround {
		p.velocity.Y -= p.config.Gravity * delta
		if p.velocity.Y < p.config.TerminalVelocity {
			p.velocity.Y = p.config.TerminalVelocity
		}
	}

	// Apply air resistance (0.91 in Minecraft)
	p.velocity.X *= 0.91
	p.velocity.Y *= 0.98
	p.velocity.Z *= 0.91

	// Update position
	p.position.X += p.velocity.X
	p.position.Y += p.velocity.Y
	p.position.Z += p.velocity.Z

	// Check for ground collision (simplified)
	// TODO: Implement proper collision detection
	if p.position.Y < 0 {
		p.position.Y = 0
		p.velocity.Y = 0
		p.onGround = true
	}

	p.ctx.Emit("physics_tick", p.position, p.velocity)
}

// Jump makes the player jump.
func (p *PhysicsPlugin) Jump() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.onGround {
		p.velocity.Y = p.config.JumpForce
		p.onGround = false
		p.ctx.Emit("physics_jump")
	}
}

// SetSprinting sets sprinting state.
func (p *PhysicsPlugin) SetSprinting(sprinting bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.isSprinting = sprinting
	p.ctx.Emit("physics_sprint_changed", sprinting)
}

// GetSpeed returns the current movement speed.
func (p *PhysicsPlugin) GetSpeed() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isSprinting {
		return p.config.SprintSpeed
	}
	return p.config.WalkSpeed
}

// CheckCollision checks for collisions at a position.
func (p *PhysicsPlugin) CheckCollision(pos *math.Vec3) bool {
	// TODO: Implement proper collision detection with blocks
	// For now, return false (no collision)
	return false
}

// Raycast performs a raycast from a position in a direction.
func (p *PhysicsPlugin) Raycast(start, direction *math.Vec3, maxDist float64) (*math.Vec3, bool) {
	// Simple raycast implementation
	step := 0.1
	current := start.Clone()

	for i := 0.0; i < maxDist; i += step {
		current.X += direction.X * step
		current.Y += direction.Y * step
		current.Z += direction.Z * step

		if p.CheckCollision(current) {
			return current, true
		}
	}

	return nil, false
}

// Event handlers

func (p *PhysicsPlugin) handlePlayerPosition(args ...interface{}) {
	if len(args) < 3 {
		return
	}

	x, ok1 := args[0].(float64)
	y, ok2 := args[1].(float64)
	z, ok3 := args[2].(float64)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	p.mu.Lock()
	p.position.Set(x, y, z)
	if len(args) > 3 {
		if onGround, ok := args[3].(bool); ok {
			p.onGround = onGround
		}
	}
	p.mu.Unlock()

	p.ctx.Emit("physics_position_update", p.position)
}

func (p *PhysicsPlugin) handleEntityVelocity(args ...interface{}) {
	if len(args) < 3 {
		return
	}

	vx, ok1 := args[0].(float64)
	vy, ok2 := args[1].(float64)
	vz, ok3 := args[2].(float64)
	if !ok1 || !ok2 || !ok3 {
		return
	}

	p.mu.Lock()
	p.velocity.Set(vx, vy, vz)
	p.mu.Unlock()

	p.ctx.Emit("physics_velocity_update", p.velocity)
}

// DistanceTo returns the distance to a position.
func (p *PhysicsPlugin) DistanceTo(pos *math.Vec3) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.position.DistanceTo(pos)
}

// DirectionTo returns the normalized direction vector to a position.
func (p *PhysicsPlugin) DirectionTo(pos *math.Vec3) *math.Vec3 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	dir := pos.Sub(p.position)
	return dir.Normalize()
}

// LookAt sets the yaw and pitch to look at a position.
func (p *PhysicsPlugin) LookAt(target *math.Vec3) (yaw, pitch float64) {
	pos := p.GetPosition()

	dx := target.X - pos.X
	dy := target.Y - pos.Y
	dz := target.Z - pos.Z

	horizontalDist := stdmath.Sqrt(dx*dx + dz*dz)

	yaw = -stdmath.Atan2(dx, dz) * 180.0 / stdmath.Pi
	pitch = -stdmath.Atan2(dy, horizontalDist) * 180.0 / stdmath.Pi

	return yaw, pitch
}

// String returns a string representation of the physics state.
func (p *PhysicsPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Physics{pos=(%.2f, %.2f, %.2f), vel=(%.2f, %.2f, %.2f), onGround=%v}",
		p.position.X, p.position.Y, p.position.Z,
		p.velocity.X, p.velocity.Y, p.velocity.Z,
		p.onGround)
}
