package storage

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Inmemory bool   `mapstructure:"inmemory"`
	Driver   string `mapstructure:"driver"`
	Ssl      string `mapstructure:"ssl"`
	Database string `mapstructure:"db"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
}
