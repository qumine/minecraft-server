package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jroimartin/gocui"
	qugrpc "github.com/qumine/qumine-server-java/internal/grpc"
	"github.com/qumine/qumine-server-java/internal/utils"
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
func (c *GRPCClient) Start(g *gocui.Gui) error {
	var err error
	time.Sleep(2 * time.Second)
	logs, err := g.View("logs")
	if err != nil {
		logrus.WithError(err).Fatal("starting grpc client failed")
	}
	go func() {
		c.ctx, c.cancel = context.WithCancel(context.Background())
		c.conn, err = grpc.Dial(c.Addr, grpc.WithInsecure())
		if err != nil {
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(logs, "failed to start grpc client")
				return nil
			})
		}
		c.client = qugrpc.NewQuMineServerClient(c.conn)
		stream, err := c.client.StreamLogs(c.ctx, &qugrpc.LogStreamRequest{Lines: -1})
		if err != nil {
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprintln(logs, "failed to stream logs")
				return nil
			})
			return
		}

		for {
			var rsp qugrpc.LogStreamResponse
			if err = stream.RecvMsg(&rsp); err != nil {
				g.Update(func(g *gocui.Gui) error {
					fmt.Fprintln(logs, "stream closed")
					return nil
				})
				return
			}
			g.Update(func(g *gocui.Gui) error {
				fmt.Fprint(logs, rsp.Line)
				return nil
			})
		}
	}()
	return nil
}

// Stop the grpc client
func (c *GRPCClient) Stop(g *gocui.Gui, v *gocui.View) error {
	logs, err := g.View("logs")
	if err != nil {
		logrus.WithError(err).Fatal("starting grpc client failed")
	}

	g.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(logs, "stopping grpc client")
		return nil
	})
	c.cancel()
	c.conn.Close()
	g.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(logs, "stopped grpc client")
		return nil
	})
	return gocui.ErrQuit
}

// Send sends the command to the server
func (c *GRPCClient) Send(g *gocui.Gui, v *gocui.View) error {
	c.client.SendCommand(c.ctx, &qugrpc.SendCommandRequest{Line: v.Buffer()})
	g.Update(func(g *gocui.Gui) error {
		v.Clear()
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return nil
	})
	return nil
}

func (c *GRPCClient) streamLogs(ctx context.Context, client qugrpc.QuMineServerClient, lines int32) {
	stream, err := client.StreamLogs(ctx, &qugrpc.LogStreamRequest{Lines: lines})
	if err != nil {
		logrus.WithError(err).Fatal("streaming logs failed")
		return
	}

	for {
		var rsp qugrpc.LogStreamResponse
		if err = stream.RecvMsg(&rsp); err != nil {
			logrus.WithError(err).Fatal("stream closed")
			return
		}
		fmt.Print(rsp.Line)
	}
}
