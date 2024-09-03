package sqlstorage

import "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"

func New() storage.Storage {
	return &store{}
}
