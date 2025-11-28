package cmd

import (
	"bufio"
	"log"
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
	Use:   "get <script_name> [output_file]",
	Short: "Download specified script to given path",
	Long: `This command downloads specified script and saves it to given location.
You can use '-' character as a file name to print script content to stdout instead.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		outFilename := scriptName
		if len(args) > 1 {
			outFilename = args[1]
		}

		content, err := c.GetScriptContent(scriptName)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		var f *os.File
		defer f.Close()

		if outFilename == "-" {
			f = os.Stdout
		} else {
			if _, err := os.Stat(outFilename); err == nil {
				if !forceWrite {
					log.Fatalf("Error: File %s exists\n", outFilename)
				}
			}

			f, err = os.Create(outFilename)
			if err != nil {
				log.Fatalf("Error: %s\n", err)
			}
		}

		buf := bufio.NewWriter(f)
		defer buf.Flush()

		if _, err = buf.WriteString(content); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}
