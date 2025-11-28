package cmd

import (
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		if scriptName == "" {
			log.Fatalln("Error: script name cannot be empty")
		}

		if err := c.ActivateScript(scriptName); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}
