package main

import (
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/config"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	internalrmqproducer "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/producer/rmq"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
)

type Config struct {
	Logger   *logger.Config              `mapstructure:"logger"`
	Database *storage.Config             `mapstructure:"database"`
	Producer *internalrmqproducer.Config `mapstructure:"rmq"`
	Duration time.Duration               `mapstructure:"duration"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
