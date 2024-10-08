package storecreator

import (
	"context"
	"fmt"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(ctx context.Context, cfg *storage.Config) (storage.Storage, error) {
	var db storage.Storage
	if cfg.Inmemory {
		fmt.Println("using inmem storage")
		db = memorystorage.New()
	} else {
		fmt.Println("using sql storage")
		db = sqlstorage.New(cfg)
	}
	err := db.Connect(ctx)
	return db, err
}
