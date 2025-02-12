package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/LMFrank/godo/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"sync"
	"time"
)

// PingCommand 结构体扩展了 cobra.Command 并添加了 CommandExecutor 字段
type PingCommand struct {
	*cobra.Command
	executor util.CommandExecutor
}

var pingCmd *PingCommand
var configFile string

func init() {
	pingCommand := &cobra.Command{
		Use:   "ping",
		Short: "Ping a specified IP or multiple IPs from a YAML file.",
		Long: `Ping a specified IP address or multiple IPs from a YAML file and print the results to the console or a file.
For example:
godo ping 8.8.8.8`,
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			configFile, _ := cmd.Flags().GetString("config")
			if configFile != "" {
				pingMultipleIPs(configFile, pingCmd.executor)
			} else if len(args) > 0 {
				ip := args[0]
				pingIP(ip, pingCmd.executor)
			} else {
				fmt.Println("Please provide an IP address or a config file")
			}
		},
	}

	pingCmd = &PingCommand{
		Command:  pingCommand,
		executor: &util.DefaultCommandExecutor{},
	}

	rootCmd.AddCommand(pingCmd.Command)
	pingCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the YAML config file containing multiple IPs.")
}

type PingConfig struct {
	Hosts []string `yaml:"hosts"`
}

func pingIP(ip string, executor util.CommandExecutor) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		return
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	err, result := executor.Execute(command)
	if err != nil {
		fmt.Printf("Ping command failed with error: %s: %v\n", ip, err)
		return
	}

	fmt.Printf("Ping results for %s:\n%s", ip, result)
}

func pingMultipleIPs(configFile string, executor util.CommandExecutor) {
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

	timestamp := time.Now().Format("20060102_15-04-05")
	dir := "./response"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	filename := fmt.Sprintf("%s/response_%s.csv", dir, timestamp)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating response file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "Result"})

	results := make(chan []string)
	var wg sync.WaitGroup

	for _, ip := range config.Hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			result := pingSingleIP(ip, executor)
			results <- []string{ip, result}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		writer.Write(result)
	}
}

func pingSingleIP(ip string, executor util.CommandExecutor) string {
	if net.ParseIP(ip) == nil {
		return fmt.Sprintf("Invalid IP address: %s\n", ip)
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	err, result := executor.Execute(command)
	if err != nil {
		return fmt.Sprintf("Error pinging %s: %v\n", ip, err)
	}

	return result
}
