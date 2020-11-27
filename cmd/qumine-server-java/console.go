package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qumine/qumine-server-java/internal/grpc/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ConsoleCommand is the subcommand for running in console mode
var ConsoleCommand = &cli.Command{
	Name:    "client",
	Aliases: []string{"c"},
	Usage:   "Start the QuMine Server Console",
	Action: func(c *cli.Context) error {
		cl := client.NewClient()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}

		go cl.Start(ctx, wg)

		<-interrupt
		logrus.Info("interrupt received, stopping")

		cancel()
		wg.Wait()
		logrus.Info("stopped")
		return nil
	},
}
