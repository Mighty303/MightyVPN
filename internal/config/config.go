package config

import (
	"log/slog"
	"os"

	"github.com/mighty303/govpn/internal/tunnel"
)

// SetupLogger initializes the structured logger with the given level
func SetupLogger(level slog.Level) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}

// SetupTUN creates and configures a TUN device
func SetupTUN(localIP, remoteIP string) *tunnel.TUN {
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
