package core

import (
	"fmt"
	stdmath "math"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// RaycastResult represents the result of a raycast.
type RaycastResult struct {
	Hit       bool
	Position  *math.Vec3
	Block     *math.BlockPos
	Face      int // 0-5 for block faces
	Entity    *Entity
	Distance  float64
}

// RaycastPlugin handles block/entity raycasting.
type RaycastPlugin struct {
	mu            sync.RWMutex
	ctx           *plugins.Context
	log           *logger.Logger
	maxDistance   float64
	liquidsVisible bool
}

// Metadata returns the plugin metadata.
func (p *RaycastPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "raycast",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Block/entity raycasting",
		Dependencies: []string{"entities"},
	}
}

// OnLoad initializes the raycast plugin.
func (p *RaycastPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.maxDistance = 6
	p.liquidsVisible = false

	p.log.Info("Raycast plugin loaded")

	return nil
}

// OnUnload cleans up the raycast plugin.
func (p *RaycastPlugin) OnUnload() error {
	p.log.Info("Raycast plugin unloaded")
	return nil
}

// SetMaxDistance sets the maximum raycast distance.
func (p *RaycastPlugin) SetMaxDistance(dist float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.maxDistance = dist
	p.log.Debug("Max distance set to %.2f", dist)
}

// SetLiquidsVisible sets whether liquids are visible to raycast.
func (p *RaycastPlugin) SetLiquidsVisible(visible bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.liquidsVisible = visible
	p.log.Debug("Liquids visible set to %v", visible)
}

// Raycast performs a raycast from a position in a direction.
func (p *RaycastPlugin) Raycast(start, direction *math.Vec3) *RaycastResult {
	p.mu.RLock()
	maxDist := p.maxDistance
	p.mu.RUnlock()

	return p.raycastWithLimit(start, direction, maxDist)
}

// RaycastWithLimit performs a raycast with a custom distance limit.
func (p *RaycastPlugin) RaycastWithLimit(start, direction *math.Vec3, maxDist float64) *RaycastResult {
	return p.raycastWithLimit(start, direction, maxDist)
}

// raycastWithLimit performs the actual raycast.
func (p *RaycastPlugin) raycastWithLimit(start, direction *math.Vec3, maxDist float64) *RaycastResult {
	result := &RaycastResult{
		Hit:      false,
		Position: start.Clone(),
	}

	// Normalize direction
	dir := direction.Normalize()

	// Step size for raycast
	stepSize := 0.1
	current := start.Clone()
	distance := 0.0

	// Raycast loop
	for distance < maxDist {
		// Move forward
		current.X += dir.X * stepSize
		current.Y += dir.Y * stepSize
		current.Z += dir.Z * stepSize
		distance += stepSize

		// Check for block collision
		blockPos := current.ToBlockPos()
		if p.checkBlockCollision(blockPos) {
			result.Hit = true
			result.Position = current
			result.Block = blockPos
			result.Face = p.calculateHitFace(dir, blockPos)
			result.Distance = distance
			break
		}

		// TODO: Check for entity collision
	}

	return result
}

// checkBlockCollision checks if a block at a position is solid.
func (p *RaycastPlugin) checkBlockCollision(pos *math.BlockPos) bool {
	// TODO: Implement proper block collision check
	// For now, assume y < 0 is solid (void)
	return pos.Y < 0
}

// calculateHitFace calculates which face of the block was hit.
func (p *RaycastPlugin) calculateHitFace(dir *math.Vec3, blockPos *math.BlockPos) int {
	// Calculate which face was hit based on direction
	blockCenter := blockPos.ToVec3()
	toBlock := blockCenter.Sub(dir)

	absX := stdmath.Abs(toBlock.X)
	absY := stdmath.Abs(toBlock.Y)
	absZ := stdmath.Abs(toBlock.Z)

	if absX > absY && absX > absZ {
		if dir.X > 0 {
			return 0 // West
		}
		return 1 // East
	}

	if absY > absX && absY > absZ {
		if dir.Y > 0 {
			return 2 // Bottom
		}
		return 3 // Top
	}

	if dir.Z > 0 {
		return 4 // North
	}
	return 5 // South
}

// RaycastEntities performs raycast against entities only.
func (p *RaycastPlugin) RaycastEntities(start, direction *math.Vec3, maxDist float64) *RaycastResult {
	result := &RaycastResult{
		Hit:      false,
		Position: start.Clone(),
	}

	// TODO: Implement entity raycast
	// Get entities plugin and check for entity intersections

	return result
}

// RaycastBlocks performs raycast against blocks only.
func (p *RaycastPlugin) RaycastBlocks(start, direction *math.Vec3, maxDist float64) *RaycastResult {
	return p.raycastWithLimit(start, direction, maxDist)
}

// GetTargetBlock returns the block the player is looking at.
func (p *RaycastPlugin) GetTargetBlock() *math.BlockPos {
	// TODO: Get player position and look direction
	// This requires integration with movement plugin
	return nil
}

// CanSee returns true if there's a clear line of sight to a position.
func (p *RaycastPlugin) CanSee(start, target *math.Vec3) bool {
	direction := target.Sub(start)
	result := p.Raycast(start, direction)
	return !result.Hit || result.Position.DistanceTo(target) < 0.5
}

// String returns a string representation of the raycast plugin.
func (p *RaycastPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Raycast{maxDist=%.2f, liquidsVisible=%v}",
		p.maxDistance, p.liquidsVisible)
}
