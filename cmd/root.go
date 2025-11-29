package cmd

import (
	"fmt"
	"io"
	"log"
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

func runInlineCmd(c *cobra.Command, args []string) {
	if err := c.RunE(c, args); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

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
	Short: "ManageSieve client",
	Long: `This tool allows you to list, download, edit and create Sieve filters on specified e-mail server.

It works in two modes: command line and interactive shell.
You need to pass at least --host, --username and --password to use any further command (exclude: help, completion and version).

If you do not specify the command, you will enter interactive mode instead.`,
	Example: `sieveman -H "imap.example.com" -u "jdoe@example.com" -p "$(qpg -qd password.txt.asc)" [command]`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		if !verbose {
			client.Logger.SetOutput(io.Discard)
		}

		// There is an option to mark flag as required using
		// cobra builtin method. However, that method cannot
		// be used in this scenario, because it is carried
		// into effect after PersistentPreRun execution.
		if host == "" {
			log.Fatalln("Error: host is not specified")
		}
		if username == "" {
			log.Fatalln("Error: username is not specified")
		}
		if password == "" {
			log.Fatalln("Error: password is not specified")
		}

		var err error
		c, err = client.NewClient(host, port)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

		err = c.AuthPLAIN(username, password)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
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
			log.Fatalf("Error: %s\n", err)
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
				log.Fatalf("Error: %s\n", err)
			}

			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) == 0 {
				continue
			}

			switch fields[0] {
			case "exit", "bye":
				return

			case "help":
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

			case "activate":
				if len(fields) == 1 {
					fmt.Println("Error: script name not specified")
					continue
				}

				runInlineCmd(activateCmd, fields[1:])

			case "deactivate":
				runInlineCmd(deactivateCmd, []string{})

			case "edit":
				if len(fields) == 1 {
					fmt.Println("Error: script name not specified")
					continue
				}

				runInlineCmd(editCmd, fields[1:])

			case "get":
				if len(fields) == 1 {
					fmt.Println("Error: script name not specified")
					continue
				}

				runInlineCmd(getCmd, fields[1:])

			case "ls":
				runInlineCmd(lsCmd, []string{})

			case "mv":
				if len(fields) == 1 {
					fmt.Println("Error: script name not specified")
					continue
				}
				if len(fields) == 2 {
					fmt.Println("Error: new name not specified")
					continue
				}

				runInlineCmd(mvCmd, fields[1:])

			case "put":
				if len(fields) == 1 {
					fmt.Println("Error: file name not specified")
					continue
				}

				runInlineCmd(putCmd, fields[1:])

			case "rm":
				if len(fields) == 1 {
					fmt.Println("Error: script name not specified")
					continue
				}

				runInlineCmd(rmCmd, fields[1:])

			case "version":
				versionCmd.Run(versionCmd, []string{})

			default:
				fmt.Println("Error: invalid command")
			}
		}
	},
}

// Root command export for use with man page generator.
func Root() *cobra.Command { return rootCmd }

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
