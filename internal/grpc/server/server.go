package server

import (
	"context"
	"net"
	"sync"

	qugrpc "github.com/qumine/qumine-server-java/internal/grpc"
	"github.com/qumine/qumine-server-java/internal/utils"
	"github.com/qumine/qumine-server-java/internal/wrapper"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// GRPCServer represents the grpc server
type GRPCServer struct {
	Wrapper *wrapper.Wrapper
	Addr    string

	grpcServer *grpc.Server
}

// NewServer creates a new GRPCServer instance
func NewServer(w *wrapper.Wrapper) *GRPCServer {
	return &GRPCServer{
		Wrapper:    w,
		Addr:       net.JoinHostPort(utils.GetEnvString("GRPC_ADDR", "127.0.0.1"), utils.GetEnvString("GRPC_PORT", "8081")),
		grpcServer: grpc.NewServer(),
	}
}

// Start the GRPCServer
func (s *GRPCServer) Start(ctx context.Context, wg *sync.WaitGroup) {
	logrus.WithField("addr", s.Addr).Info("starting grpc server")

	c, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logrus.WithError(err).Fatal("starting grpc server failed")
	}

	qugrpc.RegisterQuMineServerServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	go func() {
		wg.Add(1)
		s.grpcServer.Serve(c)
	}()
	logrus.Info("started grpc server")

	for {
		select {
		case <-ctx.Done():
			s.Stop(wg)
			return
		}
	}
}

// Stop the api
func (s *GRPCServer) Stop(wg *sync.WaitGroup) {
	logrus.Info("stopping grpc server")
	s.grpcServer.GracefulStop()
	logrus.Info("stopped grpc server")
	wg.Done()
}

func (s *GRPCServer) StreamLogs(req *qugrpc.LogStreamRequest, srv qugrpc.QuMineServer_StreamLogsServer) error {
	s.Wrapper.Console.Subscribe("client", func(line string) {
		srv.Send(&qugrpc.LogStreamResponse{
			Line: line,
		})
	})
	for {
		select {
		case <-srv.Context().Done():
			s.Wrapper.Console.Unsubscribe("client")
			return nil
		}
	}
}

func (s *GRPCServer) SendCommand(req *qugrpc.SendCommandRequest, srv qugrpc.QuMineServer_SendCommandServer) error {
	if err := s.Wrapper.Console.SendCommand(req.Line); err != nil {
		srv.Send(&qugrpc.SendCommandResponse{})
		return nil
	} else {
		return status.Errorf(codes.Internal, "method StreamLogs not implemented")
	}
}
