package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN       string `env:"DB_DSN" `
	RedisConfig RedisConfig
	ServerPort  string `env:"SERVER_PORT"`
}

type RedisConfig struct {
	// 基础连接信息
	Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`

	// 超时设置 (配置库会自动解析 "10s", "500ms" 等字符串为 time.Duration)
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`

	// 连接池设置
	PoolSize      int  `env:"REDIS_POOL_SIZE" envDefault:"10"`
	MinIdleConns  int  `env:"REDIS_MIN_IDLE_CONNS" envDefault:"5"`
	IsClusterMode bool `env:"REDIS_IS_CLUSTER_MODE" envDefault:"false"`
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
