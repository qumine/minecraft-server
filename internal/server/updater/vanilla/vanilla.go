package vanilla

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
	serverVanillaAPI  = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

// Updater is the updater for vanilla servers.
type Updater struct {
	serverVersion     string
	serverForceUpdate bool
	serverVanillaAPI  string
}

// NewVanillaUpdater creates a new vanilla updater.
func NewVanillaUpdater() *Updater {
	return &Updater{
		serverVersion:     utils.GetEnvString("SERVER_VERSION", serverVersion),
		serverForceUpdate: utils.GetEnvBool("SERVER_FORCE_UPDATE", serverForceUpdate),
		serverVanillaAPI:  utils.GetEnvString("SERVER_VANILLA_API", serverVanillaAPI),
	}
}

// Update updates the resource, if supported uses cache.
func (u *Updater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":        "VANILLA",
		"version":     u.serverVersion,
		"forceUpdate": u.serverForceUpdate,
		"vanillaAPI":  u.serverVanillaAPI,
	}).Info("checking for server updates")

	logrus.Debug("getting VersionManifest")
	versionManifest, err := getVersionManifest(u.serverVanillaAPI)
	if err != nil {
		return err
	}
	logrus.Trace("got VersionManifest")

	logrus.Debug("resolving version")
	if u.serverVersion == "latest" {
		u.serverVersion = versionManifest.Latest.Release
	}
	logrus.WithField("serverVersion", u.serverVersion).Trace("resolved version")

	logrus.Debug("resolving version details download URL")
	var versionDetailsDownloadURL string
	for i := 0; i < len(versionManifest.Versions); i++ {
		if u.serverVersion == versionManifest.Versions[i].ID {
			versionDetailsDownloadURL = versionManifest.Versions[i].URL
		}
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
	err = u.download(versionDetails.Downloads.Server.URL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	u.saveCurrentHash(versionDetails.Downloads.Client.Sha1)
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
	if versionDetails.Downloads.Server.Sha1 != currentHash {
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
