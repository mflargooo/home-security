package config

import (
	"fmt"
	"os"
)

type Config struct {
	Postgres   PostgresConfig
	Redis      RedisConfig
	ListenAddr string
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type RedisConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func (cfg PostgresConfig) FormatDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func (cfg RedisConfig) FormatDSN() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func Load() *Config {
	return &Config{
		Postgres: PostgresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DB"),
		},
		Redis: RedisConfig{
			Username: os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Database: os.Getenv("REDIS_DB"),
		},
		ListenAddr: os.Getenv("LISTEN_ADDR"),
	}
}
