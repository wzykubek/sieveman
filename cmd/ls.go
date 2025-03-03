package cmd

import (
	"fmt"
	"os"

	"go.wzykubek.xyz/sieveman/pkg/proto"

	"github.com/spf13/cobra"
)

var noIndicator bool

func init() {
	lsCmd.Flags().BoolVar(&noIndicator, "no-indicator", false, "do not show active indicator")
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available scripts",
	Long:  "This command lists all available scripts and shows indicator after activated ones.",
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		r, scripts, err := c.ListScripts()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if _, ok := r.(proto.Ok); ok {
			for _, v := range scripts {
				var ind rune
				if v.Active && !noIndicator {
					ind = '*'
				}
				fmt.Printf("%s%c\n", v.Name, ind)
			}
		}

	},
}
