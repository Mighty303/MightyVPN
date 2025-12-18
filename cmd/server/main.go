package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "UDP listen address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting GoVPN Server", "listen", *listenAddr)

	// TODO: Implement server
	// 1. Create TUN device
	// 2. Set up UDP listener
	// 3. Start packet forwarding loops

	slog.Info("Server started - Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down...")
}