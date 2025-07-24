// TODO: Refactoring of whole reading mechanism is needed.
// Responses should not have message included, only OK, NO or BYE.
// All messages should be parsed and returned as separate fields.
// Handling of both response patterns should be implemented.
// Response codes is subject to change.
package client

import (
	"bufio"
	"io"
	"strings"

	"go.wzykubek.xyz/sieveman/pkg/proto"
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
// It returns parsed response, slice of messages and error if any.
func (c *Client) ReadResponse() (response Response, rawOutputLines []string, err error) {
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return Response{}, rawOutputLines, err
		}

		if strings.HasPrefix(line, "OK") || strings.HasPrefix(line, "NO") || strings.HasPrefix(line, "BYE") {
			var bytes int
			response, bytes, err = ParseLine(strings.TrimSpace(line))
			if err != nil {
				return Response{}, rawOutputLines, err
			}

			if bytes != 0 {
				content, err := readBytes(c.Reader, bytes)
				if err != nil {
					return Response{}, rawOutputLines, err
				}
				response.Message = content
			}

			break
		} else {
			// TODO: LISTSCRIPT output does not return number of bytes
			// and in current implementation it is easier to read by
			// lines and ignore bytes.
			rawOutputLines = append(rawOutputLines, line)
		}
	}

	return response, rawOutputLines, nil
}

// TODO: Rewrite this function and capabilities parser
func (c *Client) ReadCapabilities() (cap proto.Capabilities, err error) {
	var linesCompat []string
	for {
		line, err := c.Reader.ReadString('\n')
		if err != nil {
			return cap, err
		}
		linesCompat = append(linesCompat, line)
		// Warning: Quick workaround, refactoring of parseCapabilities needed
		if strings.HasPrefix(line, "\"VERSION") {
			break
		}
	}
	cap = parseCapabilities(linesCompat)
	return cap, nil
}
