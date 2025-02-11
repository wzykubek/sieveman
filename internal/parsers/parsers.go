package parsers

import (
	"regexp"
	"strconv"
	"strings"

	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func ParseResponse(respStr string) proto.Response {
	if respStr == "" {
		return nil
	}

	// This regular expression was generated by Copilot.
	// It SHOULD parse a string containing three parts with the known pattern:
	// 		P1 (P2 ...) "P3 ..."
	//
	// where:
	//		P1		- is mandatory and does not contain whitespace
	//		(P2)	- is optional, contains parentheses, and can contain whitespace inside
	//		"P3"	- is optional, contains quotation marks, and can contain whitespace inside
	//
	// When P2 or "P3" is not present, it MUST return an empty string for those parts.
	re := regexp.MustCompile(`^(\S+)\s*(\([^)]*\))?\s*("[^"]*")?$`)
	matches := re.FindStringSubmatch(respStr)

	resp := matches[1]
	code := ParseResponseCode(matches[2])
	msg := strings.Trim(matches[3], "\"")

	switch resp {
	case "OK":
		return proto.Ok{ResponseCode: code, Msg: msg}
	case "NO":
		return proto.No{ResponseCode: code, Msg: msg}
	case "BYE":
		return proto.Bye{ResponseCode: code, Msg: msg}
	default:
		return nil
	}
}

func ParseResponseCode(codeStr string) proto.ResponseCode {
	if codeStr == "" {
		return nil
	}

	// This regular expression was generated by Copilot.
	// It SHOULD parse a string containing two parts with the known pattern in parentheses:
	// 		(P1 "P2 ...")
	//
	// where:
	//		P1		- is mandatory and does not contain whitespace
	//		"P2"	- is optional, contains quotation marks, and can contain whitespace inside
	//
	// When "P2" is not present, it MUST return an empty string for that part.
	// It also removes the surrounding parentheses.
	re := regexp.MustCompile(`^\((\S+)\s*("[^"]*")?\)$`)
	matches := re.FindStringSubmatch(codeStr)

	code := matches[1]
	msg := strings.Trim(matches[2], "\"")

	switch code {
	case "TAG":
		return proto.Tag{Msg: msg, ChildCode: nil}
	default:
		return nil
	}
}

func ParseCapabilities(messages []string) proto.Capabilities {
	cpb := proto.Capabilities{StartSSL: false}

	for _, msg := range messages {
		re := regexp.MustCompile(`"([^"]+)"`)
		matches := re.FindAllString(msg, 2)
		if matches == nil {
			return cpb
		}

		var k, v string
		if len(matches) == 1 {
			k = strings.Trim(matches[0], "\"")
		}
		if len(matches) == 2 {
			v = strings.Trim(matches[1], "\"")
		}

		switch k {
		case "IMPLEMENTATION":
			cpb.Implementation = v
		case "SASL":
			cpb.SASL = strings.Fields(v)
		case "SIEVE":
			cpb.Sieve = strings.Fields(v)
		case "STARTTLS":
			cpb.StartSSL = true
		case "MAXREDIRECTS":
			cpb.MaxRedirects, _ = strconv.Atoi(v)
		case "NOTIFY":
			cpb.Notify = strings.Fields(v)
		case "LANGUAGE":
			cpb.Language = v
		case "OWNER":
			cpb.Owner = v
		case "VERSION":
			cpb.Version = v
		default:
			continue
		}
	}

	return cpb
}
