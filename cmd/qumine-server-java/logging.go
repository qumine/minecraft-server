package main

import "github.com/sirupsen/logrus"

func init() {
	//logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func setLogLevel(level string) {
	switch level {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
		logrus.Debug("logging set to debug")
		break
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		logrus.Debug("logging set to debug")
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Debug("logging set to debug")
		break
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		logrus.Debug("logging set to debug")
		break
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		logrus.Debug("logging set to debug")
		break
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("logging set to debug")
		break
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		logrus.Debug("logging set to trace")
		break
	}
}
