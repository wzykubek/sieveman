package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getCmd.Flags().SortFlags = false
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get <script_name> <output_file*>",
	Short: "Download specified script to given path",
	Long: `This command downloads specified script and saves it to given location.
You can use '-' character as file name to print to stdout.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		scriptName := args[0]
		var outFile string = args[1]

		s, err := c.GetScript(scriptName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if outFile == "-" {
			fmt.Println(s)
		} else {
			if err := os.WriteFile(outFile, []byte(s), 0644); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	},
}
