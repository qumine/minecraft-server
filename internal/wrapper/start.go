package wrapper

import (
	"context"
	"io/ioutil"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

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
