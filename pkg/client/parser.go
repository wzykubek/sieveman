package client

import (
	"strconv"
	"strings"
)

type Response struct {
	Name    string
	Code    ResponseCode
	Message string
	Bytes   int
}

type ResponseCode struct {
	Name    string
	Message string
}

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

// parseResponse returns OK, NO or BYE string
func (p *Parser) parseResponse() (response string) {
	if p.position >= len(p.input) {
		return
	}

	start := p.position

	for p.position < len(p.input) && p.input[p.position] != ' ' {
		p.position++
	}

	response = p.input[start:p.position]

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
		code = parts[0]
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

func ParseLine(line string) (response Response, bytes int, err error) {
	p := &Parser{input: line, position: 0}

	responseName := p.parseResponse()

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
		Bytes:   bytes,
	}
	return response, bytes, nil
}
