package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deactivateCmd)
}

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate all scripts",
	Long:  "Only one script can be active at a time, so this command simply deactivates active script.",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := c.DeactivateScripts(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}
