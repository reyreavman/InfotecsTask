package database

import (
	"os"
)

// Структура, хранящая в себе данные, необходимые для подключения к БД
type Config struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
}

// Функция для загрузки конфига, в текущей реализации все параметры собираются из переменных окружения
func LoadConfig() Config {
	return Config{
		Host:         os.Getenv("POSTGRES_HOST"),
		Port:         os.Getenv("POSTGRES_PORT"),
		User:         os.Getenv("POSTGRES_USER"),
		Password:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:       os.Getenv("POSTGRES_DB_NAME"),
		SSLMode:      os.Getenv("POSTGRES_SSL_MODE"),
	}
}
