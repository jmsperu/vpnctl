package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TunnelType distinguishes WireGuard from OpenVPN configs.
type TunnelType string

const (
	WireGuard TunnelType = "wireguard"
	OpenVPN   TunnelType = "openvpn"
)

// Tunnel holds metadata about an imported VPN tunnel.
type Tunnel struct {
	Name        string     `yaml:"name"`
	Type        TunnelType `yaml:"type"`
	ConfigPath  string     `yaml:"config_path"`
	AutoConnect bool       `yaml:"auto_connect"`
	KillSwitch  bool       `yaml:"kill_switch"`
}

// Config is the top-level structure persisted in ~/.vpnctl.yml.
type Config struct {
	Tunnels []Tunnel `yaml:"tunnels"`
}

// Path returns the config file location.
func Path() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vpnctl.yml")
}

// configDir returns the directory where imported configs are stored.
func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".vpnctl.d")
}

// ConfigDir is the exported accessor.
func ConfigDir() string { return configDir() }

// Load reads the config file, returning an empty Config if it doesn't exist.
func Load() (*Config, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0600)
}

// Find returns the tunnel with the given name, or nil.
func (c *Config) Find(name string) *Tunnel {
	for i := range c.Tunnels {
		if c.Tunnels[i].Name == name {
			return &c.Tunnels[i]
		}
	}
	return nil
}

// Add imports a config file into ~/.vpnctl.d/ and saves a reference.
func (c *Config) Add(name, path string) error {
	if c.Find(name) != nil {
		return fmt.Errorf("tunnel %q already exists", name)
	}

	ext := filepath.Ext(path)
	var ttype TunnelType
	switch ext {
	case ".conf":
		ttype = WireGuard
	case ".ovpn":
		ttype = OpenVPN
	default:
		return fmt.Errorf("unsupported file extension %q (use .conf for WireGuard, .ovpn for OpenVPN)", ext)
	}

	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading source config: %w", err)
	}

	dest := filepath.Join(dir, name+ext)
	if err := os.WriteFile(dest, src, 0600); err != nil {
		return fmt.Errorf("writing config copy: %w", err)
	}

	c.Tunnels = append(c.Tunnels, Tunnel{
		Name:       name,
		Type:       ttype,
		ConfigPath: dest,
	})
	return c.Save()
}

// Remove deletes a tunnel by name, including its config file.
func (c *Config) Remove(name string) error {
	for i, t := range c.Tunnels {
		if t.Name == name {
			_ = os.Remove(t.ConfigPath)
			c.Tunnels = append(c.Tunnels[:i], c.Tunnels[i+1:]...)
			return c.Save()
		}
	}
	return fmt.Errorf("tunnel %q not found", name)
}
