package config

import (
	"flag"
	"log"
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
	configPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()

	var cfg Config

	err := cleanenv.ReadConfig(*configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
