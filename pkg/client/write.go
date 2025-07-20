package client

import "fmt"

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
