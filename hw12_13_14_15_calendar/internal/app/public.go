package app

import (
	"context"
	"errors"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
)

type ListEvents func(ctx context.Context, date time.Time) ([]storage.Event, error)

type App interface {
	CreateEvent(
		ctx context.Context,
		userID int,
		title string,
		desc string,
		start time.Time,
		stop time.Time,
		notif *time.Duration,
	) (id int, err error)
	UpdateEvent(ctx context.Context, id int, change storage.Event) error
	DeleteEvent(ctx context.Context, id int) error
	ListAllEvent(ctx context.Context) ([]storage.Event, error)
	ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func New(logger logger.Logger, storage storage.Storage) App {
	return &app{
		logger,
		storage,
	}
}

var (
	ErrNoUserID    = errors.New("no user id of the event")
	ErrEmptyTitle  = errors.New("no title of the event")
	ErrStartInPast = errors.New("start time of the event in the past")
	ErrDateBusy    = errors.New("this time is already occupied by another event")
)
