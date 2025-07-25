package cmd

import (
	"fmt"
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

		if err = c.PutScript(file, remoteName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if activate {
			if err = c.ActivateScript(remoteName); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}
