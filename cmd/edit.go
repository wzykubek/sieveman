package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <script_name>",
	Short: "Edit specified script remotely in default editor.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		scriptName := args[0]
		content, err := c.GetScriptContent(scriptName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tmpFile, err := os.CreateTemp(os.TempDir(), "sieveman")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		buf := bufio.NewWriter(tmpFile)

		_, err = buf.WriteString(content)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := buf.Flush(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		editor := os.Getenv("EDITOR")
		editorProc := exec.Command(editor, tmpFile.Name())
		editorProc.Stdin = os.Stdin
		editorProc.Stdout = os.Stdout
		editorProc.Stderr = os.Stderr

		if err = editorProc.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// TODO: Do not run if script was not modified
		if err = c.PutScript(tmpFile, scriptName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
