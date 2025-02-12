package client

import (
	"encoding/base64"
	"fmt"

	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func encCredentials(login string, password string) string {
	data := []byte("\x00" + login + "\x00" + password)
	return base64.StdEncoding.EncodeToString(data)
}

// AuthPLAIN uses PLAIN SASL to authenticate with server if that method is supported.
// It returns parsed response and error if any.
func (c *Client) AuthPLAIN(login string, password string) (proto.Response, error) {
	Logger.Println("Checking server capabilities")
	capabilities, err := c.GetCapabilities()
	if err != nil {
		return nil, err
	}

	var plainCap bool
	for _, v := range capabilities.SASL {
		if v == "PLAIN" {
			plainCap = true
		}
	}

	if !plainCap {
		Logger.Println("-> Server does not support PLAIN authentication")
		Logger.Println("Aborting authentication")
		return nil, nil
	}
	Logger.Println("-> Server supports PLAIN authentication")
	Logger.Println("Trying to authenticate")

	encCred := encCredentials(login, password)
	c.Write(fmt.Sprintf(`AUTHENTICATE "PLAIN" "%s"`, encCred))
	r, _, err := c.ReadResponse()
	if err != nil {
		return nil, err
	}
	logResponse(r)

	return r, nil
}
