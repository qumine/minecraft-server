package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

// DownloadToFile downloads a file and saves it to a given path.
func DownloadToFile(url string, path string) error {
	logrus.WithFields(logrus.Fields{
		"url":  url,
		"path": path,
	}).Debug("downloading file")

	rsp, getErr := http.Get(url)
	if getErr != nil {
		return getErr
	}

	defer rsp.Body.Close()
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		return readErr
	}

	if writeErr := WriteFile(path, body); err != nil {
		return writeErr
	}

	logrus.WithFields(logrus.Fields{
		"url":  url,
		"path": path,
	}).Info("downloaded file")
	return nil
}
