package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sieveman %s\n", version)
	},
}

func init() {
	cmdRoot.AddCommand(cmdVersion)
}
