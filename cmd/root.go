package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vpnctl",
	Short: "Multi-tunnel VPN manager for WireGuard and OpenVPN",
	Long: `vpnctl manages multiple VPN tunnels from a single CLI.

Import WireGuard (.conf) and OpenVPN (.ovpn) configs, then
start/stop them by name. Includes kill switch, DNS leak test,
and speed test.`,
}

// Execute is called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
