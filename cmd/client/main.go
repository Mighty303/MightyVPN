package main

import (
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/mighty303/govpn/internal/config"
	"github.com/mighty303/govpn/internal/forwarder"
	"github.com/mighty303/govpn/internal/util"
)

func main() {
	serverAddr := flag.String("server", "127.0.0.1:8080", "VPN server address")
	localIP := flag.String("local-ip", "10.0.0.2", "Local TUN IP address")
	remoteIP := flag.String("remote-ip", "10.0.0.1", "Remote TUN IP address")
	flag.Parse()

	config.SetupLogger(slog.LevelDebug)

	slog.Info("Starting GoVPN Client", "server", *serverAddr)

	// 1. Create and configure TUN device
	tun := config.SetupTUN(*localIP, *remoteIP)
	defer tun.Close()

	// 2. Connect to server via UDP
	conn := connectToServer(*serverAddr)
	defer conn.Close()

	// 3. Start packet forwarding loops
	go forwarder.TUNToUDP(tun, conn)
	go forwarder.UDPToTUN(conn, tun)

	slog.Info("Client started - Press Ctrl+C to stop")

	util.WaitForShutdown()
}

// connectToServer establishes UDP connection to server
func connectToServer(serverAddr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		slog.Error("Failed to resolve server address", "error", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to server", "address", serverAddr)
	return conn
}
