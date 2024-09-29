package config

import (
	"strings"

	"github.com/spf13/viper"
)

func ReadConfigFile(pathToFile, typeFile string, configuration interface{}) (interface{}, error) {
	viper.SetConfigFile(pathToFile)
	viper.SetConfigType(typeFile)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}
	return configuration, nil
}
