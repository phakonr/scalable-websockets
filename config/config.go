package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		App   App
		Redis Redis
	}

	App struct {
		ID           string
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		BodyLimit    int // bytes
	}

	Redis struct {
		Addr     string
		Username string
		Password string
	}
)

// Helper function to parse environment variables as integers.
func parseEnvInt(envValue string, defaultValue int) int {
	if envValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(envValue)
	if err != nil {
		log.Printf("error parsing environment variable as integer: %v; using default value: %d", err, defaultValue)
		return defaultValue
	}
	return value
}

func LoadConfig() Config {
	return Config{
		App: App{
			ID:           os.Getenv("APP_ID"),
			Port:         os.Getenv("APP_PORT"),
			ReadTimeout:  time.Duration(parseEnvInt(os.Getenv("APP_READ_TIMEOUT"), 60)) * time.Second,
			WriteTimeout: time.Duration(parseEnvInt(os.Getenv("APP_WRITE_TIMEOUT"), 60)) * time.Second,
			BodyLimit:    parseEnvInt(os.Getenv("APP_BODY_LIMIT"), 10490000),
		},
		Redis: Redis{
			Addr:     os.Getenv("REDIS_ADDR"),
			Username: os.Getenv("REIDS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}
