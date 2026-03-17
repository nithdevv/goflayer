package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// AnvilPlugin handles anvil interactions.
type AnvilPlugin struct {
	mu             sync.RWMutex
	ctx            *plugins.Context
	log            *logger.Logger
	openAnvilPos   *math.BlockPos
	isOpen         bool
	itemCost       int
}

// Metadata returns the plugin metadata.
func (p *AnvilPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "anvil",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Anvil interaction",
		Dependencies: []string{"inventory"},
	}
}

// OnLoad initializes the anvil plugin.
func (p *AnvilPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.isOpen = false
	p.itemCost = 0

	p.log.Info("Anvil plugin loaded")

	// Register event handlers
	p.ctx.On("anvil_data", p.handleAnvilData)

	return nil
}

// OnUnload cleans up the anvil plugin.
func (p *AnvilPlugin) OnUnload() error {
	p.log.Info("Anvil plugin unloaded")
	return nil
}

// Open opens an anvil at a position.
func (p *AnvilPlugin) Open(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isOpen {
		return fmt.Errorf("anvil already open")
	}

	p.log.Info("Opening anvil at %v", pos)
	p.openAnvilPos = pos
	p.isOpen = true

	// TODO: Send open anvil packet
	p.ctx.Emit("anvil_opening", pos)

	return nil
}

// Close closes the currently open anvil.
func (p *AnvilPlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return fmt.Errorf("no anvil open")
	}

	p.log.Info("Closing anvil")
	p.isOpen = false
	p.itemCost = 0

	// TODO: Send close anvil packet
	p.ctx.Emit("anvil_closing")

	return nil
}

// Rename renames an item.
func (p *AnvilPlugin) Rename(itemSlot int, newName string) error {
	if !p.IsOpen() {
		return fmt.Errorf("no anvil open")
	}

	p.log.Info("Renaming item in slot %d to '%s'", itemSlot, newName)

	// TODO: Implement rename logic
	p.ctx.Emit("item_rename", itemSlot, newName)

	return nil
}

// Combine combines two items.
func (p *AnvilPlugin) Combine(sourceSlot, targetSlot int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no anvil open")
	}

	p.log.Info("Combining items from slots %d and %d", sourceSlot, targetSlot)

	// TODO: Implement combine logic
	p.ctx.Emit("items_combine", sourceSlot, targetSlot)

	return nil
}

// Repair repairs an item.
func (p *AnvilPlugin) Repair(itemSlot int, materialSlot int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no anvil open")
	}

	p.log.Info("Repairing item in slot %d with material from slot %d", itemSlot, materialSlot)

	// TODO: Implement repair logic
	p.ctx.Emit("item_repair", itemSlot, materialSlot)

	return nil
}

// GetCost returns the cost in XP levels.
func (p *AnvilPlugin) GetCost() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.itemCost
}

// IsOpen returns true if an anvil is open.
func (p *AnvilPlugin) IsOpen() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOpen
}

// Event handlers

func (p *AnvilPlugin) handleAnvilData(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: Parse anvil data packet
	p.log.Debug("Anvil data updated")
}

// String returns a string representation of the anvil plugin.
func (p *AnvilPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Anvil{open=%v, cost=%d}", p.isOpen, p.itemCost)
}
