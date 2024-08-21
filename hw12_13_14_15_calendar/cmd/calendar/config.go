package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	HTTP     HTTPConf
	Database DatabaseConf
}

type LoggerConf struct {
	Level string
	File  string
}

type HTTPConf struct {
	Host string
	Port string
}

type DatabaseConf struct {
	Inmem   bool
	Connect string
}

func NewConfig(configFile string) (Config, error) {
	config := Config{}

	// viper спсобен прочитать конфиг-файл
	v := viper.New()
	// нужно для настройки viper - чтобы все нашел
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("logger.level", "INFO")

	v.SetDefault("http.host", "127.0.0.1")
	v.SetDefault("http.port", "8080")

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
