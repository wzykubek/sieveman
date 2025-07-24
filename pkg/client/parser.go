package client

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	input    string
	position int
}

// skipWhitespace skips whitespace characters in the input string
func (p *Parser) skipWhitespace() {
	if p.position >= len(p.input) {
		return
	}

	if p.input[p.position] == ' ' {
		p.position++
	}
}

// parseResponseName returns OK, NO or BYE string
func (p *Parser) parseResponseName() (response string) {
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

// parseResponseCode parses response code: `(CODE "string")`
func (p *Parser) parseReponseCode() (code string, message string) {
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

func (p *Parser) parseQuotedMessage() (message string) {
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

func (p *Parser) parseBytes() (bytes int, err error) {
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

func parseInlineResponse(line string) (response Response, bytes int, err error) {
	p := &Parser{input: line, position: 0}

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
