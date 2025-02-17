package set

import (
	"github.com/LMFrank/godo/pkg/set"
	"github.com/spf13/cobra"
)

var yumCmd = &cobra.Command{
	Use:   "yum",
	Short: "Set YUM source for CentOS systems.",
	Long:  `Set YUM source for CentOS systems, supporting CentOS 6, 7, and 8.`,
	Run: func(cmd *cobra.Command, args []string) {
		set.SetCentOSYumSource()
	},
}
