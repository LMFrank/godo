package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"os/exec"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a specified IP or multiple IPs from a YAML file.",
	Long: `Ping a specified IP address or multiple IPs from a YAML file and print the results to the console or a file.
For example:
godo ping 8.8.8.8`,
	Args: cobra.RangeArgs(0, 1),
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

type PingConfig struct {
	Hosts []string `yaml:"hosts"`
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the YAML config file "+
		"containing multiple IPs.")
}

func pingIP(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		return
	}

	cmd := exec.Command("ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Ping command failed with error: %s: %v\n", ip, err)
		return
	}

	fmt.Printf("Ping results for %s:\n%s", ip, output)
}

func pingMultipleIPs(configFile string) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	var config PingConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}

	file, err := os.Create("response.csv")
	if err != nil {
		fmt.Printf("Error creating response file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "Result"})

	for _, ip := range config.Hosts {
		result := pingSingleIP(ip)
		writer.Write([]string{ip, result})
	}

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

	return string(output)
}
