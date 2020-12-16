package papermc

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverVersion     = "latest"
	serverForceUpdate = false
	serverPapermcAPI  = "https://papermc.io/api/v2/projects/paper/"
)

// Updater is the updater for papermc servers.
type Updater struct {
	serverVersion     string
	serverForceUpdate bool
	papermcAPI        string
}

// NewPaperMCUpdater creates a new papermc updater.
func NewPaperMCUpdater() *Updater {
	return &Updater{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		papermcAPI:        utils.GetEnvString("SERVER_PAPERMC_API", serverPapermcAPI),
	}
}

// Update updates the resource, if supported uses cache.
func (u *Updater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":        "PAPERMC",
		"version":     u.serverVersion,
		"forceUpdate": u.serverForceUpdate,
		"papermcAPI":  u.papermcAPI,
	}).Info("checking for server updates")

	logrus.Debug("resolving version")
	version := ""
	if match, _ := regexp.MatchString("\\d*\\.\\d*\\.\\d", u.serverVersion); match {
		version = u.serverVersion
	} else if match, _ := regexp.MatchString("\\d*\\.\\d*", u.serverVersion); match {
		versionGroupDetailsDownloadURL := u.papermcAPI + "version_group/" + u.serverVersion
		logrus.WithField("url", versionGroupDetailsDownloadURL).Debug("getting versionGroupDetails")
		versionGroupDetails, err := getVersionGroupDetails(versionGroupDetailsDownloadURL)
		if err != nil {
			return err
		}
		logrus.WithField("versionGroupDetails", versionGroupDetails).Trace("got versionGroupDetails")

		version = versionGroupDetails.Versions[len(versionGroupDetails.Versions)-1]
	} else if u.serverVersion == "latest" {
		// TODO: Implement latest version resolver
	} else {
		return errors.New("Unsupported version")
	}
	logrus.WithField("version", version).Trace("resolved version")

	versionDetailsDownloadURL := u.papermcAPI + "versions/" + version
	logrus.WithField("url", versionDetailsDownloadURL).Debug("getting versionDetails")
	versionDetails, err := getVersionDetails(versionDetailsDownloadURL)
	if err != nil {
		return err
	}
	logrus.WithField("versionDetails", versionDetails).Trace("got versionDetails")

	buildDetailsURL := versionDetailsDownloadURL + "/builds/" + strconv.Itoa(versionDetails.Builds[len(versionDetails.Builds)-1])
	logrus.WithField("url", buildDetailsURL).Debug("getting build details")
	buildDetails, err := getBuildDetails(buildDetailsURL)
	if err != nil {
		return err
	}
	logrus.WithField("versionDetails", versionDetails).Trace("got build details")

	outdated, err := u.isOutdated(buildDetails)
	if err != nil {
		return err
	}
	if !outdated {
		logrus.Info("server already up to date")
		return nil
	}
	logrus.Info("server outdated, updating...")

	downloadURL := buildDetailsURL + "/downloads/" + buildDetails.Downloads.Application.Name
	err = u.downloadJar(downloadURL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	err = utils.WriteFileAsString("server.hash", buildDetails.Downloads.Application.Sha256)
	if err != nil {
		return err
	}
	logrus.Trace("saved new hash")

	logrus.Info("updated server")
	return nil
}

func (u *Updater) isOutdated(buildDetails *BuildDetails) (bool, error) {
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
	logrus.WithField("currentHash", currentHash).WithField("newHash", buildDetails.Downloads.Application.Sha256).Debug("comparing hashes")
	if buildDetails.Downloads.Application.Sha256 != currentHash {
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
