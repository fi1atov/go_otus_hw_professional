package memorystorage

import (
	"context"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type data map[int]storage.Event

type Storage struct {
	// TODO
	mu     sync.RWMutex //nolint:unused
	lastID int
	data   data
}

func New() storage.Storage {
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
