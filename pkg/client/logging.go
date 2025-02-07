package client

import (
	"fmt"
	"log"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

func init() {
	log.SetFlags(log.Ltime)
}

func colorize(s string, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, reset)
}

func logResponse(r Response) {
	var c string
	switch r.(type) {
	case ResponseOK:
		c = yellow
	case ResponseNO:
		c = red
	case ResponseBYE:
		c = blue
	default:
		c = reset
	}

	log.Printf("%s => %s", colorize(r.Type(), c), r.Message())
}
