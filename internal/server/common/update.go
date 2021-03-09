package common

import (
	"io/ioutil"

	"github.com/qumine/minecraft-server/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serverJarPath  = "server.jar"
	serverHashPath = "server.hash"
)

// CompareHash loads the current hash from server.hash in the current directory and compares it with the given hash.
func CompareHash(force bool, hash string) bool {
	logrus.WithField("force", force).WithField("hash", hash).Trace("hash comparing")
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
	logrus.WithField("path", serverHashPath).Trace("hash saving")
	if err := utils.WriteFileAsString(serverHashPath, hash); err != nil {
		return err
	}
	logrus.Debug("hash saved")
	return nil
}
