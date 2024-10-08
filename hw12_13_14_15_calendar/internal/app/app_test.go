package app_test

import (
	"context"
	"os"
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/app"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/storecreator"
	"github.com/stretchr/testify/suite"
)

type SuiteTest struct {
	suite.Suite
	calendar app.App
	logg     logger.Logger
	db       storage.Storage
}

// SetupTest выполняется перед каждым тестом.
func (s *SuiteTest) SetupTest() {
	ctx := context.Background()

	s.logg = logger.New("INFO", os.Stdout)

	config := storage.Config{
		Host:     "localhost",
		Port:     "8080",
		Inmemory: true,
		Driver:   "postgres",
		Ssl:      "disable",
		Database: "postgres",
		User:     "postgres",
		Password: "postgres",
	}
	s.db, _ = storecreator.New(ctx, &config)

	s.calendar = app.New(s.logg, s.db)

	_ = s.calendar.DeleteAllEvent(ctx)
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()
	_ = s.calendar.DeleteAllEvent(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) NewCommonEvent() storage.Event {
	var eventStart = time.Now().Add(time.Hour * 2) //nolint:gofumpt
	var eventStop = eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return storage.Event{
		ID:           0,
		Title:        "some event",
		Start:        eventStart,
		Stop:         eventStop,
		Description:  "the event",
		UserID:       1,
		Notification: &notification,
	}
}

func (s *SuiteTest) AddEvent(event storage.Event) (int, error) {
	ctx := context.Background()
	id, err := s.calendar.CreateEvent(
		ctx,
		event.UserID,
		event.Title,
		event.Description,
		event.Start,
		event.Stop,
		event.Notification,
	)
	return id, err
}

func (s *SuiteTest) GetAll() []storage.Event {
	ctx := context.Background()
	data, err := s.calendar.ListAllEvent(ctx)
	s.Require().NoError(err)
	return data
}

func (s *SuiteTest) EqualEvents(event1, event2 storage.Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.Unix(), event2.Start.Unix())
	s.Require().Equal(event1.Stop.Unix(), event2.Stop.Unix())
	s.Require().Equal(event1.UserID, event2.UserID)
	s.Require().Equal(event1.Notification, event2.Notification)
}
