package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BaseURL      string `yaml:"base_url"`
	DBConfig     `yaml:"db"`
	GRPCConfig   `yaml:"grpc"`
	MailerConfig `yaml:"mailer"`
}

type DBConfig struct {
	StoragePath string `yaml:"db_url" env-required:"true"`
}

type GRPCConfig struct {
	AuthAddr string `yaml:"auth_addr" env-required:"true"`
}

type MailerConfig struct {
	SMTPFrom     string `yaml:"from" env-required:"true"`
	SMTPPassword string `yaml:"password" env-required:"true"`
	SMTPHost     string `yaml:"host" env-required:"true"`
	SMTPPort     string `yaml:"port" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("env variable CONFIG_PATH is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}
