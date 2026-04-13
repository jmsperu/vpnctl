package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/config"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import <name> <path>",
	Short: "Import a WireGuard (.conf) or OpenVPN (.ovpn) config",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, path := args[0], args[1]
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Add(name, path); err != nil {
			return err
		}
		fmt.Printf("Imported %q from %s\n", name, path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
