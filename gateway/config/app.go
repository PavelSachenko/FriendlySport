package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

type Config struct {
	Port         string        `mapstructure:"server_port"`
	Host         string        `mapstructure:"server_host"`
	ReadTimeout  time.Duration `mapstructure:"server_read_timeout"`
	WriteTimeout time.Duration `mapstructure:"server_write_timeout"`

	UserServiceUrl    string `mapstructure:"user_service_url"`
	WorkoutServiceUrl string `mapstructure:"workout_service_url"`
	DietServiceUrl    string `mapstructure:"diet_service_url"`
	AimServiceUrl     string `mapstructure:"aim_service_url"`
	InviterServiceUrl string `mapstructure:"invite_service_url"`
}

var (
	instance *Config
	once     sync.Once
)

func InitConfig(logger *logrus.Logger) (error, *Config) {
	logger.Info("init config")
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

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}
