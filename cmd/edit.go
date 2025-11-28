package cmd

import (
	"bufio"
	"log"
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
	Run: func(cmd *cobra.Command, args []string) {
		scriptName := args[0]
		content, err := c.GetScriptContent(scriptName)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		tmpFile, err := os.CreateTemp(os.TempDir(), "sieveman")
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		buf := bufio.NewWriter(tmpFile)

		_, err = buf.WriteString(content)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		if err := buf.Flush(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		editor := os.Getenv("EDITOR")
		editorProc := exec.Command(editor, tmpFile.Name())
		editorProc.Stdin = os.Stdin
		editorProc.Stdout = os.Stdout
		editorProc.Stderr = os.Stderr

		if err = editorProc.Run(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		if err = c.PutScript(tmpFile, scriptName); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}
