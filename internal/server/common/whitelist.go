package common

import (
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const whitelistPath = "white-list.txt"

// ConfigureWhitelist configures the whitelist.json
func ConfigureWhitelist() error {
	if !utils.GetEnvBool("SERVER_WHITE_LIST_OVERRIDE", false) && utils.FileExists(whitelistPath) {
		logrus.Infof("%s already exist, skipping configuration", whitelistPath)
		return nil
	}

	whitelist := utils.GetEnvStringList("SERVER_WHITE_LIST", "")
	logrus.WithField("whitelist", whitelist).Infof("%s not found, configuring it now", whitelistPath)

	strings.Join(whitelist, "\n")
	if err := utils.WriteFileAsString(whitelistPath, strings.Join(whitelist, "\n")); err != nil {
		return err
	}
	logrus.Debugf("%s configured", whitelistPath)
	return nil
}
