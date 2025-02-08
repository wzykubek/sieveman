// Package client is a library for ManageSieve protocol client. It tries to be RFC 5804 compliant.
package client

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	"go.wzykubek.xyz/sieveman/internal/parsers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

type Client struct {
	tcpConn *net.TCPConn
	tlsConn *tls.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
}

func getConn(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewClient performs DNS-SRV lookup to host (RFC 5804 Section 1.8.1) to obtain IP address
// and port of ManageSieve protocol. In case of failure it reads generic A record and uses
// port as fallback (RFC 5804 Section 1.8.2).
// It returns *Client and error if any.
func NewClient(host string, port int) (*Client, error) {
	var ip net.IP
	var conn *net.TCPConn

	_, records, err := net.LookupSRV("sieve", "tcp", host)
	if err == nil {
		log.Printf("Found %d SRV records", len(records))
		for _, rec := range records {
			ips, err := net.LookupIP(rec.Target)
			if err != nil {
				log.Fatalf("Error resolving IP address: %v", err)
			}

			ip = ips[0]
			port = int(rec.Port)
			log.Printf("Resolved IPv4 from SRV (weight %d) record: %s", rec.Weight, ip)
			log.Printf("Resolved port from SRV (Weight %d) record: %d", rec.Weight, port)

			log.Printf("Connecting to %s:%d\n", host, port)
			conn, err = getConn(&net.TCPAddr{
				IP:   ip,
				Port: port,
			})

			if err == nil {
				break
			} else {
				log.Printf("Failed to connect to %s:%d\n", host, port)
				continue
			}
		}
	} else {
		// no such host
		log.Println("SRV lookup failed, resolving normal IPv4 record")
		ips, err := net.LookupIP(host)
		if err != nil {
			log.Fatalf("Error resolving IP address: %v", err)
		}

		ip = ips[0]
		log.Printf("Resolved IPv4: %s", ip)
		log.Printf("Using fallback port: %d", port)

		log.Printf("Connecting to %s:%d\n", host, port)
		conn, err = getConn(&net.TCPAddr{
			IP:   ip,
			Port: port,
		})
	}

	c := &Client{
		tcpConn: conn,
		reader:  bufio.NewReader(conn),
		writer:  bufio.NewWriter(conn),
	}

	r, _, err := c.readResponse()
	if err != nil {
		return nil, err
	}
	logResponse(r)

	if _, ok := r.(proto.Ok); ok {
		c.upgradeConn()
	}

	return c, nil
}

// Close closes both TCP and TLS connections.
// It returns error if any.
func (c *Client) Close() error {
	var err error
	if c.tcpConn != nil {
		err = c.tcpConn.Close()
	}
	if c.tlsConn != nil {
		err = c.tlsConn.Close()
	}

	return err
}

func (c *Client) upgradeConn() {
	log.Println("Trying to perform connection upgrade")

	c.writeCommand("STARTTLS")

	r, _, err := c.readResponse()
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}
	logResponse(r)

	if _, ok := r.(proto.Ok); ok {
		log.Println("Starting connection upgrade")
		c.tlsConn = tls.Client(c.tcpConn, &tls.Config{
			InsecureSkipVerify: true, // TODO
		})
		err := c.tlsConn.Handshake()
		if err != nil {
			log.Fatalf("TLS Handshake failed: %v", err)
		}

		c.reader = bufio.NewReader(c.tlsConn)
		c.writer = bufio.NewWriter(c.tlsConn)

		r, _, err = c.readResponse()
		if err != nil {
			log.Fatalf("Error reading response after TLS upgrade: %v", err)
		}
		logResponse(r)
	} else {
		log.Fatalf("Unexpected response to STARTTLS command: %v", r)
	}
}

func (c *Client) writeCommand(cmd string) error {
	_, err := fmt.Fprintf(c.writer, "%s\r\n", cmd)
	if err != nil {
		return err
	}

	c.writer.Flush()

	return nil
}

func (c *Client) readResponse() (proto.Response, []proto.Message, error) {
	var messages []proto.Message

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, messages, err
		}

		trimedLine := strings.TrimSpace(line)
		resp := parsers.ParseResponse(trimedLine)
		if resp != nil {
			return resp, messages, nil
		} else {
			msg := parsers.ParseMessage(trimedLine)
			messages = append(messages, msg)
		}
	}
}
