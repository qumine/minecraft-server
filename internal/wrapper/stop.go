package wrapper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Stop stops the wrapper by first trying to gravefully shutdown the minecraft server.
func (w *Wrapper) Stop(wg *sync.WaitGroup) {
	logrus.WithFields(logrus.Fields{
		"path":    w.cmd.Path,
		"args":    w.cmd.Args,
		"timeout": w.cmdStopTimeout,
	}).Debug("stopping wrapper")

	w.cmdKeepRunning = false
	ctx, cancel := context.WithTimeout(context.Background(), w.cmdStopTimeout)
	go w.stopWatchdog(ctx, wg)
	go w.stopWrapper(cancel, wg)
}

func (w *Wrapper) stopWatchdog(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			w.cmd.Process.Kill()
			return
		}
	}
}

func (w *Wrapper) stopWrapper(cancel context.CancelFunc, wg *sync.WaitGroup) {
	for i := 5; i > 0; i-- {
		w.Console.SendCommand(fmt.Sprintf("say Stopping server in: %d", i))
		time.Sleep(time.Second)
	}
	w.Console.SendCommand("stop")

	for {
		time.Sleep(time.Second)
		if w.cmdExited {
			break
		}
	}
	cancel()

	logrus.WithFields(logrus.Fields{
		"path":    w.cmd.Path,
		"args":    w.cmd.Args,
		"timeout": w.cmdStopTimeout,
	}).Info("stopped wrapper")
	wg.Done()
}
