package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping [IP]",
	Short: "Ping an IP address",
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		configFile, _ := cmd.Flags().GetString("config")
		if ip != "" && configFile == "" {
			err := pingSingleIP(ip)
			if err != nil {
				fmt.Println(err)
			}
		} else if ip != "" && configFile != "" {
			err2 := pingMultipleIPs(configFile)
			if err2 != nil {
				fmt.Println(err2)
			}
		} else {
			fmt.Println("please provide either an IP or a config file.")
		}
	},
}

func pingSingleIP(ip string) error {
	conn, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		fmt.Printf("Failed to ping %s: %v\n", ip, err)
	}
	err2 := conn.Close()
	if err2 != nil {
		return err2
	}
	fmt.Printf("Ping to %s successful \n", ip)
	return nil
}

func pingMultipleIPs(configFile string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringP("ip", "p", "", "IP address to ping")
}
