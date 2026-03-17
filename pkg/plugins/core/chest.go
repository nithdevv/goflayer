package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// ChestPlugin handles chest interactions.
type ChestPlugin struct {
	mu             sync.RWMutex
	ctx            *plugins.Context
	log            *logger.Logger
	openChestPos   *math.BlockPos
	isOpen         bool
}

// Metadata returns the plugin metadata.
func (p *ChestPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "chest",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Chest interaction",
		Dependencies: []string{"inventory"},
	}
}

// OnLoad initializes the chest plugin.
func (p *ChestPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.isOpen = false

	p.log.Info("Chest plugin loaded")

	// Register event handlers
	p.ctx.On("chest_open", p.handleChestOpen)
	p.ctx.On("chest_close", p.handleChestClose)

	return nil
}

// OnUnload cleans up the chest plugin.
func (p *ChestPlugin) OnUnload() error {
	p.log.Info("Chest plugin unloaded")
	return nil
}

// Open opens a chest at a position.
func (p *ChestPlugin) Open(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isOpen {
		return fmt.Errorf("chest already open")
	}

	p.log.Info("Opening chest at %v", pos)
	p.openChestPos = pos
	p.isOpen = true

	// TODO: Send open chest packet
	p.ctx.Emit("chest_opening", pos)

	return nil
}

// Close closes the currently open chest.
func (p *ChestPlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return fmt.Errorf("no chest open")
	}

	p.log.Info("Closing chest")
	p.isOpen = false

	// TODO: Send close chest packet
	p.ctx.Emit("chest_closing")

	return nil
}

// IsOpen returns true if a chest is open.
func (p *ChestPlugin) IsOpen() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOpen
}

// GetPosition returns the position of the open chest.
func (p *ChestPlugin) GetPosition() *math.BlockPos {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.openChestPos
}

// Withdraw withdraws an item from the chest.
func (p *ChestPlugin) Withdraw(slotID, count int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no chest open")
	}

	p.log.Info("Withdrawing %d items from slot %d", count, slotID)

	// TODO: Implement withdraw logic
	p.ctx.Emit("chest_withdraw", slotID, count)

	return nil
}

// Deposit deposits an item into the chest.
func (p *ChestPlugin) Deposit(slotID, count int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no chest open")
	}

	p.log.Info("Depositing %d items from slot %d", count, slotID)

	// TODO: Implement deposit logic
	p.ctx.Emit("chest_deposit", slotID, count)

	return nil
}

// Event handlers

func (p *ChestPlugin) handleChestOpen(args ...interface{}) {
	p.log.Debug("Chest opened event")
}

func (p *ChestPlugin) handleChestClose(args ...interface{}) {
	p.mu.Lock()
	p.isOpen = false
	p.openChestPos = nil
	p.mu.Unlock()

	p.log.Debug("Chest closed event")
}

// String returns a string representation of the chest plugin.
func (p *ChestPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isOpen {
		return fmt.Sprintf("Chest{open=true, pos=%v}", p.openChestPos)
	}
	return "Chest{open=false}"
}
