package common

import (
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	whitelistPath = "white-list.txt"

	envWhitelist         = "SERVER_WHITE_LIST"
	envWhitelistOverride = "SERVER_WHITE_LIST_OVERRIDE"
)

// ConfigureWhitelist configures the whitelist.json
func ConfigureWhitelist() error {
	whitelist := utils.GetEnvStringList(envOps, "")
	whitelistOverride := utils.GetEnvBool(envOpsOverride, false)
	logrus.WithFields(logrus.Fields{
		"whitelist":         whitelist,
		"whitelistOverride": whitelistOverride,
	}).Debugf("configuring %s", whitelistPath)

	if !whitelistOverride && utils.FileExists(whitelistPath) {
		logrus.Infof("configuring %s skipped, file already exists. Override this behaviour using the %s", opsPath, envWhitelistOverride)
		return nil
	}

	if err := utils.WriteFileAsString(opsPath, strings.Join(whitelist, "\n")); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"whitelist":         whitelist,
		"whitelistOverride": whitelistOverride,
	}).Infof("configured %s", whitelistPath)
	return nil
}
