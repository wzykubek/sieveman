package cmd

import (
	"fmt"
	"io"
	"os"

	"go.wzykubek.xyz/sieveman/pkg/client"

	"github.com/spf13/cobra"
)

var c *client.Client

var (
	host     string
	port     int
	username string
	password string
	verbose  bool
)

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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !verbose {
			client.Logger.SetOutput(io.Discard)
		}

		var err error
		c, err = client.NewClient(host, port)
		if err != nil {
			os.Exit(1)
		}

		err = c.AuthPLAIN(username, password)
		if err != nil {
			fmt.Println("Authentication failed!")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer c.Close()

		// TODO: Interactive part
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
