// Package proto implements ManageSieve protocol specific types especially for server replies.
package proto

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
