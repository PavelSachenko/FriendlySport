package config

import (
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Server               server
	DB                   db
	Auth                 auth
	UserPasswordHashSalt string `mapstructure:"user_password_hash_salt"`
}

var (
	instance *Config
	once     sync.Once
)

func InitConfig() (error, *Config) {
	var err error
	once.Do(func() {
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

	err = viper.Unmarshal(&cfg.Auth)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}
