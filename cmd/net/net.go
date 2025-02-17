package net

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "net",
	Short: "Network-related commands.",
	Long:  `A collection of network-related commands.`,
}

func init() {
	RootCmd.AddCommand(pingCmd)
}
