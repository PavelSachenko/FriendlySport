package config

import (
	"github.com/spf13/viper"
	"log"
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

	var sever server
	err = viper.Unmarshal(&sever)
	if err != nil {
		return err
	}
	cfg.Server = sever

	var db db
	err = viper.Unmarshal(&db)
	if err != nil {
		return err
	}
	cfg.DB = db

	var auth auth
	err = viper.Unmarshal(&auth)
	if err != nil {
		return err
	}
	cfg.Auth = auth

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	return nil
}
