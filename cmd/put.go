package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	putCmd.Flags().SortFlags = false
	// TODO: Make default flag
	rootCmd.AddCommand(putCmd)
}

var putCmd = &cobra.Command{
	Use:   "put <local_name> [remote_name]",
	Short: "Upload a local script to the server.",
	Long:  `This command uploads a local script to the server with optional name argument.`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		filename := args[0]
		remoteName := filename
		if len(args) > 1 {
			remoteName = args[1]
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		err = c.PutScript(file, remoteName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
