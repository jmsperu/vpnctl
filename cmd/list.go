package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved tunnels",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Tunnels) == 0 {
			fmt.Println("No tunnels configured.")
			return nil
		}
		fmt.Printf("%-20s %-12s %-12s %s\n", "NAME", "TYPE", "AUTO", "CONFIG")
		fmt.Println("------------------------------------------------------------")
		for _, t := range cfg.Tunnels {
			auto := "no"
			if t.AutoConnect {
				auto = "yes"
			}
			fmt.Printf("%-20s %-12s %-12s %s\n", t.Name, t.Type, auto, t.ConfigPath)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
