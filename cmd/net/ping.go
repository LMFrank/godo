package net

import (
	"fmt"
	"github.com/LMFrank/godo/pkg/net"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a specified IP or multiple IPs from a YAML file.",
	Long: `Ping a specified IP address or multiple IPs from a YAML file and print the results to the console or a file.
For example:
godo ping 8.8.8.8`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("filename")
		if configFile == "" {
			configFile, _ = cmd.Flags().GetString("f")
		}
		if configFile != "" {
			net.PingMultipleIPs(configFile)
		} else if len(args) > 0 {
			ip := args[0]
			net.PingIP(ip)
		} else {
			fmt.Println("Please provide an IP address or a config file")
		}
	},
}

func init() {
	pingCmd.Flags().StringP("hosts", "f", "", "hosts.yaml 文件路径")
}
