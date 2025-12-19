package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/mighty303/govpn/internal/config"
	"github.com/mighty303/govpn/internal/forwarder"
	"github.com/mighty303/govpn/internal/tunnel"
	"github.com/mighty303/govpn/internal/util"
)

func main() {
	listenAddr := flag.String("listen", ":8080", "UDP listen address")
	localIP := flag.String("local-ip", "10.0.0.1", "Local TUN IP address")
	remoteIP := flag.String("remote-ip", "10.0.0.2", "Remote TUN IP address")
	flag.Parse()

	config.SetupLogger(slog.LevelDebug)

	slog.Info("Starting GoVPN Server", "listen", *listenAddr)

	// 1. Create and configure TUN device
	tun := config.SetupTUN(*localIP, *remoteIP)
	defer tun.Close()

	// 2. Set up UDP listener
	conn := setupUDPListener(*listenAddr)
	defer conn.Close()

	// 3. Start packet forwarding with client tracking
	startPacketForwarding(tun, conn)

	slog.Info("Server started - Press Ctrl+C to stop")

	util.WaitForShutdown()
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

// startPacketForwarding starts bidirectional forwarding with client address tracking
func startPacketForwarding(tun *tunnel.TUN, conn *net.UDPConn) {
	var clientAddr *net.UDPAddr
	var clientMutex sync.RWMutex

	// Goroutine 1: TUN → UDP (with client address checking)
	go forwardTUNtoUDP(tun, conn, &clientAddr, &clientMutex)

	// Goroutine 2: UDP → TUN (with client address learning)
	go forwardUDPtoTUN(conn, tun, &clientAddr, &clientMutex)
}

// forwardTUNtoUDP forwards packets from TUN to the connected UDP client
func forwardTUNtoUDP(tun *tunnel.TUN, conn *net.UDPConn, clientAddr **net.UDPAddr, mutex *sync.RWMutex) {
	buf := make([]byte, forwarder.MTU)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			slog.Error("Error reading from TUN device", "error", err)
			continue
		}

		mutex.RLock()
		addr := *clientAddr
		mutex.RUnlock()

		if addr == nil {
			slog.Debug("No client connected, dropping packet", "size", n)
			continue
		}

		_, err = conn.WriteToUDP(buf[:n], addr)
		if err != nil {
			slog.Error("Error forwarding packet to client", "error", err)
			continue
		}
		slog.Debug("Forwarded packet TUN → UDP", "size", n, "client", addr.String())
	}
}

// forwardUDPtoTUN forwards packets from UDP to TUN and learns client address
func forwardUDPtoTUN(conn *net.UDPConn, tun *tunnel.TUN, clientAddr **net.UDPAddr, mutex *sync.RWMutex) {
	buf := make([]byte, forwarder.MTU)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("Error reading from UDP", "error", err)
			continue
		}

		// Update client address
		mutex.Lock()
		if *clientAddr == nil {
			slog.Info("Client connected", "address", addr.String())
			*clientAddr = addr
		} else if (*clientAddr).String() != addr.String() {
			slog.Info("Client address updated", "old", (*clientAddr).String(), "new", addr.String())
			*clientAddr = addr
		}
		mutex.Unlock()

		_, err = tun.Write(buf[:n])
		if err != nil {
			slog.Error("Error writing to TUN device", "error", err)
			continue
		}
		slog.Debug("Forwarded packet UDP → TUN", "size", n, "from", addr.String())
	}
}
