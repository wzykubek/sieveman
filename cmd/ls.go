package cmd

import (
	"fmt"
	"os"

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
	Long:  "This command lists all available scripts and shows activation indicator.",
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		scripts, err := c.GetScriptList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
