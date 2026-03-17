// Example bot demonstrates basic usage of goflayer.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nithdevv/goflayer/pkg/bot"
	"github.com/nithdevv/goflayer/pkg/plugins"
	"github.com/nithdevv/goflayer/pkg/plugins/builtin/chat"
	"github.com/nithdevv/goflayer/pkg/plugins/builtin/entities"
	"github.com/nithdevv/goflayer/pkg/plugins/builtin/game"
	"github.com/nithdevv/goflayer/pkg/plugins/builtin/physics"
)

func main() {
	// Create bot options
	options := &bot.Options{
		Host:     "localhost",
		Port:     25565,
		Username: "GoflayerBot",
		Version:  "1.20.1",
		Auth:     "offline",
		Plugins:  make(map[string]plugins.Plugin),
	}

	// Load plugins
	options.Plugins["game"] = game.NewPlugin()
	options.Plugins["entities"] = entities.NewPlugin()
	options.Plugins["physics"] = physics.NewPlugin()
	options.Plugins["chat"] = chat.NewPlugin()

	// Create bot
	b, err := bot.New(options)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Subscribe to events
	b.On("connected", func(data ...interface{}) {
		fmt.Println("Connected to server!")
	})

	b.On("disconnected", func(data ...interface{}) {
		fmt.Println("Disconnected from server")
	})

	b.On("error", func(data ...interface{}) {
		err := data[0].(error)
		log.Printf("Error: %v", err)
	})

	// Subscribe to chat
	b.On("chat", func(data ...interface{}) {
		// Handle chat messages
		fmt.Printf("Chat received: %v\n", data)
	})

	// Connect to server
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer b.Disconnect()

	fmt.Println("Bot is running. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down...")
}
