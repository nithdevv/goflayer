package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// Item represents an item in the inventory.
type Item struct {
	TypeID    int32
	Count     int32
	Damage    int16
	NBT       map[string]interface{}
	DisplayName string
}

// Slot represents an inventory slot.
type Slot struct {
	Item   *Item
	SlotID int
}

// Inventory represents an inventory.
type Inventory struct {
	Slots    []*Slot
	Size     int
	Title    string
	Type     int
}

// WindowType represents different window types.
type WindowType int

const (
	GenericInventory WindowType = iota
	Chest
	Workbench
	Furnace
	Dispenser
	EnchantmentTable
	BrewingStand
	Villager
	Beacon
	Anvil
	Hopper
	Dropper
)

// InventoryPlugin manages player inventory.
type InventoryPlugin struct {
	mu              sync.RWMutex
	ctx             *plugins.Context
	log             *logger.Logger
	playerInventory *Inventory
	openWindow      *Inventory
	windowID        int
	inventoryMap    map[int]*Inventory
}

// Metadata returns the plugin metadata.
func (p *InventoryPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "inventory",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Player inventory management",
		Dependencies: []string{},
	}
}

// OnLoad initializes the inventory plugin.
func (p *InventoryPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.playerInventory = &Inventory{
		Slots: make([]*Slot, 46), // 36 main + 1 armor + 1 offhand + 8 crafting
		Size:  46,
		Title: "Player Inventory",
		Type:  0,
	}
	p.windowID = 0
	p.inventoryMap = make(map[int]*Inventory)

	p.log.Info("Inventory plugin loaded")

	// Register event handlers
	p.ctx.On("window_open", p.handleWindowOpen)
	p.ctx.On("window_close", p.handleWindowClose)
	p.ctx.On("window_items", p.handleWindowItems)
	p.ctx.On("slot_change", p.handleSlotChange)
	p.ctx.On("held_item_change", p.handleHeldItemChange)

	return nil
}

// OnUnload cleans up the inventory plugin.
func (p *InventoryPlugin) OnUnload() error {
	p.log.Info("Inventory plugin unloaded")
	return nil
}

// GetSlot returns a slot from the player inventory.
func (p *InventoryPlugin) GetSlot(slotID int) *Slot {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if slotID < 0 || slotID >= len(p.playerInventory.Slots) {
		return nil
	}
	return p.playerInventory.Slots[slotID]
}

// SetSlot sets a slot in the player inventory.
func (p *InventoryPlugin) SetSlot(slotID int, item *Item) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if slotID < 0 || slotID >= len(p.playerInventory.Slots) {
		return fmt.Errorf("invalid slot ID: %d", slotID)
	}

	slot := &Slot{
		Item:   item,
		SlotID: slotID,
	}
	p.playerInventory.Slots[slotID] = slot

	// TODO: Send slot change packet
	p.ctx.Emit("slot_set", slotID, item)

	return nil
}

// GetHotbarSlot returns a hotbar slot (0-8).
func (p *InventoryPlugin) GetHotbarSlot(slotID int) *Slot {
	if slotID < 0 || slotID > 8 {
		return nil
	}
	return p.GetSlot(slotID)
}

// GetHeldItem returns the item in the player's hand.
func (p *InventoryPlugin) GetHeldItem() *Item {
	// TODO: Track selected hotbar slot
	return p.GetHotbarSlot(0).Item
}

// CountItem returns the count of a specific item type.
func (p *InventoryPlugin) CountItem(itemType int32) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, slot := range p.playerInventory.Slots {
		if slot != nil && slot.Item != nil && slot.Item.TypeID == itemType {
			count += int(slot.Item.Count)
		}
	}
	return count
}

// HasItem returns true if the inventory contains an item.
func (p *InventoryPlugin) HasItem(itemType int32) bool {
	return p.CountItem(itemType) > 0
}

// FindItem returns the first slot containing a specific item type.
func (p *InventoryPlugin) FindItem(itemType int32) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for i, slot := range p.playerInventory.Slots {
		if slot != nil && slot.Item != nil && slot.Item.TypeID == itemType {
			return i
		}
	}
	return -1
}

// FindItems returns all slots containing a specific item type.
func (p *InventoryPlugin) FindItems(itemType int32) []int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	slots := make([]int, 0)
	for i, slot := range p.playerInventory.Slots {
		if slot != nil && slot.Item != nil && slot.Item.TypeID == itemType {
			slots = append(slots, i)
		}
	}
	return slots
}

// ClickSlot clicks a slot in the current window.
func (p *InventoryPlugin) ClickSlot(windowID, slotID, mouseButton, mode int) error {
	p.log.Debug("Clicking slot %d in window %d", slotID, windowID)

	// TODO: Send click window packet
	p.ctx.Emit("slot_clicked", windowID, slotID, mouseButton)

	return nil
}

// DropItem drops an item from a slot.
func (p *InventoryPlugin) DropItem(slotID int, count int) error {
	p.log.Info("Dropping %d items from slot %d", count, slotID)

	// TODO: Send drop item packet
	p.ctx.Emit("item_dropped", slotID, count)

	return nil
}

// CloseWindow closes the current open window.
func (p *InventoryPlugin) CloseWindow() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.openWindow == nil {
		return fmt.Errorf("no window open")
	}

	p.log.Info("Closing window %d", p.windowID)

	// TODO: Send close window packet
	p.ctx.Emit("window_closed", p.windowID)

	p.openWindow = nil
	p.windowID = 0

	return nil
}

// GetOpenWindow returns the currently open window.
func (p *InventoryPlugin) GetOpenWindow() *Inventory {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.openWindow
}

// Event handlers

func (p *InventoryPlugin) handleWindowOpen(args ...interface{}) {
	p.log.Info("Window opened")
	// TODO: Parse window open packet
}

func (p *InventoryPlugin) handleWindowClose(args ...interface{}) {
	p.mu.Lock()
	p.openWindow = nil
	p.windowID = 0
	p.mu.Unlock()

	p.log.Info("Window closed")
}

func (p *InventoryPlugin) handleWindowItems(args ...interface{}) {
	p.log.Debug("Window items updated")
	// TODO: Parse window items packet
}

func (p *InventoryPlugin) handleSlotChange(args ...interface{}) {
	p.log.Debug("Slot changed")
	// TODO: Update slot
}

func (p *InventoryPlugin) handleHeldItemChange(args ...interface{}) {
	p.log.Debug("Held item changed")
	// TODO: Update held item
}

// String returns a string representation of the inventory plugin.
func (p *InventoryPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	itemCount := 0
	for _, slot := range p.playerInventory.Slots {
		if slot != nil && slot.Item != nil {
			itemCount++
		}
	}

	return fmt.Sprintf("Inventory{items=%d, windowOpen=%v}",
		itemCount, p.openWindow != nil)
}
