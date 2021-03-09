package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/minecraft-server/internal/api"
	grpc "github.com/qumine/minecraft-server/internal/grpc/server"
	"github.com/qumine/minecraft-server/internal/server"
	"github.com/qumine/minecraft-server/internal/wrapper"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ServerCommand is the subcommand for running in server mode
var ServerCommand = &cli.Command{
	Name:    "server",
	Aliases: []string{"s"},
	Usage:   "Start the QuMine Server",
	Action: func(c *cli.Context) error {
		srv, err := server.NewServer()
		if err != nil {
			logrus.WithError(err).Fatal("Unsupported serverType")
		}
		if err := srv.Configure(); err != nil {
			logrus.WithError(err).Fatal("server configuration failed")
		}
		if err := srv.Update(); err != nil {
			logrus.WithError(err).Fatal("server updating failed")
		}
		if err := srv.UpdatePlugins(); err != nil {
			logrus.WithError(err).Fatal("plugins updating failed")
		}

		w := wrapper.NewWrapper(srv)
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
