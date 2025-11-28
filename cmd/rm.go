package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm <script_name>",
	Short: "Remove a script with the given name",
	Long:  "Remove remote script on server. This operation is irreversible, be careful.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]

		if err := c.RemoveScript(scriptName); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}
