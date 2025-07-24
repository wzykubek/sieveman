package client

import (
	"errors"
	"fmt"

	"go.wzykubek.xyz/sieveman/internal/helpers"
)

// AuthPLAIN uses PLAIN SASL to authenticate with server if that method is supported.
// It returns parsed response and error if any.
func (c *Client) AuthPLAIN(login string, password string) error {
	Logger.Println("Checking server capabilities")
	capabilities, err := c.ReadCapabilities()
	if err != nil {
		return err
	}
	r, _, err := c.ReadResponse()
	if err != nil {
		return err
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
		return nil
	}
	Logger.Println("-> Server supports PLAIN authentication")
	Logger.Println("Trying to authenticate")

	encCred := helpers.EncCredentials(login, password)
	c.WriteLine(fmt.Sprintf(`AUTHENTICATE "PLAIN" "%s"`, encCred))
	r, _, err = c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	if r.Name != "OK" {
		return errors.New(r.Message)
	}

	return nil
}
