package main

import (
	"github.com/LMFrank/godo/cmd/net"
	"github.com/LMFrank/godo/cmd/set"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "godo",
	Short: "Godo is a convenient operation and maintenance tool based on Golang Cobra",
	Long: `Godo is a CLI tool for network monitoring and management, built on Golang Cobra.
It provides convenient commands for system administrators to perform common network tasks.
Examples:
  godo net ping 8.8.8.8
  godo set yum update`,
}

func init() {
	rootCmd.AddCommand(net.RootCmd)
	rootCmd.AddCommand(set.RootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
