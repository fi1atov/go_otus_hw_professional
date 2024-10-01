package sqlstorage

import (
	"fmt"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
)

func New(conf *storage.Config) storage.Storage {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Ssl,
	)

	return &store{
		dataSourceName: dsn,
	}
}
