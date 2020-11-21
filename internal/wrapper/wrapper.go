package wrapper

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/qumine/qumine-server-java/internal/console"
	"github.com/sirupsen/logrus"
)

type Wrapper struct {
	Console *console.Console

	cmd            *exec.Cmd
	cmdStopTimeout time.Duration
}

func NewWrapper() *Wrapper {
	cmd := exec.Command("java", "-jar", "server.jar", "nogui")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	stdin, _ := cmd.StdinPipe()
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	return &Wrapper{
		Console: console.NewConsole(stdin, stderr, stdout),

		cmd:            cmd,
		cmdStopTimeout: 15 * time.Second,
	}
}

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
