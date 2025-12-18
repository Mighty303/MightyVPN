package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"go/govpn/internal/tunnel"
)

func main() {
	serverAddr := flag.String("server", "127.0.0.1:8080", "VPN server address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting GoVPN Client", "server", *serverAddr)

	// TODO: Implement client
	// 1. Create TUN device
	// 2. Connect to server via UDP
	// 3. Start packet forwarding loops

	slog.Info("Client started - Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down...")
}

func createClientTUN(localIP, remoteIP string) {

}
