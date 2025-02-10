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
		err := c.UpgradeConn()
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

// Close closes connection.
// It returns error if any.
func (c *Client) Close() error {
	var err error
	err = c.Conn.Close()

	return err
}

// UpgradeConn upgrades existing plain TCP connection of client to TLS using StartTLS.
// It returns error if any.
func (c *Client) UpgradeConn() error {
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
func (c *Client) ReadResponse() (proto.Response, []string, error) {
	var messages []string

	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return nil, messages, err
		}

		trimedLine := strings.TrimSpace(line)
		resp := parsers.ParseResponse(trimedLine)
		if resp != nil {
			return resp, messages, nil
		} else {
			messages = append(messages, trimedLine)
		}
	}
}
