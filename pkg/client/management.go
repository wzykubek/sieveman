package client

import (
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
