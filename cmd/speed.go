package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/netutil"
	"github.com/spf13/cobra"
)

var speedCmd = &cobra.Command{
	Use:   "speed",
	Short: "Run a download speed test on the current connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running speed test (10 MB download via Cloudflare)...")
		result, err := netutil.SpeedTest()
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(speedCmd)
}
