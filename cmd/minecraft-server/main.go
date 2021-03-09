package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"

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

func main() {
	compiled, _ := time.Parse("yyyy-mm-ddThh:mm:ssZ", date)
	app := &cli.App{
		Name:        "QuMine Server",
		HelpName:    "./qumine-server",
		Usage:       "Minecraft-Server wrapper",
		Version:     version,
		Description: "QuMine Server is a simple wrapper for minecraft servers that handles basic stuff",
		Commands: []*cli.Command{
			ServerCommand,
			ConsoleCommand,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				EnvVars: []string{"DEBUG"},
			},
			&cli.BoolFlag{
				Name:    "trace",
				Aliases: []string{"t"},
				EnvVars: []string{"TRACE"},
			},
		},
		EnableBashCompletion: true,
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}
			if c.Bool("trace") {
				logrus.SetLevel(logrus.TraceLevel)
			}
			return nil
		},
		Compiled: compiled,
		Authors: []*cli.Author{
			{
				Name:  "Cedric Lewe",
				Email: "cedric@qumine.io",
			},
		},
		Copyright:              "(c) 2020 QuMine",
		UseShortOptionHandling: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
