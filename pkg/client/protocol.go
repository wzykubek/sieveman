package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

func readNBytes(reader *bufio.Reader, byteCount int) (content string, err error) {
	buffer := make([]byte, byteCount)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		return content, err
	}

	return string(buffer), nil
}

// ReadResponse is a low level method to read and parse response from server.
// It returns parsed response, slice of outputs and error if any.
// `outputs` needs to be handled depending on the command type.
// Examples:
// - In case of reading response from LISTSCRIPTS command then each output
// is a script line which need to be parsed.
// - In case of reading response from GETSCRIPT ocmmand then it will be only one item
// with whole script content.
func (c *Client) ReadResponse() (response Response, outputs []string, err error) {
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return Response{}, outputs, err
		}
		line = strings.TrimSpace(line)
		if line == "" && c.Reader.Buffered() != 0 {
			continue
		}

		p := Parser{input: line, position: 0}
		bytes, err := p.parseBytes()
		if err != nil {
			return Response{}, outputs, err
		}

		if bytes != 0 {
			out, err := readNBytes(c.Reader, bytes)
			if err != nil {
				return Response{}, outputs, err
			}

			outputs = append(outputs, out)
			continue
		}

		if strings.HasPrefix(line, `"`) {
			if err := parseCapability(&c.capabilities, line); err != nil {
				outputs = append(outputs, line)
			}
			continue
		}

		response, bytes, err = parseInlineResponse(line)
		if err != nil {
			return Response{}, outputs, err
		}

		if bytes != 0 {
			out, err := readNBytes(c.Reader, bytes)
			if err != nil {
				return Response{}, outputs, err
			}

			response.Message = out
		}

		break
	}

	return response, outputs, nil
}

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
