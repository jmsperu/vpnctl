package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/config"
	"github.com/jmsperu/vpnctl/internal/tunnel"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of all tunnels",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Tunnels) == 0 {
			fmt.Println("No tunnels configured. Use 'vpnctl import' to add one.")
			return nil
		}
		fmt.Printf("%-20s %-12s %-6s %s\n", "NAME", "TYPE", "STATE", "CONFIG")
		fmt.Println("------------------------------------------------------------")
		for i := range cfg.Tunnels {
			fmt.Println(tunnel.Status(&cfg.Tunnels[i]))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
