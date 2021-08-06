package wrapper

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var logToStatus = map[*regexp.Regexp]string{
	// Vanilla based servers.
	regexp.MustCompile(`Starting minecraft server version (.*)`): "Starting",
	regexp.MustCompile(`Done (?s)(.*)! For help, type "help"`):   "Started",
	regexp.MustCompile(`Stopping (.*) server`):                   "Stopping",

	// Bungee based servers.
	regexp.MustCompile(`Listening on .*`):     "Started",
	regexp.MustCompile(`Closing listener .*`): "Stopping",
	// Closing Server
}

// Start starts the wrapper and the minecraft server.
func (w *Wrapper) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.WithFields(logrus.Fields{
		"path": w.cmd.Path,
		"args": w.cmd.Args,
	}).Debug("starting wrapper")
	wg.Add(1)

	go w.startConsole()
	go w.startCommand()
	go w.startWatchdog()

	logrus.WithFields(logrus.Fields{
		"path": w.cmd.Path,
		"args": w.cmd.Args,
	}).Info("started wrapper")
	for {
		select {
		case <-ctx.Done():
			w.Stop(wg)
			return
		}
	}
}

func (w *Wrapper) startConsole() {
	w.Console.Subscribe("wrapper", w.onLog)
	w.Console.Start()
}

func (w *Wrapper) startCommand() {
	if err := w.cmd.Start(); err != nil {
		logrus.WithError(err).Fatal("starting wrapper failed")
	}
	w.cmdKeepRunning = true
}

func (w *Wrapper) startWatchdog() {
	for {
		time.Sleep(time.Second)
		if err := w.cmd.Wait(); err != nil {
			w.cmdExited = true
			if w.cmdKeepRunning {
				logrus.Fatal("java process stopped unexpectedly")
			}
		}
	}
}
