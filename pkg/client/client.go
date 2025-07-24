// Package client is a library for ManageSieve protocol client. It tries to be RFC 5804 compliant.
// It offers both low and top level functions and methods to be flexible and easy to use at once.
package client

import (
	"bufio"
	"errors"
	"net"
)

type Client struct {
	Conn         net.Conn
	Reader       *bufio.Reader
	Writer       *bufio.Writer
	capabilities Capabilities
}

// NewClient is a top level function to create new *Client. It handles all necessary checks,
// connects to server over plain TCP connection and performs connection upgrade to TLS.
// It returns *Client and error if any.
func NewClient(host string, port int) (*Client, error) {
	tcpConn, err := GetTCPConn(host, port)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Conn:   tcpConn,
		Reader: bufio.NewReader(tcpConn),
		Writer: bufio.NewWriter(tcpConn),
	}

	cap, err := c.ReadCapabilities()
	if err != nil {
		return nil, err
	}
	c.capabilities = cap

	r, _, err := c.ReadResponse()
	if err != nil {
		return nil, err
	}
	logResponse(r)

	if r.Name != "OK" {
		return nil, errors.New(r.Message)
	}

	if err := c.UpgradeConn(); err != nil {
		return c, err
	}

	cap, err = c.GetCapabilities()
	if err != nil {
		return nil, err
	}
	c.capabilities = cap

	return c, nil
}
