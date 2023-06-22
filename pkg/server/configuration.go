package server

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Configuration struct {
	Rcp          RcpConfig        `yaml:"RCP"`
	LogConfig    shared.LogConfig `yaml:"LogConfig"`
	StorePath    string           `yaml:"StorePath"`
	IgnoredFiles []string         `yaml:"IgnoredFiles"`
}

type RcpConfig struct {
	Port int16 `yaml:"Port"`
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
