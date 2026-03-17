package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// ExperiencePlugin tracks XP and levels.
type ExperiencePlugin struct {
	mu          sync.RWMutex
	ctx         *plugins.Context
	log         *logger.Logger
	level       int
	experience  float32
	totalXP     int
}

// Metadata returns the plugin metadata.
func (p *ExperiencePlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "experience",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "XP tracking",
		Dependencies: []string{},
	}
}

// OnLoad initializes the experience plugin.
func (p *ExperiencePlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.level = 0
	p.experience = 0
	p.totalXP = 0

	p.log.Info("Experience plugin loaded")

	// Register event handlers
	p.ctx.On("experience_change", p.handleExperienceChange)

	return nil
}

// OnUnload cleans up the experience plugin.
func (p *ExperiencePlugin) OnUnload() error {
	p.log.Info("Experience plugin unloaded")
	return nil
}

// GetLevel returns the current level.
func (p *ExperiencePlugin) GetLevel() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.level
}

// GetExperience returns the current progress to next level (0-1).
func (p *ExperiencePlugin) GetExperience() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.experience
}

// GetTotalXP returns the total XP earned.
func (p *ExperiencePlugin) GetTotalXP() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.totalXP
}

// GetXPForLevel returns the XP needed for a level.
func (p *ExperiencePlugin) GetXPForLevel(level int) int {
	if level <= 16 {
		return level * level + 6 * level
	} else if level <= 31 {
		return level*level*2 - 39*level + 288
	}
	return level*level*4 - 162*level + 2220
}

// GetXPToNextLevel returns XP needed to reach next level.
func (p *ExperiencePlugin) GetXPToNextLevel() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.GetXPForLevel(p.level + 1)
}

// GetXPUntilNextLevel returns progress until next level.
func (p *ExperiencePlugin) GetXPUntilNextLevel() int {
	totalForLevel := p.GetXPForLevel(p.level)
	totalForNext := p.GetXPForLevel(p.level + 1)
	currentTotal := totalForLevel + int(float32(totalForNext-totalForLevel)*p.experience)
	return totalForNext - currentTotal
}

// Event handlers

func (p *ExperiencePlugin) handleExperienceChange(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(args) > 0 {
		if experience, ok := args[0].(float32); ok {
			p.experience = experience
			p.log.Debug("Experience bar updated: %.2f", p.experience)
		}
	}

	if len(args) > 1 {
		if level, ok := args[1].(int); ok {
			if level != p.level {
				oldLevel := p.level
				p.level = level
				p.log.Info("Level up! %d -> %d", oldLevel, p.level)
				p.ctx.Emit("level_up", oldLevel, p.level)
			}
		}
	}

	if len(args) > 2 {
		if totalXP, ok := args[2].(int); ok {
			p.totalXP = totalXP
		}
	}

	p.ctx.Emit("experience_changed", p.experience, p.level, p.totalXP)
}

// String returns a string representation of the experience plugin.
func (p *ExperiencePlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Experience{level=%d, progress=%.2f, total=%d}",
		p.level, p.experience, p.totalXP)
}
