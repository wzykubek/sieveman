package client

type Response struct {
	Name    string
	Code    ResponseCode
	Message string
}

type ResponseCode struct {
	Name    string
	Message string
}

type Script struct {
	Name   string
	Active bool
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
