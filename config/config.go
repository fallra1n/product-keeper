package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Postgres struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	User     string        `yaml:"user"`
	Password string        `yaml:"password"`
	DBName   string        `yaml:"dbname"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
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
