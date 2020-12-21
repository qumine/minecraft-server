package common

import (
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const opsPath = "ops.txt"

// ConfigureOps configures the ops.json
func ConfigureOps() error {
	if !utils.GetEnvBool("SERVER_OPS_OVERRIDE", false) && utils.FileExists(opsPath) {
		logrus.Infof("%s already exist, skipping configuration", opsPath)
		return nil
	}

	ops := utils.GetEnvStringList("SERVER_OPS", "")
	logrus.WithField("ops", ops).Infof("%s not found, configuring it now", opsPath)

	if err := utils.WriteFileAsString(opsPath, strings.Join(ops, "\n")); err != nil {
		return err
	}
	logrus.Debugf("%s configured", opsPath)
	return nil
}
