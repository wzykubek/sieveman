package client

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"go.wzykubek.xyz/sieveman/internal/helpers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func parseScriptList(m []string) (scripts []proto.Script) {
	re := regexp.MustCompile(`"([^"]+)"(\s*ACTIVE)?`)
	for _, v := range m {
		matches := re.FindStringSubmatch(v)
		if matches != nil {
			name := matches[1]
			active := len(matches[2]) > 0
			scripts = append(scripts, proto.Script{Name: name, Active: active})
		}
	}

	return scripts
}

func (c *Client) GetScriptList() (scripts []proto.Script, err error) {
	Logger.Println("Trying to obtain script list")

	c.WriteLine("LISTSCRIPTS")
	r, m, err := c.ReadResponse()
	if err != nil {
		return scripts, err
	}
	logResponse(r)

	if r.Type() != "OK" {
		return scripts, errors.New(r.Message())
	}

	scripts = parseScriptList(m)

	return scripts, nil
}

func (c *Client) GetScriptLines(name string) (size int64, script []string, err error) {
	Logger.Println("Trying to get script content")

	c.WriteLine(fmt.Sprintf("GETSCRIPT \"%s\"", name))
	r, m, err := c.ReadResponse()
	if err != nil {
		return size, script, err
	}
	logResponse(r)

	if r.Type() != "OK" {
		return size, script, errors.New(r.Message())
	}

	size, _ = strconv.ParseInt(m[0], 10, 64) // TODO: Handle error
	script = m[1 : len(m)-1]

	return size, script, nil
}

func (c *Client) CheckSpace(name string, size int64) error {
	Logger.Println("Trying to check available space")

	c.WriteLine(fmt.Sprintf("HAVESPACE \"%s\" %d", name, size))
	r, _, err := c.ReadResponse()
	if err != nil {
		return err
	}
	logResponse(r)

	if r.Type() != "OK" {
		return errors.New(r.Message())
	}
	// TODO: Handle Quota reponse code

	return nil
}

func (c *Client) PutScript(file *os.File, name string) error {
	size, err := helpers.GetByteSize(file)
	if err != nil {
		return err
	}
	if size == 0 {
		return errors.New("File is empty")
	}

	if err = c.CheckSpace(name, size); err != nil {
		return err
	}

	Logger.Printf("Trying to put '%s' script\n", name)

	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	script := string(fileContent)

	c.WriteLine(fmt.Sprintf("PUTSCRIPT \"%s\" {%d+}\n%s", name, size, script))
	r, _, err := c.ReadResponse() // TODO: Fix handling of response when syntax is bad
	if err != nil {
		return err
	}
	logResponse(r)

	if r.Type() != "OK" {
		return errors.New(r.Message())
	}

	return nil
}
