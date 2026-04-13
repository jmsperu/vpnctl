package cmd

import (
	"fmt"

	"github.com/jmsperu/vpnctl/internal/netutil"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Connection test (public IP, DNS leak check)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking public IP...")
		ip, err := netutil.PublicIP()
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Public IP: %s\n", ip)
		}

		fmt.Println("\nDNS leak test...")
		servers, err := netutil.DNSLeakTest()
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Println("  DNS servers seen:")
			for _, s := range servers {
				fmt.Printf("    %s\n", s)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
