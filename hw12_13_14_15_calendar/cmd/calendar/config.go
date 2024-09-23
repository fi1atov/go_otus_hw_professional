package main

import (
	"strings"

	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/server/http"
	"github.com/fi1atov/go_otus_hw_professional/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/viper"
)

type Config struct {
	Logger *logger.Config `mapstructure:"logger"`
	App    *AppConf       `mapstructure:"app"`
}

type AppConf struct {
	GRPCServer *internalgrpc.Config `mapstructure:"grpc_server"`
	HTTPServer *internalhttp.Config `mapstructure:"http_server"`
	Database   *storage.Config      `mapstructure:"database"`
}

func NewConfig(configFile string) (Config, error) {
	conf := Config{}
	_, err := ReadConfigFile(configFile, "toml", &conf)
	return conf, err
}

func ReadConfigFile(pathToFile, typeFile string, configuration interface{}) (interface{}, error) {
	viper.SetConfigFile(pathToFile)
	viper.SetConfigType(typeFile)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("logger.level", "INFO")

	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.httpPort", "8080")
	viper.SetDefault("server.grpcPort", "8081")

	viper.SetDefault("database.inmemory", true)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}
