package yatopia

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	yatopiaDownloadURL = "https://api.yatopia.net/v2/latestBuild"
)

// YatopiaUpdater is the updater for yatopia servers.
type YatopiaUpdater struct {
	serverVersion       string
	serverDownloadAPI   string
	serverForceDownload bool
}

// NewYatopiaUpdater creates a new yatopia updater.
func NewYatopiaUpdater(serverVersion string, serverDownloadAPI string, serverForceDownload bool) *YatopiaUpdater {
	if len(serverDownloadAPI) < 1 {
		serverDownloadAPI = yatopiaDownloadURL
	}
	return &YatopiaUpdater{
		serverVersion:       serverVersion,
		serverDownloadAPI:   serverDownloadAPI,
		serverForceDownload: serverForceDownload,
	}
}

// Update updates the resource, if supported uses cache.
func (ysu *YatopiaUpdater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":          "YATOPIA",
		"version":       ysu.serverVersion,
		"downloadAPI":   ysu.serverDownloadAPI,
		"forceDownload": ysu.serverForceDownload,
	}).Info("checking for server updates")

	logrus.Debug("resolving version details download URL")
	versionDetailsDownloadURL := ysu.serverDownloadAPI
	if ysu.serverVersion != "latest" {
		versionDetailsDownloadURL = versionDetailsDownloadURL + "?branch=ver/" + ysu.serverVersion
	}
	logrus.WithField("url", versionDetailsDownloadURL).Trace("resolved version details download URL")

	logrus.WithField("url", versionDetailsDownloadURL).Debug("getting version details")
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}
	logrus.WithField("versionDetails", versionDetails).Trace("got version details")

	outdated, err := ysu.isOutdated(versionDetails)
	if err != nil {
		return err
	}
	if !outdated {
		logrus.Info("server already up to date")
		return nil
	}

	logrus.Info("server outdated, updating...")
	err = ysu.download(versionDetails.DownloadURL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	ysu.saveCurrentHash(versionDetails.Branch.Commit.Sha)
	logrus.Trace("saved new hash")

	logrus.Info("updated server")
	return nil
}

func (ysu *YatopiaUpdater) isOutdated(versionDetails *VersionDetails) (bool, error) {
	if ysu.serverForceDownload {
		return true, nil
	}

	if _, err := os.Stat("server.jar"); err != nil {
		return true, nil
	}
	if _, err := os.Stat("server.hash"); err != nil {
		return true, nil
	}

	currentHash, err := ysu.loadCurrentHash()
	if err != nil {
		return true, nil
	}
	if versionDetails.Branch.Commit.Sha != currentHash {
		return true, nil
	}

	return false, nil
}

func (ysu *YatopiaUpdater) download(url string) error {
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
	saveErr := ysu.saveCurrentJar(body)
	if saveErr != nil {
		logrus.WithError(saveErr).Error("saving jar failed")
		return saveErr
	}
	logrus.Trace("saved jar")
	return nil
}

func (ysu *YatopiaUpdater) loadCurrentHash() (string, error) {
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

func (ysu *YatopiaUpdater) saveCurrentJar(jar []byte) error {
	out, err := os.Create("server.jar")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(jar)
	return err
}

func (ysu *YatopiaUpdater) saveCurrentHash(hash string) error {
	out, err := os.Create("server.hash")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString(hash)
	return err
}
