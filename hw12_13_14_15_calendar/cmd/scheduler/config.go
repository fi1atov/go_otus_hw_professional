package main

import (
	"time"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/config"
	internalrmqproducer "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/producer/rmq"
)

type DatabaseConf struct {
	Inmemory bool
	Connect  string
}

type Config struct {
	Database DatabaseConf                `mapstructure:"database"`
	Producer *internalrmqproducer.Config `mapstructure:"rmq"`
	Duration time.Duration               `mapstructure:"duration"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
