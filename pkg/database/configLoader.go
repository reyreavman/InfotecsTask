package database

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

func LoadConfig() Config {
	maxOpenConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if err != nil {
		log.Fatalf("Failed to get env: %s", err.Error())	
	}
	maxIdleConns, err := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if err != nil {
		log.Fatalf("Failed to get env: %s", err.Error())
	}

	return Config{
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:       os.Getenv("POSTGRES_DB_NAME"),
		SSLMode:      os.Getenv("POSTGRES_SSL_MODE"),
		MaxOpenConns: maxOpenConns,
		MaxIdleConns: maxIdleConns,
	}
}
