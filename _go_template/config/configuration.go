package config

import (
	"log"
	"strings"
	"sync"

	"github.com/aryanugroho//pkg/db"

	"github.com/spf13/viper"
)

type AppConfiguration struct {
	ENV              string
	DBDebug          bool
	ApiPrefix        string
	ApiKey           string
	SuffixForTracing string
	Version          string
}

type ServerConfiguration struct {
	Port int
}

type DatabaseConfiguration struct {
	Master db.Config
	Slave  db.Config
}

// Configuration config
type Configuration struct {
	App         AppConfiguration
	Server      ServerConfiguration
	Database    DatabaseConfiguration
	Sentry      SentryConfiguration
	ExternalAPI ExternalAPIConfiguration
	ApiKey      ApiKeyConfiguration
	GCloud      GCloudConfig
	Consumers   ConsumerConfigs
	Publishers  PublisherConfigs
	Cron        CronConfiguration
}

type ExternalAPIConfiguration struct {
	Capt CommonExternalAPI
	DT   CommonExternalAPI
}

type CommonExternalAPI struct {
	BaseURL string
	ApiKey  string
}
type SentryConfiguration struct {
	DSN *string
}

type ApiKeyConfiguration struct {
	OPPA string
	CAPT string
}

type GCloudConfig struct {
	ProjectID string
}

type ConsumerConfigs struct {
	Enable              bool
	DecisionSession     ConsumerConfig
	DecisionTransaction ConsumerConfig
}

type PublisherConfigs struct {
	DecisionSession     PublisherConfig
	DecisionTransaction PublisherConfig
}

type ConsumerConfig struct {
	Topic                  string
	Subscription           string
	MaxOutstandingMessages int
	NumGoroutines          int
	Toggle                 bool
}

type PublisherConfig struct {
	Topic  string
	Toggle bool
}

type CronConfiguration struct {
	BlacklistDevice CronJob
}

type CronJob struct {
	Toggle   bool
	Interval string
	Limit    int
}

var (
	configuration *Configuration
	once          sync.Once
)

// All get all config
func All(opts ...Option) *Configuration {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AllowEmptyEnv(true)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&configuration); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}

		for _, opt := range opts {
			opt(configuration)
		}
	})

	return configuration
}

func Get() *Configuration {
	return configuration
}

func GetSuffixForTracing() string {
	if configuration == nil {
		return ""
	}

	return configuration.App.SuffixForTracing
}

// for testing purpose
func Set(conf *Configuration) {
	configuration = conf
}
