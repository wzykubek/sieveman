package client

import (
	"errors"
	"fmt"

	"go.wzykubek.xyz/sieveman/internal/helpers"
)

// AuthPLAIN uses PLAIN SASL to authenticate with server if that method is supported.
// It returns parsed response and error if any.
func (c *Client) AuthPLAIN(login string, password string) error {
	Logger.Println("Checking if server supports PLAIN authentication")

	var plainCap bool
	for _, v := range c.capabilities.SASL {
		if v == "PLAIN" {
			plainCap = true
		}
	}

	if !plainCap {
		Logger.Println("-> Server does not support PLAIN authentication")
		Logger.Println("Aborting authentication")
		return errors.New("Server does not support PLAIN authentication")
	}
	Logger.Println("-> Server supports PLAIN authentication")
	Logger.Println("Trying to authenticate")

	encCred := helpers.EncCredentials(login, password)
	cmd := fmt.Sprintf(`AUTHENTICATE "PLAIN" "%s"`, encCred)
	_, err := c.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}
