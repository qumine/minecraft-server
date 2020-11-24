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
	logrus.WithField("timeout", w.cmdStopTimeout).Info("stopping wrapper")

	ctx, cancel := context.WithTimeout(context.Background(), w.cmdStopTimeout)
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.cmd.Process.Kill()
				return
			}
		}
	}()
	go func() {
		for i := 5; i > 0; i-- {
			w.Console.SendCommand(fmt.Sprintf("say Stopping server in: %d", i))
			time.Sleep(time.Second)
		}
		w.Console.SendCommand("stop")

		if err := w.cmd.Wait(); err != nil {
			logrus.WithError(err).Error("stopping wrapper failed")
		}
		cancel()
		logrus.Info("stopped wrapper")
		wg.Done()
	}()
}
