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

	cmd := "LISTSCRIPTS"
	m, err := c.SendCommand(cmd)
	if err != nil {
		return scripts, err
	}

	scripts = parseScriptList(m)

	return scripts, nil
}

func (c *Client) GetScriptLines(name string) (size int64, script []string, err error) {
	Logger.Println("Trying to get script content")

	cmd := fmt.Sprintf("GETSCRIPT \"%s\"", name)
	m, err := c.SendCommand(cmd)
	if err != nil {
		return size, script, err
	}

	// TODO: Consider changing it to proto.Script
	size, _ = strconv.ParseInt(m[0], 10, 64) // TODO: Handle error
	script = m[1 : len(m)-1]

	return size, script, nil
}

func (c *Client) CheckSpace(name string, size int64) error {
	Logger.Println("Trying to check available space")

	cmd := fmt.Sprintf("HAVESPACE \"%s\" %d", name, size)
	_, err := c.SendCommand(cmd)
	if err != nil {
		return err
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

	cmd := fmt.Sprintf("PUTSCRIPT \"%s\" {%d+}\n%s", name, size, script)
	_, err = c.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ActivateScript(name string) error {
	Logger.Println("Trying to activate script")

	cmd := fmt.Sprintf("SETACTIVE \"%s\"", name)
	_, err := c.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeactivateScripts() error {
	Logger.Println("Trying to deactivate script")

	return c.ActivateScript("")
}

func (c *Client) RemoveScript(name string) error {
	Logger.Printf("Trying to remove script")

	cmd := fmt.Sprintf("DELETESCRIPT \"%s\"", name)
	_, err := c.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RenameScript(oldName, newName string) error {
	Logger.Printf("Trying to rename script")

	cmd := fmt.Sprintf("RENAMESCRIPT \"%s\" \"%s\"", oldName, newName)
	_, err := c.SendCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}
