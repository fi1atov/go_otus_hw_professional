package app

import (
	"context"
	"errors"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"time"
)

var ErrNoUserID = errors.New("no user id of the event")
var ErrEmptyTitle = errors.New("no title of the event")
var ErrStartInPast = errors.New("start time of the event in the past")
var ErrDateBusy = errors.New("this time is already occupied by another event")

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger, storage storage.Storage) App {
	return App{
		logger,
		storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, userID int, title, desc string, start, stop time.Time, notif *time.Duration) (id int, err error) {
	if userID == 0 {
		err = ErrNoUserID
		return
	}
	if title == "" {
		err = ErrEmptyTitle
		return
	}
	if start.After(stop) {
		start, stop = stop, start
	}
	if time.Now().After(start) {
		err = ErrStartInPast
		return
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, userID, start, stop, 0)
	if err != nil {
		return
	}
	if isBusy {
		err = ErrDateBusy
		return
	}

	return a.storage.Create(ctx, storage.Event{
		Title:        title,
		Start:        start,
		Stop:         stop,
		Description:  desc,
		UserID:       userID,
		Notification: notif,
	})
}

func (a *App) UpdateEvent(ctx context.Context, id int, change storage.Event) error {
	if change.Title == "" {
		return ErrEmptyTitle
	}
	if change.Start.After(change.Stop) {
		change.Start, change.Stop = change.Stop, change.Start
	}
	if time.Now().After(change.Start) {
		return ErrStartInPast
	}
	isBusy, err := a.storage.IsTimeBusy(ctx, change.UserID, change.Start, change.Stop, id)
	if err != nil {
		return err
	}
	if isBusy {
		return ErrDateBusy
	}

	return a.storage.Update(ctx, id, change)
}

func (a *App) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.Delete(ctx, id)
}

func (a *App) ListAllEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListAll(ctx)
}

func (a *App) ListDayEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListDay(ctx, date)
}

func (a *App) ListWeekEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListWeek(ctx, date)
}

func (a *App) ListMonthEvents(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListMonth(ctx, date)
}
