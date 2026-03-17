// Package world implements Minecraft chunk management and NBT parsing.
package world

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/nithdevv/goflayer/pkg/nbt"
	"github.com/nithdevv/goflayer/pkg/registry"
)

// ChunkPos represents a chunk position in the world.
type ChunkPos struct {
	X int32
	Z int32
}

// Chunk represents a 16x16xY chunk in the world.
type Chunk struct {
	mu       sync.RWMutex
	pos      ChunkPos
	sections []*ChunkSection
	biomes   []int32 // Biome IDs for the chunk
	blockEntities map[[3]int]interface{} // x, y, z -> entity data
	loaded   bool
	dirty    bool
	lightPopulated bool
	terrainPopulated bool
}

// ChunkSection represents a 16x16x16 section of a chunk.
type ChunkSection struct {
	y               int8
	blocks          []BlockState
	blocksCount     int
	palette         *Palette
	biomePalette    *Palette
	skyLight        []byte
	blockLight      []byte
}

// BlockState represents a block state with its ID and data.
type BlockState struct {
	ID    registry.BlockID
	Data  int32 // Metadata/state
}

// Palette represents a block palette for chunk compression.
type Palette struct {
	bitsPerEntry uint8
	entries      []BlockState
	paletteToId  map[BlockState]int32
	idToPalette  map[int32]BlockState
}

// NewChunk creates a new chunk at the given position.
func NewChunk(x, z int32) *Chunk {
	chunk := &Chunk{
		pos:           ChunkPos{X: x, Z: z},
		sections:      make([]*ChunkSection, 0, 24), // For Y=-64 to 320
		biomes:        make([]int32, 1024), // 4x4x4 biome grid
		blockEntities: make(map[[3]int]interface{}),
		loaded:        true,
		dirty:         false,
		lightPopulated: false,
		terrainPopulated: false,
	}

	// Initialize sections for Y=-64 to 320 (24 sections: -4 to 19)
	for i := -4; i < 20; i++ {
		chunk.sections = append(chunk.sections, NewChunkSection(int8(i)))
	}

	return chunk
}

// NewChunkSection creates a new chunk section.
func NewChunkSection(y int8) *ChunkSection {
	section := &ChunkSection{
		y:          y,
		blocks:     make([]BlockState, 4096), // 16*16*16
		blocksCount: 0,
		skyLight:   make([]byte, 2048), // 16*16*16/2 (nibble array)
		blockLight: make([]byte, 2048),
	}

	// Initialize with air
	for i := range section.blocks {
		section.blocks[i] = BlockState{ID: registry.BlockAir}
	}

	// Initialize palette with air
	section.palette = NewPalette()
	section.palette.Add(BlockState{ID: registry.BlockAir})

	return section
}

// NewPalette creates a new palette.
func NewPalette() *Palette {
	return &Palette{
		bitsPerEntry: 4,
		entries:      make([]BlockState, 0, 16),
		paletteToId:  make(map[BlockState]int32),
		idToPalette:  make(map[int32]BlockState),
	}
}

// Add adds a block state to the palette.
func (p *Palette) Add(state BlockState) int32 {
	if id, ok := p.paletteToId[state]; ok {
		return id
	}

	id := int32(len(p.entries))
	p.entries = append(p.entries, state)
	p.paletteToId[state] = id
	p.idToPalette[id] = state

	// Update bits per entry if needed
	if len(p.entries) > 1<<p.bitsPerEntry {
		p.bitsPerEntry = calculateBitsPerEntry(len(p.entries))
	}

	return id
}

// Get retrieves a block state from the palette by ID.
func (p *Palette) Get(id int32) (BlockState, bool) {
	state, ok := p.idToPalette[id]
	return state, ok
}

// Size returns the number of entries in the palette.
func (p *Palette) Size() int {
	return len(p.entries)
}

// calculateBitsPerEntry calculates the required bits per entry.
func calculateBitsPerEntry(size int) uint8 {
	switch {
	case size <= 2:
		return 1
	case size <= 4:
		return 2
	case size <= 8:
		return 3
	case size <= 16:
		return 4
	case size <= 32:
		return 5
	case size <= 64:
		return 6
	case size <= 128:
		return 7
	case size <= 256:
		return 8
	case size <= 512:
		return 9
	case size <= 1024:
		return 10
	case size <= 2048:
		return 11
	case size <= 4096:
		return 12
	case size <= 8192:
		return 13
	case size <= 16384:
		return 14
	case size <= 32768:
		return 15
	default:
		return 16
	}
}

// Position returns the chunk position.
func (c *Chunk) Position() ChunkPos {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pos
}

// IsLoaded returns whether the chunk is loaded.
func (c *Chunk) IsLoaded() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.loaded
}

// IsDirty returns whether the chunk has been modified.
func (c *Chunk) IsDirty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.dirty
}

// MarkDirty marks the chunk as dirty.
func (c *Chunk) MarkDirty() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dirty = true
}

// MarkClean marks the chunk as clean.
func (c *Chunk) MarkClean() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dirty = false
}

// GetSection returns a chunk section by Y coordinate (section index).
func (c *Chunk) GetSection(y int8) (*ChunkSection, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert Y to section index (-64 to 320 -> -4 to 19)
	sectionIndex := int(y+4) // Offset by -4 (minimum section)

	if sectionIndex < 0 || sectionIndex >= len(c.sections) {
		return nil, fmt.Errorf("section index out of bounds: %d (y=%d)", sectionIndex, y)
	}

	return c.sections[sectionIndex], nil
}

// GetBlock returns a block at local coordinates.
func (c *Chunk) GetBlock(x, y, z int32) (registry.BlockID, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if x < 0 || x >= 16 || y < -64 || y >= 320 || z < 0 || z >= 16 {
		return registry.BlockAir, fmt.Errorf("block coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	sectionY := int8(y / 16)
	section, err := c.GetSection(sectionY)
	if err != nil {
		return registry.BlockAir, err
	}

	return section.GetBlock(int(x), int(y%16), int(z))
}

// SetBlock sets a block at local coordinates.
func (c *Chunk) SetBlock(x, y, z int32, blockID registry.BlockID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if x < 0 || x >= 16 || y < -64 || y >= 320 || z < 0 || z >= 16 {
		return fmt.Errorf("block coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	sectionY := int8(y / 16)
	section, err := c.GetSection(sectionY)
	if err != nil {
		return err
	}

	err = section.SetBlock(int(x), int(y%16), int(z), BlockState{ID: blockID})
	if err != nil {
		return err
	}

	c.dirty = true
	return nil
}

// GetBlockEntity returns a block entity at the given position.
func (c *Chunk) GetBlockEntity(x, y, z int32) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := [3]int{int(x), int(y), int(z)}
	entity, ok := c.blockEntities[key]
	return entity, ok
}

// SetBlockEntity sets a block entity at the given position.
func (c *Chunk) SetBlockEntity(x, y, z int32, entity interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := [3]int{int(x), int(y), int(z)}
	c.blockEntities[key] = entity
	c.dirty = true
}

// RemoveBlockEntity removes a block entity at the given position.
func (c *Chunk) RemoveBlockEntity(x, y, z int32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := [3]int{int(x), int(y), int(z)}
	delete(c.blockEntities, key)
	c.dirty = true
}

// GetBiome returns the biome at the given position.
func (c *Chunk) GetBiome(x, y, z int32) (int32, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if x < 0 || x >= 16 || y < -64 || y >= 320 || z < 0 || z >= 16 {
		return 0, fmt.Errorf("biome coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	// Biomes are stored in a 4x4x4 grid
	bx, by, bz := x/4, y/4, z/4
	index := by*16 + bz*4 + bx

	if int(index) >= len(c.biomes) {
		return 0, fmt.Errorf("biome index out of bounds: %d", index)
	}

	return c.biomes[index], nil
}

// SetBiome sets the biome at the given position.
func (c *Chunk) SetBiome(x, y, z int32, biomeID int32) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if x < 0 || x >= 16 || y < -64 || y >= 320 || z < 0 || z >= 16 {
		return fmt.Errorf("biome coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	// Biomes are stored in a 4x4x4 grid
	bx, by, bz := x/4, y/4, z/4
	index := by*16 + bz*4 + bx

	if int(index) >= len(c.biomes) {
		return fmt.Errorf("biome index out of bounds: %d", index)
	}

	c.biomes[index] = biomeID
	c.dirty = true
	return nil
}

// GetBlock returns a block at local section coordinates.
func (s *ChunkSection) GetBlock(x, y, z int) (registry.BlockID, error) {
	if x < 0 || x >= 16 || y < 0 || y >= 16 || z < 0 || z >= 16 {
		return registry.BlockAir, fmt.Errorf("block coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	index := y*256 + z*16 + x
	return s.blocks[index].ID, nil
}

// SetBlock sets a block at local section coordinates.
func (s *ChunkSection) SetBlock(x, y, z int, state BlockState) error {
	if x < 0 || x >= 16 || y < 0 || y >= 16 || z < 0 || z >= 16 {
		return fmt.Errorf("block coordinates out of bounds: (%d, %d, %d)", x, y, z)
	}

	index := y*256 + z*16 + x
	wasAir := s.blocks[index].ID == registry.BlockAir
	s.blocks[index] = state

	if state.ID == registry.BlockAir && !wasAir {
		s.blocksCount--
	} else if state.ID != registry.BlockAir && wasAir {
		s.blocksCount++
	}

	// Update palette
	s.palette.Add(state)

	return nil
}

// IsEmpty returns whether the section is empty (all air).
func (s *ChunkSection) IsEmpty() bool {
	return s.blocksCount == 0
}

// GetY returns the Y coordinate of this section.
func (s *ChunkSection) GetY() int8 {
	return s.y
}

// ParseChunkFromNBT parses chunk data from NBT format.
func ParseChunkFromNBT(data []byte) (*Chunk, error) {
	// Decompress the data
	reader := bytes.NewReader(data)

	// Try gzip first
	gzReader, err := gzip.NewReader(reader)
	if err == nil {
		defer gzReader.Close()
		data, err = io.ReadAll(gzReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress gzip: %w", err)
		}
	} else {
		// Try zlib
		reader.Seek(0, io.SeekStart)
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer zlibReader.Close()

		data, err = io.ReadAll(zlibReader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress zlib: %w", err)
		}
	}

	// Parse NBT
	tag, err := nbt.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NBT: %w", err)
	}

	return parseChunkFromTag(tag)
}

// parseChunkFromTag parses a chunk from an NBT tag.
func parseChunkFromTag(tag nbt.Tag) (*Chunk, error) {
	compound, ok := tag.Value.(map[string]interface{})
	if !ok {
		return nil, errors.New("root tag is not a compound")
	}

	// Extract chunk position
	var chunkX, chunkZ int32
	if xVal, ok := compound["xPos"]; ok {
		chunkX = int32(xVal.(int32))
	}
	if zVal, ok := compound["zPos"]; ok {
		chunkZ = int32(zVal.(int32))
	}

	chunk := NewChunk(chunkX, chunkZ)

	// Parse sections
	if sectionsVal, ok := compound["sections"]; ok {
		sectionsList, ok := sectionsVal.([]interface{})
		if !ok {
			return nil, errors.New("sections is not a list")
		}

		for _, sectionVal := range sectionsList {
			sectionCompound, ok := sectionVal.(map[string]interface{})
			if !ok {
				continue
			}

			yVal, ok := sectionCompound["Y"].(int8)
			if !ok {
				continue
			}

			section, err := parseSectionFromCompound(sectionCompound)
			if err != nil {
				continue // Skip invalid sections
			}

			sectionIndex := int(yVal + 4)
			if sectionIndex >= 0 && sectionIndex < len(chunk.sections) {
				chunk.sections[sectionIndex] = section
			}
		}
	}

	// Parse biomes
	if biomesVal, ok := compound["Biomes"]; ok {
		if biomesData, ok := biomesVal.(map[string]interface{}); ok {
			if dataVal, ok := biomesData["data"].([]int32); ok {
				if len(dataVal) <= len(chunk.biomes) {
					copy(chunk.biomes, dataVal)
				}
			}
		}
	}

	// Parse block entities
	if entitiesVal, ok := compound["block_entities"]; ok {
		if entitiesList, ok := entitiesVal.([]interface{}); ok {
			for _, entityVal := range entitiesList {
				if entityCompound, ok := entityVal.(map[string]interface{}); ok {
					if xVal, ok := entityCompound["x"].(int32); ok {
						if yVal, ok := entityCompound["y"].(int32); ok {
							if zVal, ok := entityCompound["z"].(int32); ok {
								key := [3]int{int(xVal & 15), int(yVal), int(zVal & 15)}
								chunk.blockEntities[key] = entityCompound
							}
						}
					}
				}
			}
		}
	}

	return chunk, nil
}

// parseSectionFromCompound parses a chunk section from an NBT compound.
func parseSectionFromCompound(compound map[string]interface{}) (*ChunkSection, error) {
	yVal, ok := compound["Y"].(int8)
	if !ok {
		return nil, errors.New("missing Y coordinate")
	}

	section := NewChunkSection(yVal)

	// Parse block states
	if statesVal, ok := compound["block_states"]; ok {
		if statesCompound, ok := statesVal.(map[string]interface{}); ok {
			err := parseBlockStates(section, statesCompound)
			if err != nil {
				return nil, fmt.Errorf("failed to parse block states: %w", err)
			}
		}
	}

	// Parse biomes
	if biomesVal, ok := compound["biomes"]; ok {
		if biomesCompound, ok := biomesVal.(map[string]interface{}); ok {
			parseBiomes(section, biomesCompound)
		}
	}

	return section, nil
}

// parseBlockStates parses block states from NBT.
func parseBlockStates(section *ChunkSection, compound map[string]interface{}) error {
	// Get palette
	if paletteVal, ok := compound["palette"]; ok {
		if paletteList, ok := paletteVal.([]interface{}); ok {
			for _, paletteEntry := range paletteList {
				if entryCompound, ok := paletteEntry.(map[string]interface{}); ok {
					if nameVal, ok := entryCompound["Name"].(string); ok {
						blockID, ok := registry.GetBlockByName(nameVal)
						if !ok {
							blockID = registry.BlockAir
						}

						state := BlockState{ID: blockID}

						// Parse properties if present
						if propsVal, ok := entryCompound["Properties"]; ok {
							if propsCompound, ok := propsVal.(map[string]interface{}); ok {
								// Properties could affect the state ID
								_ = propsCompound
							}
						}

						section.palette.Add(state)
					}
				}
			}
		}
	}

	// Get block data
	if dataVal, ok := compound["data"]; ok {
		if dataLongArray, ok := dataVal.([]int64); ok {
			err := parseBlockData(section, dataLongArray)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// parseBlockData parses block data from a long array.
func parseBlockData(section *ChunkSection, data []int64) error {
	if len(data) == 0 {
		return nil
	}

	bitsPerEntry := section.palette.bitsPerEntry
	if bitsPerEntry == 0 {
		bitsPerEntry = 4
	}

	// Parse packed long array
	mask := int64(1<<bitsPerEntry) - 1
	entriesPerLong := uint8(64 / bitsPerEntry)

	blockIndex := 0
	for _, val := range data {
		for i := 0; i < int(entriesPerLong) && blockIndex < 4096; i++ {
			paletteID := int32((val >> (i * int(bitsPerEntry))) & mask)

			if state, ok := section.palette.Get(paletteID); ok {
				x := blockIndex % 16
				y := (blockIndex / 256) % 16
				z := (blockIndex / 16) % 16

				section.blocks[y*256+z*16+x] = state
			}

			blockIndex++
		}
	}

	return nil
}

// parseBiomes parses biome data from NBT.
func parseBiomes(section *ChunkSection, compound map[string]interface{}) {
	// Biome parsing implementation
	// For now, we'll use a default biome
}

// SerializeToNBT serializes the chunk to NBT format.
func (c *Chunk) SerializeToNBT(compressed bool) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	compound := make(map[string]interface{})

	// Add position
	compound["xPos"] = int32(c.pos.X)
	compound["zPos"] = int32(c.pos.Z)
	compound["yPos"] = int8(-4) // Minimum section Y

	// Add sections
	sectionsList := make([]interface{}, 0)
	for _, section := range c.sections {
		if section != nil && !section.IsEmpty() {
			sectionCompound, err := serializeSection(section)
			if err == nil {
				sectionsList = append(sectionsList, sectionCompound)
			}
		}
	}
	compound["sections"] = sectionsList

	// Add biomes
	biomesCompound := make(map[string]interface{})
	biomesCompound["data"] = c.biomes
	compound["Biomes"] = biomesCompound

	// Add block entities
	entitiesList := make([]interface{}, 0)
	for key, entity := range c.blockEntities {
		if entityCompound, ok := entity.(map[string]interface{}); ok {
			// Update coordinates
			entityCompound["x"] = int32(key[0])
			entityCompound["y"] = int32(key[1])
			entityCompound["z"] = int32(key[2])
			entitiesList = append(entitiesList, entityCompound)
		}
	}
	compound["block_entities"] = entitiesList

	// Create NBT tag
	tag := nbt.NewTag("", nbt.TagCompound, compound)

	// Serialize to bytes
	data, err := nbt.Marshal(tag)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NBT: %w", err)
	}

	if !compressed {
		return data, nil
	}

	// Compress with gzip
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, err = gzWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}
	err = gzWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return buf.Bytes(), nil
}

// serializeSection serializes a chunk section to NBT.
func serializeSection(section *ChunkSection) (map[string]interface{}, error) {
	compound := make(map[string]interface{})

	compound["Y"] = section.y

	// Serialize block states
	blockStates := make(map[string]interface{})

	// Create palette
	paletteList := make([]interface{}, 0)
	for _, state := range section.palette.entries {
		entry := make(map[string]interface{})
		if name, ok := registry.GetBlockName(state.ID); ok {
			entry["Name"] = name
		} else {
			entry["Name"] = "minecraft:air"
		}
		paletteList = append(paletteList, entry)
	}
	blockStates["palette"] = paletteList

	// Serialize block data if needed
	if section.palette.Size() > 1 {
		data, err := serializeBlockData(section)
		if err == nil {
			blockStates["data"] = data
		}
	}

	compound["block_states"] = blockStates

	return compound, nil
}

// serializeBlockData serializes block data to a long array.
func serializeBlockData(section *ChunkSection) ([]int64, error) {
	bitsPerEntry := section.palette.bitsPerEntry
	if bitsPerEntry == 0 {
		bitsPerEntry = 4
	}

	entriesPerLong := int(64 / bitsPerEntry)
	longsNeeded := (4096 + entriesPerLong - 1) / entriesPerLong
	data := make([]int64, longsNeeded)

	blockIndex := 0
	for i := range data {
		var val int64
		for j := 0; j < int(entriesPerLong) && blockIndex < 4096; j++ {
			x := blockIndex % 16
			y := (blockIndex / 256) % 16
			z := (blockIndex / 16) % 16

			state := section.blocks[y*256+z*16+x]
			paletteID := section.palette.paletteToId[state]

			val |= int64(paletteID) << (j * int(bitsPerEntry))
			blockIndex++
		}
		data[i] = val
	}

	return data, nil
}

// ReadVarInt reads a VarInt from the reader.
func ReadVarInt(r io.Reader) (int32, error) {
	var result int32
	var shift uint

	for {
		var b [1]byte
		_, err := io.ReadFull(r, b[:])
		if err != nil {
			return 0, err
		}

		value := int32(b[0] & 0x7F)
		result |= value << shift

		if (b[0] & 0x80) == 0 {
			return result, nil
		}

		shift += 7
		if shift >= 35 {
			return 0, errors.New("varint too big")
		}
	}
}

// WriteVarInt writes a VarInt to the writer.
func WriteVarInt(w io.Writer, value int32) error {
	for {
		temp := value & 0x7F
		value >>= 7

		if value != 0 {
			temp |= 0x80
		}

		err := binary.Write(w, binary.BigEndian, int8(temp))
		if err != nil {
			return err
		}

		if value == 0 {
			return nil
		}
	}
}

// ReadString reads a length-prefixed string from the reader.
func ReadString(r io.Reader) (string, error) {
	length, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	if length < 0 {
		return "", errors.New("string length is negative")
	}

	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// WriteString writes a length-prefixed string to the writer.
func WriteString(w io.Writer, value string) error {
	data := []byte(value)
	err := WriteVarInt(w, int32(len(data)))
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// ReadUByte reads an unsigned byte from the reader.
func ReadUByte(r io.Reader) (uint8, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// WriteUByte writes an unsigned byte to the writer.
func WriteUByte(w io.Writer, value uint8) error {
	return binary.Write(w, binary.BigEndian, value)
}

// ReadBool reads a boolean from the reader.
func ReadBool(r io.Reader) (bool, error) {
	b, err := ReadUByte(r)
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

// WriteBool writes a boolean to the writer.
func WriteBool(w io.Writer, value bool) error {
	var b uint8
	if value {
		b = 1
	}
	return WriteUByte(w, b)
}
