package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/qumine-server-java/internal/api"
	grpc "github.com/qumine/qumine-server-java/internal/grpc/server"
	"github.com/qumine/qumine-server-java/internal/server"
	"github.com/qumine/qumine-server-java/internal/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ServerCommand is the subcommand for running in server mode
var ServerCommand = &cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "Start the QuMine Server",
	Action: func(c *cli.Context) error {
		server, err := server.NewServer()
		if err != nil {
			logrus.WithError(err).Fatal("Unsupported serverType")
		}
		if err := server.Update(); err != nil {
			logrus.WithError(err).Fatal("server updating failed")
		}
		if err := server.Configure(); err != nil {
			logrus.WithError(err).Fatal("server configuration failed")
		}

		w := wrapper.NewWrapper()
		a := api.NewAPI(w)
		s := grpc.NewServer(w)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}

		go w.Start(ctx, wg)
		go a.Start(ctx, wg)
		go s.Start(ctx, wg)

		<-interrupt
		logrus.Info("interrupt received, stopping")

		cancel()
		wg.Wait()
		logrus.Info("stopped")
		return nil
	},
}
