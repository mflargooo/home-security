package config

import (
	"fmt"
	"os"
)

type Config struct {
	Postgres  PostgresConfig
	Redis     RedisConfig
	MediaMTX  MediaMTXConfig
	Directory DirectoryConfig
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

type MediaMTXConfig struct {
	Host string
	Port string
}

type DirectoryConfig struct {
	Path string
}

func (cfg PostgresConfig) FormatDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func (cfg RedisConfig) FormatDSN() string {
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func (cfg MediaMTXConfig) FormatDSN() string {
	return fmt.Sprintf("rtsp://%s:%s", cfg.Host, cfg.Port)
}

func (cfg *Config) StreamLiveURL(cameraID string) string {
	return fmt.Sprintf("rtsp://%s:%s/live/%s", cfg.MediaMTX.Host, cfg.MediaMTX.Port, cameraID)
}

func (cfg *Config) BufferLivePath(cameraID string) string {
	return fmt.Sprintf("/buffer/live/%s", cameraID)
}

func (cfg *Config) BufferSnapshotPath(cameraID string, sessionID string) string {
	return fmt.Sprintf("/buffer/snapshot/%s-session-%s", cameraID, sessionID)
}

func Load() Config {
	return Config{
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
		MediaMTX: MediaMTXConfig{
			Host: os.Getenv("MEDIAMTX_HOST"),
			Port: os.Getenv("MEDIAMTX_PORT"),
		},
		Directory: DirectoryConfig{
			Path: os.Getenv("APP_DIR"),
		},
	}
}
