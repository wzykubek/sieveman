package cmd

import (
	"fmt"
	"io"
	"os"

	"go.wzykubek.xyz/sieveman/pkg/client"

	"github.com/spf13/cobra"
)

var host string
var port int
var username string
var password string
var verbose bool

func init() {
	cmdRoot.PersistentFlags().StringVarP(&host, "host", "H", "", "host")
	cmdRoot.PersistentFlags().IntVarP(&port, "port", "P", 4190, "port")
	cmdRoot.PersistentFlags().StringVarP(&username, "username", "u", "", "username [email address]")
	cmdRoot.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
	cmdRoot.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")

	cmdRoot.Flags().SortFlags = false
	cmdRoot.PersistentFlags().SortFlags = false
}

var cmdRoot = &cobra.Command{
	Use:   "sieveman",
	Short: "Sieve manager",
	Long:  "Universal ManageSieve protocol client",
	Run: func(cmd *cobra.Command, args []string) {
		if !verbose {
			client.Logger.SetOutput(io.Discard)
		}

		c, err := client.NewClient(host, port)
		if err != nil {
			panic(err)
		}
		defer c.Close()

		// TODO: Interactive part
	},
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
