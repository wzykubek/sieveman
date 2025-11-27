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
	Short: "Activate a script with given name",
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
