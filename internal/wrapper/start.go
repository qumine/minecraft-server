package wrapper

import (
	"context"
	"io/ioutil"
	"os"
	"regexp"
	"sync"

	"github.com/sirupsen/logrus"
)

var logToStatus = map[string]*regexp.Regexp{
	"Starting": regexp.MustCompile(`Starting minecraft server version (.*)`),
	"Started":  regexp.MustCompile(`Done (?s)(.*)! For help, type "help"`),
	"Stopping": regexp.MustCompile(`Stopping (.*) server`),
}

// Start starts the wrapper and the minecraft server.
func (w *Wrapper) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.Info("starting wrapper")

	logrus.Info("writing eula.txt")
	ioutil.WriteFile("./eula.txt", []byte("eula=true"), os.ModeAppend)

	go func() {
		wg.Add(1)
		if err := w.cmd.Start(); err != nil {
			logrus.WithError(err).Fatal("starting wrapper failed")
		}
	}()
	go func() {
		w.Console.Subscribe("wrapper", w.handleLog)
		w.Console.Start()
	}()

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
