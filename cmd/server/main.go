package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"net"

	"github.com/mighty303/govpn/internal/tunnel"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "UDP listen address")
	localIP := flag.String("local-ip", "10.0.0.1", "Local TUN IP address")
	remoteIP := flag.String("remote-ip", "10.0.0.2", "Remote TUN IP address")
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
	
	// 1. Create TUN device
	tun, err := tunnel.NewTUN()
	if err != nil {
		slog.Error("Failed to create TUN device", "error", err)
		os.Exit(1)
	}
	defer tun.Close()

	slog.Info("TUN device created", "name", tun.Name())

	err = tun.Configure(*localIP, *remoteIP)
	if err != nil {
		slog.Error("Failed to configure TUN device", "error", err)
		os.Exit(1)
	}
	slog.Info("TUN device configured", "name", tun.Name())

	// 2. Set up UDP listener
	udpAddr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		slog.Error("Failed to resolve UDP address", "error", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		slog.Error("Failed to listen on UDP address", "error", err)
		os.Exit(1)
	}
	defer conn.Close()
	slog.Info("UDP listener started", "address", conn.LocalAddr().String())

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

	slog.Info("Server started - Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down...")
}
