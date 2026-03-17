package core

import (
	"fmt"
	stdmath "math"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// MovementPlugin provides movement API.
type MovementPlugin struct {
	mu       sync.RWMutex
	ctx      *plugins.Context
	log      *logger.Logger
	position *math.Vec3
	yaw      float32
	pitch    float32
	onGround bool
}

// Metadata returns the plugin metadata.
func (p *MovementPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "movement",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Movement API (walk, jump, look)",
		Dependencies: []string{"physics"},
	}
}

// OnLoad initializes the movement plugin.
func (p *MovementPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.position = math.NewVec3(0, 0, 0)
	p.yaw = 0
	p.pitch = 0
	p.onGround = true

	p.log.Info("Movement plugin loaded")

	// Register event handlers
	p.ctx.On("player_position", p.handlePlayerPosition)
	p.ctx.On("player_look", p.handlePlayerLook)
	p.ctx.On("player_position_look", p.handlePlayerPositionLook)

	return nil
}

// OnUnload cleans up the movement plugin.
func (p *MovementPlugin) OnUnload() error {
	p.log.Info("Movement plugin unloaded")
	return nil
}

// GetPosition returns the current position.
func (p *MovementPlugin) GetPosition() *math.Vec3 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.position.Clone()
}

// GetYaw returns the current yaw.
func (p *MovementPlugin) GetYaw() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.yaw
}

// GetPitch returns the current pitch.
func (p *MovementPlugin) GetPitch() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.pitch
}

// SetPosition sets the player position.
func (p *MovementPlugin) SetPosition(pos *math.Vec3) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.position = pos.Clone()

	// TODO: Send position packet
	p.ctx.Emit("position_update", pos)

	return nil
}

// SetLook sets the yaw and pitch.
func (p *MovementPlugin) SetLook(yaw, pitch float32) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.yaw = yaw
	p.pitch = pitch

	// TODO: Send look packet
	p.ctx.Emit("look_update", yaw, pitch)

	return nil
}

// LookAt looks at a position.
func (p *MovementPlugin) LookAt(target *math.Vec3) error {
	dx := target.X - p.position.X
	dy := target.Y - p.position.Y
	dz := target.Z - p.position.Z

	horizontalDist := stdmath.Sqrt(dx*dx + dz*dz)

	yaw := float32(-stdmath.Atan2(dx, dz) * 180.0 / stdmath.Pi)
	pitch := float32(-stdmath.Atan2(dy, horizontalDist) * 180.0 / stdmath.Pi)

	return p.SetLook(yaw, pitch)
}

// Move moves in a direction.
func (p *MovementPlugin) Move(forward, strafe, vertical float64) error {
	// Convert movement direction to world coordinates
	yawRad := float64(p.yaw) * stdmath.Pi / 180.0

	cosYaw := stdmath.Cos(yawRad)
	sinYaw := stdmath.Sin(yawRad)

	dx := forward*sinYaw + strafe*cosYaw
	dy := vertical
	dz := forward*cosYaw - strafe*sinYaw

 newPos := p.position.Add(math.NewVec3(dx, dy, dz))

	return p.SetPosition(newPos)
}

// Walk forward.
func (p *MovementPlugin) Walk(distance float64) error {
	return p.Move(distance, 0, 0)
}

// Strafe sideways.
func (p *MovementPlugin) Strafe(distance float64) error {
	return p.Move(0, distance, 0)
}

// Jump jumps.
func (p *MovementPlugin) Jump() error {
	p.mu.RLock()
	onGround := p.onGround
	p.mu.RUnlock()

	if !onGround {
		return fmt.Errorf("not on ground")
	}

	p.ctx.Emit("jump")
	return nil
}

// Sprint starts sprinting.
func (p *MovementPlugin) Sprint() error {
	p.ctx.Emit("sprint_start")
	return nil
}

// StopSprinting stops sprinting.
func (p *MovementPlugin) StopSprinting() error {
	p.ctx.Emit("sprint_stop")
	return nil
}

// Sneak starts sneaking.
func (p *MovementPlugin) Sneak() error {
	p.ctx.Emit("sneak_start")
	return nil
}

// StopSneaking stops sneaking.
func (p *MovementPlugin) StopSneaking() error {
	p.ctx.Emit("sneak_stop")
	return nil
}

// Event handlers

func (p *MovementPlugin) handlePlayerPosition(args ...interface{}) {
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
}

func (p *MovementPlugin) handlePlayerLook(args ...interface{}) {
	if len(args) < 2 {
		return
	}

	yaw, ok1 := args[0].(float32)
	pitch, ok2 := args[1].(float32)
	if !ok1 || !ok2 {
		return
	}

	p.mu.Lock()
	p.yaw = yaw
	p.pitch = pitch
	if len(args) > 2 {
		if onGround, ok := args[2].(bool); ok {
			p.onGround = onGround
		}
	}
	p.mu.Unlock()
}

func (p *MovementPlugin) handlePlayerPositionLook(args ...interface{}) {
	if len(args) < 5 {
		return
	}

	x, ok1 := args[0].(float64)
	y, ok2 := args[1].(float64)
	z, ok3 := args[2].(float64)
	yaw, ok4 := args[3].(float32)
	pitch, ok5 := args[4].(float32)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return
	}

	p.mu.Lock()
	p.position.Set(x, y, z)
	p.yaw = yaw
	p.pitch = pitch
	if len(args) > 5 {
		if onGround, ok := args[5].(bool); ok {
			p.onGround = onGround
		}
	}
	p.mu.Unlock()

	p.ctx.Emit("moved", p.position, p.yaw, p.pitch)
}

// String returns a string representation of the movement plugin.
func (p *MovementPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Movement{pos=(%.2f, %.2f, %.2f), yaw=%.1f, pitch=%.1f}",
		p.position.X, p.position.Y, p.position.Z, p.yaw, p.pitch)
}
