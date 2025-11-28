package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show app version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sieveman %s\n", version)
	},
	// Bypass hooks and avoid unnecessary connections
	PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
