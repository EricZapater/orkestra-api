package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost  string `env:"DB_HOST" envDefault:"localhost"`
	DBPort  string `env:"DB_PORT" envDefault:"5432"`
	DBUser  string `env:"DB_USER" envDefault:"postgres"`
	DBPass  string `env:"DB_PASS" envDefault:"postgres"`
	DBName  string `env:"DB_NAME" envDefault:"postgres"`
	ApiPort string `env:"API_PORT" envDefault:"8080"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"abcd1234"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
		return nil, err
	}
	return &cfg, nil
}