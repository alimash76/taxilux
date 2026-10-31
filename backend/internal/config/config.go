package config

import "os"

type Config struct {
	DBPath   string
	Port     string
	LogLevel string
}

func Load() (*Config, error) {
	return &Config{
		DBPath:   getenv("DB_PATH", "/data/taxilux.db"),
		Port:     getenv("PORT", "8080"),
		LogLevel: getenv("LOG_LEVEL", "info"),
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
