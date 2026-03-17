// Command bot - главный entry point с Graceful Shutdown
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nithdevv/goflayer/internal/logger"
	"github.com/nithdevv/goflayer/internal/protocol"
	"github.com/nithdevv/goflayer/internal/types"
	"github.com/nithdevv/goflayer/pkg/bot"
)

func main() {
	// Настройка логирования
	logger.Init(os.Stdout, logger.INFO)
	log := logger.Default()

	log.Info("========================================")
	log.Info("    Minecraft Bot - Professional")
	log.Info("========================================")

	// Создание бота с конфигурацией
	config := types.DefaultBotConfig()
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnvInt("SERVER_PORT", 25565)
	config.Player.Username = getEnv("USERNAME", "AutosellerBot")
	config.AuthMode = getEnv("AUTH_MODE", "offline")

	// Дополнительные настройки
	config.EnableReconnect = true
	config.MaxReconnects = 5
	config.WorkerCount = 4

	log.Info("Configuration:")
	log.Info("  Server: %s:%d", config.Server.Host, config.Server.Port)
	log.Info("  Username: %s", config.Player.Username)
	log.Info("  Auth: %s", config.AuthMode)

	// Создание бота
	log.Info("Creating bot instance...")
	b, err := bot.New(config)
	if err != nil {
		log.Fatal("Failed to create bot: %v", err)
	}

	// Подписка на события
	setupEventHandlers(b, log)

	// Graceful shutdown setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)

	// Запуск бота в горутине
	errCh := make(chan error, 1)
	go func() {
		if err := b.Connect(ctx); err != nil {
			errCh <- fmt.Errorf("connection failed: %w", err)
		}
	}()

	// Главный цикл
	log.Info("========================================")
	log.Info("  Bot is running! Press Ctrl+C to stop")
	log.Info("========================================")

	select {
	case err := <-errCh:
		log.Error("Fatal error: %v", err)
		b.Disconnect()
		os.Exit(1)

	case sig := <-shutdown:
		log.Info("")
		log.Info("========================================")
		log.Info("  Received signal: %v", sig)
		log.Info("  Shutting down gracefully...")
		log.Info("========================================")

		// Graceful shutdown с timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		// Сначала отменяем контекст бота
		cancel()

		// Ждём завершения или timeout
		done := make(chan struct{})
		go func() {
			b.Disconnect()
			close(done)
		}()

		select {
		case <-done:
			log.Info("✓ Shutdown completed successfully")
		case <-shutdownCtx.Done():
			log.Warn("⚠ Shutdown timed out, forcing exit")
		}
	}
}

func setupEventHandlers(b *bot.Bot, log *logger.Logger) {
	b.Events().Subscribe("connected", func(data ...interface{}) {
		log.Info("✓ Connected to server!")
	})

	b.Events().Subscribe("disconnected", func(data ...interface{}) {
		log.Warn("✗ Disconnected from server")
	})

	b.Events().Subscribe("error", func(data ...interface{}) {
		if err, ok := data[0].(error); ok {
			log.Error("Error: %v", err)
		}
	})

	b.Events().Subscribe("logged_in", func(data ...interface{}) {
		uuid := data[0].(string)
		username := data[1].(string)
		log.Info("✓ Logged in as %s (UUID: %s)", username, uuid)
	})

	b.Events().Subscribe("state_changed", func(data ...interface{}) {
		oldState := data[0].(protocol.State)
		newState := data[1].(protocol.State)
		log.Debug("State: %s -> %s", oldState, newState)
	})

	b.Events().Subscribe("reconnect_failed", func(data ...interface{}) {
		log.Error("✗ Failed to reconnect after all attempts")
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
