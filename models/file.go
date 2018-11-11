package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/h2non/filetype"
)

var fileBaseDir string

func GetFile(fileName string) (buffer *bytes.Buffer, contentType string, err error) {
	// Return error if file not exists.
	if _, err := os.Stat(fileBaseDir + "/" + fileName); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("file not exists")
	}

	file, err := ioutil.ReadFile(fileBaseDir + "/" + fileName)
	if err != nil {
		return nil, "", err
	}

	// Determine MIME type of the file.
	kind, _ := filetype.Match(file)
	return bytes.NewBuffer(file), kind.MIME.Value, nil
}

func SaveFile(buffer *bytes.Buffer, fileName string) error {
	// Return error if file exists.
	if _, err := os.Stat(fileBaseDir + "/" + fileName); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists")
	}
	return ioutil.WriteFile(fileBaseDir+"/"+fileName, buffer.Bytes(), 0644)
}

func init() {
	var ok bool
	fileBaseDir, ok = os.LookupEnv("NAGASE_FILES_DIR")
	if !ok {
		fileBaseDir = "/data/nagase/files"
	}
}
