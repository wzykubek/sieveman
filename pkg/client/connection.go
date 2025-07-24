package client

import (
	"bufio"
	"crypto/tls"
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

// Close closes connection.
// It returns error if any.
func (c *Client) Close() error {
	return c.Conn.Close()
}

// GetTCPConn performs DNS-SRV lookup to host (RFC 5804 Section 1.8.1) to obtain IP address
// and port of ManageSieve protocol. In case of failure it reads generic A record and uses
// port as fallback (RFC 5804 Section 1.8.2).
// It returns valid TCP connection and error if any.
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
// It returns TLS connection and error if any.
func GetTLSConn(plainConn net.Conn) (conn *tls.Conn, err error) {
	conn = tls.Client(plainConn, &tls.Config{
		InsecureSkipVerify: true, // TODO
	})

	err = conn.Handshake()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// UpgradeConn upgrades existing plain TCP connection of client to TLS using StartTLS.
// It returns error if any.
func (c *Client) UpgradeConn() error {
	if !c.capabilities.StartSSL {
		Logger.Println("-> Server does not support StartTLS")
		Logger.Println("Aborting connection upgrade")

		return nil // TODO: Return error
	}

	Logger.Println("-> Server supports StartTLS")
	Logger.Println("Trying to start TLS negotiation")

	c.WriteLine("STARTTLS")
	r, _, err := c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	var tlsConn *tls.Conn
	// TODO: Handle NO
	if r.Name == "OK" {
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
