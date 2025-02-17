package set

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "set",
	Short: "Set-related commands.",
	Long:  `A collection of set-related commands.`,
}

func init() {
	RootCmd.AddCommand(yumCmd)
}
