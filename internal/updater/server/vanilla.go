package server

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
	vanillaDownloadAPI = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

// VanillaUpdater the updater for vanilla servers.
type VanillaUpdater struct {
	serverVersion      string
	serverDownloadAPI  string
	serverForceDownoad bool
}

// NewVanillaUpdater creates a new vanilla updater.
func NewVanillaUpdater(serverVersion string, serverDownloadAPI string, serverForceDownoad bool) *VanillaUpdater {
	if len(serverDownloadAPI) < 1 {
		serverDownloadAPI = vanillaDownloadAPI
	}
	return &VanillaUpdater{
		serverVersion:      serverVersion,
		serverDownloadAPI:  serverDownloadAPI,
		serverForceDownoad: serverForceDownoad,
	}
}

// Update updates the resource, if supported uses cache.
func (vsu *VanillaUpdater) Update() error {
	logrus.WithFields(logrus.Fields{
		"type":          "VANILLA",
		"version":       vsu.serverVersion,
		"downloadAPI":   vsu.serverDownloadAPI,
		"forceDownload": vsu.serverForceDownoad,
	}).Info("checking for server updates")

	logrus.Debug("getting VersionManifest")
	versionManifest, err := getVersionManifest(vsu.serverDownloadAPI)
	if err != nil {
		return err
	}
	logrus.Trace("got VersionManifest")

	logrus.Debug("resolving version")
	if vsu.serverVersion == "latest" {
		vsu.serverVersion = versionManifest.Latest.Release
	}
	logrus.WithField("serverVersion", vsu.serverVersion).Trace("resolved version")

	logrus.Debug("resolving version details download URL")
	var versionDetailsURL string
	for i := 0; i < len(versionManifest.Versions); i++ {
		if vsu.serverVersion == versionManifest.Versions[i].ID {
			versionDetailsURL = versionManifest.Versions[i].URL
		}
	}
	logrus.WithField("versionDetailsURL", versionDetailsURL).Trace("resolved version details download URL")

	logrus.Debug("getting version details")
	versionDetails, err := getVersionDetails(versionDetailsURL)
	if err != nil {
		return err
	}
	logrus.WithField("versionDetails", versionDetails).Trace("got version details")

	outdated, err := vsu.isOutdated(versionDetails)
	if err != nil {
		return err
	}
	if !outdated {
		logrus.Info("server already up to date")
		return nil
	}

	logrus.Info("server outdated, updating...")
	err = vsu.download(versionDetails.Downloads.Server.URL)
	if err != nil {
		return err
	}

	logrus.Debug("saving new hash")
	vsu.saveCurrentHash(versionDetails.Downloads.Client.Sha1)
	logrus.Trace("saved new hash")

	logrus.Info("updated server")
	return nil
}

func (vsu *VanillaUpdater) isOutdated(versionDetails *VersionDetails) (bool, error) {
	if vsu.serverForceDownoad {
		return true, nil
	}

	if _, err := os.Stat("server.jar"); err != nil {
		return true, nil
	}
	if _, err := os.Stat("server.hash"); err != nil {
		return true, nil
	}

	currentHash, err := vsu.loadCurrentHash()
	if err != nil {
		return true, nil
	}
	if versionDetails.Downloads.Server.Sha1 != currentHash {
		return true, nil
	}

	return false, nil
}

func (vsu *VanillaUpdater) download(url string) error {
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
	saveErr := vsu.saveCurrentJar(body)
	if saveErr != nil {
		logrus.WithError(saveErr).Error("saving jar failed")
		return saveErr
	}
	logrus.Trace("saved jar")
	return nil
}

func (vsu *VanillaUpdater) loadCurrentHash() (string, error) {
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

func (vsu *VanillaUpdater) saveCurrentJar(jar []byte) error {
	out, err := os.Create("server.jar")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.Write(jar)
	return err
}

func (vsu *VanillaUpdater) saveCurrentHash(hash string) error {
	out, err := os.Create("server.hash")
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString(hash)
	return err
}
