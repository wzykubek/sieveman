package cmd

import (
	"fmt"
	"os"

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
			fmt.Println("Error: Script name cannot be empty.")
			os.Exit(1)
		}

		if err := c.ActivateScript(scriptName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}
