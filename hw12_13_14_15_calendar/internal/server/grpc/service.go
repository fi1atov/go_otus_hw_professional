package grpcserver

import (
	"context"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	UnimplementedCalendarServer

	app app.App
}

func NewService(app app.App) *Service {
	return &Service{
		app: app,
	}
}

func (s *Service) Create(_ context.Context, _ *Event) (*EventCreateResult, error) {
	return nil, status.Error(codes.OK, "okey")
}
