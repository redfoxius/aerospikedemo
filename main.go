package main

import (
	"aerospikedemo/internal/app"
	cfg "aerospikedemo/internal/app/config"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := cfg.NewConfig()

	handler := app.NewHandler(config)
	if err := handler.Process(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Total execution time %s", time.Since(config.StartTime))
}
