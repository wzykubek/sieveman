package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getCmd.Flags().SortFlags = false
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get <script_name> [output_file*]",
	Short: "Download specified script to given path",
	Long: `This command downloads specified script and saves it to given location.
You can use '-' character as file name to print to stdout.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		scriptName := args[0]
		outFile := scriptName
		// TODO: Add force flag and do not overwrite existing file
		if len(args) > 1 {
			outFile = args[1]
		}

		_, lines, err := c.GetScriptLines(scriptName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var f *os.File
		defer f.Close()
		if outFile == "-" {
			f = os.Stdout
		} else {
			f, err = os.Create(outFile)
			if err != nil {
				panic(err)
			}
		}

		buf := bufio.NewWriter(f)
		defer buf.Flush()

		for _, l := range lines {
			buf.WriteString(l)
		}
	},
}
