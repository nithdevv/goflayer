package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-flayer/goflayer/pkg/goflayer"
)

func main() {
	// Create a bot
	bot, err := goflayer.CreateBot(goflayer.Options{
		Host:     "localhost",
		Port:     25565,
		Username: "GoflayerBot",
		Version:  "1.20.1",
	})
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Subscribe to login event
	bot.On("login", func(data ...interface{}) {
		fmt.Println("✓ Logged in successfully!")
	})

	// Subscribe to spawn event
	bot.On("spawn", func(data ...interface{}) {
		fmt.Println("✓ Spawned in the world")
		bot.Chat("Hello from goflayer!")
	})

	// Subscribe to chat events
	bot.On("chat", func(data ...interface{}) {
		// Handle chat messages
		if len(data) >= 2 {
			username := data[0].(string)
			message := data[1].(string)
			fmt.Printf("[CHAT] %s: %s\n", username, message)
		}
	})

	// Subscribe to error events
	bot.On("error", func(data ...interface{}) {
		if len(data) >= 1 {
			err := data[0]
			fmt.Printf("[ERROR] %v\n", err)
		}
	})

	// Subscribe to disconnect event
	bot.On("disconnect", func(data ...interface{}) {
		if len(data) >= 1 {
			reason := data[0]
			fmt.Printf("Disconnected: %v\n", reason)
		}
		fmt.Println("Bot disconnected")
		os.Exit(0)
	})

	// Connect to the server
	fmt.Println("Connecting to server...")
	ctx := context.Background()
	if err := bot.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("Bot is running. Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Disconnect
	fmt.Println("\nDisconnecting...")
	if err := bot.Disconnect(); err != nil {
		log.Printf("Error disconnecting: %v", err)
	}

	fmt.Println("Bot stopped")
}
