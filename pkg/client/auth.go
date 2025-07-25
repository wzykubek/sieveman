package client

import (
	"encoding/base64"
	"errors"
	"fmt"
)

func encBase64Cred(login string, password string) string {
	data := []byte("\x00" + login + "\x00" + password)
	return base64.StdEncoding.EncodeToString(data)
}

// AuthPLAIN uses PLAIN SASL to authenticate with server if that method is supported.
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

	cred := encBase64Cred(login, password)
	cmd := fmt.Sprintf(`AUTHENTICATE "PLAIN" "%s"`, cred)
	if _, err := c.SendCommand(cmd); err != nil {
		return err
	}

	return nil
}
