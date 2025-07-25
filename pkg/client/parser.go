package client

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type parser struct {
	input    string
	position int
}

// skipWhitespace skips whitespace characters in the input string.
func (p *parser) skipWhitespace() {
	if p.position >= len(p.input) {
		return
	}

	if p.input[p.position] == ' ' {
		p.position++
	}
}

// parseResponseName returns OK, NO or BYE string.
func (p *parser) parseResponseName() (response string) {
	if p.position >= len(p.input) {
		return
	}

	start := p.position

	for p.position < len(p.input) && p.input[p.position] != ' ' {
		p.position++
	}

	response = strings.ToUpper(p.input[start:p.position])

	return response
}

// parseResponseCode parses response code: e.g. `(CODE "string")`.
func (p *parser) parseReponseCode() (code string, message string) {
	if p.position >= len(p.input) {
		return
	}

	if p.input[p.position] == '(' {
		p.position++
		start := p.position

		for p.input[p.position] != ')' {
			p.position++
		}

		// Example values:
		// TAG "some string"
		// NONEXISTENT
		// QUOTA/MAXSCRIPTS
		parentheses := p.input[start:p.position]
		parts := strings.SplitN(parentheses, " ", 2)
		code = strings.ToUpper(parts[0])
		if len(parts) == 2 {
			message = strings.Trim(parts[1], "\"")
		}

		p.position++
	}

	return code, message
}

// parseQuotedMessage parses quoted message: e.g. `"some string"`.
func (p *parser) parseQuotedMessage() (message string) {
	if p.position >= len(p.input) {
		return
	}

	if p.input[p.position] == '"' {
		p.position++
		start := p.position

		for p.position < len(p.input) && p.input[p.position] != '"' {
			p.position++
		}

		message = p.input[start:p.position]
	}

	return message
}

// parseBytes parses bytes: e.g. `{123}`.
func (p *parser) parseBytes() (bytes int, err error) {
	if p.position >= len(p.input) {
		return
	}

	if p.input[p.position] == '{' {
		p.position++
		start := p.position

		for p.position < len(p.input) && p.input[p.position] != '}' {
			p.position++
		}

		bytesStr := p.input[start:p.position]
		bytes, err = strconv.Atoi(bytesStr)
		if err != nil {
			return bytes, err
		}
	}

	return bytes, nil
}

// parseInlineResponse parses inline response: e.g. `OK (CODE "string") {123}`.
func parseInlineResponse(line string) (response Response, bytes int, err error) {
	p := &parser{input: line, position: 0}

	responseName := p.parseResponseName()

	p.skipWhitespace()

	responseCodeName, responseCodeMessage := p.parseReponseCode()
	responseCode := ResponseCode{
		Name:    responseCodeName,
		Message: responseCodeMessage,
	}

	p.skipWhitespace()

	responseMessage := p.parseQuotedMessage()

	if responseMessage == "" {
		bytes, err = p.parseBytes()
		if err != nil {
			return Response{}, bytes, err
		}
	}

	response = Response{
		Name:    responseName,
		Code:    responseCode,
		Message: responseMessage,
	}
	return response, bytes, nil
}

// parseScriptItem parses script item: e.g. `"script.sieve" ACTIVE`.
func parseScriptItem(line string) (script Script, err error) {
	re := regexp.MustCompile(`"([^"]+)"(\s*ACTIVE)?`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return script, errors.New("Can't parse script item")
	}

	name := matches[1]
	active := len(matches[2]) > 0
	script = Script{Name: name, Active: active}

	return script, nil
}

// parseCapability parses capability: e.g. `"STARTTLS"` or `"SIEVE" "imapfilter"`.
func parseCapability(cap *Capabilities, line string) error {
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindAllString(line, 2)
	if matches == nil {
		return errors.New("Invalid capability")
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
	default:
		return errors.New("Invalid capability")
	}

	return nil
}
