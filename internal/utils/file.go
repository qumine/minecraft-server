package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// FileExists returns wether a file exists or not
func FileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

// WriteFileAsString writes a string to a file.
func WriteFileAsString(path string, content string) error {
	return WriteFile(path, []byte(content))
}

// WriteFileAsJSON writes JSON to a file.
func WriteFileAsJSON(path string, content interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return WriteFile(path, []byte(b))
}

// WriteFile writes bytes to a file.
func WriteFile(path string, content []byte) error {
	logrus.WithFields(logrus.Fields{
		"path": path,
	}).Trace("writing file")

	if err := ioutil.WriteFile(path, content, os.ModePerm); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"path": path,
	}).Debug("written file")
	return nil
}
