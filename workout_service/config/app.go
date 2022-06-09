package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Server server
	DB     db
}

var (
	instance *Config
	once     sync.Once
)

func InitConfig() (error, *Config) {
	var err error
	once.Do(func() {
		log.Printf("Init config")
		instance = &Config{}
		err = instance.unmarshal()
	})
	if err != nil {
		return err, nil
	}

	return nil, instance
}

func (cfg *Config) unmarshal() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg.Server)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg.DB)
	if err != nil {
		return err
	}

	return nil
}
