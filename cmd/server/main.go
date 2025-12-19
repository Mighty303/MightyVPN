package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mighty303/govpn/internal/tunnel"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "UDP listen address")
	localIP := flag.String("local-ip", "10.0.0.1", "Local TUN IP address")
	remoteIP := flag.String("remote-ip", "10.0.0.2", "Remote TUN IP address")
	flag.Parse()

	setupLogger()

	slog.Info("Starting GoVPN Server", "listen", *listenAddr)

	// 1. Create and configure TUN device
	tun := setupTUN(*localIP, *remoteIP)
	defer tun.Close()

	// 2. Set up UDP listener
	conn := setupUDPListener(*listenAddr)
	defer conn.Close()

	// 3. Start packet forwarding loops
	startPacketForwarding(tun, conn)

	slog.Info("Server started - Press Ctrl+C to stop")

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

// setupUDPListener creates and binds a UDP listener
func setupUDPListener(listenAddr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		slog.Error("Failed to resolve UDP address", "error", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		slog.Error("Failed to listen on UDP address", "error", err)
		os.Exit(1)
	}

	slog.Info("UDP listener started", "address", conn.LocalAddr().String())
	return conn
}

// startPacketForwarding starts the goroutine that reads packets from TUN
func startPacketForwarding(tun *tunnel.TUN, conn *net.UDPConn) {
	go func() {
		buf := make([]byte, 1500) // MTU (Max Transmission Unit) size
		for {
			n, err := tun.Read(buf)
			if err != nil {
				slog.Error("Error reading from TUN device", "error", err)
				continue
			}
			slog.Debug("Read packet from TUN device", "size", n)
			// TODO: Forward packet to client via UDP
		}
	}()

	// TODO: Add goroutine to read from UDP and write to TUN
}

// waitForShutdown blocks until an interrupt signal is received
func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
