package memorystorage

import (
	"context"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type data map[int]storage.Event

type Storage struct {
	mu     sync.RWMutex //nolint:unused
	lastID int
	data   data
}

func New() *Storage {
	result := Storage{}
	result.data = make(data)
	return &result
}

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) Create(_ context.Context, event storage.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.newID()
	event.ID = id
	s.data[id] = storage.Event{
		ID:           id,
		Title:        event.Title,
		Start:        event.Start,
		Stop:         event.Stop,
		Description:  event.Description,
		UserID:       event.UserID,
		Notification: event.Notification,
	}
	return id, nil
}

func (s *Storage) Update(_ context.Context, id int, change storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.data[id]
	if !ok {
		return storage.ErrNotExistsEvent
	}

	event.Title = change.Title
	event.Start = change.Start
	event.Stop = change.Stop
	event.Description = change.Description
	event.Notification = change.Notification
	s.data[id] = event

	return nil
}

func (s *Storage) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, id)
	return nil
}

func (s *Storage) DeleteAll(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(data)
	return nil
}

func (s *Storage) ListAll(_ context.Context) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]storage.Event, 0, len(s.data))
	for _, event := range s.data {
		result = append(result, event)
	}
	return result, nil
}

func (s *Storage) ListDay(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, day := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, eventDay := event.Start.Date()
		if eventYear == year && eventMonth == month && eventDay == day {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListWeek(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, week := date.ISOWeek()
	for _, event := range s.data {
		eventYear, eventWeek := event.Start.ISOWeek()
		if eventYear == year && eventWeek == week {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) ListMonth(_ context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []storage.Event
	year, month, _ := date.Date()
	for _, event := range s.data {
		eventYear, eventMonth, _ := event.Start.Date()
		if eventYear == year && eventMonth == month {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) IsTimeBusy(_ context.Context, userID int, start, stop time.Time, excludeID int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range s.data {
		if event.UserID == userID && event.ID != excludeID && event.Start.Before(stop) && event.Stop.After(start) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) newID() int {
	s.lastID++
	return s.lastID
}
