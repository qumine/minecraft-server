package wrapper

import (
	"os/exec"
	"syscall"
	"time"

	"github.com/qumine/qumine-server-java/internal/console"
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
