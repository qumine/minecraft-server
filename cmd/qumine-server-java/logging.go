package main

import (
	"strings"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	configureLogLevel()
	configureLogger()
}

func configureLogLevel() {
	switch strings.ToUpper(utils.GetEnvString("LOG_LEVEL", "INFO")) {
	case "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
		break
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
		break
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
		break
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
		break
	}
}

func configureLogger() {
	// TODO: Configure other logger stuff here
}
