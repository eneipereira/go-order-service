package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente do sistema.")
	}

	dbUser := getEnv("POSTGRES_USER", "admin")
	dbPassword := getEnv("POSTGRES_PASSWORD", "admin")
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbName := getEnv("POSTGRES_DB", "orders_service")
	dbSSLMode := getEnv("POSTGRES_SSLMODE", "disable")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

	return Config{
		DatabaseURL: dbURL,
		ServerPort:  fmt.Sprintf(":%s", getEnv("PORT", "8080")),
	}
}