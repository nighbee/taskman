package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DatabaseURL            string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBSSLMode              string
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string

	// JWT
	JWTSecret string
	JWTExpiry string

	// Server
	Port    string
	Host    string
	GinMode string

	// CORS
	CORSAllowedOrigins []string

	// Logging
	LogLevel string

	// Migrations
	RunMigrations bool

	// Email
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		// Database configuration - prefer individual params over DATABASE_URL
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", ""),
		DBName:                 getEnv("DB_NAME", "taskman"),
		DBSSLMode:              getEnv("DB_SSL_MODE", "disable"),
		SupabaseURL:            getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		JWTSecret:              getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
		JWTExpiry:              getEnv("JWT_EXPIRY", "24h"),
		Port:                   getEnv("SERVER_PORT", getEnv("PORT", "8080")),
		Host:                   getEnv("SERVER_HOST", "localhost"),
		GinMode:                getEnv("GIN_MODE", "debug"),
		CORSAllowedOrigins:     strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"), ","),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		RunMigrations:          getEnvAsBool("RUN_MIGRATIONS", false),
		SMTPHost:               getEnv("SMTP_HOST", ""),
		SMTPPort:               getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:           getEnv("SMTP_USERNAME", ""),
		SMTPPassword:           getEnv("SMTP_PASSWORD", ""),
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
