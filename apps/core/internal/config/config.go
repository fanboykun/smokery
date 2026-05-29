package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port           string `mapstructure:"PORT"`
	DBAdapter      string `mapstructure:"DB_ADAPTER"`      // sqlite or postgres
	StorageAdapter string `mapstructure:"STORAGE_ADAPTER"` // fs or minio

	// SQLite
	SQLitePath string `mapstructure:"SQLITE_PATH"`

	// PostgreSQL
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	// MinIO
	MinioEndpoint  string `mapstructure:"MINIO_ENDPOINT"`
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`

	// Filesystem storage
	StoragePath string `mapstructure:"STORAGE_PATH"`

	RetentionDays int `mapstructure:"RETENTION_DAYS"`
}

func Load() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_ADAPTER", "sqlite")
	viper.SetDefault("STORAGE_ADAPTER", "fs")
	viper.SetDefault("SQLITE_PATH", "data/smokery.db")
	viper.SetDefault("DATABASE_URL", "postgres://smokery:smokery@localhost:5432/smokery?sslmode=disable")
	viper.SetDefault("MINIO_ENDPOINT", "localhost:9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("STORAGE_PATH", "data/artifacts")
	viper.SetDefault("RETENTION_DAYS", "30")
	viper.AutomaticEnv()
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
