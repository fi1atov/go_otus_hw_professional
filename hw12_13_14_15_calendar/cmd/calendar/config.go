package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	Database DatabaseConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host     string
	HTTPPort string
	GRPCPort string
}

type DatabaseConf struct {
	Inmemory bool
	Connect  string
}

func NewConfig(configFile string) (Config, error) {
	config := Config{}

	// viper спсобен прочитать конфиг-файл
	v := viper.New()
	// нужно для настройки viper - чтобы все нашел
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("logger.level", "INFO")

	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.httpPort", "8080")
	v.SetDefault("server.grpcPort", "8081")

	v.SetDefault("database.inmemory", true)

	if configFile != "" {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return config, fmt.Errorf("failed to read configuration: %w", err)
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return config, nil
}
