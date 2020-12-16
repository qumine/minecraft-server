package wrapper

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
)

var logToStatus = map[string]*regexp.Regexp{
	"Starting": regexp.MustCompile(`Starting minecraft server version (.*)`),
	"Started":  regexp.MustCompile(`Done (?s)(.*)! For help, type "help"`),
	"Stopping": regexp.MustCompile(`Stopping (.*) server`),
	// Closing Server
}

// Start starts the wrapper and the minecraft server.
func (w *Wrapper) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.Info("starting wrapper")

	logrus.Info("writing eula.txt")
	utils.WriteFileAsString("eula.txt", "eula=true")

	go func() {
		wg.Add(1)
		if err := w.cmd.Start(); err != nil {
			logrus.WithError(err).Fatal("starting wrapper failed")
		}
		w.cmdKeepRunning = true
	}()
	go func() {
		w.Console.Subscribe("wrapper", w.handleLog)
		w.Console.Start()
	}()
	go func() {
		for {
			if w.cmdKeepRunning {
				err := w.cmd.Wait()
				if err != nil {
					logrus.WithError(err).Fatal("java process stopped unexpectedly")
				}
			}
			time.Sleep(time.Second)
		}
	}()
	logrus.Info("started wrapper")

	for {
		select {
		case <-ctx.Done():
			w.Stop(wg)
			return
		}
	}
}

func (w *Wrapper) handleLog(line string) {
	for status, reg := range logToStatus {
		if reg.MatchString(line) {
			w.Status = status
		}
	}
}
