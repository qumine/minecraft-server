package wrapper

import (
	"os/exec"
	"syscall"
	"time"

	"github.com/qumine/qumine-server-java/internal/server"
	"github.com/qumine/qumine-server-java/internal/wrapper/console"
)

// Wrapper represents the wrapper object of the minecraft server
type Wrapper struct {
	Status  string
	Console *console.Console

	cmd            *exec.Cmd
	cmdStopTimeout time.Duration
	cmdKeepRunning bool
	cmdExited      bool
}

// NewWrapper creates a new wrapper
func NewWrapper(srv server.Server) *Wrapper {
	n, a := srv.StartupCommand()

	cmd := exec.Command(n, a...)
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
