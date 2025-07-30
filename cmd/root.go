package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.wzykubek.xyz/sieveman/pkg/client"

	"github.com/chzyer/readline"
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
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		c.Close()
	},
	Run: func(cmd *cobra.Command, args []string) {
		rl, err := readline.NewEx(&readline.Config{
			Prompt: "sieveman> ",
			AutoComplete: readline.NewPrefixCompleter(
				readline.PcItem("exit"),
				readline.PcItem("bye"),
				readline.PcItem("help"),
				readline.PcItem("activate"),
				readline.PcItem("deactivate"),
				readline.PcItem("edit"),
				readline.PcItem("get"),
				readline.PcItem("ls"),
				readline.PcItem("mv"),
				readline.PcItem("put"),
				readline.PcItem("rm"),
				readline.PcItem("version"),
			),
		})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					fmt.Println("^C")
					continue
				}
				if err == io.EOF {
					fmt.Println("EOF")
					break
				}
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}

			line = strings.TrimSpace(line)
			switch {
			case line == "exit" || line == "bye":
				return
			case line == "help":
				fmt.Println(`Available commands:
activate <script_name>              Activate a script
bye                                 Exit the shell
deactivate                          Deactivate all scripts
edit <script_name>                  Edit a script
exit                                Exit the shell
get <script_name> [file_name]       Get a script
help                                Display this help message
ls                                  List all scripts
mv <old_name> <new_name>            Rename a script
put <file_name> [script_name]       Put a script
rm <script_name>                    Remove a script
version                             Display the program version`)
			case strings.HasPrefix(line, "activate"):
				args := strings.Fields(line)
				if len(args) == 1 {
					fmt.Println("You must provide a script name.")
					continue
				}
				activateCmd.Run(activateCmd, args[1:])
			case line == "deactivate":
				deactivateCmd.Run(deactivateCmd, []string{})
			case strings.HasPrefix(line, "edit"):
				args := strings.Fields(line)
				if len(args) == 1 {
					fmt.Println("You must provide a script name.")
					continue
				}
				editCmd.Run(editCmd, args[1:])
			case strings.HasPrefix(line, "get"):
				args := strings.Fields(line)
				if len(args) == 1 {
					fmt.Println("You must provide a script name.")
					continue
				}
				getCmd.Run(getCmd, args[1:])
			case line == "ls":
				lsCmd.Run(lsCmd, []string{})
			case strings.HasPrefix(line, "mv"):
				args := strings.Fields(line)
				mvCmd.Run(mvCmd, args[1:])
				if len(args) == 1 {
					fmt.Println("You must provide a script name.")
					continue
				}
				if len(args) == 2 {
					fmt.Println("You must provide new name.")
					continue
				}
			case strings.HasPrefix(line, "put"):
				args := strings.Fields(line)
				if len(args) == 1 {
					fmt.Println("You must provide a file name.")
					continue
				}
				putCmd.Run(putCmd, args[1:])
			case strings.HasPrefix(line, "rm"):
				args := strings.Fields(line)
				if len(args) == 1 {
					fmt.Println("You must provide a script name.")
					continue
				}
				rmCmd.Run(rmCmd, args[1:])
			case line == "version":
				versionCmd.Run(versionCmd, []string{})
			default:
				fmt.Println("Invalid command.")
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
