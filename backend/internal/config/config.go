package config

import (
	"os"
)

// Config holds server configuration from env (and optional flag overrides).
type Config struct {
	GRPCAddr string
	HTTPAddr string
	DBPath   string
}

// Load reads config from environment; empty values use defaults.
func Load() *Config {
	return &Config{
		GRPCAddr: getEnv("GRPC_ADDR", ":9090"),
		HTTPAddr: getEnv("HTTP_ADDR", ":8100"),
		DBPath:   getEnv("DB_PATH", "scheduler.db"),
	}
}

func getEnv(key, defaultVal string) string {
	if envValue := os.Getenv(key); envValue != "" {
		return envValue
	}
	return defaultVal
}
