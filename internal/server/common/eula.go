package common

import (
	"errors"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const eulaPath = "eula.txt"

// ConfigureEula writes the eula.txt in the current directory if the env variable EULA is set.
func ConfigureEula() error {
	logrus.Info("eula configuring")
	if !utils.GetEnvBool("EULA", false) {
		return errors.New("EULA is not true")
	}

	utils.WriteFileAsString(eulaPath, "eula=true")

	logrus.Info("eula configured")
	return nil
}
