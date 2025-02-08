package client

import (
	"fmt"
	"log"

	"go.wzykubek.xyz/sieveman/pkg/proto"
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

func logResponse(r proto.Response) {
	var c string
	switch r.(type) {
	case proto.ResponseOK:
		c = yellow
	case proto.ResponseNO:
		c = red
	case proto.ResponseBYE:
		c = blue
	default:
		c = reset
	}

	log.Printf("%s => %s", colorize(r.Type(), c), r.Message())
}
