package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/config"
	"github.com/jmsperu/vpnctl/internal/tunnel"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down <name>",
	Short: "Disconnect a tunnel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		t := cfg.Find(args[0])
		if t == nil {
			return fmt.Errorf("tunnel %q not found", args[0])
		}

		_ = tunnel.DisableKillSwitch(t)

		if err := tunnel.Down(t); err != nil {
			return err
		}
		fmt.Printf("Tunnel %q is down\n", t.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
