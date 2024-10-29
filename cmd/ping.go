package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os/exec"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a specified IP or multiple IPs from a YAML file.",
	Long: `Ping a specified IP address and print the results to the console.
For example:
godo ping 8.8.8.8`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		pingIP(ip)
	},
}

func pingIP(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		return
	}

	cmd := exec.Command("ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ping command failed with error: %s\n", err)
		return
	}

	fmt.Printf("Ping results for %s:\n%s", ip, output)
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
