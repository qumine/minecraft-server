package custom

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Updater is the updater for custom servers.
type Updater struct {
	customURL string
}

// NewCustomUpdater creates a new custom updater.
func NewCustomUpdater() *Updater {
	return &Updater{
		customURL: os.Getenv("SERVER_CUSTOM_URL"),
	}
}

// Update updates the resource, if supported uses cache.
func (u *Updater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":      "CUSTOM",
		"customURL": u.customURL,
	}).Info("server is custom, force updating...")
	err := u.downloadJar(u.customURL)
	if err != nil {
		return err
	}
	return nil
}

func (u *Updater) downloadJar(url string) error {
	logrus.WithField("url", url).Info("downloading jar")
	rsp, getErr := http.Get(url)
	if getErr != nil {
		logrus.WithError(getErr).Error("downloading jar failed")
		return getErr
	}
	defer rsp.Body.Close()
	logrus.WithField("contentLength", rsp.ContentLength).Trace("downloaded jar")

	logrus.Debug("reading jar")
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	body, readErr := ioutil.ReadAll(rsp.Body)
	if readErr != nil {
		logrus.WithError(readErr).Error("reading jar failed")
		return readErr
	}
	logrus.WithField("body", rsp.Body).Trace("read jar")

	logrus.Debug("saving jar")
	saveErr := u.saveCurrentJar(body)
	if saveErr != nil {
		logrus.WithError(saveErr).Error("saving jar failed")
		return saveErr
	}
	logrus.Trace("saved jar")
	return nil
}

func (u *Updater) saveCurrentJar(jar []byte) error {
	if err := ioutil.WriteFile("server.jar", jar, 0); err != nil {
		return err
	}
	return nil
}
