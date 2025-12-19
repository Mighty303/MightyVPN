package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mighty303/govpn/internal/tunnel"
)

func main() {
	serverAddr := flag.String("server", "127.0.0.1:8080", "VPN server address")
	localIP := flag.String("local-ip", "10.0.0.2", "Local TUN IP address")
	remoteIP := flag.String("remote-ip", "10.0.0.1", "Remote TUN IP address")
	flag.Parse()

	setupLogger()

	slog.Info("Starting GoVPN Client", "server", *serverAddr)

	// 1. Create and configure TUN device
	tun := setupTUN(*localIP, *remoteIP)
	defer tun.Close()

	// 2. TODO: Connect to server via UDP
	// conn := connectToServer(*serverAddr)
	// defer conn.Close()

	// 3. Start packet forwarding loops
	startPacketForwarding(tun)

	slog.Info("Client started - Press Ctrl+C to stop")

	// Wait for interrupt
	waitForShutdown()

	slog.Info("Shutting down...")
}

// setupLogger initializes the structured logger
func setupLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

// setupTUN creates and configures the TUN device
func setupTUN(localIP, remoteIP string) *tunnel.TUN {
	tun, err := tunnel.NewTUN()
	if err != nil {
		slog.Error("Failed to create TUN device", "error", err)
		os.Exit(1)
	}

	slog.Info("TUN device created", "name", tun.Name())

	if err := tun.Configure(localIP, remoteIP); err != nil {
		slog.Error("Failed to configure TUN device", "error", err)
		os.Exit(1)
	}

	slog.Info("TUN device configured", "local", localIP, "remote", remoteIP)
	return tun
}

// startPacketForwarding starts the goroutine that reads packets from TUN
func startPacketForwarding(tun *tunnel.TUN) {
	go func() {
		buf := make([]byte, 1500) // MTU (Max Transmission Unit) size
		for {
			n, err := tun.Read(buf)
			if err != nil {
				slog.Error("Error reading from TUN device", "error", err)
				continue
			}
			slog.Debug("Read packet from TUN device", "size", n)
			// TODO: Forward packet to server
		}
	}()
}

// waitForShutdown blocks until an interrupt signal is received
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}

