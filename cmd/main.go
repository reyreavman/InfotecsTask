package main

import (
	"infotecstechtask/internal/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("./deployments/.env"); err != nil {
		log.Fatalf("%s", err.Error())
	}

	port := os.Getenv("APP_PORT")

	app := server.NewApp()

	if err := app.Run(port); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
