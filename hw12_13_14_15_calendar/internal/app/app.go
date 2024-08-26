package app

import (
	"context"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type app struct {
	logger  logger.Logger
	storage storage.Storage
}

func (a *app) CreateEvent(ctx context.Context, userID int, title, desc string, start, stop time.Time, notif *time.Duration) (id int, err error) {
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
	isBusy, err := a.storage.IsTimeBusyEvent(ctx, userID, start, stop, 0)
	if err != nil {
		return
	}
	if isBusy {
		err = ErrDateBusy
		return
	}

	return a.storage.CreateEvent(ctx, storage.Event{
		Title:        title,
		Start:        start,
		Stop:         stop,
		Description:  desc,
		UserID:       userID,
		Notification: notif,
	})
}

func (a *app) UpdateEvent(ctx context.Context, id int, change storage.Event) error {
	if change.Title == "" {
		return ErrEmptyTitle
	}
	if change.Start.After(change.Stop) {
		change.Start, change.Stop = change.Stop, change.Start
	}
	if time.Now().After(change.Start) {
		return ErrStartInPast
	}
	isBusy, err := a.storage.IsTimeBusyEvent(ctx, change.UserID, change.Start, change.Stop, id)
	if err != nil {
		return err
	}
	if isBusy {
		return ErrDateBusy
	}

	return a.storage.UpdateEvent(ctx, id, change)
}

func (a *app) DeleteEvent(ctx context.Context, id int) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *app) ListAllEvent(ctx context.Context) ([]storage.Event, error) {
	return a.storage.ListAllEvent(ctx)
}

func (a *app) ListDayEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListDayEvent(ctx, date)
}

func (a *app) ListWeekEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListWeekEvent(ctx, date)
}

func (a *app) ListMonthEvent(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListMonthEvent(ctx, date)
}
