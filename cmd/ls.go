package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	noIndicator bool
	onlyActive  bool
)

func init() {
	lsCmd.Flags().BoolVar(&noIndicator, "no-indicator", false, "do not show active indicator")
	lsCmd.Flags().BoolVar(&onlyActive, "active", false, "show only active script")

	lsCmd.Flags().SortFlags = false
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available scripts",
	Long:  "This command lists all available scripts and shows indicator next to the active script if any.",
	Run: func(cmd *cobra.Command, args []string) {
		scripts, err := c.GetScriptList()
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		for _, v := range scripts {
			var ind rune
			if v.Active && !noIndicator {
				ind = '*'
			}

			if onlyActive && !v.Active {
				continue
			}

			fmt.Printf("%s%c\n", v.Name, ind)
		}
	},
}
