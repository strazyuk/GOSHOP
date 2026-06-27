package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	Port            string
	GinMode         string
	DBpath          string
	JWTSecret       string
	JWTExpiryHours  int
	uploadDir       string
	MaxUploadSizeMB int64
}

func Load() *config {
	if err := godotenv.Load(); err != nil {
		log.Println("No.env file found -using ")
	}
	return &config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "release"),
		DBPath:          getEnv("DB_PATH", "./ecommerce.db"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryHours:  getEnvInt("JWT_EXPIRY_HOURS", 72),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB: int64(getEnvInt("MAX_UPLOAD_SIZE_MB", 5)),
	}
}
func getEnv(key, defaultvalue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
}
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// strconv.Atoi converts a string to an int
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
