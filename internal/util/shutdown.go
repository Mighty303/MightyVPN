package util

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// WaitForShutdown blocks until an interrupt signal is received
func WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	slog.Info("Shutting down...")
}
