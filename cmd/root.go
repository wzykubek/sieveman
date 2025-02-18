package cmd

import (
	"fmt"
	"io"
	"os"

	"go.wzykubek.xyz/sieveman/pkg/client"
	"go.wzykubek.xyz/sieveman/pkg/proto"

	"github.com/spf13/cobra"
)

var host string
var port int
var username string
var password string
var verbose bool

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "host")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "P", 4190, "port")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username [email address]")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
}

var rootCmd = &cobra.Command{
	Use:   "sieveman",
	Short: "Sieve manager",
	Long:  "Universal ManageSieve protocol client",
	Run: func(cmd *cobra.Command, args []string) {
		if !verbose {
			client.Logger.SetOutput(io.Discard)
		}

		c, err := client.NewClient(host, port)
		if err != nil {
			os.Exit(1)
		}
		defer c.Close()

		r, err := c.AuthPLAIN(username, password)
		if err != nil {
			os.Exit(1)
		}

		if _, ok := r.(proto.Ok); ok {
			fmt.Println("Authentication successful!")
		} else {
			fmt.Println("Authentication failded!")
		}

		// TODO: Interactive part
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
