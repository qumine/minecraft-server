package main

import (
	"github.com/jroimartin/gocui"
	"github.com/qumine/qumine-server-java/internal/grpc/client"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ConsoleCommand is the subcommand for running in console mode
var ConsoleCommand = &cli.Command{
	Name:    "console",
	Aliases: []string{"c"},
	Usage:   "Start the QuMine Server Console",
	Action: func(c *cli.Context) error {
		client := client.NewClient()

		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create ui")
		}
		defer g.Close()

		g.SetManagerFunc(layout)
		g.SetKeybinding("console", gocui.KeyEnter, gocui.ModNone, client.Send)
		g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, client.Stop)
		go client.Start(g)
		g.MainLoop()
		return nil
	},
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if messages, err := g.SetView("logs", -1, -1, maxX+1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetViewOnTop("logs")
		messages.Autoscroll = true
		messages.Frame = false
		messages.Wrap = true
	}

	if input, err := g.SetView("console", -1, maxY-2, maxX+1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView("console")
		input.Autoscroll = false
		input.Editable = true
		input.Frame = false
		input.Wrap = true
	}
	return nil
}
