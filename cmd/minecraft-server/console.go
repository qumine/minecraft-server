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

		cout := newConsoleOutput()
		cin := newConsoleInput(client)
		app := tview.NewApplication()
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyCtrlC {
				client.Stop(app, cout)
			}
			return event
		})
		app.SetRoot(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(cout, 0, 1, false).
			AddItem(cin, 1, 1, true), true)

		go client.Start(app, cout)
		return app.Run()
	},
}

func newLayout() *tview.Flex {
	return tview.NewFlex()
}

func newConsoleOutput() *tview.TextView {
	return tview.NewTextView()
}

func newConsoleInput(client *client.GRPCClient) *tview.InputField {
	in := tview.NewInputField()
	return in.
		SetDoneFunc(func(key tcell.Key) {
			client.Send(in.GetText())
			in.SetText("")
		}).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetPlaceholder("Enter command here...")
}
