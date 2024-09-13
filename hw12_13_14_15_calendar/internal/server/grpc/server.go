package grpcserver

import (
	"context"
	"net"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	app    app.App
	logger logger.Logger
	srv    *grpc.Server
}

func newServer(app app.App, logger logger.Logger) *server {
	s := &server{
		app:    app,
		logger: logger,
	}
	return s
}

// Start protoc --proto_path=api/ --go_out=internal/server/grpc
// --go-grpc_out=internal/server/grpc api/EventService.proto.
func (s *server) Start(addr string) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(s.logger)))
	RegisterCalendarServer(s.srv, NewService(s.app))
	reflection.Register(s.srv)

	s.logger.Info("starting grpc server on " + addr)
	return s.srv.Serve(lsn)
}

func (s *server) Stop(_ context.Context) error {
	s.srv.GracefulStop()
	return nil
}
