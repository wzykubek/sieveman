package client

import (
	"errors"
	"fmt"
	"strings"

	"go.wzykubek.xyz/sieveman/internal/parsers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func (c *Client) ListScripts() (_ proto.Response, s []proto.Script, err error) {
	Logger.Println("Trying to obtain script list")

	c.Write("LISTSCRIPTS")
	r, m, err := c.ReadResponse()
	if err != nil {
		return nil, s, err
	}
	logResponse(r)

	s = parsers.ParseScriptList(m)
	return r, s, nil
}

func (c *Client) GetScript(name string) (_ proto.Response, script string, err error) {
	Logger.Println("Trying to get script content")

	c.Write(fmt.Sprintf("GETSCRIPT \"%s\"", name))
	r, m, err := c.ReadResponse()
	if err != nil {
		return nil, script, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		script = strings.TrimSpace(strings.Join(m[1:], "\n")) // Ignore length
	} else {
		err = errors.New(r.Message())
	}

	return r, script, err
}
