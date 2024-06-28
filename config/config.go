package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env"`
	HttpServer HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Addr        string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("%w", err)
	}

	configPath := os.Getenv("CONFIG_PATH_HTTP")
	if configPath == "" {
		log.Fatal("CONFIG_PATH_HTTP is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
