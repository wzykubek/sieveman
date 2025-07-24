package client

import (
	"errors"
	"fmt"
	"os"
)

func (c *Client) GetScriptList() (scripts []Script, err error) {
	Logger.Println("Trying to obtain script list")

	cmd := "LISTSCRIPTS"
	lines, err := c.SendCommand(cmd)
	if err != nil {
		return scripts, err
	}

	for _, line := range lines {
		out, err := parseScriptItem(line)
		if err != nil {
			return scripts, err
		}
		scripts = append(scripts, out)
	}

	return scripts, nil
}

func (c *Client) GetScriptContent(name string) (content string, err error) {
	Logger.Println("Trying to get script content")

	cmd := fmt.Sprintf("GETSCRIPT \"%s\"", name)
	out, err := c.SendCommand(cmd)
	content = out[0]
	if err != nil {
		return content, err
	}

	return content, nil
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
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	size := fileInfo.Size()
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
