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

// WriteFileAsJSON writes a string to a file.
func WriteFileAsString(fileName string, content string) {

}

// WriteFileAsJSON writes JSON to a file.
func WriteFileAsJSON(fileName string, content interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, []byte(b), os.ModePerm); err != nil {
		return err
	}
	return nil
}
