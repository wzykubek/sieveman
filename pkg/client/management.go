package client

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"go.wzykubek.xyz/sieveman/internal/helpers"
	"go.wzykubek.xyz/sieveman/internal/parsers"
	"go.wzykubek.xyz/sieveman/pkg/proto"
)

func (c *Client) ListScripts() (s []proto.Script, err error) {
	Logger.Println("Trying to obtain script list")

	c.Write("LISTSCRIPTS")
	r, m, err := c.ReadResponse()
	if err != nil {
		return s, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		s = parsers.ParseScriptList(m)
	} else {
		err = errors.New(r.Message())
	}

	return s, nil
}

func (c *Client) GetScriptLines(name string) (size int64, script []string, err error) {
	Logger.Println("Trying to get script content")

	c.Write(fmt.Sprintf("GETSCRIPT \"%s\"", name))
	r, m, err := c.ReadResponse()
	if err != nil {
		return size, script, err
	}
	logResponse(r)

	if r.Type() == "OK" {
		size, _ = strconv.ParseInt(m[0], 10, 64) // TODO: Handle error
		script = m[1 : len(m)-1]
	} else {
		err = errors.New(r.Message())
	}

	return size, script, err
}

func (c *Client) HaveSpace(name string, size int64) error {
	Logger.Println("Trying to check available space")

	c.Write(fmt.Sprintf("HAVESPACE \"%s\" %d", name, size))
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

	err = c.HaveSpace(name, size)
	if err != nil {
		return err
	}

	Logger.Printf("Trying to put '%s' script\n", name)

	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}
	script := string(fileContent)

	c.Write(fmt.Sprintf("PUTSCRIPT \"%s\" {%d+}\n%s", name, size, script))
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
