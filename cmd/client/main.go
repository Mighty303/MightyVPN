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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting GoVPN Client", "server", *serverAddr)

	// 1. Create TUN device
	tun, err := tunnel.NewTUN()
	if err != nil {
		slog.Error("Failed to create TUN device", "error", err)
		os.Exit(1)
	}
	defer tun.Close()

	slog.Info("TUN device created", "name", tun.Name())

	// 2. Configure TUN device
	err = tun.Configure(*localIP, *remoteIP)
	if err != nil {
		slog.Error("Failed to configure TUN device", "error", err)
		os.Exit(1)
	}

	slog.Info("TUN device created", "name", tun.Name())
	
	// 3. Start packet forwarding loops
	go func() {
		buf := make([]byte, 1500) // MTU(Max Transmission Unit) size
		for {
			n, err := tun.Read(buf)
			if err != nil {
				slog.Error("Error reading from TUN device", "error", err)
				continue
			}
			slog.Debug("Read packet from TUN device", "size", n)
		}
	}()

	slog.Info("Client started - Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down...")
}

