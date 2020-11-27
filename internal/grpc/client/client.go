package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	qugrpc "github.com/qumine/qumine-server-java/internal/grpc"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCClient represents the grpc server
type GRPCClient struct {
	Addr string
}

// NewClient creates a new GRPCServer instance
func NewClient() *GRPCClient {
	return &GRPCClient{
		Addr: net.JoinHostPort(utils.GetEnvString("GRPC_ADDR", "127.0.0.1"), utils.GetEnvString("GRPC_PORT", "8081")),
	}
}

// Start the GRPCServer
func (c *GRPCClient) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.WithField("addr", c.Addr).Info("starting grpc client")

	conn, err := grpc.Dial(
		c.Addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		logrus.WithError(err).Fatal("starting grpc client failed")
	}

	client := qugrpc.NewQuMineServerClient(conn)
	go c.streamLogs(ctx, client, -1)

	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			c.Stop(wg)
			return
		default:
			line, _, _ := reader.ReadLine()
			client.SendCommand(ctx, &qugrpc.SendCommandRequest{Line: string(line)})
		}
	}
}

// Stop the grpc client
func (c *GRPCClient) Stop(wg *sync.WaitGroup) {
	logrus.Info("stopping grpc client")
	logrus.Info("stopped grpc client")
	wg.Done()
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
