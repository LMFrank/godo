package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os/exec"
	"strings"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a specified IP or multiple IPs from a YAML file.",
	Long: `Ping a specified IP address or multiple IPs from a YAML file and print the results to the console or a file.
For example:
godo ping 8.8.8.8`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		if configFile != "" {
			pingMultipleIPs(configFile)
		} else if len(args) > 0 {
			ip := args[0]
			pingIP(ip)
		} else {
			fmt.Println("Please provide an IP address or a config file")
		}
	},
}

var configFile string

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the YAML config file containing multiple IPs.")
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

func pingMultipleIPs(configFile string) {

}

func pingSingleIP(ip string) string {
	if net.ParseIP(ip) == nil {
		return fmt.Sprintf("Invalid IP address: %s\n", ip)
	}

	cmd := exec.Command("ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error pinging %s: %v\n", ip, err)
	}

	lines := strings.Split(string(output), "\n")
	var result string
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			result += line + "\n"
		}
	}
	return result
}
