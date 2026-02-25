package config

import (
	"fmt"
	"os"
	"strings"
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

	// CORS
	AllowedOrigins []string

	// Cart Service URL (for checkout)
	CartServiceURL string

	// JWT
	JWTAccessSecret string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServiceName:      getEnv("SERVICE_NAME", "order-service"),
		HTTPPort:         getEnv("HTTP_PORT", "8083"),
		GRPCPort:         getEnv("GRPC_PORT", "9093"),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseName:     getEnv("DATABASE_NAME", "order"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "postgres"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		AllowedOrigins:   getSliceEnv("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		CartServiceURL:   getEnv("CART_SERVICE_URL", "http://cart:8082"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "dev-access-secret-change-me"),
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
