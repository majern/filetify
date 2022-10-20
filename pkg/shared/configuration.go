package shared

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Configuration struct {
	Paths          []string  `yaml:"Paths"`
	SyncTimeoutSec int       `yaml:"SyncTimeoutSec"`
	LogConfig      LogConfig `yaml:"LogConfig"`
	IgnoredFiles   []string  `yaml:"IgnoredFiles"`
}

type LogConfig struct {
	DetailedLogs     bool `yaml:"DetailedLogs"`
	UseJsonFormatter bool `yaml:"UseJsonFormatter"`
}

var configInstanceLock sync.Once
var configInstance *Configuration

func GetConfiguration() *Configuration {
	configInstanceLock.Do(func() {
		err := viper.Unmarshal(&configInstance)
		if err != nil {
			log.Fatalf("Unable to parse configuration file. Error: %v", err)
		}
	})

	return configInstance
}
