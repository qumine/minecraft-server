package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/qumine-server-java/internal/api"
	"github.com/qumine/qumine-server-java/internal/server/operators"
	"github.com/qumine/qumine-server-java/internal/server/properties"
	su "github.com/qumine/qumine-server-java/internal/server/updater"
	"github.com/qumine/qumine-server-java/internal/server/whitelist"
	"github.com/qumine/qumine-server-java/internal/server/wrapper"
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
)

func init() {
	if helpFlag {
		showUsage()
	}
	if versionFlag {
		showVersion()
	}
}

func main() {
	properties.Configure()
	whitelist.Configure()
	operators.Configure()

	updateServer()
	updatePlugins()

	wrapper := wrapper.NewWrapper()
	api := api.NewAPI(wrapper)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	go wrapper.Start(ctx, wg)
	go api.Start(ctx, wg)

	<-c
	logrus.Info("interrupt received, stopping")

	cancel()
	wg.Wait()
	logrus.Info("stopped")
}

func showUsage() {
	flag.Usage()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%v, commit %v, built at %v", version, commit, date)
	os.Exit(0)
}

func updateServer() {
	updater, err := su.NewUpdater()
	if err != nil {
		logrus.WithError(err).Fatal("Unsupported serverType")
	}
	// TODO: If jar exists continue
	if err := updater.Update(); err != nil {
		logrus.WithError(err).Fatal("Failed to update server")
	}
	updater = nil
}

func updatePlugins() {
	// TODO: Update all Plugins
}
