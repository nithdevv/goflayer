package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// Block represents a block in the world.
type Block struct {
	Type     int
	Metadata int
	Position *math.BlockPos
}

// BlockFace represents a face of a block.
type BlockFace int

const (
	Bottom BlockFace = iota
	Top
	North
	South
	West
	East
)

// BlocksPlugin handles block interactions.
type BlocksPlugin struct {
	mu           sync.RWMutex
	ctx          *plugins.Context
	log          *logger.Logger
	targetBlock  *Block
	targetFace   BlockFace
	isDigging    bool
	digStartTime time.Time
	digDuration  time.Duration
}

// Metadata returns the plugin metadata.
func (p *BlocksPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "blocks",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Block interaction (dig, place)",
		Dependencies: []string{},
	}
}

// OnLoad initializes the blocks plugin.
func (p *BlocksPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.targetBlock = nil
	p.isDigging = false

	p.log.Info("Blocks plugin loaded")

	// Register event handlers
	p.ctx.On("block_change", p.handleBlockChange)
	p.ctx.On("multi_block_change", p.handleMultiBlockChange)

	return nil
}

// OnUnload cleans up the blocks plugin.
func (p *BlocksPlugin) OnUnload() error {
	p.log.Info("Blocks plugin unloaded")
	return nil
}

// DigAt starts digging a block at a position.
func (p *BlocksPlugin) DigAt(pos *math.BlockPos, forceLook bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isDigging {
		return fmt.Errorf("already digging")
	}

	p.log.Info("Starting to dig block at %v", pos)
	p.isDigging = true
	p.digStartTime = time.Now()
	p.digDuration = 0 // Will be determined by block type

	// TODO: Send start digging packet
	p.ctx.Emit("dig_start", pos)

	return nil
}

// Dig finishes the current digging operation.
func (p *BlocksPlugin) Dig() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isDigging {
		return fmt.Errorf("not digging")
	}

	p.log.Info("Finishing dig")
	p.isDigging = false

	// TODO: Send finish digging packet
	p.ctx.Emit("dig_complete", p.targetBlock)

	return nil
}

// CancelDig cancels the current digging operation.
func (p *BlocksPlugin) CancelDig() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isDigging {
		return fmt.Errorf("not digging")
	}

	p.log.Info("Cancelling dig")
	p.isDigging = false

	// TODO: Send cancel digging packet
	p.ctx.Emit("dig_cancelled")

	return nil
}

// Place places a block at a position.
func (p *BlocksPlugin) Place(pos *math.BlockPos, face BlockFace) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("Placing block at %v, face %d", pos, face)

	// TODO: Send place block packet
	p.ctx.Emit("block_place", pos, face)

	return nil
}

// PlaceWithHand places a block using the item in hand.
func (p *BlocksPlugin) PlaceWithHand() error {
	// TODO: Implement placing with hand-held item
	return fmt.Errorf("not implemented")
}

// ActivateBlock activates a block (right-click interaction).
func (p *BlocksPlugin) ActivateBlock(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("Activating block at %v", pos)

	// TODO: Send block activation packet
	p.ctx.Emit("block_activate", pos)

	return nil
}

// GetBlockAt returns the block at a position.
func (p *BlocksPlugin) GetBlockAt(x, y, z int) *Block {
	// TODO: Query from world state
	return &Block{
		Type:     0,
		Metadata: 0,
		Position: math.NewBlockPos(x, y, z),
	}
}

// CanSeeBlock checks if a block is visible from the player's position.
func (p *BlocksPlugin) CanSeeBlock(pos *math.BlockPos) bool {
	// TODO: Implement line-of-sight check
	return true
}

// Event handlers

func (p *BlocksPlugin) handleBlockChange(args ...interface{}) {
	p.log.Debug("Block changed")
	p.ctx.Emit("block_updated", args)
}

func (p *BlocksPlugin) handleMultiBlockChange(args ...interface{}) {
	p.log.Debug("Multiple blocks changed")
	p.ctx.Emit("blocks_updated", args)
}

// String returns a string representation of the blocks plugin.
func (p *BlocksPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Blocks{digging=%v}", p.isDigging)
}
