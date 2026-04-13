package tunnel

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jmsperu/vpnctl/internal/config"
)

// Up brings a tunnel up.
func Up(t *config.Tunnel) error {
	switch t.Type {
	case config.WireGuard:
		return run("wg-quick", "up", t.ConfigPath)
	case config.OpenVPN:
		return run("openvpn", "--config", t.ConfigPath, "--daemon", "--log", "/tmp/vpnctl-"+t.Name+".log")
	default:
		return fmt.Errorf("unknown tunnel type %q", t.Type)
	}
}

// Down tears a tunnel down.
func Down(t *config.Tunnel) error {
	switch t.Type {
	case config.WireGuard:
		return run("wg-quick", "down", t.ConfigPath)
	case config.OpenVPN:
		// Find and kill the openvpn process for this config.
		out, err := exec.Command("pgrep", "-f", t.ConfigPath).Output()
		if err != nil {
			return fmt.Errorf("no running openvpn process found for %s", t.Name)
		}
		for _, pid := range strings.Fields(strings.TrimSpace(string(out))) {
			_ = exec.Command("kill", pid).Run()
		}
		return nil
	default:
		return fmt.Errorf("unknown tunnel type %q", t.Type)
	}
}

// IsUp checks whether a tunnel appears to be active.
func IsUp(t *config.Tunnel) bool {
	switch t.Type {
	case config.WireGuard:
		// wg show lists active interfaces; our config name is embedded in the interface.
		out, err := exec.Command("wg", "show", "interfaces").Output()
		if err != nil {
			return false
		}
		iface := extractWGInterface(t.ConfigPath)
		return strings.Contains(string(out), iface)
	case config.OpenVPN:
		err := exec.Command("pgrep", "-f", t.ConfigPath).Run()
		return err == nil
	}
	return false
}

// Status returns a human-readable status line for a tunnel.
func Status(t *config.Tunnel) string {
	up := IsUp(t)
	state := "down"
	if up {
		state = "up"
	}
	return fmt.Sprintf("%-20s %-12s %-6s %s", t.Name, t.Type, state, t.ConfigPath)
}

// EnableKillSwitch adds firewall rules to block non-VPN traffic.
func EnableKillSwitch(t *config.Tunnel) error {
	switch t.Type {
	case config.WireGuard:
		iface := extractWGInterface(t.ConfigPath)
		cmds := [][]string{
			{"iptables", "-I", "OUTPUT", "!", "-o", iface, "-m", "mark", "!", "--mark", "0xca6c", "-j", "DROP"},
			{"iptables", "-I", "OUTPUT", "-o", "lo", "-j", "ACCEPT"},
		}
		for _, c := range cmds {
			if err := run(c[0], c[1:]...); err != nil {
				return fmt.Errorf("kill switch rule failed: %w", err)
			}
		}
		return nil
	default:
		return fmt.Errorf("kill switch only supported for WireGuard tunnels currently")
	}
}

// DisableKillSwitch removes the kill switch rules.
func DisableKillSwitch(t *config.Tunnel) error {
	switch t.Type {
	case config.WireGuard:
		iface := extractWGInterface(t.ConfigPath)
		_ = exec.Command("iptables", "-D", "OUTPUT", "!", "-o", iface, "-m", "mark", "!", "--mark", "0xca6c", "-j", "DROP").Run()
		_ = exec.Command("iptables", "-D", "OUTPUT", "-o", "lo", "-j", "ACCEPT").Run()
		return nil
	default:
		return nil
	}
}

// extractWGInterface derives the interface name from the config path.
// wg-quick uses the filename (sans extension) as the interface name.
func extractWGInterface(path string) string {
	base := path
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	if strings.HasSuffix(base, ".conf") {
		base = base[:len(base)-5]
	}
	return base
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s: %w\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return nil
}
