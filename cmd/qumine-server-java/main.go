package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/qumine-server-java/internal/api"
	su "github.com/qumine/qumine-server-java/internal/updater/server"
	"github.com/qumine/qumine-server-java/internal/wrapper"
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
	logLevel    string

	serverType         string
	serverVersion      string
	serverDownloadAPI  string
	serverForceDownoad bool

	forceDownloadPluginsFlag bool
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "Show this page")
	flag.BoolVar(&versionFlag, "version", false, "Show the current version")
	flag.StringVar(&logLevel, "log-level", "info", "Enable debugging log level")

	flag.StringVar(&serverType, "server-type", "vanilla", "Which server type to use.")
	flag.StringVar(&serverVersion, "server-version", "latest", "Which server version to use.")
	flag.StringVar(&serverDownloadAPI, "server-download-api", "", "Url to the server download api.")
	flag.BoolVar(&serverForceDownoad, "server-force-download", false, "Force the download of the server jar")

	flag.BoolVar(&forceDownloadPluginsFlag, "force-download-plugins", false, "Force the download of the server plugins")
	flag.Parse()

	if helpFlag {
		showUsage()
	}
	if versionFlag {
		showVersion()
	}
}

func main() {
	// configure logging
	setLogLevel(logLevel)

	// build everything
	// serverUpdater
	// pluginUpdater
	wrapper := wrapper.NewWrapper()
	api := api.NewAPI(wrapper)

	var updater su.Updater
	updater = su.NewVanillaUpdater(serverVersion, serverDownloadAPI, serverForceDownoad)
	updater.Update()

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
