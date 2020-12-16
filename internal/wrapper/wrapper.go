package wrapper

import (
	"os/exec"
	"syscall"
	"time"
)

// Wrapper represents the wrapper object of the minecraft server
type Wrapper struct {
	Status  string
	Console *Console

	cmd            *exec.Cmd
	cmdStopTimeout time.Duration
	cmdKeepRunning bool
}

// NewWrapper creates a new wrapper
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
		Console: NewConsole(stdin, stderr, stdout),

		cmd:            cmd,
		cmdStopTimeout: 15 * time.Second,
	}
}
