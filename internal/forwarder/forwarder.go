package forwarder

import (
	"log/slog"
	"net"

	"github.com/mighty303/govpn/internal/tunnel"
)

const MTU = 1500 // Standard MTU size

// TUNToUDP reads packets from TUN device and sends them via UDP
func TUNToUDP(tun *tunnel.TUN, conn net.Conn) {
	buf := make([]byte, MTU)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			slog.Error("Error reading from TUN device", "error", err)
			continue
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			slog.Error("Error writing to UDP", "error", err)
			continue
		}
		slog.Debug("Forwarded packet TUN → UDP", "size", n)
	}
}

// UDPToTUN reads packets from UDP and writes them to TUN device
func UDPToTUN(conn net.Conn, tun *tunnel.TUN) {
	buf := make([]byte, MTU)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("Error reading from UDP", "error", err)
			continue
		}

		_, err = tun.Write(buf[:n])
		if err != nil {
			slog.Error("Error writing to TUN device", "error", err)
			continue
		}
		slog.Debug("Forwarded packet UDP → TUN", "size", n)
	}
}
