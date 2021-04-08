package main

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/qumine/minecraft-server/internal/grpc/client"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

// ConsoleCommand is the subcommand for running in console mode
var ConsoleCommand = &cli.Command{
	Name:    "console",
	Aliases: []string{"c"},
	Usage:   "Start the QuMine Server Console",
	Action: func(c *cli.Context) error {
		client := client.NewClient()

		app := tview.NewApplication()
		cout := newConsoleOutput()
		cin := newConsoleInput(app, client, cout)
		app.SetRoot(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(cout, 0, 1, false).
			AddItem(cin, 1, 1, true), true)

		go client.Start(app, cout)
		return app.Run()
	},
}

func newConsoleOutput() *tview.TextView {
	return tview.NewTextView()
}

func newConsoleInput(app *tview.Application, client *client.GRPCClient, cout *tview.TextView) *tview.InputField {
	in := tview.NewInputField()
	in.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			client.Stop(app, cout)
			app.Stop()
		}
		return event
	})
	return in.
		SetDoneFunc(func(key tcell.Key) {
			client.Send(in.GetText())
			in.SetText("")
		}).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholder("Enter command here...")
}
