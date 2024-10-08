package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	JWTSecret string
	Port      string
}

var (
	once     sync.Once
	instance *Config
)

func New() *Config {
	once.Do(func() {
		// Load .env file if it exists
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading env file")
			return
		}

		instance = &Config{
			MongoURI:  getEnv("MONGO_URI"),
			JWTSecret: getEnv("JWT_SECRET"),
			Port:      getEnv("PORT"),
		}
	})
	return instance
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func GetJWTSecret() string {
	instance := New()
	return instance.JWTSecret
}
