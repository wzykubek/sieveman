package helpers

import (
	"encoding/base64"
	"os"
)

func GetByteSize(file *os.File) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func EncCredentials(login string, password string) string {
	data := []byte("\x00" + login + "\x00" + password)
	return base64.StdEncoding.EncodeToString(data)
}
