package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// FileExists returns wether a file exists or not
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return true
	}
	return false
}

// WriteFileAsString writes a string to a file.
func WriteFileAsString(fileName string, content string) error {
	return WriteFile(fileName, []byte(content))
}

// WriteFileAsJSON writes JSON to a file.
func WriteFileAsJSON(fileName string, content interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return WriteFile(fileName, []byte(b))
}

// WriteFile writes bytes to a file.
func WriteFile(fileName string, content []byte) error {
	if err := ioutil.WriteFile(fileName, content, os.ModePerm); err != nil {
		return err
	}
	return nil
}
