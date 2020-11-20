package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	version = "dev"
	commit  = "none"
	date    = "uknown"
)

var (
	helpFlag    bool
	versionFlag bool
	debugFlag   bool
	traceFlag   bool
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show this page")
	flag.BoolVar(&versionFlag, "version", false, "Show the current version")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debugging log level")
	flag.BoolVar(&traceFlag, "trace", false, "Enable trace log level")
	flag.Parse()
}

func main() {
	if helpFlag {
		showHelp()
	}

	if versionFlag {
		showVersion()
	}

	if debugFlag {
		enableDebug()
	}
	if debugFlag {
		enableTrace()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c
}

func showHelp() {
	flag.Usage()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%v, commit %v, built at %v", version, commit, date)
	os.Exit(0)
}

func enableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debug("debugging enabled")
}

func enableTrace() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Debug("tracing enabled")
}
