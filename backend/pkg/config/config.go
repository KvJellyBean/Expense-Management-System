package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret      string
	JWTExpiryHours int

	ServerPort string

	PaymentAPIURL string

	WorkerPoolSize   int
	WorkerMaxRetries int
}

func Load() *Config {
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	workerPoolSize, _ := strconv.Atoi(getEnv("WORKER_POOL_SIZE", "5"))
	workerMaxRetries, _ := strconv.Atoi(getEnv("WORKER_MAX_RETRIES", "3"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "expense_user"),
		DBPassword: getEnv("DB_PASSWORD", "expense_pass"),
		DBName:     getEnv("DB_NAME", "expense_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours: jwtExpiry,

		ServerPort: getEnv("SERVER_PORT", "8080"),

		PaymentAPIURL: getEnv("PAYMENT_API_URL", "https://1620e98f-7759-431c-a2aa-f449d591150b.mock.pstmn.io"),

		WorkerPoolSize:   workerPoolSize,
		WorkerMaxRetries: workerMaxRetries,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
