// Package inventory implements player inventory handling for Minecraft.
// This plugin handles managing the player's inventory slots.
package inventory

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
	"github.com/nithdevv/goflayer/pkg/protocol/codec"
)

// Plugin handles inventory operations.
type Plugin struct {
	*plugins.BasePlugin

	// Player inventory items (36 slots: 0-8 hotbar, 9-35 main inventory)
	items   map[int]*Item
	mu      sync.RWMutex

	// Protocol client for sending packets
	client *protocol.Client
}

// Item represents an item in an inventory slot.
type Item struct {
	// Item ID
	ID int32

	// Item name/identifier
	Name string

	// Display name
	DisplayName string

	// Count/stack size
	Count byte

	// Slot number
	Slot int

	// NBT data (raw)
	NBT []byte
}

// NewPlugin creates a new inventory plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("inventory", "1.0.0")
	return &Plugin{
		BasePlugin: base,
		items:      make(map[int]*Item),
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Get protocol client
	p.client = b.Client()

	// Subscribe to inventory packets
	b.On("packet", p.handlePacket)

	return nil
}

// handlePacket handles incoming packets and dispatches to inventory handlers.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet, ok := data[0].(*protocol.Packet)
	if !ok {
		return
	}

	// Only handle play state packets
	if packet.State.String() != "Play" {
		return
	}

	// Set Slot (0x15 in 1.20.1) - handles both window and inventory updates
	if packet.ID == 0x15 {
		p.handleSetSlot(packet)
	}
}

// handleSetSlot handles the Set Slot packet.
func (p *Plugin) handleSetSlot(packet *protocol.Packet) {
	reader := codec.NewReader(bytes.NewReader(packet.Data))

	// Read window ID (0 = player inventory)
	windowID, _ := reader.ReadByte()

	// Only process player inventory updates
	if windowID != 0 {
		return
	}

	// Read slot number
	slot, _ := reader.ReadVarInt()

	// Read item data
	item := p.parseSlot(reader)

	p.mu.Lock()
	if item != nil {
		item.Slot = int(slot)
		p.items[int(slot)] = item
	} else {
		// Item removed from slot
		delete(p.items, int(slot))
	}
	p.mu.Unlock()

	p.Emit("slotChanged", int(slot), item)
	p.Bot().Emit("slotChanged", int(slot), item)
}

// parseSlot parses a single item slot from reader.
func (p *Plugin) parseSlot(reader *codec.Reader) *Item {
	// Read present flag (boolean)
	present, err := reader.ReadBool()
	if err != nil || !present {
		return nil
	}

	item := &Item{}

	// Read item ID (varint)
	itemID, err := reader.ReadVarInt()
	if err != nil {
		return nil
	}
	item.ID = itemID

	// Read item count
	count, err := reader.ReadByte()
	if err != nil {
		return nil
	}
	item.Count = count

	// TODO: Read NBT data for item name/enchantments
	// For now, just set a basic name based on ID
	item.Name = fmt.Sprintf("item_%d", itemID)

	return item
}

// Items returns all items in the player inventory.
func (p *Plugin) Items() []*Item {
	p.mu.RLock()
	defer p.mu.RUnlock()

	items := make([]*Item, 0, len(p.items))
	for _, item := range p.items {
		items = append(items, item)
	}
	return items
}

// FindItem finds an item by name in the inventory.
func (p *Plugin) FindItem(itemName string) *Item {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, item := range p.items {
		if item.Name == itemName || contains(item.Name, itemName) {
			return item
		}
	}
	return nil
}

// CountItem returns the total count of items matching the name.
func (p *Plugin) CountItem(itemName string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, item := range p.items {
		if item.Name == itemName || contains(item.Name, itemName) {
			count += int(item.Count)
		}
	}
	return count
}

// Slot returns the item in a specific slot.
func (p *Plugin) Slot(slot int) *Item {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.items[slot]
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && indexOf(s, substr) >= 0))
}

// indexOf returns the index of a substring or -1.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
