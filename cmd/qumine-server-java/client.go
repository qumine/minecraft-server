package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ClientCommand is the subcommand for running in server mode
var ClientCommand = &cli.Command{
	Name:    "client",
	Aliases: []string{"c"},
	Usage:   "Start the QuMine Server Client",
	Action: func(c *cli.Context) error {
		logrus.Warn("Client mode not yet supported")
		return nil
	},
}
