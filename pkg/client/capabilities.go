package client

import (
	"regexp"
	"strconv"
	"strings"

	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func parseCapabilities(messages []string) (capb proto.Capabilities) {
	capb.StartSSL = false

	for _, msg := range messages {
		re := regexp.MustCompile(`"([^"]+)"`)
		matches := re.FindAllString(msg, 2)
		if matches == nil {
			return capb
		}

		var k, v string
		if len(matches) >= 1 {
			k = strings.Trim(matches[0], "\"")
		}
		if len(matches) >= 2 {
			v = strings.Trim(matches[1], "\"")
		}

		switch k {
		case "IMPLEMENTATION":
			capb.Implementation = v
		case "SASL":
			capb.SASL = strings.Fields(v)
		case "SIEVE":
			capb.Sieve = strings.Fields(v)
		case "STARTTLS":
			capb.StartSSL = true
		case "MAXREDIRECTS":
			capb.MaxRedirects, _ = strconv.Atoi(v)
		case "NOTIFY":
			capb.Notify = strings.Fields(v)
		case "LANGUAGE":
			capb.Language = v
		case "OWNER":
			capb.Owner = v
		case "VERSION":
			capb.Version = v
		}
	}

	return capb
}

// GetCapabilities reads server capabilities.
// It returns proto.Capabilities and error if any.
func (c *Client) GetCapabilities() (capb proto.Capabilities, err error) {
	cmd := "CAPABILITY"
	m, err := c.SendCommand(cmd)
	if err != nil {
		return capb, err
	}

	capb = parseCapabilities(m)

	return capb, nil
}
