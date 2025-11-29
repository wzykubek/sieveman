package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mvCmd)
}

var mvCmd = &cobra.Command{
	Use:   "mv <old_name> <new_name>",
	Short: "Rename script",
	Long:  "This command renames remote scripts on server.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldName := args[0]
		newName := args[1]

		if err := c.RenameScript(oldName, newName); err != nil {
			return err
		}

		return nil
	},
}
