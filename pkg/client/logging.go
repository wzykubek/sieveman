package client

import (
	"fmt"
	"log"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func init() {
	log.SetFlags(log.Ltime)
}

func colorize(s string, c string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}

func logResponse(r Response) {
	var c string
	switch r.(type) {
	case ResponseOK:
		c = Yellow
	case ResponseNO:
		c = Red
	case ResponseBYE:
		c = Blue
	default:
		c = Reset
	}

	log.Printf("%s => %s", colorize(r.Type(), c), r.Message())
}
