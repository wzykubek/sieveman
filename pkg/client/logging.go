package client

import (
	"fmt"
	"log"
	"os"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

var Logger log.Logger = *log.New(os.Stderr, "", 3)

func colorize(s string, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, reset)
}

func logResponse(r Response) {
	var c string
	switch r.Name {
	case "OK":
		c = yellow
	case "NO":
		c = red
	case "BYE":
		c = blue
	default:
		c = reset
	}

	Logger.Printf("%s => %s", colorize(r.Name, c), r.Message)
}
