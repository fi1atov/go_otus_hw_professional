package internalhttp

import (
	"context"
	"errors"
	"time"

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

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration `json:"notification,omitempty"`
}

type CreateResult struct {
	ID int
}

type OkResult struct {
	Ok bool
}

type DeleteRequest struct {
	ID int
}

type ListRequest struct {
	Date time.Time
}

type ListResult []Event

func (event Event) Validate() error {
	if event.UserID < 1 {
		return errors.New("user ID не может быть меньше 1")
	}

	return nil
}
