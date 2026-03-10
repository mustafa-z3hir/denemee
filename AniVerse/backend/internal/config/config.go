package config

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "aniverse"),
		DBPassword: getEnv("DB_PASSWORD", "aniverse123"),
		DBName:     getEnv("DB_NAME", "aniverse"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		JWTSecret:  getEnv("JWT_SECRET", "aniverse-secret-key-2024"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	// SQLite kullan (dosya tabanlı, kurulum gerektirmez)
	return gorm.Open(sqlite.Open("aniverse.db"), &gorm.Config{})
}

func ConnectRedis(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})
}
