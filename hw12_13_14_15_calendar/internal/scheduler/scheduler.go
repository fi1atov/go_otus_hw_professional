package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/producer"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

type App struct {
	logger   logger.Logger
	duration time.Duration
	p        producer.Producer
	storage  storage.Storage
}

func New(logger logger.Logger, p producer.Producer, storage storage.Storage, duration time.Duration) *App {
	return &App{
		logger:   logger,
		p:        p,
		storage:  storage,
		duration: duration,
	}
}

func (a *App) Publish(ctx context.Context, data []byte) error {
	return a.p.Publish(ctx, data)
}

func (a *App) PublishEvents(ctx context.Context) {
	events, err := a.storage.GetEventsReminder(ctx)
	if err != nil {
		a.logger.Error("fail get events", zap.Error(err))
	}
	for _, event := range events {
		var notification sql.NullInt64
		a.logger.Info(fmt.Sprintf("publish event %v, %s for user %v", event.ID, event.Title, event.UserID))
		encoded, err := json.Marshal(event)
		if err != nil {
			a.logger.Error("fail marshal event", zap.Error(err))
			return
		}
		if err := a.Publish(ctx, encoded); err != nil {
			a.logger.Error("fail publish event", zap.Error(err))
		}
		event.Notification = (*time.Duration)(&notification.Int64)
		err = a.storage.UpdateEvent(ctx, event.ID, event)
		if err != nil {
			a.logger.Error("fail update event", zap.Error(err))
			return
		}
	}
}

func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.duration)
	defer ticker.Stop()
	a.logger.Info("Starting scheduler")
	defer a.logger.Info("Stopping scheduler")
	for {
		go func() {
			a.PublishEvents(ctx)
			a.logger.Info(fmt.Sprintf("Next operation start after %s", a.duration))
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
