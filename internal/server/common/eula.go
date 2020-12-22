package common

import (
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	eulaPath = "eula.txt"

	envEula = "EULA"
)

// ConfigureEula writes the eula.txt in the current directory if the env variable EULA is set.
func ConfigureEula() error {
	eula := utils.GetEnvString(envEula, "false")
	logrus.WithFields(logrus.Fields{
		"eula": eula,
	}).Debugf("configuring %s", eulaPath)

	utils.WriteFileAsString(eulaPath, "eula="+eula)

	logrus.WithFields(logrus.Fields{
		"eula": eula,
	}).Infof("configured %s", eulaPath)
	return nil
}
