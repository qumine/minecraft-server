package yatopia

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverYatopiaAPI  = "https://api.yatopia.net/v2/latestBuild"
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
	err = u.download(versionDetails.DownloadURL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	u.saveCurrentHash(versionDetails.Branch.Commit.Sha)
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
	if versionDetails.Branch.Commit.Sha != currentHash {
		return true, nil
	}

	return false, nil
}

func (u *Updater) download(url string) error {
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

func (u *Updater) loadCurrentHash() (string, error) {
	file, openErr := os.Open("server.jar")
	if openErr != nil {
		return "", openErr
	}
	defer file.Close()
	hash := sha1.New()
	if _, copyErr := io.Copy(hash, file); copyErr != nil {
		return "", copyErr
	}
	hashInBytes := hash.Sum(nil)[:20]
	return hex.EncodeToString(hashInBytes), nil
}

func (u *Updater) saveCurrentJar(jar []byte) error {
	out, err := os.Create("server.jar")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(jar)
	return err
}

func (u *Updater) saveCurrentHash(hash string) error {
	out, err := os.Create("server.hash")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString(hash)
	return err
}
