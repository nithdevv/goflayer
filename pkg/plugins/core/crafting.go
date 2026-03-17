package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/math"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// Recipe represents a crafting recipe.
type Recipe struct {
	RecipeID int32
	Type     string
	Width    int
	Height   int
	Ingredients [][]int32
	Result *Item
}

// CraftingPlugin handles crafting table interactions.
type CraftingPlugin struct {
	mu             sync.RWMutex
	ctx            *plugins.Context
	log            *logger.Logger
	openTablePos   *math.BlockPos
	isOpen         bool
	recipes        map[int32]*Recipe
}

// Metadata returns the plugin metadata.
func (p *CraftingPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "crafting",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Crafting table interaction",
		Dependencies: []string{"inventory"},
	}
}

// OnLoad initializes the crafting plugin.
func (p *CraftingPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.isOpen = false
	p.recipes = make(map[int32]*Recipe)

	p.log.Info("Crafting plugin loaded")

	// Register event handlers
	p.ctx.On("crafting_recipe", p.handleCraftingRecipe)

	return nil
}

// OnUnload cleans up the crafting plugin.
func (p *CraftingPlugin) OnUnload() error {
	p.log.Info("Crafting plugin unloaded")
	return nil
}

// Open opens a crafting table at a position.
func (p *CraftingPlugin) Open(pos *math.BlockPos) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isOpen {
		return fmt.Errorf("crafting table already open")
	}

	p.log.Info("Opening crafting table at %v", pos)
	p.openTablePos = pos
	p.isOpen = true

	// TODO: Send open crafting table packet
	p.ctx.Emit("crafting_table_opening", pos)

	return nil
}

// Close closes the currently open crafting table.
func (p *CraftingPlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isOpen {
		return fmt.Errorf("no crafting table open")
	}

	p.log.Info("Closing crafting table")
	p.isOpen = false

	// TODO: Send close crafting table packet
	p.ctx.Emit("crafting_table_closing")

	return nil
}

// Craft crafts an item by recipe ID.
func (p *CraftingPlugin) Craft(recipeID int32, count int) error {
	p.log.Info("Crafting recipe %d, count %d", recipeID, count)

	// TODO: Implement crafting logic
	p.ctx.Emit("item_crafted", recipeID, count)

	return nil
}

// CraftItems crafts items from ingredients.
func (p *CraftingPlugin) CraftItems(ingredients [][]int32, result *Item, count int) error {
	p.log.Info("Crafting item %s, count %d", result.DisplayName, count)

	// TODO: Implement crafting logic
	p.ctx.Emit("item_crafted_custom", ingredients, result, count)

	return nil
}

// GetRecipes returns all known recipes.
func (p *CraftingPlugin) GetRecipes() []*Recipe {
	p.mu.RLock()
	defer p.mu.RUnlock()

	recipes := make([]*Recipe, 0, len(p.recipes))
	for _, recipe := range p.recipes {
		recipes = append(recipes, recipe)
	}
	return recipes
}

// FindRecipe finds a recipe by result item type.
func (p *CraftingPlugin) FindRecipe(itemType int32) *Recipe {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, recipe := range p.recipes {
		if recipe.Result.TypeID == itemType {
			return recipe
		}
	}
	return nil
}

// CanCraft checks if we have enough ingredients to craft a recipe.
func (p *CraftingPlugin) CanCraft(recipe *Recipe, count int) bool {
	// TODO: Check inventory for ingredients
	return true
}

// Event handlers

func (p *CraftingPlugin) handleCraftingRecipe(args ...interface{}) {
	p.log.Debug("Crafting recipe received")
	// TODO: Parse recipe packet
}

// String returns a string representation of the crafting plugin.
func (p *CraftingPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Crafting{open=%v, recipes=%d}", p.isOpen, len(p.recipes))
}
