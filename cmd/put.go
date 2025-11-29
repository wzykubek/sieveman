package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var activate bool

func init() {
	putCmd.Flags().SortFlags = false
	putCmd.Flags().BoolVarP(&activate, "activate", "a", false, "mark the uploaded script as active")
	rootCmd.AddCommand(putCmd)
}

var putCmd = &cobra.Command{
	Use:   "put <local_name> [remote_name]",
	Short: "Upload local script to the server",
	Long: `You can use this command to upload local script.
If remote name is not specified, script will be uploaded with the same name as local file.

When you want to make changes to remote script and upload it immediately you should consider using edit command instead.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		remoteName := filename
		if len(args) > 1 {
			remoteName = args[1]
		}

		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		if err = c.PutScript(file, remoteName); err != nil {
			return err
		}

		if activate {
			if err = c.ActivateScript(remoteName); err != nil {
				return err
			}
		}

		return nil
	},
}
