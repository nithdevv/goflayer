// Package core provides essential plugins for goflayer.
package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/internal/types"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// GamePlugin tracks game state.
type GamePlugin struct {
	mu     sync.RWMutex
	ctx    *plugins.Context
	log    *logger.Logger
	state  types.GameState
	level  string
	gamemode int
	dimension string
	difficulty int
	hardcore bool
}

// Metadata returns the plugin metadata.
func (p *GamePlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "game",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Game state tracking (game mode, dimension, difficulty, etc.)",
		Dependencies: []string{},
	}
}

// OnLoad initializes the game plugin.
func (p *GamePlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.state = types.StateDisconnected

	p.log.Info("Game plugin loaded")

	// Register event handlers
	p.ctx.On("login", p.handleLogin)
	p.ctx.On("respawn", p.handleRespawn)
	p.ctx.On("game_state_change", p.handleGameStateChange)
	p.ctx.On("difficulty", p.handleDifficulty)
	p.ctx.On("disconnect", p.handleDisconnect)

	return nil
}

// OnUnload cleans up the game plugin.
func (p *GamePlugin) OnUnload() error {
	p.log.Info("Game plugin unloaded")
	return nil
}

// GetState returns the current game state.
func (p *GamePlugin) GetState() types.GameState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetLevel returns the level name.
func (p *GamePlugin) GetLevel() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.level
}

// GetGamemode returns the player's game mode.
// 0 = Survival, 1 = Creative, 2 = Adventure, 3 = Spectator
func (p *GamePlugin) GetGamemode() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.gamemode
}

// GetDimension returns the current dimension.
// Common values: "minecraft:overworld", "minecraft:the_nether", "minecraft:the_end"
func (p *GamePlugin) GetDimension() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dimension
}

// GetDifficulty returns the difficulty.
// 0 = Peaceful, 1 = Easy, 2 = Normal, 3 = Hard
func (p *GamePlugin) GetDifficulty() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.difficulty
}

// IsHardcore returns true if hardcore mode is enabled.
func (p *GamePlugin) IsHardcore() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.hardcore
}

// IsInPlay returns true if in play state.
func (p *GamePlugin) IsInPlay() bool {
	return p.GetState() == types.StatePlay
}

// IsSurvival returns true if in survival mode.
func (p *GamePlugin) IsSurvival() bool {
	return p.GetGamemode() == 0
}

// IsCreative returns true if in creative mode.
func (p *GamePlugin) IsCreative() bool {
	return p.GetGamemode() == 1
}

// IsAdventure returns true if in adventure mode.
func (p *GamePlugin) IsAdventure() bool {
	return p.GetGamemode() == 2
}

// IsSpectator returns true if in spectator mode.
func (p *GamePlugin) IsSpectator() bool {
	return p.GetGamemode() == 3
}

// Event handlers

func (p *GamePlugin) handleLogin(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.state = types.StatePlay
	p.log.Info("Entered game state")
	p.ctx.Emit("game_entered")
}

func (p *GamePlugin) handleRespawn(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Update dimension from respawn packet
	if len(args) > 0 {
		if dimension, ok := args[0].(string); ok {
			p.dimension = dimension
			p.log.Info("Respawned in dimension: %s", dimension)
		}
	}

	p.ctx.Emit("respawned")
}

func (p *GamePlugin) handleGameStateChange(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	reason, ok := args[0].(int)
	if !ok {
		return
	}

	p.mu.Lock()

	switch reason {
	case 3: // Game mode change
		if len(args) > 1 {
			if gamemode, ok := args[1].(int); ok {
				p.gamemode = gamemode
				p.log.Info("Game mode changed to: %d", gamemode)
				p.ctx.Emit("gamemode_changed", gamemode)
			}
		}
	}

	p.mu.Unlock()
}

func (p *GamePlugin) handleDifficulty(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(args) > 0 {
		if difficulty, ok := args[0].(int); ok {
			p.difficulty = difficulty
			p.log.Debug("Difficulty: %d", difficulty)
		}
	}

	if len(args) > 1 {
		if hardcore, ok := args[1].(bool); ok {
			p.hardcore = hardcore
			p.log.Debug("Hardcore: %v", hardcore)
		}
	}
}

func (p *GamePlugin) handleDisconnect(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.state = types.StateDisconnected
	p.log.Info("Disconnected from game")
	p.ctx.Emit("game_exited")
}

// String returns a string representation of the game state.
func (p *GamePlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	gamemodeName := "unknown"
	switch p.gamemode {
	case 0:
		gamemodeName = "survival"
	case 1:
		gamemodeName = "creative"
	case 2:
		gamemodeName = "adventure"
	case 3:
		gamemodeName = "spectator"
	}

	return fmt.Sprintf("Game{state=%s, gamemode=%s, dimension=%s, difficulty=%d, hardcore=%v}",
		p.state, gamemodeName, p.dimension, p.difficulty, p.hardcore)
}
