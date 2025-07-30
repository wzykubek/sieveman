package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var forceWrite bool

func init() {
	getCmd.Flags().SortFlags = false
	getCmd.Flags().BoolVarP(&forceWrite, "force", "f", false, "force overwrite existing file")
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get <script_name> [output_file*]",
	Short: "Download specified script to given path",
	Long: `This command downloads specified script and saves it to given location.
You can use '-' character as file name to print to stdout.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		outFilename := scriptName
		if len(args) > 1 {
			outFilename = args[1]
		}

		content, err := c.GetScriptContent(scriptName)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		var f *os.File
		defer f.Close()

		if outFilename == "-" {
			f = os.Stdout
		} else {
			if _, err := os.Stat(outFilename); err == nil {
				if !forceWrite {
					fmt.Printf("Error: File %s exists\n", outFilename)
					os.Exit(1)
				}
			}

			f, err = os.Create(outFilename)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
		}

		buf := bufio.NewWriter(f)
		defer buf.Flush()

		if _, err = buf.WriteString(content); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
}
