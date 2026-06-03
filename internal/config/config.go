package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// 服务器配置
	Port        string
	Environment string

	// 数据库配置
	DatabaseURL         string
	DatabaseMaxOpenConns int
	DatabaseMaxIdleConns int
	DatabaseMaxLifetime  time.Duration

	// 认证配置
	JWTSecret           string
	JWTAccessExpiresIn  time.Duration
	JWTRefreshExpiresIn time.Duration

	// Redis配置
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// 天气API配置
	WeatherAPIKey string

	// SOS配置
	SOSPhone string
}

func Load() *Config {
	return &Config{
		Port:                 getEnv("PORT", "8080"),
		Environment:          getEnv("ENV", "development"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gowild_db?sslmode=disable"),
		DatabaseMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DatabaseMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DatabaseMaxLifetime:  time.Duration(getEnvInt("DB_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
		JWTSecret:            getEnv("JWT_SECRET", "gowild-jwt-secret-key-2024-dev-only"),
		JWTAccessExpiresIn:   time.Duration(getEnvInt("JWT_ACCESS_EXPIRES_HOURS", 24)) * time.Hour,
		JWTRefreshExpiresIn:  time.Duration(getEnvInt("JWT_REFRESH_EXPIRES_DAYS", 30)) * 24 * time.Hour,
		RedisAddr:            getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              getEnvInt("REDIS_DB", 0),
		WeatherAPIKey:        getEnv("WEATHER_API_KEY", ""),
		SOSPhone:             getEnv("SOS_PHONE", ""),
	}
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
