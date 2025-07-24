package client

import (
	"bufio"
	"io"
	"strings"
)

func readBytes(reader *bufio.Reader, byteCount int) (content string, err error) {
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
// TODO: Capabilities reading
func (c *Client) ReadResponse() (response Response, outputs []string, err error) {
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return Response{}, outputs, err
		}

		// TODO: Responses does not need to be uppercase
		if strings.HasPrefix(line, "OK") || strings.HasPrefix(line, "NO") || strings.HasPrefix(line, "BYE") {
			var bytes int
			response, bytes, err = parseInlineResponse(strings.TrimSpace(line))
			if err != nil {
				return Response{}, outputs, err
			}

			if bytes != 0 {
				content, err := readBytes(c.Reader, bytes)
				if err != nil {
					return Response{}, outputs, err
				}
				response.Message = content
			}

			break
		} else {
			// TODO: Test in different scenarios for LISTSCRIPTS
			// e.g. with byte size in the middle, ReadResponse should handle that
			p := Parser{input: line, position: 0}
			bytes, err := p.parseBytes()
			if err != nil {
				return Response{}, outputs, err
			}
			if bytes != 0 {
				out, err := readBytes(c.Reader, bytes)
				if err != nil {
					return Response{}, outputs, err
				}
				outputs = append(outputs, out)
				continue
			}
			outputs = append(outputs, line)
		}
	}

	return response, outputs, nil
}

func (c *Client) ReadCapabilities() (cap Capabilities, err error) {
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return cap, err
		}
		if err := parseCapability(&cap, line); err != nil {
			return cap, err
		}
		if cap.Version != "" {
			break
		}
	}

	return cap, nil
}
