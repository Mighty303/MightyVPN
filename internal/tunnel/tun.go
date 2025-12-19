package tunnel

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/songgao/water"
)

// TUN wraps the water TUN device with helper methods
type TUN struct {
	iface *water.Interface
}

// NewTUN creates and configures a new TUN device
func NewTUN() (*TUN, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}

	iface, err := water.New(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN device: %w", err)
	}

	return &TUN{iface: iface}, nil
}

// Name returns the TUN device name
func (t *TUN) Name() string {
	return t.iface.Name()
}

// Read reads a packet from the TUN device
func (t *TUN) Read(buf []byte) (int, error) {
	return t.iface.Read(buf)
}

// Write writes a packet to the TUN device
func (t *TUN) Write(buf []byte) (int, error) {
	return t.iface.Write(buf)
}

// Close closes the TUN device
func (t *TUN) Close() error {
	return t.iface.Close()
}

// Configure configures the TUN device with the specified IP addresses
func (t *TUN) Configure(localIP, remoteIP string) error {
	if runtime.GOOS == "darwin" {
		return t.configureMacOS(localIP, remoteIP)
	}
	return t.configureLinux(localIP, remoteIP)
}

func (t *TUN) configureMacOS(localIP, remoteIP string) error {
	// Set interface address and bring it up
	if err := t.runCommand("ifconfig", t.Name(), localIP, remoteIP, "up"); err != nil {
		return fmt.Errorf("failed to configure interface: %w", err)
	}

	// Add route to VPN subnet (ignore error if route already exists)
	_ = t.runCommand("route", "add", "-net", "10.0.0.0/24", "-interface", t.Name())

	return nil
}

func (t *TUN) configureLinux(localIP, remoteIP string) error {
	// Add IP address to interface
	if err := t.runCommand("ip", "addr", "add", localIP+"/24", "dev", t.Name()); err != nil {
		return err
	}

	// Bring interface up
	if err := t.runCommand("ip", "link", "set", "dev", t.Name(), "up"); err != nil {
		return err
	}

	return nil
}

// runCommand executes a shell command and returns an error if it fails
func (t *TUN) runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed [%s %v]: %w", name, args, err)
	}
	return nil
}