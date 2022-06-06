package console

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

// Start will start the console, it will follow the logs
func (c *Console) Start(ctx context.Context) {
	go c.startStderr()
	go c.startStdout()
	<-ctx.Done()
}

func (c *Console) startStderr() {
	for {
		if l, err := readLine(c.stderr); err == nil {
			go func() {
				for _, s := range c.subscriber {
					s(l)
				}
			}()
			fmt.Print(l)
		}
	}
}

func (c *Console) startStdout() {
	for {
		if l, err := readLine(c.stdout); err == nil {
			go func() {
				for _, s := range c.subscriber {
					s(l)
				}
			}()
			fmt.Print(l)
		}
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	if reader == nil {
		return "", io.EOF
	}
	l, err := reader.ReadString('\n')
	if err == io.EOF {
		return "", io.EOF
	}
	return l, nil
}
