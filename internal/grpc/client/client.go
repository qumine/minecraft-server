package client

import (
	"context"
	"fmt"
	"net"
	"time"

	qugrpc "github.com/qumine/minecraft-server/internal/grpc"
	"github.com/qumine/minecraft-server/internal/utils"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCClient represents the grpc server
type GRPCClient struct {
	Addr string

	conn   *grpc.ClientConn
	client qugrpc.QuMineServerClient
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient creates a new GRPCServer instance
func NewClient() *GRPCClient {
	return &GRPCClient{
		Addr: net.JoinHostPort(utils.GetEnvString("GRPC_ADDR", "127.0.0.1"), utils.GetEnvString("GRPC_PORT", "8081")),
	}
}

// Start the GRPCServer
func (c *GRPCClient) Start(app *tview.Application, out *tview.TextView) error {
	var err error
	time.Sleep(2 * time.Second)
	if err != nil {
		logrus.WithError(err).Fatal("starting grpc client failed")
	}
	go func() {
		c.ctx, c.cancel = context.WithCancel(context.Background())
		c.conn, err = grpc.Dial(c.Addr, grpc.WithInsecure())
		if err != nil {
			fmt.Fprintln(out, "failed to start grpc client")
			app.Draw()
			return
		}
		c.client = qugrpc.NewQuMineServerClient(c.conn)
		stream, err := c.client.StreamLogs(c.ctx, &qugrpc.LogStreamRequest{Lines: -1})
		if err != nil {
			fmt.Fprintln(out, "failed to stream logs")
			app.Draw()
			return
		}

		for {
			var rsp qugrpc.LogStreamResponse
			if err = stream.RecvMsg(&rsp); err != nil {
				fmt.Fprintln(out, "stream closed")
				app.Draw()
				return
			}
			out.Write([]byte(rsp.Line))
			out.ScrollToEnd()
			app.Draw()
		}
	}()
	return nil
}

// Stop the grpc client
func (c *GRPCClient) Stop(app *tview.Application, out *tview.TextView) error {
	fmt.Fprintln(out, "CONSOLE: stopping grpc client")
	app.Draw()
	c.cancel()
	if err := c.conn.Close(); err != nil {
		return err
	}

	fmt.Fprintln(out, "CONSOLE: stopping grpc client")
	app.Draw()
	return nil
}

// Send sends the command to the server
func (c *GRPCClient) Send(line string) error {
	c.client.SendCommand(c.ctx, &qugrpc.SendCommandRequest{Line: line})
	return nil
}
