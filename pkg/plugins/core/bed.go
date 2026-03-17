package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// BedPlugin handles bed interactions.
type BedPlugin struct {
	mu          sync.RWMutex
	ctx         *plugins.Context
	log         *logger.Logger
	bedPos      *math.BlockPos
	isSleeping  bool
}

// Metadata returns the plugin metadata.
func (p *BedPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "bed",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Bed interaction",
		Dependencies: []string{},
	}
}

// OnLoad initializes the bed plugin.
func (p *BedPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.isSleeping = false

	p.log.Info("Bed plugin loaded")

	// Register event handlers
	p.ctx.On("sleep_start", p.handleSleepStart)
	p.ctx.On("sleep_stop", p.handleSleepStop)

	return nil
}

// OnUnload cleans up the bed plugin.
func (p *BedPlugin) OnUnload() error {
	p.log.Info("Bed plugin unloaded")
	return nil
}

// Sleep sleeps in a bed at a position.
func (p *BedPlugin) Sleep(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isSleeping {
		return fmt.Errorf("already sleeping")
	}

	p.log.Info("Sleeping in bed at %v", pos)
	p.bedPos = pos
	p.isSleeping = true

	// TODO: Send use bed packet
	p.ctx.Emit("bed_sleeping", pos)

	return nil
}

// Wake wakes up from sleeping.
func (p *BedPlugin) Wake() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isSleeping {
		return fmt.Errorf("not sleeping")
	}

	p.log.Info("Waking up")
	p.isSleeping = false

	// TODO: Send wake packet
	p.ctx.Emit("bed_wake")

	return nil
}

// IsSleeping returns true if currently sleeping.
func (p *BedPlugin) IsSleeping() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isSleeping
}

// GetBedPosition returns the bed position.
func (p *BedPlugin) GetBedPosition() *math.BlockPos {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.bedPos
}

// Event handlers

func (p *BedPlugin) handleSleepStart(args ...interface{}) {
	p.mu.Lock()
	p.isSleeping = true
	p.mu.Unlock()

	p.log.Info("Started sleeping")
}

func (p *BedPlugin) handleSleepStop(args ...interface{}) {
	p.mu.Lock()
	p.isSleeping = false
	p.bedPos = nil
	p.mu.Unlock()

	p.log.Info("Stopped sleeping")
}

// String returns a string representation of the bed plugin.
func (p *BedPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isSleeping {
		return fmt.Sprintf("Bed{sleeping=true, pos=%v}", p.bedPos)
	}
	return "Bed{sleeping=false}"
}
