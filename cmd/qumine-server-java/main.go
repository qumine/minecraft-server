package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/qumine/qumine-server-java/internal/api"
	"github.com/qumine/qumine-server-java/internal/server/operators"
	"github.com/qumine/qumine-server-java/internal/server/properties"
	su "github.com/qumine/qumine-server-java/internal/server/updater"
	"github.com/qumine/qumine-server-java/internal/server/whitelist"
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
	compiled, _ := time.Parse("yyyy-mm-ddThh:mm:ssZ", date)
	app := &cli.App{
		Name:     "QuMine Server",
		HelpName: "./qumine-server",
		Usage:    "Minecraft-Server wrapper",

		Version:     version,
		Description: "QuMine Server is a simple wrapper for minecraft servers that handles basic stuff",
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start the QuMine Server",
				Action: func(c *cli.Context) error {
					properties.Configure()
					whitelist.Configure()
					operators.Configure()

					updateServer()
					updatePlugins()

					w := wrapper.NewWrapper()
					a := api.NewAPI(w)

					interrupt := make(chan os.Signal, 1)
					signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
					ctx, cancel := context.WithCancel(context.Background())
					wg := &sync.WaitGroup{}

					go w.Start(ctx, wg)
					go a.Start(ctx, wg)

					<-interrupt
					logrus.Info("interrupt received, stopping")

					cancel()
					wg.Wait()
					logrus.Info("stopped")
					return nil
				},
			},
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Start the QuMine Server Client",
				Action: func(c *cli.Context) error {
					logrus.Warn("Client mode not yet supported")
					return nil
				},
			},
		},
		EnableBashCompletion: true,
		Compiled:             compiled,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Cedric Lewe",
				Email: "cedric@qumine.io",
			},
		},
		Copyright: "(c) 2020 QuMine",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
