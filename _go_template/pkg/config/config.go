package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var c config

type config struct {
	port       string
	address    string
	statsdHost string
	statsdPort string
}

func Load(path string) error {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
	c = config{
		port:       viper.GetString("APP_PORT"),
		address:    viper.GetString("APP_ADDRESS"),
		statsdHost: viper.GetString("STATSD_HOST"),
		statsdPort: viper.GetString("STATSD_PORT"),
	}
	return nil
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.address, c.port)
}

func GetStatsDAddress() string {
	return fmt.Sprintf("%s:%s", c.statsdHost, c.statsdPort)
}
