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

	// JWT
	JWTAccessSecret       string
	JWTRefreshSecret      string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration

	// CORS
	AllowedOrigins []string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServiceName:           getEnv("SERVICE_NAME", "identity-service"),
		HTTPPort:              getEnv("HTTP_PORT", "8080"),
		GRPCPort:              getEnv("GRPC_PORT", "9090"),
		DatabaseHost:          getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:          getEnv("DATABASE_PORT", "5432"),
		DatabaseName:          getEnv("DATABASE_NAME", "identity"),
		DatabaseUser:          getEnv("DATABASE_USER", "postgres"),
		DatabasePassword:      getEnv("DATABASE_PASSWORD", "postgres"),
		RedisHost:             getEnv("REDIS_HOST", "localhost"),
		RedisPort:             getEnv("REDIS_PORT", "6379"),
		KafkaBrokers:          getEnv("KAFKA_BROKERS", "localhost:9092"),
		JWTAccessSecret:       getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:      getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
		JWTRefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		AllowedOrigins:        getSliceEnv("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
	}

	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
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
