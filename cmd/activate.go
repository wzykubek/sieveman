package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(activateCmd)
}

var activateCmd = &cobra.Command{
	Use:   "activate <script_name>",
	Short: "Activate a script",
	Long:  "Activate a script with given name. Keep in mind that in most cases only one script will be active at a time, so when there was an active script, it will be deactivated by a server.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptName := args[0]
		if scriptName == "" {
			return fmt.Errorf("script name cannot be empty\n")
		}

		if err := c.ActivateScript(scriptName); err != nil {
			return err
		}

		return nil
	},
}
