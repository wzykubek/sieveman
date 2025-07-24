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

func (c *Client) SendCommand(cmd string) (outputs []string, err error) {
	if err := c.WriteLine(cmd); err != nil {
		return outputs, err
	}

	resp, outputs, err := c.ReadResponse()
	if err != nil {
		return outputs, err
	}
	logResponse(resp)

	if resp.Name != "OK" {
		return outputs, errors.New(resp.Message)
	}

	return outputs, nil
}
