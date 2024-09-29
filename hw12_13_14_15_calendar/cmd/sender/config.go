package main

import (
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/config"
	internalrmqconsumer "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/mq/consumer/rmq"
)

type Config struct {
	Consumer *internalrmqconsumer.Config `mapstructure:"rmq"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := config.ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}
