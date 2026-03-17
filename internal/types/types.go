// Package types содержит базовые типы для всего проекта.
package types

import (
	"time"
)

// GameState представляет текущее состояние бота.
type GameState int

const (
	StateDisconnected GameState = iota
	StateConnecting
	StateHandshaking
	StateLoggingIn
	StatePlay
	StateReconnecting
)

func (s GameState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateHandshaking:
		return "handshaking"
	case StateLoggingIn:
		return "logging_in"
	case StatePlay:
		return "play"
	case StateReconnecting:
		return "reconnecting"
	default:
		return "unknown"
	}
}

// ServerInfo содержит информацию о сервере.
type ServerInfo struct {
	Host     string
	Port     int
	Protocol int
}

// PlayerInfo содержит информацию об игроке.
type PlayerInfo struct {
	UUID     string
	Username string
}

// BotConfig содержит конфигурацию бота.
type BotConfig struct {
	Server         ServerInfo
	Player         PlayerInfo
	AuthMode       string // "offline" or "microsoft"

	// Networking
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	// Reconnection
	EnableReconnect  bool
	MaxReconnects    int
	ReconnectDelay   time.Duration
	ReconnectBackoff float64 // multiplier for delay

	// Workers
	WorkerCount int // number of packet workers
}

// DefaultBotConfig возвращает конфигурацию по умолчанию.
func DefaultBotConfig() BotConfig {
	return BotConfig{
		Server: ServerInfo{
			Host:     "localhost",
			Port:     25565,
			Protocol: 763, // 1.20.1
		},
		Player: PlayerInfo{
			Username: "Bot",
		},
		AuthMode: "offline",

		ConnectTimeout: 10 * time.Second,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Second,

		EnableReconnect:  true,
		MaxReconnects:    5,
		ReconnectDelay:   2 * time.Second,
		ReconnectBackoff: 1.5,

		WorkerCount: 4,
	}
}
