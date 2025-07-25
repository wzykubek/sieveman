package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deactivateCmd)
}

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate all scripts",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := c.DeactivateScripts(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
