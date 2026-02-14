package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ServiceName string
	HTTPPort    string
	GRPCPort    string

	// Database
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string

	// Redis
	RedisHost string
	RedisPort string

	// Kafka
	KafkaBrokers string

	// CORS
	AllowedOrigins []string

	// PIM Provider
	PIMProvider string
	PIMURL      string
	PIMAPIKey   string

	// Search Provider
	SearchProvider string
	SearchURL      string
	SearchAPIKey   string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServiceName:      getEnv("SERVICE_NAME", "catalog-service"),
		HTTPPort:         getEnv("HTTP_PORT", "8081"),
		GRPCPort:         getEnv("GRPC_PORT", "9091"),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseName:     getEnv("DATABASE_NAME", "catalog"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "postgres"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
		AllowedOrigins:   getSliceEnv("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		PIMProvider:      getEnv("PIM_PROVIDER", "mock"),
		PIMURL:           getEnv("PIM_URL", ""),
		PIMAPIKey:        getEnv("PIM_API_KEY", ""),
		SearchProvider:   getEnv("SEARCH_PROVIDER", "mock"),
		SearchURL:        getEnv("SEARCH_URL", ""),
		SearchAPIKey:     getEnv("SEARCH_API_KEY", ""),
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
}

func (c *Config) RedisURL() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		var result []string
		for _, s := range strings.Split(value, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				result = append(result, s)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1"
	}
	return defaultValue
}
