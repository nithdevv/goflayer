package core

import (
	"fmt"
	"sync"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/pkg/plugins"
)

// HealthPlugin tracks health and food.
type HealthPlugin struct {
	mu           sync.RWMutex
	ctx          *plugins.Context
	log          *logger.Logger
	health       float32
	maxHealth    float32
	food         int
	saturation   float32
	exhaustion   float32
}

// Metadata returns the plugin metadata.
func (p *HealthPlugin) Metadata() plugins.PluginMetadata {
	return plugins.PluginMetadata{
		Name:        "health",
		Version:     "1.0.0",
		Author:      "goflayer",
		Description: "Health/food tracking",
		Dependencies: []string{},
	}
}

// OnLoad initializes the health plugin.
func (p *HealthPlugin) OnLoad(ctx *plugins.Context) error {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.health = 20
	p.maxHealth = 20
	p.food = 20
	p.saturation = 5
	p.exhaustion = 0

	p.log.Info("Health plugin loaded")

	// Register event handlers
	p.ctx.On("health_update", p.handleHealthUpdate)
	p.ctx.On("food_update", p.handleFoodUpdate)
	p.ctx.On("respawn", p.handleRespawn)

	return nil
}

// OnUnload cleans up the health plugin.
func (p *HealthPlugin) OnUnload() error {
	p.log.Info("Health plugin unloaded")
	return nil
}

// GetHealth returns the current health.
func (p *HealthPlugin) GetHealth() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

// GetMaxHealth returns the maximum health.
func (p *HealthPlugin) GetMaxHealth() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.maxHealth
}

// GetFood returns the current food level.
func (p *HealthPlugin) GetFood() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.food
}

// GetSaturation returns the current saturation level.
func (p *HealthPlugin) GetSaturation() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.saturation
}

// IsDead returns true if health is 0.
func (p *HealthPlugin) IsDead() bool {
	return p.GetHealth() <= 0
}

// IsHungry returns true if food is below 18.
func (p *HealthPlugin) IsHungry() bool {
	return p.GetFood() < 18
}

// IsStarving returns true if food is 0.
func (p *HealthPlugin) IsStarving() bool {
	return p.GetFood() <= 0
}

// GetHealthPercentage returns health as a percentage.
func (p *HealthPlugin) GetHealthPercentage() float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.maxHealth == 0 {
		return 0
	}
	return (p.health / p.maxHealth) * 100
}

// GetFoodPercentage returns food as a percentage.
func (p *HealthPlugin) GetFoodPercentage() float32 {
	food := p.GetFood()
	return (float32(food) / 20) * 100
}

// WaitUntilHealthy waits until health is above a threshold.
func (p *HealthPlugin) WaitUntilHealthy(threshold float32) error {
	p.log.Info("Waiting for health above %.1f", threshold)

	// TODO: Implement wait logic with timeout
	return nil
}

// WaitUntilFull waits until health is full.
func (p *HealthPlugin) WaitUntilFull() error {
	return p.WaitUntilHealthy(p.GetMaxHealth())
}

// Event handlers

func (p *HealthPlugin) handleHealthUpdate(args ...interface{}) {
	if len(args) < 1 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if health, ok := args[0].(float32); ok {
		p.health = health
		p.log.Debug("Health updated: %.1f", p.health)
		p.ctx.Emit("health_changed", p.health)
	}

	if len(args) > 1 {
		if maxHealth, ok := args[1].(float32); ok {
			p.maxHealth = maxHealth
		}
	}

	if p.health <= 0 {
		p.ctx.Emit("death")
	}
}

func (p *HealthPlugin) handleFoodUpdate(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(args) > 0 {
		if food, ok := args[0].(int); ok {
			p.food = food
			p.log.Debug("Food updated: %d", p.food)
		}
	}

	if len(args) > 1 {
		if saturation, ok := args[1].(float32); ok {
			p.saturation = saturation
		}
	}

	p.ctx.Emit("food_changed", p.food, p.saturation)
}

func (p *HealthPlugin) handleRespawn(args ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.health = p.maxHealth
	p.food = 20
	p.saturation = 5

	p.log.Info("Respawned - health and food restored")
	p.ctx.Emit("respawned")
}

// String returns a string representation of the health plugin.
func (p *HealthPlugin) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return fmt.Sprintf("Health{health=%.1f/%.1f, food=%d, saturation=%.1f}",
		p.health, p.maxHealth, p.food, p.saturation)
}
