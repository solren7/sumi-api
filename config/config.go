package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN      string `env:"DBDSN" `
	RedisAddr  string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisDB    int    `env:"REDIS_DB"`
	ServerPort string `env:"SERVER_PORT"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
