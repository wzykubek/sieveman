package client

import (
	"bufio"
	"crypto/tls"
	"errors"
	"net"
)

func resolveIPv4(host string) (net.IP, error) {
	Logger.Println("Resolving IPv4")

	ips, err := net.LookupIP(host)
	if err != nil {
		Logger.Printf("-> Error resolving IP address: %v", err)
		return nil, err
	}

	Logger.Printf("-> Resolved %s IP", ips[0])

	return ips[0], nil
}

func dialTCP(ip net.IP, port int) (*net.TCPConn, error) {
	Logger.Printf("Connecting to %s:%d\n", ip, port)
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		Logger.Printf("-> Failed to connect to %s:%d\n", ip, port)
		return nil, err
	}

	return conn, nil
}

func attemptConn(host string, port int) (*net.TCPConn, error) {
	ip, err := resolveIPv4(host)
	if err != nil {
		return nil, err
	}

	conn, err := dialTCP(ip, port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Close closes existing connection with server.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// GetTCPConn performs DNS-SRV lookup to host (RFC 5804 Section 1.8.1) to obtain IP address
// and port of ManageSieve protocol. In case of failure it reads generic A record and uses
// port as fallback (RFC 5804 Section 1.8.2).
func GetTCPConn(host string, port int) (conn *net.TCPConn, err error) {
	Logger.Printf("Trying to lookup SRV records of %s host\n", host)

	_, records, err := net.LookupSRV("sieve", "tcp", host)
	if err != nil {
		Logger.Println("-> SRV lookup failed")
		conn, err = attemptConn(host, port)
		if err != nil {
			return nil, err
		}
	}

	Logger.Printf("-> Found SRV records")

	for _, rec := range records {
		conn, err = attemptConn(rec.Target, int(rec.Port))
		if err == nil {
			break
		}
	}

	if conn == nil {
		Logger.Println("-> SRV records did not yield a connection")
		conn, err = attemptConn(host, port)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// GetTLSConn connects to server using TLS.
func GetTLSConn(plainConn net.Conn, serverName string) (conn *tls.Conn, err error) {
	conn = tls.Client(plainConn, &tls.Config{
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	})

	err = conn.Handshake()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// UpgradeConn is checking if server supports StartTLS and performs connection upgrade.
func (c *Client) UpgradeConn(serverName string) error {
	Logger.Println("Checking if server supports StartTLS")

	if !c.capabilities.StartSSL {
		Logger.Println("-> Server does not support StartTLS")
		Logger.Println("Aborting connection upgrade")

		return errors.New("Server does not support StartTLS")
	}

	Logger.Println("-> Server supports StartTLS")
	Logger.Println("Trying to start TLS negotiation")

	if _, err := c.SendCommand("STARTTLS"); err != nil {
		return err
	}

	Logger.Println("Starting TLS connection")

	tlsConn, err := GetTLSConn(c.Conn, serverName)
	if err != nil {
		return err
	}

	c.Conn = tlsConn
	c.Reader = bufio.NewReader(tlsConn)
	c.Writer = bufio.NewWriter(tlsConn)

	r, _, err := c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	if r.Name != "OK" {
		return errors.New(r.Message)
	}

	return nil
}
