package client

import (
	"errors"
	"fmt"
	"strconv"

	"go.wzykubek.xyz/sieveman/internal/parsers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func (c *Client) ListScripts() (s []proto.Script, err error) {
	Logger.Println("Trying to obtain script list")

	c.Write("LISTSCRIPTS")
	r, m, err := c.ReadResponse()
	if err != nil {
		return s, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		s = parsers.ParseScriptList(m)
	} else {
		err = errors.New(r.Message())
	}

	return s, nil
}

func (c *Client) GetScriptLines(name string) (size int64, script []string, err error) {
	Logger.Println("Trying to get script content")

	c.Write(fmt.Sprintf("GETSCRIPT \"%s\"", name))
	r, m, err := c.ReadResponse()
	if err != nil {
		return size, script, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		size, _ = strconv.ParseInt(m[0], 10, 64) // TODO: Handle error
		script = m[1 : len(m)-1]
	} else {
		err = errors.New(r.Message())
	}

	return size, script, err
}
