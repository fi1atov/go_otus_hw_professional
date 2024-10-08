//go:generate protoc -I "/usr/local/include/" --proto_path=../../../api/ --go_out=. --go-grpc_out=. ../../../api/EventService.proto
package grpcserver

import (
	"context"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
)

type Server interface {
	Start(addr string) error
	Stop(ctx context.Context) error
}

func NewServer(app app.App, logger logger.Logger) Server {
	return newServer(app, logger)
}
