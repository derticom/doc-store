package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel    string     `yaml:"log_level"`
	AdminToken  string     `yaml:"admin_token"`
	PostgresURL string     `yaml:"postgres_url"`
	RedisURL    string     `yaml:"redis_url"`
	Minio       MinIO      `yaml:"minio_server"`
	Server      HTTPServer `yaml:"http_server"`
}

type MinIO struct {
	Address   string `yaml:"address"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	UseSSL    bool   `yaml:"use_ssl"`
}

type HTTPServer struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
