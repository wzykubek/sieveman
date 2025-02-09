package client

import (
	"fmt"
	"log"
	"os"

	"go.wzykubek.xyz/sieveman/pkg/proto"
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

func logResponse(r proto.Response) {
	var c string
	switch r.(type) {
	case proto.Ok:
		c = yellow
	case proto.No:
		c = red
	case proto.Bye:
		c = blue
	default:
		c = reset
	}

	Logger.Printf("%s => %s", colorize(r.Type(), c), r.Message())
}
