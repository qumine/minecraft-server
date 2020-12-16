package common

import (
	"io/ioutil"
	"net/http"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverJarPath  = "server.jar"
	serverHashPath = "server.hash"
)

// DownloadServerJar downloads a jar from a given url and saves it as server.jar in the current directory.
func DownloadServerJar(url string) error {
	logrus.WithField("url", url).Info("serverJar downloading")

	rsp, getErr := getFromURL(url)
	if getErr != nil {
		return getErr
	}

	body, readErr := readBodyFromResponse(rsp)
	if readErr != nil {
		return readErr
	}

	saveErr := saveServerJar(body)
	if saveErr != nil {
		return saveErr
	}

	logrus.Debug("serverJar downloaded")
	return nil
}

// CompareHash loads the current hash from server.hash in the current directory and compares it with the given hash.
func CompareHash(force bool, hash string) bool {
	logrus.WithField("force", force).WithField("hash", hash).Info("hash comparing")
	if force {
		return false
	}

	currentHash, err := ioutil.ReadFile(serverHashPath)
	if err != nil {
		return false
	}

	outdated := hash == string(currentHash)
	logrus.WithField("outdated", outdated).Debug("hash compared")
	return outdated
}

// SaveHash saves the new hash as server.hash in the current directory.
func SaveHash(hash string) error {
	logrus.WithField("path", serverHashPath).Info("hash saving")
	if err := utils.WriteFileAsString(serverHashPath, hash); err != nil {
		return err
	}
	logrus.Debug("hash saved")
	return nil
}

func getFromURL(url string) (*http.Response, error) {
	logrus.WithField("url", url).Info("getting from URL")
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	logrus.WithField("rsp", rsp).Trace("got from URL")
	return rsp, nil
}

func readBodyFromResponse(rsp *http.Response) ([]byte, error) {
	defer rsp.Body.Close()
	logrus.WithField("rsp", rsp).Debug("reading body from response")
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	logrus.WithField("body", body).Trace("read body from response")
	return body, nil
}

func saveServerJar(content []byte) error {
	logrus.WithField("path", serverJarPath).Debug("saving jar")
	if err := utils.WriteFile(serverJarPath, content); err != nil {
		return err
	}
	logrus.Trace("saved jar")
	return nil
}
