package client

import (
	"errors"
	"fmt"
)

// WriteLine is a low level method to write a line to Writer.
// It returns error if any.
func (c *Client) WriteLine(str string) error {
	_, err := fmt.Fprintf(c.Writer, "%s\r\n", str)
	if err != nil {
		return err
	}

	c.Writer.Flush()

	return nil
}

func (c *Client) SendCommand(cmd string) (out []string, err error) {
	if err := c.WriteLine(cmd); err != nil {
		return out, err
	}

	resp, out, err := c.ReadResponse()
	if err != nil {
		return out, err
	}
	logResponse(resp)

	// TODO: Reponse codes should cause errors
	// Almost all response codes are returned with NO/BYE response, but
	// some are returned with OK response: (TAG, WARNINGS, SASL)
	if resp.Name != "OK" {
		return out, errors.New(resp.Message)
	}

	return out, nil
}
