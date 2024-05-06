package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("empty config path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path doesn't exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	return os.Getenv("CONFIG_PATH")
}
