package client

import (
	"crypto/tls"
	"net"
)

func resolveIPv4(host string) (net.IP, error) {
	Logger.Println("Resolving IPv4")

	ips, err := net.LookupIP(host)
	if err != nil {
		Logger.Printf("-> Error resolving IP address: %v", err)
		return nil, err
	} else {
		Logger.Printf("-> Resolved %s IP", ips[0])
	}

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

// GetTCPConn performs DNS-SRV lookup to host (RFC 5804 Section 1.8.1) to obtain IP address
// and port of ManageSieve protocol. In case of failure it reads generic A record and uses
// port as fallback (RFC 5804 Section 1.8.2).
// It returns valid TCP connection and error if any.
func GetTCPConn(host string, port int) (*net.TCPConn, error) {
	var conn *net.TCPConn

	Logger.Printf("Trying to lookup SRV records of %s host\n", host)
	_, records, err := net.LookupSRV("sieve", "tcp", host)
	if err == nil {
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
	} else {
		Logger.Println("-> SRV lookup failed")
		conn, err = attemptConn(host, port)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// GetTLSConn connects to server using TLS.
// It returns TLS connection and error if any.
func GetTLSConn(plainConn net.Conn) (*tls.Conn, error) {
	var tlsConn *tls.Conn

	tlsConn = tls.Client(plainConn, &tls.Config{
		InsecureSkipVerify: true, // TODO
	})
	err := tlsConn.Handshake()
	if err != nil {
		return nil, err
	}
	return tlsConn, nil
}
