// Package window implements window/inventory handling for Minecraft.
// This plugin handles opening windows, clicking items, and managing inventory.
package window

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/protocol"
	"github.com/nithdevv/goflayer/pkg/protocol/codec"
	"github.com/nithdevv/goflayer/pkg/protocol/states"
)

// Plugin handles window/inventory operations.
type Plugin struct {
	*plugins.BasePlugin

	// Current open window
	currentWindow *Window
	mu            sync.RWMutex

	// Window ID counter
	nextWindowID byte

	// Protocol client for sending packets
	client *protocol.Client
}

// Window represents an open window (inventory, chest, shop, etc.).
type Window struct {
	// Window ID (0 = player inventory)
	ID byte

	// Window type (generic_9x3, anvil, etc.)
	Type string

	// Window title
	Title string

	// Number of slots in the window
	SlotCount int

	// Items in the window
	Items map[int]*Item

	// Window type ID (from protocol)
	TypeID byte
}

// Item represents an item in a window slot.
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

// NewPlugin creates a new window plugin.
func NewPlugin() *Plugin {
	base := plugins.NewBasePlugin("window", "1.0.0")
	return &Plugin{
		BasePlugin:   base,
		currentWindow: nil,
		nextWindowID:  2, // 0 = player inventory, 1 = generic
	}
}

// Load loads the plugin.
func (p *Plugin) Load(b plugins.Bot) error {
	if err := p.BasePlugin.Load(b); err != nil {
		return err
	}

	// Get protocol client
	p.client = b.Client()

	// Subscribe to window packets (1.20.1 Play state)
	b.On("packet", p.handlePacket)

	return nil
}

// handlePacket handles incoming packets and dispatches to window handlers.
func (p *Plugin) handlePacket(data ...interface{}) {
	packet, ok := data[0].(*protocol.Packet)
	if !ok {
		return
	}

	// Only handle play state packets
	if packet.State.String() != "Play" {
		return
	}

	// Window Open (0x2E in 1.20.1)
	if packet.ID == 0x2E {
		p.handleWindowOpen(packet)
		return
	}

	// Window Close (0x13 in 1.20.1)
	if packet.ID == 0x13 {
		p.handleWindowClose(packet)
		return
	}

	// Window Items (0x14 in 1.20.1)
	if packet.ID == 0x14 {
		p.handleWindowItems(packet)
		return
	}

	// Set Slot (0x15 in 1.20.1)
	if packet.ID == 0x15 {
		p.handleSetSlot(packet)
		return
	}
}

// handleWindowOpen handles the Open Window packet.
func (p *Plugin) handleWindowOpen(packet *protocol.Packet) {
	reader := codec.NewReader(bytes.NewReader(packet.Data))

	// Read window ID
	windowID, _ := reader.ReadByte()

	// Read window type
	windowType, _ := reader.ReadVarInt()

	// Read window title (chat)
	// For now, just use a placeholder
	title := "Unknown"

	// Get slot count based on window type
	slotCount := p.getSlotCountForType(int(windowType))

	window := &Window{
		ID:        windowID,
		TypeID:    byte(windowType),
		Type:      p.getWindowTypeName(byte(windowType)),
		Title:     title,
		SlotCount: slotCount,
		Items:     make(map[int]*Item),
	}

	p.mu.Lock()
	p.currentWindow = window
	p.mu.Unlock()

	p.Emit("windowOpen", window)
	p.Bot().Emit("windowOpen", window)

	fmt.Printf("[Window] Opened: %s (ID: %d, Type: %s, Slots: %d)\n",
		title, windowID, window.Type, slotCount)
}

// handleWindowClose handles the Close Window packet (server confirming close).
func (p *Plugin) handleWindowClose(packet *protocol.Packet) {
	reader := codec.NewReader(bytes.NewReader(packet.Data))

	windowID, _ := reader.ReadByte()

	p.mu.Lock()
	closedWindow := p.currentWindow
	if p.currentWindow != nil && p.currentWindow.ID == windowID {
		p.currentWindow = nil
	}
	p.mu.Unlock()

	p.Emit("windowClose", windowID)
	p.Bot().Emit("windowClose", windowID)

	if closedWindow != nil {
		fmt.Printf("[Window] Closed: %s (ID: %d)\n", closedWindow.Title, windowID)
	}
}

// handleWindowItems handles the Window Items packet (initial window contents).
func (p *Plugin) handleWindowItems(packet *protocol.Packet) {
	reader := codec.NewReader(bytes.NewReader(packet.Data))

	// Read window ID
	windowID, _ := reader.ReadByte()

	// Read state (0 = full window update)
	state, _ := reader.ReadVarInt()

	// Read item count
	count, _ := reader.ReadVarInt()

	p.mu.RLock()
	window := p.currentWindow
	p.mu.RUnlock()

	if window == nil || window.ID != windowID {
		return
	}

	// Clear existing items if state is 0
	if state == 0 {
		window.Items = make(map[int]*Item)
	}

	// Parse items
	for i := 0; i < int(count); i++ {
		item := p.parseSlot(reader)
		if item != nil {
			item.Slot = i
			window.Items[i] = item
		}
	}

	p.Emit("windowItems", window)
	p.Bot().Emit("windowItems", window)

	fmt.Printf("[Window] Updated %d items in window %d\n", count, windowID)
}

// handleSetSlot handles the Set Slot packet (single item update).
func (p *Plugin) handleSetSlot(packet *protocol.Packet) {
	reader := codec.NewReader(bytes.NewReader(packet.Data))

	// Read window ID
	windowID, _ := reader.ReadByte()

	// Read slot number
	slot, _ := reader.ReadVarInt()

	// Read item data
	item := p.parseSlot(reader)

	p.mu.RLock()
	window := p.currentWindow
	p.mu.RUnlock()

	if windowID == 0 || (window != nil && window.ID == windowID) {
		if item != nil {
			item.Slot = int(slot)
			if window != nil {
				window.Items[int(slot)] = item
			}
			p.Emit("setSlot", windowID, slot, item)
		} else {
			// Item removed
			if window != nil {
				delete(window.Items, int(slot))
			}
			p.Emit("setSlot", windowID, slot, nil)
		}
	}
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

// CurrentWindow returns the currently open window.
func (p *Plugin) CurrentWindow() *Window {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.currentWindow
}

// ClickWindow clicks a slot in the current window.
// mode: 0 = click, 1 = shift click, 2 = number key, etc.
// button: 0 = left click, 1 = right click, etc.
func (p *Plugin) ClickWindow(slot int, mode byte, button byte) error {
	p.mu.RLock()
	window := p.currentWindow
	p.mu.RUnlock()

	if window == nil {
		return fmt.Errorf("no window open")
	}

	// Create window click packet
	// Packet ID: 0x0E (Window Click) for 1.20.1
	buf := &bytes.Buffer{}
	writer := codec.NewWriter(buf)

	// Window ID
	writer.WriteByte(window.ID)

	// Slot (varint)
	writer.WriteVarInt(int32(slot))

	// Button
	writer.WriteByte(button)

	// Mode
	writer.WriteByte(mode)

	// Variable for each click type (usually 1 for left click, 0 for right click)
	varintValue := int32(1)
	if button == 1 {
		varintValue = 0
	}
	writer.WriteVarInt(varintValue)

	packet := protocol.NewPacket(0x0E, states.Play, states.Serverbound)
	packet.Data = buf.Bytes()

	// Send packet
	if err := p.client.Write(packet); err != nil {
		return fmt.Errorf("failed to send click packet: %w", err)
	}

	fmt.Printf("[Window] Clicked slot %d (mode: %d, button: %d)\n", slot, mode, button)
	return nil
}

// CloseWindow closes the current window.
func (p *Plugin) CloseWindow(window *Window) error {
	if window == nil {
		p.mu.RLock()
		window = p.currentWindow
		p.mu.RUnlock()
	}

	if window == nil {
		return fmt.Errorf("no window to close")
	}

	// Create window close packet
	// Packet ID: 0x10 (Close Window) for 1.20.1
	buf := &bytes.Buffer{}
	buf.WriteByte(window.ID)

	packet := protocol.NewPacket(0x10, states.Play, states.Serverbound)
	packet.Data = buf.Bytes()

	p.mu.Lock()
	if p.currentWindow == window {
		p.currentWindow = nil
	}
	p.mu.Unlock()

	// Send packet
	if err := p.client.Write(packet); err != nil {
		return fmt.Errorf("failed to send close packet: %w", err)
	}

	fmt.Printf("[Window] Closed window %d\n", window.ID)
	return nil
}

// ContainerItems returns items in the window container (excluding player inventory).
func (p *Plugin) ContainerItems(window *Window) []*Item {
	if window == nil {
		p.mu.RLock()
		window = p.currentWindow
		p.mu.RUnlock()
	}

	if window == nil {
		return nil
	}

	items := make([]*Item, 0)
	playerInvStart := window.SlotCount - 36

	for _, item := range window.Items {
		// Only include container items, not player inventory
		if item.Slot >= 0 && item.Slot < playerInvStart {
			items = append(items, item)
		}
	}
	return items
}

// FindItem finds an item by name in the window.
func (p *Plugin) FindItem(window *Window, itemName string) *Item {
	if window == nil {
		p.mu.RLock()
		window = p.currentWindow
		p.mu.RUnlock()
	}

	if window == nil {
		return nil
	}

	for _, item := range window.Items {
		if item.Name == itemName || contains(item.Name, itemName) {
			return item
		}
	}
	return nil
}

// getSlotCountForType returns the slot count for a window type.
func (p *Plugin) getSlotCountForType(windowType int) int {
	// Common window types
	switch windowType {
	case 0: // Generic 9x1
		return 9 + 36 // Add player inventory slots
	case 1: // Generic 9x2
		return 18 + 36
	case 2: // Generic 9x3
		return 27 + 36
	case 3: // Generic 9x4
		return 36 + 36
	case 4: // Generic 9x5
		return 45 + 36
	case 5: // Generic 9x6
		return 54 + 36
	case 9: // Generic 3x3
		return 9 + 36
	case 10: // Anvil
		return 3 + 36
	case 11: // Beacon
		return 1 + 36
	case 12: // Blast Furnace
		return 3 + 36
	case 13: // Brewing Stand
		return 4 + 36
	case 14: // Crafting
		return 10 + 36
	case 15: // Enchantment
		return 2 + 36
	case 16: // Furnace
		return 3 + 36
	case 17: // Grindstone
		return 3 + 36
	case 18: // Hopper
		return 5 + 36
	case 19: // Lectern
		return 1 + 36
	case 20: // Loom
		return 4 + 36
	case 21: // Merchant
		return 3 + 36
	case 22: // Shulker Box
		return 27 + 36
	case 23: // Smoker
		return 3 + 36
	case 24: // Cartography
		return 3 + 36
	case 25: // Stonecutter
		return 2 + 36
	default:
		return 54 + 36 // Default to large chest + player inventory
	}
}

// getWindowTypeName returns the name of a window type.
func (p *Plugin) getWindowTypeName(typeID byte) string {
	switch typeID {
	case 0:
		return "minecraft:generic_9x1"
	case 1:
		return "minecraft:generic_9x2"
	case 2:
		return "minecraft:generic_9x3"
	case 3:
		return "minecraft:generic_9x4"
	case 4:
		return "minecraft:generic_9x5"
	case 5:
		return "minecraft:generic_9x6"
	case 9:
		return "minecraft:generic_3x3"
	case 10:
		return "minecraft:anvil"
	case 11:
		return "minecraft:beacon"
	case 12:
		return "minecraft:blast_furnace"
	case 13:
		return "minecraft:brewing_stand"
	case 14:
		return "minecraft:crafting"
	case 15:
		return "minecraft:enchantment"
	case 16:
		return "minecraft:furnace"
	case 17:
		return "minecraft:grindstone"
	case 18:
		return "minecraft:hopper"
	case 19:
		return "minecraft:lectern"
	case 20:
		return "minecraft:loom"
	case 21:
		return "minecraft:merchant"
	case 22:
		return "minecraft:shulker_box"
	case 23:
		return "minecraft:smoker"
	case 24:
		return "minecraft:cartography_table"
	case 25:
		return "minecraft:stonecutter"
	default:
		return fmt.Sprintf("minecraft:unknown_%d", typeID)
	}
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
