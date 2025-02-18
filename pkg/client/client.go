// Package client is a library for ManageSieve protocol client. It tries to be RFC 5804 compliant.
// It offers both low and top level functions and methods to be flexible and easy to use at once.
package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"go.wzykubek.xyz/sieveman/internal/parsers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

type Client struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
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

	r, _, err := c.ReadResponse()
	if err != nil {
		return nil, err
	}
	logResponse(r)

	if _, ok := r.(proto.Ok); ok {
		if err := c.UpgradeConn(); err != nil {
			return c, err
		}
	}

	return c, nil
}

// Close closes connection.
// It returns error if any.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// UpgradeConn upgrades existing plain TCP connection of client to TLS using StartTLS.
// It returns error if any.
func (c *Client) UpgradeConn() error {
	Logger.Println("Checking server capabilities")
	capabilities, err := c.GetCapabilities()
	if err != nil {
		return err
	}

	if !capabilities.StartSSL {
		Logger.Println("-> Server does not support StartTLS")
		Logger.Println("Aborting connection upgrade")
		return nil
	}

	Logger.Println("-> Server supports StartTLS")
	Logger.Println("Trying to start TLS negotiation")
	c.Write("STARTTLS")

	r, _, err := c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	var tlsConn *tls.Conn
	if _, ok := r.(proto.Ok); ok {
		Logger.Println("Starting TLS connection")
		tlsConn, err = GetTLSConn(c.Conn)
		if err != nil {
			return err
		}
	}

	c.Conn = tlsConn
	c.Reader = bufio.NewReader(tlsConn)
	c.Writer = bufio.NewWriter(tlsConn)

	r, _, err = c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	return nil
}

// Write is a low level method to write a line to Writer.
// It returns error if any.
func (c *Client) Write(str string) error {
	_, err := fmt.Fprintf(c.Writer, "%s\r\n", str)
	if err != nil {
		return err
	}

	c.Writer.Flush()

	return nil
}

// ReadResponse is a low level method to read and parse response from server.
// It returns parsed response, slice of messages and error if any.
func (c *Client) ReadResponse() (_ proto.Response, messages []string, err error) {
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return nil, messages, err
		}

		trimedLine := strings.TrimSpace(line)
		if resp := parsers.ParseResponse(trimedLine); resp != nil {
			return resp, messages, nil
		} else {
			messages = append(messages, trimedLine)
		}
	}
}

// GetCapabilities reads server capabilities.
// It returns proto.Capabilities and error if any.
func (c *Client) GetCapabilities() (capb proto.Capabilities, err error) {
	err = c.Write("CAPABILITY")
	if err != nil {
		fmt.Println(err)
	}

	r, messages, err := c.ReadResponse()
	if err != nil {
		return capb, err
	}

	if _, ok := r.(proto.Ok); ok {
		capb = parsers.ParseCapabilities(messages)
	}

	return capb, nil
}
