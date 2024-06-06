package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DatabaseConfig
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "mydatabase"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Port: getEnv("PORT", "8080"),
	}
}

func LoadConfigTest() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST_TEST", "localhost"),
			Port:     getEnv("DB_PORT_TEST", "5432"),
			User:     getEnv("DB_USER_TEST", "postgres"),
			Password: getEnv("DB_PASSWORD_TEST", ""),
			DBName:   getEnv("DB_NAME_TEST", "mydatabase"),
			SSLMode:  getEnv("DB_SSLMODE_TEST", "disable"),
		},
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
