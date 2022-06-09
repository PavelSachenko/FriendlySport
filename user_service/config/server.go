package config

import "time"

type server struct {
	Port         string        `mapstructure:"server_port"`
	Host         string        `mapstructure:"server_host"`
	ReadTimeout  time.Duration `mapstructure:"server_read_timeout"`
	WriteTimeout time.Duration `mapstructure:"server_write_timeout"`
}
