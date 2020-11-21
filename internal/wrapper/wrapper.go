package wrapper

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Wrapper struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stderr *bufio.Reader
	stdout *bufio.Reader
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
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdin),
		stderr: bufio.NewReader(stderr),
		stdout: bufio.NewReader(stdout),
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
		for {
			line, err := w.stdout.ReadString('\n')
			if err == io.EOF {
				return
			} else if err != nil {
				logrus.WithError(err).Fatal("failed to read from server")
			}
			fmt.Print(line)
		}
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
	logrus.Info("stopping wrapper")
	w.SendCommand("say Stopping server... 3")
	time.Sleep(time.Second)
	w.SendCommand("say Stopping server... 2")
	time.Sleep(time.Second)
	w.SendCommand("say Stopping server... 1")
	time.Sleep(time.Second)
	w.SendCommand("stop")

	if err := w.cmd.Wait(); err != nil {
		logrus.WithError(err).Error("stopping wrapper failed")
	}
	wg.Done()
}

func (w *Wrapper) SendCommand(c string) error {
	logrus.WithField("c", c).Info("sending command")
	if _, err := w.stdin.WriteString(fmt.Sprintf("%s\r\n", c)); err != nil {
		logrus.WithError(err).Error("failed")
	}
	return w.stdin.Flush()
}

func (w *Wrapper) tailLogs() {
	for {
		line, err := w.stdout.ReadString('\n')
		if err == io.EOF {
			return
		} else if err != nil {
			logrus.WithError(err).Fatal("failed to read from server")
		}
		fmt.Print(line)
	}
}
