package yatopia

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverYatopiaAPI  = "https://api.yatopiamc.org/v2/latestBuild"
)

// Updater is the updater for yatopia servers.
type Updater struct {
	serverVersion     string
	serverForceUpdate bool
	yatopiaAPI        string
}

// NewYatopiaUpdater creates a new yatopia updater.
func NewYatopiaUpdater() *Updater {
	return &Updater{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		yatopiaAPI:        utils.GetEnvString("SERVER_YATOPIA_API", serverYatopiaAPI),
	}
}

// Update updates the resource, if supported uses cache.
func (u *Updater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":        "YATOPIA",
		"version":     u.serverVersion,
		"forceUpdate": u.serverForceUpdate,
		"yatopiaApi":  u.yatopiaAPI,
	}).Info("checking for server updates")

	logrus.Debug("resolving version details download URL")
	versionDetailsDownloadURL := u.yatopiaAPI
	if u.serverVersion != "latest" {
		versionDetailsDownloadURL = versionDetailsDownloadURL + "?branch=ver/" + u.serverVersion
	}
	logrus.WithField("url", versionDetailsDownloadURL).Trace("resolved version details download URL")

	logrus.WithField("url", versionDetailsDownloadURL).Debug("getting version details")
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}
	logrus.WithField("versionDetails", versionDetails).Trace("got version details")

	outdated, err := u.isOutdated(versionDetails)
	if err != nil {
		return err
	}
	if !outdated {
		logrus.Info("server already up to date")
		return nil
	}

	logrus.Info("server outdated, updating...")
	err = u.downloadJar(versionDetails.DownloadURL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	err = utils.WriteFileAsString("server.hash", versionDetails.Branch.Commit.Sha)
	if err != nil {
		return err
	}
	logrus.Trace("saved new hash")

	logrus.Info("updated server")
	return nil
}

func (u *Updater) isOutdated(versionDetails *VersionDetails) (bool, error) {
	if u.serverForceUpdate {
		return true, nil
	}

	if _, err := os.Stat("server.jar"); err != nil {
		return true, nil
	}
	if _, err := os.Stat("server.hash"); err != nil {
		return true, nil
	}

	currentHash, err := u.loadCurrentHash()
	if err != nil {
		return true, nil
	}
	logrus.WithField("currentHash", currentHash).WithField("newHash", versionDetails.Branch.Commit.Sha).Debug("comparing hashes")
	if versionDetails.Branch.Commit.Sha != currentHash {
		return true, nil
	}

	return false, nil
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
	saveErr := utils.WriteFile("server.jar", body)
	if saveErr != nil {
		logrus.WithError(saveErr).Error("saving jar failed")
		return saveErr
	}
	logrus.Trace("saved jar")
	return nil
}

func (u *Updater) loadCurrentHash() (string, error) {
	hash, err := ioutil.ReadFile("server.hash")
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
