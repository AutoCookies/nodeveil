package config

import "os"

type Config struct {
	HTTPAddr string
	DBPath   string
}

func Load() Config {
	return Config{
		HTTPAddr: getenv("NODEVEIL_HTTP_ADDR", "127.0.0.1:8080"),
		DBPath:   getenv("NODEVEIL_DB_PATH", DefaultDBPath()),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
