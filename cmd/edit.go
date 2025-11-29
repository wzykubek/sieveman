package cmd

import (
	"bufio"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit <script_name>",
	Short: "Edit specified script remotely in default editor",
	Long:  "This command combines get and put together. Use it to directly edit remote scripts in your default editor specified with EDITOR environment variable.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptName := args[0]
		content, err := c.GetScriptContent(scriptName)
		if err != nil {
			return err
		}

		tmpFile, err := os.CreateTemp(os.TempDir(), "sieveman")
		if err != nil {
			return err
		}

		buf := bufio.NewWriter(tmpFile)

		_, err = buf.WriteString(content)
		if err != nil {
			return err
		}

		if err := buf.Flush(); err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		editorProc := exec.Command(editor, tmpFile.Name())
		editorProc.Stdin = os.Stdin
		editorProc.Stdout = os.Stdout
		editorProc.Stderr = os.Stderr

		if err = editorProc.Run(); err != nil {
			return err
		}

		if err = c.PutScript(tmpFile, scriptName); err != nil {
			return err
		}

		return nil
	},
}
