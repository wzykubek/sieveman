package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var activate bool

func init() {
	putCmd.Flags().SortFlags = false
	putCmd.Flags().BoolVarP(&activate, "activate", "a", false, "set the uploaded script as active")
	rootCmd.AddCommand(putCmd)
}

var putCmd = &cobra.Command{
	Use:   "put <local_name> [remote_name]",
	Short: "Upload a local script to the server",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		remoteName := filename
		if len(args) > 1 {
			remoteName = args[1]
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		defer file.Close()

		if err = c.PutScript(file, remoteName); err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		if activate {
			if err = c.ActivateScript(remoteName); err != nil {
				log.Fatalf("Error: %s\n", err)
			}
		}
	},
}
