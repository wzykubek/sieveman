package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mvCmd)
}

var mvCmd = &cobra.Command{
	Use:   "mv <old_name> <new_name>",
	Short: "Rename script",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]

		if err := c.RenameScript(oldName, newName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}
