package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/config"
	"github.com/jmsperu/vpnctl/internal/tunnel"
	"github.com/spf13/cobra"
)

var killSwitch bool

var upCmd = &cobra.Command{
	Use:   "up <name>",
	Short: "Connect a tunnel",
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

		if err := tunnel.Up(t); err != nil {
			return err
		}
		fmt.Printf("Tunnel %q is up\n", t.Name)

		if killSwitch || t.KillSwitch {
			if err := tunnel.EnableKillSwitch(t); err != nil {
				fmt.Printf("Warning: kill switch failed: %v\n", err)
			} else {
				fmt.Println("Kill switch enabled")
			}
		}
		return nil
	},
}

func init() {
	upCmd.Flags().BoolVar(&killSwitch, "kill-switch", false, "Enable kill switch (block traffic if VPN drops)")
	rootCmd.AddCommand(upCmd)
}
