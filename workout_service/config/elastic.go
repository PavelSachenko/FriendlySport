package config

type elastic struct {
	Host     string `mapstructure:"elastic_host"`
	Port     string `mapstructure:"elastic_port"`
	Username string `mapstructure:"elastic_username"`
	Password string `mapstructure:"elastic_password"`
	Network  string `mapstructure:"elastic_network"`
}
