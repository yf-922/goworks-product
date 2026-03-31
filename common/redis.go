package common

import (
	"os"
	"strconv"
)

const (
	defaultRedisHost = "127.0.0.1"
	defaultRedisPort = 6379
)

type RedisConfig struct {
	Host string
	Port int
}

func NewRedisConfig() RedisConfig {
	port := defaultRedisPort
	if value := os.Getenv("REDIS_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			port = parsed
		}
	}
	return RedisConfig{
		Host: envOrDefault("REDIS_HOST", defaultRedisHost),
		Port: port,
	}
}
