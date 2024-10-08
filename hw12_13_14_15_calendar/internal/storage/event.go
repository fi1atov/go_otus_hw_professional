package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	Base
	Events
}

type Base interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type Events interface {
	CreateEvent(ctx context.Context, event Event) (int, error)
	UpdateEvent(ctx context.Context, id int, change Event) error
	DeleteEvent(ctx context.Context, id int) error
	DeleteAllEvent(ctx context.Context) error
	ListAllEvent(ctx context.Context) ([]Event, error)
	ListDayEvent(ctx context.Context, date time.Time) ([]Event, error)
	ListWeekEvent(ctx context.Context, date time.Time) ([]Event, error)
	ListMonthEvent(ctx context.Context, date time.Time) ([]Event, error)
	IsTimeBusyEvent(ctx context.Context, userID int, start, stop time.Time, excludeID int) (bool, error)
	GetEventsReminder(ctx context.Context) ([]Event, error)
	DeleteEventsBeforeDate(ctx context.Context, date time.Time) error
}

type Event struct {
	ID           int
	Title        string
	Start        time.Time
	Stop         time.Time
	Description  string
	UserID       int
	Notification *time.Duration
}

func (e *Event) Transfer() Event {
	return Event{
		Title:        e.Title,
		Start:        e.Start,
		Stop:         e.Stop,
		Description:  e.Description,
		UserID:       e.UserID,
		Notification: e.Notification,
	}
}

var ErrNotExistsEvent = errors.New("no such event")
