package client

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func parseCapabilities(messages []string) (cap Capabilities) {
	cap.StartSSL = false

	for _, msg := range messages {
		re := regexp.MustCompile(`"([^"]+)"`)
		matches := re.FindAllString(msg, 2)
		if matches == nil {
			return cap
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
			cap.Implementation = v
		case "SASL":
			cap.SASL = strings.Fields(v)
		case "SIEVE":
			cap.Sieve = strings.Fields(v)
		case "STARTTLS":
			cap.StartSSL = true
		case "MAXREDIRECTS":
			cap.MaxRedirects, _ = strconv.Atoi(v)
		case "NOTIFY":
			cap.Notify = strings.Fields(v)
		case "LANGUAGE":
			cap.Language = v
		case "OWNER":
			cap.Owner = v
		case "VERSION":
			cap.Version = v
		}
	}

	return cap
}

// GetCapabilities reads server capabilities.
// It returns proto.Capabilities and error if any.
func (c *Client) GetCapabilities() (cap Capabilities, err error) {
	cmd := "CAPABILITY"
	err = c.WriteLine(cmd)
	if err != nil {
		return cap, err
	}
	cap, err = c.ReadCapabilities()
	if err != nil {
		return cap, err
	}
	r, _, err := c.ReadResponse()
	if err != nil {
		return cap, err
	}
	if r.Name != "OK" {
		return cap, errors.New(r.Message)
	}

	return cap, nil
}
