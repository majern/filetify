package client

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Configuration struct {
	Paths          []string         `yaml:"Paths"`
	SyncTimeoutSec int              `yaml:"SyncTimeoutSec"`
	LogConfig      shared.LogConfig `yaml:"LogConfig"`
	IgnoredFiles   []string         `yaml:"IgnoredFiles"`
	Rcp            RcpConfig        `yaml:"RCP"`
}

type RcpConfig struct {
	Endpoint string `yaml:"Endpoint"`
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
