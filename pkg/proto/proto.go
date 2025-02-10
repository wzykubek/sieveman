// Package proto implements ManageSieve protocol specific types especially for server replies.
package proto

type Response interface {
	Type() string
	Code() ResponseCode
	Message() string
}

type ResponseCode interface {
	Type() string
	Message() string
	// Child returns additional hierarchical response code if any. In most cases it is nil.
	Child() ResponseCode
}

type Capabilities struct {
	Implementation string
	SASL           []string
	Sieve          []string
	StartSSL       bool
	MaxRedirects   int
	Notify         []string
	Language       string
	Owner          string
	Version        string
}
