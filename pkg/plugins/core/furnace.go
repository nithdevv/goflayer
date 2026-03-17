package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// FurnaceData represents furnace state.
type FurnaceData struct {
	FuelProgress      float32
	CookProgress      float32
	FuelLeft          int16
	MaxFuelTime       int16
	CookTime          int16
	MaxCookTime       int16
}

// FurnacePlugin handles furnace interactions.
type FurnacePlugin struct {
	mu             sync.RWMutex
	ctx            *plugins.Context
	log            *logger.Logger
	openFurnacePos *math.BlockPos
	isOpen         bool
	furnaceData    *FurnaceData
}

// Metadata returns the plugin metadata.
func (p *FurnacePlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "furnace",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Furnace interaction",
		Dependencies: []string{"inventory"},
	}
}

// OnLoad initializes the furnace plugin.
func (p *FurnacePlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.isOpen = false
	p.furnaceData = &FurnaceData{}

	p.log.Info("Furnace plugin loaded")

	// Register event handlers
	p.ctx.On("furnace_data", p.handleFurnaceData)

	return nil
}

// OnUnload cleans up the furnace plugin.
func (p *FurnacePlugin) OnUnload() error {
	p.log.Info("Furnace plugin unloaded")
	return nil
}

// Open opens a furnace at a position.
func (p *FurnacePlugin) Open(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isOpen {
		return fmt.Errorf("furnace already open")
	}

	p.log.Info("Opening furnace at %v", pos)
	p.openFurnacePos = pos
	p.isOpen = true

	// TODO: Send open furnace packet
	p.ctx.Emit("furnace_opening", pos)

	return nil
}

// Close closes the currently open furnace.
func (p *FurnacePlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return fmt.Errorf("no furnace open")
	}

	p.log.Info("Closing furnace")
	p.isOpen = false

	// TODO: Send close furnace packet
	p.ctx.Emit("furnace_closing")

	return nil
}

// IsOpen returns whether a furnace is open.
func (p *FurnacePlugin) IsOpen() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOpen
}

// PutFuel puts fuel in the furnace.
func (p *FurnacePlugin) PutFuel(slotID, count int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no furnace open")
	}

	p.log.Info("Putting fuel from slot %d, count %d", slotID, count)

	// TODO: Implement fuel logic
	p.ctx.Emit("furnace_fuel_added", slotID, count)

	return nil
}

// PutInput puts items to smelt in the furnace.
func (p *FurnacePlugin) PutInput(slotID, count int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no furnace open")
	}

	p.log.Info("Putting input from slot %d, count %d", slotID, count)

	// TODO: Implement input logic
	p.ctx.Emit("furnace_input_added", slotID, count)

	return nil
}

// TakeOutput takes smelted items from the furnace.
func (p *FurnacePlugin) TakeOutput(count int) error {
	if !p.IsOpen() {
		return fmt.Errorf("no furnace open")
	}

	p.log.Info("Taking %d items from output", count)

	// TODO: Implement output logic
	p.ctx.Emit("furnace_output_taken", count)

	return nil
}

// GetData returns the current furnace data.
func (p *FurnacePlugin) GetData() *FurnaceData {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return a copy
	return &FurnaceData{
		FuelProgress: p.furnaceData.FuelProgress,
		CookProgress: p.furnaceData.CookProgress,
		FuelLeft:     p.furnaceData.FuelLeft,
		MaxFuelTime:  p.furnaceData.MaxFuelTime,
		CookTime:     p.furnaceData.CookTime,
		MaxCookTime:  p.furnaceData.MaxCookTime,
	}
}

// IsCooking returns true if the furnace is currently cooking.
func (p *FurnacePlugin) IsCooking() bool {
	data := p.GetData()
	return data.CookProgress > 0 && data.CookProgress < 1
}

// WaitUntilDone waits until cooking is complete.
func (p *FurnacePlugin) WaitUntilDone(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if !p.IsCooking() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("cooking not complete after timeout")
}

// Event handlers

func (p *FurnacePlugin) handleFurnaceData(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: Parse furnace data packet
	p.log.Debug("Furnace data updated")
}

// String returns a string representation of the furnace plugin.
func (p *FurnacePlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isOpen {
		return fmt.Sprintf("Furnace{open=true, cooking=%v}", p.IsCooking())
	}
	return "Furnace{open=false}"
}
