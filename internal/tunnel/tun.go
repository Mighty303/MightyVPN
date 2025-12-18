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

// ConfigureClient configures the TUN device for client mode
func (t *TUN) ConfigureClient(localIP, remoteIP string) error {
	if runtime.GOOS == "darwin" {
		return t.configureMacOS(localIP, remoteIP)
	}
	return t.configureLinux(localIP, remoteIP)
}

// ConfigureServer configures the TUN device for server mode
func (t *TUN) ConfigureServer(localIP, remoteIP string) error {
	if runtime.GOOS == "darwin" {
		return t.configureMacOS(localIP, remoteIP)
	}
	return t.configureLinux(localIP, remoteIP)
}

func (t *TUN) configureMacOS(localIP, remoteIP string) error {
	cmd := exec.Command("ifconfig", t.Name(), localIP, remoteIP, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure interface: %w", err)
	}

	cmd = exec.Command("route", "add", "-net", "10.0.0.0/24", "-interface", t.Name())
	_ = cmd.Run() // Ignore error if route exists

	return nil
}

func (t *TUN) configureLinux(localIP, remoteIP string) error {
	cmds := [][]string{
		{"ip", "addr", "add", localIP + "/24", "dev", t.Name()},
		{"ip", "link", "set", "dev", t.Name(), "up"},
	}

	for _, cmd := range cmds {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("command failed %v: %w", cmd, err)
		}
	}

	return nil
}