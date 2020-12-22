package common

import (
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	opsPath = "ops.txt"

	envOps         = "SERVER_OPS"
	envOpsOverride = "SERVER_OPS_OVERRIDE"
)

// ConfigureOps configures the ops.json
func ConfigureOps() error {
	ops := utils.GetEnvStringList(envOps, "")
	opsOverride := utils.GetEnvBool(envOpsOverride, false)
	logrus.WithFields(logrus.Fields{
		"ops":         ops,
		"opsOverride": opsOverride,
	}).Debugf("configuring %s", opsPath)

	if !opsOverride && utils.FileExists(opsPath) {
		logrus.Infof("configuring %s skipped, file already exists. Override this behaviour using the %s", opsPath, envOpsOverride)
		return nil
	}

	if err := utils.WriteFileAsString(opsPath, strings.Join(ops, "\n")); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"ops":         ops,
		"opsOverride": opsOverride,
	}).Infof("configured %s", opsPath)
	return nil
}
