package client

import (
	"errors"
	"fmt"
	"strings"

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

func (c *Client) GetScript(name string) (script string, err error) {
	Logger.Println("Trying to get script content")

	c.Write(fmt.Sprintf("GETSCRIPT \"%s\"", name))
	r, m, err := c.ReadResponse()
	if err != nil {
		return script, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		script = strings.TrimSpace(strings.Join(m[1:], "\n")) // Ignore length
	} else {
		err = errors.New(r.Message())
	}

	return script, err
}
