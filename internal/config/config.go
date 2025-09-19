package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port         string
	MongoURI     string
	DatabaseName string
	Environment  string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	config := &Config{
		Port:         getEnvWithDefault("PORT", "3000"),
		MongoURI:     getEnvWithDefault("MONGO_URI", "MONGO_URI=mongodb+srv://root:bfw9sJjTXxfleyu7@localdevincidents.voqqf7r.mongodb.net/localdevincidents?retryWrites=true&w=majority&ssl=true&sslInsecure=true"),
		DatabaseName: getEnvWithDefault("DATABASE_NAME", "localdevincidents"),
		Environment:  getEnvWithDefault("ENVIRONMENT", "development"),
	}

	// Log loaded configuration (excluding sensitive data)
	log.Printf("Configuration loaded:")
	log.Printf("- Port: %s", config.Port)
	log.Printf("- Database Name: %s", config.DatabaseName)
	log.Printf("- Environment: %s", config.Environment)
	log.Printf("- MongoDB URI: %s", maskURI(config.MongoURI))

	return config
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// maskURI masks sensitive information in URI for logging
func maskURI(uri string) string {
	if len(uri) > 20 {
		return uri[:20] + "..."
	}
	return uri
}
